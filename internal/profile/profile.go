package profile

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/reviz0r/http-api/pkg/profile"
)

// UserService .
type UserService struct {
	pk    int64
	users []*profile.User
}

// Create .
func (s *UserService) Create(ctx context.Context, in *profile.CreateRequest) (*profile.CreateResponse, error) {
	if in.GetUser().GetName() != "" && in.GetUser().GetEmail() != "" {
		s.users = append(s.users, in.GetUser())
	} else {
		return nil, status.Error(codes.InvalidArgument, "user cannot be nil")
	}

	s.pk++
	return &profile.CreateResponse{Id: s.pk}, nil
}

// ReadAll .
func (s *UserService) ReadAll(ctx context.Context, in *profile.ReadAllRequest) (*profile.ReadAllResponse, error) {
	return &profile.ReadAllResponse{Users: s.users}, nil
}

// Read .
func (s *UserService) Read(ctx context.Context, in *profile.ReadRequest) (*profile.ReadResponse, error) {
	return &profile.ReadResponse{User: s.users[in.GetId()-1]}, nil
}

// Update .
func (s *UserService) Update(context.Context, *profile.UpdateRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}

// Replace .
func (s *UserService) Replace(context.Context, *profile.ReplaceRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}

// Delete .
func (s *UserService) Delete(context.Context, *profile.DeleteRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}
