package internal

import (
	libFormat "fmt"
	libOS "os"
	libReflect "reflect"
	libString "strings"

	libKong "github.com/alecthomas/kong"
	libMobyClient "github.com/moby/moby/client"
)

type Instance struct {
	PID int
}

func useInstanceMapper(client *libMobyClient.Client) libKong.MapperFunc {
	return func(ctx *libKong.DecodeContext, field libReflect.Value) error {
		reference, exception := ctx.Scan.PopValue("")

		if exception != nil {
			return exception
		}

		value := reference.String()

		containerPID, exception := G_toPIDFromContainer(client, value)

		if exception != nil {
			return exception
		}

		instance := Instance{
			PID: containerPID,
		}

		field.Set(libReflect.ValueOf(instance))
		return nil
	}
}

type Mount struct {
	Source      string
	Destination string
}

func useMountMapper() libKong.MapperFunc {
	return func(ctx *libKong.DecodeContext, field libReflect.Value) error {
		reference, exception := ctx.Scan.PopValue("")

		if exception != nil {
			return exception
		}

		value := reference.String()

		source, destination, isOK := libString.Cut(value, ":")

		if !isOK {
			return libFormat.Errorf("You provided mount as '%s' which is not a valid format", value)
		}

		if _, exception := libOS.Stat(source); exception != nil {
			return libFormat.Errorf("You provided mount as '%s' which is not an accessible source", value)
		}

		mount := Mount{
			Source:      source,
			Destination: destination,
		}

		field.Set(libReflect.Append(field, libReflect.ValueOf(mount)))

		return nil
	}
}

type Unmount struct {
	Destination string
}

func useUnmountMapper() libKong.MapperFunc {
	return func(ctx *libKong.DecodeContext, field libReflect.Value) error {
		reference, exception := ctx.Scan.PopValue("")

		if exception != nil {
			return exception
		}

		value := reference.String()

		destination := value

		mount := Unmount{
			Destination: destination,
		}

		field.Set(libReflect.Append(field, libReflect.ValueOf(mount)))

		return nil
	}
}

type Grammar struct {
	Instance Instance  `name:"instance" placeholder:"<CONTAINER ID | CONTAINER NAME | PID>" required:"" arg:""`
	Mount    []Mount   `name:"mount"    placeholder:"[--mount <SOURCE>:<DEST>]..."`
	Unmount  []Unmount `name:"unmount"  placeholder:"[--unmount <DEST>]..." help:"A destination you'd like to remove in the picked container instance"`
}

func (self *Grammar) AfterApply(ctx *libKong.Context) error {
	if len(self.Mount) == 0 && len(self.Unmount) == 0 {
		return libFormat.Errorf("You must specify at least one of `mount` or `unmount` with the proper value.")
	}

	return nil
}
