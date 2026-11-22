# Bitpack

Bitpack is a language-neutral project packaging, development, and deployment
tool. It can manage web apps, APIs, backend services, CLI tools, mobile
projects, and multi-language repositories. Bitpack itself is implemented
primarily in Go, with Python and shell also used internally; those are not
limits on the projects it can manage.

Bitpack runs the commands and toolchains you configure. It does not automatically
install project runtimes or build/push container images.

## Get started

Install Bitpack once for your Windows user:

```powershell
.\install.ps1
```

The installer builds Bitpack and adds its Go binary directory to your user
`PATH`. Restart your terminal host (restart VS Code too if you use its
integrated terminal), then use the command directly:

```powershell
bitpack help
bitpack forge .\Bitfile
bitpack assemble my-app
bitpack verify my-app
bitpack ignite my-app
```

Alternatively, run from source with `go run .\cmd\bitpack help`.

Commands declared by a Bitfile run through the operating-system shell with the
current user's permissions. Only run project files and commands you trust.
Required language toolchains must already be installed and available on
`PATH`.

## Project and workflow files

### `Bitfile`

Describes a project workspace, environment, secret names, and its assemble,
verify, and ignite commands. `using` is descriptive rather than a language
allowlist. For example:

```text
blueprint "my-app" v1
using application
bring "./"

env <|
    APP_ENV = production
|>

assemble <|
    cargo build --release
|>

verify <|
    cargo test
|>

ignite <|
    ./target/release/my-app
|>
```

Replace those commands with the tools your project uses, such as `npm`, `go`,
`dotnet`, `flutter`, or scripts in a multi-language repository. `forge` saves
the parsed project information as a local blueprint. `assemble` and `verify`
run their configured command blocks. `ignite` runs its configured command
block in the foreground; press Ctrl+C to stop it.

Secrets are declared by name, not value. Bitpack resolves each declared secret
from an environment variable named `BITPACK_SECRET_<NAME>`, with the name
uppercased and non-alphanumeric characters replaced by underscores. For
example, `db-password` is read from `BITPACK_SECRET_DB_PASSWORD`. Values are
passed to the configured process and are not intentionally printed.

### `.bpack`

Describes a container assembly: modules, image references, replicas, ports,
volumes, dependencies, firewall rules, autoscaling, and policies. Apply an
assembly to the Kubernetes cluster selected by the current kubeconfig context:

```powershell
bitpack assemble .\examples\bitpack.bpack --dry-run
bitpack assemble .\examples\bitpack.bpack --namespace production
bitpack assemblies
bitpack health social-network
bitpack stop social-network
bitpack restart social-network
bitpack destroy social-network
```

Kubernetes deployment requires `kubectl` and a configured, reachable cluster.
Bitpack applies resources to that existing cluster; it does not create an
AKS/GKE/EKS cluster, build images, or push them to a registry. Each module's
`blueprint` (or `image`) value is used as a container image reference, which
must already be accessible to Kubernetes.

The compiler creates Deployments and Services, PersistentVolumeClaims,
NetworkPolicies, HorizontalPodAutoscalers, and Kubernetes Secrets as needed.
Volume claims default to `1Gi` when no `size` is supplied. Modules with an
autoscale rule receive a `100m` CPU request; autoscaling also requires a working
cluster metrics API. NetworkPolicy enforcement depends on the cluster's network
plugin. `depends` is recorded as module metadata but does not wait for another
module to become ready.

For each secret name used by a module, Bitpack resolves the corresponding
`BITPACK_SECRET_<NAME>` environment variable and creates a Kubernetes Secret
reference for the workload. Secret values are sent to the Kubernetes API;
encryption at rest and access control depend on cluster configuration.

### `.bpdev`

Defines a local development workspace, environment, secrets, startup commands,
and shell commands. Run or inspect it with:

```powershell
bitpack dev inspect .\examples\project.bpdev
bitpack dev run .\examples\project.bpdev
```

Commands execute locally through the system shell. Secret environment variables
use the same `BITPACK_SECRET_<NAME>` convention. Declared services, ports, and
watch paths are currently descriptive only: the local runner does not provision
services or monitor files.

### `.bpbatch`

Defines a multi-stage workflow for forge, verify, publish, assemble, and ignite
targets. Stages run in order. With `parallel: true`, targets within each stage
run concurrently; output is collected in target order. `failfast: true` stops
later work after a failure, but cannot cancel targets already running in a
parallel stage.

```powershell
bitpack batch .\examples\release.bpbatch
bitpack batch-status
```

### `.bpackignore` and `.bpkg`

`pack` creates a gzip-compressed tar source package using the Bitfile's project
workspace. `.bpackignore` patterns exclude files from the package; defaults
exclude Bitpack state, VCS metadata, `node_modules`, and common environment
secret files. Project-specific patterns can exclude build output, caches, or
other generated files.

```powershell
bitpack pack .\Bitfile .\my-app.bpkg
bitpack unpack .\my-app.bpkg .\my-app
```

`unpack` refuses to overwrite an existing destination and rejects unsafe archive
paths and unsupported archive entry types. A `.bpkg` is a source archive, not a
container image.

## Current process and infrastructure scope

Kubernetes Deployments provide restart and replica supervision for `.bpack`
assemblies. `stop`, `restart`, `scale`, `logs`, `health`, `rollout`, `rollback`,
and `destroy` operate on recorded Kubernetes assemblies and their resources.
`ignite` for a Bitfile remains a foreground local command; local background
process supervision is not implemented.

The `cluster create` command records local cluster metadata; it does not create
cloud infrastructure. Several monitoring and registry commands are also local
metadata workflows rather than external provider integrations. Check command
output for the active provider and operation before relying on it.

## Develop Bitpack

Run the Go test suite and static checks from the repository root:

```powershell
go test ./...
go vet ./...
```
