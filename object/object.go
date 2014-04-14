// Storage-access API
package object

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"code.google.com/p/go-uuid/uuid"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
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
	StatObject(path string) (Object, error)
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
	return c, c.Exec("PRAGMA foreign_keys = ON;")
}

// CreateSqlite uses the provided directory to create a new, empty object store
// ready for AddObject calls
func CreateSqlite() (*Sqlite, error) {
	db := new(Sqlite)
	db.file = ".bp/object.db"
	_, err := os.Stat(db.file)
	if err != nil {
		os.MkdirAll(filepath.Split(db.file), 0755)
	} else if VerifySchema() {
		return fmt.Errorf("Database already exists.")
	}
	c, err := db.conn()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	err = c.Exec(schema)
	if err != nil {
		return nil, err
	}
	return db, nil
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

// AddObject adds a path to be tracked to the store
func (s *Sqlite) AddObject(path string, flags int, o Object) error {
	if o.mode&(os.ModeType&^os.ModeDir&^os.ModeSymlink) != 0 {
		return fmt.Errorf("Cannot add irregular file '%s'.", path)
	}
	c, err := s.conn()
	if err != nil {
		return err
	}
	defer c.Close()
	uuid := uuid.NewUUID().String()
	err = c.Exec("INSERT INTO global VALUES (?, ?, ?, ?, ?, ?, ?);", uuid, path, flags, int(o.mode), o.modTime.UnixNano(), o.size, o.hash)
	return err
}

// ViewObject queries the store for the given path's object values
func (s *Sqlite) ViewObject(path string) (o Object, err error) {
	c, err := s.conn()
	if err != nil {
		return
	}
	defer c.Close()
	q, err := c.Query("SELECT mode, modtime, size, hash FROM global WHERE path == ?", path)
	if err == nil {
		var mode, modTime int64
		q.Scan(&mode, &modTime, &o.size, &o.hash)
		o.mode = os.FileMode(mode)
		o.modTime = time.Unix(0, modTime)
	}
	return
}

func (s *Sqlite) RemoveObject(path string) error {
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
