package main

import (
	libFormat "fmt"
	libOS "os"

	libSTDCopy "github.com/moby/moby/api/pkg/stdcopy"
	libMobyClient "github.com/moby/moby/client"
	libInternal "github.com/rebuggable/FDs-being-Smuggled/interlocutor/internal"
)

func main() {
	client, exception := libMobyClient.New(libMobyClient.WithHostFromEnv())

	if exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		libOS.Exit(1)
	}

	grammar, exception := libInternal.G_actuateProgramme(client)

	if exception != nil {
		libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		libOS.Exit(1)
	}

	if grammar.Instance[1] != nil {
		newGrammar := grammar.Clone()

		for i := range newGrammar.Mount {
			newGrammar.Mount[i].Destination = libFormat.Sprintf("/mnt/%s", libInternal.G_toIdempotentID(newGrammar.Mount[i].Source))
		}

		if attachment, exception := libInternal.G_actuateEphemeralContainer(client, newGrammar); exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
			libOS.Exit(1)
		} else {
			if _, exception := libSTDCopy.StdCopy(libOS.Stdout, libOS.Stderr, attachment.Reader); exception != nil {
				libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
			}

			defer attachment.Close()
		}

		attachment, exception := libInternal.G_actuateEngineHandoff(client, grammar)

		if exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
			libOS.Exit(1)
		}

		defer attachment.Close()

		if _, exception := libSTDCopy.StdCopy(libOS.Stdout, libOS.Stderr, attachment.Reader); exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		}
	} else {
		attachment, exception := libInternal.G_actuateEphemeralContainer(client, grammar)

		if exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
			libOS.Exit(1)
		}

		defer attachment.Close()

		if _, exception := libSTDCopy.StdCopy(libOS.Stdout, libOS.Stderr, attachment.Reader); exception != nil {
			libFormat.Fprintf(libOS.Stderr, "%v\n", exception)
		}
	}
}
