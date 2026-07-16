package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/frhan23/dashboard-go/internal/entity"
)

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type ArticleRepository interface {
	Create(ctx context.Context, article *entity.Article) error
	GetByID(ctx context.Context, id int) (entity.Article, error)
	Fetch(ctx context.Context, limit, offset int, status string) ([]entity.Article, int, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id int) error
}

type articleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) Create(ctx context.Context, article *entity.Article) error {
	query := `INSERT INTO posts (title, content, category, status, created_date, updated_date)
	          VALUES (?, ?, ?, ?, NOW(), NOW())`
	res, err := r.db.ExecContext(ctx, query, article.Title, article.Content, article.Category, article.Status)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	article.ID = int(id)
	return nil
}

func (r *articleRepository) Fetch(ctx context.Context, limit, offset int, status string) ([]entity.Article, int, error) {
	var query string
	var countQuery string
	var args []any
	var countArgs []any

	if status != "" {
		query = `SELECT id, title, content, category, created_date, updated_date, status
			         FROM posts WHERE status = ? ORDER BY updated_date DESC LIMIT ? OFFSET ?`
		args = append(args, status, limit, offset)

		countQuery = `SELECT COUNT(1) FROM posts WHERE status = ?`
		countArgs = append(countArgs, status)
	} else {
		query = `SELECT id, title, content, category, created_date, updated_date, status
			         FROM posts WHERE status != 'trash' ORDER BY updated_date DESC LIMIT ? OFFSET ?`
		args = append(args, limit, offset)

		countQuery = `SELECT COUNT(1) FROM posts WHERE status != 'trash'`
	}

	var totalRows int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalRows)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting rows: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying articles: %w", err)
	}
	defer rows.Close()

	var articles []entity.Article
	for rows.Next() {
		var a entity.Article
		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.CreatedDate, &a.UpdatedDate, &a.Status)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %w", err)
		}
		articles = append(articles, a)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return articles, totalRows, nil
}

func (r *articleRepository) GetByID(ctx context.Context, id int) (entity.Article, error) {
	query := `SELECT id, title, content, category, created_date, updated_date, status
	          FROM posts WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var a entity.Article
	err := row.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.CreatedDate, &a.UpdatedDate, &a.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Article{}, errors.New("article not found")
		}
		return entity.Article{}, err
	}

	return a, nil
}

func (r *articleRepository) Update(ctx context.Context, article *entity.Article) error {
	query := `UPDATE posts SET title = ?, content = ?, category = ?, status = ?, updated_date = NOW()
	          WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, article.Title, article.Content, article.Category, article.Status, article.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no article updated (ID might not exist)")
	}

	return nil
}

func (r *articleRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE posts SET status = 'trash', updated_date = NOW() WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no article deleted (ID might not exist)")
	}

	return nil
}
