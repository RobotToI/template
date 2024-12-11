-- +goose Up
-- +goose StatementBegin
CREATE TABLE template_example (
    id SERIAL NOT NULL PRIMARY KEY,
    first_column TEXT,
    second_column TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE template_example;
-- +goose StatementEnd