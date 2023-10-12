package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	_ "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/migrations"
	"github.com/pressly/goose/v3"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	args := os.Args[1:]
	if len(args) < 2 {
		flag.Usage()
		return
	}
	command := args[1]

	file, fileErr := os.Open(configFile)
	if fileErr != nil {
		log.Fatal("Error opening config file.")
	}

	cfg, configErr := config.NewConfig(file)
	if configErr != nil {
		log.Fatal("Error parsing config file.")
	}

	db, err := goose.OpenDBWithDriver("postgres", cfg.Db.Dsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	var arguments []string
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, cfg.Db.MigrationsDir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
