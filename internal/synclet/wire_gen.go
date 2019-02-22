// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package synclet

import (
	"context"
	"github.com/windmilleng/tilt/internal/container"
	"github.com/windmilleng/tilt/internal/docker"
	"github.com/windmilleng/tilt/internal/k8s"
)

// Injectors from wire.go:

func WireSynclet(ctx context.Context, env k8s.Env, runtime container.Runtime) (*Synclet, error) {
	dockerEnv, err := docker.ProvideEnv(ctx, env, runtime)
	if err != nil {
		return nil, err
	}
	client, err := docker.ProvideDockerClient(ctx, dockerEnv)
	if err != nil {
		return nil, err
	}
	version, err := docker.ProvideDockerVersion(ctx, client)
	if err != nil {
		return nil, err
	}
	cli, err := docker.DefaultClient(ctx, client, version)
	if err != nil {
		return nil, err
	}
	synclet := NewSynclet(cli)
	return synclet, nil
}
