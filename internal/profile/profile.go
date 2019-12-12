package profile

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/reviz0r/http-api/internal/profile/models"
	"github.com/reviz0r/http-api/pkg/profile"
)

// UserService .
type UserService struct {
	*sql.DB
}

// Create .
func (s *UserService) Create(ctx context.Context, in *profile.CreateRequest) (*profile.CreateResponse, error) {
	user := &models.User{
		Name:  in.GetUser().GetName(),
		Email: in.GetUser().GetEmail(),
	}

	err := user.Insert(ctx, s.DB, boil.Infer())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Create: %w", err.Error())
	}

	return &profile.CreateResponse{Id: user.ID}, nil
}

// ReadAll .
func (s *UserService) ReadAll(ctx context.Context, in *profile.ReadAllRequest) (*profile.ReadAllResponse, error) {
	users, err := models.Users(qm.Select(in.GetFields().GetPaths()...)).All(ctx, s.DB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.ReadAll: %w", err.Error())
	}

	pbUsers := make([]*profile.User, len(users))
	for i, user := range users {
		pbUsers[i] = &profile.User{Name: user.Name, Email: user.Email}
	}

	return &profile.ReadAllResponse{Users: pbUsers}, nil
}

// Read .
func (s *UserService) Read(ctx context.Context, in *profile.ReadRequest) (*profile.ReadResponse, error) {
	user, err := models.FindUser(ctx, s.DB, in.GetId(), in.GetFields().GetPaths()...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Read: %w", err.Error())
	}

	pbUser := &profile.User{
		Name:  user.Name,
		Email: user.Email,
	}

	return &profile.ReadResponse{User: pbUser}, nil
}

// Update .
func (s *UserService) Update(ctx context.Context, in *profile.UpdateRequest) (*empty.Empty, error) {
	log.Println(in.GetFields().GetPaths())
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}

// Replace .
func (s *UserService) Replace(context.Context, *profile.ReplaceRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}

// Delete .
func (s *UserService) Delete(ctx context.Context, in *profile.DeleteRequest) (*empty.Empty, error) {
	res, err := s.DB.ExecContext(ctx, "delete from users where id = $1", in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Delete: %w", err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Delete: %w", err.Error())
	}

	if rows != 1 {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}

	return new(empty.Empty), nil
}