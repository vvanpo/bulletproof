
// Schema for object store:
// CREATE TABLE trusted (key TEXT PRIMARY KEY NOT NULL, alias TEXT);
// CREATE TABLE object (uuid TEXT PRIMARY KEY NOT NULL, hash TEXT, modtime INTEGER);
// CREATE TABLE global (path TEXT PRIMARY KEY NOT NULL, object TEXT REFERENCES object (uuid), flags INTEGER);
// CREATE TABLE local (path TEXT UNIQUE, object TEXT REFERENCES object (uuid), flags INTEGER, override TEXT REFERENCES global (path) UNIQUE, CHECK (CASE WHEN path ISNULL THEN override NOTNULL END));
package bp

import (
	"path/filepath"
	"code.google.com/p/go-sqlite/go1/sqlite3"
)

// Object flags
const (
	recursive int = 1 << iota
	follow
	dataOnly
	metadataOnly
	encrypt
)

// Unique file/dir object
type object struct {
}

type config struct {
}

func (s *Session) updateCfg() error {
	path := filepath.Join(s.root, ".bp/global")
	c, err := sqlite3.Open(path)
	defer c.Close()
	if err != nil { return err }
	
	return nil
}

