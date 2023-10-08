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

CREATE TABLE "saving" (
      "account_number" varchar(256) PRIMARY KEY NOT NULL,
      "account_type" varchar(256) NOT NULL,
      "branch_code" varchar(256) NOT NULL,
      "short_name" varchar(300) NOT NULL,
      "currency" varchar(3) NOT NULL,
      "cbal" varchar(256) NOT NULL,
      "hold" varchar(256) NOT NULL,
      "opening_date" varchar(256) NOT NULL,
      "product_group" varchar(256) NOT NULL,
      "product_name" varchar(256) NOT NULL,
      "status" varchar(1) NOT NULL
);

CREATE TABLE "brinjournalseq" (
      "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY NOT NULL,
      "branch_code" varchar(256) NOT NULL,
      "joirnalseq" integer NOT NULL
);

CREATE INDEX ON "saving" ("account_number");

CREATE INDEX ON "brinjournalseq" ("branch_code");

INSERT INTO accounts(uname, pass, access_token, refresh_token) VALUES
                                                                   ('deo', '123', null, null),
                                                                   ('reo', '123', null, null);

INSERT INTO saving (account_number, account_type, branch_code, short_name, currency, cbal, hold, opening_date, product_group, product_name, status) VALUES
                ('045202000001809', 'S', '0999', 'DEOTAMA GUNANSA', 'USD', '500.00', '00.00', '2020-10-10 10:00:00', '10001', 'Britama Saving', '1'),
                ('045202000001808', 'S', '0999', 'REO SAHOBBY', 'USD', '500.00', '00.00', '2020-10-10 10:00:00', '10001', 'Britama Saving', '1');


INSERT INTO brinjournalseq(branch_code, joirnalseq) VALUES ('09999', 1);
