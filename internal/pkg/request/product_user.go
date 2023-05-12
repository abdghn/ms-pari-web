/*
 * Created on 27/05/22 09.30
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package request

type ProductUser struct {
	ProductID int `json:"product_id"`
	UserID    int `json:"user_id"`
	CompanyID int `json:"company_id"`
	RoleID    int `json:"role_id"`
}
