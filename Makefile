run-docker: 
	sudo docker run --network host test-api

build-docker:
	docker build .