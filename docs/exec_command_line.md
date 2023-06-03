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
