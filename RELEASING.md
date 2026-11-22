# Releasing Bitpack

Pushing a version tag such as `v0.1.17` runs
`.github/workflows/release.yml`. The workflow tests the Go CLI, creates
Windows, macOS, and Linux archives for amd64 and arm64, publishes a GitHub
Release with checksums, and attaches the VS Code `.vsix` package.

The release also prepares a winget manifest for
`BitpackDevelopment.Bitpack`. To submit it automatically:

1. Fork `microsoft/winget-pkgs` under the `CoCaptinCrazy` account.
2. Add a `WINGET_GITHUB_TOKEN` Actions secret containing a classic GitHub
   personal access token with permission to push to that fork and open pull
   requests. The release workflow skips winget submission when this secret is
   absent.
3. Merge the generated pull request in `microsoft/winget-pkgs` before
   recommending `winget install BitpackDevelopment.Bitpack`.

To publish the extension as `bitpackdevelopment.bitpack-tools`:

1. Register the `bitpackdevelopment` publisher in the Visual Studio
   Marketplace.
2. Configure Marketplace Trusted Publishing for GitHub owner
   `CoCaptinCrazy`, repository `secret-bitpack`, workflow `release.yml`, with
   no GitHub environment.
3. Set the repository Actions variable
   `VSCODE_MARKETPLACE_PUBLISHING_ENABLED` to `true`.

Without Marketplace publishing configured, the extension is still packaged
and attached to each GitHub Release for manual VSIX installation. It runs the
separately installed `bitpack` executable and can be configured with the
`bitpack.executable` VS Code setting.

For contributors who have cloned the source, `.\install.ps1` installs the
current checkout into the Go binary directory and adds it to the user `PATH`.
