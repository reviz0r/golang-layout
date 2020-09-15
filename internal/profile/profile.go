package profile

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/reviz0r/golang-layout/internal/profile/models"
	"github.com/reviz0r/golang-layout/pkg/profile"
)

// Module register user service in DI container
var Module = fx.Invoke(RegisterUserService)

// UserService .
type UserService struct {
	*sql.DB
}

// RegisterUserService .
func RegisterUserService(s *grpc.Server, db *sql.DB) {
	profile.RegisterUserServiceServer(s, &UserService{DB: db})
}

// Create .
func (s *UserService) Create(ctx context.Context, in *profile.CreateRequest) (*profile.CreateResponse, error) {
	user := userFromProto(in.GetUser())
	user.ID = 0

	err := user.Insert(ctx, s.DB, boil.Infer())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Create: %s", err.Error())
	}

	return &profile.CreateResponse{Id: user.ID}, nil
}

// ReadAll .
func (s *UserService) ReadAll(ctx context.Context, in *profile.ReadAllRequest) (*profile.ReadAllResponse, error) {
	// add fields to span
	span := opentracing.SpanFromContext(ctx).
		// SetOperationName("operation ReadAll"). // it's not needed - grpc already set operation name
		SetTag("blablabla", 123).      // I don't see it in span log
		SetBaggageItem("ping", "pong") // this works fine

	// just add some logging to tracer
	span.LogFields(log.String("field", "some_value"), log.Int("some_int", 1))
	span.LogKV("field_1", "value_1", "field_2", "value_2")

	// add fields to request logger
	ctxlogrus.AddFields(ctx, logrus.Fields{"my.custom.field": "some_value"})

	// Extract a single request-scoped logrus.Logger and log messages.
	l := ctxlogrus.Extract(ctx)

	var offset = in.GetOffset()
	var limit = in.GetLimit()
	if in.GetLimit() == 0 {
		limit = 100
	}
	if in.GetLimit() > 1000 {
		limit = 1000
	}

	l.Trace("selecting users")
	users, err := models.Users(
		qm.Select(in.GetFields().GetPaths()...),
		qm.Limit(int(limit)),
		qm.Offset(int(offset)),
	).All(ctx, s.DB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.ReadAll: %s", err.Error())
	}

	l.Trace("selecting total count")
	total, err := models.Users().Count(ctx, s.DB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.ReadAll: %s", err.Error())
	}

	l.Trace("marshal users to protobuf")
	pbUsers := make([]*profile.User, len(users))
	for i, user := range users {
		pbUsers[i] = userToProto(user)
	}

	l.Trace("return response")
	return &profile.ReadAllResponse{Users: pbUsers, Limit: limit, Offset: offset, Total: int32(total)}, nil
}

// Read .
func (s *UserService) Read(ctx context.Context, in *profile.ReadRequest) (*profile.ReadResponse, error) {
	user, err := models.FindUser(ctx, s.DB, in.GetId(), in.GetFields().GetPaths()...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Read: %s", err.Error())
	}

	pbUser := userToProto(user)

	return &profile.ReadResponse{User: pbUser}, nil
}

// Update .
func (s *UserService) Update(ctx context.Context, in *profile.UpdateRequest) (*empty.Empty, error) {
	if len(in.GetFields().GetPaths()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "UserService.Update: fields must be specified")
	}

	user := userFromProto(in.GetUser())
	user.ID = in.GetId()

	rows, err := user.Update(ctx, s.DB, boil.Whitelist(in.GetFields().GetPaths()...))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Update: %s", err.Error())
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
	user := &models.User{ID: in.GetId()}

	rows, err := user.Delete(ctx, s.DB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UserService.Delete: %s", err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}
	if rows > 1 {
		return nil, status.Errorf(codes.Internal, "UserService.Delete: expect deleting 1 row, but deleted %d rows", rows)
	}

	return new(empty.Empty), nil
}
