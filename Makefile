DOCKER_IMAGE=gcr.io/moonrhythm-containers/webserver

docker:
	buildctl build \
		--frontend dockerfile.v0 \
		--local dockerfile=. \
		--local context=. \
		--output type=image,name=$(DOCKER_IMAGE),push=true
