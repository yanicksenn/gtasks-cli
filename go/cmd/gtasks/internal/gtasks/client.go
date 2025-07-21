package gtasks

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/auth"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/config"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// Client is the interface for interacting with the tasks service, abstracting
// away the online and offline implementations.
type Client interface {
	ListTaskLists() (*tasks.TaskLists, error)
	CreateTaskList(opts CreateTaskListOptions) (*tasks.TaskList, error)
	GetTaskList(opts GetTaskListOptions) (*tasks.TaskList, error)
	UpdateTaskList(opts UpdateTaskListOptions) (*tasks.TaskList, error)
	DeleteTaskList(opts DeleteTaskListOptions) error

	ListTasks(opts ListTasksOptions) (*tasks.Tasks, error)
	CreateTask(opts CreateTaskOptions) (*tasks.Task, error)
	GetTask(opts GetTaskOptions) (*tasks.Task, error)
	UpdateTask(opts UpdateTaskOptions) (*tasks.Task, error)
	CompleteTask(opts CompleteTaskOptions) (*tasks.Task, error)
	UncompleteTask(opts UncompleteTaskOptions) (*tasks.Task, error)
	DeleteTask(opts DeleteTaskOptions) error
}

// onlineClient is a client that interacts with the real Google Tasks API.
type onlineClient struct {
	service *tasks.Service
}

// NewClient creates a new client based on the --offline flag.
func NewClient(cmd *cobra.Command, ctx context.Context) (Client, error) {
	offline, _ := cmd.Flags().GetBool("offline")
	if offline {
		return newOfflineClient()
	}
	return newOnlineClient(ctx)
}

func newOnlineClient(ctx context.Context) (*onlineClient, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		return nil, err
	}

	httpClient, err := authenticator.GetClient(ctx, cfg.ActiveAccount)
	if err != nil {
		return nil, err
	}

	service, err := tasks.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &onlineClient{service: service}, nil
}
