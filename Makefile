GC ::= go

DOCKER_CMD ::= docker
DOCKER_NETWORK ::= akademi
DOCKER_PREFIX ::= akademi_
SWARM_PEERS ::= 10

akademi:
	cd src/cmd && ${GC} build -o ../../akademi .

.PHONY: docker_image
docker_image: akademi
	${DOCKER_CMD} build -t akademi .

.PHONY: docker_clean
docker_clean:
	${DOCKER_CMD} rmi akademi

.PHONY: swarm_start
swarm_start: docker_image
	${DOCKER_CMD} ps -a | awk '{ print $$1,$$2 }' | grep akademi | awk '{print $$1 }' | xargs -I {} ${DOCKER_CMD} stop {}
	${DOCKER_CMD} network ls | grep ${DOCKER_NETWORK} || ${DOCKER_CMD} network create ${DOCKER_NETWORK}

	${DOCKER_CMD} run --rm -d --network=${DOCKER_NETWORK} --name ${DOCKER_PREFIX}bootstrap -p 3856:3856 akademi
	for i in $$(seq ${SWARM_PEERS}); do\
		${DOCKER_CMD} run --rm -d --network=${DOCKER_NETWORK} --name ${DOCKER_PREFIX}$$i akademi;\
	done

.PHONY: swarm_stop
swarm_stop:
	${DOCKER_CMD} ps -a | awk '{ print $$1,$$2 }' | grep akademi | awk '{print $$1 }' | xargs -I {} ${DOCKER_CMD} stop {}
	${DOCKER_CMD} network ls | grep ${DOCKER_NETWORK} && ${DOCKER_CMD} network rm ${DOCKER_NETWORK} || return 0

.PHONY: clean
clean:
	rm akademi