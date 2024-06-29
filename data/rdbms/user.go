package data

import (
	"supertal-tha-parking-app/data"
	"supertal-tha-parking-app/model"
	"supertal-tha-parking-app/utils"

	"github.com/iamrz1/rutils"
	"gorm.io/gorm"
)

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) data.UserStore {
	return &userStore{db: db}
}

func (u *userStore) Register(req *model.UserCreateReq) (*model.UserCreateRes, error) {
	user := model.User{
		Name:              req.Name,
		Username:          req.Username,
		EncryptedPassword: utils.GetEncodedPassword(req.Password),
	}
	isParkingManger := false
	if req.IsParkingManager != nil && *req.IsParkingManager {
		isParkingManger = true
	}

	user.IsParkingManager = &isParkingManger

	err := u.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	response := model.UserCreateRes{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}

	return &response, nil
}

func (u *userStore) Login(req *model.LoginReq) (*model.LoginRes, error) {
	user := model.User{
		Username: req.Username,
	}

	err := u.db.Where(&user).First(&user).Error
	if err != nil {
		return nil, err
	}

	loggedIn := utils.VerifyPassword(req.Password, user.EncryptedPassword)
	if !loggedIn {
		return nil, rutils.NewValidationError("Invalid credentials", nil)
	}
	role := utils.RoleUser
	if *user.IsParkingManager {
		role = utils.RoleManger
	}

	access, refresh := utils.GenerateTokens(req.Username, role)
	response := model.LoginRes{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	return &response, nil
}

func (u *userStore) GetUser(username string) (*model.UserCreateRes, error) {
	var user model.User

	err := u.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	response := model.UserCreateRes{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}

	return &response, err
}
