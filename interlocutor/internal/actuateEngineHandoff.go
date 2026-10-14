package internal

import (
	libContext "context"
	libFormat "fmt"

	libMobyClient "github.com/moby/moby/client"
	libLO "github.com/samber/lo"
)

func actuateEngineHandoff(client *libMobyClient.Client, grammar Grammar) (libMobyClient.ExecAttachResult, error) {
	container, exception := toContainerFrom(client, grammar.Instance[0].ContainerID)

	if exception != nil {
		return libMobyClient.ExecAttachResult{}, exception
	}

	innerInstance := grammar.Instance[1]

	arguments := []string{
		"/usr/local/bin/orchestrator",
		G_chooseEither(innerInstance.ContainerID, innerInstance.ContainerName),
	}

	arguments = libLO.Reduce(grammar.Mount, func(accumulated []string, mount Mount, _ int) []string {
		libFormat.Printf("Dest= %s\n", mount.Destination)
		return append(accumulated, []string{"--mount", libFormat.Sprintf("/mnt/%s:%s", toIdempotentID(mount.Source), mount.Destination)}...)
	}, arguments)

	arguments = libLO.Reduce(grammar.Unmount, func(accumulated []string, unmount Unmount, _ int) []string {
		return append(accumulated, []string{"--unmount", unmount.Destination}...)
	}, arguments)

	exec, exception := client.ExecCreate(
		libContext.Background(),
		container.ID,
		libMobyClient.ExecCreateOptions{
			Cmd: arguments,

			AttachStdout: true,
			AttachStderr: true,
			TTY:          false,
		},
	)

	if exception != nil {
		return libMobyClient.ExecAttachResult{}, exception
	}

	attachment, exception := client.ExecAttach(
		libContext.Background(),
		exec.ID,
		libMobyClient.ExecAttachOptions{
			TTY: false,
		},
	)

	resolved, exception := client.ExecStart(
		libContext.Background(),
		exec.ID,
		libMobyClient.ExecStartOptions{},
	)
	_ = resolved

	if exception != nil {
		return libMobyClient.ExecAttachResult{}, exception
	}

	return attachment, nil
}

var G_actuateEngineHandoff = actuateEngineHandoff
