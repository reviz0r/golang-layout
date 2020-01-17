package profile_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	internal "github.com/reviz0r/golang-layout/internal/profile"
	pkg "github.com/reviz0r/golang-layout/pkg/profile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Profile Suite")
}

var _ = Describe("Profile", func() {
	const bufSize = 1024 * 1024

	var (
		// Mock DB
		db   *sql.DB
		mock sqlmock.Sqlmock

		// Server fake port
		lis *bufconn.Listener
		s   *grpc.Server

		// Client fake connection
		conn   *grpc.ClientConn
		client pkg.UserServiceClient
	)

	BeforeSuite(func() {
		// create mock db
		db, mock, _ = sqlmock.New()

		// create local grpc server
		lis = bufconn.Listen(bufSize)
		s = grpc.NewServer()
		pkg.RegisterUserServiceServer(s, &internal.UserService{DB: db})
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("Server exited with error: %v", err)
			}
		}()
		bufDialer := func(string, time.Duration) (net.Conn, error) {
			return lis.Dial()
		}

		// connect to local grpc server
		var err error
		conn, err = grpc.DialContext(context.Background(), "bufnet",
			grpc.WithDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to dial bufnet: %s", err.Error())
		}

		client = pkg.NewUserServiceClient(conn)
	})

	AfterSuite(func() {
		conn.Close()
		s.Stop()
		db.Close()
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		q := `^INSERT INTO "users" (.+) VALUES (.+) RETURNING "id"$`

		It("can create user", func() {
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(q).WithArgs("user", "user@example.com").WillReturnRows(rows)

			res, err := client.Create(context.Background(),
				&pkg.CreateRequest{User: &pkg.User{Name: "user", Email: "user@example.com"}})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.GetId()).To(Equal(int64(1)))
		})

		It("gives Internal error if cannot create user", func() {
			mock.ExpectQuery(q).WithArgs("user", "user@example.com").WillReturnError(errors.New("some error"))

			res, err := client.Create(context.Background(),
				&pkg.CreateRequest{User: &pkg.User{Name: "user", Email: "user@example.com"}})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.Internal))
			Expect(grpcStatus.Message()).To(Equal("UserService.Create: %!w(string=models: unable to insert into users: some error)"))
		})
	})

	Describe("ReadAll", func() {
		qCount := `^SELECT COUNT(.+) FROM "users";$`

		It("can get all users", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
			mock.ExpectQuery(`^SELECT (.+) FROM "users" LIMIT 100;$`).WillReturnRows(rows)
			mock.ExpectQuery(qCount).WillReturnRows(countRows)

			res, err := client.ReadAll(context.Background(), &pkg.ReadAllRequest{})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.GetUsers()).To(HaveLen(1))
			Expect(res.GetUsers()).To(Equal([]*pkg.User{{Name: "user", Email: "user@example.com"}}))
			Expect(res.GetLimit()).To(Equal(int32(100)))
			Expect(res.GetOffset()).To(Equal(int32(0)))
			Expect(res.GetTotal()).To(Equal(int32(1)))
		})

		It("gives Internal error if cannot get all users", func() {
			mock.ExpectQuery(`^SELECT (.+) FROM "users" LIMIT 100;$`).
				WillReturnError(errors.New("some error"))

			res, err := client.ReadAll(context.Background(), &pkg.ReadAllRequest{})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.Internal))
			Expect(grpcStatus.Message()).To(Equal("UserService.ReadAll: %!w(string=models: failed to assign all query results to User slice: bind failed to execute query: some error)"))
		})

		It("gives Internal error if cannot get users count", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			mock.ExpectQuery(`^SELECT (.+) FROM "users" LIMIT 100;$`).WillReturnRows(rows)
			mock.ExpectQuery(qCount).WillReturnError(errors.New("some error"))

			res, err := client.ReadAll(context.Background(), &pkg.ReadAllRequest{})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.Internal))
			Expect(grpcStatus.Message()).To(Equal("UserService.ReadAll: %!w(string=models: failed to count users rows: some error)"))
		})

		It("can get all users with limit 10", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
			mock.ExpectQuery(`^SELECT (.+) FROM "users" LIMIT 10;$`).WillReturnRows(rows)
			mock.ExpectQuery(qCount).WillReturnRows(countRows)

			res, err := client.ReadAll(context.Background(), &pkg.ReadAllRequest{Limit: 10})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.GetUsers()).To(HaveLen(1))
			Expect(res.GetUsers()).To(Equal([]*pkg.User{{Name: "user", Email: "user@example.com"}}))
			Expect(res.GetLimit()).To(Equal(int32(10)))
			Expect(res.GetOffset()).To(Equal(int32(0)))
			Expect(res.GetTotal()).To(Equal(int32(1)))
		})

		It("can get all users with limit 1000", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
			mock.ExpectQuery(`^SELECT (.+) FROM "users" LIMIT 1000;$`).WillReturnRows(rows)
			mock.ExpectQuery(qCount).WillReturnRows(countRows)

			res, err := client.ReadAll(context.Background(), &pkg.ReadAllRequest{Limit: 1000})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.GetUsers()).To(HaveLen(1))
			Expect(res.GetUsers()).To(Equal([]*pkg.User{{Name: "user", Email: "user@example.com"}}))
			Expect(res.GetLimit()).To(Equal(int32(1000)))
			Expect(res.GetOffset()).To(Equal(int32(0)))
			Expect(res.GetTotal()).To(Equal(int32(1)))
		})

		It("can get all users with limit 10000 (but really 1000)", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
			mock.ExpectQuery(`^SELECT (.+) FROM "users" LIMIT 1000;$`).WillReturnRows(rows)
			mock.ExpectQuery(qCount).WillReturnRows(countRows)

			res, err := client.ReadAll(context.Background(), &pkg.ReadAllRequest{Limit: 10000})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.GetUsers()).To(HaveLen(1))
			Expect(res.GetUsers()).To(Equal([]*pkg.User{{Name: "user", Email: "user@example.com"}}))
			Expect(res.GetLimit()).To(Equal(int32(1000)))
			Expect(res.GetOffset()).To(Equal(int32(0)))
			Expect(res.GetTotal()).To(Equal(int32(1)))
		})
	})

	Describe("Read", func() {
		q := `^select (.+) from "users" where "id"=\$1$`

		It("can get user by id", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			mock.ExpectQuery(q).WithArgs(1).WillReturnRows(rows)

			res, err := client.Read(context.Background(), &pkg.ReadRequest{Id: 1})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.GetUser()).To(Equal(&pkg.User{Name: "user", Email: "user@example.com"}))
		})

		It("gives Internal error if cannot get user", func() {
			mock.ExpectQuery(q).WithArgs(1).WillReturnError(errors.New("some error"))

			res, err := client.Read(context.Background(), &pkg.ReadRequest{Id: 1})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.Internal))
			Expect(grpcStatus.Message()).To(Equal("UserService.Read: %!w(string=models: unable to select from users: bind failed to execute query: some error)"))
		})

		It("gives NotFound error if user does not exist", func() {
			mock.ExpectQuery(q).WithArgs(2).WillReturnError(sql.ErrNoRows)

			res, err := client.Read(context.Background(), &pkg.ReadRequest{Id: 2})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.NotFound))
			Expect(grpcStatus.Message()).To(Equal(codes.NotFound.String()))
		})
	})

	Describe("Update", func() {
		q := `^UPDATE "users" SET "name"=\$1,"email"=\$2 WHERE "id"=\$3$`

		It("can update user by id", func() {
			mock.ExpectExec(q).WithArgs("user1", "user1@example.com", 1).WillReturnResult(sqlmock.NewResult(0, 1))

			res, err := client.Update(context.Background(), &pkg.UpdateRequest{
				Id:     1,
				User:   &pkg.User{Name: "user1", Email: "user1@example.com"},
				Fields: &field_mask.FieldMask{Paths: []string{"name", "email"}},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
		})

		It("gives InvalidArgument error if fields not specified", func() {
			res, err := client.Update(context.Background(), &pkg.UpdateRequest{
				Id:   1,
				User: &pkg.User{Name: "user1", Email: "user1@example.com"},
			})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.InvalidArgument))
			Expect(grpcStatus.Message()).To(Equal("UserService.Update: fields must be specified"))
		})

		It("gives Internal error if cannot update user", func() {
			mock.ExpectExec(q).WithArgs("user1", "user1@example.com", 1).WillReturnError(errors.New("some error"))

			res, err := client.Update(context.Background(), &pkg.UpdateRequest{
				Id:     1,
				User:   &pkg.User{Name: "user1", Email: "user1@example.com"},
				Fields: &field_mask.FieldMask{Paths: []string{"name", "email"}},
			})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.Internal))
			Expect(grpcStatus.Message()).To(Equal("UserService.Update: %!w(string=models: unable to update users row: some error)"))
		})

		It("gives NotFound error if updated 0 rows", func() {
			mock.ExpectExec(q).WithArgs("user1", "user1@example.com", 1).WillReturnResult(sqlmock.NewResult(0, 0))

			res, err := client.Update(context.Background(), &pkg.UpdateRequest{
				Id:     1,
				User:   &pkg.User{Name: "user1", Email: "user1@example.com"},
				Fields: &field_mask.FieldMask{Paths: []string{"name", "email"}},
			})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.NotFound))
			Expect(grpcStatus.Message()).To(Equal(codes.NotFound.String()))
		})

		It("gives Internal error if updated more than 1 row", func() {
			mock.ExpectExec(q).WithArgs("user1", "user1@example.com", 1).WillReturnResult(sqlmock.NewResult(0, 2))

			res, err := client.Update(context.Background(), &pkg.UpdateRequest{
				Id:     1,
				User:   &pkg.User{Name: "user1", Email: "user1@example.com"},
				Fields: &field_mask.FieldMask{Paths: []string{"name", "email"}},
			})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.Internal))
			Expect(grpcStatus.Message()).To(Equal("UserService.Update: expect updating 1 row, but updated 2 rows"))
		})
	})

	Describe("Delete", func() {
		q := `^DELETE FROM "users" WHERE "id"=\$1$`

		It("can delete user by id", func() {
			mock.ExpectExec(q).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

			res, err := client.Delete(context.Background(), &pkg.DeleteRequest{Id: 1})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
		})

		It("gives Internal error if cannot delete user", func() {
			mock.ExpectExec(q).WithArgs(1).WillReturnError(errors.New("some error"))

			res, err := client.Delete(context.Background(), &pkg.DeleteRequest{Id: 1})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.Internal))
			Expect(grpcStatus.Message()).To(Equal("UserService.Delete: %!w(string=models: unable to delete from users: some error)"))
		})

		It("gives NotFound error if updated 0 rows", func() {
			mock.ExpectExec(q).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))

			res, err := client.Delete(context.Background(), &pkg.DeleteRequest{Id: 1})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.NotFound))
			Expect(grpcStatus.Message()).To(Equal(codes.NotFound.String()))
		})

		It("gives Internal error if updated more than 1 row", func() {
			mock.ExpectExec(q).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 2))

			res, err := client.Delete(context.Background(), &pkg.DeleteRequest{Id: 1})

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.Internal))
			Expect(grpcStatus.Message()).To(Equal("UserService.Delete: expect deleting 1 row, but deleted 2 rows"))
		})
	})
})
