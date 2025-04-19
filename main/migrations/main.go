package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"log"
)

type Migration struct {
	StoragePath    string `env:"STORAGE_PATH"`
	MigrationsPath string `env:"MIGRATIONS_PATH"`
}

func main() {
	var migration Migration
	err := cleanenv.ReadEnv(&migration)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", migration.StoragePath)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migration.MigrationsPath,
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Migrations did run successfully")
}

//err = m.Down()
//if err != nil && err != migrate.ErrNoChange {
//	log.Fatal(err)
//}
