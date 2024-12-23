-- +goose Up
-- +goose StatementBegin
CREATE TABLE RepairT(
  id BIGSERIAL PRIMARY KEY,
  driverID BigInt references DriverT(id) NOT NULL,
  pic varchar,
  place varchar NOT null,
  create_date Timestamp NOT NULL DEFAULT NOW(),
  approved_date Timestamp,
  deleted_date Timestamp,
  last_modified_date Timestamp NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
Drop Table if exists RepairT;
-- +goose StatementEnd