TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=postmark

default: build

all: fmt test vet errcheck build

build: fmtcheck
	go installMac

installMac:
	rm -Rf ~/.terraform.d/plugins/hashicorp.com/edu/postmark/terraform-provider-postmark
	go build -o terraform-provider-postmark
	mv terraform-provider-postmark ~/.terraform.d/plugins/hashicorp.com/edu/postmark

installLinux:
	rm -Rf ~/.terraform.d/plugins/hashicorp.com/edu/postmark/0.2/$OS_ARCH
	go build -o terraform-provider-postmark
	export OS_ARCH="$(go env GOHOSTOS)_$(go env GOHOSTARCH)"
	mkdir -p ~/.terraform.d/plugins/hashicorp.com/edu/postmark/0.2/$OS_ARCH
	mv terraform-provider-postmark ~/.terraform.d/plugins/hashicorp.com/edu/postmark/0.2/$OS_ARCH

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile