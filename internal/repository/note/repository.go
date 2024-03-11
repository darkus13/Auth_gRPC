package note

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/darkus13/Auth_gRPC/pkg/user_api_v1"
)

const (
	dbDSN    = "host=localhost port=54321 dbname=auth user=darkus password=andrej sslmode=disable"
	grpcPort = 50051
)

type User struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	PassConfirm string    `db:"password_confirm"`
	Role        string    `db:"role"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Структура сервера
type server struct {
	user_api_v1.UnimplementedUserV1Server
	db *pgxpool.Pool
	qb sq.StatementBuilderType
}

func (s *server) Get(ctx context.Context, req *user_api_v1.GetRequest) (*user_api_v1.GetResponse, error) {
	ID := req.GetId()

	SelectBuilder := s.qb.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": ID})

	query, args, err := SelectBuilder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	row, err := s.db.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to get user from query: %v", err)
	}

	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByNameLax[User])
	if err != nil {
		log.Fatalf("failed to collect user from db: %v", err)
	}

	roleNum := user_api_v1.Role_value[user.Role]

	userData := user_api_v1.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user_api_v1.Role(roleNum),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}

	return &userData, nil
}

func (s *server) Create(ctx context.Context, req *user_api_v1.CreateRequest) (*user_api_v1.CreateResponse, error) {
	Name := req.GetName()
	Email := req.GetEmail()
	Password := req.GetPassword()
	ConfirmPassword := req.GetPasswordConfirm()
	Role := req.GetRole().String()

	InsertBuilder := s.qb.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "password_confirm", "role").
		Values(Name, Email, Password, ConfirmPassword, Role).
		Suffix("RETURNING id")

	query, args, err := InsertBuilder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int64

	err = s.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with ID: %d", userID)

	return &user_api_v1.CreateResponse{
		Id: userID,
	}, nil
}

func (s *server) Update(ctx context.Context, req *user_api_v1.UpdateRequest) (*emptypb.Empty, error) {

	ID := req.GetId()
	Name := req.GetName()
	Email := req.GetEmail()
	Role := req.GetRole()

	UpdateBuilder := s.qb.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("name", Name).
		Set("email", Email).
		Set("role", Role).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": ID})

	query, agrs, err := UpdateBuilder.ToSql()
	if err != nil {
		log.Fatalf("failed to builder update: %v", err)
	}

	res, err := s.db.Exec(ctx, query, agrs...)
	if err != nil {
		log.Fatalf("failed to update db row: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return &emptypb.Empty{}, nil
}

// Delete Cоздадим удаление пользователя.
func (s *server) Delete(ctx context.Context, req *user_api_v1.DeleteRequest) (*emptypb.Empty, error) {

	ID := req.GetId()

	DeleteBuilder := s.qb.Delete("users").
		Where(sq.Eq{"id": ID})

	query, args, err := DeleteBuilder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	row, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	log.Printf("delete %d rows", row.RowsAffected())

	return &emptypb.Empty{}, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pgxConfig, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		log.Fatalf("failed to patde config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("ping to postgres failed: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	user_api_v1.RegisterUserV1Server(s, &server{
		db: pool,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
