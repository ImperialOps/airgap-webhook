.SILENT: build
build:
	docker build -t airgap-webhook .

.SILENT: run
run: build
	docker run -p 8000:8000 airgap-webhook
