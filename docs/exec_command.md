# Table: exec_command

Execute a command locally or on a remote machine and return as a single row.

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
