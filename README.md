# K3SAIR 🏴‍☠️️ ('Corsair')

`k3sair` is a cli for the installation of k3s in an Air-Gapped environment.

The idea is born, during the installation attempt in my company. So we are using this cli too, for our own installtions.

It is inspired by [k3sup](https://github.com/alexellis/k3sup), which does a great work.

### Known Limitation 😵

`k3sair` is still under develoment and supports at the moment only amd64 architecture and no version selection. It is
always the binary you provide.

And there is no HA Setup. The `install` command is for a single control plane server.

### Prerequisite 📚

You should have access to a http server hosting the files from [k3s](https://github.com/k3s-io/k3s) release page.

- k3s
- k3s-airgap-images-`<arch>`.tar.gz (See Known Limitation)

### Usage ⚙️

#### Install 💾

```bash
k3sair install \
--ssh-key /ssh/cluster \
--arch amd64 \
--base "https://repo.local/" \
--ip 127.0.0.1 \
--user core
```

#### Join 🚪

```bash
k3sair join \
--ssh-key /ssh/cluster \
--arch amd64 \
--base "https://repo.local/" \
--ip 193.148.164.208 \
--control-plane-ip 193.148.165.6 \
--user core
```

### Contributing 🤝

#### Contributing via GitHub

TBA

#### License

Apache License, Version 2.0

### Roadmap 🛣️

- GitHub Actions
- Release via goreleaser
- HA Support
- ...