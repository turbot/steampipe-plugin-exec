## Local Squid Proxy setup
```sh
docker-compose up -d
```
- host/port: http://127.0.0.1:8080
- username: `luis`
- password: `c@v41c@nt3`

#### To add a new user and password:
```sh
#Access proxy docker container
docker-compose exec squid bash

# INSTALL htpasswd command
apt update && apt install apache2-utils

# ADD a user
docker-compose exec squid htpasswd -c /etc/squid/passwords <username>
```
