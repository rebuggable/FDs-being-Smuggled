package internal

import (
	libContext "context"

	libMobyClient "github.com/moby/moby/client"
)

func actuateEphemeralContainer(client *libMobyClient.Client, grammar Grammar) (libMobyClient.ContainerAttachResult, error) {
	resolved, exception := G_createEphemeralContainer(client, grammar)

	if exception != nil {
		return libMobyClient.ContainerAttachResult{}, exception
	}

	attachment, exception := client.ContainerAttach(
		libContext.Background(),
		resolved.ID,
		libMobyClient.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		},
	)

	if exception != nil {
		return libMobyClient.ContainerAttachResult{}, exception
	}

	container, exception := client.ContainerStart(
		libContext.Background(),
		resolved.ID,
		libMobyClient.ContainerStartOptions{},
	)
	_ = container

	if exception != nil {
		return libMobyClient.ContainerAttachResult{}, exception
	}

	return attachment, nil
}

var G_actuateEphemeralContainer = actuateEphemeralContainer
