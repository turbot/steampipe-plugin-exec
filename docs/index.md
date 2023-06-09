---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/exec.svg"
brand_color: "#8467F3"
display_name: "Exec"
short_name: "exec"
description: "Steampipe plugin to run & query shell commands on local and remote servers."
og_description: "Run & query shell commands with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/exec-social-graphic.png"
---

# Exec + Steampipe

TODO

### Configuration

There are multiple ways to configure the plugin depending on the type of connection you want to use. Bellows are some examples of how to configure the most common connections.

#### Local connection

```hcl
connection "exec_local" {
  plugin = "exec"
}
```

#### Remote Linux connection through SSH with private key

```hcl
connection "exec_remote_linux" {
  plugin = "exec"

  host = "my-remote-linux-host"
  user = "my-username"
  private_key = "~/.ssh/my-remote-linux-host.pem"
}
```

#### Remote Linux connection through SSH with in-file private key

```hcl
connection "exec_remote_linux" {
  plugin = "exec"

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
  
  host = "my-remote-linux-host"
  user = "my-username"
  password = "my-password"
}
```

#### Remote Windows connection through WinRM

```hcl
connection "exec_remote_windows" {
  plugin = "exec"

  protocol = "winrm"
  host = "18.228.214.45"
  port = 5986
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
