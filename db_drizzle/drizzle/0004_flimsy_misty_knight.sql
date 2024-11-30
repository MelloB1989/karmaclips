CREATE TABLE IF NOT EXISTS "models" (
	"model_id" varchar PRIMARY KEY NOT NULL,
	"type" varchar PRIMARY KEY NOT NULL,
	"provider" varchar DEFAULT '',
	"pre_prompt" text DEFAULT '',
	"banner" text DEFAULT '',
	"description" text DEFAULT '',
	"credits_per_gen" integer DEFAULT 0,
	CONSTRAINT "models_model_id_unique" UNIQUE("model_id"),
	CONSTRAINT "models_type_unique" UNIQUE("type")
);
