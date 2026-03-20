package repository

import (
	"database/sql"

	"video-platform/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(user model.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (user_id, public_key, created_at)
		 VALUES (?, ?, ?)
		 ON CONFLICT(user_id) DO UPDATE SET public_key=excluded.public_key`,
		user.ID, user.PublicKey, user.CreatedAt,
	)
	return err
}

func (r *UserRepository) FindByID(id string) (model.User, bool) {
	var u model.User
	err := r.db.QueryRow(
		`SELECT user_id, public_key, created_at FROM users WHERE user_id = ?`,
		id,
	).Scan(&u.ID, &u.PublicKey, &u.CreatedAt)
	if err != nil {
		return model.User{}, false
	}
	return u, true
}

func (r *UserRepository) FindByPublicKey(publicKey string) (model.User, bool) {
	var u model.User
	err := r.db.QueryRow(
		`SELECT user_id, public_key, created_at FROM users WHERE public_key = ?`,
		publicKey,
	).Scan(&u.ID, &u.PublicKey, &u.CreatedAt)
	if err != nil {
		return model.User{}, false
	}
	return u, true
}
