package main

import (
	"context"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/peterbourgon/ff/v4"
	"github.com/simonklee/sourdough/query"
)

type ListCmdOptions struct {
	Root *RootCmdOptions
}

type ListCmd struct {
	Opts ListCmdOptions

	root    *RootCmd
	Flags   *ff.FlagSet
	Command *ff.Command
}

func NewListCmd(parent *RootCmd) *ListCmd {
	var cmd ListCmd
	cmd.Opts.Root = &parent.Opts
	cmd.root = parent
	cmd.Flags = ff.NewFlagSet("list").SetParent(parent.Flags)
	cmd.Command = &ff.Command{
		Name:      "list",
		Usage:     CmdLabel + " list [flags]",
		ShortHelp: "list recipes",
		Flags:     cmd.Flags,
		Exec:      ListCmdExec(&cmd.Opts),
	}
	cmd.root.Command.Subcommands = append(cmd.root.Command.Subcommands, cmd.Command)
	return &cmd
}

func ListCmdExec(opts *ListCmdOptions) CmdExec {
	return func(ctx context.Context, args []string) error {
		db, err := opts.Root.SetupStore(ctx)
		if err != nil {
			return err
		}

		recipes, err := db.ListRecipes(ctx)
		if err != nil {
			return err
		}

		if len(recipes) == 0 {
			return nil
		}

		return listRecipes(opts.Root.Stdout, opts.Root.OutputFormat(), recipes)
	}
}

func listRecipes(w io.Writer, format OutputFormat, recipes []query.Recipe) error {
	// Create a new table writer
	tw := table.NewWriter()

	// Append a header row
	tw.AppendHeader(table.Row{"ID", "Name"})

	// Append data rows
	for _, recipe := range recipes {
		tw.AppendRow(table.Row{recipe.ID, recipe.Name})
	}

	return renderTable(w, format, tw)
}
