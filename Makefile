generate:
	sqlc -f db/sqlc.yaml -x generate


test: generate
	go test ./...


run: generate
	go run .


precommit:
	pre-commit run -a


migrate:
	tern migrate -c db/migrations/tern.conf  -m db/migrations
