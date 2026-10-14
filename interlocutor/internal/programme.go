package internal

import (
	libIO "io"
	libOS "os"
	libReflect "reflect"
	libString "strings"

	libKong "github.com/alecthomas/kong"
	libMobyClient "github.com/moby/moby/client"
)

func actuateProgramme(client *libMobyClient.Client) (Grammar, error) {
	var grammar Grammar

	arguments := libOS.Args[1:]

	instanceMapper := useInstanceMapper(client)
	mountMapper := useMountMapper()
	unmountMapper := useUnmountMapper()

	programme, exception := libKong.New(
		&grammar,
		libKong.TypeMapper(libReflect.TypeOf(Pair[Instance]{}), instanceMapper),
		libKong.TypeMapper(libReflect.TypeOf([]Mount{}), mountMapper),
		libKong.TypeMapper(libReflect.TypeOf([]Unmount{}), unmountMapper),
		libKong.UsageOnError(),
	)

	if exception != nil {
		return grammar, exception
	}

	if isPiped, exception := isBeingPiped(); exception != nil {
		return grammar, exception
	} else if isPiped {
		value, exception := libIO.ReadAll(libOS.Stdin)

		if exception != nil {
			return grammar, exception
		}

		arguments = append(arguments, libString.Fields(string(value))...)
	}

	ctx, exception := programme.Parse(arguments)
	_ = ctx

	if exception != nil {
		return grammar, exception
	}

	return grammar, nil
}

var G_actuateProgramme = actuateProgramme
