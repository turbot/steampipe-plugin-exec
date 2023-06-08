# Defines a local connection
connection "exec_local" {
  plugin = "exec"

  # Working directory to use when running commands.
  # Default is the current directory where steampipe is running.
  # For example:
  # working_dir = "/tmp"

  # Shell interpreter to use for commands.
  # Defaults: Windows is "cmd /C", otherwise "sh -c".
  # For example, to use bash:
  # interpreter = [ "/bin/bash", "-c" ]
}

# Defines a remote connection to a Linux host
connection "exec_linux" {
  plugin = "exec"

  # Protocol for the connection, set ssh for Linux hosts.
  protocol = "ssh"

  # Host to connect to, either an IP address or a hostname.
  host = "my-remote-linux-host"

  # Username for the remote host connection.
  user = "ubuntu"

  # Credentials, either password or private_key.
  password = "my_password"
  private_key = <<EOK
-----BEGIN RSA PRIVATE KEY-----
... snipped ...
-----END RSA PRIVATE KEY-----
EOK

  # Optional - Proxy settings
  # Enables the plugin to connect to host through a HTTP proxy.
  # If username and password are not set, then it will try to connect to the proxy without authentication.

  # proxy_host = "127.0.0.1"
  # proxy_port = 8080
  # proxy_user_name = "my_proxy_user"
  # proxy_user_password = "proxy_password"


  # Optional - Bastion connection settings
  # Enables connecting to host through a bastion host. The plugin will connect to bastion_host first, and then connect from there to host.

  # bastion_user = "ec2-user"
  # bastion_host = "52.67.108.24"
  # bastion_port = 22
  # bastion_password = "my_password"
  # bastion_private_key = <<EOK
-----BEGIN RSA PRIVATE KEY-----
... snipped ...
-----END RSA PRIVATE KEY-----
EOK

}

# Defines a remote connection to a Windows host
connection "exec_windows" {
  plugin = "exec"

  # Protocol for the connection, set winrm for Windows hosts.
  protocol = "winrm"

  # Host to connect to, either an IP address or a hostname.
  host = "18.228.214.45"

  # Optional - Port to connect to, defaults to 5985
  # port = 5986

  # Username for the remote host connection.
  user = "Administrator"

  # Password for the remote host connection.
  password = "rh=PM76t54nouv&dqwe3cNM7J1(*skZhh*"

  # Optional - WinrM over HTTPS instead of HTTP. Defaults to false
  # https = true

  # Optional - Ignore certificate for HTTPS connection. Defaults to false
	# insecure = true
}
