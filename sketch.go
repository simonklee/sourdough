package main

import (
	"fmt"
	"os"
	"strconv"
)

type UnitType int

const (
	Weight UnitType = iota
	Volume
	Fixed
)

type Node struct {
	label        string
	unitType     UnitType
	value        float64
	percentage   float64
	dependencies []*Node
}

type Recipe struct {
	nodes []*Node
}

func mains() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <Artisan|Blend> <Water|Flour> <amount>")
		return
	}

	variant := os.Args[1]
	param := os.Args[2]
	amount, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Println("Invalid amount. Please enter a valid number.")
		return
	}

	flour := &Node{label: "Total Flour", unitType: Weight}
	// Define the recipe DAG for Artisan Light (change percentages for Balanced Blend)
	var recipe Recipe
	if variant == "Artisan" {
		recipe = Recipe{
			nodes: []*Node{
				flour,
				{label: "Whole Grain Flour", unitType: Weight, percentage: 0.125, dependencies: []*Node{flour}},
				{label: "White Flour", unitType: Weight, percentage: 0.875, dependencies: []*Node{flour}},
				{label: "Salt", unitType: Weight, percentage: 0.018, dependencies: []*Node{flour}},
				{label: "Sourdough Starter", unitType: Weight, percentage: 0.15, dependencies: []*Node{flour}},
				{label: "Water", unitType: Weight, percentage: 0.77, dependencies: []*Node{flour}},
			},
		}
	} else if variant == "Blend" {
		recipe = Recipe{
			nodes: []*Node{
				flour,
				{label: "Whole Grain Flour", unitType: Weight, percentage: 0.5, dependencies: []*Node{flour}},
				{label: "White Flour", unitType: Weight, percentage: 0.5, dependencies: []*Node{flour}},
				{label: "Salt", unitType: Weight, percentage: 0.018, dependencies: []*Node{flour}},
				{label: "Sourdough Starter", unitType: Weight, percentage: 0.15, dependencies: []*Node{flour}},
				{label: "Water", unitType: Weight, percentage: 0.77, dependencies: []*Node{flour}},
			},
		}
	} else {
		fmt.Println("Invalid variant. Use either Artisan or Blend.")
		return
	}

	if param == "Water" {
		flour.value = amount / recipe.nodes[5].percentage // assuming water is the last node in the recipe
	} else if param == "Flour" {
		flour.value = amount
	} else {
		fmt.Println("Invalid parameter. Use either Water or Flour.")
		return
	}

	// calculatedRecipe := calculate(&recipe)
	// displayRecipe(&calculatedRecipe)
}
