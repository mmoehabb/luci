import * as fs from "fs"
import { argv } from "process";
import { exec } from "child_process";

import type { AnnotatedAction, Config, VagueAction } from "./types";
import { dig, isAnnotatedAction, isVagueAction } from "./util";
import { printAction, printUsage } from "./print";

function main() {
  if (!fs.existsSync("luci.config.json")) {
    console.error("luci.config.json file not found!");
    return;
  }

  try {
    const config: Config = JSON.parse(fs.readFileSync("luci.config.json", "utf8"));
    if (argv.length <= 2) printUsage(config);
    else act(config, argv.slice(2));
  } catch (err) {
    console.error(err);
  }
}

function act(config: Config, inputs: Array<string>) {
  const action = dig(config, inputs);
  if (!action) return printAction(config, inputs.slice(0, inputs.length-1));

  if (isVagueAction(action)) {
    const executed = execAction(action as VagueAction);
    if (executed) return;
  }
  
  if (isAnnotatedAction(action)) {
    const val = (action as AnnotatedAction).value;
    if (isVagueAction(val)) {
      const executed = execAction(val as VagueAction);
      if (executed) return;
    }
  }

  printAction(config, inputs);
}

function execAction(action: VagueAction): boolean {
  if (typeof action === "string") {
    exec(action, (err, stdout) => {
      if (err) console.error(err);
      console.log(stdout);
    });
    return true;
  }

  if (typeof action === "object" && typeof (action as any).length === "number") {
    const cmds = action as Array<string>;
    exec(cmds.join(" && "), (err, stdout) => {
      if (err) console.error(err);
      console.log(stdout);
    })
    return true;
  }

  return false;
}

main();
