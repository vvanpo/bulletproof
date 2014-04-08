
// Schema for object store:
// CREATE TABLE trusted (key TEXT PRIMARY KEY NOT NULL, alias TEXT);
// CREATE TABLE object (uuid TEXT PRIMARY KEY NOT NULL, hash TEXT, size INTEGER, modtime INTEGER);
// CREATE TABLE global (path TEXT UNIQUE, object TEXT REFERENCES object (uuid), flags INTEGER, override TEXT UNIQUE,
//					CHECK (CASE WHEN path ISNULL THEN override NOTNULL END),
//					CHECK (override NOT IN (SELECT path FROM global)));
// CREATE TABLE local (path TEXT UNIQUE, object TEXT REFERENCES object (uuid), flags INTEGER, override TEXT UNIQUE,
//					CHECK (CASE WHEN path ISNULL THEN override NOTNULL END),
//					CHECK (override NOT IN (SELECT path FROM local)));
// CREATE VIEW files AS (SELECT path FROM global EXCEPT SELECT override AS path FROM local)
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

