### Pattern 1:
dynamic targets generations for repeated actions per service or target using target definition

```makefile
SERVICES=core event consumer # list of services

build: $(SERVICES) ## Build all services

$(SERVICES): %: build/%

build/%: ## build the target service
	@mkdir -p bin
	@echo "Building service $* with version $(VERSION)..."

	@CGO_ENABLED=$(CGO_ENABLED) $(GO_BIN) build \
	-ldflags "$(BUILD_VERSION_LD_FLAGS)" \
	-o bin/$* \
	./cmd/$*
```

### Pattern 2:
dynamic target generation for repeated actions per service or target using variable definitions
[[Helper Methods|@Helper Methods]]
```makefile
SERVICES=core event consumer # list of services
BUILD_SERVICES=$(call gen_targets,build)

build: $(BUILD_SERVICES) ## Build all services

build/%: ## build the target service
	@mkdir -p bin
	@echo "Building service $* with version $(VERSION)..."

	@CGO_ENABLED=$(CGO_ENABLED) $(GO_BIN) build \
	-ldflags "$(BUILD_VERSION_LD_FLAGS)" \
	-o bin/$* \
	./cmd/$*
```



HELPER METHODS

### Generate Targets:
`gen_targets`
```makefile
define gen_targets
$(foreach service,$(SERVICES),$(1)/$(service))
endef
```

### Help message for dynamic targets:
Need to be updated as per the required targets and messages
```makefile
define help_message
	$(eval a1= $(shell echo $(1) | cut -d'/' -f1))
	$(eval a2= $(shell echo $(1) | cut -d'/' -f2))

	$(if $(findstring backend/,$(1)),@printf "$(BLUE)%-40s$(NC) %s\n" "$(1)" " run terraform $(a2) in the $(a1) directory",)
	$(if $(findstring dev/,$(1)),@printf "$(BLUE)%-40s$(NC) %s\n" "$(1)" " run terraform $(a2) for the $(a1) env",)
	$(if $(findstring workspace/,$(1)),@printf "$(BLUE)%-40s$(NC) %s\n" "$(1)" " select workspace $(a2)",)
	$(if $(findstring test/,$(1)),@printf "$(BLUE)%-40s$(NC) %s\n" "$(1)" " run tests for the package $(a2)",)
	@echo $(1) $(2)
endef
```