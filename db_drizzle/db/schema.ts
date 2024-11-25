import {
  pgTable,
  json,
  varchar,
  timestamp,
  unique,
  text,
  integer,
} from "drizzle-orm/pg-core";
import { sql } from "drizzle-orm";

export const users = pgTable("users", {
  id: varchar("id").primaryKey().unique(),
  email: varchar("email").default(""),
  name: varchar("name").default(""),
  password: varchar("password").default(""),
  phone: varchar("phone").unique(),
  location: varchar("location"),
  referral_code: varchar("referral_code").unique(),
});

export const referrals = pgTable("referrals", {
  id: varchar("id").primaryKey().unique(),
  createdBy: varchar("created_by").references(() => users.id),
  timestamp: timestamp("timestamp").default(sql`now()`),
});

export const generations = pgTable("generations", {
  id: varchar("id").primaryKey().unique(),
  createdBy: varchar("created_by").references(() => users.id),
  creditsUsed: integer("credits_used").default(0),
  timestamp: timestamp("timestamp").default(sql`now()`),
  mediaUri: text("media_uri").default(""),
  type: varchar("type").default(""),
  meta: json("meta").default({}),
});
