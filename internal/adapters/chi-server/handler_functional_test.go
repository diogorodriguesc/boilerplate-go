//go:build functional

package chiserver

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMyIntegration(t *testing.T) {
	ctx := context.Background()

	// 1. Define the container request
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	// 2. Start the container
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start redis: %s", err)
	}

	// 3. Clean up when the test finishes
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate container: %s", err)
		}
	}()

	// 4. Get the dynamic host and port
	endpoint, _ := redisC.Endpoint(ctx, "")

	t.Logf("Redis is running at: %s", endpoint)
	// Run your actual logic here...
}
