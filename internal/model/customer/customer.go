package customer

import (
	"github.com/JanFant/LicenseServer/internal/app/db"
	u "github.com/JanFant/LicenseServer/internal/app/utils"
	"github.com/JanFant/TLServer/logger"
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

//validate проверка данных клиента
func (customer *Customer) validate() error {
	return validation.ValidateStruct(customer,
		validation.Field(&customer.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&customer.Address, validation.Required),
		validation.Field(&customer.Phone, validation.Required, is.Int, validation.Length(11, 11)),
		validation.Field(&customer.Email, is.Email, validation.Required),
		validation.Field(&customer.Servers, validation.Length(0, 0)),
	)
}

//Create создание клиента
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
	_, err = db.GetDB().Exec(`INSERT INTO public.customers (name, address, phone, email,servers) VALUES ($1, $2, $3, $4, $5)`,
		customer.Name, customer.Address, customer.Phone, customer.Email, pq.Array(customer.Servers))
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	return u.Message(http.StatusOK, "customer created")
}

//Delete удалить клиента
func (customer *Customer) Delete() u.Response {
	_, err := db.GetDB().Exec(`DELETE FROM public.customers WHERE id=$1`, customer.ID)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	return u.Message(http.StatusOK, "customer deleted")
}

//Update обновить данные клиента
func (customer *Customer) Update() u.Response {
	_, err := db.GetDB().Exec(`UPDATE public.customers SET name=$1, address=$2, phone=$3, email=$4 WHERE id=$5`,
		customer.Name, customer.Address, customer.Phone, customer.Email, customer.ID)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	return u.Message(http.StatusOK, "customer update")
}

func (customer *Customer) Get(id int) error {
	rows, err := db.GetDB().Query(`SELECT id, name, address, servers, phone, email FROM public.customers WHERE id=$1`, id)
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Address, pq.Array(&customer.Servers), &customer.Phone, &customer.Email)
		if err != nil {
			return err
		}
	}
	return nil
}

//GetAllCustomers получить всех клиента
func GetAllCustomers() []Customer {
	rows, err := db.GetDB().Query(`SELECT id, name, address, servers, phone, email FROM public.customers`)
	if err != nil {
		logger.Error.Printf("|Message: %v", err.Error())
		return make([]Customer, 0)
	}
	var customers []Customer
	for rows.Next() {
		var temp Customer
		err := rows.Scan(&temp.ID, &temp.Name, &temp.Address, pq.Array(&temp.Servers), &temp.Phone, &temp.Email)
		if err != nil {
			logger.Error.Printf("|Message: %v", err.Error())
			return make([]Customer, 0)
		}
		customers = append(customers, temp)
	}
	if len(customers) == 0 {
		customers = make([]Customer, 0)
	}
	return customers
}
