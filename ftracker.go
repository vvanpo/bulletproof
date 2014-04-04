package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"crypto/md5"
	"gopkg.in/yaml.v1"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

// Directory and file name components
//   e.g. {"home", "victor", "my-tracked-dir", "my-file"}
//     or {"..", "my-tracked-dir", "my-file"}
type Path []string

// String representation of Path
func (p Path) String() string {
	return path.Join(p...)
}

// A representation of an inode, i.e. a unique file object
type node struct {
	size    []int64 // Size in bytes
	mode    os.FileMode
	modtime time.Time
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
	symlinks []map[link]Path
	// List of directories nested within this directory
	dirs []dir
}

// Per-instance file tree object
type Instance struct {
	// Root directory
	root *dir
	// Absolute pathname to root
	path Path
	// Absolute pathname mapping to fsnotify watcher objects
	watchers map[string]fsnotify.Watcher
}

// NewInstance sets up the data structure for the given path by reading the configuration files
func NewInstance(p Path) *Instance {
	i := new(Instance)
	i.root = new(dir)
	b, err := yaml.Marshal(i)
	print(err)
	os.Create("example.yaml")

	return i
}

func (i *Instance) AddPath(p Path) error {
	return nil
}

func main() {
	root := strings.Split(os.Args[1], "/")
	i := NewInstance(root)
	println(i)
}
