package profile

import (
	"github.com/reviz0r/golang-layout/internal/profile/models"
	"github.com/reviz0r/golang-layout/pkg/profile"
)

func userFromProto(in *profile.User) *models.User {
	if in == nil {
		return nil
	}

	return &models.User{
		Name:  in.GetName(),
		Email: in.GetEmail(),
	}
}

func userToProto(in *models.User) *profile.User {
	if in == nil {
		return nil
	}

	return &profile.User{
		Name:  in.Name,
		Email: in.Email,
	}
}
