
package bp

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
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
		node string
	}
	private map[string]struct{
		path string
		flags int
	}
}

func (s *Session) updateCfg() error {
	path := filepath.Join(s.root, ".bp/global")
	text, _ := ioutil.ReadFile(path)
	r := bytes.NewReader(text)
	scanner := bufio.NewScanner(r)
	global := make(map[string]struct{
		path string; flags int; node string
	})
	err := parseCfg(global, scanner)
	return err
}

func parseCfg(files map[string]struct{
			path string; flags int; node string
		}, s *bufio.Scanner) error {
	var tag string
	for i := 1; s.Scan(); i++ {
		f := strings.Fields(s.Text())
		if len(f) == 0 || strings.HasPrefix(f[0], "#") { continue }
		if f[0] != "-" {
			if len(f) > 1 && !strings.HasSuffix(f[0], ":") && f[1] != ":" {
				return fmt.Errorf("Invalid configuration, line %d", i)
			}
			tag = f[0]
			//files[tag]
		} else {
			if len(f) == 2 {
				files[tag] = struct { f[1], 0, ""}
			}
		}
	}
	return nil
}
