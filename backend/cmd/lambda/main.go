package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/supporttickr/backend/internal/config"
	"github.com/supporttickr/backend/internal/routes"
	"github.com/supporttickr/backend/internal/store"
)

var adapter *httpadapter.HandlerAdapterV2

func init() {
	log.Println("Lambda cold start: initializing SupportDesk API...")

	cfg := config.Load()

	ctx := context.Background()
	st, err := store.NewStore(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	handler := routes.Setup(st, cfg)
	adapter = httpadapter.NewV2(handler)

	log.Println("Lambda initialization complete")
}

func handleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// API Gateway HTTP API includes the stage (e.g. "prod") as a path prefix when using named stages.
	// Our routes expect /api/... not /prod/api/..., so strip the stage prefix.
	stage := req.RequestContext.Stage
	if stage != "" && stage != "$default" {
		prefix := "/" + stage
		if strings.HasPrefix(req.RawPath, prefix) {
			req.RawPath = req.RawPath[len(prefix):]
			if req.RawPath == "" {
				req.RawPath = "/"
			}
		}
	}
	return adapter.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handleRequest)
}
