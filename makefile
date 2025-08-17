run-gen: shutdown
	docker compose up -d
	docker exec broker /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic demoservice-orders
	exec go run . -gen=true


run: shutdown
	docker compose up -d
	docker exec broker /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic demoservice-orders
	exec go run .

shutdown:
	docker compose down

docs-generate:
	swag init -g /web/backend/app.go

lint:
	golangci-lint run .

recreate_topic:
	docker compose up -d
	docker exec broker /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic demoservice-orders
	docker exec broker /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic demoservice-orders
