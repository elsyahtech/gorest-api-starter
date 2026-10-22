package library

import "github.com/elsyahtech/gorest-api-starter/app/model"

// UserLibrary contains user-related application logic.
type UserLibrary struct {
	UserModel *model.UserModel
}

// DepartmentLibrary contains department-related application logic.
type DepartmentLibrary struct {
	DepartmentModel *model.DepartmentModel
}

// OrderLibrary contains order-related application logic.
type OrderLibrary struct {
	OrderModel *model.OrderModel
}

// NewUserLibrary creates a user library backed by the given model.
func NewUserLibrary(userModel *model.UserModel) *UserLibrary {
	return &UserLibrary{
		UserModel: userModel,
	}
}

// NewDepartmentLibrary creates a department library backed by the given model.
func NewDepartmentLibrary(departmentModel *model.DepartmentModel) *DepartmentLibrary {
	return &DepartmentLibrary{
		DepartmentModel: departmentModel,
	}
}

// NewOrderLibrary creates an order library backed by the given model.
func NewOrderLibrary(orderModel *model.OrderModel) *OrderLibrary {
	return &OrderLibrary{
		OrderModel: orderModel,
	}
}
