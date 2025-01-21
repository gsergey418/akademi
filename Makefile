GC ::= go

DOCKER_CMD ::= docker
DOCKER_NETWORK ::= akademi
DOCKER_PREFIX ::= akademi_
DOCKER_BOOTSTRAP_PREFIX ::= akademi_bootstrap_
BOOTSTRAP_NODES ::= 3
SWARM_PEERS ::= 10

.PHONY: docker, docker_clean, swarm, swarm_stop, clean, cleanall, test

akademi: pb
	cd src/cmd && ${GC} build -o ../../akademi .

pb:
	protoc --go_out=src/ src/pb/message.proto

test:
	find . -type d -name tests -exec sh -c "cd {}; go test -v -count=1 ." \;

docker: akademi
	${DOCKER_CMD} build -t akademi:latest .

docker_clean:
	${DOCKER_CMD} rmi akademi || exit 0

swarm: docker
	${DOCKER_CMD} ps | awk '{ print $$1,$$3 }' | grep akademi | awk '{print $$1 }' | xargs -I {} ${DOCKER_CMD} stop {}
	${DOCKER_CMD} ps -a | awk '{ print $$1,$$3 }' | grep akademi | awk '{print $$1 }' | xargs -I {} ${DOCKER_CMD} rm {}
	${DOCKER_CMD} network ls | grep ${DOCKER_NETWORK} || ${DOCKER_CMD} network create ${DOCKER_NETWORK}

	for i in $$(seq ${BOOTSTRAP_NODES}); do\
		${DOCKER_CMD} run -d --network=${DOCKER_NETWORK} --name ${DOCKER_BOOTSTRAP_PREFIX}$$i akademi /bin/akademi --no-bootstrap;\
	done
	for i in $$(seq ${SWARM_PEERS}); do\
		${DOCKER_CMD} run -d --network=${DOCKER_NETWORK} --name ${DOCKER_PREFIX}$$i akademi;\
	done

swarm_stop:
	${DOCKER_CMD} ps | awk '{ print $$1,$$3 }' | grep akademi | awk '{print $$1 }' | xargs -I {} ${DOCKER_CMD} stop {}
	${DOCKER_CMD} ps -a | awk '{ print $$1,$$3 }' | grep akademi | awk '{print $$1 }' | xargs -I {} ${DOCKER_CMD} rm {}
	${DOCKER_CMD} network ls | grep ${DOCKER_NETWORK} && ${DOCKER_CMD} network rm ${DOCKER_NETWORK} || exit 0

cleanall: swarm_stop docker_clean clean

clean:
	rm -f akademi