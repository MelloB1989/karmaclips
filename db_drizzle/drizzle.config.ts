import type { Config } from "drizzle-kit";

export default {
  schema: "./db/schema.ts",
  out: "./drizzle",
  dialect: "postgresql", // 'pg' | 'mysql2' | 'better-sqlite' | 'libsql' | 'turso'
  dbCredentials: {
    host: "192.168.29.73",
    user: "mellob",
    password: "mellob1989",
    database: "karmaclips",
    ssl: false,
  },
} satisfies Config;
