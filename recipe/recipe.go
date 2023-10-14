package recipe

import (
	"fmt"
	"regexp"
	"strconv"
)

type Unit string

const (
	UnitGrams       Unit = "g"
	UnitKilos       Unit = "kg"
	UnitLitres      Unit = "l"
	UnitMillilitres Unit = "ml"
	UnitTeaspoons   Unit = "tsp"
	UnitTablespoons Unit = "tbsp"
	UnitCups        Unit = "cup"
	UnitPinches     Unit = "pinch"
	UnitHandfuls    Unit = "handful"
)

func ParseUnit(value string) (Unit, error) {
	switch value {
	case "g", "grams":
		return UnitGrams, nil
	case "kg", "kilos", "kilograms":
		return UnitKilos, nil
	case "l", "litres":
		return UnitLitres, nil
	case "ml", "millilitres":
		return UnitMillilitres, nil
	case "tsp", "teaspoons":
		return UnitTeaspoons, nil
	case "tbsp", "tablespoons":
		return UnitTablespoons, nil
	case "cup", "cups":
		return UnitCups, nil
	case "pinch", "pinches":
		return UnitPinches, nil
	case "handful", "handfuls":
		return UnitHandfuls, nil
	default:
		return "", fmt.Errorf("invalid unit: %s", value)
	}
}

type UnitCategory string

const (
	UnitCategoryWeight  UnitCategory = "weight"
	UnitCategoryVolume  UnitCategory = "volume"
	UnitCategoryCount   UnitCategory = "count"
	UnitCategoryUnknown UnitCategory = ""
)

func FromUnitCategory(cat UnitCategory) (Unit, error) {
	switch cat {
	case UnitCategoryWeight:
		return UnitGrams, nil
	case UnitCategoryVolume:
		return UnitLitres, nil
	case UnitCategoryCount:
		return UnitTeaspoons, nil
	default:
		return "", fmt.Errorf("invalid unit type: %s", cat)
	}
}

func DefaultFromUnitCategory(cat UnitCategory) Unit {
	unit, err := FromUnitCategory(cat)
	if err != nil {
		return UnitGrams
	}
	return unit
}

func (u Unit) String() string {
	return string(u)
}

func (u Unit) Format(value float64) string {
	return FormatValue(value, u)
}

func (u Unit) IsWeight() bool {
	switch u {
	case UnitGrams, UnitKilos:
		return true
	default:
		return false
	}
}

func (u Unit) IsVolume() bool {
	switch u {
	case UnitLitres, UnitMillilitres:
		return true
	default:
		return false
	}
}

func (u Unit) IsCount() bool {
	switch u {
	case UnitTeaspoons, UnitTablespoons, UnitCups, UnitPinches, UnitHandfuls:
		return true
	default:
		return false
	}
}

func (u Unit) IsUnknown() bool {
	return !u.IsWeight() && !u.IsVolume() && !u.IsCount()
}

func (u Unit) IsCompatible(other Unit) bool {
	if u.IsUnknown() || other.IsUnknown() {
		return false
	}

	return u.IsWeight() == other.IsWeight() &&
		u.IsVolume() == other.IsVolume() &&
		u.IsCount() == other.IsCount()
}

// Convert converts the given value from one unit to another.
func (u Unit) Convert(value float64, to Unit) (float64, error) {
	return Convert(value, u, to)
}

// ConvertIngredient converts the given value from one unit to another.
func (u Unit) ConvertIngredient(value float64, to Unit, ingredient Kind) (float64, error) {
	return ConvertIngredient(value, u, to, ingredient)
}

