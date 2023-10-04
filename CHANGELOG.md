## v0.1.1 [2023-10-04]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#19](https://github.com/turbot/steampipe-plugin-exec/pull/19))

## v0.1.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#17](https://github.com/turbot/steampipe-plugin-exec/pull/17))

## v0.0.4 [2023-09-29]

_Breaking changes_

- Removed the `output` column in the `exec_command` table. This column has been replaced by the `stdout_output` and `stderr_output` columns. ([#13](https://github.com/turbot/steampipe-plugin-exec/pull/13))

_What's new?_

- Added `stdout_output` and `stderr_output` columns to the `exec_command` table. ([#13](https://github.com/turbot/steampipe-plugin-exec/pull/13))
- Added `stream` column to the `exec_command_line` table. ([#13](https://github.com/turbot/steampipe-plugin-exec/pull/13))
- Added plugin limiter `exec_global` with `MaxConcurrency` set to 15 in an effort to reduce abuse reports due to large number of concurrent remote connections. ([#13](https://github.com/turbot/steampipe-plugin-exec/pull/13))

_Bug fixes_

- Results from the `exec_command` table should now be consistent when using local and remote connections. ([#13](https://github.com/turbot/steampipe-plugin-exec/pull/13))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v560-2023-09-27) which adds support for rate limiters. ([#13](https://github.com/turbot/steampipe-plugin-exec/pull/13))
- Recompiled plugin with Go 1.21. ([#13](https://github.com/turbot/steampipe-plugin-exec/pull/13))

## v0.0.3 [2023-08-11]

_Bug fixes_

- Fixed the plugin to return an appropriate error message when the config file is missing required arguments. ([#7](https://github.com/turbot/steampipe-plugin-exec/pull/7))

## v0.0.2 [2023-07-20]

_Bug fixes_

- Fixed the incorrect github repository reference in `docs/index.md` file.

## v0.0.1 [2023-07-13]

_What's new?_

- New tables added
  - [exec_command](https://hub.steampipe.io/plugins/turbot/exec/tables/exec_command)
  - [exec_command_line](https://hub.steampipe.io/plugins/turbot/exec/tables/exec_command_line)
