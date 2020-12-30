package database

import (
	"database/sql"
	"log"

	// Justifying comment
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"

	// Justifying comment
	_ "github.com/golang-migrate/migrate/source/file"
)

// Db database pointer
var Db *sql.DB

// InitDB initializes the DB
func InitDB() {
	// Use root:dbpass@tcp(172.17.0.2)/hackernews, if you're using Windows.
	db, err := sql.Open("mysql", "root:dbpass@tcp(localhost)/hackernews")
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	Db = db
}

// Migrate moves stuff
func Migrate() {
	if err := Db.Ping(); err != nil {
		log.Fatal(err)
	}
	driver, _ := mysql.WithInstance(Db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/mysql",
		"mysql",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

}
