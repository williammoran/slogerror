# slogerror

[![Go Reference](https://pkg.go.dev/badge/github.com/williammoran/slogerror.svg)](https://pkg.go.dev/github.com/williammoran/slogerror)

I've made a video about the library, in case you prefer that to reading:
(https://youtu.be/5sFlWf9AdZY)

## The Problem

In addition to making it easy to integrate with log aggregators,
Go's slog package also makes it intuitive to build context as to
where important events are happening within the code, primarily
through the use of `logger.With()`.

However, when an error occurs, that context needs to be reconstructed
in the creation of the error, since there's no way to transfer the
logging context that's been constructed to the error. Take this
example:

```go
for stateName, cities := range states {
	l0 := l.With(slog.String("state name", stateName))
	for cityName, streets := range cities {
		l1 := l0.With(slog.String("city name", cityName))
		for _, streetName := range streets {
			l2 := l1.With(slog.String("street name", streetName))
			l2.Debug("processing street")
			err := fmt.Errorf("error on %s in %s, %s", streetName, cityName, stateName)
			// ...
		}
	}
}
```

Creating the error requires recreating the context that already
exists in the logger. This is redundant and error prone, and
can be quite onerous when the context is complex or extensive.
Sometimes it's not even possible because the code may not
have access to all the values that are required to
construct a full context. Since the structured log package is
designed to build that context up already, it makes sense to
make use of it.

## A Solution

`slogerror` allows you to create error variables that derive
contextual information from an `slog.Logger`

```go
// Use the slogerror error handler to track log attributes
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
			err := slogerror.Error(l2, "error on street")
			// ...
		}
	}
}
```

`slogerror.Error()` includes all context from the Logger in the error
message automatically. See the examples directory for more information.

###### Install

```sh
go get github.com/williammoran/slogerror
```

###### 
