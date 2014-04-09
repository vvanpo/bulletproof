package bp

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"crypto/md5"
	"fmt"
	"crypto/sha1"
	"os"
	"io"
	"io/ioutil"
	"time"
)

/*
// Objects database
type db struct {
	// Convenience field, is always s.root + ".bp/object.db"
	path string
}*/

func (s *Session) dbConn() (*sqlite3.Conn, error) {
	return sqlite3.Open(s.absolutePath(".bp/object.db"))
}

// createDatabase creates a new, empty object store
func (s *Session) createDatabase() error {
	c, err := s.dbConn()
	if err != nil { return err }
	defer c.Close()
	err = c.Exec(SCHEMA)
	if err != nil { return err }
	return nil
}

// object represents a unique file object (i.e. an inode)
type object struct {
	size int64	// In bytes, only for regular files
	mode os.FileMode
	modTime time.Time
	hash string	// Only for regular files
}

// Object flags
const (
	encrypt int = 1 << iota
	dataOnly
	metadataOnly
	follow
)

// getObject grabs the values stored in the database for a specified path
func (s *Session) getObject(path string) (o object, err error) {
	c, err := s.dbConn()
	if err != nil { return }
	defer c.Close()
	if err != nil { return }
	return
}

// getFile grabs the object values from the live file
func (s *Session) getFile(path string) (o object, err error) {
	fi, err := os.Lstat(s.absolutePath(path))
	if err != nil { return }
	o.mode = fi.Mode()
	o.modTime = fi.ModTime().UTC()
	if o.mode.IsRegular() {
		o.size = fi.Size()
		var buf []byte
		buf, err = ioutil.ReadFile(s.absolutePath(path))
		hash := md5.Sum(buf)
		o.hash = fmt.Sprintf("%x", hash[:])
	}
	return
}

// addPath grabs the file info from the specified path and stores it in the
// database.  The local parameter determines which table it is to be stored in.
func (s *Session) addPath(path string, local bool) error {
	o, err := s.getFile(path)
	if err != nil { return err }
	c, err := s.dbConn()
	if err != nil { return err }
	defer c.Close()
	// Ensure path doesn't already exist


	var uuid string
	pathhash := sha1.Sum(path)
	q, err := c.Query("SELECT path FROM files AS f INNER JOIN objects AS o ON f.object == o.uuid WHERE o.mode == ? && o.modtime == ?", int(o.mode), o.modTime.Unix())
	for ; err == nil; err = q.Next() {
		var p string
		q.Scan(&p)
		if os.Samefile(
	}
	if err == io.EOF {
		// Need to add a new object
		uuid = uuid.NewUUID().String()
		if o.mode.IsRegular() {
			err = c.Exec("INSERT INTO object VALUES (?, ?, ?, ?, ?, ?)",
					uuid, pathhash, o.size, int(o.mode), o.modTime.Unix(), o.hash)
		} else {
			err = c.Exec("INSERT INTO object (uuid, id, mode, modtime) VALUES (?, ?, ?, ?)",
					uuid, pathhash, int(o.mode), o.modTime.Unix())
		}
		return err
	} else if err != nil {
		return err
	}
}

// verifyObject returns true if a the stored object data is consistent with the
// the current values for the object
func (s *Session) verifyObject(path string) (bool, error) {
	o, err := s.getObject(path)
	if err != nil { return false, err }
	f, err := s.getFile(path)
	if err != nil { return false, err }
	if o != f {
		return false, nil
	}
	return true, nil
}

