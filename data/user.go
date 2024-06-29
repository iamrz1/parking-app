package data

import "supertal-tha-parking-app/model"

type UserStore interface {
	Register(req *model.UserCreateReq) (*model.UserCreateRes, error)
	Login(req *model.LoginReq) (*model.LoginRes, error)
	GetUser(username string) (*model.UserCreateRes, error)
}
