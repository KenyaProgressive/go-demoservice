run: shutdown
	docker compose up -d
	docker exec broker /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic demoservice-orders
	go run .

shutdown:
	docker compose down
