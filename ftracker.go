
package bphome

import (
    "os"
    "strings"
    "time"
    "code.google.com/p/go.exp/fsnotify"
)

// Directory and file name components. Relative paths require "." or ".." as the first element.
// e.g. {"home", "victor", "my-tracked-dir", "my-file"}
//   or {"..", "my-tracked-dir", "my-file"}
type path []string

func (p path) String() string {
    return strings.Join(p, os.PathSeparator)
}

// tracked is a per-user structure of all tracked files and directories for that
// user
type tracked struct {
    // Root directory
    root path
}

// A representation of an inode
type node struct {
    size []int64        // Size in bytes
    mode os.FileMode
    modtime time.time
    seLinux string      // SELinux context, if one exists
}

// Represents a file, i.e. a link to an inode
type link *node

type symlink struct {
    // The symlink's attributes
    this link
    // The linked file's path, either relative or absolute
    pathname path
}

// The currently linked file of the symlink
func (s *symlink) link() link {
    p := s.path.abs()
    return nil
}

func (f *file) name() {
    os.filepath

func (d *dir) addFile() {
    return
}

func searchAll(root dir) {

    return
}
