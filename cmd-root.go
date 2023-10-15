package main

import (
	"context"
	"io"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffval"
	"github.com/simonklee/sourdough/query"
)

type RootCmdOptions struct {
	Stdout         io.Writer
	Stderr         io.Writer
	Stdin          io.Reader
	Verbose        bool
	FormatTerm     bool
	FormatHTML     bool
	FormatMarkdown bool
}

func (cfg *RootCmdOptions) SetupStore(ctx context.Context) (*query.Queries, error) {
	return InitStore(ctx, defaultDBPath())
}

type OutputFormat string

const (
	OutputFormatTerm     OutputFormat = "terminal"
	OutputFormatMarkdown OutputFormat = "markdown"
	OutputFormatHTML     OutputFormat = "html"
)

func (cfg *RootCmdOptions) OutputFormat() OutputFormat {
	if cfg.FormatHTML {
		return OutputFormatHTML
	}
	if cfg.FormatMarkdown {
		return OutputFormatMarkdown
	}
	return OutputFormatTerm
}

type RootCmd struct {
	Opts    RootCmdOptions
	Flags   *ff.FlagSet
	Command *ff.Command
}

const CmdLabel = "sourdough"

func NewRootCmd(stdin io.Reader, stdout, stderr io.Writer) *RootCmd {
	cmd := RootCmd{
		Opts: RootCmdOptions{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		},
	}
	cmd.Flags = ff.NewFlagSet(CmdLabel)
	_, _ = cmd.Flags.AddFlag(ff.FlagConfig{
		ShortName: 'v',
		LongName:  "verbose",
		Value:     ffval.NewValue(&cmd.Opts.Verbose),
		Usage:     "log verbose output",
		NoDefault: true,
	})
	cmd.Flags.BoolVar(&cmd.Opts.FormatTerm, 0, "term", "output in terminal format")
	cmd.Flags.BoolVar(&cmd.Opts.FormatHTML, 0, "html", "output in HTML format")
	cmd.Flags.BoolVar(&cmd.Opts.FormatMarkdown, 0, "markdown", "output in Markdown format")
	cmd.Command = &ff.Command{
		Name:      CmdLabel,
		ShortHelp: "sourdough is a CLI tool for managing recipes and baking sourdough bread",
		Usage:     CmdLabel + " [flags] <subcommand> ...",
		Flags:     cmd.Flags,
		Exec: func(ctx context.Context, args []string) error {
			return ff.ErrHelp
		},
	}

	return &cmd
}
