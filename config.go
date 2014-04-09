
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
	"fmt"
	"path/filepath"
	"code.google.com/p/go-sqlite/go1/sqlite3"
)

func (s *Session) updateConf() error {
	db := filepath.Join(s.root, ".bp/objects.db")
	c, err := sqlite3.Open(db)
	defer c.Close()
	if err != nil { return err }
	sql := "SELECT path FROM global"
	row := make(sqlite3.RowMap)
	for i, err := c.Query(sql); err == nil; err = i.Next() {
		var path string
		i.Scan(path, row)
		fmt.Println(path, row)
	}
	return nil
}

