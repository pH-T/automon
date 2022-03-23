# automon

Sysmon installation wrapper - quickly (un)install Sysmon with specific config.

## Hardcoded Configs
```
Available Sysmon configs:
[0] https://raw.githubusercontent.com/NextronSystems/aurora-helpers/master/sysmon-config/aurora-sysmon-config.xml
[1] https://raw.githubusercontent.com/Neo23x0/sysmon-config/master/sysmonconfig-export.xml
[2] https://raw.githubusercontent.com/SwiftOnSecurity/sysmon-config/master/sysmonconfig-export.xml
```

### Example:

* `automon --listconfigs` --> lists all known Sysmon config URLs (use `configURL` for other URL or create a issue)
* `automon --config 0` --> fresh Sysmon installation with config [0]
* `automon --config 0 --force` --> uninstalls old Sysmon and installs new Sysmon with config [0]
* `automon --sysmondownload` --> downloads and unzips Sysmon

## Usage

Downloaded files are written to current working directory!

```
  -arch string
        Which Sysmon version to use: 64 or 32 (default "64")
  -config int
        Which config should be used (default -1)
  -configURL string
        URL to download config
  -force
        Uninstalls Sysmon before installing
  -listconfigs
        Lists hardcoded Sysmon config URLs
  -sysmonURL string
        URL to download Sysmon zip (default "https://download.sysinternals.com/files/Sysmon.zip")
  -sysmondownload
        Just downloads Sysmon
  -uninstall
        Uninstall Sysmon
```

## Build

`make win`