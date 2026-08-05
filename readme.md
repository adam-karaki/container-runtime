# myrun

A lightweight container runtime written in Go using Linux namespaces and cgroups.

## Features

* PID, UTS, and mount namespace isolation
* Custom root filesystem with `pivot_root`
* `/proc` mounting
* CPU and memory limits with cgroups v2
* Configurable hostname

## Usage

```bash
./myrun run \
  --hostname test \
  --rootfs /path/to/rootfs \
  --memory 256M \
  --cpu 0.5 \
  /bin/sh
```

Linux and cgroups v2 are required.

## How it works

`myrun` re-executes itself in a new set of Linux namespaces, sets up the root filesystem, mounts `/proc`, and runs the requested command. Resource limits are applied through cgroups v2.
