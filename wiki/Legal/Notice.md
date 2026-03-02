# Notice

## Upstream Attribution

SoloDev is a derivative work built from [Gitness by Harness](https://github.com/harness/gitness), an open-source DevOps platform licensed under the Apache License, Version 2.0.

The following components are inherited from the upstream Gitness project:
- Git hosting and source code management
- CI/CD pipeline execution (Drone-based)
- Gitspaces (cloud developer environments)
- Artifact registry (OCI-compliant)
- Authentication and authorization framework
- Web UI application shell
- CLI framework
- Database migration system

## Attribution

```
Gitness
Copyright 2023 Harness, Inc.

Licensed under the Apache License, Version 2.0
```

## SoloDev Additions

All SoloDev-specific additions (AI Layer, Signal System, Remediation System, Developer Experience modules, Agent System, and associated documentation) are original work, also licensed under Apache License 2.0.

## Derivative Work Notice

This project is presented as SoloDev, a fork that extends Gitness in an AI-native, solo-builder-focused direction. It is not presented as an unrelated clean-room rewrite. The provenance is intentional and explicit:

- Upstream attribution is preserved
- Apache-2.0 licensing remains in place
- Derivative-work notices are maintained
- The [NOTICE](https://github.com/EolaFam1828/SoloDev/blob/main/NOTICE) file in the repository contains full attribution

## Third-Party Dependencies

SoloDev's dependencies and their licenses are tracked in `go.mod` (Go dependencies) and `web/package.json` (JavaScript/TypeScript dependencies). All dependencies are compatible with Apache-2.0 licensing.
