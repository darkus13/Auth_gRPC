package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/darkus13/Auth_gRPC/internal/repository/auth/model"
	decs "github.com/darkus13/Auth_gRPC/pkg/user_api_v1"
)

func ToUserFromService(auth *model.Auth) *decs.User {
	var updatedAt *timestamppb.Timestamp
	if auth.UpdatedAt.Valid {
		updatedAt = timestamppb.New(auth.UpdatedAt.Time)
	}

	infoPtr := &auth.Info

	return &decs.User{
		Id:        auth.ID,
		Info:      ToUserInfoFromServices(infoPtr),
		CreatedAt: timestamppb.New(auth.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromServices(auth *model.AuthInfo) *decs.UserInfo {
	return &decs.UserInfo{
		Name:     auth.Name,
		Email:    auth.Email,
		Password: auth.Password,
		Role:     decs.Role(auth.Role),
	}
}

func ToServiceFromAuth(auth *decs.CreateRequest) *model.AuthInfo {
	return &model.AuthInfo{
		Name:        auth.Name,
		Email:       auth.Email,
		Password:    auth.Password,
		PassConfirm: auth.PasswordConfirm,
		Role:        int32(auth.Role), // Преобразование enum в строку
	}
}
