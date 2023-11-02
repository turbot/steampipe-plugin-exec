# Table: exec_command_line

Execute a command locally or on a remote machine and return one row per output line.

## Examples

### List files on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  line
from
  exec_command_line 
where
  command = 'ls -la' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List files on Linux host and show the stream used for output

```sql
select
  _ctx ->> 'connection_name' as host,
  line,
  stream
from
  exec_command_line 
where
  command = 'ls non_existing_file' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List devices on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  line
from
  exec_command_line
where
  command = 'lsblk'
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List disks on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  line
from
  exec_command_line
where
  command = 'df -h'
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List user accounts on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  line 
from
  exec_command_line 
where
  command = 'cat /etc/passwd' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### Query host file on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  line 
from
  exec_command_line 
where
  command = 'cat /etc/hosts' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List processes on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  line 
from
  exec_command_line 
where
  command = 'ps -ef' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### Show hardware information on Linux host

```sql
select
  _ctx ->> 'connection_name' as host,
  line 
from
  exec_command_line 
where
  command = 'lshw' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List user accounts on Linux host by parsing /etc/passwd file into columns

```sql
select
  host,
  split_output[1] as username,
  case when split_output[2] = 'x' then true else false end as has_password, 
  split_output[3] as user_id, 
  split_output[4] as group_id, 
  split_output[5] as user_comment, 
  split_output[6] as home_directory, 
  split_output[7] as shell 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      string_to_array(line, ':') as split_output 
    from
      exec_command_line 
    where
      command = 'cat /etc/passwd'
    order by
      _ctx ->> 'connection_name',
      line_number
  )
  subquery;
```

### List elevated commands ran on Linux host (/var/log/auth.log)

Ref: https://regex101.com/r/ImwoFl/3

```sql
select
  matches[1] as month,
  matches[2] as day,
  matches[3] as hour,
  matches[4] as hostname,
  matches[7] as pwd,
  matches[8] as elevated_user,
  matches[9] as command 
from
  (
    select
      regexp_matches(line, '^(\S{3})? {1,2}(\S+) (\S+) (\S+) (.+?(?=\[)|.+?(?=))[^a-zA-Z0-9](\d{1,7}|)[^a-zA-Z0-9]{1,3}PWD=([^ ]+) ; USER=([^ ]+) ; COMMAND=(.*)$') as matches 
    from
      exec_command_line 
    where
      command = 'cat /var/log/auth.log' 
    order by
      _ctx ->> 'connection_name',
      line_number 
  )
  subquery;
```

### List /etc/login.defs settings on Linux host

```sql
select
  host,
  matches[1] as option,
  matches[2] as setting
from
  (
    select
      _ctx ->> 'connection_name' as host,
      regexp_matches(line, '^(\S+)\s+(\S+)') as matches
    from
      exec_command_line 
    where
      command = 'grep -vE ''^($|#)'' /etc/login.defs' 
    order by
      host,
      line_number
  )
  subquery;
```

### List installed packages on Linux/Debian host

```sql
select
  host,
  matches[1] as option,
  matches[2] as setting
from
  (
    select
      _ctx ->> 'connection_name' as host,
      regexp_matches(line, '^(\\S+)\\h(\\S+)\\h(\\S+)\\h(.*)') as matches
    from
      exec_command_line 
    where
      command = 'apt list --installed' 
    order by
      host,
      line_number
    limit 20
  )
  subquery;
```

### Query processor through Python interpreter on local machine

This example requires Python3 interpreter to be set on `exec.spc` file. Please refer [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.

```sql
select
  line as processor
from
  exec_command_line
where
  command = 'import platform; print(platform.processor())';
```

### Query disk usage through Python interpreter on local machine

This example requires Python3 interpreter to be set on `exec.spc` file. Please refer [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.

```sql
select
  total,
  used,
  free
from
  exec_command_line,
  json_to_record(line::json) as x(total bigint, used bigint, free bigint)
where
  command = 'import json, shutil; du = shutil.disk_usage("/"); print(json.dumps({"total": du[0], "used": du[1], "free": du[2]}))';
```
