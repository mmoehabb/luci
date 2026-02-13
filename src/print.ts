import chalk from "chalk";
import type { ActionRecord, AnnotatedAction, Config } from "./types";
import { dig, getShellType, isAction, isActionRecord, isAnnotatedAction } from "./util";

const log = console.log;

export function printConfigHeader(config: Config) {
    if (config.title) log(chalk.green(`*** ${config.title} ***`));
    if (config.description) log(chalk.gray(`> ${config.description}`));
    log(chalk.yellow("\nUsage:\n"));
}

export function printUsage(config: Config) {
  const shell = getShellType();
  printConfigHeader(config);
  for (const action in config[shell]) {
    printAction(config, [action]);
  }
}

export function printAction(config: Config, inputs: Array<string>, level: number = 0) {
  const action = dig(config, inputs);
  if (!action) return printUsage(config);

  if (isAnnotatedAction(action)) {
    const annAction = action as AnnotatedAction;
    log(chalk.bgBlack(`${indent(level)}luci ${inputs.join(' ')}`));
    if (annAction.title) log(chalk.blue(`${indent(level+1)}** ${annAction.title} **`));
    if (annAction.description) log(chalk.gray(`${indent(level+1)}> ${annAction.description}`));
    if (isActionRecord(annAction.value)) {
      for (const action in annAction.value as ActionRecord) {
        printAction(config, [...inputs, action], level + 1);
      }
    }
    return;
  }

  if (isAction(action)) {
    return log(`${indent(level)}luci ${inputs.join(' ')}`);
  }

  for (const key in action as ActionRecord) {
    printAction(config, [...inputs, key]);
  }
}

function indent(count: number) {
  let res = '';
  for (let i = 0; i < count; i++) res += '\t';
  return res;
}
