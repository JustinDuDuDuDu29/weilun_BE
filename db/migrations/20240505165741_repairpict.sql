-- +goose Up
-- +goose StatementBegin
Create Table RepairPicT(
    id BIGSERIAL PRIMARY KEY,
    repair_id BigInt references repairT(id) NOT NULL,
    pic varchar,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    approved_date Timestamp,
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop Table if exists RepairPicT;
-- +goose StatementEnd
