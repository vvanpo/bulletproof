
package bp

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"
)

type config struct {
}

func getConfig(root string) (*config, error) {
	pathname := filepath.Join(root, ".bp/config")
	text, err := ioutil.ReadFile(pathname)
	r := bytes.NewReader(text)
	s := bufio.NewScanner(r)
	err = s.Err()
	c := new(config)
	return c, err
}
