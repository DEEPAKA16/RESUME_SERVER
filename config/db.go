package config

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	// Optional: Load .env only in local dev
	if os.Getenv("RENDER") != "true" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("No .env file found, using environment variables")
		}
	}

	// Load PEM from environment variable content
	pemContent := os.Getenv("DB_SSL_CA_CONTENT")
	if pemContent == "" {
		log.Fatal("DB_SSL_CA_CONTENT environment variable not set")
	}
	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM([]byte(pemContent)); !ok {
		log.Fatal("Failed to append PEM content")
	}

	// Register TLS config
	err := mysql.RegisterTLSConfig("tidb", &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		log.Fatalf("Failed to register TLS config: %v", err)
	}

	// DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=tidb",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	fmt.Println("✅ Successfully connected to TiDB Cloud!")
}
