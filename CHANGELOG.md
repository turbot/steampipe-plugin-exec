## v0.0.4 [TDB]

_Bug fixes_

- Added support for configuring the maximum concurrency limitation, which defaults to 15.
- Removed the column `output` of table `exec_command` in favor of `stdout_output` and `stderr_output` columns.
- Added the columns `stdout_output` and `stderr_output` to table `exec_command` to better reflect the streams of command output.

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