var scaleFactor = map[pair]float64{
	{UnitMillilitres, UnitLitres}:    0.001,
	{UnitLitres, UnitMillilitres}:    1000,
	{UnitGrams, UnitKilos}:           0.001,
	{UnitKilos, UnitGrams}:           1000,
	{UnitTeaspoons, UnitTablespoons}: 0.333333,
	{UnitTablespoons, UnitTeaspoons}: 3,
	{UnitCups, UnitLitres}:           0.236588,
	{UnitLitres, UnitCups}:           4.22675,
	{UnitPinches, UnitTeaspoons}:     0.333333,
	{UnitTeaspoons, UnitPinches}:     3,
	{UnitHandfuls, UnitCups}:         0.5,
	{UnitCups, UnitHandfuls}:         2,
}

type pair struct {
	From Unit
	To   Unit
}

// Convert converts the given value from one unit to another.
func Convert(value float64, from Unit, to Unit) (float64, error) {
	if from == to {
		return value, nil
	}

	factor, ok := scaleFactor[pair{from, to}]
	if !ok {
		return 0, fmt.Errorf("unsupported conversion from %s to %s", from, to)
	}

	return value * factor, nil
}

type Kind string

const (
	KindWater     Kind = "water"
	KindFlour     Kind = "flour"
	KindSalt      Kind = "salt"
	KindSugar     Kind = "sugar"
	KindYeast     Kind = "yeast"
	KindOil       Kind = "oil"
	KindEgg       Kind = "egg"
	KindMilk      Kind = "milk"
	KindButter    Kind = "butter"
	KindSourdough Kind = "sourdough"
	KindOther     Kind = ""
)

// Implement from SQL Driver Valuer interface
func (it Kind) Value() (any, error) {
	return string(it), nil
}

// Implement from SQL Scanner interface
func (it *Kind) Scan(src any) error {
	if src == nil {
		*it = KindOther
		return nil
	}

	switch src := src.(type) {
	case string:
		*it = Kind(src)
	default:
		return fmt.Errorf("invalid ingredient type: %v, %T", src, src)
	}

	return nil
}

var ingredientTypeDensity = map[Kind]float64{
	// Water density is 1 g/ml
	KindWater: 1,

	// Flour density ranged from 0.57 ± 0.0 g/ml to 0.75 ± 0.0 g/ml
	// https://www.researchgate.net/figure/Bulk-density-g-ml-of-grains-flour_fig2_364274783
	KindFlour: 0.66,

	// Salt (Sodium chloride) has a density of 2.16 g/ml
	KindSalt: 2.16,

	// Sugar density is 0.85 g/ml
	KindSugar: 0.85,

	// Yeast density is 0.36 g/ml
	KindYeast: 0.36,

	// Oil density is 0.92 g/ml
	KindOil: 0.92,

	// Egg density is 1.03 g/ml
	KindEgg: 1.03,

	// Milk density is 1.03 g/ml
	KindMilk: 1.03,

	// Butter density is 0.96 g/ml
	KindButter: 0.96,

	// Sourdough density is ~0.95 g/ml given the hydration is 120%
	KindSourdough: 0.95,
}

func ConvertIngredient(value float64, from Unit, to Unit, ingredient Kind) (float64, error) {
	// Check unit compatibility first
	if from.IsCompatible(to) {
		return Convert(value, from, to)
	}

	// For incompatible units, check if ingredient density is available
	density, found := ingredientTypeDensity[ingredient]
	if !found {
		return Convert(value, from, to) // fallback to default conversion
	}

	// Normalize value to a standard gram unit first.
	norm, err := normalizeToGrams(value, from, density)
	if err != nil {
		return 0, err
	}

	// Convert the normalized value to the target unit.
	return convertNormalized(norm, to, density)
}

func normalizeToGrams(value float64, from Unit, density float64) (float64, error) {
	switch {
	case from.IsWeight():
		return from.Convert(value, UnitGrams)
	case from.IsVolume():
		normalizedVolume, err := from.Convert(value, UnitMillilitres)
		if err != nil {
			return 0, err
		}
		return normalizedVolume * density, nil
	default:
		return 0, fmt.Errorf("unknown unit: %s", from)
	}
}

