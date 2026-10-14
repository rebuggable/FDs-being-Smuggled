package main

import (
	libFormat "fmt"
	libOS "os"

	libMobyClient "github.com/moby/moby/client"
	libInternal "github.com/rebuggable/FDs-being-Smuggled/orchestrator/internal"
)

func main() {
	client, exception := libMobyClient.New(libMobyClient.WithHostFromEnv())

	if exception != nil && libInternal.G_isDinD() {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		libOS.Exit(1)
	}

	grammar, exception := libInternal.G_actuateProgramme(client)

	if exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		libOS.Exit(1)
	}

	if len(grammar.Mount) > 0 {
		if exception := libInternal.G_actuateMount(grammar); exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
			libOS.Exit(1)
		}
	}

	if len(grammar.Unmount) > 0 {
		if exception := libInternal.G_actuateUnmount(grammar); exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
			libOS.Exit(1)
		}
	}
}
