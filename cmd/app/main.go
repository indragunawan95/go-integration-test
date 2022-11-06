package main

import (
	"log"
	"net"

	user_handler "integration-test/internal/handler/grpc"
	user_repo "integration-test/internal/repo/user"
	user_usecase "integration-test/internal/usecase/user"
	pb "integration-test/proto/integration-test/user"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbURL := "postgres://postgres:postgres@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	userRepo := user_repo.New(db)
	userUsecase := user_usecase.New(userRepo)
	userHandler := user_handler.New(userUsecase)

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to start grpc")
	}
}
