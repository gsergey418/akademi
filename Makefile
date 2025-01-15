akademi:
	cd src/cmd && go build -o ../../akademi .

docker: akademi
	docker build -t akademi .

clean:
	rm akademi && docker rmi akademi