docker:
	docker build -t kanshiroron/test-webserver .

run:
	go run .


.PHONY: docker run