package main

import (
	"testing"
)

func TestConvertIngredient(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		from        Unit
		to          Unit
		ingredient  IngredientType
		expected    float64
		expectError bool
	}{
		{
			name:       "Convert grams to kilograms",
			value:      1000,
			from:       UnitGrams,
			to:         UnitKilos,
			ingredient: IngredientFlour,
			expected:   1,
		},
		{
			name:       "Convert kilograms to grams",
			value:      1,
			from:       UnitKilos,
			to:         UnitGrams,
			ingredient: IngredientFlour,
			expected:   1000,
		},
		{
			name:       "Convert millilitres to litres",
			value:      1000,
			from:       UnitMillilitres,
			to:         UnitLitres,
			ingredient: IngredientWater,
			expected:   1,
		},
		{
			name:       "Convert 1l water to 1000ml",
			value:      1,
			from:       UnitLitres,
			to:         UnitMillilitres,
			ingredient: IngredientWater,
			expected:   1000,
		},
		{
			name:       "Convert 1000g water to 1000ml",
			value:      1000,
			from:       UnitGrams,
			to:         UnitMillilitres,
			ingredient: IngredientWater,
			expected:   1000,
		},
		{
			name:       "Convert 1000ml water to 1kg",
			value:      1000,
			from:       UnitMillilitres,
			to:         UnitKilos,
			ingredient: IngredientWater,
			expected:   1,
		},
		{
			name:       "Convert 1000g flour to litres",
			value:      1000,
			from:       UnitGrams,
			to:         UnitLitres,
			ingredient: IngredientFlour,
			expected:   1.5151515151515151,
		},
		// Salt (Sodium chloride) has a density of 2.16 g/ml
		{
			name:       "Convert 1kg salt to litres",
			value:      1,
			from:       UnitKilos,
			to:         UnitLitres,
			ingredient: IngredientSalt,
			expected:   0.46296296296296297,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ConvertIngredient(test.value, test.from, test.to, test.ingredient)
			if (err != nil) != test.expectError {
				t.Fatalf("expected error: %v, got: %v", test.expectError, err)
			}
			if err == nil && result != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}

func TestParseDependency(t *testing.T) {
	tests := []struct {
		name    string
		dep     string
		want    Dependency
		wantErr bool
	}{
		{
			name: "valid dependency",
			dep:  "water 1000g",
			want: Dependency{
				Label: "water",
				Value: 1000,
				Unit:  UnitGrams,
			},
		},
		{
			name: "valid dependency",
			dep:  "water 1kg",
			want: Dependency{
				Label: "water",
				Value: 1,
				Unit:  UnitKilos,
			},
		},
		{
			name: "valid dependency",
			dep:  "water 1l",
			want: Dependency{
				Label: "water",
				Value: 1,
				Unit:  UnitLitres,
			},
		},
		{
			name: "valid dependency",
			dep:  "water 1.5l",
			want: Dependency{
				Label: "water",
				Value: 1.5,
				Unit:  UnitLitres,
			},
		},
		{
			name: "valid dependency",
			dep:  "water 1.5kg",
			want: Dependency{
				Label: "water",
				Value: 1.5,
				Unit:  UnitKilos,
			},
		},
		{
			name: "valid dependency",
			dep:  "water 1.5g",
			want: Dependency{
				Label: "water",
				Value: 1.5,
				Unit:  UnitGrams,
			},
		},
		{
			name:    "invalid dependency",
			dep:     "water 1.5",
			wantErr: true,
		},
		{
			name:    "invalid dependency",
			dep:     "water 1.5kgg",
			wantErr: true,
		},
		{
			name:    "invalid dependency",
			dep:     "water 1.5ll",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDependency(tt.dep)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDependency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Label != tt.want.Label {
				t.Errorf("ParseDependency() got = %v, want %v", got, tt.want)
			}

			if got.Value != tt.want.Value {
				t.Errorf("ParseDependency() got = %v, want %v", got, tt.want)
			}

			if got.Unit != tt.want.Unit {
				t.Errorf("ParseDependency() got = %v, want %v", got, tt.want)
			}
		})
	}
}
