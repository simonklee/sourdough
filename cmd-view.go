package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/peterbourgon/ff/v4"
	"github.com/simonklee/sourdough/query"
)

type ViewCmdOptions struct {
	Dependencies []string

	Root *RootCmdOptions
}

type ViewCmd struct {
	Opts ViewCmdOptions

	root    *RootCmd
	Flags   *ff.FlagSet
	Command *ff.Command
}

func NewViewCmd(parent *RootCmd) *ViewCmd {
	var cmd ViewCmd
	cmd.Opts.Root = &parent.Opts
	cmd.root = parent
	cmd.Flags = ff.NewFlagSet("view").SetParent(parent.Flags)
	cmd.Flags.StringListVar(&cmd.Opts.Dependencies, 'd', "dependency", "dependency (e.g. --dependency \"total_floor 450g\")")

	cmd.Command = &ff.Command{
		Name:      "view",
		Usage:     CmdLabel + " view [flags] <recipe>",
		ShortHelp: "view recipe",
		LongHelp: `  By default the recipes template (relative values) will be
  displayed. To view the amount of ingredients required for 
  a specific portion use the --dependencies flag to specify 
  the dependencies.
  
  Example:
  
     $ sourdough view --dependency "total_flour 450g" 1

`,
		Flags: cmd.Flags,
		Exec:  ViewCmdExec(&cmd.Opts),
	}
	cmd.root.Command.Subcommands = append(cmd.root.Command.Subcommands, cmd.Command)
	return &cmd
}

func ViewCmdExec(opts *ViewCmdOptions) CmdExec {
	return func(ctx context.Context, args []string) error {
		if len(args) == 0 {
			return errors.New("requires a recipe ID or name")
		}

		value := args[0]
		recipeID, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid recipe ID: %w", err)
		}

		dependencies, err := ParseDependencies(opts.Dependencies)
		if err != nil {
			return err
		}

		db, err := opts.Root.SetupStore(ctx)
		if err != nil {
			return err
		}

		recipe, err := db.GetRecipe(ctx, recipeID)
		if err != nil && errors.Is(err, query.ErrNotFound) {
			fmt.Fprintf(opts.Root.Stdout, "recipe %d not found\n", recipeID)
			return nil
		} else if err != nil {
			return err
		}

		ingredients, err := db.ListRecipeIngredients(ctx, recipeID)
		if err != nil {
			return err
		}

		var portionIngredients []PortionIngredient
		if len(dependencies) > 0 {
			portionIngredients, err = Calculate(ingredients, dependencies)
			if err != nil {
				return err
			}
		}

		return RecipeView{
			Recipe:      recipe,
			Ingredients: ingredients,
			Portions:    portionIngredients,
		}.Render(ctx, opts.Root.Stdout, opts.Root.OutputFormat())
	}
}

type RecipeView struct {
	Recipe      query.Recipe
	Ingredients []query.ListRecipeIngredientsRow
	Portions    []PortionIngredient
}

func (r RecipeView) Render(ctx context.Context, w io.Writer, format OutputFormat) error {
	// Create a new table writer
	tw := table.NewWriter()

	// Set table style
	tw.SetStyle(table.StyleLight)

	// Set table title
	tw.SetTitle(fmt.Sprintf("Recipe: %s", r.Recipe.Name))

	// Configure columns
	tw.AppendHeader(table.Row{"#", "Ingredient", "Unit", "Percentage", "Dependency"})
	for _, ingredient := range r.Ingredients {
		// v, _ := strconv.ParseFloat(fmt.Sprintf("%f", ingredient.Percentage*100), 64)
		tw.AppendRow(table.Row{
			ingredient.ID,
			ingredient.Name,
			ingredient.UnitType,
			ingredient.Percentage,
			ingredient.Dependency,
		})
	}

	if err := renderTable(w, format, tw); err != nil {
		return err
	}

	if len(r.Portions) > 0 {
		tw := table.NewWriter()
		tw.SetStyle(table.StyleLight)
		tw.SetTitle(fmt.Sprintf("Ingredients for: %s", r.Recipe.Name))
		tw.SetColumnConfigs([]table.ColumnConfig{
			{Number: 3, Align: text.AlignRight},
		})

		tw.AppendHeader(table.Row{"#", "Ingredient", "Amount"})
		for i, portion := range r.Portions {
			tw.AppendRow(table.Row{
				i + 1,
				portion.Name,
				portion.Unit.Appropriate(portion.Amount).Format(),
			})
		}

		if err := renderTable(w, format, tw); err != nil {
			return err
		}
	}

	return nil
}
