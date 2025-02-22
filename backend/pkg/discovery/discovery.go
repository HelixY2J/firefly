package discovery

import (
	"context"
	"fmt"
	"math/rand"
)

type Registry interface {
	Register(ctx context.Context, instanceID, serviceName, hostport string) error
	Unregister(ctx context.Context, instanceID string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceID, serviceName string) error
}

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.Intn(100000))
}
