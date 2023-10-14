package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/peterbourgon/ff/v4"
	"github.com/simonklee/sourdough/query"
	"github.com/simonklee/sourdough/recipe"
)

type AddCmdOptions struct {
	Name string

	Root *RootCmdOptions
}

type AddCmd struct {
	Opts AddCmdOptions

	root    *RootCmd
	Flags   *ff.FlagSet
	Command *ff.Command
}

func NewAddCmd(parent *RootCmd) *AddCmd {
	var cmd AddCmd
	cmd.Opts.Root = &parent.Opts
	cmd.root = parent
	cmd.Flags = ff.NewFlagSet("add").SetParent(parent.Flags)
	cmd.Flags.StringVar(&cmd.Opts.Name, 'n', "name", "", "name of the recipe")

	cmd.Command = &ff.Command{
		Name:      "add",
		Usage:     CmdLabel + " add [flags]",
		ShortHelp: "add a new recipe",
		Flags:     cmd.Flags,
		Exec:      AddCmdExec(&cmd.Opts),
	}
	cmd.root.Command.Subcommands = append(cmd.root.Command.Subcommands, cmd.Command)
	_ = newAddIngredientCmd(&cmd)

	return &cmd
}

func AddCmdExec(opts *AddCmdOptions) CmdExec {
	return func(ctx context.Context, args []string) error {
		db, err := opts.Root.SetupStore(ctx)
		if err != nil {
			return err
		}

		// Create new recipe
		_, err = db.CreateRecipe(ctx, opts.Name)
		if err != nil {
			return err
		}

		return err
	}
}

type AddIngredientCmdOptions struct {
	Name           string  `validate:"required"`
	UnitType       string  `validate:"required,oneof=weight volume count teaspoon"`
	RecipeID       int     `validate:"gt=0,required"`
	Percentage     float64 `validate:"gte=0,lte=1,required"`
	Dependency     string
	IngredientType string

	Parent *AddCmdOptions
}

type addIngredientCmd struct {
	Opts AddIngredientCmdOptions

	parent  *AddCmd
	Flags   *ff.FlagSet
	Command *ff.Command
}

func newAddIngredientCmd(parent *AddCmd) *addIngredientCmd {
	var cmd addIngredientCmd
	cmd.Opts.Parent = &parent.Opts
	cmd.parent = parent
	cmd.Flags = ff.NewFlagSet("ingredient")
	cmd.Flags.IntVar(&cmd.Opts.RecipeID, 'r', "recipe", 0, "recipe ID")
	cmd.Flags.StringVar(&cmd.Opts.Name, 'n', "name", "", "name of the ingredient")
	cmd.Flags.StringEnumVar(&cmd.Opts.UnitType, 'u', "unit", "unit type of the ingredient", "weight", "volume", "count", "teaspoon")
	cmd.Flags.Float64Var(&cmd.Opts.Percentage, 'p', "percentage", 0.0, "percentage of the ingredient")
	cmd.Flags.StringVar(&cmd.Opts.Dependency, 'd', "dependency", "", "dependency of the ingredient")
	cmd.Flags.StringVar(&cmd.Opts.IngredientType, 't', "type", "", "type of the ingredient")
	cmd.Command = &ff.Command{
		Name:      "ingredient",
		Usage:     CmdLabel + " add ingredient <recipe> [flags]",
		ShortHelp: "add a new ingredient to a recipe",
		Flags:     cmd.Flags,
		Exec:      addIngredientCmdExec(&cmd.Opts),
	}
	cmd.parent.Command.Subcommands = append(cmd.parent.Command.Subcommands, cmd.Command)

	return &cmd
}

func addIngredientCmdExec(opts *AddIngredientCmdOptions) CmdExec {
	return func(ctx context.Context, args []string) error {
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(opts)
		if err != nil {
			return err
		}

		db, err := opts.Parent.Root.SetupStore(ctx)
		if err != nil {
			return err
		}

		// Check if recipe exists
		r, err := db.GetRecipe(ctx, int64(opts.RecipeID))
		if err != nil {
			return fmt.Errorf("recipe %d does not exist: %w", opts.RecipeID, err)
		}

		if opts.Parent.Root.Verbose {
			fmt.Fprintf(opts.Parent.Root.Stdout, "Adding ingredient %s to recipe %s\n", opts.Name, r.Name)
			fmt.Fprintf(opts.Parent.Root.Stdout, "Ingredient: %s, Unit: %s, Percentage: %f, Dependency: %s\n", opts.Name, opts.UnitType, opts.Percentage, opts.Dependency)
		}

		_, err = addRecipeIngredient(ctx, db, AddIngredientParams{
			Name:           opts.Name,
			RecipeID:       r.ID,
			UnitType:       opts.UnitType,
			Percentage:     opts.Percentage,
			Dependency:     opts.Dependency,
			IngredientType: recipe.IngredientType(opts.IngredientType),
		})

		return err
	}
}

type AddIngredientParams struct {
	Name           string
	RecipeID       int64
	UnitType       string
	Percentage     float64
	Dependency     string
	IngredientType recipe.IngredientType
}

// addRecipeIngredient adds a new ingredient to a recipe. If the ingredient
// does not exist, it will be created.
func addRecipeIngredient(ctx context.Context, db *query.Queries, args AddIngredientParams) (*query.RecipeIngredient, error) {
	ingredient, err := db.GetIngredientByName(ctx, args.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get ingredient %s: %w", args.Name, err)
	}

	if ingredient.ID == 0 {
		ingredient, err = db.CreateIngredient(ctx, query.CreateIngredientParams{
			Name:           args.Name,
			IngredientType: args.IngredientType,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create ingredient %s: %w", args.Name, err)
		}
	}

	// Create new recipe ingredient
	params := query.CreateRecipeIngredientParams{
		RecipeID:     args.RecipeID,
		IngredientID: ingredient.ID,
		UnitType:     args.UnitType,
		Percentage:   args.Percentage,
		Dependency:   args.Dependency,
	}

	ri, err := db.CreateRecipeIngredient(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create recipe ingredient: %w", err)
	}

	return &ri, nil
}
