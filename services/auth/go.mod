module github.com/didar-platform/services/auth

go 1.23.3

require (
	github.com/didar-platform/pkg/authkit v0.0.0
	github.com/didar-platform/pkg/platform v0.0.0
	github.com/go-chi/chi/v5 v5.3.1
	github.com/go-chi/cors v1.2.2
	github.com/jackc/pgx/v5 v5.7.5
	github.com/joho/godotenv v1.5.1
	github.com/kavenegar/kavenegar-go v0.0.0-20240205151018-77039f51467d
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace (
	github.com/didar-platform/pkg/authkit => ../../pkg/authkit
	github.com/didar-platform/pkg/platform => ../../pkg/platform
)
