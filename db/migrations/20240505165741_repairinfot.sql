-- +goose Up
-- +goose StatementBegin
Create Table RepairInfoT(
    id BIGSERIAL PRIMARY KEY,
    repairID BigInt references RepairT(id) NOT NULL,
    itemName varchar(20) NOT NULL,
    quantity int DEFAULT(1) NOT NULL,
    totalPrice BigInt DEFAULT(0) NOT NULL,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop Table if exists RepairInfoT;
-- +goose StatementEnd
