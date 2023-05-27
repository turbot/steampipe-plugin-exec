# Table: exec_remote_command_result

Execute a command locally or on a remote machine and return as a single row.

## Examples

### List files on multiple Linux hosts

```sql
select * from ubuntu.exec_remote_command_result where command = 'ls -la'
```
### List files on Windows hosts

```sql
select * from windows.exec_remote_command_result where command = 'dir'
```

### Query package.json dependencies on multiple hosts

```sql
SELECT
  dep.key AS dependency,
  dep.value AS version,
  _ctx->>'connection_name' AS host
FROM
  ubuntu.exec_remote_command_result,
  json_each_text(line::json->'dependencies') AS dep(key, value)
where
  command = 'cat package.json';
```

### List Linux devices

```sql
select * from ubuntu.exec_remote_command_result where command = 'lsblk'
```

### List Linux users accounts

```sql
select * from ubuntu.exec_remote_command_result where command = 'cat /etc/passwd'
```

### Query Linux host files on multiple hosts

```sql
select line from ubuntu.exec_remote_command_result where command = 'cat /etc/hosts'
```
