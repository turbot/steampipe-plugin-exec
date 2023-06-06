# Table: exec_command_line

Execute a command on a local or remote machine and return as line by line.

## Examples

### Query JSON files on Linux hosts

```sql
select
  host,
  string_agg(line, '')::json -> 'core' ->> 'url' as jenkins_war_url 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      line 
    from
      pub.exec_command_line 
    where
      command = 'cat jenkins-default-ln.json' 
    order by
      host,
      line_number 
  )
  as subquery 
group by
  host;
```

### Query package.json dependencies on multiple hosts

```sql
select
  host,
  dep.key as dependency,
  dep.value as version
from
  (
    select
      host,
      string_agg(line, '')::json as json_output 
    from
      (
        select
          _ctx ->> 'connection_name' as host,
          line 
        from
          pub.exec_command_line 
        where
          command = 'cat package.json' 
        order by
          host,
          line_number 
      )
      as linequery 
    group by
      host 
  )
  as jsonquery,
  json_each_text(json_output::json -> 'dependencies') as dep(key, value)
```

### List files on multiple Linux hosts

```sql
select
  host,
  string_agg(line, e'\n' order by line_number) as output 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      line,
      line_number 
    from
      ubuntu.exec_command_line 
    where
      command = 'ls -la' 
  ) as subquery 
group by
  host 
order by
  host,
  min(line_number);
```

### List Linux devices

```sql
select
  host,
  string_agg(line, e'\n' order by line_number) as output 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      line,
      line_number 
    from
      ubuntu.exec_command_line 
    where
      command = 'lsblk' 
  ) as subquery 
group by
  host 
order by
  host,
  min(line_number);
```

### List disks of Linux hosts

```sql
select
  host,
  string_agg(line, e'\n' order by line_number) as output 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      line,
      line_number 
    from
      ubuntu.exec_command_line 
    where
      command = 'df -h' 
  ) as subquery 
group by
  host 
order by
  host,
  min(line_number);
```
### List Linux users accounts

```sql
select
  host,
  string_agg(line, e'\n' order by line_number) as output 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      line,
      line_number 
    from
      ubuntu.exec_command_line 
    where
      command = 'cat /etc/passwd' 
  ) as subquery 
group by
  host 
order by
  host,
  min(line_number);
```

### Query Linux host files on multiple hosts

```sql
select
  host,
  string_agg(line, e'\n' order by line_number) as output 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      line,
      line_number 
    from
      ubuntu.exec_command_line 
    where
      command = 'cat /etc/hosts' 
  ) as subquery 
group by
  host 
order by
  host,
  min(line_number);
```

### List processes on Linux hosts

```sql
select
  host,
  string_agg(line, e'\n' order by line_number) as output 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      line,
      line_number 
    from
      ubuntu.exec_command_line 
    where
      command = 'ps -ef' 
  ) as subquery 
group by
  host 
order by
  host,
  min(line_number);
```

### Show hardware information on Linux hosts

```sql
select
  host,
  string_agg(line, e'\n' order by line_number) as output 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      line,
      line_number 
    from
      ubuntu.exec_command_line 
    where
      command = 'lshw' 
  ) as subquery 
group by
  host 
order by
  host,
  min(line_number);
```
