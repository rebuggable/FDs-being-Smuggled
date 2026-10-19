include Makefiles/binaries.mk

define actuateEngine
	$(DOCKER_BIN) network inspect mule-network >/dev/null 2>&1 || $(DOCKER_BIN) network create mule-network ; \
	$(DOCKER_BIN) run -d \
		--privileged \
		--name mule-engine \
		--network mule-network \
		--network-alias docker \
		-e DOCKER_TLS_CERTDIR=/certs \
		-v mule-docker-certs-ca:/certs/ca \
		-v mule-docker-certs-client:/certs/client \
		a-mule-engine
endef
