package mysql

import (
	"context"
	"database/sql"

	"github.com/YudhistiraTA/sharing-vision-be/model"
)

type ArticleRepository struct {
	DB *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{DB: db}
}

func (r *ArticleRepository) FindByID(ctx context.Context, id int) (res model.Article, err error) {
	query := "SELECT id, title, content, category, created_at, updated_at, status FROM posts WHERE id = ?"
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return model.Article{}, err
	}
	row := stmt.QueryRowContext(ctx, id)
	res = model.Article{}
	err = row.Scan(&res.ID, &res.Title, &res.Content, &res.Category, &res.CreatedAt, &res.UpdatedAt, &res.Status)

	return
}

func (r *ArticleRepository) Delete(ctx context.Context, id int) (err error) {
	query := "UPDATE posts SET status = 'Trash' WHERE id = ?"
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, id)

	return
}

func (r *ArticleRepository) FindAll(ctx context.Context, limit int, offset int, status string) (res []model.Article, err error) {
	baseQuery := "SELECT id, title, content, category, created_at, updated_at, status FROM posts"
	var query string
	var args []interface{}

	if status != "" {
		query = baseQuery + " WHERE status LIKE ? ORDER BY updated_at DESC LIMIT ? OFFSET ?"
		args = append(args, status, limit, offset)
	} else {
		query = baseQuery + " ORDER BY updated_at DESC LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article model.Article
		err = rows.Scan(&article.ID, &article.Title, &article.Content, &article.Category, &article.CreatedAt, &article.UpdatedAt, &article.Status)
		if err != nil {
			return
		}
		res = append(res, article)
	}

	return
}

func (r *ArticleRepository) Create(ctx context.Context, article *model.Article) (err error) {
	query := "INSERT INTO posts (title, content, category, status) VALUES (?, ?, ?, ?)"
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, article.Title, article.Content, article.Category, article.Status)

	return
}

func (r *ArticleRepository) Update(ctx context.Context, id int, article *model.Article) (err error) {
	query := "UPDATE posts SET title = ?, content = ?, category = ?, status = ? WHERE id = ?"
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, article.Title, article.Content, article.Category, article.Status, id)

	return
}
