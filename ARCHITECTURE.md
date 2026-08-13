# Architecture

Fabric Network Orchestrator is a CLI product. Its supported interfaces are the
`fno` commands and the YAML, JSON, and TOML configuration formats; the Go
packages are implementation details and therefore live under `internal/`.

## Dependency direction

```text
cmd/cli -> internal/cli -> internal/application -> config and infrastructure
cmd/experiment-runner -> internal/experiment
```

- `cmd/` contains executable entry points only.
- `internal/cli` owns Cobra command and flag wiring. Command handlers load a
  configuration, invoke one application workflow, and report its result.
- `internal/application` coordinates artifact generation, identities, images,
  network lifecycle, and chaincode deployment. It adds operation context to
  returned errors and does not terminate the process.
- Focused infrastructure packages such as `compose`, `configtx`, `network`,
  `chaincode`, `executor`, `file`, and `yaml` perform rendering and external
  operations. They must not import `cli` or `application`.
- `config`, `spec`, and `validate` own the configuration model, defaults,
  decoding, and validation rules.

External commands are accessed through `internal/executor.Executor`. Production
workflows use `DefaultExecutor`; tests inject fakes so Docker and Fabric binaries
are not required.

The dependency-direction test in `internal/architecture` prevents lower-level
packages from importing the CLI or application layer.
