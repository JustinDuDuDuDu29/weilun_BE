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
    role smallint NOT NULL,
    -- 0~100:super admin
    -- 101~200:company
    -- 300+:driver
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW()
  );

CREATE TABLE CMPT(
    id BIGSERIAL PRIMARY KEY,
    name varchar(25) NOT NULL Unique,
    -- phoneNum phoneType NOT NULL Unique,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
  );

CREATE TABLE CMPInChargeT(
    id BIGSERIAL PRIMARY KEY,
    userI BIGSERIAL references UserT(id),
    cmpID BIGSERIAL references CMPT(id),
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW(), 

    PRIMARY KEY(userID, cmpID)
  );

CREATE TABLE DriverT(
    id BIGSERIAL PRIMARY KEY references UserT(id) ,
    name varchar(25) NOT NULL,
    phoneNum phoneType NOT NULL Unique,
    belongCMP BIGSERIAL references CMPT(id),
    -- BLABLABLA
    percentage smallint NOT NULL DEFAULT 20, 
    -- BLABLABLA
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
  );


CREATE TABLE JobListT(
    id BIGSERIAL PRIMARY KEY,
    from_loc varchar(40) NOT NULL,
    mid varchar(40),
    to_loc varchar(40) NOT NULL,
    
    price  smallint NOT NULL,
    estimated smallint NOT NULL,
    remaining smallint NOT NULL,

    belongCMP BIGSERIAL references CMPT(id) NOT NULL,
    source varchar(40) NOT NULL,
    jobDate Timestamp NOT NULL DEFAULT NOW(),
    memo varchar(60),
    -- BLABLABLA
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
  );

 
CREATE TABLE JobTakenT(
    id BIGSERIAL PRIMARY KEY,
    jobID BIGSERIAL references JobListT(id),
    driverID BIGSERIAL references DriverT(id),
    percentage smallint,
    is_Finished Boolean NOT NULL,
    finished_date Timestamp,
    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
  );


CREATE TABLE revenueT(
    id BIGSERIAL PRIMARY KEY,
    jobTakenID BIGSERIAL references JobTakenT(id),
    driverEarn smallint NOT NULL,

    create_date Timestamp NOT NULL DEFAULT NOW(),
    deleted_date Timestamp,
    last_modified_date Timestamp NOT NULL DEFAULT NOW() 
  );


