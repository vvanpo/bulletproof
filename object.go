package bp

import (
	//	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"crypto/md5"
	"path/filepath"
	"os"
	"io/ioutil"
	"time"
)

// Objects database
type db struct {
	// Convenience field, is always s.root + ".bp/object.db"
	path string
}

// createDatabase creates a new, empty object store
func (s *Session) createDatabase() error {
	path := filepath.Join(s.root, ".bp/object.db")
	c, err := sqlite3.Open(path)
	if err != nil {
		return err
	}
	defer c.Close()
	err = c.Exec(SCHEMA)
	if err != nil { return err }
	s.db = new(db)
	s.db.path = path
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
	return
}

// getFile grabs the object values from the live file
func (s *Session) getFile(path string) (o object, err error) {
	fi, err := os.Lstat(s.absolutePath(path))
	if err != nil { return }
	o.mode = fi.Mode()
	o.modTime = fi.ModTime()
	if o.mode.IsRegular() {
		o.size = fi.Size()
		var buf []byte
		buf, err = ioutil.ReadFile(s.absolutePath(path))
		hash := md5.Sum(buf)
		o.hash = string(hash[:])
	}
	return
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

