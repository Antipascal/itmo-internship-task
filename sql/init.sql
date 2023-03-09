CREATE TABLE user_auth (
    access_token VARCHAR(256) PRIMARY KEY NOT NULL,
    ISU INTEGER NOT NULL
);

CREATE TABLE user_info (
    ISU INTEGER PRIMARY KEY NOT NULL,
    given_name varchar(100) NOT NULL,
    middle_name varchar(100) NOT NULL,
    family_name varchar(100),
    email varchar(50),
    phone_number varchar(20)
);

CREATE TABLE admins (
    ISU INTEGER PRIMARY KEY NOT NULL
);