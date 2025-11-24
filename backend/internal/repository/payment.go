package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/fajrinajiseno/mygolangapp/internal/entity"
)

//go:generate mockgen -source payment.go -destination mocks/payment_mock.go -package=mock
type PaymentRepository interface {
	GetPayments(status, sortExpr string, limit, offset int) ([]*entity.Payment, int, int, int, error)
	Review(id string) (string, error)
}

type Payment struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *Payment {
	return &Payment{db: db}
}

func (r *Payment) GetPayments(status, sortExpr string, limit, offset int) ([]*entity.Payment, int, int, int, error) {
	// whitelist sort columns to avoid SQL injection
	allowedSort := map[string]string{
		"amount":     "amount",
		"created_at": "created_at",
	}

	// default sort
	col := "created_at"
	dir := "ASC"

	// parse sort expression (e.g., "-created_at" => DESC)
	if sortExpr != "" {
		s := strings.TrimSpace(sortExpr)
		if strings.HasPrefix(s, "-") {
			dir = "DESC"
			s = strings.TrimPrefix(s, "-")
		}
		if v, ok := allowedSort[s]; ok {
			col = v
		}
	}

	q := "SELECT id, merchant, amount, status, created_at FROM payments"
	args := []interface{}{}
	if status != "" {
		q += " WHERE status = ?"
		args = append(args, status)
	}
	if col != "" && dir != "" {
		q += fmt.Sprintf(" ORDER BY %s %s", col, dir)
	}
	if limit > 0 {
		q += " LIMIT ?"
		args = append(args, limit)
	}
	if offset > 0 {
		q += " OFFSET ?"
		args = append(args, offset)
	}

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, 0, 0, 0, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}
	defer rows.Close()
	res := []*entity.Payment{}
	for rows.Next() {
		var p entity.Payment
		if err := rows.Scan(&p.ID, &p.Merchant, &p.Amount, &p.Status, &p.CreatedAt); err != nil {
			return nil, 0, 0, 0, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
		}
		res = append(res, &p)
	}

	var total, totalSuccess, totalFailed int

	row := r.db.QueryRow("SELECT COUNT(1) FROM payments")
	err = row.Scan(&total)
	if err != nil {
		return nil, 0, 0, 0, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}

	row = r.db.QueryRow("SELECT COUNT(1) FROM payments WHERE status = 'completed'")
	err = row.Scan(&totalSuccess)
	if err != nil {
		return nil, 0, 0, 0, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}

	row = r.db.QueryRow("SELECT COUNT(1) FROM payments WHERE status = 'failed'")
	err = row.Scan(&totalFailed)
	if err != nil {
		return nil, 0, 0, 0, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}

	return res, total, totalSuccess, totalFailed, nil
}

func (r *Payment) Review(id string) (string, error) {
	row := r.db.QueryRow("SELECT id FROM payments WHERE id = ?", id)
	var p entity.Payment
	if err := row.Scan(&p.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", entity.ErrorNotFound("payment not found")
		}
		return "", entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}
	return "Success Review", nil
}
