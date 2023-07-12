![image](https://hub.steampipe.io/images/plugins/turbot/exec-social-graphic.png)

# Exec Plugin for Steampipe

Use SQL to run commands locally or on remote Linux and Windows hosts. Then get the result as a table.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/exec)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/exec/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
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

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-exec/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Exec Plugin](https://github.com/turbot/steampipe-plugin-exec/labels/help%20wanted)
