# Defines a local connection
connection "exec_local" {
  plugin = "exec"
}

# Defines a remote connection
connection "exec" {
  plugin = "exec"

  # Working directory to use when running commands.
  # Default is the current directory where steampipe is running.
  # For example:
  # working_dir = "/tmp"

  # Shell interpreter to use for commands.
  # Defaults: Windows is "cmd /C", otherwise "sh -c".
  # For example, to use bash:
  # interpreter = [ "/bin/bash", "-c" ]

  # Protocol for the connection, either ssh or winrm.
  protocol = "ssh"

  # Host to connect to, either an IP address or a hostname.
  host = "my-remote-linux-host"

  # Port to connect to.
  # port = 5986

  # Username for the remote host connection.
  user = "ubuntu"

  # Credentials, either password or private_key.
  password = "my_password"
  private_key = <<EOK
-----BEGIN RSA PRIVATE KEY-----
... snipped ...
-----END RSA PRIVATE KEY-----
EOK

  # Optional - Setting the following enables the plugin connect to host through a HTTP proxy. If username and password are not set, then it will try to connect to the proxy without authentication.
  # proxy_host = "127.0.0.1"
  # proxy_port = 8080
  # proxy_user_name = "my_proxy_user"
  # proxy_user_password = "proxy_password"


  ## Linux only

  # Optional - Setting the following enables the bastion Host connection. The plugin will connect to bastion_host first, and then connect from there to host.
  # bastion_user = "ec2-user"
  # bastion_host = "52.67.108.24"
  # bastion_port = 22

  # Credentials, either password or private_key.
  # bastion_password = "my_password"
  # bastion_private_key = <<EOK
-----BEGIN RSA PRIVATE KEY-----
... snipped ...
-----END RSA PRIVATE KEY-----
EOK

  ## Windows only

  # WinrM over HTTPS. Defaults to false
  # https = true

  # Ignore certificate for HTTPS connection. Defaults to false
	# insecure = true
}
