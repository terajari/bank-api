CREATE TABLE "accounts" (
  "id" varchar(100) PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" varchar(100) PRIMARY KEY,
  "account_id" varchar(100) NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" varchar(100) PRIMARY KEY,
  "sender_id" varchar(100) NOT NULL,
  "receiver_id" varchar(100) NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("sender_id");

CREATE INDEX ON "transfers" ("receiver_id");

CREATE INDEX ON "transfers" ("sender_id", "receiver_id");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("sender_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("receiver_id") REFERENCES "accounts" ("id");