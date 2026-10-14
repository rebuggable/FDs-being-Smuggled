package internal

import (
	libFormat "fmt"
	libOS "os"
	libReflect "reflect"
	libString "strings"

	libKong "github.com/alecthomas/kong"
	libMobyClient "github.com/moby/moby/client"
	libLO "github.com/samber/lo"
)

type Instance struct {
	ContainerID   string
	ContainerName string
	PID           int
}

func useInstanceMapper(client *libMobyClient.Client) libKong.MapperFunc {
	return func(ctx *libKong.DecodeContext, field libReflect.Value) error {
		reference, exception := ctx.Scan.PopValue("")

		if exception != nil {
			return exception
		}

		instancePair := Pair[Instance]{}

		value := reference.String()

		valueA, valueB, isPaired := libString.Cut(value, "/")

		container, exception := G_toContainerFrom(client, valueA)

		if exception != nil {
			return exception
		}

		instancePair[0] = &Instance{
			ContainerID:   container.ID,
			ContainerName: container.Name,
			PID:           container.State.Pid,
		}

		if isPaired {
			instancePair[1] = &Instance{
				ContainerID:   libLO.Ternary(G_isContainerID(valueB), valueB, ""),
				ContainerName: libLO.Ternary(G_isContainerName(valueB), valueB, ""),
				PID:           0,
			}
		}

		field.Set(libReflect.ValueOf(instancePair))
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
			Source:      G_toAbsoluteDir(source),
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
	Instance Pair[Instance] `name:"instance" placeholder:"<CONTAINER ID | CONTAINER NAME>" required:"" arg:"" help:"A running container instance picked as target"`
	Mount    []Mount        `name:"mount"    placeholder:"[--mount <SOURCE>:<DEST>]..."                       help:"A map from source to destination used to mount in the picked container instance"`
	Unmount  []Unmount      `name:"unmount"  placeholder:"[--unmount <DEST>]..."                              help:"A destination you'd like to remove in the picked container instance"`
}

func (self *Grammar) AfterApply(ctx *libKong.Context) error {
	if len(self.Mount) == 0 && len(self.Unmount) == 0 {
		return libFormat.Errorf("You must specify at least one of `mount` or `unmount` with the proper value.")
	}
	return nil
}

func (g Grammar) Clone() Grammar {
	clone := g

	clone.Mount = make([]Mount, len(g.Mount))
	copy(clone.Mount, g.Mount)

	clone.Unmount = make([]Unmount, len(g.Unmount))
	copy(clone.Unmount, g.Unmount)

	for i, instance := range g.Instance {
		if instance != nil {
			value := *instance
			clone.Instance[i] = &value
		}
	}

	return clone
}
