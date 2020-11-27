package custMain

var (
	typeError          = "error"
	typeClose          = "close"
	typeCustInfo       = "custInfo"
	typeCustUpdate     = "custUpdate"
	typeCreateCustomer = "createCustomer"
	typeDeleteCustomer = "deleteCustomer"
	typeUpdateCustomer = "updateCustomer"
	typeCreateLicense  = "createLicense"
	errParseType       = "сервер не смог обработать запрос"
)

type CustMess struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func newCustomerMess(mType string, data map[string]interface{}) CustMess {
	var resp CustMess
	resp.Type = mType
	if data != nil {
		resp.Data = data
	} else {
		resp.Data = make(map[string]interface{})
	}
	return resp
}

//ErrorMessage структура ошибки
type ErrorMessage struct {
	Error string `json:"error"`
}
