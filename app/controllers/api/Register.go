package api

type RigisterRequest struct {
	Name     string `json: "name"`
	Password string `json: "pw"`
	Id       int64  `json: "uid"`
}

type RigisterResponse struct {
	Name string `json: "name"`
	Id   int64  `json: "uid"`
}
