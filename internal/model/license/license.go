package license

import (
	"github.com/JanFant/LicenseServer/internal/app/db"
	u "github.com/JanFant/LicenseServer/internal/app/utils"
	"github.com/JanFant/LicenseServer/internal/model/customer"
	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"net/http"
	"time"
)

var key = "yreRmn6JKVv1md1Yh1PptBIjtGrL8pRjo8sAp5ZPlR6zK8xjxnzt6mGi6mtjWPJ6lz1HbhgNBxfSReuqP9ijLQ4JiWLQ4ADHefWVgtTzeI35pqB6hsFjOWufdAW8UEdK9ajm3T76uQlucUP2g4rUV8B9gTMoLtkn5Pxk6G83YZrvAIR7ddsd5PreTwGDoLrS6bdsbJ7u"

//License информация о лицензии клиента (БД?)
type License struct {
	Id        int       `json:"id",sql:"id"`                 //уникальный номер сервера
	NumDev    int       `json:"numdev",sql:"numdev"`         //количество устройств
	NumAcc    int       `json:"numacc",sql:"numacc"`         //колическво аккаунтов
	YaKey     string    `json:"yakey",sql:"yakey"`           //ключ яндекса
	TokenPass string    `json:"tokenpass",sql:"tokenpass"`   //пароль для шифрования токена https запросов
	TechEmail []string  `json:"tech_email",sql:"tech_email"` //почта для отправки сообщений в тех поддержку
	EndTime   time.Time `json:"endtime",sql:"endtime"`       //время окончания лицензии
	Token     string    `json:"token",sql:"token"`           //созданный токен
}

//LicenseToken токен лицензии клиента
type Token struct {
	NumDevice int    //количество устройств
	YaKey     string //ключ яндекса
	TokenPass string //пароль для шифрования токена https запросов
	NumAcc    int
	Name      string   //название фирмы
	Phone     string   //телефон фирмы
	Id        int      //уникальный номер сервера
	TechEmail []string //почта для отправки сообщений в тех поддержку
	Email     string   //почта фирмы
	jwt.StandardClaims
}

func (license *License) validate() error {
	return validation.ValidateStruct(license,
		validation.Field(&license.NumDev, validation.Required, validation.Min(1), validation.Max(1000)),
		validation.Field(&license.NumAcc, validation.Required, validation.Min(1), validation.Max(1000)),
		validation.Field(&license.YaKey, validation.Required),
		validation.Field(&license.EndTime, validation.Required),
		validation.Field(&license.TechEmail, validation.Required),
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
	row := db.GetDB().QueryRow(`INSERT INTO public.license (numdev, numacc, yakey, tokenpass, endtime, token, tech_email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		license.NumDev, license.NumAcc, license.YaKey, license.TokenPass, string(pq.FormatTimestamp(license.EndTime)), license.Token, pq.Array(license.TechEmail))
	if err := row.Scan(&idLicense); err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	_, err = db.GetDB().Exec(`UPDATE public.customers SET servers = array_append(servers, $1) WHERE id = $2`, idLicense, idCustomer)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}

	return u.Message(http.StatusOK, "license record created")
}

func (license *License) CreateToken(clientID, tokenID int) u.Response {
	var customerInfo customer.Customer
	err := customerInfo.Get(clientID)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	if customerInfo.ID == 0 {
		return u.Message(http.StatusBadRequest, "this client doesn't exist")
	}

	err = db.GetDB().QueryRow("SELECT id, numdev, numacc, yakey, tokenpass, token, tech_email, endtime FROM public.license WHERE id = $1", tokenID).Scan(
		&license.Id, &license.NumDev, &license.NumAcc, &license.YaKey, &license.TokenPass, &license.Token, pq.Array(&license.TechEmail), &license.EndTime)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}

	have := false
	for _, server := range customerInfo.Servers {
		if server == int64(license.Id) {
			have = true
		}
	}
	if !have {
		return u.Message(http.StatusInternalServerError, "this client doesn't own a license")
	}

	//создаем токен
	tk := &Token{
		Name:      customerInfo.Name,
		YaKey:     license.YaKey,
		Email:     customerInfo.Email,
		NumDevice: license.NumDev,
		Phone:     customerInfo.Phone,
		TokenPass: license.TokenPass,
		TechEmail: license.TechEmail,
		NumAcc:    license.NumAcc,
		Id:        license.Id}
	//врямя выдачи токена
	tk.IssuedAt = time.Now().Unix()
	//время когда закончится действие токена
	tk.ExpiresAt = license.EndTime.Unix()

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(key))

	//сохраняем токен в БД
	_, err = db.GetDB().Exec(`UPDATE  public.license SET token = $1 WHERE id = $2`, tokenString, license.Id)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}

	//Формируем ответ
	resp := u.Message(http.StatusOK, "LicenseToken")
	resp.Obj["token"] = tokenString
	resp.Obj["tk"] = tk
	return resp
}

func GetAllLicenseInfo(id int) u.Response {
	var customerInfo customer.Customer
	err := customerInfo.Get(id)
	if err != nil {
		return u.Message(http.StatusInternalServerError, err.Error())
	}
	if customerInfo.ID == 0 {
		return u.Message(http.StatusBadRequest, "this client doesn't exist")
	}
	var allLicense []License
	if len(customerInfo.Servers) > 0 {
		query, args, err := sqlx.In("SELECT id, numdev, numacc, yakey, tokenpass, token, tech_email, endtime FROM public.license WHERE id IN (?)", customerInfo.Servers)
		if err != nil {
			return u.Message(http.StatusInternalServerError, err.Error())
		}
		query = db.GetDB().Rebind(query)
		rows, err := db.GetDB().Queryx(query, args...)
		if err != nil {
			return u.Message(http.StatusInternalServerError, err.Error())
		}
		for rows.Next() {
			var temp License
			err := rows.Scan(&temp.Id, &temp.NumDev, &temp.NumAcc, &temp.YaKey, &temp.TokenPass, &temp.Token, pq.Array(&temp.TechEmail), &temp.EndTime)
			if err != nil {
				return u.Message(http.StatusInternalServerError, err.Error())
			}
			allLicense = append(allLicense, temp)
		}
	}
	if len(allLicense) == 0 {
		allLicense = make([]License, 0)
	}
	resp := u.Message(http.StatusOK, "all license info")
	resp.Obj["customer"] = customerInfo
	resp.Obj["licenses"] = allLicense
	return resp
}
