package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	user_handler "integration-test/internal/handler/grpc"
	user_repo "integration-test/internal/repo/user"
	user_usecase "integration-test/internal/usecase/user"
	pb "integration-test/proto/integration-test/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserIntegrationSuite struct {
	suite.Suite
	// clientUserService pb.UserServiceClient
	grpcServer *grpc.Server
	dialer     func(context.Context, string) (net.Conn, error)
}

func (suite *UserIntegrationSuite) SetupSuite() {
	// TODO change this to testcontainers
	dbURL := "postgres://postgres:postgres@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)
	grpcServer := grpc.NewServer()
	userRepo := user_repo.New(db)
	userUsecase := user_usecase.New(userRepo)
	userHandler := user_handler.New(userUsecase)

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			fmt.Println("error", err)
			panic(err)
		}
	}()
	// suite.clientUserService = client
	suite.grpcServer = grpcServer
	suite.dialer = func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

}

func (suite *UserIntegrationSuite) TestExample() {

	input := &pb.User{
		Id:   "test-3",
		Name: "Indra",
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(suite.dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		suite.T().Fatal(err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	out, err := client.CreateUser(ctx, input)
	st, ok := status.FromError(err)
	fmt.Println(st.Code(), ok)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), input.Id, out.Id)
	assert.Equal(suite.T(), input.Name, out.Name)
}

func (suite *UserIntegrationSuite) TearDownAllSuite() {
	suite.grpcServer.Stop()
}

func Test_IntegrationExample(t *testing.T) {
	integrationTest := os.Getenv("ALLOW_INTEGRATION_TEST")
	if integrationTest == "true" {
		suite.Run(t, new(UserIntegrationSuite))
	} else {
		t.Skip("Skip Integration test Test_IntegrationExample")
	}
}
