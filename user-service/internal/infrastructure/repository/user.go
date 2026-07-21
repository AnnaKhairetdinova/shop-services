package repository

import (
	"context"

	"github.com/AnnaKhairetdinova/user-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	query := `INSERT INTO users (uuid, name, email, created_at) VALUES ($1, $2, $3, $4) RETURNING uuid, name, email, created_at`

	var createdUser domain.User
	err := r.db.QueryRow(ctx, query, user.UUID, user.Name, user.Email, user.CreatedAt).Scan(&createdUser.UUID, &createdUser.Name, &createdUser.Email, &createdUser.CreatedAt)

	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (r *postgresRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	query := `SELECT uuid, name, email, created_at FROM users WHERE uuid = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, uuid).Scan(&user.UUID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *postgresRepository) List(ctx context.Context) ([]domain.User, error) {
	query := `SELECT uuid, name, email, created_at FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		err := rows.Scan(&u.UUID, &u.Name, &u.Email, &u.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
