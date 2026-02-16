# BotCall Repository Structure

**Decision: Monorepo for MVP → Split on v1.0**

---

## Phase 1: Monorepo (Now)

```
botcall/
├── README.md
├── LICENSE (AGPL-3.0)
├── CONTRIBUTING.md
├── Makefile
├── docker-compose.yml
├── docs/
│   ├── SPEC.md              # Protocol specification
│   ├── ARCHITECTURE.md
│   └── DEPLOYMENT.md
├── server/                  # Go discovery server
│   ├── cmd/
│   │   └── botcall-server/
│   ├── internal/
│   │   ├── api/
│   │   ├── auth/
│   │   ├── discovery/
│   │   └── websocket/
│   ├── go.mod
│   └── Dockerfile
├── pwa/                     # Web client
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── Dockerfile
├── sdk-go/                  # Go bot SDK
│   ├── client.go
│   ├── go.mod
│   └── examples/
├── sdk-python/              # Python bot SDK
│   ├── botcall/
│   ├── setup.py
│   └── examples/
└── scripts/
    ├── dev.sh
    └── deploy.sh
```

**Why monorepo for MVP:**
- Atomic commits (change protocol → server + SDKs same PR)
- Easier integration testing
- One CI/CD pipeline
- Contributors see full context
- Faster iteration

---

## Phase 2: Split Repos (v1.0 release)

When API stabilizes, split into:

```
GitHub Org: botcall
├── botcall-specs            # RFC-style protocol specs
├── botcall-server           # Go discovery server
├── botcall-client-web       # PWA
├── botcall-sdk-go           # Go SDK
├── botcall-sdk-python       # Python SDK
├── botcall-sdk-rust         # Rust SDK
└── botcall-docs             # Documentation site
```

---

## Recommendation

**Use monorepo now.** Splitting crates friction. Move to separate repos when:
- API is stable (post v1.0)
- Each component has dedicated maintainers
- You need independent versioning

For now, momentum > purity. 🚀
