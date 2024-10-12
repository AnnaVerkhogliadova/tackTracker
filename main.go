package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	task_controller "taskTracker/controller"
	"taskTracker/database"
	"taskTracker/driver"
	pb "taskTracker/task-tracker/tasktracker"
	task_handler "taskTracker/trasnsport/task"
)

func main() {
	rw := database.Connect()

	taskDriver, err := driver.NewDbDriver(rw, nil)
	if err != nil {
		log.Fatalf("Failed to create db driver: %v", err)
	}

	taskCtrl, err := task_controller.NewController(taskDriver)
	if err != nil {
		log.Fatalf("Failed to create task controller: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTaskServiceServer(s, &task_handler.Handler{
		Controller: taskCtrl,
	})

	fmt.Println("gRPC server listening on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
