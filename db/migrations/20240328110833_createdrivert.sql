-- +goose Up
-- +goose StatementBegin
CREATE DOMAIN nationalIDNumberType  as varchar(10)
CHECK (
    VALUE ~ '(^[A-Z]\d{9}$)'
);

CREATE TABLE DriverT(
    id bigint PRIMARY KEY references UserT(id),
    plateNum varchar NOT NULL unique ,
    -- BLABLABLA
    insurances varchar,
    registration varchar,
    driverLicense varchar,
    truckLicense varchar,
    nationalIDNumber nationalIDNumberType NOT NULL Unique,
    percentage smallint NOT NULL DEFAULT 20, 
    lastAlert bigint references alertT(id),
    approved_date Timestamp
  );
CREATE OR REPLACE FUNCTION test()
    RETURNS trigger AS
$$
BEGIN
    Update usert set last_modified_date = NOW() where id = OLD.ID;
 RETURN NEW;
END;

$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER updateDM 
AFTER UPDATE ON drivert 
FOR EACH ROW
 WHEN (pg_trigger_depth() < 1)  -- !
EXECUTE PROCEDURE test();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS DriverT;
DROP DOMAIN IF EXISTS nationalIDNumberType;
-- +goose StatementEnd
