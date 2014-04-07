package bp

import (
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

