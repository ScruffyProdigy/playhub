# Kubernetes secrets (local only — never commit)

All `*.yaml` here except `*.example.yaml` are **gitignored**. Copy examples, fill in real values, keep them on your machine (backup via 1Password, etc.).

## First-time setup

```bash
cp k8s/secrets/pg-auth.example.yaml       k8s/secrets/pg-auth.yaml
cp k8s/secrets/pg-dsn.example.yaml        k8s/secrets/pg-dsn.yaml
cp k8s/secrets/jwks-secret.example.yaml   k8s/secrets/jwks-secret.yaml   # or: cd backend && go run ./scripts/jwks.go -namespace joinquest -yamlout ../k8s/secrets/jwks-secret.yaml
cp k8s/secrets/lobby-smtp.example.yaml    k8s/secrets/lobby-smtp.yaml
cp k8s/secrets/lobby-game-service.example.yaml k8s/secrets/lobby-game-service.yaml
```

Edit each file:

- **`metadata.namespace`** — must match where you deploy (`joinquest`, `playhub`, `playhub-staging`, …).
- **`lobby-smtp.yaml`** — paste your Resend API key in `SMTP_PASSWORD`.
- **`lobby-game-service.yaml`** — optional; sets Lobby’s `LOBBY_GAME_SERVICE_TOKEN` (sent to games as `lobby.serviceToken` on provision). Games store the token per match — no game-side env var.

## Apply to a cluster

```bash
kubectl config use-context YOUR_CONTEXT
kubectl apply -f k8s/secrets/
./scripts/deploy-joinquest.sh    # Lobby only
./scripts/deploy-gke.sh          # Lobby + demo-game-rps (GKE)
```

Deploy scripts apply secrets from this folder and wire SMTP into `lobby-backend` when `lobby-smtp.yaml` exists.

## Local dev (no Kubernetes)

`./scripts/dev.sh` reads **`k8s/secrets/lobby-smtp.yaml`** and sets `SMTP_*` env vars automatically (same key as GKE). You can still override in `.env` if needed.

## Per environment

Use **one set of secret files per machine** and change `namespace:` (and passwords/DSN host) when switching clusters — or maintain separate copies outside the repo (e.g. `~/secure/lobby-secrets-joinquest/`).

Do **not** use the same production DB password or JWT key as demo.
