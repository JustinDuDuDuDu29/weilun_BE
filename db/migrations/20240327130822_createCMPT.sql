-- +goose Up
-- +goose StatementBegin
CREATE TABLE CMPT(
    id BIGSERIAL PRIMARY KEY,
    name varchar(25) NOT NULL Unique,
    -- phoneNum phoneType NOT NULL Unique,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop TABLE If Exists CMPT;
-- +goose StatementEnd
