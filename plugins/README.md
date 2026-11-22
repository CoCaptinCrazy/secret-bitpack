# Python extensions

Python is an optional extension runtime. Go remains the application entry point
and only runs an extension when the user explicitly invokes:

```powershell
go run .\cmd\bp plugin run <name> [arguments...]
```

Create `plugins/<name>.py` for each extension. The script receives the
arguments after its name and inherits the terminal's standard input and output.
Go sets `BITPACK_HOST=go` and `BITPACK_EXTENSION=<name>`.

Extensions are executable code. Only run plugins you trust. Extension names
must contain only letters, digits, underscores, and hyphens; scripts must be
regular files located inside this directory.
