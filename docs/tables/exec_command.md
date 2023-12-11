---
title: "Steampipe Table: exec_command - Query Exec Commands using SQL"
description: "Allows users to query Exec Commands, specifically the command execution details, providing insights into command outputs and potential anomalies."
---

# Table: exec_command - Query Exec Commands using SQL

The Exec Command is a feature that enables the execution of arbitrary commands in the context of the current session. It is a powerful tool that can be used to run scripts, utilities, and other command-line tasks. With Exec Command, you can execute commands and scripts in a secure, controlled environment, and capture the output for further processing or analysis.

## Table Usage Guide

The `exec_command` table provides insights into the execution of arbitrary commands within the current session context. As a Systems Administrator, explore command-specific details through this table, including command outputs, exit codes, and associated metadata. Utilize it to uncover information about command execution, such as error messages, the duration of command execution, and the verification of command outputs.

## Examples

### Query JSON files on Linux host
Explore the configuration of your Linux host to identify the URL of your Jenkins WAR file. This could be useful for troubleshooting or for confirming the source of your Jenkins installation.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output::jsonb -> 'core' ->> 'url' as jekins_war_url
from
  exec_command
where
  command = 'cat jenkins-default.json';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  json_extract(json(stdout_output), '$.core.url') as jekins_war_url
from
  exec_command
where
  command = 'cat jenkins-default.json';
```

### Query package.json dependencies on Linux host
Explore the dependencies and their versions in your package.json file on a Linux host. This is useful to understand the versions of libraries your project is using, which can help in debugging or updating your project.

```sql+postgres
select
  _ctx->>'connection_name' as host,
  dep.key as dependency,
  dep.value as version
from
  exec_command,
  json_each_text(stdout_output::json->'dependencies') as dep(key, value)
where
  command = 'cat package.json';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  dep.key as dependency,
  dep.value as version
from
  exec_command,
  json_each(json_extract(stdout_output, '$.dependencies')) as dep
where
  command = 'cat package.json';
```

### List files on Linux host
Explore the contents of a Linux host by listing all files within it. This can be useful for assessing the current file structure or identifying any unexpected or suspicious files.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output
from
  exec_command 
where
  command = 'ls -la';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output
from
  exec_command 
where
  command = 'ls -la';
```

### List devices on Linux host
Explore the connected devices on a Linux host system. This query is useful for system administrators who need to monitor the devices linked to their Linux servers.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output
from
  exec_command
where
  command = 'lsblk';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output
from
  exec_command
where
  command = 'lsblk';
```

### List disks on Linux host
Explore the disk usage on a Linux host to manage storage efficiently by identifying areas with high usage. This allows for proactive cleanup and allocation of resources, enhancing system performance.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output
from
  exec_command
where
  command = 'df -h';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output
from
  exec_command
where
  command = 'df -h';
```

### List user accounts on Linux host
Explore which user accounts exist on a Linux host to better manage system access and security. This can be particularly useful in maintaining control over who has access to your system and ensuring unauthorized users are not present.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output 
from
  exec_command 
where
  command = 'cat /etc/passwd';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output 
from
  exec_command 
where
  command = 'cat /etc/passwd';
```

### Query host file on Linux host
Explore the host file details on a Linux host to understand the mappings between domain names and IP addresses. This can help in troubleshooting network connectivity issues or verifying the correct setup of network services.

```sql+postgres
select
  stdout_output,
  _ctx ->> 'connection_name' as host 
from
  exec_command 
where
  command = 'cat /etc/hosts';
```

```sql+sqlite
select
  stdout_output,
  json_extract(_ctx, '$.connection_name') as host 
from
  exec_command 
where
  command = 'cat /etc/hosts';
```

### List processes on Linux host
Explore the active processes on a Linux host to understand the system's performance and resource allocation. This can help in identifying potential bottlenecks or issues that might be affecting the system's efficiency.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output 
from
  exec_command 
where
  command = 'ps -ef';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output 
from
  exec_command 
where
  command = 'ps -ef';
```

