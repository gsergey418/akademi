GC ::= go

DOCKER_CMD ::= docker
DOCKER_NETWORK ::= akademi
DOCKER_PREFIX ::= akademi_
DOCKER_BOOTSTRAP_PREFIX ::= akademi_bootstrap_
BOOTSTRAP_NODES ::= 3
SWARM_PEERS ::= 100

.PHONY: docker, clean, test

akademi: pb
	cd src/cmd && ${GC} build -o ../../akademi .

pb:
	protoc --go_out=src/ src/pb/message.proto

test: akademi
	./akademi daemon --no-bootstrap &
	find . -type d -name tests -exec sh -c "cd {}; go test -v -count=1 ." \;
	pkill akademi


docker: akademi
	${DOCKER_CMD} build -t akademi:latest .

clean:
	rm -f akademi