CREATE TABLE "quiz" (
  "id" varchar(255) PRIMARY KEY,
  "type" int NOT NULL default 1,
  "owner" varchar(255) NOT NULL,
  "content" text,
  "hash_content" varchar(255) NOT NULL,
  "answer" text,
  "hash_answer" varchar(255),
  "reward" bigint,
  "winner" varchar (255),
  "duration" bigint not null,
  "duration_voting" bigint,
  "timestamp_created" bigint NOT NULL DEFAULT 0,
  "status" int NOT NULL default -1,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "answer" (
  "id" varchar(255) PRIMARY KEY,
  "index" int not null,
  "quiz_id" varchar(255) NOT NULL,
  "owner" varchar(255) NOT NULL,
  "content" text,
  "vote" int not null default 0,
  "hash_content" varchar(255) NOT NULL,
  "timestamp_created" bigint NOT NULL DEFAULT 0,
  "status" int NOT NULL default -1,
  "created_at" timestamp DEFAULT (now())
);

ALTER TABLE "answer" ADD FOREIGN KEY ("quiz_id") REFERENCES "quiz" ("id");

CREATE TABLE "chain_config" (
    "id" serial primary key,
    "chain_id" varchar(255) NOT NULL,
    "contract_address" varchar(255) NOT NULL,
    "rpc_url" varchar(500) not null,
    "wss_url" varchar(500) not null,
    "block_number" bigint NOT NULL DEFAULT 0,
    "step_number" bigint NOT NULL DEFAULT -1,
    "created_at" timestamp DEFAULT (now()),
    "updated_at" timestamp DEFAULT (now())
);