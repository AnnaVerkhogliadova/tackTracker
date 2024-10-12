package task

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	task_controller "taskTracker/controller"
	"taskTracker/model"
	pb "taskTracker/task-tracker/tasktracker"
	"taskTracker/trasnsport/error_handler"
)

type Handler struct {
	Controller task_controller.ITaskController
	pb.UnimplementedTaskServiceServer
}

//
//func NewHandler(controller task_controller.ITaskController) *Handler {
//	return &Handler{
//		Controller: controller,
//	}
//}

func (h *Handler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      uint64(req.Status),
	}

	taskId, err := h.Controller.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTaskResponse{
		TaskId: taskId,
	}, nil
}

func (h *Handler) GetTasks(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	task, err := h.Controller.Get(ctx, req.TaskId)

	if err != nil {
		return nil, error_handler.HandleError(err)
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

func (h *Handler) DeleteTask(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := h.Controller.Delete(ctx, req.TaskId)

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
