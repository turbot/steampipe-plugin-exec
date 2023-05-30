# Table: exec_local_command_result

Execute a command on the local machine and return as a single row.

## Examples

### List local files

```sql
select * from exec_local_command_result where command = 'ls -la'
```

### Query JSON files

```sql
select line::jsonb -> 'core' ->> 'url' jekins_war_url from exec_local_command_result where command = 'cat jenkins-default.json'
```

### Query package.json dependencies

```sql
SELECT
  dep.key AS dependency,
  dep.value AS version,
  _ctx->>'connection_name' AS host
FROM
  exec_local_command_result,
  json_each_text(line::json->'dependencies') AS dep(key, value)
where
  command = 'cat package.json';
```

### List local devices on a Linux machine

```sql
select * from exec_local_command_result where command = 'lsblk'
```

### List local disks on a Linux machine

```sql
select * from exec_local_command_result where command = 'df -h'
```

### List local disks on a Mac OSX

```sql
select * from exec_local_command_result where command = 'diskutil list'
```

### List Linux users accounts

```sql
select * from exec_local_command_result where command = 'cat /etc/passwd'
```

### Query Linux host file

```sql
select line, line_number, _ctx->>'connection_name' AS host from exec_local_command_result where command = 'cat /etc/hosts'
```

### List local processes

```sql
select * from exec_local_command_result where command = 'ps -ef'
```

### List logged in users on Linux hosts

```sql
select * from exec_local_command_result where command = 'w'
```

### Show hardware information on Linux hosts

```sql
select * from exec_local_command_result where command = 'lshw'
```
