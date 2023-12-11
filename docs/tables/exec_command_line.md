---
title: "Steampipe Table: exec_command_line - Query Exec Command Lines using SQL"
description: "Allows users to query Exec Command Lines, providing insights into command line execution results and potential issues."
---

# Table: exec_command_line - Query Exec Command Lines using SQL

Exec is a service that enables the execution of arbitrary commands on the local system. It provides a mechanism to run commands and scripts, and capture their output for further processing. Exec helps users perform system-level operations, diagnose issues, and automate tasks.

## Table Usage Guide

The `exec_command_line` table provides insights into command line execution within the local system. As a system administrator or DevOps engineer, explore command-specific details through this table, including execution results, exit status, and associated metadata. Utilize it to uncover information about command execution, such as those with unexpected results, the status of executed commands, and the verification of command outputs.

## Examples

### List files on Linux host
Explore the list of files on a Linux host to better manage your system resources and understand file distribution. This can be particularly useful for identifying unnecessary files that may be taking up valuable space.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  line
from
  exec_command_line 
where
  command = 'ls -la' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  line
from
  exec_command_line 
where
  command = 'ls -la' 
order by
  json_extract(_ctx, '$.connection_name'),
  line_number;
```

### List files on Linux host and show the stream used for output
This query is used to analyze the output stream from a specific command executed on a Linux host. It can help in identifying any potential issues or errors that occurred during the execution of the command.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  line,
  stream
from
  exec_command_line 
where
  command = 'ls non_existing_file' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  line,
  stream
from
  exec_command_line 
where
  command = 'ls non_existing_file' 
order by
  json_extract(_ctx, '$.connection_name'),
  line_number;
```

### List devices on Linux host
This query allows you to analyze the devices connected to a Linux host. It's particularly useful for understanding the composition of your hardware resources and managing them effectively.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  line
from
  exec_command_line
where
  command = 'lsblk'
order by
  _ctx ->> 'connection_name',
  line_number;
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  line
from
  exec_command_line
where
  command = 'lsblk'
order by
  json_extract(_ctx, '$.connection_name'),
  line_number;
```

### List disks on Linux host
Explore which disks are available on a Linux host to gain insights into storage usage and capacity. This information can help you manage your storage resources more effectively.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  line
from
  exec_command_line
where
  command = 'df -h'
order by
  _ctx ->> 'connection_name',
  line_number;
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  line
from
  exec_command_line
where
  command = 'df -h'
order by
  json_extract(_ctx, '$.connection_name'),
  line_number;
```

### List user accounts on Linux host
Explore the user accounts present on a Linux host to gain a comprehensive understanding of who has access to the system. This is useful for auditing purposes and ensuring only authorized users have access.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  line 
from
  exec_command_line 
where
  command = 'cat /etc/passwd' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  line 
from
  exec_command_line 
where
  command = 'cat /etc/passwd' 
order by
  json_extract(_ctx, '$.connection_name'),
  line_number;
```

### Query host file on Linux host
Analyze the host file on a Linux host to understand the static hostname-to-IP mappings. This is useful for troubleshooting network issues and verifying the correct resolution of hostnames.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  line 
from
  exec_command_line 
where
  command = 'cat /etc/hosts' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  line 
from
  exec_command_line 
where
  command = 'cat /etc/hosts' 
order by
  json_extract(_ctx, '$.connection_name'),
  line_number;
```

### List processes on Linux host
Explore the active processes across various Linux hosts. This is particularly useful for system administrators to monitor and manage system processes, ensuring optimal performance and security.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  line 
from
  exec_command_line 
where
  command = 'ps -ef' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  line 
from
  exec_command_line 
where
  command = 'ps -ef' 
order by
  json_extract(_ctx, '$.connection_name'),
  line_number;
```

### Show hardware information on Linux host
Discover detailed hardware information on your Linux host. This query helps you gain insights into your system's hardware configuration, which can be beneficial for troubleshooting, system upgrades, or general maintenance.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  line 
from
  exec_command_line 
where
  command = 'lshw' 
order by
  _ctx ->> 'connection_name',
  line_number;
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  line 
from
  exec_command_line 
where
  command = 'lshw' 
order by
  json_extract(_ctx, '$.connection_name'),
  line_number;
```

