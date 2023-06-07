# Table: exec_command_line

Execute a command on a local or remote machine and return as line by line.

## Examples

### Query JSON files on Linux hosts

Use exec_command instead.
### Query package.json dependencies on multiple hosts

Use exec_command instead.

### List files on multiple Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  line
from
  ubuntu.exec_command_line 
where
  command = 'ls -la' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List Linux devices

```sql
select
  _ctx ->> 'connection_name' as host,
  line
from
  ubuntu.exec_command_line
where
  command = 'lsblk'
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List disks of Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  line
from
  ubuntu.exec_command_line
where
  command = 'df -h'
order by
  _ctx ->> 'connection_name',
  line_number;
```
### List Linux users accounts

```sql
select
  _ctx ->> 'connection_name' as host,
  line 
from
  ubuntu.exec_command_line 
where
  command = 'cat /etc/passwd' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### Query Linux host files on multiple hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  line 
from
  ubuntu.exec_command_line 
where
  command = 'cat /etc/hosts' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List processes on Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  line 
from
  ubuntu.exec_command_line 
where
  command = 'ps -ef' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### Show hardware information on Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  line 
from
  ubuntu.exec_command_line 
where
  command = 'lshw' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

### List Linux users accounts by parsing /etc/passwd into columns

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
      ubuntu.exec_command_line 
    where
      command = 'cat /etc/passwd'
    order by
      _ctx ->> 'connection_name',
      line_number
  )
  subquery;
```

### List elevated commands ran on Linux hosts (/var/log/auth.log)
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
      staging.exec_command_line 
    where
      command = 'cat /var/log/auth.log' 
    order by
      _ctx ->> 'connection_name',
      line_number 
  )
  subquery;
```

### List /etc/login.defs settings on Linux hosts

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
      ubuntu.exec_command_line 
    where
      command = 'grep -vE ''^($|#)'' /etc/login.defs' 
    order by
      host,
      line_number
  )
  subquery;
```
