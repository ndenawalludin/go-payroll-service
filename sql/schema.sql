CREATE TABLE employees
(
    id          SERIAL PRIMARY KEY,
    code        VARCHAR(50) UNIQUE  NOT NULL,
    full_name   VARCHAR(255)        NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    base_salary BIGINT              NOT NULL,
    allowance   BIGINT              NOT NULL DEFAULT 0,
    is_active   BOOLEAN             NOT NULL DEFAULT TRUE,
    hire_date   DATE                NOT NULL,
    created_at  TIMESTAMP           NOT NULL,
    updated_at  TIMESTAMP           NOT NULL
);

CREATE TABLE payroll_periods
(
    id         SERIAL PRIMARY KEY,
    code       VARCHAR(50) UNIQUE NOT NULL,
    start_date DATE               NOT NULL,
    end_date   DATE               NOT NULL,
    closed     BOOLEAN            NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP          NOT NULL,
    updated_at TIMESTAMP          NOT NULL
);

CREATE TABLE payslips
(
    id                SERIAL PRIMARY KEY,
    employee_id       INTEGER NOT NULL REFERENCES employees (id),
    payroll_period_id INTEGER NOT NULL REFERENCES payroll_periods (id),
    base_salary       BIGINT  NOT NULL,
    allowance         BIGINT  NOT NULL,
    deduction         BIGINT  NOT NULL,
    net_salary        BIGINT  NOT NULL
);
