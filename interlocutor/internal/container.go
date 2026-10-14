package internal

import (
	libContext "context"
	libHEX "encoding/hex"
	libFormat "fmt"

	libMobyContainer "github.com/moby/moby/api/types/container"
	libMobyClient "github.com/moby/moby/client"
)

func isContainerID(value string) bool {
	if !(len(value) == DOCKER_CONTAINER_ID_FIXED_LENGTH || len(value) == DOCKER_CONTAINER_ID_TRUNCATED) {
		return false
	}

	if _, exception := libHEX.DecodeString(value); exception != nil {
		return false
	} else {
		return true
	}
}

var G_isContainerID = isContainerID

func isContainerName(value string) bool {
	return len(value) != 0 && !isContainerID(value)
}

var G_isContainerName = isContainerName

func toContainerFromFilter(client *libMobyClient.Client, filter libMobyClient.Filters) (libMobyContainer.InspectResponse, error) {
	resolved, exception := client.ContainerList(
		libContext.Background(),
		libMobyClient.ContainerListOptions{
			Filters: filter, Limit: 1, All: true,
		},
	)

	if exception != nil {
		return libMobyContainer.InspectResponse{}, exception
	}

	entities := resolved.Items

	if len(entities) == 0 {
		return libMobyContainer.InspectResponse{}, libFormat.Errorf("You're looking for a non-existent container instance")
	}

	entity := entities[0]

	inspection, exception := client.ContainerInspect(
		libContext.Background(),
		entity.ID,
		libMobyClient.ContainerInspectOptions{},
	)

	if exception != nil {
		return libMobyContainer.InspectResponse{}, exception
	}

	return inspection.Container, nil
}

var G_toContainerFromFilter = toContainerFromFilter

func toContainerFromID(client *libMobyClient.Client, containerID string) (libMobyContainer.InspectResponse, error) {
	filter := make(libMobyClient.Filters)
	filter.Add("id", containerID)

	return toContainerFromFilter(client, filter)
}

var G_toContainerFromID = toContainerFromID

func toContainerFromName(client *libMobyClient.Client, containerName string) (libMobyContainer.InspectResponse, error) {
	filter := make(libMobyClient.Filters)
	filter.Add("name", containerName)

	return toContainerFromFilter(client, filter)
}

var G_toContainerFromName = toContainerFromName

func toPIDFromContainerID(client *libMobyClient.Client, containerID string) (int, error) {
	container, exception := toContainerFromID(client, containerID)

	if exception != nil {
		return 0, exception
	}

	return container.State.Pid, nil
}

var G_toPIDFromContainerID = toPIDFromContainerID

func toPIDFromContainerName(client *libMobyClient.Client, containerName string) (int, error) {
	container, exception := toContainerFromName(client, containerName)

	if exception != nil {
		return 0, exception
	}

	return container.State.Pid, nil
}

var G_toPIDFromContainerName = toPIDFromContainerName

func toContainerFrom(client *libMobyClient.Client, value string) (libMobyContainer.InspectResponse, error) {
	if isContainerID(value) {
		return toContainerFromID(client, value)
	} else if isContainerName(value) {
		return toContainerFromName(client, value)
	} else {
		return libMobyContainer.InspectResponse{}, libFormat.Errorf("You've provided an unrecognized input value")
	}
}

var G_toContainerFrom = toContainerFrom
