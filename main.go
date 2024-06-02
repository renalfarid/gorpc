package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	pb "golangrpc/proto"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedEmployeeServiceServer
	db *sql.DB
}

func (s *server) CreateEmployee(ctx context.Context, req *pb.CreateEmployeeRequest) (*pb.CreateEmployeeResponse, error) {
	// Insert employee into the SQLite database
	stmt, err := s.db.Prepare("INSERT INTO employees(first_name, last_name, email, position, salary) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(req.FirstName, req.LastName, req.Email, req.Position, req.Salary)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &pb.CreateEmployeeResponse{
		Employee: &pb.Employee{
			Id:        int32(id),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Position:  req.Position,
			Salary:    req.Salary,
		},
	}, nil
}

func main() {
	db, err := sql.Open("sqlite3", "./employees.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Create the employees table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS employees (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        first_name TEXT,
        last_name TEXT,
        email TEXT,
        position TEXT,
        salary REAL
    )`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterEmployeeServiceServer(s, &server{db: db})
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
