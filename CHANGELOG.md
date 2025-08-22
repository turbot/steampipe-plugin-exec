## v1.2.0 [2025-08-22]

_Bug fixes_

- Fixed the example query in the plugin documentation to use the correct column name `exec_output` instead of `output`. ([#63](https://github.com/turbot/steampipe-plugin-exec/pull/63)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!!)

_Dependencies_

- Recompiled plugin with Go version `1.24`.
- Recompiled plugin with [steampipe-plugin-sdk v5.13.0](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5130-2025-07-21) that addresses critical and high vulnerabilities in dependent packages.

## v1.1.1 [2025-04-18]

_Bug fixes_

- Fixed Linux AMD64 plugin build failures for `Postgres 14 FDW`, `Postgres 15 FDW`, and `SQLite Extension` by upgrading GitHub Actions runners from `ubuntu-20.04` to `ubuntu-22.04`.

## v1.1.0 [2025-04-17]

_Dependencies_

- Recompiled plugin with Go version `1.23.1`. ([#51](https://github.com/turbot/steampipe-plugin-exec/pull/51))
- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.11.5/CHANGELOG.md#v5115-2025-03-31) that addresses critical and high vulnerabilities in dependent packages. ([#51](https://github.com/turbot/steampipe-plugin-exec/pull/51))

## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#40](https://github.com/turbot/steampipe-plugin-exec/pull/40))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#40](https://github.com/turbot/steampipe-plugin-exec/pull/40))

## v0.2.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#27](https://github.com/turbot/steampipe-plugin-exec/pull/27))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#27](https://github.com/turbot/steampipe-plugin-exec/pull/27))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-exec/blob/main/docs/LICENSE). ([#27](https://github.com/turbot/steampipe-plugin-exec/pull/27))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#26](https://github.com/turbot/steampipe-plugin-exec/pull/26))

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
