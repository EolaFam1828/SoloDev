# Local Deployment

Running SoloDev locally for development and evaluation.

## Docker Compose (Recommended)

```bash
git clone https://github.com/EolaFam1828/SoloDev.git
cd SoloDev
docker compose up -d
```

| Port | Service |
|------|---------|
| 3000 | Web UI + REST API |
| 3022 | SSH (Git over SSH) |

Data persists in a Docker volume (`solodev-data`).

## From Source

```bash
make dep && make tools
cd web && yarn install && yarn build && cd ..
make build
./gitness server .local.env
```

## Database

Local deployment uses SQLite. The database file is stored in the data directory. No external database configuration is needed.

## LLM Configuration

For AI remediation, configure an LLM provider in the server configuration:
- **Ollama** (local, no API key) — Install Ollama, pull a model, point SoloDev to `http://localhost:11434`
- **Remote providers** — Set the API key for Anthropic, OpenAI, or Gemini

SoloDev functions without an LLM provider. Remediation tasks will be created but not processed.

## Resource Requirements

| Resource | Minimum |
|----------|---------|
| RAM | 4 GB allocated to Docker |
| Disk | 2 GB free |
| CPU | 2 cores |

## Troubleshooting

**Build fails during Docker Compose up**
Ensure Docker has at least 4 GB RAM allocated. On Apple Silicon, cross-compilation takes longer.

**Port 3000 in use**
Change the port mapping in `docker-compose.yml`: `"3001:3000"`.

**Permission denied on Docker socket**
The container needs `/var/run/docker.sock` for Gitspaces and pipeline execution. On Linux, add your user to the `docker` group.
