package pgsql

import (
	"context"
	"database/sql"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
	"github.com/latif-ecommerce-microservices/user-service/pkg/customerror"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/model"
	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/table"
)

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) UpdateUser(ctx context.Context, user model.Users) (model.Users, error) {
	t := table.Users

	stmt := t.
		UPDATE(
			t.Name,
			t.Email,
			t.PhoneNumber,
			t.UpdatedAt,
			t.DeletedAt,
		).
		MODEL(user).
		WHERE(t.ID.EQ(postgres.UUID(user.ID))).
		RETURNING(t.AllColumns)

	err := stmt.QueryContext(ctx, r.db, &user)
	if err != nil {
		return model.Users{},
			customerror.InternalServerError.
				WithCause(err).
				WithStackTrace().
				WithLocator(customerror.WhereAmI())
	}

	return user, nil
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

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.Users, error) {
	t := table.Users
	stmt := t.
		SELECT(t.AllColumns).
		WHERE(t.Email.EQ(postgres.String(email))).
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

func (r *userRepository) CreateUser(ctx context.Context, user model.Users) (model.Users, error) {
	t := table.Users

	stmt := t.INSERT(t.AllColumns).MODEL(user).RETURNING(t.AllColumns)
	err := stmt.QueryContext(ctx, r.db, &user)
	if err != nil {
		return model.Users{}, customerror.InternalServerError.WithCause(err).WithStackTrace().WithLocator(customerror.WhereAmI())
	}

	return user, nil
}

func NewUserRepository(db *sql.DB) repository.UserRepositoryProvider {
	return &userRepository{db: db}
}
