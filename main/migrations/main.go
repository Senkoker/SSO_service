package main

import (
	"database/sql"
	"flag"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

//команда для запуска go run main.go --database_path=postgres://postgres:Masterchyrka890@localhost:5432/postgres?sslmode=disable --migration_path=C:/Golang_social_project/GRPC_Service_sso/internal/migration_db

func main() {
	var storagePath, migrationsPath string
	flag.StringVar(&storagePath, "database_path", "", "indicates the path of db")
	flag.StringVar(&migrationsPath, "migration_path", "", "indicates the path of migrations")
	flag.Parse()
	if storagePath == "" {
		panic("migration is reqiure")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
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
