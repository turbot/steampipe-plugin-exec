# Table: exec_remote_command_result

Execute a command on the remote machine and return as a single row or line by line.

## Examples

### List files on multiple Linux hosts

```sql
select * from ubuntu.exec_remote_command_result where command = 'ls -la'
```

### List files on multiple hosts, line by line

```sql
select * from ubuntu.exec_remote_command_result where command = 'ls -la' and line_by_line order by _ctx ->> 'connection_name', line_number
```

### List files on Windows hosts

```sql
select * from windows.exec_remote_command_result where command = 'dir'
```
