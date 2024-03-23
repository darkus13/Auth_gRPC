package converter

import (
	"github.com/darkus13/Auth_gRPC/internal/repository/auth/model"
	modelRepo "github.com/darkus13/Auth_gRPC/internal/repository/auth/model"
)

func ToAuthFromRepo(auth *modelRepo.Auth) *model.Auth {
	return &model.Auth{
		ID:        auth.ID,
		Info:      ToAuthInfoFromRepo(auth.Info),
		CreatedAt: auth.CreatedAt,
		UpdatedAt: auth.UpdatedAt,
	}
}
func ToAuthInfoFromRepo(info modelRepo.AuthInfo) model.AuthInfo {
	return model.AuthInfo{
		Name:        info.Name,
		Email:       info.Email,
		Password:    info.Password,
		PassConfirm: info.PassConfirm,
		Role:        info.Role,
	}
}
