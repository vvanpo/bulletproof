package object

import (
	"fmt"
	"os"
)

var schema string

func init() {
	schema = `BEGIN TRANSACTION;` +
		// The public keys of trusted nodes
		`CREATE TABLE trusted (key TEXT PRIMARY KEY NOT NULL, alias TEXT);` +
		// The global table lists the user-specified paths to track
		`CREATE TABLE global (uuid TEXT PRIMARY KEY NOT NULL, path TEXT UNIQUE NOT NULL, flags INTEGER, mode INTEGER NOT NULL, modtime INTEGER NOT NULL, size INTEGER, hash TEXT);` +
		// The event table records update events to files.  Events that have been
		// synced across all nodes can be cleaned, and multiple updates to the same
		// file can be merged.
		`CREATE TABLE event (uuid TEXT PRIMARY KEY NOT NULL, time INTEGER NOT NULL, file TEXT REFERENCES global (uuid) NOT NULL);` +

		fmt.Sprintf("CREATE VIEW folder AS SELECT * FROM global WHERE (mode & %d) != 0;", os.ModeDir) +
		fmt.Sprintf("CREATE VIEW symlink AS SELECT * FROM global WHERE (mode & %d) != 0;", os.ModeSymlink) +
		fmt.Sprintf("CREATE VIEW file AS SELECT * FROM global WHERE (mode & %d) == 0;", os.ModeType) +
		`COMMIT;`
}
