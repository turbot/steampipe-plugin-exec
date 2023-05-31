# Table: exec_run

Execute a command on a local or remote machine and return as line by line.

## Examples

### List files on Linux hosts

```sql
select * from ubuntu.exec_run where command = 'ls -la' order by _ctx ->> 'connection_name'
select * from ubuntu.exec_run where command = 'ls -la' and line_by_line order by _ctx ->> 'connection_name', line_number
select * from exec_local.exec_run where command = 'ls -la'
select * from exec_local.exec_run where command = 'ls -la' and line_by_line order by line_number
```

### Query JSON files on Linux hosts

```sql
select line::jsonb -> 'core' ->> 'url' jekins_war_url from pub.exec_run where command = 'cat jenkins-default.json'
```

### Query package.json dependencies on multiple hosts

```sql
SELECT
  dep.key AS dependency,
  dep.value AS version,
  _ctx->>'connection_name' AS host
FROM
  ubuntu.exec_run,
  json_each_text(line::json->'dependencies') AS dep(key, value)
where
  command = 'cat package.json';
```

### List Linux devices

```sql
select * from ubuntu.exec_run where command = 'lsblk'
```

### List disks of Linux hosts

```sql
select * from ubuntu.exec_run where command = 'df -h' order by _ctx ->> 'connection_name', line_number
```

### List local disks on a Mac OSX

```sql
select * from exec_local.exec_run where command = 'diskutil list' order by line_number
```

### List Linux users accounts

```sql
select * from ubuntu.exec_run where command = 'cat /etc/passwd' order by _ctx ->> 'connection_name', line_number
```

### Query Linux host files on multiple hosts

```sql
select line, line_number, _ctx->>'connection_name' AS host from ubuntu.exec_run where command = 'cat /etc/hosts' order by _ctx ->> 'connection_name', line_number
```

### List processes on Linux hosts

```sql
select * from ubuntu.exec_run where command = 'ps -ef' order by _ctx ->> 'connection_name', line_number
```

### List local processes

```sql
select * from exec_local.exec_run where command = 'ps -ef' order by line_number
```

### List processes on Windows hosts

```sql
select * from windows.exec_run where command = 'tasklist' order by _ctx ->> 'connection_name', line_number
```

### List logged in users on Linux hosts

```sql
select * from ubuntu.exec_run where command = 'w' order by _ctx ->> 'connection_name', line_number
```

### Show hardware information on Linux hosts

```sql
select * from ubuntu.exec_run where command = 'lshw' order by _ctx ->> 'connection_name', line_number
```

### Show hardware information on Windows hosts

```sql
select * from windows.exec_run where command = 'wmic computersystem get model,name,manufacturer,systemtype' order by _ctx ->> 'connection_name', line_number
```
