package router_os

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/router_tracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"io"
	"testing"
	"time"
)

func TestRouterOS(t *testing.T) {
	ctx := context.Background()
	routerOSContainer, errGenericContainer := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "evilfreelancer/docker-routeros:7.17",
			ExposedPorts: []string{"8728/tcp"},
			HostConfigModifier: func(config *container.HostConfig) {
				config.CapAdd = []string{"NET_ADMIN"}
				config.Privileged = true
				config.Devices = []container.DeviceMapping{
					{PathOnHost: "/dev/net/tun", PathInContainer: "/dev/net/tun"},
					{PathOnHost: "/dev/kvm", PathInContainer: "/dev/kvm"},
				}
			},
		},
		Started: true,
	})
	require.NoError(t, errGenericContainer, "error starting router OS container")

	state, errGetState := routerOSContainer.State(ctx)
	require.NoError(t, errGetState)
	t.Log("Router OS state running:", state.Running)

	time.Sleep(10 * time.Second)

	logs, err := routerOSContainer.Logs(ctx)
	if err != nil {
		t.Fatalf("Failed to get container logs: %v", err)
	}
	defer logs.Close()

	bytes, err := io.ReadAll(logs)
	if err != nil {
		t.Fatalf("Failed to read container logs: %v", err)
	}

	fmt.Printf("Container logs:\n%s", string(bytes))

	ports, errPorts := routerOSContainer.Ports(ctx)
	require.NoError(t, errPorts)
	fmt.Println("Ports:", ports)

	endpoint, errPortEndpoint := routerOSContainer.PortEndpoint(ctx, "8728/tcp", "")
	require.NoError(t, errPortEndpoint)
	t.Log("Endpoint:", endpoint)

	tracker, err := router_tracker.NewRouterTracker()
	if err != nil {
		t.Fatal(err)
	}

	macAddresses, errGet := tracker.GetMacAddresses(ctx, endpoint, "admin", "")
	if errGet != nil {
		t.Fatal(errGet)
	}
	fmt.Println("MAC Addresses: ", macAddresses)

	assert.Len(t, macAddresses, 0)

	testcontainers.CleanupContainer(t, routerOSContainer)
}