### Show hardware information on Linux host
Analyze the hardware configuration of a Linux host to understand its components and specifications. This can be useful for system administrators who need to assess the current hardware setup or plan for upgrades.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output 
from
  exec_command 
where
  command = 'lshw';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output 
from
  exec_command 
where
  command = 'lshw';
```

### Query configuration file for rsyslog on Linux host
Gain insights into the configuration of the rsyslog service on a Linux host. This is useful for understanding the current logging settings and identifying any potential issues or misconfigurations.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output
from
  exec_command
where
  command = 'cat /etc/rsyslog.conf';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output
from
  exec_command
where
  command = 'cat /etc/rsyslog.conf';
```

### Query host IP addresses on Linux host
Explore which Linux hosts have specific IP addresses. This query is useful for network management and troubleshooting, allowing you to quickly identify which hosts are using which IP addresses.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output
from
  exec_command
where
  command = 'ip addr'
order by
  host;
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output
from
  exec_command
where
  command = 'ip addr'
order by
  host;
```

### List files on Windows host
Explore the contents of a Windows host by listing all files present. This can be useful for auditing file contents or tracking down specific files.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output
from
  windows.exec_command 
where
  command = 'dir';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output
from
  windows_exec_command 
where
  command = 'dir';
```

### List network info on Windows host
Explore the network information on a Windows host to gain insights into the status and details of all active network connections. This can be useful for troubleshooting network issues or for routine network monitoring.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output 
from
  windows.exec_command 
where
  command = 'ipconfig /all';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output 
from
  windows_exec_command 
where
  command = 'ipconfig /all';
```

### List disks on a local Mac OSX
Explore the disk configuration of a local Mac OSX to gain insight into the system's storage setup. This is particularly useful for system administrators seeking to understand the disk utilization of their machines.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  stdout_output
from
  exec_command 
where
  command = 'diskutil list';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  stdout_output
from
  exec_command 
where
  command = 'diskutil list';
```

### Handle failing commands
Determine the areas in which commands are failing by analyzing the output of those commands. This can be especially useful for diagnosing and troubleshooting issues in your system.

```sql+postgres
select
  _ctx ->> 'connection_name' as host,
  case
    when exit_code = 0 then stdout_output
    else stderr_output
  end as output
from
  exec_command 
where
  command = 'ls non_existing_file';
```

```sql+sqlite
select
  json_extract(_ctx, '$.connection_name') as host,
  case
    when exit_code = 0 then stdout_output
    else stderr_output
  end as output
from
  exec_command 
where
  command = 'ls non_existing_file';
```

### Query network interfaces through Python interpreter on local machine
This query allows you to pinpoint the specific network interfaces on your local machine using a Python interpreter. In a practical setting, this can be useful for identifying potential network issues or for understanding the configuration of your local machine's network interfaces.
This example requires Python3 interpreter to be set on `exec.spc` file. Please refer [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.


```sql+postgres
select
  index,
  name
from
  exec_command,
  json_to_recordset(stdout_output::json) as x(index int, name text)
where
  command = 'import json, socket; print(json.dumps([{"index": interface[0], "name": interface[1]} for interface in socket.if_nameindex()]))';
```

```sql+sqlite
Error: SQLite does not support json_to_recordset function.
```

### Query hostname through Perl interpreter on local machine
Explore the system's hostname using the Perl interpreter on your local machine. This is useful for identifying the specific machine you're working on, especially in a networked environment with multiple machines.
This example requires Perl interpreter to be set on `exec.spc` file. Please refer [this](index.md#local-connection-using-a-specific-interpreter) on how to set it up.


```sql+postgres
select
  stdout_output as hostname
from
  exec_command
where
  command = 'use Sys::Hostname; my $hostname = hostname; print "$hostname\n";';
```

```sql+sqlite
select
  stdout_output as hostname
from
  exec_command
where
  command = 'use Sys::Hostname; my $hostname = hostname; print "$hostname\n";';
```