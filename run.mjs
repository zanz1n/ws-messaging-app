import { spawnSync } from "child_process";

spawnSync("pnpm", ["build:client"], { stdio: "inherit" });
spawnSync("go", ["run", "."], { stdio: "inherit" });
