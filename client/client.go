package main

import (
	"context"
	"log"
	"time"

	pb "golangrpc/proto" // Adjust this import path if necessary

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Updated to use grpc.DialContext
	//conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEmployeeServiceClient(conn)

	r, err := c.CreateEmployee(ctx, &pb.CreateEmployeeRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Position:  "Developer",
		Salary:    60000,
	})
	if err != nil {
		log.Fatalf("could not create employee: %v", err)
	}
	log.Printf("Employee created: %v", r.Employee)
}
