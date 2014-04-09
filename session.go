package bp

import (
	"code.google.com/p/go.exp/fsnotify"
	"fmt"
	"log"
	"path/filepath"
)

// Per-instance session object
type Session struct {
	// Absolute pathname to root
	root string
	// Database
	*db
	// File watcher instance
	watcher *fsnotify.Watcher
	// List of watched paths
	watch []string
}

func NewSession(root string) *Session {
	s := new(Session)
	s.root = root
	var err error
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to initialize inotify instance:\n%s", err)
	}
	return s
}

func (s *Session) addWatch(path string) error {
	err := s.watcher.Watch(filepath.Join(s.root, path))
	if err != nil {
		return fmt.Errorf("fsnotify: %s", err)
	}
	return nil
}
