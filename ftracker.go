package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// A representation of an inode, i.e. a unique file object
type node struct {
	size    int64 // Size in bytes
	mode    os.FileMode
	modTime time.Time
	seLinux string         // SELinux context, if one exists
	hash    [md5.Size]byte // MD5 hash of file data
	links   []*link        // Back-references to files.  This is ignored if node is a directory
}

// Represents a file, i.e. a link to an inode
type link struct {
	// Base name
	name string
	*node
}

func (l link) isSymlink() bool {
	if l.node.mode|os.ModeSymlink != 0 {
		return true
	}
	return false
}

type dir struct {
	// A directory is also just a link, but most systems don't support multi-linking directories and neither will we
	name string
	*node
	// List of files
	files []link
	// List of symlinks within this directory, mapping the symlink files themselves to a pathname
	symlinks map[link]string
	// List of directories nested within this directory
	dirs []dir
}

// Per-instance file tree object
type Session struct {
	// Absolute pathname to root
	pathname string
	// Root directory
	root *dir
	// Absolute pathname mapping to fsnotify watcher objects
	watchers map[string]fsnotify.Watcher
}

// NewSession sets up the data structure for the given pathname by reading the configuration files in that path
func NewSession(pathname string) *Session {
	return nil
}

