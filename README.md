# redisproxy
[![Circle CI](https://circleci.com/gh/nikogura/redisproxy.svg?style=shield)](https://circleci.com/gh/nikogura/gomason)

[![Go Report Card](https://goreportcard.com/badge/github.com/nikogura/redisproxy)](https://goreportcard.com/report/github.com/nikogura/redisproxy)

An experimental LRU Cache for Redis in golang and docker.

If I were doing this for real, I'd probably just put Nginx in front of Redis and get on with my life.

Still, it's fun, and informative to do it the hard way on occasion.

# Architecture Overview
Redisproxy implements a simple HTTP interface over the top of a very simple LRU cache which is itself a front end for Redis.

The test demonstrator is implemented in a pair of docker containers, built and run by the 'docker-compose' command.

It's intended for use on Linux or Mac machines.  No attempt at Windows compatability has been made at this time.

## Packages

There is a single package hierarchy under *github.com/nikogura/redisproxy/proxy*.  Within that package you'll find subpackages for *cache*, *cmd*, and *service*.

Within each package you will find files of the pattern:

* foo.go
* foo_test.go
* foo_fixtures.go

The \*_test.go is a normal Go pattern for tests.  The \*_fixtures.go file is my own convention.  I prefer to keep the actual test data out of the tests as much as possible and keep them in their own place.

Tests are Code, and generally you want to keep Code and Data separate.  This is a goal.  It is not always practical to follow to its full extension.

### Cache

The cache package contains the cache itself, and the code for entries within the cache.

### Service

The service package contains the code that runs the actual http proxy service and hosts the cache.

### Cmd

The cmd package is a built in feature of the Cobra command framework.  I used Cobra because it's clean, easy, saves time, and generally does a whiz-bang job of making not only command line parsing easy, but also making it easy to have useful and accurate help messages.

# Requirements

* Docker or Docker for Mac

* Access to DockerHub (to pull images)

* Make

* Curl

* Bash

# What it does

In English, that means you make requests of the redisproxy app, and it queries Redis on your behalf, caching the results.  The cache is limited in the number of entries it can hold.  When that limit is exceeded, the Least Recently Used entry in the cache is purged to make room.

Entries have a configured expiration time, and even if they are found in the cache after their time is up, a new value will be fetched from the upstream Redis server.

The containers attach to the host OS on ports 5050 (proxy application) and 6379 (redis).  If those ports are in use, the default config will fail.

For simplicity's sake, I didn't implement a background garbage collector-like cache purging mechanism as I've seen some do in this case.  Instead, as the last part of a fetch, the cache size is measured, and the oldest entry in the cache is purged if the cache is found to be too large.

Since a fetch is performed, and then the cache size is checked, and reduced, the maximum number of items in the cache at any one time will be the configured maximum plus one.  The 'plus one state' will be momentary, but must be kept in mind if this were adapted to a space-restricted environment.

# Algorithmic Complexity

According to 'gocyclo', it's 100%.  I'm not sure what that's really worth however.

Determining whether an item is in the cache aught to be constant time, as it's just a hash lookup.  Likewise retrieving a cached element is quick, as it's just a hash lookup

Fetch of a non-cached element is of course going to be dependent on the internals of the Redis client and the network.  That won't be constant, but we're not really trying to rate that part of the problem.

Determining the age of an element should likewise be pretty quick.

The required LRU functionality's complexity is going to be entirely dependent on how "container/list" is implemented.  Given that it's a Go builtin, I would expect it's fairly fast.

Due to the purge mechanism, once the cache fills to capacity, there will be an additional overhead of 2 remove operations on every fetch.  At the scale this demo is intended to run at, that was judged to be an acceptable trade off for not implementing a periodic background expiration purger routine.

At a greater scale, a background purger that proactively gets rid of the stale entries might be just the ticket.

# Configuration

The following values may be configured in the Makefile for testing:

* The port to run the proxy on *default: 5050*

* Expiration in seconds for the items in the cache *default: 5*

* Address of redis cache *default: redis*

* Size of cache (number of items) *default: 3* (makes the cache size handling easy to test)

* Whether the containers stay attached in the foreground. *default: false*


# Running

If it is your pleasure to actually build, install, and run the proxy directly, you may do so via:

    go install
    
If your $GOPATH/bin is in your $PATH, you can run the proxy via:

    redisproxy run -c (SIZE) -e (EXPIRATION) -p (PORT) -r (REDIS)
    
If you run into trouble, run:

    redisproxy help
    
# Testing

## One Click Validation

To test, run:

    make test
    
2 containers will spin up.  A generic redis container and the actual proxy. the generic redis container will be filled with some very limited test data.
  
Unit tests for the *cache* package are run as part of the proxy build.  Tests for the *service* package are run as well.  Basically the same tests are performed on the cache first by itself in memory, and then through the http service.  Technically that makes them integration tests, not unit tests, but who am I to quibble?  I prefer to test the pieces individually, and then again in groups.  The more help I can give myself in isolating the exact point of the error, the faster the fixes will be.

You'll see the results on the screen as they come out.  If one were to bomb, ```make test``` will fail.

Full blown integration tests are performed against the container via the Makefile.  You should see TAP (Test Anything Protocol) output for those.

Obviously, if you prefer to test against an existing cache, modify the REDIS in the makefile.  If you don't add a port, port 6379 will be assumed.

Testing and verification wise, it will test itself, hitting the cache repeatedly to verify that the cache size limitation is implemented.  You should be able to see that it's working as intended.  

The requirement for 'no additonal software' puts the kybosh on any proper test fixture.  Sorry.

## Go Test

Assuming you have a fully setup GOPATH, when run from the checked out [git repository](https://github.com/nikogura/redisproxy), the following will allow you to test the code:

    cd /place/where/you/checked/out/the/code

    go get github.com/kardianos/govendor
    
    govendor sync
    
    go test -v ./...
    
## Gomason Testing

Gomason is a tool for doing local tests in a clean, ephemeral GOPATH.  Gomason at this time requires you to use SSH keys to connect with Github. ('git clone git@github.com:nikogura/redisproxy.git' rather than 'git clone https://github.com/nikogura/redisproxy.git')

Install it via:

    go get github.com/nikogura/gomason
    
    go install github.com/nikogura/gomason
    
Test via:

    cd /place/where/you/checked/out/the/code
    
    gomason test

## Cache Limit Validation
If you want to see the cache limiting behavior happening before your eyes, you'll need 2 terminals.  

If that's your desire, change ```FOREGROUND=false``` to ```FOREGROUND=true``` in the Makefile.

Then in terminal 1 run:

    make test
      
In terminal 2 run:

    make cache-test
    
You will, of course need to use a ctrl-C in terminal 1 to end the test.

## Concurrent Validation

I'm figuring you have more advanced tests you'll want to run than the little bit I have here.  The requirement of 'no additonal software' is also a little rough when we're talking multiple concurrent tests.  Most browsers don't do that.

To that end, I've provided some apache benchmark tests.  If you're on Linux, you probably have it available.  If you're on a Mac 'brew install apache-httpd' will give you the tool.

To run them, first change ```FOREGROUND=false``` to ```FOREGROUND=true``` in the Makefile.

For a nice concurrent test, run:

    make ab-nice
     
For a nastier concurrent test, run:

    make ab-nasty
    
If you really want to be a jerk, run:

    make ab-omg-what-are-you-trying-to-prove

The latter is highly dependent on your hardware.  

Again, I'm guessing you have more advanced test harnesses for verification.  I've provided a humble mechanism for the case where you do not.  It does require some setup, but them's the breaks as they say.

# Time to Completion

This took longer than would be ideal, but it forced me to dig deeper into some Go constructs than I've been before.  The first time is always slow.

* The cache itself and it's attendant fixtures and tests took about 3 hours.

* The http server, which I knew about, but hadn't messed with before, took another hour or so.  Refactoring it to be able to use it in in-code integration tests took another hour.  It wasn't purely necessary, but I'm happier with the result.

* The dockerfile and docker compose stuff was maybe 30 minutes.

* Re-educating myself on the wild world of 'make' burned a good hour.  It's been a long time.  My IDE is configured to use spaces rather than tabs, and OMG is that a pain in the butt when dealing with Make.  Good ol' vim to the rescue, except I had that configured similarly.  Had to run down a hotkey combo to force a tab character.

* I burned another couple hours chasing down a rabbit hole on channels and such as I could *not* work out how to get the multiple goroutines of the http server to share a single cache.  No sir.  Not for nothing.  Channels, selects, you name it.  I tried every frigging combination under the sun that would compile.

I could do channels that blocked, channels that didn't block, but didn't wait for results, every failure condition I could imagine, but not the one friggin result I needed.  

I combed the internet looking for examples, I invented new ways to phrase my problem.  Finally I found someone else's project that was close enough to the simplicity I craved, and got very angry.  He seemed to be able to pull off what was eluding me.  He got something that looked like my first attempt to work before I went down the channel rabbit hole.  

Turned out it was the magical struct.  Duh.  Of course, once I understood it I saw it in a bunch of other threads.  It was there, I just couldn't absorb it until I wore a forehead shaped hole in the wall.  Interesting though.  I've been meaning to investigate channels.

* Documentation was perhaps an hour or so sprinkled throughout the project.

# Unimplemented Requirements

I chose not to transparently proxy RESP.  Honestly, if we want to go that far, let's just use Nginx or HAProxy.  

It's already proven, sets up in minutes.  There's a whole lot of testing and support around it.  Right tool for the right job.  

I might have taken a stab at it if I hadn't gone down the channel rabbit hole, but enough is enough.

