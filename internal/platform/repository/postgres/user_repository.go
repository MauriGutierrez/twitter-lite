package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"ualaTwitter/internal/domain/user"
)

type UserRepository struct {
	conn *pgx.Conn
}

func NewPostgresUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (r *UserRepository) Create(ctx context.Context, u user.User) error {
	_, err := r.conn.Exec(ctx, `INSERT INTO users (id, name, document) VALUES ($1, $2, $3)`,
		u.ID, u.Name, u.Document)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (user.User, error) {
	var u user.User
	err := r.conn.QueryRow(ctx, `SELECT id, name FROM users WHERE id = $1`, id).Scan(&u.ID, &u.Name)
	return u, err
}
