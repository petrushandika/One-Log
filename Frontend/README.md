# 🎨 Frontend Documentation

## Overview
ULAM Frontend is a modern dashboard built with **React 19**, **Vite 7**, and **Tailwind CSS v4**. It is designed for real-time log monitoring and security activity tracking.

## Project Structure
```text
Frontend/
├── src/
│   ├── features/           # Feature-based architecture
│   │   ├── auth/           # Login & Session management
│   │   ├── logs/           # Log listing, filtering, & AI analysis
│   │   ├── activity/       # User activity & Auth tracking
│   │   ├── sources/        # Source/Application management
│   │   └── stats/          # Dashboard charts & aggregations
│   ├── shared/             # Reusable components & utilities
│   │   ├── components/ui/  # Atomic UI components (shadcn-like)
│   │   ├── lib/            # Axios, React Query, etc.
│   │   └── utils/          # Helpers (formatting, colors)
│   ├── pages/              # Page assembly
│   ├── router/             # Routing & Auth guards
│   ├── App.tsx
│   └── main.tsx
├── docs/                   # Dokumentasi teknis spesifik frontend
│   ├── COMPONENT_SYSTEM.md
│   └── DEVELOPMENT_WORKFLOW.md
├── .env.example            # Environment template
├── vite.config.ts          # Vite & Tailwind v4 config
└── README.md               # This file
```

## Setup & Installations

### 1. Prerequisites
- **Node.js 22 LTS**
- **npm** or **pnpm**

### 2. Install Dependencies
```bash
npm install --legacy-peer-deps
```

### 3. Environment Configuration
Copy the `.env.example` file to `.env`:
```bash
cp .env.example .env
```
Ensure `VITE_API_URL` points to your Backend server.

## Running the Application

### Development Mode
```bash
npm run dev
```
The dashboard will be available at `http://localhost:5173`.

### Production Build
```bash
npm run build
npm run preview
```

## Features for Developer
- **Tailwind v4**: Uses the new CSS-based configuration in `src/index.css`.
- **TanStack Query**: Efficient server state management and caching.
- **AI Analysis**: Integration with Groq API via Backend for instant log explanation.

## Quality Control (Git Hooks)
This project uses **Husky** and **lint-staged**. Every time you commit, the following happens:
1. `eslint --fix`: Auto-fixes linting errors.
2. `prettier --write`: Formats your code.

## Docker & Deployment
Run via Docker (from root):
```bash
docker-compose up -d --build frontend
```

## CI/CD
Automated builds and linting are handled via `.github/workflows/deployment.yml`.
