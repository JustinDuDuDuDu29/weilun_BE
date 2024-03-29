-- +goose Up
-- +goose StatementBegin
CREATE TABLE LoginT(
    id BIGSERIAL PRIMARY KEY,
    userID bigint references UserT(id),
    create_date Timestamp NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop TABLE IF Exists LoginT;
-- +goose StatementEnd
