default:
	go run cmd/main.go

migrate:
	go run cmd/migrate/main.go


install-vegeta:
	go install github.com/tsenart/vegeta/v12@latest

test-vegeta:
	echo "POST http://localhost:8081/messages/" | vegeta attack -duration=15s -rate=1000/s -header 'Content-Type: application/json' -body <(cat req.json)
