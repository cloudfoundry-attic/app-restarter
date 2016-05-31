app-restarter
=============

A CF cli plugin for restarting apps.

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
