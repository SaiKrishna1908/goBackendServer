-- +goose Up
-- +goose StatementBegin
alter table workout_entries alter column reps drop not null;
alter table workout_entries alter column sets drop not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
