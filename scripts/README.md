# Testing

Start local docker images with these commands:
```shell
# Start a server for SSH in terminal #1
docker run -p 22:22 ghcr.io/s1ntaxe770r/image:latest

# In a second terminal, start a second server for SSH
docker run -p 23:22 ghcr.io/s1ntaxe770r/image:latest
```

Configure steampipe for the different servers and create an aggregator:
```hcl
connection "exec22" {
  plugin = "exec"
  protocol = "ssh"
  host = "127.0.0.1"
  user = "test"
  password = "test"
}

connection "exec23" {
  plugin = "exec"
  protocol = "ssh"
  host = "127.0.0.1"
  user = "test"
  password = "test"
  port = 23
}

connection "exec" {
  plugin = "exec"
  type = "aggregator"
  connections = [ "*" ]
}
```

Run some queries:
```sql
select * from exec_remote_command_result where command = 'uname -a'
select * from exec_remote_command_result where command = 'cat /etc/passwd'
select * from exec_remote_command_result where command = 'ls /tmp'
```
