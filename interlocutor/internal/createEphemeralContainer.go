package internal

import (
	libContext "context"
	libFormat "fmt"

	libMobyContainer "github.com/moby/moby/api/types/container"
	libMobyMount "github.com/moby/moby/api/types/mount"
	libMobyClient "github.com/moby/moby/client"
	libLO "github.com/samber/lo"
)

func createEphemeralContainer(client *libMobyClient.Client, grammar Grammar) (libMobyClient.ContainerCreateResult, error) {

	arguments := []string{libFormat.Sprintf("%d", grammar.Instance[0].PID)}

	arguments = libLO.Reduce(grammar.Mount, func(accumulated []string, mount Mount, _ int) []string {
		return append(accumulated, []string{"--mount", libFormat.Sprintf("/mnt/%s:%s", G_toIdempotentID(mount.Source), mount.Destination)}...)
	}, arguments)

	arguments = libLO.Reduce(grammar.Unmount, func(accumulated []string, unmount Unmount, _ int) []string {
		return append(accumulated, []string{"--unmount", unmount.Destination}...)
	}, arguments)

	return client.ContainerCreate(
		libContext.Background(),
		libMobyClient.ContainerCreateOptions{
			Config: &libMobyContainer.Config{
				Image: EPHEMERAL_IMAGE_NAME,

				Tty:          false,
				OpenStdin:    false,
				AttachStdin:  false,
				AttachStdout: true,
				AttachStderr: true,

				Cmd: arguments,
			},

			HostConfig: &libMobyContainer.HostConfig{
				PidMode: libMobyContainer.PidMode("host"),

				CapDrop: []string{"ALL"},
				CapAdd:  []string{"SYS_ADMIN", "SYS_CHROOT", "SYS_PTRACE"},

				Mounts: libLO.Map(grammar.Mount, func(mount Mount, _ int) libMobyMount.Mount {
					return libMobyMount.Mount{
						Type:   libMobyMount.TypeBind,
						Source: mount.Source,
						Target: libFormat.Sprintf("/mnt/%s", G_toIdempotentID(mount.Source)),
					}
				}),

				AutoRemove: true,
			},
		},
	)
}

var G_createEphemeralContainer = createEphemeralContainer
