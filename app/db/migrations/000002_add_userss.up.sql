CREATE TABLE "users" (
  "name" varchar PRIMARY KEY,
  "password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_change_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
  "create_at" timestamptz NOT NULL DEFAULT (now())
);

-- 单引号引用字符串。用双引号引用的字符串被解释为一个识别符。 所以上面時間要單引號
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("name");


-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");