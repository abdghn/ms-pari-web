/*
 * Created on 01/04/22 14.58
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package config

import (
	"fmt"
	"log"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func DbConnect(dbUser, dbPass, dbHost, dbPort, dbName string) *gorm.DB {
	consStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName)

	db, err := gorm.Open("mysql", consStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		log.Fatal("Error when connect db " + consStr + " : " + err.Error())
		return nil
	}

	db.Debug().AutoMigrate(
		model.User{},
		model.Role{},
		model.Company{},
		model.Giro{},
		model.Product{},
		model.ProductUser{},
		model.TransactionPreOrder{},
		model.TransactionPreOrderUser{},
	)
	return db
}
