// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.checkLikeExistsStmt, err = db.PrepareContext(ctx, checkLikeExists); err != nil {
		return nil, fmt.Errorf("error preparing query CheckLikeExists: %w", err)
	}
	if q.checkRecentViewStmt, err = db.PrepareContext(ctx, checkRecentView); err != nil {
		return nil, fmt.Errorf("error preparing query CheckRecentView: %w", err)
	}
	if q.cleanupOldViewsStmt, err = db.PrepareContext(ctx, cleanupOldViews); err != nil {
		return nil, fmt.Errorf("error preparing query CleanupOldViews: %w", err)
	}
	if q.createSessionStmt, err = db.PrepareContext(ctx, createSession); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSession: %w", err)
	}
	if q.createSnippetStmt, err = db.PrepareContext(ctx, createSnippet); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSnippet: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.decrementLikesCountStmt, err = db.PrepareContext(ctx, decrementLikesCount); err != nil {
		return nil, fmt.Errorf("error preparing query DecrementLikesCount: %w", err)
	}
	if q.deleteExpiredSessionsStmt, err = db.PrepareContext(ctx, deleteExpiredSessions); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteExpiredSessions: %w", err)
	}
	if q.deleteLikeStmt, err = db.PrepareContext(ctx, deleteLike); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteLike: %w", err)
	}
	if q.deleteSavedSnippetStmt, err = db.PrepareContext(ctx, deleteSavedSnippet); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSavedSnippet: %w", err)
	}
	if q.deleteSessionStmt, err = db.PrepareContext(ctx, deleteSession); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSession: %w", err)
	}
	if q.deleteSnippetStmt, err = db.PrepareContext(ctx, deleteSnippet); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSnippet: %w", err)
	}
	if q.getLikedSnippetsStmt, err = db.PrepareContext(ctx, getLikedSnippets); err != nil {
		return nil, fmt.Errorf("error preparing query GetLikedSnippets: %w", err)
	}
	if q.getSavedSnippetsStmt, err = db.PrepareContext(ctx, getSavedSnippets); err != nil {
		return nil, fmt.Errorf("error preparing query GetSavedSnippets: %w", err)
	}
	if q.getSessionStmt, err = db.PrepareContext(ctx, getSession); err != nil {
		return nil, fmt.Errorf("error preparing query GetSession: %w", err)
	}
	if q.getSnippetStmt, err = db.PrepareContext(ctx, getSnippet); err != nil {
		return nil, fmt.Errorf("error preparing query GetSnippet: %w", err)
	}
	if q.getSnippetsStmt, err = db.PrepareContext(ctx, getSnippets); err != nil {
		return nil, fmt.Errorf("error preparing query GetSnippets: %w", err)
	}
	if q.getSnippetsByAuthorStmt, err = db.PrepareContext(ctx, getSnippetsByAuthor); err != nil {
		return nil, fmt.Errorf("error preparing query GetSnippetsByAuthor: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.getUserByEmailStmt, err = db.PrepareContext(ctx, getUserByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByEmail: %w", err)
	}
	if q.getUserByUsernameStmt, err = db.PrepareContext(ctx, getUserByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByUsername: %w", err)
	}
	if q.incrementLikesCountStmt, err = db.PrepareContext(ctx, incrementLikesCount); err != nil {
		return nil, fmt.Errorf("error preparing query IncrementLikesCount: %w", err)
	}
	if q.incrementViewsStmt, err = db.PrepareContext(ctx, incrementViews); err != nil {
		return nil, fmt.Errorf("error preparing query IncrementViews: %w", err)
	}
	if q.likeSnippetStmt, err = db.PrepareContext(ctx, likeSnippet); err != nil {
		return nil, fmt.Errorf("error preparing query LikeSnippet: %w", err)
	}
	if q.recordViewStmt, err = db.PrepareContext(ctx, recordView); err != nil {
		return nil, fmt.Errorf("error preparing query RecordView: %w", err)
	}
	if q.saveSnippetStmt, err = db.PrepareContext(ctx, saveSnippet); err != nil {
		return nil, fmt.Errorf("error preparing query SaveSnippet: %w", err)
	}
	if q.updateLikesCountStmt, err = db.PrepareContext(ctx, updateLikesCount); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateLikesCount: %w", err)
	}
	if q.updateSessionExpiryStmt, err = db.PrepareContext(ctx, updateSessionExpiry); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSessionExpiry: %w", err)
	}
	if q.updateSnippetStmt, err = db.PrepareContext(ctx, updateSnippet); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSnippet: %w", err)
	}
	if q.updateUserAvatarStmt, err = db.PrepareContext(ctx, updateUserAvatar); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserAvatar: %w", err)
	}
	if q.updateUserInfoStmt, err = db.PrepareContext(ctx, updateUserInfo); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserInfo: %w", err)
	}
	if q.updateUserPasswordStmt, err = db.PrepareContext(ctx, updateUserPassword); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserPassword: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.checkLikeExistsStmt != nil {
		if cerr := q.checkLikeExistsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing checkLikeExistsStmt: %w", cerr)
		}
	}
	if q.checkRecentViewStmt != nil {
		if cerr := q.checkRecentViewStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing checkRecentViewStmt: %w", cerr)
		}
	}
	if q.cleanupOldViewsStmt != nil {
		if cerr := q.cleanupOldViewsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing cleanupOldViewsStmt: %w", cerr)
		}
	}
	if q.createSessionStmt != nil {
		if cerr := q.createSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSessionStmt: %w", cerr)
		}
	}
	if q.createSnippetStmt != nil {
		if cerr := q.createSnippetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSnippetStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.decrementLikesCountStmt != nil {
		if cerr := q.decrementLikesCountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing decrementLikesCountStmt: %w", cerr)
		}
	}
	if q.deleteExpiredSessionsStmt != nil {
		if cerr := q.deleteExpiredSessionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteExpiredSessionsStmt: %w", cerr)
		}
	}
	if q.deleteLikeStmt != nil {
		if cerr := q.deleteLikeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteLikeStmt: %w", cerr)
		}
	}
	if q.deleteSavedSnippetStmt != nil {
		if cerr := q.deleteSavedSnippetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSavedSnippetStmt: %w", cerr)
		}
	}
	if q.deleteSessionStmt != nil {
		if cerr := q.deleteSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSessionStmt: %w", cerr)
		}
	}
	if q.deleteSnippetStmt != nil {
		if cerr := q.deleteSnippetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSnippetStmt: %w", cerr)
		}
	}
	if q.getLikedSnippetsStmt != nil {
		if cerr := q.getLikedSnippetsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLikedSnippetsStmt: %w", cerr)
		}
	}
	if q.getSavedSnippetsStmt != nil {
		if cerr := q.getSavedSnippetsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSavedSnippetsStmt: %w", cerr)
		}
	}
	if q.getSessionStmt != nil {
		if cerr := q.getSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSessionStmt: %w", cerr)
		}
	}
	if q.getSnippetStmt != nil {
		if cerr := q.getSnippetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSnippetStmt: %w", cerr)
		}
	}
	if q.getSnippetsStmt != nil {
		if cerr := q.getSnippetsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSnippetsStmt: %w", cerr)
		}
	}
	if q.getSnippetsByAuthorStmt != nil {
		if cerr := q.getSnippetsByAuthorStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSnippetsByAuthorStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.getUserByEmailStmt != nil {
		if cerr := q.getUserByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByEmailStmt: %w", cerr)
		}
	}
	if q.getUserByUsernameStmt != nil {
		if cerr := q.getUserByUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByUsernameStmt: %w", cerr)
		}
	}
	if q.incrementLikesCountStmt != nil {
		if cerr := q.incrementLikesCountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing incrementLikesCountStmt: %w", cerr)
		}
	}
	if q.incrementViewsStmt != nil {
		if cerr := q.incrementViewsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing incrementViewsStmt: %w", cerr)
		}
	}
	if q.likeSnippetStmt != nil {
		if cerr := q.likeSnippetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing likeSnippetStmt: %w", cerr)
		}
	}
	if q.recordViewStmt != nil {
		if cerr := q.recordViewStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing recordViewStmt: %w", cerr)
		}
	}
	if q.saveSnippetStmt != nil {
		if cerr := q.saveSnippetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing saveSnippetStmt: %w", cerr)
		}
	}
	if q.updateLikesCountStmt != nil {
		if cerr := q.updateLikesCountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateLikesCountStmt: %w", cerr)
		}
	}
	if q.updateSessionExpiryStmt != nil {
		if cerr := q.updateSessionExpiryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateSessionExpiryStmt: %w", cerr)
		}
	}
	if q.updateSnippetStmt != nil {
		if cerr := q.updateSnippetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateSnippetStmt: %w", cerr)
		}
	}
	if q.updateUserAvatarStmt != nil {
		if cerr := q.updateUserAvatarStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserAvatarStmt: %w", cerr)
		}
	}
	if q.updateUserInfoStmt != nil {
		if cerr := q.updateUserInfoStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserInfoStmt: %w", cerr)
		}
	}
	if q.updateUserPasswordStmt != nil {
		if cerr := q.updateUserPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserPasswordStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                        DBTX
	tx                        *sql.Tx
	checkLikeExistsStmt       *sql.Stmt
	checkRecentViewStmt       *sql.Stmt
	cleanupOldViewsStmt       *sql.Stmt
	createSessionStmt         *sql.Stmt
	createSnippetStmt         *sql.Stmt
	createUserStmt            *sql.Stmt
	decrementLikesCountStmt   *sql.Stmt
	deleteExpiredSessionsStmt *sql.Stmt
	deleteLikeStmt            *sql.Stmt
	deleteSavedSnippetStmt    *sql.Stmt
	deleteSessionStmt         *sql.Stmt
	deleteSnippetStmt         *sql.Stmt
	getLikedSnippetsStmt      *sql.Stmt
	getSavedSnippetsStmt      *sql.Stmt
	getSessionStmt            *sql.Stmt
	getSnippetStmt            *sql.Stmt
	getSnippetsStmt           *sql.Stmt
	getSnippetsByAuthorStmt   *sql.Stmt
	getUserStmt               *sql.Stmt
	getUserByEmailStmt        *sql.Stmt
	getUserByUsernameStmt     *sql.Stmt
	incrementLikesCountStmt   *sql.Stmt
	incrementViewsStmt        *sql.Stmt
	likeSnippetStmt           *sql.Stmt
	recordViewStmt            *sql.Stmt
	saveSnippetStmt           *sql.Stmt
	updateLikesCountStmt      *sql.Stmt
	updateSessionExpiryStmt   *sql.Stmt
	updateSnippetStmt         *sql.Stmt
	updateUserAvatarStmt      *sql.Stmt
	updateUserInfoStmt        *sql.Stmt
	updateUserPasswordStmt    *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                        tx,
		tx:                        tx,
		checkLikeExistsStmt:       q.checkLikeExistsStmt,
		checkRecentViewStmt:       q.checkRecentViewStmt,
		cleanupOldViewsStmt:       q.cleanupOldViewsStmt,
		createSessionStmt:         q.createSessionStmt,
		createSnippetStmt:         q.createSnippetStmt,
		createUserStmt:            q.createUserStmt,
		decrementLikesCountStmt:   q.decrementLikesCountStmt,
		deleteExpiredSessionsStmt: q.deleteExpiredSessionsStmt,
		deleteLikeStmt:            q.deleteLikeStmt,
		deleteSavedSnippetStmt:    q.deleteSavedSnippetStmt,
		deleteSessionStmt:         q.deleteSessionStmt,
		deleteSnippetStmt:         q.deleteSnippetStmt,
		getLikedSnippetsStmt:      q.getLikedSnippetsStmt,
		getSavedSnippetsStmt:      q.getSavedSnippetsStmt,
		getSessionStmt:            q.getSessionStmt,
		getSnippetStmt:            q.getSnippetStmt,
		getSnippetsStmt:           q.getSnippetsStmt,
		getSnippetsByAuthorStmt:   q.getSnippetsByAuthorStmt,
		getUserStmt:               q.getUserStmt,
		getUserByEmailStmt:        q.getUserByEmailStmt,
		getUserByUsernameStmt:     q.getUserByUsernameStmt,
		incrementLikesCountStmt:   q.incrementLikesCountStmt,
		incrementViewsStmt:        q.incrementViewsStmt,
		likeSnippetStmt:           q.likeSnippetStmt,
		recordViewStmt:            q.recordViewStmt,
		saveSnippetStmt:           q.saveSnippetStmt,
		updateLikesCountStmt:      q.updateLikesCountStmt,
		updateSessionExpiryStmt:   q.updateSessionExpiryStmt,
		updateSnippetStmt:         q.updateSnippetStmt,
		updateUserAvatarStmt:      q.updateUserAvatarStmt,
		updateUserInfoStmt:        q.updateUserInfoStmt,
		updateUserPasswordStmt:    q.updateUserPasswordStmt,
	}
}
