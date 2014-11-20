package main

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"code.google.com/p/rsc/fuse"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"log"
	"os"
	"path"
)

/*
 * main parses the command-line in the form of:
 * ./fs <file> <mount-point>
 */
func main() {
	flag.Parse()
	f := new(FS)
	master := path.Clean(flag.Arg(1))
	f.authenticate([]byte("password"))
	if master == "." {
		log.Fatal("Missing file")
	}
	mountPoint := flag.Arg(2)
	if mountPoint == "" {
		log.Fatal("Missing mount point")
	}
	if !f.validate() {
		log.Fatal("Invalid file")
	}
	var err error
	f.File, err = os.OpenFile(master, os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	// TODO password
	if err = f.authenticate([]byte("password")); err != nil {
		log.Fatal(err)
	}
	c, err := fuse.Mount(mountPoint)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Serve(FS{}); err != nil {
		log.Fatal(err)
	}
}

type FS struct {
	*os.File // master file descriptor
	ciph cipher.AEAD
}

// validate ensures the master file is a valid filesystem
func (f FS) validate() bool {
	return true
}

func (f *FS) authenticate(pass []byte) (err error) {
	key := pbkdf2.Key(pass, nil, 4096, 32, sha256.New)
	ciph, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	f.ciph, err = cipher.NewGCM(ciph)
	if err != nil {
		return
	}
	return nil
}

func (f *FS) Root() (fuse.Node, fuse.Error) {
	b := readBlock(0)
	return nil, nil
}

/*
func (FS) Statfs(r *fuse.StatfsResponse, intr fuse.Intr) fuse.Error {
}
*/

/*
type Handle struct{}

func (Handle) Flush() fuse.Error {
}
*/

type Node struct {
	attr fuse.Attr
}

func (n *Node) Attr() fuse.Attr {
	return n.attr
}
