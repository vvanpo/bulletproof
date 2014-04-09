package bp

var SCHEMA string

func init() {
	SCHEMA = `PRAGMA foreign_keys=ON;
	BEGIN TRANSACTION;
	CREATE TABLE trusted (key TEXT PRIMARY KEY NOT NULL, alias TEXT);
	CREATE TABLE object (uuid TEXT PRIMARY KEY NOT NULL, size INTEGER NOT NULL,
			mode INTEGER NOT NULL, modtime INTEGER NOT NULL,
			hash TEXT NOT NULL);
	CREATE TABLE global (path TEXT UNIQUE, object TEXT REFERENCES object (uuid),
			flags INTEGER, override TEXT UNIQUE,
			CHECK (CASE WHEN path ISNULL THEN override NOTNULL END));
	CREATE TABLE local (path TEXT UNIQUE, object TEXT REFERENCES object (uuid),
			flags INTEGER, override TEXT UNIQUE,
			CHECK (CASE WHEN path ISNULL THEN override NOTNULL END));
	CREATE VIEW files AS SELECT path, object, flags, override FROM global
			WHERE(path NOT IN (SELECT override FROM local))
			UNION SELECT path, object, flags, override FROM local;
	COMMIT;`
}
