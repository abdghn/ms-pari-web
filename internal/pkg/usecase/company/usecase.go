package company

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/company"
)

type Usecase interface {
	Create(company *model.Company) (*model.Company, error)
	ReadAll() (*[]model.Company, error)
	ReadById(id int) (*model.Company, error)
	Update(id int, company *model.Company) (*model.Company, error)
	Delete(id int) error
}

type usecase struct {
	repository company.Repository
}

func NewUsecase(repository company.Repository) Usecase {
	return &usecase{repository}
}

func (e *usecase) Create(company *model.Company) (*model.Company, error) {
	return e.repository.Create(company)
}

func (e *usecase) ReadAll() (*[]model.Company, error) {
	return e.repository.ReadAll()
}

func (e *usecase) ReadById(id int) (*model.Company, error) {
	return e.repository.ReadById(id)
}

func (e *usecase) Update(id int, company *model.Company) (*model.Company, error) {
	return e.repository.Update(id, company)
}

func (e *usecase) Delete(id int) error {
	return e.repository.Delete(id)
}
