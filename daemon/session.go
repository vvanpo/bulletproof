package main

import (
	"github.com/vvanpo/bulletproof/object"
	"code.google.com/p/go.exp/fsnotify"
	"fmt"
	"os"
	"path/filepath"
	"log"
)

// Per-instance session object
type session struct {
	// Absolute pathname to root
	root string
	// The object store instance
	store object.ObjectStore
	// File watcher instance
	watcher *fsnotify.Watcher
	// Map of all watched paths to events
	event map[string]*fsnotify.FileEvent
}

func newSession(root string) *session {
	s := new(session)
	s.root = root
	os.Chdir(root)
	var err error
	s.store, err = object.CreateSqlite()
	if err != nil {
		log.Fatalf("Failed to initialize object store:\n%s", err)
	}
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to initialize inotify instance:\n%s", err)
	}
	return s
}

func (s *session) start() error {
	err := filepath.Walk("", walkFn)
	return err
}

func (s *session) addWatch(path string) error {
	err := s.watcher.Watch(path)
	if err != nil {
		return fmt.Errorf("fsnotify: %s", err)
	}
	return nil
}

// walkFn visits every file, dir, and symlink in the s.root tree
func walkFn(path string, info os.FileInfo, err error) error {
	fmt.Println(path)
	return nil
}
