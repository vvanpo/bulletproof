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

func (n *node) isSymlink() bool {
	if n.mode|os.ModeSymlink != 0 {
		return true
	}
	return false
}

// A file, i.e. a hard-link to a node
type file struct {
	*node
	name string
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
	// Absolute pathname mapping to fsnotify watcher objects
	watchers map[string]fsnotify.Watcher
}

func NewSession(pathname string) *Session {
	return nil
}

