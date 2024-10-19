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
	status := mapTaskStatusToUint64(req.Status)
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

func mapTaskStatusToUint64(status *pb.TaskStatus) *uint64 {
	if status == nil {
		return nil
	}

	var mappedValue uint64
	switch *status {
	case pb.TaskStatus_STATUS_UNSPECIFIED:
		mappedValue = 0
	case pb.TaskStatus_STATUS_STOPPED:
		mappedValue = 1
	case pb.TaskStatus_STATUS_ACTIVE:
		mappedValue = 2
	case pb.TaskStatus_STATUS_NOT_ACTIVE:
		mappedValue = 3
	}

	return &mappedValue
}
