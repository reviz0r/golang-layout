package profile

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/reviz0r/http-api/pkg/profile"
)

// UserService .
type UserService struct {
	*sql.DB
}

// Create .
func (s *UserService) Create(ctx context.Context, in *profile.CreateRequest) (*profile.CreateResponse, error) {
	var id int64

	err := s.DB.QueryRowContext(ctx, "insert into users (name, email) values ($1, $2) returning id", in.GetUser().GetName(), in.GetUser().GetEmail()).Scan(&id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Create: %w", err.Error())
	}

	return &profile.CreateResponse{Id: id}, nil
}

// ReadAll .
func (s *UserService) ReadAll(ctx context.Context, in *profile.ReadAllRequest) (*profile.ReadAllResponse, error) {
	log.Println(in.GetFields().GetPaths())

	rows, err := s.DB.QueryContext(ctx, "select name, email from users limit 100")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.ReadAll: %w", err.Error())
	}

	var users []*profile.User

	for rows.Next() {
		var user profile.User

		err := rows.Scan(&user.Name, &user.Email)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "UserService.ReadAll: %w", err.Error())
		}

		users = append(users, &user)
	}

	return &profile.ReadAllResponse{Users: users}, nil
}

// Read .
func (s *UserService) Read(ctx context.Context, in *profile.ReadRequest) (*profile.ReadResponse, error) {
	log.Println(in.GetFields().GetPaths())

	var user profile.User

	err := s.DB.QueryRowContext(ctx, "select name, email from users where id = $1", in.GetId()).Scan(&user.Name, &user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Read: %w", err.Error())
	}

	return &profile.ReadResponse{User: &user}, nil
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
