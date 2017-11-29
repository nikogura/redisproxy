#**********************************************************
# The following values can be tweaked in the docker test configuration.
#
# By default, running 'make test' will spin up 2 instances.  A generic redis container and the actual proxy.  The generic redis will be set up with some minimal test data.

# Testing and verification wise, it will test itself, hitting the cache repeatedly.  You should be able to see that it's working as intended.  The requirement of no software installation prevents a real test harness.  Sorry
#
# I'm figuring you have more advanced tests you'll want to run than the little bit I have here.
#
# I've also provided some apache benchmark tests.  If you're on linux, you probably have it available.  If you're on a Mac 'brew install apache-httpd' will give you the tool, and then you can run 'make ab-nice' 'make ab-nasty' and 'make ab-omg-what-are-you-trying-to-prove' if you so desire.  The latter is highly dependent on your hardware.
#
# If you want to see the cache limiting behavior happening before your eyes, you'll need 2 terminals.  
#
# If that's your desire, Change FOREGROUND=false to FOREGROUND=true.
#
# Then in terminal 1 run 'make test'.  In terminal 2 run 'make cache-test'.

#**********************************************************
## Start Configuration

# The port to run the proxy on
PORT=5050

# Expiration in seconds for the items in the cache
EXPIRATION=5

# Address of redis cache
REDIS=redis

# Size of cache (number of items)
SIZE=3

# Stay attached to the terminal (Easier to see what's happening with the proxy that way)
FOREGROUND=false


## End Configuration
#**********************************************************

UNAME := $(shell uname)

SED := $(shell which sed)

RCLI := $(shell which redis-cli)

ifeq ($(UNAME), Linux)
	SEDOPT="-i"
endif

ifeq ($(UNAME), Darwin)
	SEDOPT="-i.bak"
endif

ifeq ($(RCLI), '')
	RCLI="docker exec redisproxy_proxy_1 redis-cli -h redis"

endif

TEST_TARGETS = Dockerfile docker-compose.yml run-compose

ifeq ($(FOREGROUND), true)
	DCOPTS= up
else
	DCOPTS= up -d
	TEST_TARGETS += cache-test
endif

test:  $(TEST_TARGETS)

clean: clean-docker clean-compose clean-image

run-compose:
	docker-compose $(DCOPTS)

Dockerfile:
	@cp dockerfile.tmpl Dockerfile

	@$(SED) $(SEDOPT) 's/__PORT__/$(PORT)/' Dockerfile
	@$(SED) $(SEDOPT) 's/__SIZE__/$(SIZE)/' Dockerfile
	@$(SED) $(SEDOPT) 's/__EXPIRATION__/$(EXPIRATION)/' Dockerfile
	@$(SED) $(SEDOPT) 's/__REDIS__/$(REDIS)/' Dockerfile

	@rm -f Dockerfile.bak


docker-compose.yml:
	@cp compose.tmpl docker-compose.yml

	@$(SED) $(SEDOPT) 's/__PORT__/$(PORT)/g' docker-compose.yml

	@rm -f docker-compose.yml.bak

clean-docker:
	@rm -f Dockerfile

clean-compose:
	@docker-compose down 
	@rm -f docker-compose.yml

clean-image:
	@docker rmi -f redisproxy_proxy

cache-test:
	@echo "Running Integration Tests"
	@echo ""
	@echo "GET http://localhost:$(PORT)/foo expecting \"foo\""
	@./integration_test.sh foo $(PORT)
	@echo ""
	@echo "GET http://localhost:$(PORT)/bar expecting \"bar\""
	@./integration_test.sh bar $(PORT)
	@echo ""
	@echo "GET http://localhost:$(PORT)/wip expecting \"wip\""
	@./integration_test.sh wip $(PORT)
	@echo ""
	@echo "GET http://localhost:$(PORT)/zoz expecting \"zoz\""
	@./integration_test.sh zoz $(PORT)
	@echo ""
	@echo "GET http://localhost:$(PORT)/bar expecting \"bar\""
	@./integration_test.sh bar $(PORT)
	@echo ""
	@echo "GET http://localhost:$(PORT)/foo expecting \"foo\""
	@./integration_test.sh foo $(PORT)
	@echo ""
	@echo "GET http://localhost:$(PORT)/wip expecting \"wip\""
	@./integration_test.sh wip $(PORT)
	@echo ""
	@echo "GET http://localhost:$(PORT)/zoz expecting \"zoz\""
	@./integration_test.sh zoz $(PORT)
	@echo ""
	@echo "Done"

ab-nice:
	ab -n 100 -c 10 http://localhost:$(PORT)/foo

ab-nasty:
	ab -n 1000 -c 30 http://localhost:$(PORT)/foo

ab-omg-what-are-you-trying-to-prove:
	ab -n 1000 -c 50 http://localhost:5050/foo
