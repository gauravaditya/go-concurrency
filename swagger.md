How to add swagger documentation in golang

- Add swago dependency
	> ```shell
	> go install github.com/swaggo/swag/cmd/swag@latest
	> ```
- Add fiber swagger dependency
	> ```shell
	> go get github.com/gofiber/swagger
	> ```
- Add swagger handler and imports in the file
	>  ```go
	>  _ "github.com/gauravaditya/go-monorepo/docs/core" // path to swagger docs json
	>	
	>  // Swagger UI endpoint
	>  app.Get("/swagger/*", swagger.HandlerDefault)
	>  ```


Make helper commands to generate documentation:
```shell

SERVICES=core event consumer
# path to go file that contains the swagger documentation comments
swagger_api_file_dirs=$(foreach route,$(shell fd 'routes.go' internal/ | uniq), $(dir $(route)))
# generate swagger targets for a monorepo setup
SWAGGER_DOCS_SERVICES_TARGET=$(foreach service,$(SERVICES),docs/$(service))

run/docs/%:
	@echo "Generating Swagger for $(@D) -> $(@F)... \n"

	$(eval api_path := $(subst run/docs/,internal/,$@))
	$(eval route_paths := $(filter internal/$(@F)/%,$(swagger_api_file_dirs)))
	$(foreach route,$(route_paths),swag init -o $(subst internal/,docs/,$(route)) -d $(route),$(api_path) -g routes.go -pdl 1 --parseInternal;)

$(SWAGGER_DOCS_SERVICES_TARGET): %: run/%

docs: $(SWAGGER_DOCS_SERVICES_TARGET)
```