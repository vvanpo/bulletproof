
package filetracker

import (
	"path/filepath"
)

type Config struct {
}

func GetConfig(root string) (*Config, error) {
	pathname := filepath.Join(root, ".bp/config")
	c := new(Config)
	return c, err
}
