# Bitpack CLI for VS Code

This extension runs Bitpack commands in the first open workspace folder and
displays output in the **Bitpack** output channel.

Install Bitpack separately and ensure its executable is on `PATH`. If needed,
set `bitpack.executable` to the executable's name or absolute path. Run
**Bitpack: Run Command** from the Command Palette, then enter a command such as
`modules`, `forge .\Bitfile`, or `assemble app --dry-run`.

Arguments are passed directly to the Bitpack process without a shell. The
extension does not bundle Bitpack or execute arbitrary shell commands.
