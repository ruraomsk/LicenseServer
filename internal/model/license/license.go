package license

import (
	"github.com/JanFant/LicenseServer/internal/app/db"
	u "github.com/JanFant/LicenseServer/internal/app/utils"
	"github.com/JanFant/LicenseServer/internal/model/customer"
	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
	"net/http"
	"time"
)

var key = "asdqweqwe123dzsd12312cxq"

//License информация о лицензии клиента (БД?)
type License struct {
	Id        int       `json:"id",sql:"id"`               //уникальный номер сервера
	NumDevice int       `json:"numDev",sql:"numdev"`       //количество устройств
	YaKey     string    `json:"yaKey",sql:"yakey"`         //ключ яндекса
	TokenPass string    `json:"tokenPass",sql:"tokenpass"` //пароль для шифрования токена https запросов
	EndTime   time.Time `json:"endtime",sql:"endtime"`     //время окончания лицензии
	Token     string    `json:"token",sql:"token"`         //созданный токен
}

//LicenseToken токен лицензии клиента
type Token struct {
	NumDevice int    //количество устройств
	YaKey     string //ключ яндекса
	TokenPass string //пароль для шифрования токена https запросов
	Name      string //название фирмы
	Phone     string //телефон фирмы
	Id        int    //уникальный номер сервера
	Email     string //почта фирмы
	jwt.StandardClaims
}

func (license *License) validate() error {
	return validation.ValidateStruct(license,
		validation.Field(&license.NumDevice, validation.Required, validation.Min(0), validation.Max(1000)),
		validation.Field(&license.YaKey, validation.Required),
		validation.Field(&license.EndTime, validation.Required),
	)
}

func (license *License) CreateLicense(idCustomer int) u.Response {
	err := license.validate()
	if err != nil {
		return u.Message(http.StatusBadRequest, err.Error())
	}
	//генерация ключа
	license.TokenPass = u.GenerateRandomKey(100)
	var idLicense int
	row := db.GetDB().QueryRow(`INSERT INTO public.license (numdev, yakey, tokenpass, endtime) VALUES ($1, $2, $3, $4) RETURNING id`,
		license.NumDevice, license.YaKey, license.TokenPass, string(pq.FormatTimestamp(license.EndTime)))
	if err := row.Scan(&idLicense); err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	_, err = db.GetDB().Exec(`UPDATE public.customers SET servers = array_append(servers, $1) WHERE id = $2`, idLicense, idCustomer)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}

	return u.Message(http.StatusOK, "license record created")
}

func (license *License) CreateToken(id int) u.Response {
	var customerInfo customer.Customer
	err := customerInfo.Get(id)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	//создаем токен
	tk := &Token{
		Name:      customerInfo.Name,
		YaKey:     license.YaKey,
		Email:     customerInfo.Email,
		NumDevice: license.NumDevice,
		Phone:     customerInfo.Phone,
		TokenPass: license.TokenPass,
		Id:        license.Id}
	//врямя выдачи токена
	tk.IssuedAt = time.Now().Unix()
	//время когда закончится действие токена
	tk.ExpiresAt = license.EndTime.Unix()

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(key))

	//сохраняем токен в БД
	//GetDB().Exec("update public.accounts set token = ? where login = ?", account.Token, account.Login)

	//Формируем ответ
	resp := u.Message(http.StatusOK, "LicenseToken")
	resp.Obj["token"] = tokenString
	resp.Obj["tk"] = tk
	return resp
}
