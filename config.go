
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

