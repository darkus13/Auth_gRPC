package auth

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/darkus13/Auth_gRPC/internal/converter"
	"github.com/darkus13/Auth_gRPC/internal/repository"
	"github.com/darkus13/Auth_gRPC/internal/repository/auth/model"
	decs "github.com/darkus13/Auth_gRPC/pkg/user_api_v1"
)

const (
	dbDSN                = "host=localhost port=54321 dbname=auth user=darkus password=andrej sslmode=disable"
	grpcPort             = 50051
	usersId              = "id"
	usersName            = "name"
	usersEmail           = "email"
	usersRole            = "role"
	usersCreatedat       = "created_at"
	usersUpdatedat       = "updated_at"
	tableUsers           = "users"
	usersPassword        = "password"
	usersPasswordConfirm = "password_confirm"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.AuthInfo) (int64, error) {
	builder := sq.Insert(tableUsers).
		Columns(usersName, usersEmail, usersPassword, usersPasswordConfirm, usersRole).
		Values(info.Name, info.Email, info.Password, info.PassConfirm, info.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*decs.User, error) {
	builderGet := sq.Select(usersId, usersName, usersEmail, usersRole, usersCreatedat, usersUpdatedat).
		From(tableUsers).
		Where(sq.Eq{usersId: id}).
		Limit(1)

	query, args, err := builderGet.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	var auth model.Auth
	err = r.db.QueryRow(ctx, query, args...).Scan(&auth.ID, &auth.CreatedAt, &auth.UpdatedAt, &auth.Info)
	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	return converter.ToUserFromService(&auth), nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableUsers).
		Where(sq.Eq{usersId: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to execute delete query: %v", err)
		return err
	}
	return nil
}

func (r *repo) Update(ctx context.Context, id int64, info *model.AuthInfo) error {
	builderUpdate := sq.Update(tableUsers).
		PlaceholderFormat(sq.Dollar).
		Set(usersName, info.Name).
		Set(usersEmail, info.Email).
		Set(usersRole, info.Role).
		Set(usersUpdatedat, time.Now()).
		Where(sq.Eq{usersId: id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to execute update query: %v", err)
		return err
	}
	return nil
}
