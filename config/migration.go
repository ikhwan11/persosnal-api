package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(command string, args ...string) {
	// ambil config DB
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// DSN untuk migrate
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s",
		user, pass, host, port, name)

	switch command {
	case "up":
		m, err := migrate.New("file://migrations", dsn)
		if err != nil {
			log.Fatal(err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("✅ Migration UP success")
		os.Exit(0)

	case "down":
		m, err := migrate.New("file://migrations", dsn)
		if err != nil {
			log.Fatal(err)
		}
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("✅ Migration DOWN success")
		os.Exit(0)

	case "refresh":
		m, err := migrate.New("file://migrations", dsn)
		if err != nil {
			log.Fatal(err)
		}
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("✅ Migration REFRESH success")
		os.Exit(0)

	case "create":
		if len(args) < 1 {
			fmt.Println("Usage: go run main.go create <migration_name>")
			os.Exit(1)
		}
		createMigrationFile(args[0])
		os.Exit(0)

	default:
		fmt.Println("Unknown command. Use: up | down | refresh | create <name>")
		os.Exit(1)
	}
}

func createMigrationFile(name string) {
	// bikin folder migrations kalau belum ada
	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(migrationsDir, os.ModePerm); err != nil {
			log.Fatal("❌ Gagal membuat folder migrations:", err)
		}
	}

	// bikin nama file timestamped
	timestamp := time.Now().Format("20060102150405")
	safeName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	upFile := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.up.sql", timestamp, safeName))
	downFile := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.down.sql", timestamp, safeName))

	// ambil nama tabel dari migration name
	// contoh: "create_users_table" -> "users"
	tableName := safeName
	if strings.Contains(safeName, "create_") && strings.Contains(safeName, "_table") {
		tableName = strings.TrimPrefix(safeName, "create_")
		tableName = strings.TrimSuffix(tableName, "_table")
	}

	// isi default SQL
	upContent := fmt.Sprintf(`-- Migration UP: %s
CREATE TABLE %s (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    image VARCHAR(255) NOT NULL,
    about TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);`, safeName, tableName)

	downContent := fmt.Sprintf(`-- Migration DOWN: %s
DROP TABLE IF EXISTS %s;`, safeName, tableName)

	// tulis file
	if err := os.WriteFile(upFile, []byte(upContent), 0o644); err != nil {
		log.Fatal("❌ Gagal membuat file up migration:", err)
	}
	if err := os.WriteFile(downFile, []byte(downContent), 0o644); err != nil {
		log.Fatal("❌ Gagal membuat file down migration:", err)
	}

	fmt.Println("✅ Migration created:")
	fmt.Println("   ", upFile)
	fmt.Println("   ", downFile)

	os.Exit(0)
}
