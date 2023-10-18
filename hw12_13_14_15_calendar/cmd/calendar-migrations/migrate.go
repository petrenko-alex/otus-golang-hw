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
	os.Exit(run())
}

func run() int {
	flag.Parse()

	args := os.Args[1:]
	if len(args) < 2 {
		flag.Usage()
		return 0
	}
	command := args[1]

	file, fileErr := os.Open(configFile)
	if fileErr != nil {
		log.Println("Error opening config file.")

		return 1
	}

	cfg, configErr := config.NewConfig(file)
	if configErr != nil {
		log.Println("Error parsing config file.")

		return 1
	}

	db, err := goose.OpenDBWithDriver("postgres", cfg.Db.Dsn)
	if err != nil {
		log.Printf("goose: failed to open DB: %v\n\n", err)

		return 1
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
		log.Printf("goose %v: %v\n", command, err)

		return 1
	}

	return 0
}
