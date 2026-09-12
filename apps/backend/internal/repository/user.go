package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xanity-07/openmat/internal/enums"
	"github.com/xanity-07/openmat/internal/lib/utils"
	"github.com/xanity-07/openmat/internal/model/user"
	"github.com/xanity-07/openmat/internal/server"
)

type UserRepository struct {
	server *server.Server
}

func NewUserRepository(server *server.Server) *UserRepository {
	return &UserRepository{server: server}
}

func (r *UserRepository) CreateUser(ctx context.Context, payload *user.CreateUserPayload) (*user.User, error) {
	stmt := `
		INSERT INTO
		    users (
		           id,
		           email,
		           password_hash,
		           role,
		           created_at,
		           updated_at
		    )
			VALUES (
			        @id,
			        @email,
			        @password_hash,
			        @role,
			        @created_at,
			        @updated_at
			)
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":            utils.GenerateID(11),
		"email":         payload.Email,
		"password_hash": payload.Password,
		"role":          enums.USER,
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	createdUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table users: %w", err)
	}

	return &createdUser, nil
}

func (r *UserRepository) CheckUserExistsEmail(ctx context.Context, email string) (bool, error) {
	stmt := `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE email = @email
		) AS exists
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"email": email,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	exists, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	stmt := `
		SELECT
		    id,
		    email,
		    password_hash,
		    role,
		    created_at,
		    updated_at
		FROM
			users
		WHERE
			email = @email
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"email": email,
	})
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	foundUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table users: %w", err)
	}

	return &foundUser, nil
}

func (r *UserRepository) GetUsers(ctx context.Context, q *user.GetUsersQuery) ([]user.User, error) {
	stmt := `
		SELECT
			id,
			email,
			role,
			password_hash,
			created_at,
			updated_at
		FROM
			users
	`

	conditions := []string{}
	args := pgx.NamedArgs{}

	if q.Search != nil {
		conditions = append(conditions, "role ILIKE @search OR email ILIKE @search")
		args["search"] = "%" + *q.Search + "%"
	}

	if q.Role != nil {
		conditions = append(conditions, "role = @role")
		args["role"] = *q.Role
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause += " WHERE " + strings.Join(conditions, " AND ")
	}

	stmt += whereClause

	var count int
	countStmt := "SELECT COUNT(*) FROM users " + whereClause
	err := r.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to get count: %w", err)
	}

	if q.Order != nil {
		stmt += "ORDER BY " + *q.Order
		if q.Sort != nil && *q.Sort != "desc" {
			stmt += " DESC"
		} else {
			stmt += " ASC"
		}
	} else {
		stmt += "ORDER BY created_at DESC"
	}

	stmt += " OFFSET @offset LIMIT @limit"
	args["limit"] = *q.Limit
	args["offset"] = (*q.Page - 1) * *q.Limit

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute fetch users query: %w", err)
	}

	userList, err := pgx.CollectRows(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table user: %w", err)
	}

	return userList, nil
}
