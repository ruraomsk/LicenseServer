package model

import (
	"encoding/json"
	"errors"
	"github.com/JanFant/LicenseServer/internal/app/db"
	"github.com/JanFant/easyLog"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/lib/pq"
)

//Customer структура покупателя
type Customer struct {
	ID       int       `json:"id" ,sql:"id"`
	Name     string    `json:"name" ,sql:"name"`
	Address  string    `json:"address" ,sql:"address"`
	Phone    string    `json:"phone" ,sql:"phone"`
	Email    string    `json:"email" ,sql:"email"`
	Servers  []int64   `json:"servers" ,sql:"servers"`
	Licenses []License `json:"licenses" ,sql:"licenses"`
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
func (customer *Customer) Create() error {
	customer.Servers = make([]int64, 0)
	if err := customer.validate(); err != nil {
		return err
	}
	var id int
	row, err := db.GetDB().NamedQuery(`SELECT id FROM public.customers WHERE email = :email OR name = :name`, customer)
	if err != nil {
		return errors.New("ошибка связи с БД")
	}
	for row.Next() {
		_ = row.Scan(&id)
		if id > 0 {
			return errors.New("пользователь с таким именем уже существует")
		}
	}
	_, err = db.GetDB().Exec(`INSERT INTO public.customers (name, address, phone, email,servers) VALUES ($1, $2, $3, $4, $5)`,
		customer.Name, customer.Address, customer.Phone, customer.Email, pq.Array(customer.Servers))
	if err != nil {
		return errors.New("ошибка связи с БД")
	}
	return nil
}

//Delete удалить клиента
func (customer *Customer) Delete() error {
	_, err := db.GetDB().Exec(`DELETE FROM public.customers WHERE id=$1`, customer.ID)
	if err != nil {
		return err
	}
	return nil
}

//Update обновить данные клиента
func (customer *Customer) Update() error {
	var exists bool
	err := db.GetDB().QueryRow(`SELECT exists (SELECT id FROM public.customers WHERE id = $1)`, customer.ID).Scan(&exists)
	if err != nil {
		return errors.New("ошибка связи с БД")
	}
	if exists {
		_, err = db.GetDB().Exec(`UPDATE public.customers SET name=$1, address=$2, phone=$3, email=$4 WHERE id=$5`,
			customer.Name, customer.Address, customer.Phone, customer.Email, customer.ID)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("такой записи не существует")
	}
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

//GetAllInfo получить всех клиента
func GetAllInfo() []Customer {
	rows, err := db.GetDB().Query(`SELECT  
											cust.id, cust.name, cust.address, cust.servers, cust.phone, cust.email, 
											json_strip_nulls(json_agg(json_build_object('id',lic.id,
																						'numdev',lic.numdev,
																						'numacc',lic.numacc,
																						'yakey',lic.yakey,
																						'tokenpass',lic.tokenpass,
																						'token',lic.token,
																						'tech_email',lic.tech_email,
																						'endtime',lic.endtime))) as licenses
											FROM public.customers as cust
											LEFT JOIN public.license as lic ON lic.id = ANY(cust.servers)
											GROUP BY cust.id;`)
	if err != nil {
		easyLog.Error.Printf("|Message: %v", err.Error())
		return make([]Customer, 0)
	}
	var customers []Customer
	for rows.Next() {
		var (
			temp        Customer
			licensesStr string
		)
		err := rows.Scan(&temp.ID, &temp.Name, &temp.Address, pq.Array(&temp.Servers), &temp.Phone, &temp.Email, &licensesStr)
		if err != nil {
			easyLog.Error.Printf("|Message: %v", err.Error())
			return make([]Customer, 0)
		}
		err = json.Unmarshal([]byte(licensesStr), &temp.Licenses)
		if err != nil {
			easyLog.Error.Printf("|Message: %v", err.Error())
		}
		if temp.Licenses[0].Id == 0 {
			temp.Licenses = make([]License, 0)
		}
		customers = append(customers, temp)
	}
	if len(customers) == 0 {
		customers = make([]Customer, 0)
	}
	return customers
}
