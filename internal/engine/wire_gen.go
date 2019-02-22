// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package engine

import (
	"context"
	"github.com/google/wire"
	"github.com/windmilleng/tilt/internal/build"
	"github.com/windmilleng/tilt/internal/docker"
	"github.com/windmilleng/tilt/internal/dockercompose"
	"github.com/windmilleng/tilt/internal/dockerfile"
	"github.com/windmilleng/tilt/internal/k8s"
	"github.com/windmilleng/tilt/internal/output"
	"github.com/windmilleng/tilt/internal/synclet"
	"github.com/windmilleng/wmclient/pkg/analytics"
	"github.com/windmilleng/wmclient/pkg/dirs"
)

// Injectors from wire.go:

func provideBuildAndDeployer(ctx context.Context, docker2 docker.Client, kClient k8s.Client, dir *dirs.WindmillDir, env k8s.Env, updateMode UpdateModeFlag, sCli synclet.SyncletClient, dcc dockercompose.DockerComposeClient) (BuildAndDeployer, error) {
	syncletManager := NewSyncletManagerForTests(kClient, sCli)
	syncletBuildAndDeployer := NewSyncletBuildAndDeployer(syncletManager, kClient)
	containerUpdater := build.NewContainerUpdater(docker2)
	memoryAnalytics := analytics.NewMemoryAnalytics()
	localContainerBuildAndDeployer := NewLocalContainerBuildAndDeployer(containerUpdater, memoryAnalytics)
	console := build.DefaultConsole()
	writer := output.CaptureAll()
	labels := _wireLabelsValue
	dockerImageBuilder := build.NewDockerImageBuilder(docker2, console, writer, labels)
	imageBuilder := build.DefaultImageBuilder(dockerImageBuilder)
	cacheBuilder := build.NewCacheBuilder(docker2)
	runtime := k8s.ProvideContainerRuntime(ctx, kClient)
	dockerEnv, err := docker.ProvideEnv(ctx, env, runtime)
	if err != nil {
		return nil, err
	}
	clock := build.ProvideClock()
	execCustomBuilder := build.NewExecCustomBuilder(docker2, dockerEnv, clock)
	engineUpdateMode, err := ProvideUpdateMode(updateMode, env, runtime)
	if err != nil {
		return nil, err
	}
	imageBuildAndDeployer := NewImageBuildAndDeployer(imageBuilder, cacheBuilder, execCustomBuilder, kClient, env, memoryAnalytics, engineUpdateMode, clock, runtime)
	engineImageAndCacheBuilder := NewImageAndCacheBuilder(imageBuilder, cacheBuilder, execCustomBuilder, engineUpdateMode)
	dockerComposeBuildAndDeployer := NewDockerComposeBuildAndDeployer(dcc, docker2, engineImageAndCacheBuilder, clock)
	buildOrder := DefaultBuildOrder(syncletBuildAndDeployer, localContainerBuildAndDeployer, imageBuildAndDeployer, dockerComposeBuildAndDeployer, env, engineUpdateMode, runtime)
	compositeBuildAndDeployer := NewCompositeBuildAndDeployer(buildOrder)
	return compositeBuildAndDeployer, nil
}

var (
	_wireLabelsValue = dockerfile.Labels{}
)

func provideImageBuildAndDeployer(ctx context.Context, docker2 docker.Client, kClient k8s.Client, dir *dirs.WindmillDir) (*ImageBuildAndDeployer, error) {
	console := build.DefaultConsole()
	writer := output.CaptureAll()
	labels := _wireLabelsValue
	dockerImageBuilder := build.NewDockerImageBuilder(docker2, console, writer, labels)
	imageBuilder := build.DefaultImageBuilder(dockerImageBuilder)
	cacheBuilder := build.NewCacheBuilder(docker2)
	env := _wireEnvValue
	runtime := k8s.ProvideContainerRuntime(ctx, kClient)
	dockerEnv, err := docker.ProvideEnv(ctx, env, runtime)
	if err != nil {
		return nil, err
	}
	clock := build.ProvideClock()
	execCustomBuilder := build.NewExecCustomBuilder(docker2, dockerEnv, clock)
	memoryAnalytics := analytics.NewMemoryAnalytics()
	updateModeFlag := _wireUpdateModeFlagValue
	updateMode, err := ProvideUpdateMode(updateModeFlag, env, runtime)
	if err != nil {
		return nil, err
	}
	imageBuildAndDeployer := NewImageBuildAndDeployer(imageBuilder, cacheBuilder, execCustomBuilder, kClient, env, memoryAnalytics, updateMode, clock, runtime)
	return imageBuildAndDeployer, nil
}

