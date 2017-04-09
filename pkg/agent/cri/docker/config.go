package docker

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"strconv"
)

func GetConfig(spec *types.ContainerSpec) *container.Config {

	var volumes map[string]struct{}
	var ports nat.PortSet

	ports = make(map[nat.Port]struct{})
	volumes = make(map[string]struct{})

	for _, p := range spec.Ports {
		port := nat.Port(strconv.Itoa(p.ContainerPort))
		ports[port] = struct{}{}
	}

	cfg := &container.Config{
		Hostname:     spec.Network.Hostname,
		Domainname:   spec.Network.Domain,
		Env:          spec.Envs,
		ExposedPorts: ports,
		Volumes:      volumes,
		Labels:       spec.Labels,
		Cmd:          strslice.StrSlice(spec.Command),
		Entrypoint:   strslice.StrSlice(spec.Entrypoint),
		Image:        spec.Image.Name,
	}
	return cfg
}

func GetHostConfig(spec *types.ContainerSpec) *container.HostConfig {

	rPolicy := container.RestartPolicy{
		Name:              spec.RestartPolicy.Name,
		MaximumRetryCount: spec.RestartPolicy.Attempt,
	}

	resources := container.Resources{
		Memory:    spec.Quota.Memory,
		CPUShares: spec.Quota.CPUShares,
	}

	var ports nat.PortMap
	var mounts []mount.Mount
	var binds []string

	for _, v := range spec.Volumes {
		mnt := mount.Mount{
			Type:   mount.TypeVolume,
			Source: v.Volume,
			Target: v.MountPath,
		}
		mounts = append(mounts, mnt)
		binds = append(binds, v.Volume)
	}

	log := container.LogConfig{}
	ports = make(nat.PortMap)

	cfg := &container.HostConfig{
		Binds:         binds,
		LogConfig:     log,
		NetworkMode:   container.NetworkMode(spec.Network.Mode),
		PortBindings:  ports,
		DNS:           spec.DNS.Server,
		DNSOptions:    spec.DNS.Options,
		DNSSearch:     spec.DNS.Search,
		RestartPolicy: rPolicy,
		Resources:     resources,
		Mounts:        mounts,
	}
	return cfg
}

func GetNetworkConfig(spec *types.ContainerSpec) *network.NetworkingConfig {
	cfg := &network.NetworkingConfig{
		EndpointsConfig: make(map[string]*network.EndpointSettings),
	}

	endpoint := &network.EndpointSettings{
		NetworkID: spec.Network.Network,
	}
	cfg.EndpointsConfig["lo"] = endpoint

	return cfg
}
