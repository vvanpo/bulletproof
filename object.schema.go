package bp

var SCHEMA string

func init() {
	SCHEMA = `PRAGMA foreign_keys=ON;
	BEGIN TRANSACTION;`
	// The public keys of trusted nodes
	+ `CREATE TABLE trusted (key TEXT PRIMARY KEY NOT NULL, alias TEXT);`
	// The global and local tables list the user-specified paths to track
	+ `CREATE TABLE global (path TEXT UNIQUE, flags INTEGER, override TEXT UNIQUE, CHECK (CASE WHEN path ISNULL THEN override NOTNULL END));`
	+ `CREATE TABLE local (path TEXT UNIQUE, flags INTEGER, override TEXT UNIQUE, CHECK (CASE WHEN path ISNULL THEN override NOTNULL END));`
	// The file object pointed to by the global and local tables.
	+ `CREATE TABLE file (uuid TEXT PRIMARY KEY NOT NULL, size INTEGER, mode INTEGER NOT NULL, modtime INTEGER NOT NULL, hash TEXT);`
	// The event table records update events to files.  Events that have been
	// synced across all nodes can be cleaned, and multiple updates to the same
	// file can be merged.
	+ `CREATE TABLE event (uuid TEXT PRIMARY KEY NOT NULL, time INTEGER NOT NULL, file TEXT REFERENCES file (id) NOT NULL);`

	+ `COMMIT;`
}
