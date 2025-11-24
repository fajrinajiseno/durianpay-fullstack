package config

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	JwtSecret   = []byte(getEnv("JWT_SECRET", "dev-secret-replace-me"))
	HttpAddress = getEnv("HTTP_ADDR", ":8080")
)

type contextUserId string

const (
	ContextUserID contextUserId = "user_id"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
