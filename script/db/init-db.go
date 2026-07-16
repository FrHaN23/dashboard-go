package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/frhan23/dashboard-go/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to the database.")

	log.Println("Initializing database tables...")
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS posts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		content TEXT NOT NULL,
		category VARCHAR(100) NOT NULL,
		created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		status VARCHAR(100) NOT NULL
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create 'posts' table: %v", err)
	}
	log.Println("Table 'posts' is ready ✔")

	var count int
	err = db.QueryRow("SELECT COUNT(1) FROM posts").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check existing data: %v", err)
	}

	if count == 0 {
		log.Println("Seeding initial data...")
		seedQuery := `
		INSERT INTO posts (title, content, category, status) VALUES
		(
			'Belajar Dasar Golang REST API untuk Pemula',
			'Lorem ipsum dolor sit amet consectetur adipiscing elit quisque faucibus ex sapien vitae pellentesque sem placerat in id cursus mi pretium tellus duis convallis tempus leo eu aenean sed diam urna tempor pulvinar vivamus fringilla lacus nec metus bibendum egestas iaculis massa nisl malesuada lacinia integer nunc posuere ut hendrerit semper vel class aptent taciti sociosqu ad litora torquent per conubia nostra inceptos himenaeos orci varius natoque penatibus et magnis dis parturient montes nascetur ridiculus mus donec rhoncus eros lobortis nulla molestie mattis scelerisque maximus eget fermentum odio phasellus non purus est efficitur laoreet mauris pharetra vestibulum fusce dictum risus.',
			'Golang',
			'publish'
		),
		(
			'Tips Menulis Clean Code di Bahasa Pemrograman Go',
			'Lorem ipsum dolor sit amet consectetur adipiscing elit quisque faucibus ex sapien vitae pellentesque sem placerat in id cursus mi pretium tellus duis convallis tempus leo eu aenean sed diam urna tempor pulvinar vivamus fringilla lacus nec metus bibendum egestas iaculis massa nisl malesuada lacinia integer nunc posuere ut hendrerit semper vel class aptent taciti sociosqu ad litora torquent per conubia nostra inceptos himenaeos orci varius natoque penatibus et magnis dis parturient montes nascetur ridiculus mus donec rhoncus eros lobortis nulla molestie mattis scelerisque maximus eget fermentum odio phasellus non purus est efficitur laoreet mauris pharetra vestibulum fusce dictum risus.',
			'Software Engineering',
			'draft'
		),
		(
			'Artikel Ini Sudah Tidak Relevan Lagi',
			'Lorem ipsum dolor sit amet consectetur adipiscing elit quisque faucibus ex sapien vitae pellentesque sem placerat in id cursus mi pretium tellus duis convallis tempus leo eu aenean sed diam urna tempor pulvinar vivamus fringilla lacus nec metus bibendum egestas iaculis massa nisl malesuada lacinia integer nunc posuere ut hendrerit semper vel class aptent taciti sociosqu ad litora torquent per conubia nostra inceptos himenaeos orci varius natoque penatibus et magnis dis parturient montes nascetur ridiculus mus donec rhoncus eros lobortis nulla molestie mattis scelerisque maximus eget fermentum odio phasellus non purus est efficitur laoreet mauris pharetra vestibulum fusce dictum risus.',
			'Lain-lain',
			'trash'
		);`

		_, err = db.Exec(seedQuery)
		if err != nil {
			log.Fatalf("Failed to seed data: %v", err)
		}
		log.Println("Database seeding completed successfully ✔")
	} else {
		log.Println("Database already has data. Skipping seed.")
	}

	fmt.Println("\nDatabase Init & Seeding Process finished successfully!")
}
