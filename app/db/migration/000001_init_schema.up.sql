CREATE TABLE "quiz" (
  "id" varchar(255) PRIMARY KEY,
  "owner" varchar(255) NOT NULL,
  "content" text,
  "hash_content" varchar(255) NOT NULL,
  "answer" text,
  "hash_answer" varchar(255),
  "timestamp_created" bigint NOT NULL DEFAULT 0,
  "status" int NOT NULL default -1,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "answer" (
  "id" varchar(255) PRIMARY KEY,
  "quiz_id" varchar(255) NOT NULL,
  "owner" varchar(255) NOT NULL,
  "content" text,
  "hash_content" varchar(255) NOT NULL,
  "timestamp_created" bigint NOT NULL DEFAULT 0,
  "status" int NOT NULL default -1,
  "created_at" timestamp DEFAULT (now())
);

ALTER TABLE "answer" ADD FOREIGN KEY ("quiz_id") REFERENCES "quiz" ("id");

CREATE TABLE "event_log" (
    "chain_id" varchar(255) NOT NULL,
    "contract_address" varchar(255) NOT NULL,
    "block_number" bigint NOT NULL DEFAULT 0,
    "step_number" bigint NOT NULL DEFAULT 1000,
    "created_at" timestamp DEFAULT (now()),
    "updated_at" timestamp DEFAULT (now()),
    PRIMARY KEY ("chain_id", "contract_address")
);