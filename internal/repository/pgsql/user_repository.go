package pgsql

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/model"
	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/table"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Users, error) {
	t := table.Users

	stmt := t.
		SELECT(t.AllColumns).
		WHERE(t.ID.EQ(postgres.UUID(id))).
		LIMIT(1)

	var user model.Users
	err := stmt.QueryContext(ctx, r.db, &user)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
