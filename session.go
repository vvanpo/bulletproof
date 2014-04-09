
package bp

import (
	"code.google.com/p/go.exp/fsnotify"
	"fmt"
)

// Per-instance session object
type Session struct {
	// Absolute pathname to root
	root string
	// Absolute pathname mapping to fsnotify watcher objects
	watcher fsnotify.Watcher
}

func NewSession(root string) *Session {
	s := new(Session)
	s.root = root
	s.watcher, err = fsnotify.NewWatcher()
	s.updateConf()
	return s
}

func (s *Session) setWatcher(path string) error {
	_, exists := s.watchers[path]
	if exists {
		return fmt.Errorf("Watcher for file %s already exists", path)
	}
	w, err := fsnotify.NewWatcher()
	if err != nil { return err }
	err = w.Watch
	s.watchers[path] = w
	return nil
}
