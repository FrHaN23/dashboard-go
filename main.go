package main

import (
	"database/sql"
	"fmt"
	"log"

	ah "github.com/frhan23/dashboard-go/internal/module/article/handler"
	ar "github.com/frhan23/dashboard-go/internal/module/article/repository"
	au "github.com/frhan23/dashboard-go/internal/module/article/usecase"

	"github.com/frhan23/dashboard-go/internal/api"
	"github.com/frhan23/dashboard-go/internal/config"
	srv "github.com/frhan23/dashboard-go/internal/service/http"
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
		log.Fatal(err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}

	articleRepo := ar.NewArticleRepository(db)
	articleUC := au.NewArticleUsecase(articleRepo)
	articleH := ah.NewArticleHandler(articleUC)

	apiHandler := &api.APIHandler{
		Article: articleH,
	}

	server := srv.NewServer(apiHandler)

	addr := cfg.HTTPAddr
	log.Printf("starting server on %s", addr)
	server.Start(addr)
}

func initDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		content TEXT NOT NULL,
		category VARCHAR(100) NOT NULL,
		created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		status VARCHAR(100) NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}
