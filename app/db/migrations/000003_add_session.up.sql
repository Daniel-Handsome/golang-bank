CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "username" varchar NOT NULL,
    "refresh_token" varchar NOT NULL,
    "user_agent"  varchar NOT NULL,
    "client_id" varchar NOT NULL,
    "isblacked" boolean NOT NULL DEFAULT false,
    "expires_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");