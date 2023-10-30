package services

import (
	"context"
	"golang/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source service.go -destination mockmodels/service_mock.go -package mockmodels

type Service interface {
<<<<<<< HEAD
	// CreatInventory(ctx context.Context, ni NewInventory, userId uint) (Inventory, error)
	// ViewInventory(ctx context.Context, userId string) ([]Inventory, float64, error)
	CreateCompany(ctx context.Context, newComp models.Company) (models.Company, error)
=======
	GetAllJobs(ctx context.Context) ([]models.Job,error)
	GetJobById(ctx context.Context, jobId string) (models.Job,error)
	FetchJobByCompanyId(ctx context.Context, companyId string) ([]models.Job, error)
	JobByCompanyId(jobs []models.Job,compId string,)([]models.Job,error)
	FetchCompanyByID(ctx context.Context, companyId string) (models.Company, error) 
	ViewCompanies(ctx context.Context)([]models.Company,error)
	CreateCompany(ctx context.Context, newComp models.Company) (models.Company, error) 
>>>>>>> f458359910e6e3cd468d0e422509099cf050d8c6
	CreateUser(ctx context.Context, nu models.NewUser) (models.User, error)
	Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims, error)
	AutoMigrate() error
}

type Store struct {
	Service
}

func NewStore(s Service) Store {
	return Store{Service: s}
}
