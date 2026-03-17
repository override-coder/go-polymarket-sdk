###############################################################################
###                                Linting                                  ###
###############################################################################

golangci_lint_cmd=golangci-lint
golangci_version=v1.60.3

lint-install:
	@echo "--> Installing golangci-lint $(golangci_version)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)

lint:
	echo "--> Running linter"
	$(MAKE) lint-install
	$(golangci_lint_cmd) run --build-tags=$(GO_BUILD) --fix --out-format=tab --timeout=10m
	@if [ $$(find . -name '*.go' -type f | xargs grep 'nolint\|#nosec' | wc -l) -ne 24 ]; then \
        echo $$(find . -name '*.go' -type f | xargs grep 'nolint\|#nosec' | wc -l); \
		echo "--> increase or decrease nolint, please recheck them"; \
		echo "--> list nolint: \`find . -name '*.go' -type f | xargs grep 'nolint\|#nosec'\`"; exit 1;\
	fi

format: lint

.PHONY: format lint

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	@echo "--> Running tests"
	go test -mod=readonly ./...

test-count:
	go test -mod=readonly -cpu 1 -count 1 -cover ./... | grep -v 'types\|cli\|no test files'

.PHONY: test
