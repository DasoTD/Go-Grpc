CREATE TABLE "accounts" (
  "username" varchar PRIMARY KEY,
  "firstname" varchar NOT NULL,
  "lastname" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,  
  "created_at" timestamptz NOT NULL DEFAULT (now())
);