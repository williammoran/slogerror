package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/williammoran/slogerror"
)

var states = map[string]map[string][]string{
	"PA": {
		"Pittsburgh":   {"Main Street", "Babcock Boulevard"},
		"Philadelphia": {"Main Street", "Broad Street"},
	},
	"OH": {
		"Columbus": {"High Street"},
	},
}

// This example shows how the sloglogger incorporates logging context
// into error messages.

func main() {
	// Use the slogerror handler as middleware to track log
	// attributes.
	l := slog.New(slogerror.NewHandler(slog.NewTextHandler(os.Stderr, nil)))
	for stateName, cities := range states {
		l0 := l.With(slog.String("state name", stateName))
		for cityName, streets := range cities {
			l1 := l0.With(slog.String("city name", cityName))
			for _, streetName := range streets {
				l2 := l1.With(slog.String("street name", streetName))
				l2.Info("processing street")
				// There's no need to include details of the street
				// in the error message, since the logging context will
				// automatically add it.
				err := slogerror.Error(l2, "error on this street")
				// Obvoiusly, this isn't a good way to handle an error,
				// but it makes for an easy to see example of how
				// slog context is included in the error message.
				fmt.Println(err.Error())
			}
		}
	}
}
