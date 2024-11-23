-- +goose Up
-- +goose StatementBegin
Create Table GasInfoT(
  id BIGSERIAL PRIMARY KEY,
  gasID BigInt references GasT(id) NOT NULL,
  itemName varchar NOT NULL,
  quantity int DEFAULT(1) NOT NULL,
  totalPrice BigInt NOT NULL,
  create_date Timestamp NOT NULL DEFAULT NOW(),
  deleted_date Timestamp,
  last_modified_date Timestamp NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop Table if exists GasInfoT;
-- +goose StatementEnd
