package profile_test

import (
	"context"
	"database/sql"
	"log"
	"net"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc"
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
		db, mock, _ = sqlmock.New()

		lis = bufconn.Listen(bufSize)
		s = grpc.NewServer()
		pkg.RegisterUserServiceServer(s, &internal.UserService{DB: db})
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("Server exited with error: %v", err)
			}
		}()
	})

	AfterSuite(func() {
		s.Stop()
		db.Close()
	})

	BeforeEach(func() {
		bufDialer := func(string, time.Duration) (net.Conn, error) {
			return lis.Dial()
		}

		var err error
		conn, err = grpc.DialContext(context.Background(), "bufnet",
			grpc.WithDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			panic("Failed to dial bufnet: " + err.Error())
		}
		client = pkg.NewUserServiceClient(conn)
	})

	AfterEach(func() {
		conn.Close()
	})

	Describe("Create", func() {
		It("can create user", func() {
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+) RETURNING "id"`).
				WithArgs("user", "user@example.com").
				WillReturnRows(rows)

			res, err := client.Create(context.Background(),
				&pkg.CreateRequest{User: &pkg.User{Name: "user", Email: "user@example.com"}})

			Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.GetId()).To(Equal(int64(1)))
		})
	})

	Describe("ReadAll", func() {
		It("can get all users", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			mock.ExpectQuery(`SELECT (.+) FROM "users" LIMIT 100`).WillReturnRows(rows)

			res, err := client.ReadAll(context.Background(), &pkg.ReadAllRequest{})

			Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.GetUsers()).To(HaveLen(1))
			Expect(res.GetUsers()).To(Equal([]*pkg.User{{Name: "user", Email: "user@example.com"}}))
		})
	})

	Describe("Read", func() {
		It("can get user by id", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			mock.ExpectQuery(`select (.+) from "users" where "id"=\$1`).
				WithArgs(1).
				WillReturnRows(rows)

			res, err := client.Read(context.Background(), &pkg.ReadRequest{Id: 1})

			Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.GetUser()).To(Equal(&pkg.User{Name: "user", Email: "user@example.com"}))
		})
	})

	Describe("Update", func() {
		It("can update user by id", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			mock.ExpectQuery(`select (.+) from "users" where "id"=\$1`).
				WithArgs(1).
				WillReturnRows(rows)
			mock.ExpectExec(`UPDATE "users" SET "name"=\$1,"email"=\$2 WHERE "id"=\$3`).
				WithArgs("user1", "user1@example.com", 1).
				WillReturnResult(sqlmock.NewResult(0, 1))

			res, err := client.Update(context.Background(), &pkg.UpdateRequest{
				Id:     1,
				User:   &pkg.User{Name: "user1", Email: "user1@example.com"},
				Fields: &field_mask.FieldMask{Paths: []string{"name", "email"}},
			})

			Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
		})
	})

	Describe("Delete", func() {
		It("can delete user by id", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user", "user@example.com")
			mock.ExpectQuery(`select (.+) from "users" where "id"=\$1`).
				WithArgs(1).
				WillReturnRows(rows)
			mock.ExpectExec(`DELETE FROM "users" WHERE "id"=\$1`).
				WithArgs(1).
				WillReturnResult(sqlmock.NewResult(0, 1))

			res, err := client.Delete(context.Background(), &pkg.DeleteRequest{Id: 1})

			Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
		})
	})
})
