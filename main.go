package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	task_controller "taskTracker/controller"
	"taskTracker/database"
	"taskTracker/driver"
	pb "taskTracker/task-tracker/tasktracker"
	task_handler "taskTracker/trasnsport/task"
)

func main() {
	rw := database.Connect()
	log.Info().Msg("connection to database")

	taskDriver, err := driver.NewDbDriver(rw, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create db driver")
	}

	taskCtrl, err := task_controller.NewController(taskDriver)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create task controller")
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	s := grpc.NewServer()

	pb.RegisterTaskServiceServer(s, &task_handler.Handler{
		Controller: taskCtrl,
	})

	fmt.Println("gRPC server listening on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
