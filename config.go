
package bp

import (
	"fmt"
	"path/filepath"
	"code.google.com/p/go-sqlite/go1/sqlite3"
)

func (s *Session) updateConf() error {
	db := filepath.Join(s.root, ".bp/objects.db")
	c, err := sqlite3.Open(db)
	defer c.Close()
	if err != nil { return err }
	sql := "SELECT path FROM global"
	row := make(sqlite3.RowMap)
	for i, err := c.Query(sql); err == nil; err = i.Next() {
		var path string
		i.Scan(path, row)
		fmt.Println(path, row)
	}
	return nil
}

