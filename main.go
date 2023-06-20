package main

import (
	"path"

	"github.com/felipekzig/gnd/internal/cli"
	"github.com/felipekzig/gnd/internal/domain"
	"github.com/mitchellh/go-homedir"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	h := getHomeDir()
	path := path.Join(h, ".gnd.db")

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	ts := domain.NewTaskService(db)

	cli.Execute(ts)
}

func getHomeDir() string {
	h, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return h
}
