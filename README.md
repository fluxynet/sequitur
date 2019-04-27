![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen.svg)
[![GoDoc](https://img.shields.io/badge/go-doc%20-blue.svg)](https://godoc.org/github.com/fluxynet/sequitur)

# Sequitur

Sequitur is a simple way of handling errors in Go.

## Overview
- [Sequitur](#sequitur)
- [Overview](#overview)
- [Features](#features)
- [Usage](#usage)
- [Examples](#examples)
    + [Simple Sequence](#simple-sequence)
    + [Error handling](#error-handling)
    + [Deferred Error handling](#deferred-error-handling)
    + [Panic handling](#panic-handling)
    + [Context handling](#context-handling)
    + [Logging with logrus](#logging-with-logrus)
- [Acknowledgements](#acknowledgements)


## Features

-   Simple syntax, no more `if value, err = action(); err != nil {`.
-   Each action must have a description which makes for user-friendly error messages.
-   Automatic panic handling. If any `Action` panics, it is treated as an error of type `sequitur.ErrorPanic`.
-   Logging can be done fairly easily using Logrus (more mechanisms to be added).
-   Context handling. If context expires, sequence stops.

## Usage

A `Sequence` is a series of `Action` that are executed.
Each `Action` is just a `func() error`

If any `Action` returns an `error`, the `Sequence` is stopped.

The `Catch` function of `Sequence` is called.
If there are no errors, the `Then` function of the `Sequence` is called.

```go
import "github.com/fluxynet/sequitur"
//...
seq := sequitur.Linear() //currently only Linear sequence is supported

seq.Do("short action1 name", func() error {
    //...
})

seq.Do("short action2 name", func() error {
    //...
})

//executed if any action returned error
//name is one of the names above, e.g. short action1 name
//it could be used to provide a user friendly description of the error to the user
//it could even be passed via an i18n lib
//err can be used for logging or for logical branching
seq.Catch(func(name string, err error) {
    //...
})

seq.Then(func() { //executed if no action returned error
    //...
})
```

## Examples

### Simple Sequence

```go
import (
    "fmt"
    "github.com/fluxynet/sequitur"
    "io/ioutil"
    "log"
)

func main() {
    var data []byte

    seq := sequitur.Linear()

    seq.Do("reading file from disk", func() error {
        var err error
        data, err = ioutil.ReadFile("foo.bar")
        return err
    })

    seq.Do("writing file to disk", func() error {
        return ioutil.WriteFile("bar.foo", data, 0644)
    })

    seq.Catch(func(name string, err error) {
        fmt.Printf("An error occurred during %s.", name)
        log.Println(err)
    })

    seq.Then(func() {
        fmt.Println("Wrote %d bytes", len(data))
    })
}
```

### Error handling

```go
import (
    "errors"
    "fmt"
    "github.com/fluxynet/sequitur"
    "net/http"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        X int `json:"x"`
        Y int `json:"y"`
    }

    var output struct {
        Sum   int `json:"sum"`
    }

    seq := sequitur.Linear()

    var b []byte
    sequence.Do("reading request body", func() error {
		b, err = ioutil.ReadAll(r.Body)
		return err
	})

	sequence.Do("decoding request body", func() error {
		return json.Unmarshal(b, &input)
    })

    output.Sum = input.X + input.Y //no need to use action, does not yield error

    var j []byte
	sequence.Do("writing response", func() error {
		j, err = json.Marshal(response)
		return err
	})

	sequence.Then(func() {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
    })

	seq.Catch(func(name string, err error) {
        msg := fmt.Sprintf("An error occurred during %s.", name, err.String())
        log.Println(err)

        w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error":"` + msg + `"}`))
	})
}
```

### Deferred Error handling

```go
import (
    "errors"
    "fmt"
    "github.com/fluxynet/sequitur"
    "net/http"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        X int `json:"x"`
        Y int `json:"y"`
    }

    var output struct {
        Sum   int `json:"sum"`
    }

    seq := sequitur.Linear()
    defer sequence.Catch(catchError(w, r)) //all my errors go here

    var b []byte
    sequence.Do("reading request body", func() error {
		b, err = ioutil.ReadAll(r.Body)
		return err
	})

	sequence.Do("decoding request body", func() error {
		return json.Unmarshal(b, &input)
    })

    output.Sum = input.X + input.Y //no need to use action, does not yield error

    var j []byte
	sequence.Do("writing response", func() error {
		j, err = json.Marshal(response)
		return err
	})

	sequence.Then(func() {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
    })
}

//some more 'elaborate' error handling
func catchError(w http.ResponseWriter, r *http.Request) sequitur.Consequence {
	return func(name string, err error) {
		var (
			msg    string
			status int
        )

        log.Println(err)

		switch err {
		default:
			msg = "Error when " + name
			if strings.HasPrefix(err.Error(), "invalid") {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
		case ErrorInvalidToken:
			msg = "Invalid Token"
			status = http.StatusForbidden
		case ErrorNoToken:
			msg = "You must login to proceed"
			status = http.StatusForbidden
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(`{"error":"` + msg + `"}`))
	}
}
```

### Panic handling

```go
import (
    "fmt"
    "log"
    "github.com/fluxynet/sequitur"
)

func main() {
    var result string

	seq.Do("letter a", func() error {
		result += "a"
		return nil
	})

	seq.Do("misbehave", func() error {
		panic("foobar")
	})

	seq.Do("letter b", func() error {
		result += "b"
		return nil
	})

	seq.Catch(func(name string, err error) {
        result += "z"
        if err == sequitur.ErrPanic {
            log.Println("Panic averted during: "+name)
        } else {
            log.Println(name, err)
        }
    })

    seq.Then(func() {
        result += "c"
    })

    fmt.Println(result) //az
}
```

### Context handling

```go
ctxC, cancel := context.WithCancel(context.Background())

cancel()
//cancel and the sequence will stop
//although an ongoing action will never yield
//think of it as a "critical section"
```

### Logging with logrus

```go
seq := sequitur.WithLogrus(sequitur.Linear())

//or if using logrus instance
seq := sequitur.WithLogrus(sequitur.Linear(), myloggerInstance)
```

## Acknowledgements

This library is heavily inspired by a blog post from Martin Kühl entitled [Rob Pike Reinvented Monads](https://www.innoq.com/en/blog/golang-errors-monads/).
