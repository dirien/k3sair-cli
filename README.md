# K3SAIR 🏴‍☠️️ ('Corsair')

`k3sair` is a cli for the installation of k3s in an Air-Gapped environment.

The idea is born, during the installation attempt in my company. So we are using this cli too, for our own
installations. It is build completely on zero-trust, `k3sair` is not saving anything.

It is inspired by [k3sup](https://github.com/alexellis/k3sup), which does a great work.

### TL;DR 🚀

Install via homebrew:

```bash
brew tap dirien/homebrew-dirien
brew install k3sair
```

Linux or Windows user, can directly download (or use `curl`/`wget`) the binary via
the [release page](https://github.com/dirien/k3sair-cli/releases).

### Known Limitation 😵

`k3sair` is still under development and supports at the moment only amd64 architecture and no version selection. It is
always the binary you provide.

And there is no HA Setup. The `install` command is for a single control plane server.

### Prerequisite 📚

You should have access to a http server hosting the files from [k3s](https://github.com/k3s-io/k3s) release page. We
use [Artifactory](https://jfrog.com/).

- k3s
- k3s-airgap-images-`<arch>`.tar.gz (See Known Limitation)

### Usage ⚙️

#### Install 💾

[![asciicast](https://asciinema.org/a/420950.svg)](https://asciinema.org/a/420950)

```bash
k3sair install -h

Usage:
  k3sair install [flags]

Flags:
      --arch string      Enter the target sever os architecture (amd64 supported atm)
      --base string      Enter the on site proxy repository url (e.g Artifactory)
  -h, --help             help for install
      --ip string        Public IP or FQDN of node
      --mirror string    Mirrored Registry. (Default: '')
      --ssh-key string   The ssh key to use for remote login
      --sudo              Use sudo for installation. (Default: true) (default true)
      --user string      Username for SSH login (Default: root (default "root")

Examples:
    k3sair install \
    --ssh-key /ssh/cluster \
    --arch amd64 \
    --base "https://repo.local/" \
    --ip 127.0.0.1 \
    --user core
```

#### Join 🚪

[![asciicast](https://asciinema.org/a/420961.svg)](https://asciinema.org/a/420961)

```bash
k3sair join -h   

Usage:
  k3sair join [flags]

Flags:
      --arch string               Enter the target sever os architecture (amd64 supported atm)
      --base string               Enter the on site proxy repository url (e.g Artifactory)
      --control-plane-ip string   Public IP or FQDN of an existing k3s server
  -h, --help                      help for join
      --ip string                 Public IP or FQDN of node
      --mirror string             Mirrored Registry. (Default: '')
      --ssh-key string            The ssh key to use for remote login
      --sudo                      Use sudo for installation. (Default: true) (default true)
      --tls-san string            Add additional hostname or IP as a Subject Alternative Name in the TLS cert
      --user string               Username for SSH login (Default: root (default "root")
      

Examples:
k3sair join \
--ssh-key /ssh/cluster \
--arch amd64 \
--base "https://repo.local/" \
--ip 127.0.0.2 \
--control-plane-ip 127.0.0.1 \
--user core
```

#### Kubeconfig

[![asciicast](https://asciinema.org/a/420965.svg)](https://asciinema.org/a/420965)

```bash
k3sair kubeconfig -h
Get the kubeconfig from the k3s control plane server

Usage:
  k3sair kubeconfig [flags]

Flags:
  -h, --help             help for kubeconfig
      --ip string        Public IP or FQDN of node
      --ssh-key string   The ssh key to use for remote login
      --sudo              Use sudo for installation. (Default: true) (default true)
      --user string      Username for SSH login (Default: root (default "root")

Examples:
k3sair kubeconfig \
--ssh-key ~/.ssh/id_rsa
--ip 127.0.0.1
```

### Contributing 🤝

#### Contributing via GitHub

TBA

#### License

Apache License, Version 2.0

### Roadmap 🛣️

- [x] K3s private registry support [link](https://rancher.com/docs/k3s/latest/en/installation/private-registry/)
- [x] tls-san support
- [x] Kubeconfig: change server address 127.0.0.1 to host ip/name
- [x] INSTALL_K3S_EXEC support
- [x] GitHub Actions
- [x] Release via goreleaser
- [x] Rework codebase
- [ ] HA Support
- [ ] Tests
- ...

### Libraries & Tools 🔥

- https://github.com/fatih/color
- https://github.com/melbahja/goph
- https://github.com/spf13/cobra
- https://github.com/goreleaser