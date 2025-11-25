package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/fajrinajiseno/mygolangapp/internal/entity"
)

//go:generate mockgen -source payment.go -destination mock/payment_mock.go -package=mock
type PaymentRepository interface {
	GetPayments(status, id string, sortExpr string, limit, offset int) ([]*entity.Payment, *entity.PaymentSummary, error)
	Review(id string) (string, error)
}

type Payment struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *Payment {
	return &Payment{db: db}
}

func (r *Payment) GetPayments(status, id string, sortExpr string, limit, offset int) ([]*entity.Payment, *entity.PaymentSummary, error) {
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
	where := []string{}
	args := []interface{}{}
	if status != "" {
		where = append(where, "status = ?")
		args = append(args, status)
	}
	if id != "" {
		where = append(where, "id = ?")
		args = append(args, id)
	}
	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
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
		return nil, nil, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}
	defer rows.Close()
	res := []*entity.Payment{}
	for rows.Next() {
		var p entity.Payment
		if err := rows.Scan(&p.ID, &p.Merchant, &p.Amount, &p.Status, &p.CreatedAt); err != nil {
			return nil, nil, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
		}
		res = append(res, &p)
	}

	var totalByFiler, total, totalCompleted, totalFailed, totalPending int

	qt := "SELECT COUNT(1) FROM payments"
	whereT := []string{}
	argsT := []interface{}{}
	if status != "" {
		whereT = append(whereT, "status = ?")
		argsT = append(argsT, status)
	}
	if id != "" {
		whereT = append(whereT, "id = ?")
		argsT = append(argsT, id)
	}
	if len(whereT) > 0 {
		qt += " WHERE " + strings.Join(where, " AND ")
	}
	row := r.db.QueryRow(qt, argsT...)
	err = row.Scan(&totalByFiler)
	if err != nil {
		return nil, nil, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}

	row = r.db.QueryRow("SELECT COUNT(1) FROM payments")
	err = row.Scan(&total)
	if err != nil {
		return nil, nil, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}

	row = r.db.QueryRow("SELECT COUNT(1) FROM payments WHERE status = 'completed'")
	err = row.Scan(&totalCompleted)
	if err != nil {
		return nil, nil, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}

	row = r.db.QueryRow("SELECT COUNT(1) FROM payments WHERE status = 'failed'")
	err = row.Scan(&totalFailed)
	if err != nil {
		return nil, nil, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}

	row = r.db.QueryRow("SELECT COUNT(1) FROM payments WHERE status = 'pending'")
	err = row.Scan(&totalPending)
	if err != nil {
		return nil, nil, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}

	return res, &entity.PaymentSummary{
		TotalByFiler:   totalByFiler,
		Total:          total,
		TotalCompleted: totalCompleted,
		TotalFailed:    totalFailed,
		TotalPending:   totalPending,
	}, nil
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
