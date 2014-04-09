package bp

import (
	//	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	//	"crypto/md5"
	"path/filepath"
	//	"os"
	//	"time"
)

type db struct {
}

// newDatabase creates a new, empty object store, backing up the older if
// necessary
func (s *Session) newDatabase() error {
	path := filepath.Join(s.root, ".bp/object.db")
	c, err := sqlite3.Open(path)
	if err != nil {
		return err
	}
	defer c.Close()
	return c.Exec(SCHEMA)
}

// Object flags
const (
	encrypt int = 1 << iota
	dataOnly
	metadataOnly
	follow
)

// getObject grabs the values stored in the database for a specified path
func (s *Session) getObject(path string) {
	return
}

// verifyObject returns true if a the stored object data is consistent with the
// the current values for the object
func (s *Session) verifyObject(path string) (bool, error) {
	s.getObject(path)
	return true, nil
}
