# sourdough

**sourdough** is a CLI tool for managing recipes and facilitating the baking of
sourdough bread and buns. It comes with a suite of subcommands to handle recipe
listing, addition, and viewing. Whether you're a beginner or seasoned sourdough
baker, **sourdough** streamlines your baking process.

## Installation

```sh
go install -tags fts5 github.com/simonklee/sourdough@latest
```

or clone the repository and run

```sh
just install
```

## Usage

Invoke **sourdough** in your terminal as follows:

```bash
sourdough [flags] <subcommand> ...
```

### Subcommands

- `list` - List all saved recipes.
- `add` - Add a new recipe.
- `view` - View a specific recipe.

### Flags

Global flags available across subcommands:

- `-v, --verbose` - Log verbose output.
- `--term` - Output in terminal format.
- `--html` - Output in HTML format.
- `--markdown` - Output in Markdown format.

## Commands Detail

### list

Lists all saved recipes.

```bash
sourdough list [flags]
```

### add

Adds a new recipe. It contains a subcommand to add a new ingredient to a recipe.

```bash
sourdough add [flags]
```

#### Subcommands

- `ingredient` - Add a new ingredient to a recipe.

```bash
sourdough add ingredient <recipe> [flags]
```

#### Flags (add)

- `-n, --name STRING` - Name of the recipe.

#### Flags (ingredient)

- `-r, --recipe INT` - Recipe ID (default: 0).
- `-n, --name STRING` - Name of the ingredient.
- `-u, --unit STRING` - Preferred output unit category (default: weight).
- `-p, --percentage FLOAT64` - Percentage of the ingredient (default: 0).
- `-d, --dependency STRING` - Dependency of the ingredient.
- `-k, --kind STRING` - Kind of the ingredient.

### view

View a specific recipe. By default, the recipe template (relative values) will be displayed. For viewing the amount of ingredients required for a specific portion, use the `--dependencies` flag.

```bash
sourdough view [flags] <recipe>
```

#### Flags (view)

- `-d, --dependency STRING` - Dependency (e.g. --dependency "total_flour 450g").
- `-i, --ingredients` - Only display ingredients.

## Examples

Add a recipe:

```bash
sourdough add --name "Balanced Blend Buns"
sourdough add ingredient -p .5 --name 'White Flour' --recipe 1 --dependency total_flour -k flour
sourdough add ingredient -p .5 --name 'Whole Grain Flour' --recipe 1 --dependency total_flour -k flour
sourdough add ingredient -p .15 --name 'Sourdough Starter' --recipe 1 --dependency total_flour -k sourdough
sourdough add ingredient -p .77 --name 'Water' --recipe 1 --dependency total_flour -k water
sourdough add ingredient -p .018 --name 'Salt' --recipe 1 --dependency total_flour -k salt
```

Viewing a recipe with specified dependency:

```bash
$ sourdough view --dependency "total_flour 450g" 1

┌─────────────────────────────────────────────────────────────────────────────────┐
│ Recipe: Balanced Blend Buns                                                     │
├────┬───────────────────┬───────────┬─────────────────┬────────────┬─────────────┤
│  # │ INGREDIENT        │ KIND      │ PREFER CATEGORY │ PERCENTAGE │ DEPENDENCY  │
├────┼───────────────────┼───────────┼─────────────────┼────────────┼─────────────┤
│  6 │ White Flour       │ flour     │ weight          │        0.5 │ total_flour │
│  7 │ Whole Grain Flour │ flour     │ weight          │        0.5 │ total_flour │
│  8 │ Sourdough Starter │ sourdough │ weight          │       0.15 │ total_flour │
│  9 │ Water             │ water     │ weight          │       0.77 │ total_flour │
│ 10 │ Salt              │ salt      │ weight          │      0.018 │ total_flour │
└────┴───────────────────┴───────────┴─────────────────┴────────────┴─────────────┘
┌─────────────────────────────────┐
│ Ingredients for: Balanced Blend │
│ Buns                            │
├───┬───────────────────┬─────────┤
│ # │ INGREDIENT        │ AMOUNT  │
├───┼───────────────────┼─────────┤
│ 1 │ White Flour       │ 225.00g │
│ 2 │ Whole Grain Flour │ 225.00g │
│ 3 │ Sourdough Starter │  67.50g │
│ 4 │ Water             │ 346.50g │
│ 5 │ Salt              │   8.10g │
└───┴───────────────────┴─────────┘
```

## License

**sourdough** is distributed under the MIT License. See `LICENSE` for more
information.
