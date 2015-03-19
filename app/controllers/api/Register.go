package api

type RigisterRequest struct {
	Name     string `json: "name"`
	Password string `json: "pw"`
	Id       int64  `json: "uid"`
}

type RigisterResponse struct {
	Id            int64  `json: "uid"`
	Username      string `json: "username"`
	EmailProvider string `json:"ep"`
	Email         string `json:"email"`
}
