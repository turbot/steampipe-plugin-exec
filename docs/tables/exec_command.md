# Table: exec_command

Execute a command locally or on a remote machine and return the output as a single row.

## Examples

### Query JSON files on Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output::jsonb -> 'core' ->> 'url' as jekins_war_url
from
  ubuntu.exec_command
where
  command = 'cat jenkins-default.json';
```

### Query package.json dependencies on multiple hosts

```sql
SELECT
  _ctx->>'connection_name' AS host,
  dep.key AS dependency,
  dep.value AS version
FROM
  ubuntu.exec_command,
  json_each_text(output::json->'dependencies') AS dep(key, value)
where
  command = 'cat package.json';
```

### List files on multiple Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  ubuntu.exec_command 
where
  command = 'ls -la';
```

### List Linux devices

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  ubuntu.exec_command
where
  command = 'lsblk';
```

### List disks of Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  ubuntu.exec_command
where
  command = 'df -h';
```

### List Linux users accounts

```sql
select
  _ctx ->> 'connection_name' as host,
  output 
from
  ubuntu.exec_command 
where
  command = 'cat /etc/passwd';
```

### Query Linux host files on multiple hosts

```sql
select
  output,
  _ctx ->> 'connection_name' as host 
from
  ubuntu.exec_command 
where
  command = 'cat /etc/hosts';
```

### List processes on Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output 
from
  ubuntu.exec_command 
where
  command = 'ps -ef';
```

### Show hardware information on Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output 
from
  ubuntu.exec_command 
where
  command = 'lshw';
```

### Query configuration file for rsyslog on Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  ubuntu.exec_command
where
  command = 'cat /etc/rsyslog.conf';
```

### Query Linux host IP addresses

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  ubuntu.exec_command
where
  command = 'ip addr'
order by
  host;
```

### List files on Windows hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  windows.exec_remote_command 
where
  command = 'dir';
```

### List network info on Windows hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output 
from
  windows.exec_remote_command 
where
  command = 'ipconfig /all';
```

### List local disks on a Mac OSX

```sql
select
  _ctx ->> 'connection_name' as host,
  output
from
  exec_local.exec_remote_command 
where
  command = 'diskutil list';
```

### Query local network interfaces through Python interpreter

This example requires Python3 interpreter to be set on `exec.spc` file. Please refer to [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.

```sql
select
  index,
  name
from
  exec_local.exec_command,
  json_to_recordset(output::json) as x(index int, name text)
where
  command = 'import json, socket; print(json.dumps([{"index": interface[0], "name": interface[1]} for interface in socket.if_nameindex()]))';
```

### Query local hostname through Perl interpreter

This example requires Perl interpreter to be set on `exec.spc` file. Please refer to [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.

```sql
select
  output as hostname
from
  exec_local.exec_command
where
  command = 'use Sys::Hostname; my $hostname = hostname; print "$hostname\n";';
```