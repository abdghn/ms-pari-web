package request

type TransactionPreOrderUser struct {
	TransactionPreOrderID int `json:"transaction_pre_order_id"`
	UserID                int `json:"user_id"`
	CompanyID             int `json:"company_id"`
	RoleID                int `json:"role_id"`
}
