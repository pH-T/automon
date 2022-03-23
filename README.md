# autosysmon

## Wrapper for Sysmon

Quickly install Sysmon with specific config!

## Hardcoded Configs
```
Available Sysmon configs:
[0] https://raw.githubusercontent.com/NextronSystems/aurora-helpers/master/sysmon-config/aurora-sysmon-config.xml
[1] https://raw.githubusercontent.com/Neo23x0/sysmon-config/master/sysmonconfig-export.xml
[2] https://raw.githubusercontent.com/SwiftOnSecurity/sysmon-config/master/sysmonconfig-export.xml
```

## Usage

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
        Lists known Sysmon configs
  -sysmonURL string
        URL to download Sysmon (default "https://download.sysinternals.com/files/Sysmon.zip")
  -sysmondownload
        Just downloads Sysmon
  -uninstall
        Uninstall Sysmon

```