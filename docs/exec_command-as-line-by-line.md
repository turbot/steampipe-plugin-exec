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
  regexp_split_to_table(output, E'\n') as line
from
  ubuntu.exec_command 
where
  command = 'ls -la';
```

### List Linux devices

```sql
select
  _ctx ->> 'connection_name' as host,
  regexp_split_to_table(output, E'\n') as line
from
  ubuntu.exec_command
where
  command = 'lsblk';
```

### List disks of Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  regexp_split_to_table(output, E'\n') as line
from
  ubuntu.exec_command
where
  command = 'df -h';
```

### List Linux users accounts

```sql
select
  _ctx ->> 'connection_name' as host,
  regexp_split_to_table(output, E'\n') as line 
from
  ubuntu.exec_command 
where
  command = 'cat /etc/passwd';
```

### Query Linux host files on multiple hosts

```sql
select
  regexp_split_to_table(output, E'\n') as line,
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
  regexp_split_to_table(output, E'\n') as line 
from
  ubuntu.exec_command 
where
  command = 'ps -ef';
```

### Show hardware information on Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  regexp_split_to_table(output, E'\n') as line 
from
  ubuntu.exec_command 
where
  command = 'lshw';
```
