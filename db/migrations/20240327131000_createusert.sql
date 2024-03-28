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
    userName varchar(20)  NOT NULL Unique,
    pwd  varchar(76) NOT NULL,
    name varchar(25) NOT NULL,
    belongCMP bigint references CMPT(id),
    phoneNum phoneType NOT NULL Unique,
    role smallint NOT NULL,
    -- 0~100:super admin
    -- 101~200:company
    -- 300+:driver
    loginTimes integer NOT NULL DEFAULT 0,
    initPwdChanged Boolean NOT NULL DEFAULT false,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW(),

    CHECK (loginTimes>=0)
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop TABLE If Exists UserT;
-- +goose StatementEnd
