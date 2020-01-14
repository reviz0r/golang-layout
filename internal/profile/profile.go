package profile

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/reviz0r/golang-layout/internal/profile/models"
	"github.com/reviz0r/golang-layout/pkg/profile"
)

// UserService .
type UserService struct {
	*sql.DB
}

// Create .
func (s *UserService) Create(ctx context.Context, in *profile.CreateRequest) (*profile.CreateResponse, error) {
	user := userFromProto(in.GetUser())

	err := user.Insert(ctx, s.DB, boil.Infer())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Create: %w", err.Error())
	}

	return &profile.CreateResponse{Id: user.ID}, nil
}

// ReadAll .
func (s *UserService) ReadAll(ctx context.Context, in *profile.ReadAllRequest) (*profile.ReadAllResponse, error) {
	var offset = in.GetOffset()
	var limit = in.GetLimit()
	if in.GetLimit() == 0 {
		limit = 100
	}
	if in.GetLimit() > 1000 {
		limit = 1000
	}

	users, err := models.Users(
		qm.Select(in.GetFields().GetPaths()...),
		qm.Limit(int(limit)),
		qm.Offset(int(offset)),
	).All(ctx, s.DB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.ReadAll: %w", err.Error())
	}

	total, err := models.Users().Count(ctx, s.DB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.ReadAll: %w", err.Error())
	}

	pbUsers := make([]*profile.User, len(users))
	for i, user := range users {
		pbUsers[i] = userToProto(user)
	}

	return &profile.ReadAllResponse{Users: pbUsers, Limit: limit, Offset: offset, Total: int32(total)}, nil
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

	pbUser := userToProto(user)

	return &profile.ReadResponse{User: pbUser}, nil
}

// Update .
func (s *UserService) Update(ctx context.Context, in *profile.UpdateRequest) (*empty.Empty, error) {
	if len(in.GetFields().GetPaths()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "UserService.Update: fields must be specified")
	}

	user, err := models.FindUser(ctx, s.DB, in.GetId(), models.UserColumns.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Update: %w", err.Error())
	}

	{
		user.Name = in.GetUser().GetName()
		user.Email = in.GetUser().GetEmail()
	}

	rows, err := user.Update(ctx, s.DB, boil.Whitelist(in.GetFields().GetPaths()...))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Update: %w", err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}
	if rows > 1 {
		return nil, status.Errorf(codes.Internal, "UserService.Update: expect updating 1 row, but updated %d rows", rows)
	}

	return new(empty.Empty), nil
}

// Delete .
func (s *UserService) Delete(ctx context.Context, in *profile.DeleteRequest) (*empty.Empty, error) {
	user, err := models.FindUser(ctx, s.DB, in.GetId(), models.UserColumns.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Delete: %w", err.Error())
	}

	rows, err := user.Delete(ctx, s.DB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Delete: %w", err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}
	if rows > 1 {
		return nil, status.Errorf(codes.Internal, "UserService.Delete: expect deleting 1 row, but deleted %d rows", rows)
	}

	return new(empty.Empty), nil
}
