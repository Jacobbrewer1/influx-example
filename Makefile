run: setup
	@echo "Running example"
	go test ./...
	make clean
setup:
	@echo "Setting up"
	docker-compose --env-file influxdb_example.env up -d
	@echo "Done"
clean:
	@echo "Cleaning up"
	docker-compose --env-file influxdb_example.env down
	@echo "Done"
