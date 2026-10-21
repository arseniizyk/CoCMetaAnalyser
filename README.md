# CoCMetaAnalyser

A modern tool and interactive dashboard that retrieves and visualizes Hero Equipment meta data from top Clash of Clans Legendary League players.

## 📊 [Live Interactive Dashboard](https://arseniizyk.github.io/CoCMetaAnalyser/)

## 📂 Data Storage

Meta statistics are stored as formatted JSON files in the [`frontend/public/meta`](frontend/public/meta) and [`frontend/public/metapairs`](frontend/public/metapairs) directories.

Data is collected monthly for top 10,000 players in Legendary League.

## 🚀 How to Run Locally

1. **Set Up API Access**
   - Register at the [Clash of Clans Developer Portal](https://developer.clashofclans.com/)
   - Create an API key for your IP address

2. **Configure Environment**
   - Copy `.env.example` to `.env` in the repository root:
     ```ini
     COC_API=your_token_here
     ```

3. **CLI Usage**
   - Collect both single items and pairs in a single API pass:
     ```bash
     go run ./cmd/app all -season 2025-05 -limit 10000 -filename meta10k
     ```
   - Automatically backfill all missing seasons from CoC API:
     ```bash
     go run ./cmd/app backfill -from 2025-01 -limit 10000 -filename meta10k
     ```
   - List available Legendary League seasons:
     ```bash
     go run ./cmd/app seasons
     ```

## 🤖 Automated CI/CD Workflows

- **Monthly Auto-Collector (`collect-meta.yml`)**:
  - Automatically runs on the 1st of every month via GitHub Actions cron.
  - Retrieves statistics for the previous month, commits the data to `main`.
  - Automatically triggers GitHub Pages rebuild via `deploy.yml`.
  - Can also be triggered manually via GitHub Actions **Run workflow** (supports single season collection or automated backfill).

- **GitHub Pages Deployment (`deploy.yml`)**:
  - Automatically builds and deploys the Vite + React dashboard upon changes to `main`.

## 💻 Frontend Development

```bash
cd frontend
npm install
npm run dev
```
