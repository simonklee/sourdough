/* DROP TABLE IF EXISTS recipe_ingredients; */
/* DROP TABLE IF EXISTS recipes; */
/* DROP TABLE IF EXISTS ingredients; */
CREATE TABLE IF NOT EXISTS ingredients (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  kind TEXT NULL
);

CREATE TABLE IF NOT EXISTS recipes (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT CHECK (LENGTH(name) > 0) NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_ingredients (
  id INTEGER NOT NULL PRIMARY KEY,
  recipe_id INTEGER NOT NULL,
  ingredient_id INTEGER NOT NULL,
  prefer_unit_category TEXT CHECK (prefer_unit_category IN ('weight', 'volume', 'count', 'teaspoon')) NOT NULL,
  percentage REAL NOT NULL CHECK ((
    percentage BETWEEN 0 AND 1
  ) AND (
    percentage > 0
  )),
  dependency TEXT NOT NULL CHECK (dependency IN ('total_flour')),
  FOREIGN KEY (recipe_id) REFERENCES recipes (
    id
  ) ON DELETE CASCADE,
  FOREIGN KEY (ingredient_id) REFERENCES ingredients (
    id
  ) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_recipe_id ON recipe_ingredients (recipe_id);

CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_ingredient_id ON recipe_ingredients (ingredient_id);