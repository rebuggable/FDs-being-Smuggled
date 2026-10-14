package internal

import (
	libFormat "fmt"
	libOS "os"
	libRuntime "runtime"

	libLO "github.com/samber/lo"
	libUNIX "golang.org/x/sys/unix"
)

func actuateUnmount(grammar Grammar) error {
	namespaceFD, exception := libUNIX.PidfdOpen(
		grammar.Instance.PID,
		0,
	)

	if exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)

		libOS.Exit(1)
	}

	defer libUNIX.Close(namespaceFD)

	namespaceOriginalFD, exception := libUNIX.PidfdOpen(
		libOS.Getpid(),
		0,
	)

	if exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)

		libOS.Exit(1)
	}

	libRuntime.LockOSThread()
	defer libRuntime.UnlockOSThread()

	if exception := libUNIX.Unshare(libUNIX.CLONE_FS); exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)

		libOS.Exit(1)
	}

	if exception := libUNIX.Setns(namespaceFD, libUNIX.CLONE_NEWNS); exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)

		libOS.Exit(1)
	}

	defer func() {
		if exception := libUNIX.Setns(namespaceOriginalFD, libUNIX.CLONE_NEWNS); exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		}
		defer libUNIX.Close(namespaceOriginalFD)
	}()

	libLO.ForEach(grammar.Unmount, func(unmount Unmount, index int) {
		destination := unmount.Destination

		if exception := libUNIX.Unmount(destination, 0); exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v", exception)
		} else {
			libFormat.Printf("Succeeded at unmounting  %s\n", destination)
		}
	})

	return nil
}

var G_actuateUnmount = actuateUnmount
