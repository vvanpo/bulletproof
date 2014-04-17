// Storage-access API
package object

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"code.google.com/p/go-uuid/uuid"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"path/filepath"
)

// Object represents a unique file, directory or symlink
type Object struct {
	mode    os.FileMode
	modTime time.Time
	size    int64  // Size in bytes, only for regular files
	hash    string // Only for regular files
}

// Test if an Object is equal to another
func (o *Object) Equal(p Object) bool {
	if o.mode != p.mode || o.size != p.size || o.hash != p.hash {
		return false
	}
	return o.modTime.Equal(p.modTime)
}

// Path flags for defining how an object gets synced
const (
	Encrypt int = 1 << iota
	DataOnly
	MetadataOnly
	Follow
)

// This interface is used to define the access methods to a particular object
// store
type ObjectStore interface {
	VerifySchema() (bool, error)
	StatObject(path string) (Object, error)
	IsObject(path string) (bool, error)
	AddObject(path string, flags int, o Object) error
	ViewObject(path string) (Object, error)
	RemoveObject(path string) error
	VerifyObject(path string) (bool, error)
}

// Our implementation of ObjectStore uses sqlite as a back-end
type Sqlite struct{
	file string
}

func (s *Sqlite) conn() (c *sqlite3.Conn, err error) {
	c, err = sqlite3.Open(s.file)
	if err != nil { return }
	return c, c.Exec("PRAGMA foreign_keys = ON")
}

// CreateSqlite uses the provided directory to create a new object store ready
// for AddObject calls
func CreateSqlite() (*Sqlite, error) {
	s := new(Sqlite)
	s.file = ".bp/object.db"
	_, err := os.Stat(s.file)
	if err != nil {
		os.MkdirAll(filepath.Dir(s.file), 0755)
	} else if v, err := s.VerifySchema(); v {
		if err != nil { return nil, err }
		return s, nil
	}
	c, err := s.conn()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	err = c.Exec(schema)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Verifies that the database exists and conforms to the schema outlined in
// object.schema.go
func (s *Sqlite) VerifySchema() (bool, error) {
	c, err := s.conn()
	if err != nil { return false, err }
	defer c.Close()
	q, e := c.Query(`SELECT sql FROM sqlite_master WHERE type == 'table' OR type == 'view';`)
	for ; e == nil; e = q.Next() {
		var sql string
		err := q.Scan(&sql)
		if err != nil { return false, err }
		match := strings.Index(schema, sql)
		if match < 0 { return false, err }
	}
	if e != io.EOF { return false, e }
	return true, nil
}

// StatObject queries the current object values from the filesystem
func (s *Sqlite) StatObject(path string) (o Object, err error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return
	}
	o.mode = fi.Mode()
	o.modTime = fi.ModTime().UTC()
	if o.mode.IsRegular() {
		o.size = fi.Size()
		var buf []byte
		buf, err = ioutil.ReadFile(path)
		hash := md5.Sum(buf)
		o.hash = fmt.Sprintf("%x", hash[:])
	}
	return
}

// IsObject returns whether the specified path is in the object store
func (s *Sqlite) IsObject(path string) (b bool, err error) {
	c, err := s.conn()
	if err != nil { return }
	defer c.Close()
	q, err := c.Query("SELECT count(*) FROM global WHERE path == ?", path)
	if err == nil {
		var count int
		q.Scan(&count)
		if count > 0 {
			b = true
		}
		q.Close()
	} else if err == io.EOF {
		err = fmt.Errorf("Failed to retrieve object '%s'.", path)
	}
	return

}

// AddObject adds a path to be tracked to the store
func (s *Sqlite) AddObject(path string, flags int, o Object) error {
	if o.mode&(os.ModeType&^os.ModeDir&^os.ModeSymlink) != 0 {
		return fmt.Errorf("Cannot add irregular file '%s'.", path)
	}
	c, err := s.conn()
	if err != nil { return err }
	defer c.Close()
	uuid := uuid.NewUUID().String()
	err = c.Exec("INSERT INTO global VALUES (?, ?, ?, ?, ?, ?, ?)", uuid, path, flags, int(o.mode), o.modTime.UnixNano(), o.size, o.hash)
	if err != nil { return err }
	return err
}

// ViewObject queries the store for the given path's object values
func (s *Sqlite) ViewObject(path string) (o Object, err error) {
	c, err := s.conn()
	if err != nil { return }
	defer c.Close()
	q, err := c.Query("SELECT mode, modtime, size, hash FROM global WHERE path == ?", path)
	if err == nil {
		var mode, modTime int64
		q.Scan(&mode, &modTime, &o.size, &o.hash)
		o.mode = os.FileMode(mode)
		o.modTime = time.Unix(0, modTime)
		q.Close()
	} else if err == io.EOF {
		err = fmt.Errorf("Failed to retrieve object '%s'.", path)
	}
	return
}

func (s *Sqlite) RemoveObject(path string) error {
	c, err := s.conn()
	if err != nil { return err }
	defer c.Close()
	err = c.Exec("DELETE FROM global WHERE path == ?", path)
	if err != nil { return err }
	if c.RowsAffected() == 0 {
		return fmt.Errorf("Object '%s' failed to be removed from table.", path)
	}
	return nil
}

// VerifyObject returns true if a the stored object data is consistent with the
// the filesystem
func (s *Sqlite) VerifyObject(path string) (bool, error) {
	o, err := s.ViewObject(path)
	if err != nil {
		return false, err
	}
	f, err := s.StatObject(path)
	if err != nil {
		return false, err
	}
	return o.Equal(f), nil
}

// VerifyAllObjects returns a map of all paths in the database, with the value
// being true or false depending on whether the stored object is consistent
// with the live object
func (s *Sqlite) VerifyAllObjects() (p map[string]bool, err error) {
	c, err := s.conn()
	if err != nil { return err }
	defer c.Close()
	q, e := c.Query("SELECT path FROM global")
	for ; e == nil; e = q.Next() {
		var path string
		err = t.Scan(&path)
		if err != nil { return nil, err }
		if v, err := s.VerifyObject(path); v {
			p[path] = true
		} else {
			p[path] = false
		}
	}
	if e != io.EOF { return nil, e }
}
