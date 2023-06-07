# Table: exec_command

Execute a command locally or on a remote machine and return as a single row.

## Examples

### Query JSON files on Linux hosts

does not apply

### Query package.json dependencies on multiple hosts

does not apply

### List files on multiple Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  regexp_split_to_table(output, '\n') as line
from
  ubuntu.exec_command 
where
  command = 'ls -la';
```

### List Linux devices

```sql
select
  _ctx ->> 'connection_name' as host,
  regexp_split_to_table(output, '\n') as line
from
  ubuntu.exec_command
where
  command = 'lsblk';
```

### List disks of Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  regexp_split_to_table(output, '\n') as line
from
  ubuntu.exec_command
where
  command = 'df -h';
```

### List Linux users accounts

```sql
select
  _ctx ->> 'connection_name' as host,
  regexp_split_to_table(output, '\n') as line 
from
  ubuntu.exec_command 
where
  command = 'cat /etc/passwd';
```

### Query Linux host files on multiple hosts

```sql
select
  regexp_split_to_table(output, '\n') as line,
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
  regexp_split_to_table(output, '\n') as line 
from
  ubuntu.exec_command 
where
  command = 'ps -ef';
```

### Show hardware information on Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  regexp_split_to_table(output, '\n') as line 
from
  ubuntu.exec_command 
where
  command = 'lshw';
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
      string_to_array(unnest(string_to_array(output, '\n')), ':') as split_output 
    from
      ubuntu.exec_command
    where
      command = 'cat /etc/passwd'
    order by
      _ctx ->> 'connection_name'
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
      regexp_matches(unnest(string_to_array(output, E'\n')), '^(\S{3})? {1,2}(\S+) (\S+) (\S+) (.+?(?=\[)|.+?(?=))[^a-zA-Z0-9](\d{1,7}|)[^a-zA-Z0-9]{1,3}PWD=([^ ]+) ; USER=([^ ]+) ; COMMAND=(.*)$') as matches 
    from
      staging.exec_command 
    where
      command = 'cat /var/log/auth.log' 
  )
  subquery;
```
