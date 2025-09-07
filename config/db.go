package config

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Register SSL CA for TiDB
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(os.Getenv("DB_SSL_CA"))
	//db file
	if err != nil {
		log.Fatalf("Failed to read SSL CA file: %v", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append PEM.")
	}

	// Register TLS config
	err = mysql.RegisterTLSConfig("tidb", &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		log.Fatalf("Failed to register TLS config: %v", err)
	}

	// DSN with tls=tidb
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

	// Check connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	fmt.Println("✅ Successfully connected to TiDB Cloud!")
}
