install: build
	cp sourdough ~/.local/bin/sourdough

build: generate test
  go build -o sourdough -tags fts5 .

test:
  @go test ./...

generate:
  sqlc generate

fmt:
	@gofumpt -w $(fd -e go)

bootstrap: install
  rm -rf ~/.config/sourdough

  sourdough add --name "Artisan Light Buns"
  sourdough add ingredient -p .875 --name 'White Flour' --recipe 1 --dependency total_flour -k flour
  sourdough add ingredient -p .125 --name 'Whole Grain Flour' --recipe 1 --dependency total_flour -k flour
  sourdough add ingredient -p .15 --name 'Sourdough Starter' --recipe 1 --dependency total_flour -k sourdough
  sourdough add ingredient -p .77 --name 'Water' --recipe 1 --dependency total_flour -k water
  sourdough add ingredient -p .018 --name 'Salt' --recipe 1 --dependency total_flour -k salt

  sourdough add --name "Balanced Blend Buns"
  sourdough add ingredient -p .5 --name 'White Flour' --recipe 2 --dependency total_flour -k flour
  sourdough add ingredient -p .5 --name 'Whole Grain Flour' --recipe 2 --dependency total_flour -k flour
  sourdough add ingredient -p .15 --name 'Sourdough Starter' --recipe 2 --dependency total_flour -k sourdough
  sourdough add ingredient -p .77 --name 'Water' --recipe 2 --dependency total_flour -k water
  sourdough add ingredient -p .018 --name 'Salt' --recipe 2 --dependency total_flour -k salt

tools:
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@main
