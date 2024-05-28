-- +goose Up
-- +goose StatementBegin
CREATE TABLE repairT(
    id BIGSERIAL PRIMARY KEY,
    type varchar NOT NULL, 
    driverID BigInt references DriverT(id) NOT NULL,
    repairInfo JSONB NOT NULL,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    approved_date Timestamp,
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop Table if exists repairT;
-- +goose StatementEnd
