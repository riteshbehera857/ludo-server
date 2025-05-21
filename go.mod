module metagame/gameserver

go 1.23.2

require github.com/gorilla/websocket v1.5.3 // indirect

require (
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.9.0
	go.mongodb.org/mongo-driver v1.17.1
	ludo v0.0.0
	messaging v0.0.0

)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lobby v0.0.0-00010101000000-000000000000 // indirect
	rng v0.0.0 // indirect
)

replace ludo => ./ludo

replace lobby => ./lobby

replace messaging => ./messaging

replace rng => ./rng
