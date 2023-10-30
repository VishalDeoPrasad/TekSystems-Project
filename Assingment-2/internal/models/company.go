package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name string `json:"companyName" validate:"required"`
	City string `json:"city" validate:"required"`
<<<<<<< HEAD
}
type Job struct {
	gorm.Model
	Name       string `json:"jobTitle"`
	Field      string `json:"field"`
	Experience uint   `json:"experience"`
=======
	Jobs []Job  `json:"jobs,omitempty" gorm:"foreignKey:CompanyId"`
}
type Job struct {
	gorm.Model
	Name       string `json:"title"`
	Field      string `json:"field"`
	Experience uint   `json:"experience"`
	CompanyId  uint64 `json:"companyId"`
>>>>>>> f458359910e6e3cd468d0e422509099cf050d8c6
}
