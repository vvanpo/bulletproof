
package bphome

import (
    "os"
    "strings"
    "time"
    "code.google.com/p/go.exp/fsnotify"
)

// Directory and file name components
// Relative paths require "." or ".." as the first element
// e.g. {"home", "victor", "my-tracked-dir", "my-file"}
//   or {"..", "my-tracked-dir", "my-file"}
type path []string

// String representation of path
func (p path) String() string {
    return strings.Join(p, os.PathSeparator)
}

// A representation of an inode, i.e. a unique file object
type node struct {
    size []int64        // Size in bytes
    mode os.FileMode
    modtime time.time
    seLinux string      // SELinux context, if one exists
}

// Represents a file, i.e. a link to an inode
type link struct {
    // Base name
    name string
    *node
}

type dir struct {
    // A directory is also just a link, but most systems don't support multi-linking directories and neither will we
    name string
    *node
    files []link
    // List of symlinks within this directory, mapping the symlink files themselves to a pathname
    symlinks []map[link]path
}

var fileTree struct {
    // Root directory
    root dir
    // Absolute pathname to root
    path
    // Absolute pathnames mapping to fsnotify watcher objects
    watchers map[path]fsnotify.Watcher
}


