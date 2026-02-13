export type Config = {
  title: string;
  description: string;
  bash:  Record<ShellType, ActionRecord>;
  zshell: Record<ShellType, ActionRecord>;
  bat: Record<ShellType, ActionRecord>;
};

export type ShellType = "bash" | "zshell" | "bat";

export type ActionRecord = Record<string, Action>;

export type Action = AnnotatedAction | VagueAction;

export type VagueAction = string | Array<string>;

export type AnnotatedAction = {
  title: string;
  description: string;
  value: VagueAction | ActionRecord;
}
