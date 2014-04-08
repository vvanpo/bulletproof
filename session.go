
package bp

import (
	"code.google.com/p/go.exp/fsnotify"
//	"code.google.com/p/go-uuid/uuid"
)

// Per-instance session object
type Session struct {
	// Absolute pathname to root
	root string
	// Configuration structure
	config
	// Absolute pathname mapping to fsnotify watcher objects
	watchers map[string]fsnotify.Watcher
}

func NewSession(root string) *Session {
	s := new(Session)
	s.root = root
	s.updateCfg()
	return s
}
