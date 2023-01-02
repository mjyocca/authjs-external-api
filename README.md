# Authjs External Api

[Auth.js](https://authjs.dev/) social sign in example with an external API backend written in [GO](https://go.dev/). The GO API server uses JWE (Json Web Encryption) A256GCM middleware that aligns with the Auth.js/Next-auth Package when JWT is enabled.

Authjs Providers:
  - Github
  - Google

## Frontend (Nextjs)

Setup

Generate secret to use with both **frontend** and **backend** used as the `NEXTAUTH_SECRET`

```
openssl rand -base64 32
```

Create env file `.env.local`

```
GITHUB_ID=[INFO HERE]
GITHUB_SECRET=[INFO HERE]

GOOGLE_CLIENT_ID=[INFO HERE]
GOOGLE_CLIENT_SECRET=[INFO HERE]

NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=[INFO HERE]

EXTERNAL_API_ENDPOINT=http://localhost:8000
```

Start Frontend Service

*Need pnpm installed*
```bash
pnpm install \
pnpm dev
```

## Backend (Go)

Create env file `.env.local`

```
NEXTAUTH_SECRET=[INFO HERE]
```

Start Backend Service

```bash
cd backend \
docker compose up
```