-- +goose Up
-- +goose StatementBegin
CREATE TABLE CMPInChargeT(

    id BIGSERIAL PRIMARY KEY,
    
    userID BIGINT references UserT(id),
    cmpID BIGINT references CMPT(id),

    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW(), 

    Unique(userID, cmpID)

  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS CMPInChargeT;
-- +goose StatementEnd
