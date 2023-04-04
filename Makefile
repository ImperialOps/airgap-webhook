.SILENT: run
run:
	go run .

.SILENT: docker_build
docker_build:
	docker build -t airgap-webhook .

.SILENT: docker_run
docker_run: docker_build
	docker run -p 8000:8000 airgap-webhook
