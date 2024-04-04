-- +goose Up
-- +goose StatementBegin
CREATE DOMAIN nationalIDNumberType  as varchar(10)
CHECK (
    VALUE ~ '(^[A-Z]\d{9}$)'
);

CREATE TABLE DriverT(
    id BIGSERIAL PRIMARY KEY,
    userId bigint references UserT(id),
    -- BLABLABLA
    percentage smallint NOT NULL DEFAULT 20, 
    nationalIDNumber nationalIDNumberType NOT NULL
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS DriverT;
DROP DOMAIN IF EXISTS nationalIDNumberType;
-- +goose StatementEnd
