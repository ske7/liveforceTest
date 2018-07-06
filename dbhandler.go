package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

// DB connection handle
var db *sql.DB

// SetupDB make DB connection and create necessary table
func SetupDB() (*sql.DB, error) {
	// credentials and db config
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "12qwas"
		dbname   = "postgres"
	)

	dataSource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	// create user table
	sql := `
          DO $$
          BEGIN
          IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_type') THEN
            CREATE TYPE role_type AS
            ENUM ('admin', 'user');
          END IF;
          END $$;

          CREATE TABLE IF NOT EXISTS public.users (
              id SERIAL NOT NULL PRIMARY KEY,
              name varchar(255) NOT NULL,
              email varchar(255) NOT NULL,
              password varchar(255) NOT NULL,
              date_created TIMESTAMPTZ NOT NULL,
              role role_type NOT NULL DEFAULT 'user',
            CONSTRAINT user_email_key UNIQUE (email)
          );
          `
	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	// insert admin record if not exists else do nothing
	sql = `INSERT INTO public.users (name, email, password, date_created, role)
         VALUES ($1, $2, $3, CURRENT_TIMESTAMP, $4)
	       ON CONFLICT (email) DO NOTHING`
	_, err = db.Exec(sql, "admin", "adminemail@site.com", HashPassword([]byte("admpwd")), "admin")
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

// GetDB return DB connection handle
func GetDB() *sql.DB {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
