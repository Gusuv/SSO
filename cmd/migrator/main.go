package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Gusuv/sso/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	var cfgPath string
	var migrationPath string
	var step int

	flag.StringVar(&cfgPath, "config", "", "path to config")
	flag.StringVar(&migrationPath, "migration-path", "", "path to migration files")
	flag.IntVar(&step, "step", 0, "step for migration")
	flag.Parse()

	if cfgPath == "" {
		fmt.Println("config path is requitred")
		os.Exit(1)
	}
	c := config.MustLoad(cfgPath)

	dbUrl := c.DbUrl()

	absPath, err := filepath.Abs(migrationPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m, err := migrate.New("file://"+absPath, dbUrl)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var migrateError error
	if step != 0 {
		migrateError = m.Steps(step)
	} else {
		migrateError = m.Up()
	}

	if migrateError != nil && !errors.Is(migrateError, migrate.ErrNoChange) {
		fmt.Println(migrateError)
		os.Exit(1)
	}

	fmt.Println("migrations applied successfully")
}
