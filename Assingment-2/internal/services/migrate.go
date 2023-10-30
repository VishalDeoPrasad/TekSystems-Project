package services

import (
	"errors"
	"golang/internal/models"

	"gorm.io/gorm"
)

/*
The struct has one field, db, which is a pointer to a GORM (gorm.DB) database
connection. This field holds the database connection instance associated with
this DbConnStruct.
*/
type DbConnStruct struct {
	db *gorm.DB
}

// This is a constructor function for creating instances of the DbConnStruct type.
// It takes a *gorm.DB as an argument, which is a GORM database connection.
func NewConn(dbInstance *gorm.DB) (*DbConnStruct, error) {
	if dbInstance == nil {
		return nil, errors.New("provide the databse instance")
	}
	/*
		it creates a new DbConnStruct instance and initializes its db field with 
		the provided database connection.*/
	return &DbConnStruct{db: dbInstance}, nil
}

func (s *DbConnStruct) AutoMigrate() error {

	// AutoMigrate function will ONLY create tables, missing columns and missing indexes, and WON'T change existing column's type or delete unused columns
	/*
	AutoMigrate(...): This is a method of the migrator object. It's used to 
	automatically create or update database tables based on the provided model 
	structs. In your case, it's being used to migrate the User, Company, and 
	Job tables to match the corresponding model structs. If the tables already 
	exist, it will add any missing columns or indexes as necessary, but it won't 
	change existing column types*/
	err := s.db.Migrator().AutoMigrate(&models.User{}, &models.Company{}, &models.Job{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return err
	}
	return nil
}
