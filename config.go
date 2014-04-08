
package bp

import (
	"path/filepath"
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
	return nil
}

