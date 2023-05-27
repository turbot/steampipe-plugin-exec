# Table: exec_command

Execute a command locally or on a remote machine and return as a single row.

## Examples

### List files on multiple Linux hosts

```sql
select * from ubuntu.exec_command where command = 'ls -la'
```
### List files on Windows hosts

```sql
select * from windows.exec_command where command = 'dir'
```

### Query JSON files on Linux hosts

```sql
select output::jsonb -> 'core' ->> 'url' jekins_war_url from pub.exec_command where command = 'cat jenkins-default.json'
```

### Query package.json dependencies on multiple hosts

```sql
SELECT
  dep.key AS dependency,
  dep.value AS version,
  _ctx->>'connection_name' AS host
FROM
  ubuntu.exec_command,
  json_each_text(output::json->'dependencies') AS dep(key, value)
where
  command = 'cat package.json';
```

### List Linux devices

```sql
select * from ubuntu.exec_command where command = 'lsblk'
```

### List disks of Linux hosts

```sql
select * from ubuntu.exec_command where command = 'df -h' order by _ctx ->> 'connection_name'
```

### List local disks on a Mac OSX

```sql
select * from exec_local.exec_command where command = 'diskutil list'
```


### List Linux users accounts

```sql
select * from ubuntu.exec_command where command = 'cat /etc/passwd'
```

### Query Linux host files on multiple hosts

```sql
select output, _ctx->>'connection_name' AS host from ubuntu.exec_command where command = 'cat /etc/hosts'
```
