package app

// controller ajax返回
type Api struct {
	Ok   bool
	Code int
	Msg  string
	Id   string
	List interface{}
	Item interface{}
	Next NextJson
}

type NextJson struct {
	Name string
	Link string
}

func NewResp() Api {
	return Api{Ok: false}
}

func NewOk() Api {
	return Api{Ok: true}
}

func NewOkMsg(message string) Api {
	return Api{Ok: true, Msg: message}
}
