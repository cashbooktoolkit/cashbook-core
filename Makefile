#
# Build the cashbook executable.
#

##############
# Go is on the path somewhere
GO := go

# gb is local to tools
GB := $(shell pwd)/tools/bin/gb

# App/executable being build
APP := cashbook

GO_SRCS = $(shell find src/ -type f -name '*.go')
###############

build : tools/bin/gb bin/$(APP)
	@echo "Build complete"

bin/$(APP) : $(GO_SRCS)
	$(GB) build -ldflags "-X main.build_sha=`git rev-parse HEAD`" $(APP)

test :
	$(GB) test ...

# Tools is a private GOPATH containing build and dev support tools (just one for now)
# We're using github.com/constabulary/gb/
tools/bin/gb :
	@echo "Installing github.com/constabulary/gb/ into tools..."
	GOPATH=$(shell pwd)/tools $(GO) get  github.com/constabulary/gb/...
	@echo "OPTIONAL: Copy tools/bin/* to somewhere on your path for managing dependencies with gb vendor"
	@echo ""