var (
	_wireEnvValue            = k8s.Env(k8s.EnvDockerDesktop)
	_wireUpdateModeFlagValue = UpdateModeFlag(UpdateModeAuto)
)

func provideDockerComposeBuildAndDeployer(ctx context.Context, dcCli dockercompose.DockerComposeClient, dCli docker.Client, dir *dirs.WindmillDir) (*DockerComposeBuildAndDeployer, error) {
	console := build.DefaultConsole()
	writer := output.CaptureAll()
	labels := _wireLabelsValue
	dockerImageBuilder := build.NewDockerImageBuilder(dCli, console, writer, labels)
	imageBuilder := build.DefaultImageBuilder(dockerImageBuilder)
	cacheBuilder := build.NewCacheBuilder(dCli)
	env := _wireK8sEnvValue
	portForwarder := k8s.ProvidePortForwarder()
	clientConfig := k8s.ProvideClientConfig()
	namespace := k8s.ProvideConfigNamespace(clientConfig)
	kubeContext, err := k8s.ProvideKubeContext(clientConfig)
	if err != nil {
		return nil, err
	}
	kubectlRunner := k8s.ProvideKubectlRunner(kubeContext)
	client := k8s.ProvideK8sClient(ctx, env, portForwarder, namespace, kubectlRunner, clientConfig)
	runtime := k8s.ProvideContainerRuntime(ctx, client)
	dockerEnv, err := docker.ProvideEnv(ctx, env, runtime)
	if err != nil {
		return nil, err
	}
	clock := build.ProvideClock()
	execCustomBuilder := build.NewExecCustomBuilder(dCli, dockerEnv, clock)
	updateModeFlag := _wireEngineUpdateModeFlagValue
	updateMode, err := ProvideUpdateMode(updateModeFlag, env, runtime)
	if err != nil {
		return nil, err
	}
	engineImageAndCacheBuilder := NewImageAndCacheBuilder(imageBuilder, cacheBuilder, execCustomBuilder, updateMode)
	dockerComposeBuildAndDeployer := NewDockerComposeBuildAndDeployer(dcCli, dCli, engineImageAndCacheBuilder, clock)
	return dockerComposeBuildAndDeployer, nil
}

var (
	_wireK8sEnvValue               = k8s.Env(k8s.EnvNone)
	_wireEngineUpdateModeFlagValue = UpdateModeFlag(UpdateModeAuto)
)

// wire.go:

var DeployerBaseWireSet = wire.NewSet(build.DefaultConsole, output.CaptureAll, wire.Value(dockerfile.Labels{}), wire.Value(UpperReducer), docker.ProvideEnv, build.DefaultImageBuilder, build.NewCacheBuilder, build.NewDockerImageBuilder, build.NewExecCustomBuilder, wire.Bind(new(build.CustomBuilder), new(build.ExecCustomBuilder)), NewImageBuildAndDeployer, build.NewContainerUpdater, NewSyncletBuildAndDeployer,
	NewLocalContainerBuildAndDeployer,
	NewDockerComposeBuildAndDeployer,
	NewImageAndCacheBuilder,
	DefaultBuildOrder, wire.Bind(new(BuildAndDeployer), new(CompositeBuildAndDeployer)), NewCompositeBuildAndDeployer,
	ProvideUpdateMode,
	NewGlobalYAMLBuildController,
)

var DeployerWireSetTest = wire.NewSet(
	DeployerBaseWireSet,
	NewSyncletManagerForTests,
)

var DeployerWireSet = wire.NewSet(
	DeployerBaseWireSet,
	NewSyncletManager,
)
