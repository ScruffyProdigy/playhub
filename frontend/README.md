# JoinQuest Frontend

Player-facing **JoinQuest** UI: sign-in, game catalog, queues, and launch into third-party games when a match is ready.

This is not your game client — integrated titles run on their own `play_url`. JoinQuest handles lobby plumbing; authors implement fun. See **[Product vision](../docs/vision.md)** and **[Lobby handoff](../docs/lobby-protocol-handoff.md)**.

Built with React + Vite. Testing: [README_TESTING.md](./README_TESTING.md).

## Local dev

From repo root: `./scripts/dev.sh` (frontend `:5173`, API `:8080/graphql`).

## Vite plugins

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) — Babel (or [oxc](https://oxc.rs) with [rolldown-vite](https://vite.dev/guide/rolldown)) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) — SWC for Fast Refresh

## React Compiler

Not enabled in this template (dev/build cost). See [installation docs](https://react.dev/learn/react-compiler/installation).

## ESLint

For production apps, consider the [TypeScript template](https://github.com/vitejs/vite/tree/main/packages/create-vite/template-react-ts) and [`typescript-eslint`](https://typescript-eslint.io).
