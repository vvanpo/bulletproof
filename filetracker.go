package filetracker

import (
	"code.google.com/p/go.exp/fsnotify"
	"crypto/md5"
	"os"
	"time"
)

// A representation of an inode, i.e. a unique file object
type node struct {
	size    int64 // Size in bytes
	mode    os.FileMode
	modTime time.Time
	seLinux string         // SELinux context, if one exists
	hash    [md5.Size]byte // MD5 hash of file data
}

// A file, i.e. a hard-link to a node
type file struct {
	name string
	*node
	symlinkPath string	// Empty if regular file
}

type dir struct {
	name string
	*node
	file []file
	// List of sub-directories
	subdir []dir
}

// Per-instance session object
type Session struct {
	// Absolute pathname to root
	pathname string
	// Root directory
	root *dir
	// Configuration structure
	config
	// Absolute pathname mapping to fsnotify watcher objects
	watchers map[string]fsnotify.Watcher
}

func NewSession(pathname string) *Session {
	c, err := GetConfig(pathname)
	s := new(Session)
	return s
}

type config struct {
	trackedFiles = []string
}

func GetConfig(pathname string) (*config, err) {
	return nil, nil
}
