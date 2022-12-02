# golang-kubectx-tray

Small Menu Bar App to quickly identify the current Kubernetes Context:

![k8s icon](docs/mac_tray.png "k8s icon")

Clicking on it we will show the full name of the current context (first line) and the current namespace (second line):

![k8s icon](docs/mac_tray_context_namespace.png "k8s icon")

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Configuration](#api-implementation)
- [Installation](#installation)
    - [Build and install on Mac](#build-and-install-on-mac)
    - [Build and install on Linux](#build-and-install-on-linux)
    - [Build and install on Windows](#build-and-install-on-windows)

## Configuration

Using the **~/.kube/kct-config** configuration file we can choose an icon and some small text to be shown on the tray:

```
contexts:
- match: preprod.k8s
  title: preprod
  icon: yellow
- match: devel.k8s
  title: dev
  icon: green
- match: kind-
  title: kind
```

The current list of available icons is:

* **kube**: ![k8s icon](icons/kube.png "k8s icon")
* **green**: ![Green k8s icon](icons/green.png "Green k8s icon")
* **yellow**: ![Yellow k8s icon](icons/yellow.png "Yellow k8s icon")
* **red**: ![Red k8s icon](icons/red.png "Red k8s icon")
* **loki**: ![Loki icon](icons/loki.png "Loki icon")
* **odin**: ![Odin icon](icons/odin.png "Odin icon")
* **greenproc**: ![Green processor icon](icons/proc_green.png "Green processor icon")
* **yellowproc**: ![Yellow processor icon](icons/proc_yellow.png "Yellow processor icon")
* **redproc**: ![Red processor icon](icons/proc_red.png "Red processor icon")

Please submit a PR if you want to add more, in the **icons** folder you can find the icon plus the [2goarray](https://github.com/cratonica/2goarray) version of it.

## Installation

It's currently only tested on **MacOS**.

### Build and install on Mac

To build this application you can simply run `make`. It will create the binary file under **KubeCtxTray.app** so you can just drag it into your applications folder.

### Build and install on Linux

It's not tested, but it should work file. Please follow [getlantern/systray](https://github.com/getlantern/systray) instructions for installing it's dependencies.

### Build and install on Windows

Please submit a PR if you get it working on Windows.