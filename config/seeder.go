package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"my-personal-web/database"
)

// RunSeeder jalankan perintah seed
func RunSeeder(command string, name string) {
	switch command {
	case "create":
		if name == "" {
			log.Fatal("⚠️  Seeder name required. Example: go run main.go seed create users_seeder")
		}
		CreateSeeder(name)

	case "run":
		if name == "" {
			log.Fatal("⚠️  Seeder name required. Example: go run main.go seed run users_seeder")
		}
		runSeeder(name)

	case "refresh":
		refreshSeeders()

	default:
		fmt.Println("Unknown command. Use: create | run | refresh")
	}
}

// buat file seeder baru
func CreateSeeder(name string) {
	seederDir := "seeders"
	if _, err := os.Stat(seederDir); os.IsNotExist(err) {
		os.MkdirAll(seederDir, os.ModePerm)
	}

	// nama file dengan timestamp
	timestamp := time.Now().Format("20060102150405")
	filename := filepath.Join(seederDir, fmt.Sprintf("%s_%s.sql", timestamp, name))

	// isi default SQL
	content := fmt.Sprintf(`-- Seeder: %s
INSERT INTO users (name, slug, username, password, image, about)
VALUES ('Example User', 'example-user', 'example', 'hashed_password', 'default.png', 'Seeder example');
`, name)

	// tulis ke file
	err := os.WriteFile(filename, []byte(content), 0o644)
	if err != nil {
		log.Fatalf("⚠️  Error creating seeder file: %v", err)
	}

	fmt.Println("✅ Seeder created:", filename)
}

// jalankan 1 file seeder
func runSeeder(name string) {
	filePath := fmt.Sprintf("seeders/%s.sql", name)
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("⚠️  Cannot read seeder file: %v", err)
	}

	queries := strings.Split(string(sqlBytes), ";")
	for _, query := range queries {
		q := strings.TrimSpace(query)
		if q == "" {
			continue
		}
		res := database.DB.Exec(q)
		if res.Error != nil {
			log.Fatalf("⚠️ Failed execute query: %v", res.Error)
		}
		log.Printf("Rows affected: %d", res.RowsAffected)

	}

	fmt.Printf("✅ Seeder %s executed successfully\n", name)
}

// jalankan semua seeder
func refreshSeeders() {
	files, err := os.ReadDir("seeders")
	if err != nil {
		log.Fatalf("⚠️  Cannot read seeders directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			name := strings.TrimSuffix(file.Name(), ".sql")
			runSeeder(name)
		}
	}

	fmt.Println("✅ All seeders refreshed successfully")
}
