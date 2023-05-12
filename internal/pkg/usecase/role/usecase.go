package role

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/role"
)

type Usecase interface {
	Create(role *model.Role) (*model.Role, error)
	ReadAll() (*[]model.Role, error)
	ReadById(id int) (*model.Role, error)
	Update(id int, role *model.Role) (*model.Role, error)
	Delete(id int) error
}

type usecase struct {
	repository role.Repository
}

func NewUsecase(repository role.Repository) Usecase {
	return &usecase{repository}
}

func (e *usecase) Create(role *model.Role) (*model.Role, error) {
	return e.repository.Create(role)
}

func (e *usecase) ReadAll() (*[]model.Role, error) {
	return e.repository.ReadAll()
}

func (e *usecase) ReadById(id int) (*model.Role, error) {
	return e.repository.ReadById(id)
}

func (e *usecase) Update(id int, role *model.Role) (*model.Role, error) {
	return e.repository.Update(id, role)
}

func (e *usecase) Delete(id int) error {
	return e.repository.Delete(id)
}