func convertNormalized(value float64, to Unit, density float64) (float64, error) {
	switch {
	case to.IsWeight():
		return UnitGrams.Convert(value, to)
	case to.IsVolume():
		volumeInMillilitres := value / density
		return UnitMillilitres.Convert(volumeInMillilitres, to)
	default:
		return 0, fmt.Errorf("unknown unit: %s", to)
	}
}

// Tuple is a value and unit pair.
type Tuple struct {
	Value float64
	Unit  Unit
}

// Appropriate converts the given value to the most appropriate unit.
// It's like a "humanize" function.
//
// For example:
//   - if the value is 1000g, the tuple will be {1, UnitKilos}.
//   - if the value is 1.5l, the tuple will be {1.5, UnitLitres}.
//   - if the value is 1500ml, the tuple will be {1.5, UnitLitres}.
func (u Unit) Appropriate(value float64) Tuple {
	switch u {
	case UnitGrams:
		if value >= 1000 {
			return Tuple{Value: value / 1000, Unit: UnitKilos}
		}
	case UnitKilos:
		if value < 1 {
			return Tuple{Value: value * 1000, Unit: UnitGrams}
		}
	case UnitLitres:
		if value < 1 {
			return Tuple{Value: value * 1000, Unit: UnitMillilitres}
		}
	case UnitMillilitres:
		if value >= 1000 {
			return Tuple{Value: value / 1000, Unit: UnitLitres}
		}
	case UnitTeaspoons:
		if value >= 3 {
			return Tuple{Value: value / 3, Unit: UnitTablespoons}
		}
	case UnitTablespoons:
		if value < 1 {
			return Tuple{Value: value * 3, Unit: UnitTeaspoons}
		}
	case UnitCups:
		if value >= 4 {
			return Tuple{Value: value / 4, Unit: UnitLitres}
		}
	case UnitPinches:
		if value >= 3 {
			return Tuple{Value: value / 3, Unit: UnitTeaspoons}
		}
	}

	return Tuple{Value: value, Unit: u}
}

func (u Unit) Tuple(value float64) Tuple {
	return Tuple{Value: value, Unit: u}
}

func (t Tuple) Appropriate() Tuple {
	return t.Unit.Appropriate(t.Value)
}

func (t Tuple) Convert(to Unit) (Tuple, error) {
	value, err := t.Unit.Convert(t.Value, to)
	if err != nil {
		return Tuple{}, err
	}
	return Tuple{Value: value, Unit: to}, nil
}

func (t Tuple) ConvertIngredient(to Unit, ingredient Kind) (Tuple, error) {
	value, err := ConvertIngredient(t.Value, t.Unit, to, ingredient)
	if err != nil {
		return Tuple{}, err
	}
	return Tuple{Value: value, Unit: to}, nil
}

func (t Tuple) Format() string {
	return FormatValue(t.Value, t.Unit)
}

// FormatValue formats a value with the given unit.
//
// The value is formatted as a string with the following rules:
//
//   - if the unit is a weight, the value is formatted as a decimal with 2
//     decimal places.
//   - if the unit is a volume, the value is formatted as a decimal with 1
//     decimal place.
//   - if the unit is a count, the value is formatted as a decimal with no
//     decimal places.
//   - if the unit is unknown, the value is formatted as a decimal with 2
//     decimal places.
//
// The values unit is appended to the end of the string.
func FormatValue(value float64, unit Unit) string {
	switch {
	case unit.IsWeight():
		return fmt.Sprintf("%.2f%s", value, pluralizeUnit(unit, value))
	case unit.IsVolume():
		return fmt.Sprintf("%.1f%s", value, pluralizeUnit(unit, value))
	case unit.IsCount():
		return fmt.Sprintf("%.0f%s", value, pluralizeUnit(unit, value))
	default:
		return fmt.Sprintf("%.2f%s", value, pluralizeUnit(unit, value))
	}
}

func pluralizeUnit(unit Unit, value float64) string {
	if value == 1 {
		return unit.String()
	}

	switch unit {
	case UnitCups:
		return "cups"
	case UnitPinches:
		return "pinches"
	case UnitHandfuls:
		return "handfuls"
	default:
		return unit.String()
	}
}

