# Agent Dev Environment

This project contains a sandboxed environment. Currently implemented through a fat docker image using ubuntu and a few utilities, and the go code that serves an http api.

The http api can then be consumed by orchestration workflows to enable fully autonomous and git-aware coding workflows.

The environment expects to have all configuration available in the cloned repos (i.e. secrets are encrypted and can be decrypted by pulumi config passphrase).

# Rutime view

This is deployed through pulumi and github workflows. The main non-obvious bits are:
 - Headless service deployment. This is done so the agent orchestration can spin up multiple pods of this fat docker image and still access it through the network.
 - Secrets. There are a few important secrets, you can check these in the pulumi templates. Github token is used to interact with github repositories and pulumi config passphrase so the agent can decrypt secrets in these repositories.

# Development
## Mise

[mise](https://mise.jdx.dev/) is used to manage tool versions and abstract common tasks. It is installed in the Docker image and available at runtime.

## Prerequisites

- [mise](https://mise.jdx.dev/) installed.

## Running locally

On your first run:
```bash
# Install tools
mise install

# Configure project
mise setup-project
```

The project is aimed to be run through e2e tests. So you can use either e2e command in the mise.toml (running the e2e with or without docker).
