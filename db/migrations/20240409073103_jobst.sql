-- +goose Up
-- +goose StatementBegin
CREATE TABLE JobsT(
    id BIGSERIAL PRIMARY KEY,
    from_loc varchar(40) NOT NULL,
    mid varchar(40),
    to_loc varchar(40) NOT NULL,
    
    price  smallint NOT NULL,
    estimated smallint NOT NULL,
    remaining smallint NOT NULL check(remaining>=0),

    belongCMP BigInt references CMPT(id) NOT NULL,
    source varchar(40) NOT NULL,
    jobDate Timestamp NOT NULL DEFAULT NOW(),
    memo varchar(60),
    -- BLABLABLA
    create_date Timestamp NOT NULL DEFAULT NOW(),
    end_date Timestamp,
    deleted_date Timestamp,
    finished_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop Table if exists JobsT;
-- +goose StatementEnd
