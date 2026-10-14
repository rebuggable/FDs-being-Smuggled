package internal

import (
	libFormat "fmt"
	libOS "os"
	libRuntime "runtime"

	libLO "github.com/samber/lo"
	libUNIX "golang.org/x/sys/unix"
)

func actuateMount(grammar Grammar) error {
	mappedFDs, exception := libLO.MapErr(grammar.Mount, func(mount Mount, index int) (int, error) {
		fd, exception := libUNIX.OpenTree(
			libUNIX.AT_FDCWD,
			mount.Source,
			libUNIX.OPEN_TREE_CLONE|libUNIX.OPEN_TREE_CLOEXEC,
		)

		return fd, exception
	})

	if exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		libOS.Exit(1)
	}

	namespaceFD, exception := libUNIX.PidfdOpen(
		grammar.Instance.PID,
		0,
	)

	if exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		libOS.Exit(1)
	}

	namespaceOriginalFD, exception := libUNIX.PidfdOpen(
		libOS.Getpid(),
		0,
	)

	if exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)

		libLO.ForEach(mappedFDs, func(fd int, _ int) {
			libUNIX.Close(fd)
		})

		libOS.Exit(1)
	}

	defer libUNIX.Close(namespaceFD)

	libRuntime.LockOSThread()
	defer libRuntime.UnlockOSThread()

	if exception := libUNIX.Unshare(libUNIX.CLONE_FS); exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)

		libLO.ForEach(mappedFDs, func(fd int, _ int) {
			libUNIX.Close(fd)
		})

		libOS.Exit(1)
	}

	if exception := libUNIX.Setns(namespaceFD, libUNIX.CLONE_NEWNS); exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)

		libLO.ForEach(mappedFDs, func(fd int, _ int) {
			libUNIX.Close(fd)
		})

		libOS.Exit(1)
	}

	defer func() {
		if exception := libUNIX.Setns(namespaceOriginalFD, libUNIX.CLONE_NEWNS); exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		}
		defer libUNIX.Close(namespaceOriginalFD)
	}()

	libLO.ForEach(mappedFDs, func(fd int, index int) {
		destination := grammar.Mount[index].Destination

		if _, exception := libOS.Stat(destination); exception != nil {
			if exception := libOS.MkdirAll(destination, 0755); exception != nil {
				libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
			}
		}

		if exception := libUNIX.MoveMount(
			fd,
			"",
			libUNIX.AT_FDCWD,
			destination,
			libUNIX.MOVE_MOUNT_F_EMPTY_PATH,
		); exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		} else {
			libFormat.Printf("Succeeded at mounting %s\n", destination)
		}

		libUNIX.Close(fd)
	})

	return nil
}

var G_actuateMount = actuateMount
