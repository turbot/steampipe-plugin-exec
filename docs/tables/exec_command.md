# Table: exec_command

Execute a command locally or on a remote machine and return the output as a single row.

## Examples

### Query JSON files on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  output::jsonb -> 'core' ->> 'url' as jekins_war_url
from
  exec_command
where
  command = 'cat jenkins-default.json';
```

### Query package.json dependencies on Linux host

```sql
SELECT
  _ctx->>'connection_name' as host,
  dep.key as dependency,
  dep.value as version
FROM
  exec_command,
  json_each_text(output::json->'dependencies') as dep(key, value)
where
  command = 'cat package.json';
```

### List files on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  exec_command 
where
  command = 'ls -la';
```

### List devices on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  exec_command
where
  command = 'lsblk';
```

### List disks on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  exec_command
where
  command = 'df -h';
```

### List user accounts on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  output 
from
  exec_command 
where
  command = 'cat /etc/passwd';
```

### Query host file on Linux host

```sql
select
  output,
  _ctx ->> 'connection_name' as host 
from
  exec_command 
where
  command = 'cat /etc/hosts';
```

### List processes on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  output 
from
  exec_command 
where
  command = 'ps -ef';
```

### Show hardware information on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  output 
from
  exec_command 
where
  command = 'lshw';
```

### Query configuration file for rsyslog on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  exec_command
where
  command = 'cat /etc/rsyslog.conf';
```

### Query host IP addresses on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  exec_command
where
  command = 'ip addr'
order by
  host;
```

### List files on Windows host

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  windows.exec_remote_command 
where
  command = 'dir';
```

### List network info on Windows host

```sql
select
  _ctx ->> 'connection_name' as host,
  output 
from
  windows.exec_remote_command 
where
  command = 'ipconfig /all';
```

### List disks on a local Mac OSX

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  exec_remote_command 
where
  command = 'diskutil list';
```

### Query network interfaces through Python interpreter on local machine

This example requires Python3 interpreter to be set on `exec.spc` file. Please refer [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.

```sql
select
  index,
  name
from
  exec_command,
  json_to_recordset(output::json) as x(index int, name text)
where
  command = 'import json, socket; print(json.dumps([{"index": interface[0], "name": interface[1]} for interface in socket.if_nameindex()]))';
```

### Query hostname through Perl interpreter on local machine

This example requires Perl interpreter to be set on `exec.spc` file. Please refer [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.

```sql
select
  output as hostname
from
  exec_command
where
  command = 'use Sys::Hostname; my $hostname = hostname; print "$hostname\n";';
```
