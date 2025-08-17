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
