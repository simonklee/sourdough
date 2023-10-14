// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package query

import (
	"context"
)

type Querier interface {
	CreateIngredient(ctx context.Context, arg CreateIngredientParams) (Ingredient, error)
	CreateRecipe(ctx context.Context, name string) (Recipe, error)
	CreateRecipeIngredient(ctx context.Context, arg CreateRecipeIngredientParams) (RecipeIngredient, error)
	DeleteIngredient(ctx context.Context, id int64) error
	DeleteRecipe(ctx context.Context, id int64) error
	DeleteRecipeIngredient(ctx context.Context, id int64) error
	GetIngredient(ctx context.Context, id int64) (Ingredient, error)
	GetIngredientByName(ctx context.Context, name string) (Ingredient, error)
	GetIngredients(ctx context.Context) ([]Ingredient, error)
	GetRecipe(ctx context.Context, id int64) (Recipe, error)
	ListRecipeIngredients(ctx context.Context, recipeID int64) ([]ListRecipeIngredientsRow, error)
	ListRecipes(ctx context.Context) ([]Recipe, error)
	ListRecipesByIngredient(ctx context.Context, id int64) ([]Recipe, error)
	UpdateIngredient(ctx context.Context, arg UpdateIngredientParams) (Ingredient, error)
	UpdateRecipe(ctx context.Context, arg UpdateRecipeParams) (Recipe, error)
	UpdateRecipeIngredient(ctx context.Context, arg UpdateRecipeIngredientParams) (RecipeIngredient, error)
}

var _ Querier = (*Queries)(nil)
