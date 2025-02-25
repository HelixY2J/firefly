package consul

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
)

type Registry struct {
	client *api.Client
}

func NewRegistry(addr string) (*Registry, error) {
	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client}, nil
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostport string) error {
	parts := strings.Split(hostport, ":")
	if len(parts) != 2 {
		return errors.New("invalid host:port format")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	host := parts[0]
	log.Printf("Registering service %s with Consul at %s:%d", serviceName, host, port)

	// Get hostname for metadata
	hostname, _ := os.Hostname()

	err = r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: host,
		Port:    port,
		Tags:    []string{"v1", "firefly", hostname},
		Meta: map[string]string{
			"hostname": hostname,
			"version":  "1.0",
		},
		Check: &api.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "20s",
		},
	})

	if err != nil {
		log.Printf(" Failed to register service in Consul: %v", err)
		return err
	}

	log.Printf("Successfully registered service %s with Consul", serviceName)
	return nil

}

func (r *Registry) Unregister(ctx context.Context, instanceID string) error {
	log.Printf("Unregistering service %s", instanceID)
	return r.client.Agent().ServiceDeregister(instanceID)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	// Use QueryOptions to enable network-wide search
	qOpts := &api.QueryOptions{
		AllowStale: true,  // Allow reading from non-leader nodes
		WaitTime:   time.Second * 2,  // Wait up to 2 seconds for consistent result
	}

	entries, _, err := r.client.Health().Service(serviceName, "", true, qOpts)
	if err != nil {
		return nil, err
	}

	var instances []string
	for _, entry := range entries {
		// Only include healthy instances
		if entry.Checks.AggregatedStatus() == api.HealthPassing {
			instances = append(instances, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
		}
	}
	return instances, nil
}

func (r *Registry) HealthCheck(instanceID string, servicename string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}
