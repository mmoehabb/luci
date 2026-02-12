import * as os from "os"
import * as fs from "fs"
import { argv } from "process";
import chalk from "chalk";
import { exec } from "child_process";

type CmdObj = {
  title: string;
  description: string;
  cmd: Command;
};

type Command = string | Array<string> | CmdObj | {
  [key: string]: Command;
};

type ShellType = "bash" | "zshell" | "bat";

type Shell = Record<string, Command>;

type Config = {
  title: string;
  description: string;
  bash: Shell;
  zshell: Shell;
  bat: Shell;
};

function getShellType(): ShellType {
  const platform = os.platform();
  if (platform === "linux") return "bash";
  if (platform === "darwin") return "zshell";
  if (platform === "win32") return "bat";
  throw Error("Shell couldn't be recognized!");
}

function printConfigHeader(config: Config) {
    if (config.title) console.log(chalk.green(`*** ${config.title} ***`));
    if (config.description) console.log(chalk.gray(`> ${config.description}`));

    console.log(chalk.yellow("\nUsage:\n"));
}

function printUsage(config: Config) {
  const shell = getShellType();
  printConfigHeader(config);
  for (const cmd in config[shell]) printCommand(config, [cmd]);
}

function printCommand(config: Config, inputs: Array<string>, level: number = 0) {
  const command = getCommand(config, inputs);
  if (!command) return printUsage(config);

  if (typeof command === "string") {
    return console.log(`${indent(level)}luci ${inputs.join(' ')}`);
  }

  if (typeof command === "object" && (command as any).length) {
    return console.log(`${indent(level)}luci ${inputs.join(' ')}`);
  }

  // NOTE: the json file assumed to not include custom command named "cmd".
  if (typeof command === "object" && (command as CmdObj).cmd) {
    const cmdobj = command as CmdObj;
    console.log(chalk.bgBlack(`${indent(level)}luci ${inputs.join(' ')}`));
    if (cmdobj.title) console.log(chalk.blue(`${indent(level+1)}** ${cmdobj.title} **`));
    if (cmdobj.description) console.log(chalk.gray(`${indent(level+1)}> ${cmdobj.description}`));
    if (typeof cmdobj.cmd === "object" && !(cmdobj.cmd as any).length) {
      for (const innerCommand in cmdobj.cmd) {
        printCommand(config, [...inputs, innerCommand], level + 1);
      }
    }
    return;
  }

  for (const cmd in command) {
    printCommand(config, [...inputs, cmd]);
  }
}

function indent(count: number) {
  let res = '';
  for (let i = 0; i < count; i++) res += '\t';
  return res;
}

function getCommand(config: Config, inputs: Array<string>): Command | null {
  const shell = getShellType();
  let command: any = config[shell];
  for (const input of inputs) {
    if (!command) return null;
    if (typeof command === "string") return command;
    if (typeof command === "object" && command.length) return command;
    if (typeof command === "object" && typeof command.cmd === "string") return command;
    if (typeof command === "object" && typeof command.cmd === "object" && command.cmd.length) return command;
    if (typeof command === "object" && typeof command.cmd === "object") command = command.cmd[input];
    else command = command[input];
  }
  return command || null;
}

function execCommand(config: Config, inputs: Array<string>) {
  const command = getCommand(config, inputs);
  if (!command) return printCommand(config, inputs.slice(0, inputs.length-1));

  if (typeof command === "string") {
    exec(command, (err, stdout) => {
      if (err) console.error(err);
      console.log(stdout);
    });
    return;
  }

  if (typeof command === "object" && (command as any).length) {
    const cmds = command as Array<string>;
    exec(cmds.join(" && "), (err, stdout) => {
      if (err) console.error(err);
      console.log(stdout);
    })
    return;
  }
  
  if (typeof command === "object" && (command as CmdObj).cmd) {
    const cmd = (command as CmdObj).cmd;

    if (typeof cmd === "string") {
      exec(cmd, (err, stdout) => {
        if (err) console.error(err);
        console.log(stdout);
      });
      return;
    }

    if (typeof cmd === "object" && (cmd as any).length) {
      const cmds = cmd as Array<string>;
      exec(cmds.join(" && "), (err, stdout) => {
        if (err) console.error(err);
        console.log(stdout);
      })
      return;
    }
  }

  printCommand(config, inputs);
}

function main() {
  const config: Config = JSON.parse(fs.readFileSync("luci.config.json", "utf8"));
  if (argv.length <= 2) printUsage(config);
  else execCommand(config, argv.slice(2));
}

main();
