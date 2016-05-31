app-restarter
=============

A CF cli plugin for restarting apps.

## Install

1. Visit the [Releases Page](https://github.com/cloudfoundry-incubator/app-restarter/releases) and copy a link to the plugin for your workstation's architecture.
2. Run `cf install-plugin <link-to-release-copied-above>`

## Usage

```bash
cf restart-apps
cf restart-apps -o org-name
cf restart-apps -s space-name
```

By default, the plugin will wait 60 seconds between restarting apps. Set `CF_STARTUP_TIMEOUT` in 
the shell to specify the number of seconds to wait.

```bash
CF_STARTUP_TIMEOUT=5 cf restart-apps
```
