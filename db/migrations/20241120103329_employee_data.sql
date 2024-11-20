-- +goose Up
-- +goose StatementBegin
INSERT INTO employees (first_name, last_name, email, hire_date, salary) VALUES
    ('John', 'Doe', 'john.doe@company.com', '2023-01-15', 75000.00),
    ('Jane', 'Smith', 'jane.smith@company.com', '2023-02-01', 82000.00),
    ('Michael', 'Johnson', 'michael.j@company.com', '2023-03-10', 65000.00),
    ('Sarah', 'Williams', 'sarah.w@company.com', '2023-04-05', 71000.00),
    ('Robert', 'Brown', 'robert.b@company.com', '2023-05-20', 68000.00);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM employees WHERE email IN (
    'john.doe@company.com',
    'jane.smith@company.com', 
    'michael.j@company.com',
    'sarah.w@company.com',
    'robert.b@company.com'
);
-- +goose StatementEnd
