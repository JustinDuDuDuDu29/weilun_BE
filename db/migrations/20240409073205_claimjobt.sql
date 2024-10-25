-- +goose Up
-- +goose StatementBegin
CREATE TABLE ClaimJobT(
    id BIGSERIAL PRIMARY KEY,
    jobID BigInt references JobsT(id) NOT NULL,
    driverID BigInt references DriverT(id) NOT NULL,
    -- percentage integer,
    finishPic varchar,
    memo  varchar,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    deleted_by BigInt references UserT(id),
    finished_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW(),
    approved_By BigInt references UserT(id),
    approved_date Timestamp 
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop Table if exists ClaimJobT;
-- +goose StatementEnd
