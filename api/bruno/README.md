# TrieOH API Collection

Bruno collection for all TrieOH services.

## Setup

1. Install [Bruno](https://www.usebruno.com/)
2. Open → point at this `bruno/` folder
3. Select the `local` environment
4. Run `Login` under `identityx`, copy the token into the `access_token` env var

## Ports (local defaults)

| Service    | Port |
|------------|------|
| identityx  | 8001 |
| informd    | 8002 |
| payssage   | 8003 |
| univents   | 8004 |

## CI

```sh
npx @usebruno/cli run bruno/ --env local
```
