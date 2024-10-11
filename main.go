package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	task_controller "taskTracker/controller"
	"taskTracker/database"
	"taskTracker/driver"
	"taskTracker/model"
	pb "taskTracker/task-tracker/tasktracker"
)

type handler struct {
	controller task_controller.ITaskController
	pb.UnimplementedTaskServiceServer
}

func (h *handler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      uint64(req.Status),
	}

	taskId, err := h.controller.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTaskResponse{
		TaskId: taskId,
	}, nil
}

func (h *handler) GetTasks(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	task, err := h.controller.Get(ctx, req.TaskId)

	if err != nil {
		return nil, err
	}

	var status pb.TaskStatus
	if task.Status != 0 {
		status = pb.TaskStatus(task.Status)
	}

	var createdDate *timestamppb.Timestamp
	if &task.CreatedAt != nil {
		createdDate = timestamppb.New(task.CreatedAt)
	}

	pbTask := pb.Task{
		TaskId:      task.ID,
		Title:       task.Title,
		Status:      status,
		Description: task.Description,
		CreateDate:  createdDate,
	}

	return &pb.GetResponse{
		Tasks: &pbTask,
	}, nil
}

func (h *handler) DeleteTask(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := h.controller.Delete(ctx, req.TaskId)

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

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

	pb.RegisterTaskServiceServer(s, &handler{
		controller: taskCtrl,
	})

	fmt.Println("gRPC server listening on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
