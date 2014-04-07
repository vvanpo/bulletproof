
package bp

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"
)

const (
	recursive int = 1 << iota
	follow
	dataOnly
	metadataOnly
	encrypt
)

type config struct {
	sync map[string]struct{
		path string
		flags int
	}
	private map[string]struct{
		path string
		flags int
	}
}

func (s *Session) update() {
	path := filepath.Join(s.root, ".bp/global")
	text, err := ioutil.ReadFile(path)
	r := bytes.NewReader(text)
	scanner := bufio.NewScanner(r)
	err = scanner.Err()
	return
}
