/**
 * Ensures dev dependencies (vite, etc.) exist before dev/build.
 * Wails may skip frontend:dev:install when node_modules exists and package.json.md5 matches,
 * but node_modules can still be incomplete, causing "vite is not recognized" and Svelte IDE errors.
 */
import { spawnSync } from "node:child_process";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const frontendRoot = path.resolve(__dirname, "..");
const viteMarker = path.join(frontendRoot, "node_modules", "vite", "package.json");

if (!fs.existsSync(viteMarker)) {
  const r = spawnSync("npm", ["ci"], { cwd: frontendRoot, stdio: "inherit", shell: true });
  if ((r.status ?? 1) !== 0) {
    process.exit(r.status ?? 1);
  }
}

process.exit(0);
