OBJS = $(shell find cmd -mindepth 1 -type d -execdir printf '%s\n' {} +)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
BASEPKG = github.com/allinbits/emeris-cns-server
EXTRAFLAGS :=

.PHONY: $(OBJS) clean generate-swagger

all: $(OBJS)

clean:
	@rm -rf build cns/docs/swagger.* cns/docs/docs.go

generate-swagger:
	go generate ${BASEPKG}/cns/docs
	@rm cns/docs/docs.go

generate-mocks:
	@rm mocks/*.go || true
	mockery --srcpkg sigs.k8s.io/controller-runtime/pkg/client --name Client

$(OBJS):
	go build -o build/$@ -ldflags='-X main.Version=${BRANCH}-${COMMIT}' ${EXTRAFLAGS} ${BASEPKG}/cmd/$@
