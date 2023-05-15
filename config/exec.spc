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

  # Protocol for the connection, e.g. TODO
  protocol = "ssh"

  # Host to connect to, TODO - explain
  host = "localhost"

  # Username for login, TODO - explain
  user = "ubuntu"

  # Credentials, TODO - methods
  password = "my_password"
  private_key = <<EOK
-----BEGIN RSA PRIVATE KEY-----
... snipped ...
-----END RSA PRIVATE KEY-----
EOK

  # Windows only - WinrM over HTTPS. Defaults to false
  #https = true

  # Windows only - Ignore certificate for HTTPS connection. Defaults to false
	#insecure = true

}
