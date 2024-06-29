package model

type UserCreateReq struct {
	Name             string
	Username         string
	Password         string
	IsParkingManager *bool
}

type UserCreateRes struct {
	ID       uint
	Username string
	Name     string
}

type LoginReq struct {
	Username string
	Password string
}

type LoginRes struct {
	AccessToken  string
	RefreshToken string
}
