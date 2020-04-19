package customer

import (
	"github.com/JanFant/LicenseServer/internal/app/db"
	u "github.com/JanFant/LicenseServer/internal/app/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/lib/pq"
	"net/http"
)

//Customer структура покупателя
type Customer struct {
	ID      int     `json:"id",sql:"id"`
	Name    string  `json:"name",sql:"name"`
	Address string  `json:"address",sql:"address"`
	Phone   string  `json:"phone",sql:"phone"`
	Email   string  `json:"email",sql:"email"`
	Servers []int64 `json:"servers",sql:"servers"`
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

func (customer *Customer) Create() u.Response {
	if err := customer.validate(); err != nil {
		return u.Message(http.StatusBadRequest, err.Error())
	}
	var id int
	row, err := db.GetDB().NamedQuery(`SELECT id FROM public.customers WHERE email = :email OR name = :name`, customer)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	for row.Next() {
		_ = row.Scan(&id)
		if id > 0 {
			return u.Message(http.StatusBadRequest, "this customer has already been created")
		}
	}
	_, err = db.GetDB().Exec(`INSERT INTO public.customers (name, address, phone, email,servers) VALUES ($1, $2, $3, $4, $5)`, customer.Name, customer.Address, customer.Phone, customer.Email, pq.Array(customer.Servers))
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	return u.Message(http.StatusOK, "customer created")
}

func GetAllCustomers() u.Response {
	rows, err := db.GetDB().Query(`SELECT id, name, address, servers, phone, email FROM public.customers`)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	var customers []Customer
	for rows.Next() {
		var temp Customer
		err := rows.Scan(&temp.ID, &temp.Name, &temp.Address, pq.Array(&temp.Servers), &temp.Phone, &temp.Email)
		if err != nil {
			return u.Message(http.StatusInternalServerError, err.Error())
		}
		customers = append(customers, temp)
	}
	if len(customers) == 0 {
		customers = make([]Customer, 0)
	}
	resp := u.Message(http.StatusOK, "all customers")
	resp.Obj["customers"] = customers
	return resp
}
