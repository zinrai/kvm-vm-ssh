# KVM SSH Connection and Port Forwarding Tool

`kvm-ssh` is a command-line tool designed to simplify SSH connections to KVM (Kernel-based Virtual Machine) instances and provide easy port forwarding capabilities.

## Features

- Easily connect to a specific KVM virtual machine via SSH
- Local and remote port forwarding with simple port number specification
- Pass arbitrary SSH options via `-ssh-opts`
- Sensible defaults for common options

## Notes

- Ensure that the `ssh` command is available on your system.
- The tool reads VM information from `/var/lib/libvirt/dnsmasq/<bridge_name>.status`.

## Installation

To install `kvm-ssh`, clone the repository and build the tool:

```bash
$ go build
```

Make sure to place the built `kvm-ssh` binary in your system's PATH for easy access.

## Usage

```
kvm-ssh [options] <vm_name>
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-bridge` | `virbr0` | Bridge name |
| `-user` | `$USER` | SSH user |
| `-local` | | Local forward ports (comma-separated) |
| `-remote` | | Remote forward ports (comma-separated) |
| `-ssh-opts` | | Additional SSH options |

## Examples

Connect to a VM:

```bash
$ kvm-ssh ubuntu-vm
```

Connect with a specific user and bridge:

```bash
$ kvm-ssh -user john -bridge br0 debian-vm
```

Local port forwarding:

```bash
$ kvm-ssh -local 2375,40413 bookworm64-docker
```

Remote port forwarding:

```bash
$ kvm-ssh -remote 3000 dev-vm
```

Both local and remote forwarding with SSH options:

```bash
$ kvm-ssh -local 8080 -remote 3000 -ssh-opts "-o StrictHostKeyChecking=no" dev-vm
```

## License

This project is licensed under the [MIT License](./LICENSE).
