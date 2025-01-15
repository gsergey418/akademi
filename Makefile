akademi:
	cd src/cmd && go build -o ../../akademi .

.PHONY: docker_image
docker_image: akademi
	docker build -t akademi .

.PHONY: docker_clean
docker_clean:
	docker rmi akademi

.PHONY: swarm_start
swarm_start: docker_image
	docker ps -a | awk '{ print $$1,$$2 }' | grep akademi | awk '{print $$1 }' | xargs -I {} docker stop {}
	docker network ls | grep akademi || docker network create akademi

	docker run --rm -d --network=akademi --name akademi_bootstrap -p 3856:3856 akademi
	for i in $$(seq 10); do\
		docker run --rm -d --network=akademi --name akademi_$$i akademi;\
	done

.PHONY: swarm_stop
swarm_stop:
	docker ps -a | awk '{ print $$1,$$2 }' | grep akademi | awk '{print $$1 }' | xargs -I {} docker stop {}
	docker network ls | grep akademi && docker network rm akademi || return 0

.PHONY: clean
clean:
	rm akademi