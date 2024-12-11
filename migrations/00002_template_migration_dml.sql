-- +goose Up
-- +goose StatementBegin
-- Insert initial values
INSERT INTO
    template_example (
        first_column,
        second_column,
        created_at,
        updated_at
    )
VALUES (
        'root',
        '',
        '2024-11-29 11:30:30',
        '2024-11-29 11:30:30'
    ),
    (
        'under root',
        'child',
        '2024-11-29 11:30:30',
        '2024-11-29 11:30:30'
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM template_example
WHERE
    first_column IN ('root', 'under root');
-- +goose StatementEnd