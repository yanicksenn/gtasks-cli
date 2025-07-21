package gtasks

import (
	"context"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

func newTestClient(serverURL string) (Client, error) {
	service, err := tasks.NewService(context.Background(),
		option.WithEndpoint(serverURL),
		option.WithoutAuthentication(),
		option.WithHTTPClient(&http.Client{}),
	)
	if err != nil {
		return nil, err
	}

	return &onlineClient{service: service}, nil
}