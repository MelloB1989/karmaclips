CREATE TABLE IF NOT EXISTS "generations" (
	"id" varchar PRIMARY KEY NOT NULL,
	"created_by" varchar,
	"credits_used" integer DEFAULT 0,
	"timestamp" timestamp DEFAULT now(),
	"media_uri" text DEFAULT '',
	"type" varchar DEFAULT '',
	"meta" json DEFAULT '{}'::json,
	CONSTRAINT "generations_id_unique" UNIQUE("id")
);
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "generations" ADD CONSTRAINT "generations_created_by_users_id_fk" FOREIGN KEY ("created_by") REFERENCES "public"."users"("id") ON DELETE no action ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
