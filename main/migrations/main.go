package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
)

type Migration struct {
	Host           string `env:"HOST_POSTGRES"`
	DbName         string `env:"DB_NAME"`
	UserName       string `env:"USER_NAME"`
	UserPass       string `env:"USER_PASS"`
	MigrationsPath string `env:"MIGRATIONS_PATH"`
}

const (
	driver  = "pgx"
	address = "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable"
)

func main() {
	var migration Migration
	var command string
	flag.StringVar(&command, "command", "", "write command for lead migration")
	flag.Parse()
	err := cleanenv.ReadEnv(&migration)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(command)
	db, err := goose.OpenDBWithDriver(driver, fmt.Sprintf(address, migration.Host, migration.UserName, migration.UserPass, migration.DbName))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close()
	ctx := context.Background()
	err = goose.RunContext(ctx, command, db, migration.MigrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migration completed successfully")
}
