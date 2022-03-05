include help.mk

.PHONY: install test
.DEFAULT_GOAL := help

HUB_USER 	 = esequielvirtuoso
HUB_REPO     = oauth-go-lib
BUILD        = $(shell git rev-parse --short HEAD)
NAME         = $(shell basename $(CURDIR))
IMAGE        = $(HUB_USER)/$(HUB_REPO):$(BUILD)

install: clean ##@dev Download dependencies via go mod.
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor

test: ##@check Run tests and coverage.
	docker build --progress=plain \
		--tag $(IMAGE) \
		--target=test \
		--file=Dockerfile .

	-mkdir coverage
	docker create --name $(NAME)-$(BUILD) $(IMAGE)
	docker cp $(NAME)-$(BUILD):/index.html ./coverage/.
	docker rm -vf $(NAME)-$(BUILD)
