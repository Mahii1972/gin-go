-- name: GetEmployee :one
SELECT * FROM employees 
WHERE id = $1;

-- name: ListEmployees :many
SELECT * FROM employees
ORDER BY last_name, first_name;

-- name: CreateEmployee :one
INSERT INTO employees (first_name, last_name, email, hire_date, salary)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateEmployee :one
UPDATE employees
SET first_name = $2,
    last_name = $3,
    email = $4,
    hire_date = $5,
    salary = $6
WHERE id = $1
RETURNING *;

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;
