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

func (h *Handler) SetStatus(ctx context.Context, req *pb.SetStatusRequest) (*emptypb.Empty, error) {
	var status *uint64
	if req.Status != 0 {
		ft := uint64(req.Status.Number())
		status = &ft
	}

	err := h.Controller.SetStatus(ctx, req.TaskId, status)

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
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

func (h *Handler) GetListTasks(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	var status *uint64
	if req.Status != nil {
		ft := uint64(req.Status.Number())
		status = &ft
	}

	tasks, err := h.Controller.GetList(ctx, status)

	if err != nil {
		return nil, error_handler.HandleError(err)
	}

	result := make([]*pb.Task, len(tasks))

	for i, task := range tasks {
		var createdDate *timestamppb.Timestamp
		if &task.CreatedAt != nil {
			createdDate = timestamppb.New(task.CreatedAt)
		}

		var statusReq pb.TaskStatus
		if task.Status != 0 {
			statusReq = pb.TaskStatus(task.Status)
		}

		result[i] = &pb.Task{
			TaskId:      task.ID,
			Title:       task.Title,
			Status:      statusReq,
			Description: task.Description,
			CreateDate:  createdDate,
		}
	}

	return &pb.GetListResponse{
		Tasks: result,
	}, nil
}
