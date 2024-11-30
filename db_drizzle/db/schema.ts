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

export const ai_services = pgTable("ai_services", {
  aid: varchar("aid").primaryKey().unique(),
  type: varchar("type").primaryKey().unique(),
  provider: varchar("provider").default(""),
  pre_prompt: text("pre_prompt").default(""),
  banner: text("banner").default(""),
  description: text("description").default(""),
  credits_per_gen: integer("credits_per_gen").default(0),
});

export const wallet = pgTable("wallet", {
  wallet_id: varchar("wallet_id").primaryKey().unique(),
  user_id: varchar("user_id").references(() => users.id),
  balance: integer("balance").default(0),
});

export const transactions = pgTable("transactions", {
  trx_id: varchar("trx_id").primaryKey().unique(),
  user_id: varchar("user_id").references(() => users.id),
  wallet_id: varchar("wallet_id").references(() => wallet.wallet_id),
  amount: integer("amount").default(0),
  description: text("description").default(""),
  timestamp: timestamp("timestamp").default(sql`now()`),
});

export const projects = pgTable("projects", {
  project_id: varchar("project_id").primaryKey().unique(),
  type: varchar("type").primaryKey().unique(),
  user_id: varchar("user_id").references(() => users.id),
  name: varchar("name").default(""),
  description: text("description").default(""),
  meta: json("meta").default({}),
  states: json("states").default({}),
});
