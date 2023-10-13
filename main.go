package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

func main() {
	var (
		ctx    = context.Background()
		args   = os.Args[1:]
		stdin  = os.Stdin
		stdout = os.Stdout
		stderr = os.Stderr
	)

	err := cli(ctx, args, stdin, stdout, stderr)

	switch {
	case err == nil:
	case errors.Is(err, ff.ErrNoExec):
	case errors.Is(err, ff.ErrHelp):
	case err != nil:
		fmt.Fprintf(stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type CmdExec func(context.Context, []string) error

func cli(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) (err error) {
	root := NewRootCmd(stdin, stdout, stderr)
	_ = NewListCmd(root)
	_ = NewAddCmd(root)
	_ = NewViewCmd(root)

	defer func() {
		if errors.Is(err, ff.ErrHelp) {
			fmt.Fprintf(stderr, "\n%s\n", ffhelp.Command(root.Command))
		}
	}()

	if err = root.Command.Parse(args); err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	if err = root.Command.Run(ctx); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
