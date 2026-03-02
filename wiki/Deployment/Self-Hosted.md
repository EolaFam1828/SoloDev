# Self-Hosted Deployment

Production-style deployment of SoloDev on your own infrastructure.

## Architecture

```
┌────────────────┐    ┌──────────────┐    ┌──────────────┐
│ Reverse Proxy  │───▶│ SoloDev      │───▶│ PostgreSQL   │
│ (nginx/Caddy)  │    │ Binary       │    │              │
│ TLS term.      │    │ Port 3000    │    │ Port 5432    │
└────────────────┘    └──────────────┘    └──────────────┘
```

## Prerequisites

- Linux server (Ubuntu 22.04+ or similar)
- PostgreSQL 14+
- Docker (for Gitspaces and pipeline execution)
- Reverse proxy (nginx, Caddy, or similar) for TLS
- Domain name (for HTTPS)

## Database Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE solodev;
CREATE USER solodev WITH PASSWORD '<strong-password>';
GRANT ALL PRIVILEGES ON DATABASE solodev TO solodev;
```

Configure the database connection in the server environment.

## Binary Deployment

Build or download the SoloDev binary:

```bash
make build
```

Run as a systemd service:

```ini
[Unit]
Description=SoloDev
After=network.target postgresql.service

[Service]
Type=simple
User=solodev
ExecStart=/opt/solodev/gitness server /opt/solodev/.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Reverse Proxy

### nginx Example

```nginx
server {
    listen 443 ssl;
    server_name solodev.example.com;

    ssl_certificate /etc/ssl/certs/solodev.crt;
    ssl_certificate_key /etc/ssl/private/solodev.key;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Backups

Back up the PostgreSQL database regularly:

```bash
pg_dump solodev > solodev-backup-$(date +%Y%m%d).sql
```

Back up the data directory (Git repositories, artifacts).

## Security Considerations

- Use strong passwords for database and admin accounts
- Enable TLS via reverse proxy
- Restrict Docker socket access
- Keep the binary and dependencies updated
- Configure firewall rules to expose only ports 443 and 22 (SSH)
