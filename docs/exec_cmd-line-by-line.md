# Table: exec_cmd

Execute a command on a local or remote machine and return as line by line.

## Examples

### Query JSON files on Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output_lines::jsonb -> 'core' ->> 'url' as jenkins_war_url
from
  pub.exec_cmd,
  jsonb_array_elements_text(output) as output_lines
where
  command = 'cat jenkins-default.json';
```

### Query package.json dependencies on multiple hosts

TODO

### List files on multiple Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output_line 
from
  ubuntu.exec_cmd,
  jsonb_array_elements_text(output) as output_line 
where
  command = 'ls -la';
```

### List Linux devices

```sql
select
  _ctx ->> 'connection_name' as host,
  output_lines
from
  ubuntu.exec_cmd,
  jsonb_array_elements_text(output) as output_lines
where
  command = 'lsblk';
```

### List disks of Linux hosts

```sql
select
  _ctx ->> 'connection_name' as host,
  output_lines
from
  ubuntu.exec_cmd,
  jsonb_array_elements_text(output) as output_lines
where
  command = 'df -h';
```

### List Linux users accounts

```sql
select
  _ctx ->> 'connection_name' as host,
  output_lines
from
  ubuntu.exec_cmd,
  jsonb_array_elements_text(output) as output_lines
where
  command = 'cat /etc/passwd';
```

### Query Linux host files on multiple hosts
```sql
select
  _ctx ->> 'connection_name' as host,
  output_lines
from
  ubuntu.exec_cmd,
  jsonb_array_elements_text(output) as output_lines
where
  command = 'cat /etc/hosts';
```

### List processes on Linux hosts

```sql 
select
  _ctx ->> 'connection_name' as host,
  output_lines
from
  ubuntu.exec_cmd,
  jsonb_array_elements_text(output) as output_lines
where
  command = 'ps -ef';
```

### Show hardware information on Linux hosts
```sql
select
  _ctx ->> 'connection_name' as host,
  output_lines
from
  ubuntu.exec_cmd,
  jsonb_array_elements_text(output) as output_lines
where
  command = 'lshw';
```
