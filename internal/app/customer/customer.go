package customer

import (
	"errors"
	"github.com/JanFant/LicenseServer/internal/app/db"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

//Customer структура покупателя
type Customer struct {
	Name    string `json:"name",sql:"name"`
	Address string `json:"address",sql:"address"`
	Phone   string `json:"phone",sql:"phone"`
	Email   string `json:"email",sql:"email"`
	Servers []int  `json:"servers",sql:"servers"`
}

func (customer *Customer) validate() error {
	return validation.ValidateStruct(customer,
		validation.Field(&customer.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&customer.Address, validation.Required),
		validation.Field(&customer.Phone, validation.Required, is.Int, validation.Length(11, 11)),
		validation.Field(&customer.Email, is.Email, validation.Required),
		validation.Field(&customer.Servers, validation.Length(0, 0)),
	)
}

func (customer *Customer) Create() error {
	if err := customer.validate(); err != nil {
		return err
	}
	var id int
	err := db.GetDB().Select(&id, "SELECT id FROM public.customers WHERE email = ? OR name = ?", customer.Email, customer.Name)
	if err != nil {
		return err
	}
	if id != 0 {
		return errors.New("customer must be unique")
	}

	db.GetDB().MustExec("INSERT INTO public.customers (name,address,phone,email) VALUES (:name,:address,phone,email)", &customer)
	return nil
}