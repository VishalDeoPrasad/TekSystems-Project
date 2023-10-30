// * Database intialization and configuration
package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// it returns two values: a pointer to a gorm.DB object (a database connection) and an error.
func Open() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	/*
		This line establishes a database connection using GORM.
		it's using the PostgreSQL database driver, postgres.Open(dsn),
		and passing an empty gorm.Config{}.

		By passing &gorm.Config{}, you are effectively using the
		default configuration settings for GORM.*/
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	/*
		db: This is a variable of type *gorm.DB, which represents a GORM database
		connection. It's a pointer to the database connection that you'll use to
		interact with the database.
		err: This is an error variable that will store any potential error that might 
		occur during the process of opening the database connection. If the 
		connection is successful, err will be nil. If there's an issue, err 
		will contain details about the error.*/
	if err != nil {
		return nil, err
	}
	return db, nil
}
