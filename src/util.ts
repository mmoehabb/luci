import os from "node:os";
import type { Action, ActionRecord, Config, ShellType } from "./types";

export function isActionRecord(obj: any): boolean {
  return typeof obj === "object" && typeof obj.length !== "number" && !Object.values(obj).map(v => isAction(v)).includes(false);
}

export function isAction(obj: any): boolean {
  return isVagueAction(obj) || isAnnotatedAction(obj);
}

export function isVagueAction(obj: any): boolean {
  return (
    typeof obj === "string" ||
    (typeof obj === "object" && obj.length && typeof obj[0] === "string")
  );
}

export function isAnnotatedAction(obj: any): boolean {
  return typeof obj === "object" && obj.value;
}

export function getShellType(): ShellType {
  const platform = os.platform();
  if (platform === "linux") return "bash";
  if (platform === "darwin") return "zshell";
  if (platform === "win32") return "bat";
  throw Error("Shell couldn't be recognized!");
}

export function dig(config: Config, inputs: Array<string>): Action | ActionRecord | null {
  const shell = getShellType();
  let obj: any = config[shell];
  for (const input of inputs) {
    if (!obj) return null;
    else if (isAnnotatedAction(obj) && isActionRecord(obj.value)) obj = obj.value[input];
    else if (isAction(obj)) return obj;
    else obj = obj[input];
  }
  return obj || null;
}
