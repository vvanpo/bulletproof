
package filetracker

import (
	"code.google.com/p/gcfg"
	"path/filepath"
)

// Filled with gcfg-parsed configuration file (code.google.com/p/gcfg)
type config struct {
	Sync map[string]*struct{
		v []string
	}
	Local map[string]*struct{}
}

func GetConfig(root string) (*config, error) {
	pathname := filepath.Join(root, ".bp/config")
	c := new(config)
	err := gcfg.ReadFileInto(c, pathname)
	return c, err
}
