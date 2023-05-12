/*
 * Created on 07/04/22 06.07
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package auth

import (
	"errors"
	"regexp"
	"time"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/company"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/giro"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/role"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/user"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/request"
)

type Usecase interface {
	Register(user request.User) (*model.User, error)
	BulkRegister(users request.Users) ([]*model.User, error)
	Login(user *model.User) (*model.User, error)
	ValidateGiro(code string) (*model.Giro, error)
	GetToken(clientKey, secretKey string) (key *request.OpenKey, err error)
}

type usecase struct {
	userRepository    user.Repository
	giroRepository    giro.Repository
	companyRepository company.Repository
	roleRepository    role.Repository
}

func NewUsecase(userRepository user.Repository, giroRepository giro.Repository, roleRepository role.Repository, companyRepository company.Repository) Usecase {
	return &usecase{userRepository: userRepository, giroRepository: giroRepository, roleRepository: roleRepository, companyRepository: companyRepository}
}

func (e *usecase) Register(u request.User) (*model.User, error) {
	helper.HashPassword(&u.Password)

	r, err := e.roleRepository.ReadByName(u.Role)
	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	c, err := e.companyRepository.ReadById(u.CompanyID)
	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	newUser := &model.User{
		RoleID:            r.ID,
		CompanyID:         c.ID,
		Name:              u.Name,
		Email:             u.Email,
		VerificationLevel: u.VerificationLevel,
		Password:          u.Password,
	}

	m, err := e.userRepository.Create(newUser)
	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	m.Password = ""
	m.RoleName = u.Role
	m.CompanyName = c.Name

	return m, nil
}

func (e *usecase) BulkRegister(users request.Users) ([]*model.User, error) {
	listUsers := make([]*model.User, 0)
	for _, u := range users {

		helper.HashPassword(&u.Password)

		r, err := e.roleRepository.ReadByName(u.Role)
		if err != nil {
			helper.CommonLogger().Error(err)
			return nil, err
		}

		c, err := e.companyRepository.ReadById(u.CompanyID)
		if err != nil {
			helper.CommonLogger().Error(err)
			return nil, err
		}

		newUser := &model.User{
			RoleID:            r.ID,
			CompanyID:         c.ID,
			Name:              u.Name,
			Email:             u.Email,
			VerificationLevel: u.VerificationLevel,
			Password:          u.Password,
		}

		m, err := e.userRepository.Create(newUser)
		if err != nil {
			helper.CommonLogger().Error(err)
			return nil, err
		}

		m.Password = ""
		m.RoleName = u.Role
		m.CompanyName = c.Name
		listUsers = append(listUsers, m)
	}

	return listUsers, nil
}

func (e *usecase) Login(user *model.User) (*model.User, error) {
	return e.userRepository.ReadByEmail(user.Email)
}

func (e *usecase) ValidateGiro(code string) (*model.Giro, error) {
	return e.giroRepository.ReadByCode(code)
}

func (e *usecase) GetToken(clientKey, secretKey string) (key *request.OpenKey, err error) {
	var (
		expduration = 15 //minutes
		expTime     = (time.Now().Add(time.Minute * time.Duration(expduration))).Format("2006-01-02 15:04:05")
	)

	// validation
	if clientKey == "" || secretKey == "" {
		return nil, errors.New("empty client or secret key")
	}

	if helper.ClientKey != clientKey || helper.SecretKey != secretKey {
		return nil, errors.New("incorrect client or secret key")
	}

	b1, err := helper.RsaEncrypt([]byte("PARI%" + expTime))
	if err != nil {
		helper.CommonLogger().Error(err)
		panic(err)
	}

	token := helper.Base64Enc(b1)
	re := regexp.MustCompile(`\r?\n`)
	token = re.ReplaceAllString(token, "")
	mod := request.OpenKey{
		Token:     token,
		ExpiredAt: expTime,
	}

	return &mod, nil
}
