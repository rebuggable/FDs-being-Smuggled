package internal

import (
	libOS "os"
	libFilePath "path/filepath"

	libUUID "github.com/google/uuid"
	libLO "github.com/samber/lo"
)

type Pair[T any] [2]*T

func isBeingPiped() (bool, error) {
	if describe, exception := libOS.Stdin.Stat(); exception != nil {
		return false, exception
	} else {
		return (describe.Mode() & libOS.ModeCharDevice) == 0, nil
	}
}

func chooseEither(valueA string, valueB string) string {
	return libLO.Ternary(len(valueA) > 0, valueA, valueB)
}

var G_chooseEither = chooseEither

func toAbsoluteDir(value string) string {
	if dir, exception := libFilePath.Abs(value); exception != nil {
		return value
	} else {
		return dir
	}
}

var G_toAbsoluteDir = toAbsoluteDir

func toIdempotentID(value string) string {
	return libUUID.NewSHA1(UUID_NAMESPACE, []byte(value)).String()
}

var G_toIdempotentID = toIdempotentID
