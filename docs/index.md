---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/exec.svg"
brand_color: "#002342"
display_name: "Exec"
short_name: "exec"
description: "Steampipe plugin to run & query shell commands on local and remote servers."
og_description: "Run & query shell commands with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/exec-social-graphic.png"
---

# Exec + Steampipe

Execute commands locally or on remote Linux and Windows hosts through SSH or WinRM.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

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

## Documentation

- **[Table definitions & examples →](/plugins/turbot/exec/tables)**

## Get started

### Install

Download and install the latest Exec plugin:

```bash
steampipe plugin install exec
```

### Credentials

| Item        | Description                                                                                                                                                                                                                   |
|-------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Credentials | Local connection does not require credentials. For Linux remote connections, the plugin supports SSH private key authentication or password. For Windows remote connections, the plugin supports password authentication.     |
| Permissions | As the connection to host relies on SSH or WinRM, the user must have the necessary permissions to connect to the host.                                                                                                        |
| Radius      | Each connection represents a single Exec Installation.                                                                                                                                                                        |
| Resolution  | 1. With configuration provided in connection in steampipe _**.spc**_ config file.<br />2. An exec.yaml file in a .exec folder in the current user's home directory _**(~/.exec/exec.yaml or %userprofile\.exec\exec.yaml)**_. |

### Configuration

There are multiple ways to configure the plugin depending on the type of connection you want to use. Bellows are some examples of how to configure the most common connections.

#### Local connection

```hcl
connection "exec_local" {
  plugin = "exec"
}
```

#### Remote Linux connection through SSH

```hcl
connection "exec_remote_linux" {
  plugin = "exec"

  protocol = "ssh"
  host = "my-remote-linux-host"
  user = "my-username"
  private_key = "~/.ssh/my-remote-linux-host.pem"
}
```

#### Remote Windows connection through WinRM

```hcl
connection "exec_remote_windows" {
  plugin = "exec"

  protocol = "winrm"
  host = "18.228.214.45"
  port = 5985
  user = "Administrator"
  password = "rh=PM76t54nouv&dqwe3cNM7J1(*skZhh*"
}
```

#### Remote Windows connection through WinRM ignoring certificate validation

```hcl
connection "exec_remote_windows" {
  plugin = "exec"

  protocol = "winrm"
  host = "18.228.214.45"
  https = true
  port = 5986
  insecure = true
  user = "Administrator"
  password = "rh=PM76t54nouv&dqwe3cNM7J1(*skZhh*"
}
```

#### Remote Linux connection through SSH with in-file private key

```hcl
connection "exec_remote_linux" {
  plugin = "exec"

  protocol = "ssh"
  host = "my-remote-linux-host"
  user = "my-username"
  private_key = <<EOK
-----BEGIN RSA PRIVATE KEY-----
... snipped ...
-----END RSA PRIVATE KEY-----
EOK
}
```

#### Remote Linux connection through SSH with password

```hcl
connection "exec_remote_linux" {
  plugin = "exec"
  
  protocol = "ssh"
  host = "my-remote-linux-host"
  user = "my-username"
  password = "my-password"
}
```

#### Remote Linux connection through SSH via bastion host

```hcl
connection "production" {
  plugin = "exec"

  protocol = "ssh"
  host = "172.31.40.195"
  user = "ubuntu"
  private_key = "~/.ssh/my-remote-linux-host.pem"
  
  bastion_user = "ubuntu"
  bastion_host = "52.67.221.206"
  bastion_port = 22
  bastion_private_key = "~/.ssh/my-bastion-host.pem"
}
```

#### Remote Linux connection through SSH over proxy

```hcl
connection "staging" {
  plugin = "exec"
  
  protocol = "ssh"
  host = "52.67.221.206"
  user = "ubuntu"
  private_key = "~/.ssh/my-remote-linux-host.pem"

  proxy_host = "10.10.10.200"
  proxy_port = 8080
  proxy_user_name = "luis"
  proxy_user_password = "c@v41c@nt3"
}
```

#### Remote Linux connection through SSH via bastion host over proxy

```hcl
connection "production" {
  plugin = "exec"

  protocol = "ssh"
  host = "172.31.40.195"
  user = "ubuntu"
  private_key = "~/.ssh/my-remote-linux-host.pem"

  bastion_user = "ubuntu"
  bastion_host = "52.67.221.206"
  bastion_port = 22
  bastion_private_key = "~/.ssh/my-bastion-host.pem"

  proxy_host = "10.10.10.200"
  proxy_port = 8080
  proxy_user_name = "luis"
  proxy_user_password = "c@v41c@nt3"
}
```

#### Remote connection to multiple hosts

Here we are using **[aggregators](https://steampipe.io/docs/managing/connections#using-aggregators)** to run the same query on multiple hosts. The `staging_servers` connection is an aggregator connection that will run the same query on all connections that match the `*-staging` pattern. The `server1-staging` and `server2-staging` connections are the actual connections to the remote hosts.

So, for example, if you run `steampipe query "select * from staging_servers.exec_command where command = 'uname -a';"` it will run the query on both `server1-staging` and `server2-staging` connections.

You can still run queries on individual connections, for example `steampipe query "select * from server1-staging.exec_command where command = 'uname -a';"` will only run the query on `server1-staging`.

```hcl
connection "staging_servers" {
  plugin = "exec"
  
  type = "aggregator"
  connections = [ "*-staging" ]
}

connection "server1-staging" {
  plugin = "exec"

  protocol = "ssh"
  host = "my-remote-linux-host"
  user = "my-username"
  private_key = "~/.ssh/my-remote-linux-host.pem"
}

connection "server2-staging" {
  plugin = "exec"

  protocol = "ssh"
  host = "my-remote-linux-host"
  user = "my-username"
  private_key = "~/.ssh/my-remote-linux-host.pem"
}
```

#### Local connection using a specific interpreter

##### ZSH interpreter

```hcl
connection "exec_local" {
  plugin = "exec"
  interpreter = ["/bin/zsh", "-c"]
}
```

##### Python3 interpreter

```hcl
connection "exec_local" {
  plugin = "exec"
  interpreter = ["/bin/python3", "-c"]
}
```

##### Perl interpreter

```hcl
connection "exec_local" {
  plugin = "exec"
  interpreter = ["/bin/perl", "-e"]
}
```

> Pro-tip: If you are going to use multiple interpreters, you can name the connection name to reflect the interpreter used. For example, `connection "exec_local_bash" { ... }` or `connection "exec_local_python" { ... }` etc.

## Limiting concurrent connections

For each query (command) executed, this plugin opens a new SSH/WinRM connection. If you are running a lot of queries against the same host, these connection attempts may be seen as abusive activity.

To reduce the chance of getting flagged, the plugin has a default `max_concurrency` limiter set to `15`. However, this limiter can be toggled by defining a `limiter` resource in your `exec.spc` configuration file:

```hcl
connection "exec_local" {
  plugin      = "exec"
  working_dir = "."
}

plugin "exec" {
  limiter "exec_max_concurrency_limiter" {
    max_concurrency = 15
  }
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-exec
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
