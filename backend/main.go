package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/fajrinajiseno/mygolangapp/internal/api"
	"github.com/fajrinajiseno/mygolangapp/internal/config"
	ah "github.com/fajrinajiseno/mygolangapp/internal/module/auth/handler"
	ar "github.com/fajrinajiseno/mygolangapp/internal/module/auth/repository"
	au "github.com/fajrinajiseno/mygolangapp/internal/module/auth/usecase"
	ph "github.com/fajrinajiseno/mygolangapp/internal/module/payment/handler"
	pr "github.com/fajrinajiseno/mygolangapp/internal/module/payment/repository"
	pu "github.com/fajrinajiseno/mygolangapp/internal/module/payment/usecase"
	srv "github.com/fajrinajiseno/mygolangapp/internal/service/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	db, err := sql.Open("sqlite3", "dashboard.db?_foreign_keys=1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}

	JwtExpiredDuration, err := time.ParseDuration(config.JwtExpired)
	if err != nil {
		panic(err)
	}

	userRepo := ar.NewUserRepo(db)
	paymentRepo := pr.NewPaymentRepo(db)

	authUC := au.NewAuthUsecase(userRepo, config.JwtSecret, JwtExpiredDuration)
	paymentUC := pu.NewPaymentUsecase(paymentRepo, userRepo)

	authH := ah.NewAuthHandler(paymentUC, authUC)
	paymentH := ph.NewPaymentHandler(paymentUC)

	apiHandler := &api.APIHandler{
		Auth:    authH,
		Payment: paymentH,
	}

	server := srv.NewServer(apiHandler, config.OpenapiYamlLocation)

	addr := config.HttpAddress
	log.Printf("starting server on %s", addr)
	server.Start(addr)
}

func initDB(db *sql.DB) error {
	// create tables if not exists
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  email TEXT NOT NULL UNIQUE,
		  password_hash TEXT NOT NULL,
		  role TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS payments (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  merchant TEXT NOT NULL,
		  status TEXT NOT NULL,
		  amount REAL NOT NULL,
		  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	// seed payments if empty
	var cnt int
	row := db.QueryRow("SELECT COUNT(1) FROM payments")
	if err := row.Scan(&cnt); err != nil {
		return err
	}
	if cnt == 0 {
		payments := []struct {
			merchant  string
			amount    float64
			status    string
			createdAt string
		}{
			{"merchant 1", 100.0, "pending", time.Now().Format(time.RFC3339)},
			{"merchant 2", 200.0, "completed", time.Now().Format(time.RFC3339)},
			{"merchant 3", 150.5, "failed", time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
			{"merchant 4", 100.0, "pending", time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
			{"merchant 5", 200.0, "completed", time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
			{"merchant 6", 150.5, "failed", time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
			{"merchant 7", 100.0, "pending", time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
			{"merchant 8", 200.0, "completed", time.Now().Add(-48 * time.Hour).Format(time.RFC3339)},
			{"merchant 9", 150.5, "failed", time.Now().Add(-48 * time.Hour).Format(time.RFC3339)},
			{"merchant 10", 100.0, "pending", time.Now().Add(-48 * time.Hour).Format(time.RFC3339)},
			{"merchant 11", 200.0, "completed", time.Now().Add(-48 * time.Hour).Format(time.RFC3339)},
			{"merchant 12", 150.5, "failed", time.Now().Add(-48 * time.Hour).Format(time.RFC3339)},
		}
		for _, p := range payments {
			if _, err := db.Exec("INSERT INTO payments(merchant, amount, status, created_at) VALUES (?, ?, ?, ?)", p.merchant, p.amount, p.status, p.createdAt); err != nil {
				return err
			}
		}
	}

	// seed admin user if not exists
	row = db.QueryRow("SELECT COUNT(1) FROM users")
	if err := row.Scan(&cnt); err != nil {
		return err
	}
	if cnt == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO users(email, password_hash, role) VALUES (?, ?, ?)", "cs@test.com", string(hash), "cs"); err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO users(email, password_hash, role) VALUES (?, ?, ?)", "operation@test.com", string(hash), "operation"); err != nil {
			return err
		}
	}

	const dbLifetime = time.Minute * 5
	db.SetConnMaxLifetime(dbLifetime)
	return nil
}
