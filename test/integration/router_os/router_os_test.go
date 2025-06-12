package router_os

import (
	"context"
	"fmt"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/router"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func TestRouterOS(t *testing.T) {
	ctx := context.Background()
	routerOSContainer, errGenericContainer := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "evilfreelancer/docker-routeros:7.17",
			ExposedPorts: []string{"8728/tcp"},
			WaitingFor:   wait.ForExposedPort(),
			//Cmd:          []string{"/proto/api.proto"},
			//HostConfigModifier: func(config *container.HostConfig) {
			//	config.AutoRemove = true
			//	config.Binds = []string{currPath + "/api_proto:/proto"}
			//},
		},
		Started: true,
	})
	require.NoError(t, errGenericContainer, "error starting router OS container")
	testcontainers.CleanupContainer(t, routerOSContainer)

	tracker, err := router.NewRouterTracker()
	if err != nil {
		t.Fatal(err)
	}

	macAddresses, errGet := tracker.GetMacAddresses(ctx, "127.0.0.1:8729", "admin", "")
	if errGet != nil {
		t.Fatal(errGet)
	}
	fmt.Println("MAC Addresses: ", macAddresses)
}
