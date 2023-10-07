DROP DATABASE IF EXISTS brimoveapp;

CREATE DATABASE brimoveapp;

\c brimoveapp;

CREATE TABLE "accounts" (
                            "id" INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
                            "uname" varchar(256) NOT NULL,
                            "pass" varchar(256) NOT NULL,
                            "access_token" varchar(256) NULL,
                            "refresh_token" varchar(256) NULL
);