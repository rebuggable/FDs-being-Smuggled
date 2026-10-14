package internal

import (
	libOS "os"

	libUUID "github.com/google/uuid"
)

func isBeingPiped() (bool, error) {
	if describe, exception := libOS.Stdin.Stat(); exception != nil {
		return false, exception
	} else {
		return (describe.Mode() & libOS.ModeCharDevice) == 0, nil
	}
}

func isDinD() bool {
	return libOS.Getenv("DOCKER_MODE") == "DinD"
}

var G_isDinD = isDinD

func isDooD() bool {
	return libOS.Getenv("DOCKER_MODE") == "DooD"
}

var G_isDooD = isDooD

func toIdempotentID(value string) string {
	return libUUID.NewSHA1(UUID_NAMESPACE, []byte(value)).String()
}

var G_toIdempotentID = toIdempotentID
