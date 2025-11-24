package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-payroll-service/cmd/internal/payroll/model/domain"
	"go-payroll-service/cmd/internal/payroll/util"
	"time"
)

type EmployeeRepository interface {
	List(ctx context.Context) ([]domain.Employee, error)
	Create(ctx context.Context, employee domain.Employee) (domain.Employee, error)
	GetByID(ctx context.Context, id int64) (domain.Employee, error)
	Update(ctx context.Context, employee domain.Employee) (domain.Employee, error)
	Delete(ctx context.Context, id int64) error
}

type employeeRepository struct {
	db *sql.DB
}

func (r *employeeRepository) List(ctx context.Context) ([]domain.Employee, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, code, full_name, email, base_salary, allowance, is_active,
		       hire_date, created_at, updated_at
		FROM employees
		ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Employee
	for rows.Next() {
		var e domain.Employee
		if err := rows.Scan(
			&e.ID, &e.Code, &e.FullName, &e.Email,
			&e.BaseSalary, &e.Allowance, &e.IsActive,
			&e.HireDate, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, rows.Err()
}

func (r *employeeRepository) Create(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	e.IsActive = true

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO employees(code, full_name, email, base_salary, allowance, is_active, hire_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		e.Code, e.FullName, e.Email, e.BaseSalary, e.Allowance, e.IsActive, e.HireDate, e.CreatedAt, e.UpdatedAt,
	).Scan(&e.ID)
	if err != nil {
		return domain.Employee{}, err
	}
	return e, nil
}

func (r employeeRepository) GetByID(ctx context.Context, id int64) (domain.Employee, error) {
	var e domain.Employee
	err := r.db.QueryRowContext(ctx, `
		SELECT id, code, full_name, email, base_salary, allowance, is_active, hire_date, created_at, updated_at
		FROM employees
		WHERE id = $1`, id,
	).Scan(
		&e.ID, &e.Code, &e.FullName, &e.Email,
		&e.BaseSalary, &e.Allowance, &e.IsActive,
		&e.HireDate, &e.CreatedAt, &e.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Employee{}, util.ErrNotFound
	}

	if err != nil {
		return domain.Employee{}, err
	}

	return e, nil
}

func (r employeeRepository) Update(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	e.UpdatedAt = time.Now()

	res, err := r.db.ExecContext(ctx, `
		UPDATE employees
		SET full_name=$1, email=$2, base_salary=$3, allowance=$4, is_active=$5, hire_date=$6, updated_at=$7
		WHERE id = $8`,
		e.FullName, e.Email, e.BaseSalary, e.Allowance, e.IsActive,
		e.HireDate, e.UpdatedAt, e.ID,
	)

	if err != nil {
		return domain.Employee{}, err
	}

	aff, err := res.RowsAffected()
	if err == nil && aff == 0 {
		return domain.Employee{}, util.ErrNotFound
	}

	return e, nil
}

func (r employeeRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM employees WHERE id = $1`, id)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err == nil && aff == 0 {
		return util.ErrNotFound
	}

	return nil
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}
