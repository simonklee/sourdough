-- name: GetRecipe :one
SELECT r.id, r.name
FROM recipes r
WHERE r.id = ? 
LIMIT 1;

-- name: ListRecipes :many
SELECT r.id, r.name
FROM recipes r
ORDER BY r.id DESC;

-- name: CreateRecipe :one
INSERT INTO recipes (name)
VALUES (?)
RETURNING *;

-- name: UpdateRecipe :one
UPDATE recipes 
SET 
  name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteRecipe :exec
DELETE 
FROM recipes
WHERE id = ?; 

-- name: ListRecipesByIngredient :many
SELECT r.id, r.name
FROM recipes r
JOIN recipe_ingredients ri ON ri.recipe_id = r.id
JOIN ingredients i ON i.id = ri.ingredient_id
WHERE i.id = ?
ORDER BY r.id DESC;

-- name: ListRecipeIngredients :many
SELECT ri.id, ri.recipe_id, i.name, ri.unit_type, ri.percentage, ri.dependency, i.ingredient_type
FROM recipe_ingredients ri
JOIN ingredients i ON i.id = ri.ingredient_id
WHERE ri.recipe_id = ?;

-- name: CreateRecipeIngredient :one
INSERT INTO recipe_ingredients (recipe_id, ingredient_id, unit_type, percentage, dependency)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateRecipeIngredient :one
UPDATE recipe_ingredients
SET 
  unit_type = ?,
  percentage = ?,
  dependency = ?,
  ingredient_id = ?
WHERE id = ?
RETURNING *;

-- name: DeleteRecipeIngredient :exec
DELETE
FROM recipe_ingredients
WHERE id = ?;

-- name: GetIngredients :many
SELECT i.id, i.name, i.ingredient_type
FROM ingredients i
ORDER BY i.id DESC;

-- name: GetIngredient :one
SELECT i.id, i.name, i.ingredient_type
FROM ingredients i
WHERE i.id = ?;

-- name: GetIngredientByName :one
SELECT i.id, i.name, i.ingredient_type
FROM ingredients i
WHERE i.name LIKE ?
LIMIT 1;

-- name: CreateIngredient :one
INSERT INTO ingredients (name, ingredient_type)
VALUES (?, ?)
RETURNING *;

-- name: UpdateIngredient :one
UPDATE ingredients
SET 
  name = ?,
  ingredient_type = ?
WHERE id = ?
RETURNING *;

-- name: DeleteIngredient :exec
DELETE
FROM ingredients
WHERE id = ?;
