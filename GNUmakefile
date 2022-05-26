TEST?=netbox/*.go
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
version=0.0.1
provider=netbox
provider_path=/opt/terraform/providers/$(provider)

export NETBOX_VERSION=v3.1.9
export NETBOX_SERVER_URL=http://localhost:8001
export NETBOX_API_TOKEN=0123456789abcdef0123456789abcdef01234567
export NETBOX_TOKEN=$(NETBOX_API_TOKEN)


build_macos:
	go build -o terraform-provider-$(provider)_v$(version)

	mkdir -p $(provider_path)
	mv ./terraform-provider-$(provider)_v$(version) $(provider_path)

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc: docker-up
	@echo "⌛ Startup acceptance tests on $(NETBOX_SERVER_URL)"
	TF_ACC=1 go test -v -cover $(TEST)

.PHONY: test
test: 
	go test $(TEST) $(TESTARGS) -timeout=120s -parallel=4 -cover

# Run dockerized Netbox for acceptance testing
.PHONY: docker-up
docker-up: 
	@echo "⌛ Startup Netbox $(NETBOX_VERSION) and wait for service to become ready"	
	docker-compose -f docker/docker-compose.yml up --build wait
	docker-compose -f docker/docker-compose.yml logs
	@echo "🚀 Netbox is up and running!"

.PHONY: docker-logs
docker-logs:
	docker-compose -f docker/docker-compose.yml logs
	
.PHONY: docker-down
docker-down: 
	docker-compose -f docker/docker-compose.yml down

#! Development 
# The following make goals are only for local usage 
.PHONY: fmt
fmt:
	go fmt
	go fmt netbox/*.go
