![image](https://hub.steampipe.io/images/plugins/turbot/exec-social-graphic.png)

# Exec Plugin for Steampipe

Use SQL to run commands locally or on remote Linux and Windows hosts. Then get the result as a table.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/exec)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/exec/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-exec/issues)

## Quick start

### Install

Download and install the latest Steampipe plugin:

```bash
steampipe plugin install exec
```

Configure your [credentials](https://hub.steampipe.io/plugins/turbot/exec#credentials) and [config file](https://hub.steampipe.io/plugins/turbot/exec#configuration).

Configure your credential details in `~/.steampipe/config/exec.spc`:

```hcl
connection "exec" {
  plugin = "exec"

  protocol = "ssh"
  host = "my-remote-linux-host"
  user = "my-username"
  private_key = "~/.ssh/my-remote-linux-host.pem"
}
```

Run steampipe:

```shell
steampipe query
```

List disks on a Linux host

```sql
select
  output
from
  exec_command 
where
  command = 'df -h';
```

```
+------------------------------------------------------+
| output                                               |
+------------------------------------------------------+
| Filesystem      Size  Used Avail Use% Mounted on     |
| /dev/root       7.6G  3.4G  4.3G  44% /              |
| tmpfs           483M     0  483M   0% /dev/shm       |
| tmpfs           194M  872K  193M   1% /run           |
| tmpfs           5.0M     0  5.0M   0% /run/lock      |
| /dev/xvda15     105M  5.3M  100M   5% /boot/efi      |
| tmpfs            97M  4.0K   97M   1% /run/user/1001 |
+------------------------------------------------------+
```

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/index) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs//steampipe_sqlite/index) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/index) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-exec.git
cd steampipe-plugin-exec
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/exec.spc
```

Try it!

```
steampipe query
> .inspect exec
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Steampipe](https://steampipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #steampipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Exec Plugin](https://github.com/turbot/steampipe-plugin-exec/labels/help%20wanted)
