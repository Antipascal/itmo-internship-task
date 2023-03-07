CREATE TABLE UserAuthToken (
    AccessToken VARCHAR(550) PRIMARY KEY,
    ISU INTEGER
);

CREATE TABLE UserInfo (
    ISU INTEGER PRIMARY KEY,
    GivenName varchar(100),
    MiddleName varchar(100),
    FamilyName varchar(100),
    Email varchar(50),
    PhoneNumber varchar(20)
);