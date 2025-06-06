package schemas

type Auth struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}
type Config struct {
	Username string `json:"username" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

func (a Auth) Read(p []byte) (n int, err error) {
	panic("unimplemented")
}
