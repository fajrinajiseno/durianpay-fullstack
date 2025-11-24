package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newMockRepo(t *testing.T) (*Payment, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	repo := NewPaymentRepo(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestGetPayments_Success(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "merchant", "amount", "status", "created_at"}).
		AddRow("p1", "m1", 100.0, "pending", time.Now()).
		AddRow("p2", "m2", 200.0, "completed", time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, merchant, amount, status, created_at FROM payments WHERE status = ? ORDER BY created_at ASC LIMIT ? OFFSET ?")).
		WithArgs("completed", 10, 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(1) FROM payments")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(1) FROM payments WHERE status = 'completed'")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(1) FROM payments WHERE status = 'failed'")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	items, total, totalSuccess, totalFailed, err := repo.GetPayments("completed", "created_at", 10, 1)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, 2, total)
	assert.Equal(t, 1, totalSuccess)
	assert.Equal(t, 0, totalFailed)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetPayments_QueryError(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	// Simulate DB query error on main SELECT
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, merchant, amount, status, created_at FROM payments ORDER BY created_at ASC")).
		WillReturnError(errors.New("db select failed"))

	_, _, _, _, err := repo.GetPayments("", "created_at", 0, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestReview_Success(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM payments WHERE id = ?")).
		WithArgs("p1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("p1"))

	msg, err := repo.Review("p1")
	assert.NoError(t, err)
	assert.Equal(t, "Success Review", msg)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestReview_NotFound(t *testing.T) {
	repo, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM payments WHERE id = ?")).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.Review("missing")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment not found")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
