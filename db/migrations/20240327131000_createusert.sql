-- +goose Up
-- +goose StatementBegin

CREATE DOMAIN phoneType AS TEXT
CHECK(
-- (^\+[1-9]{1}[0-9]{3,14}$)|(^09\d{2}(\d{6}|-\d{3}-\d{3})$) ???

   VALUE ~ '(^\+[1-9]{1}[0-9]{3,14}$)'
OR VALUE ~ '^09\d{2}(\d{6}|-\d{3}-\d{3})$'
);

CREATE TABLE UserT(
    id BIGSERIAL PRIMARY KEY,
    phoneNum phoneType NOT NULL Unique,
    pwd  varchar(76) NOT NULL,
    name varchar(25) NOT NULL,
    belongCMP bigint NOT NULL references CMPT(id),
    seed varchar(20), 
    role smallint NOT NULL,
    -- 0~100:super admin
    -- 101~200:company
    -- 300+:driver
    initPwdChanged Boolean NOT NULL DEFAULT false,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW()
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop TABLE If Exists UserT;
Drop DOMAIN If Exists phoneType;
-- +goose StatementEnd
