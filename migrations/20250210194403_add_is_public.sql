-- +goose Up
-- +goose StatementBegin
ALTER TABLE PODS ADD COLUMN is_public BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE PODS DROP COLUMN is_public;
-- +goose StatementEnd
