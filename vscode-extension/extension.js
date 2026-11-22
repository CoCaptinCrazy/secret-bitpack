const { spawn } = require("node:child_process");
const vscode = require("vscode");
const { parseCommandLine } = require("./args");

function activate(context) {
  const output = vscode.window.createOutputChannel("Bitpack");
  context.subscriptions.push(output);
  context.subscriptions.push(
    vscode.commands.registerCommand("bitpack.runCommand", async () => {
      const workspace = vscode.workspace.workspaceFolders?.[0];
      if (!workspace) {
        vscode.window.showErrorMessage("Open a workspace folder before running Bitpack.");
        return;
      }

      const input = await vscode.window.showInputBox({
        prompt: "Enter a Bitpack command and arguments",
        value: "help",
        placeHolder: 'For example: forge .\\Bitfile or assemble app --dry-run',
      });
      if (input === undefined) {
        return;
      }

      let args;
      try {
        args = parseCommandLine(input);
      } catch (error) {
        vscode.window.showErrorMessage(error.message);
        return;
      }
      if (args.length === 0) {
        vscode.window.showErrorMessage("Enter a Bitpack command.");
        return;
      }

      const executable = vscode.workspace
        .getConfiguration("bitpack")
        .get("executable", "bitpack");
      output.show(true);
      output.appendLine(`> ${executable} ${input}`);

      const process = spawn(executable, args, {
        cwd: workspace.uri.fsPath,
        shell: false,
        windowsHide: true,
      });
      process.stdout.on("data", (chunk) => output.append(chunk.toString()));
      process.stderr.on("data", (chunk) => output.append(chunk.toString()));
      process.on("error", (error) => {
        output.appendLine(`Could not start Bitpack: ${error.message}`);
      });
      process.on("close", (code, signal) => {
        output.appendLine(
          code === 0
            ? "Bitpack finished successfully."
            : `Bitpack exited with code ${code}${signal ? ` (signal ${signal})` : ""}.`,
        );
      });
    }),
  );
}

function deactivate() {}

module.exports = { activate, deactivate };