### List user accounts on Linux host by parsing /etc/passwd file into columns
Explore which user accounts exist on a Linux host and gain insights into their configuration details, such as whether they have a password, their user and group IDs, and their home directory and shell settings. This is useful for understanding the security and organization of your Linux system.

```sql+postgres
select
  host,
  split_output[1] as username,
  case when split_output[2] = 'x' then true else false end as has_password, 
  split_output[3] as user_id, 
  split_output[4] as group_id, 
  split_output[5] as user_comment, 
  split_output[6] as home_directory, 
  split_output[7] as shell 
from
  (
    select
      _ctx ->> 'connection_name' as host,
      string_to_array(line, ':') as split_output 
    from
      exec_command_line 
    where
      command = 'cat /etc/passwd'
    order by
      _ctx ->> 'connection_name',
      line_number
  )
  subquery;
```

```sql+sqlite
Error: SQLite does not support split or string_to_array functions.
```

### List elevated commands ran on Linux host (/var/log/auth.log)
Identify instances where elevated commands were run on a Linux host. This is useful for security auditing and detecting potentially malicious activities.
Ref: https://regex101.com/r/ImwoFl/3


```sql+postgres
select
  matches[1] as month,
  matches[2] as day,
  matches[3] as hour,
  matches[4] as hostname,
  matches[7] as pwd,
  matches[8] as elevated_user,
  matches[9] as command 
from
  (
    select
      regexp_matches(line, '^(\S{3})? {1,2}(\S+) (\S+) (\S+) (.+?(?=\[)|.+?(?=))[^a-zA-Z0-9](\d{1,7}|)[^a-zA-Z0-9]{1,3}PWD=([^ ]+) ; USER=([^ ]+) ; COMMAND=(.*)$') as matches 
    from
      exec_command_line 
    where
      command = 'cat /var/log/auth.log' 
    order by
      _ctx ->> 'connection_name',
      line_number 
  )
  subquery;

```sql+sqlite
Error: SQLite does not support regex operations.
```

### List /etc/login.defs settings on Linux host
Determine the settings of a Linux host's login configuration. This query is useful for system administrators seeking to understand and manage user login preferences on a Linux system.

```sql+postgres
select
  host,
  matches[1] as option,
  matches[2] as setting
from
  (
    select
      _ctx ->> 'connection_name' as host,
      regexp_matches(line, '^(\S+)\s+(\S+)') as matches
    from
      exec_command_line 
    where
      command = 'grep -vE ''^($|#)'' /etc/login.defs' 
    order by
      host,
      line_number
  )
  subquery;
```

```sql+sqlite
Error: SQLite does not support regex operations.
```

### List installed packages on Linux/Debian host
Determine the software packages installed on a Linux/Debian host. This is particularly useful for system administrators needing to manage software installations and updates across multiple systems.

```sql+postgres
select
  host,
  matches[1] as option,
  matches[2] as setting
from
  (
    select
      _ctx ->> 'connection_name' as host,
      regexp_matches(line, '^(\\S+)\\h(\\S+)\\h(\\S+)\\h(.*)') as matches
    from
      exec_command_line 
    where
      command = 'apt list --installed' 
    order by
      host,
      line_number
    limit 20
  )
  subquery;
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### Query processor through Python interpreter on local machine
Analyze your local machine's processor type by executing a Python command. This can be useful for system diagnostics or when tailoring software to specific hardware.
This example requires Python3 interpreter to be set on `exec.spc` file. Please refer [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.


```sql+postgres
select
  line as processor
from
  exec_command_line
where
  command = 'import platform; print(platform.processor())';
```

```sql+sqlite
select
  line as processor
from
  exec_command_line
where
  command = 'import platform; print(platform.processor())';
```

### Query disk usage through Python interpreter on local machine
Analyze the settings to understand the usage of disk space on your local machine through the Python interpreter. This can help you manage your storage effectively by identifying how much space is used and how much is still available.
This example requires Python3 interpreter to be set on `exec.spc` file. Please refer [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.

```sql+postgres
select
  total,
  used,
  free
from
  exec_command_line,
  json_to_record(line::json) as x(total bigint, used bigint, free bigint)
where
  command = 'import json, shutil; du = shutil.disk_usage("/"); print(json.dumps({"total": du[0], "used": du[1], "free": du[2]}))';
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

