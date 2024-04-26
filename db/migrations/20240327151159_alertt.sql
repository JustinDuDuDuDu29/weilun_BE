-- +goose Up
-- +goose StatementBegin
Create TABLE alertT (
    id BIGSERIAL PRIMARY KEY,
    alert varchar not NULL,
    belongCMP bigint references CMPT(id) not NULL,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop Table if exists alertT;
-- +goose StatementEnd
