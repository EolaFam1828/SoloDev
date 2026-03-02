# Quick Start

Get SoloDev running in minutes.

## 1. Clone and Start

```bash
git clone https://github.com/EolaFam1828/SoloDev.git
cd SoloDev
docker compose up -d
```

Open [http://localhost:3000](http://localhost:3000).

## 2. Create an Account

Register the first user. This user becomes the admin.

## 3. Create a Space

A space is an organizational container (like a GitHub organization). Click **New Space**, give it a name, and enter it.

## 4. Create a Repository

Inside the space, click **New Repository**. Initialize with a README or push an existing project:

```bash
git remote add solodev http://localhost:3000/git/<space>/<repo>.git
git push solodev main
```

## 5. View the Dashboard

Navigate to the SoloDev dashboard from the sidebar menu. It shows summary cards for:
- Pipelines
- Security
- Quality Gates
- Error Tracker
- Remediation
- Health Monitor
- Feature Flags
- Tech Debt

## 6. Next Steps

- [First Pipeline](First-Pipeline) — Run a CI/CD pipeline
- [First Remediation](First-Remediation) — See the AI remediation loop in action
- [Architecture Overview](../Architecture/Overview) — Understand how the system works
