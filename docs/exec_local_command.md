# Table: exec_local_command

Execute a command on the local machine and return as line by line.

## Examples

### List local files

```sql
select * from exec_local_command where command = 'ls -la' order by line_number
```

### Query JSON files

```sql
select line::jsonb -> 'core' ->> 'url' jekins_war_url from exec_local_command where command = 'cat jenkins-default.json'
```

### Query package.json dependencies

```sql
SELECT
  dep.key AS dependency,
  dep.value AS version,
  _ctx->>'connection_name' AS host
FROM
  exec_local_command,
  json_each_text(line::json->'dependencies') AS dep(key, value)
where
  command = 'cat package.json';
```

### List local devices on a Linux machine

```sql
select * from exec_local_command where command = 'lsblk'
```

### List local disks on a Linux machine

```sql
select * from exec_local_command where command = 'df -h' order by line_number
```

### List local disks on a Mac OSX

```sql
select * from exec_local_command where command = 'diskutil list' order by line_number
```

### List Linux users accounts

```sql
select * from exec_local_command where command = 'cat /etc/passwd' order by line_number
```

### Query Linux host file

```sql
select line, line_number, _ctx->>'connection_name' AS host from exec_local_command where command = 'cat /etc/hosts' order by line_number
```

### List local processes

```sql
select * from exec_local_command where command = 'ps -ef' order by line_number
```

### List logged in users on Linux hosts

```sql
select * from exec_local_command where command = 'w' order by line_number
```

### Show hardware information on Linux hosts

```sql
select * from exec_local_command where command = 'lshw' order by line_number
```
