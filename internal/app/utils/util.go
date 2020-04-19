package utils

type Response struct {
	Code int
	Obj  map[string]interface{}
}

//Message создает map для ответа пользователю
func Message(code int, message string) Response {
	return Response{Code: code, Obj: map[string]interface{}{"message": message}}
}
