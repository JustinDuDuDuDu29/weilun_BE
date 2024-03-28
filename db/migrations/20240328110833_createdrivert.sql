-- +goose Up
-- +goose StatementBegin
CREATE DOMAIN nationalIDNumberType  as varchar(10)
CHECK (
    VALUE ~ '(^[A-Z]\d{9}$)'
);

CREATE TABLE DriverT(
    id BIGSERIAL PRIMARY KEY,
    userId bigint references UserT(id) ,
    -- BLABLABLA
    percentage smallint NOT NULL DEFAULT 20, 
    nationalIDNumber nationalIDNumberType NOT NULL,
    -- BLABLABLA
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS DriverT;
-- +goose StatementEnd