// Dependency is a dependency of a recipe.
type Dependency struct {
	Label string
	Value Tuple
}

// Ingredient is an ingredient of a recipe.
type Ingredient struct {
	Name         string
	UnitCategory UnitCategory
	Percentage   float64
	Dependency   string
}

// RecipeIngredient is an ingredient of a recipe template.
type RecipeIngredient struct {
	Name               string
	Kind               Kind
	PreferUnitCategory UnitCategory
	Percentage         float64
	Dependency         string
}

// Recipe is a recipe template.
type Recipe struct {
	Name        string
	Ingredients []RecipeIngredient
}

// PortionIngredient is an ingredient with the amount and unit. It's used to
// represent the ingredients in a portion of a recipe. This is the final result
// of a recipe calculation.
type PortionIngredient struct {
	Name               string
	Kind               Kind
	PreferUnitCategory UnitCategory
	Value              Tuple
}

// ParseDependencies parses a list of dependency strings into a list of
// Dependency structs.
func ParseDependencies(dependencies []string) ([]Dependency, error) {
	var deps []Dependency
	for _, dep := range dependencies {
		d, err := ParseDependency(dep)
		if err != nil {
			return nil, err
		}

		deps = append(deps, d)
	}

	return deps, nil
}

// ParseDependency parses a dependency string into a Dependency struct.
//
// The dependency string should be in the format of:
//
//	<dependency label> <dependency value><dependency unit>
//
// For example:
//
//	"water 1000g"
//	"water 1.5l"
//	"water .5kg"
//
// The dependency label can be anything.
func ParseDependency(dep string) (Dependency, error) {
	matches := depRegexp.FindStringSubmatch(dep)
	if len(matches) != 4 {
		return Dependency{}, fmt.Errorf("invalid dependency string: %s", dep)
	}

	amount, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return Dependency{}, fmt.Errorf("invalid dependency value: %w", err)
	}

	unit, err := ParseUnit(matches[3])
	if err != nil {
		return Dependency{}, fmt.Errorf("invalid dependency unit: %w", err)
	}

	label := matches[1]

	var kind Kind

	switch label {
	case "total_flour":
		kind = KindFlour
	default:
		kind = Kind(label)
	}

	value := Tuple{
		Value: amount,
		Unit:  unit,
	}
	if v, err := value.ConvertIngredient(UnitGrams, kind); err == nil {
		value = v
	}

	return Dependency{
		Label: label,
		Value: value,
	}, nil
}

var depRegexp = regexp.MustCompile(`(?P<label>\w+)\s+(?P<value>\d+\.?\d*)(?P<unit>\w+)`)

func findDependency(dependencies []Dependency, label string) (Dependency, error) {
	for _, dep := range dependencies {
		if dep.Label == label {
			return dep, nil
		}
	}

	return Dependency{}, fmt.Errorf("dependency not found: %s", label)
}

// Calculate calculates the portion ingredients based on the given
// dependencies.
func (r Recipe) Calculate(dependencies []Dependency) ([]PortionIngredient, error) {
	return Calculate(r.Ingredients, dependencies)
}

// Calculate calculates the portion ingredients based on the given templates
// and dependencies.
func Calculate(templates []RecipeIngredient, dependencies []Dependency) ([]PortionIngredient, error) {
	var ingredients []PortionIngredient
	for _, template := range templates {
		dep, err := findDependency(dependencies, template.Dependency)
		if err != nil {
			return nil, fmt.Errorf("failed to find dependency for %s: %w", template.Name, err)
		}

		ingredients = append(ingredients, PortionIngredient{
			Name: template.Name,
			Kind: template.Kind,
			Value: Tuple{
				Value: dep.Value.Value * template.Percentage,
				Unit:  dep.Value.Unit,
			},
			PreferUnitCategory: template.PreferUnitCategory,
		})
	}

	return ingredients, nil
}
