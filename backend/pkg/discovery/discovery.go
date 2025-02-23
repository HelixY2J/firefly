package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(ctx context.Context, instanceID, serviceName, hostport string) error
	Unregister(ctx context.Context, instanceID string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceID, serviceName string) error
}

func GenerateInstanceID(serviceName string, port int) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s-%d-%d", serviceName, port, rand.Intn(100000))
}
