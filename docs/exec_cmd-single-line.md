# Table: exec_cmd

Execute a command on a local or remote machine return output as an array of lines.

## Examples

### Query JSON files on Linux hosts

TODO

### Query package.json dependencies on multiple hosts

TODO

### List files on multiple Linux hosts

```sql
select
  host,
  string_agg(output_lines, E'\n') as lines
from (
  select
    _ctx ->> 'connection_name' as host,
    output_lines
  from
    ubuntu.exec_cmd,
    jsonb_array_elements_text(output) as output_lines
  where
    command = 'ls -la'
) as substring
group by
host;
```

### List Linux devices

```sql
select
  host,
  string_agg(output_lines, E'\n') as lines
from (
  select
    _ctx ->> 'connection_name' as host,
    output_lines
  from
    ubuntu.exec_cmd,
    jsonb_array_elements_text(output) as output_lines
  where
    command = 'lsblk'
) as substring
group by
host;
```

### List disks of Linux hosts

```sql
select
  host,
  string_agg(output_lines, E'\n') as lines
from (
  select
    _ctx ->> 'connection_name' as host,
    output_lines
  from
    ubuntu.exec_cmd,
    jsonb_array_elements_text(output) as output_lines
  where
    command = 'df -h'
) as substring
group by
host;
```

### List Linux users accounts

```sql
select
  host,
  string_agg(output_lines, E'\n') as lines
from (
  select
    _ctx ->> 'connection_name' as host,
    output_lines
  from
    ubuntu.exec_cmd,
    jsonb_array_elements_text(output) as output_lines
  where
    command = 'cat /etc/passwd'
) as substring
group by
host;
```

### Query Linux host files on multiple hosts

```sql
select
  host,
  string_agg(output_lines, E'\n') as lines
from (
  select
    _ctx ->> 'connection_name' as host,
    output_lines
  from
    ubuntu.exec_cmd,
    jsonb_array_elements_text(output) as output_lines
  where
  command = 'cat /etc/hosts'
) as substring
group by
host;
```

### List processes on Linux hosts

```sql
select
  host,
  string_agg(output_lines, E'\n') as lines
from (
  select
    _ctx ->> 'connection_name' as host,
    output_lines
  from
    ubuntu.exec_cmd,
    jsonb_array_elements_text(output) as output_lines
  where
  command = 'ps -ef'
) as substring
group by
host;
```

### Show hardware information on Linux hosts

```sql
select
  host,
  string_agg(output_lines, E'\n') as lines
from (
  select
    _ctx ->> 'connection_name' as host,
    output_lines
  from
    ubuntu.exec_cmd,
    jsonb_array_elements_text(output) as output_lines
  where
  command = 'lshw'
) as substring
group by
host;
```
