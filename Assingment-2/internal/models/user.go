// * User model
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	PassHash string `json:"-"`
}
type NewUser struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
<<<<<<< HEAD
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required`
=======
	Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password string `json:"password" validate:"required"`
>>>>>>> f458359910e6e3cd468d0e422509099cf050d8c6
}
