package tui

import (
	"errors"
	"fmt"
)

var errNothingSelected = errors.New("nothing selected")

func errUnknownCommand(c string) error { return fmt.Errorf("unknown command: %s", c) }
func errMissingArg(name string) error  { return fmt.Errorf("missing argument: %s", name) }
func errInvalidArg(name string) error  { return fmt.Errorf("invalid argument: %s", name) }
