# The port to run the proxy on
PORT=5050

# Expiration in seconds for the items in the cache
EXPIRATION=5

# Address of redis cache
REDIS=redis

# Size of cache (number of items)
SIZE=50

UNAME := $(shell uname)

SED := $(shell which sed)


ifeq ($(UNAME), Linux)
	SEDOPT="-i"
endif

ifeq ($(UNAME), Darwin)
	SEDOPT="-i.bak"
endif

test: test-local Dockerfile docker-compose.yml run-compose

clean: clean-docker clean-compose clean-image

run-compose:
	docker-compose up	

test-local:
	go test -v ./...

Dockerfile:
	cp dockerfile.tmpl Dockerfile

	$(SED) $(SEDOPT) 's/__PORT__/$(PORT)/' Dockerfile
	$(SED) $(SEDOPT) 's/__SIZE__/$(SIZE)/' Dockerfile
	$(SED) $(SEDOPT) 's/__EXPIRATION__/$(EXPIRATION)/' Dockerfile
	$(SED) $(SEDOPT) 's/__REDIS__/$(REDIS)/' Dockerfile

	rm -f Dockerfile.bak


docker-compose.yml:
	cp compose.tmpl docker-compose.yml

	$(SED) $(SEDOPT) 's/__PORT__/$(PORT)/g' docker-compose.yml

	rm -f docker-compose.yml.bak

clean-docker:
	rm -f Dockerfile

clean-compose:
	rm -f docker-compose.yml

clean-image:
	docker rmi -f redisproxy_proxy
