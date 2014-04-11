package bp

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"crypto/md5"
	"fmt"
	"os"
	"io/ioutil"
	"time"
)

// object represents a unique file, directory or symlink
type object struct {
	mode os.FileMode
	modTime time.Time
	size int64	// Size in bytes, only for regular files
	hash string	// Only for regular files
}

// Object flags
const (
	encrypt int = 1 << iota
	dataOnly
	metadataOnly
	follow
)

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

func (s *Session) dbConn() (*sqlite3.Conn, error) {
	c, err := sqlite3.Open(s.absolutePath(".bp/object.db"))
	if err != nil { return nil, err }
	return c, c.Exec("PRAGMA foreign_keys = ON;")
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

// getObject grabs the object values stored in the database for a specified path
func (s *Session) getObject(path string) (o object, err error) {
	c, err := s.dbConn()
	if err != nil { return }
	defer c.Close()
	if err != nil { return }
	return
}

// addObject adds a path to the global table
func (s *Session) addObject(path string, flags int) error {
	o, err := s.getFile(path)
	if err != nil { return err }
	if o.mode & (os.ModeType &^ os.ModeDir &^ os.ModeSymlink) != 0 {
		return fmt.Errorf("Cannot add irregular file '%s'.", path)
	}
	c, err := s.dbConn()
	if err != nil { return err }
	defer c.Close()
	uuid := uuid.NewUUID().String()
	err = c.Exec("INSERT INTO global VALUES (?, ?, ?, ?, ?, ?, ?);", uuid, path, flags, int(o.mode), o.modTime.UnixNano(), o.size, o.hash)
	return err
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

