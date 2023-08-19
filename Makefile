DOCKER_REGISTRY = ghcr.io
DOCKER_NAMESPACE = darchie4
DOCKER_IMAGE = devops-hand-in-g09a

check-vars:
	@[ "${PACKAGE}" ] || ( echo "PACKAGE is not set"; exit 1 )
	@[ "${VERSION}" ] || ( echo "VERSION is not set"; exit 1 )
	@[ -d "${PACKAGE}" ] || ( echo "package '${PACKAGE}' does not exist"; exit 1 )

publish: check-vars
	docker build -t ${DOCKER_REGISTRY}/${DOCKER_NAMESPACE}/${DOCKER_IMAGE}-${PACKAGE}:${VERSION} ${PACKAGE}
	docker push ${DOCKER_REGISTRY}/${DOCKER_NAMESPACE}/${DOCKER_IMAGE}-${PACKAGE}:${VERSION}

build:
	@go build frontend backend

test:
	@go test frontend backend