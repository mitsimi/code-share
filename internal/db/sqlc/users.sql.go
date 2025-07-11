// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id,
    username,
    email,
    password_hash
) VALUES (
    ?, ?, ?, ?
) 
RETURNING id, username, avatar, email, password_hash, created_at, updated_at
`

type CreateUserParams struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.queryRow(ctx, q.createUserStmt, createUser,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.PasswordHash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Avatar,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, avatar, email, password_hash, created_at, updated_at FROM users
WHERE id = ?
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.queryRow(ctx, q.getUserStmt, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Avatar,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, username, avatar, email, password_hash, created_at, updated_at FROM users
WHERE email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Avatar,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, avatar, email, password_hash, created_at, updated_at FROM users
WHERE username = ?
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.queryRow(ctx, q.getUserByUsernameStmt, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Avatar,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserAvatar = `-- name: UpdateUserAvatar :one
UPDATE users
SET 
    avatar = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, username, avatar, email, password_hash, created_at, updated_at
`

type UpdateUserAvatarParams struct {
	Avatar sql.NullString `json:"avatar"`
	ID     string         `json:"id"`
}

func (q *Queries) UpdateUserAvatar(ctx context.Context, arg UpdateUserAvatarParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserAvatarStmt, updateUserAvatar, arg.Avatar, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Avatar,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserInfo = `-- name: UpdateUserInfo :one
UPDATE users
SET 
    username = ?,
    email = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, username, avatar, email, password_hash, created_at, updated_at
`

type UpdateUserInfoParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserInfoStmt, updateUserInfo, arg.Username, arg.Email, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Avatar,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET 
    password_hash = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
`

type UpdateUserPasswordParams struct {
	PasswordHash string `json:"password_hash"`
	ID           string `json:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.exec(ctx, q.updateUserPasswordStmt, updateUserPassword, arg.PasswordHash, arg.ID)
	return err
}
