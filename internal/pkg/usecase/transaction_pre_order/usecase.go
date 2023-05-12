package transaction_pre_order

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/enum"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/role"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/transaction_pre_order"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/transaction_pre_order_user"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/user"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/request"
	"github.com/jinzhu/gorm"
)

type Usecase interface {
	Create(transactionPreOrder *request.TransactionPreOrder) (*model.TransactionPreOrder, error)
	ReadAll() (*[]model.TransactionPreOrder, error)
	ReadAllBy(req request.TransactionPreOrderPaged) (*[]model.TransactionPreOrder, error)
	ReadById(id int) (*model.TransactionPreOrder, error)
	ReadBy(req request.TransactionPreOrderDetail) (*helper.TransactionPreOrderResponse, error)
	Update(id int, transactionPreOrder *model.TransactionPreOrder) (*model.TransactionPreOrder, error)
	Delete(id int) error
	Count(req request.TransactionPreOrderPaged) int
	Summary(companyId int) (interface{}, error)
	Verification(transactionPreOrderUser *request.TransactionPreOrderUser) (*helper.TransactionPreOrderResponse, error)
}

type usecase struct {
	transactionPreOrderRepository     transaction_pre_order.Repository
	transactionPreOrderUserRepository transaction_pre_order_user.Repository
	userRepository                    user.Repository
	roleRepository                    role.Repository
}

func NewUsecase(transactionPreOrderRepository transaction_pre_order.Repository, transactionPreOrderUserRepository transaction_pre_order_user.Repository, userRepository user.Repository, roleRepository role.Repository) Usecase {
	return &usecase{transactionPreOrderRepository, transactionPreOrderUserRepository, userRepository, roleRepository}
}

func (e *usecase) Create(transactionPreOrder *request.TransactionPreOrder) (*model.TransactionPreOrder, error) {

	m := &model.TransactionPreOrder{
		PariProductID:     transactionPreOrder.PariProductId,
		PariTransactionID: transactionPreOrder.PariTransactionId,
		ProductID:         transactionPreOrder.ProductID,
		CompanyID:         transactionPreOrder.CompanyID,
		Quantity:          transactionPreOrder.Quantity,
		Status:            transactionPreOrder.Status,
		ActualPrice:       transactionPreOrder.ActualPrice,
		BuyerName:         transactionPreOrder.BuyerName,
		BuyerAddress:      transactionPreOrder.BuyerAddress,
		BuyerContact:      transactionPreOrder.BuyerContact,
	}

	return e.transactionPreOrderRepository.Create(m)
}

func (e *usecase) ReadAll() (*[]model.TransactionPreOrder, error) {
	return e.transactionPreOrderRepository.ReadAll()
}

func (e *usecase) ReadAllBy(req request.TransactionPreOrderPaged) (*[]model.TransactionPreOrder, error) {
	criteria := make(map[string]interface{})
	criteria["company_id"] = req.CompanyID

	if req.Status != "" {
		criteria["status"] = req.Status
	}

	if req.Commodity != "" {
		criteria["commodity"] = req.Commodity
	}

	fmt.Println(req)

	return e.transactionPreOrderRepository.ReadAllBy(criteria, req.Search, req.Page, req.Size)
}

func (e *usecase) ReadById(id int) (*model.TransactionPreOrder, error) {
	return e.transactionPreOrderRepository.ReadById(id)
}

func (e *usecase) ReadBy(req request.TransactionPreOrderDetail) (*helper.TransactionPreOrderResponse, error) {
	var isVerifiedByUser bool
	productModel, err := e.transactionPreOrderRepository.ReadById(req.ID)

	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	if req.UserID != 0 {
		countTransactionPreOrderUser := e.transactionPreOrderUserRepository.Count(map[string]interface{}{"transaction_pre_order_id": req.ID, "user_id": req.UserID})
		if countTransactionPreOrderUser > 0 {
			isVerifiedByUser = true
		}
	}

	result := &helper.TransactionPreOrderResponse{TransactionPreOrder: productModel, IsVerifiedByUser: isVerifiedByUser}

	return result, nil
}

func (e *usecase) Update(id int, transactionPreOrder *model.TransactionPreOrder) (*model.TransactionPreOrder, error) {
	return e.transactionPreOrderRepository.Update(id, transactionPreOrder)
}

func (e *usecase) Delete(id int) error {
	return e.transactionPreOrderRepository.Delete(id)
}

func (e *usecase) Verification(request *request.TransactionPreOrderUser) (*helper.TransactionPreOrderResponse, error) {

	productModel, err := e.transactionPreOrderRepository.ReadById(request.TransactionPreOrderID)
	if err != nil {
		helper.CommonLogger().Error(err)
	}

	productUser, err := e.transactionPreOrderUserRepository.ReadBy(map[string]interface{}{"transaction_pre_order_id": request.TransactionPreOrderID, "user_id": request.UserID, "company_id": request.CompanyID})
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("failed finding transaction pre order user: ", err)
	}

	r, err := e.roleRepository.ReadById(request.RoleID)
	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	// checking whether transactionPreOrderUser exists or not
	if productUser == nil {
		pu := &model.TransactionPreOrderUser{TransactionPreOrderID: request.TransactionPreOrderID, UserID: request.UserID, CompanyID: request.CompanyID}
		_, err := e.transactionPreOrderUserRepository.Create(pu)
		if err != nil {
			helper.CommonLogger().Error(err)
			return nil, err
		}
	}

	// checking product has been approved by all user in company
	countTransactionPreOrderUser := e.transactionPreOrderUserRepository.Count(map[string]interface{}{"company_id": request.CompanyID, "transaction_pre_order_id": request.TransactionPreOrderID})
	countUser := e.userRepository.Count(map[string]interface{}{"company_id": request.CompanyID, "role_id": r.ID})
	if countUser == countTransactionPreOrderUser {
		productModel.Status = enum.Approved
		_, err := e.transactionPreOrderRepository.Update(productModel.ID, productModel)
		if err != nil {
			helper.CommonLogger().Error(err)
			return nil, err
		}
	}

	result := &helper.TransactionPreOrderResponse{TransactionPreOrder: productModel, IsVerifiedByUser: true}
	return result, nil
}

func (e *usecase) Count(req request.TransactionPreOrderPaged) int {
	criteria := make(map[string]interface{})
	criteria["company_id"] = req.CompanyID

	if req.Status != "" {
		criteria["status"] = req.Status
	}

	if req.Commodity != "" {
		criteria["commodity"] = req.Commodity
	}

	return e.transactionPreOrderRepository.Count(criteria)
}

func (e *usecase) Summary(companyId int) (interface{}, error) {

	allTransactionPreOrder := e.transactionPreOrderRepository.Count(map[string]interface{}{"company_id": companyId})
	processingTransactionPreOrder := e.transactionPreOrderRepository.Count(map[string]interface{}{"company_id": companyId, "status": "processing"})
	approvedTransactionPreOrder := e.transactionPreOrderRepository.Count(map[string]interface{}{"company_id": companyId, "status": "approved"})
	rejectedTransactionPreOrder := e.transactionPreOrderRepository.Count(map[string]interface{}{"company_id": companyId, "status": "rejected"})

	return map[string]interface{}{
		"all_product":        allTransactionPreOrder,
		"processing_product": processingTransactionPreOrder,
		"approved_product":   approvedTransactionPreOrder,
		"rejected_product":   rejectedTransactionPreOrder,
	}, nil
}
