package bp

var SCHEMA string

func init() {
	SCHEMA = `PRAGMA foreign_keys=ON;
	BEGIN TRANSACTION;
	CREATE TABLE trusted (key TEXT PRIMARY KEY NOT NULL, alias TEXT);`

	+ `CREATE TABLE object (uuid TEXT PRIMARY KEY NOT NULL, id TEXT REFERENCES object_id (id) NOT NULL, size INTEGER, mode INTEGER NOT NULL, modtime INTEGER NOT NULL, hash TEXT);`
	// Object ids are created by hashing the concatentation of all (hashed)
	// paths.  Ids get updated whenever links are added or removed to the
	// object.  Objects with at least one link from a 'global' path, use only
	// the global paths to compute the id.  Only when there are only local
	// paths pointing to the object are they computed using the local table.
	+ `CREATE TABLE object_id (id TEXT PRIMARY KEY NOT NULL, object TEXT REFERENCES object (uuid) NOT NULL, last_id TEXT REFERENCES object_id (id));`
	// The events table records modifications, most importantly ones that
	// result in relabeled object ids.  Events that have been synced with all
	// nodes can be cleaned, and events that are between and syncs that don't
	// involve id relabels can be combined.
	+ `CREATE TABLE event (uuid TEXT PRIMARY KEY NOT NULL, time INTEGER NOT NULL, object TEXT REFERENCES object_id (id) NOT NULL);`

	+ `CREATE TABLE global (path TEXT UNIQUE, object TEXT REFERENCES object (uuid), flags INTEGER, override TEXT UNIQUE, CHECK (CASE WHEN path ISNULL THEN override NOTNULL END));`

	+ `CREATE TABLE local (path TEXT UNIQUE, object TEXT REFERENCES object (uuid), flags INTEGER, override TEXT UNIQUE, CHECK (CASE WHEN path ISNULL THEN override NOTNULL END));`

	+ `CREATE VIEW files AS SELECT path, object, flags, override FROM global
			WHERE (path NOT IN (SELECT override FROM local))
			UNION SELECT path, object, flags, override FROM local;
	COMMIT;`
}
