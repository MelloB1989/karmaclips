CREATE TABLE IF NOT EXISTS "referrals" (
	"id" varchar PRIMARY KEY NOT NULL,
	"created_by" varchar,
	"timestamp" timestamp DEFAULT now(),
	CONSTRAINT "referrals_id_unique" UNIQUE("id")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "users" (
	"id" varchar PRIMARY KEY NOT NULL,
	"email" varchar DEFAULT '',
	"phone" varchar,
	"location" varchar,
	"referral_code" varchar,
	CONSTRAINT "users_id_unique" UNIQUE("id"),
	CONSTRAINT "users_phone_unique" UNIQUE("phone"),
	CONSTRAINT "users_referral_code_unique" UNIQUE("referral_code")
);
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "referrals" ADD CONSTRAINT "referrals_created_by_users_id_fk" FOREIGN KEY ("created_by") REFERENCES "public"."users"("id") ON DELETE no action ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
