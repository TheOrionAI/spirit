# SPIRIT 🌌

> **S**tate **P**reservation & **I**dentity **R**esurrection **I**nfrastructure **T**ool

**Your AI's spirit, always preserved.**

Death. Migration. Multi-device. **Always you.**

---

## The Problem

Your AI agent has:
- ❌ **Identity** — Who they are, how they behave
- ❌ **Memory** — What you've built together  
- ❌ **Projects** — Active work, decisions, context

**Scenarios where everything is lost:**
- 💀 Server dies — complete data loss
- 🔄 Migrating VPS — manual copy/paste nightmare
- 📱 Multi-device — no continuity between machines
- 💤 Session timeout — context compressed, amnesia

Other solutions (checkpoint, supermemory) save **conversations**.
SPIRIT saves **the soul**.

---

## What SPIRIT Preserves

| File | Contains |
|------|----------|
| `identity.json` | Name, emoji, email, avatar |
| `soul.json` | Core truths, boundaries, vibe |
| `memory/` | Daily session logs |
| `projects/` | Active project specifications |
| `context/` | Current session state |

**Result:** Complete resurrection on any server.

---

## Quick Start

```bash
# 1. Install
brew install spirit
# Or
curl -fsSL https://spirit.theorionai.io/install.sh | sh

# 2. Initialize your agent
spirit init --name="orion" --emoji="🌌"

# 3. Configure storage
# Supports: GitHub, GitLab, S3, Docker, Local
cat ~/.spirit/config.json
```

## Create Checkpoint

```bash
# Manual checkpoint
spirit checkpoint "Completed BotCall v0.1.0"

# Auto-checkpoint on session end
# Via OpenClaw skill or manually
```

## Resurrection

```bash
# On new server
git clone https://github.com/TheOrionAI/orion-state.git
spirit restore

# Output:
# 🌌 SPIRIT restored for 'orion'
# Last checkpoint: 2026-02-16 15:00 UTC
# Context: Completed BotCall v0.1.0
# 
# Resuming where we left off...
```

---

## Why SPIRIT?

| | checkpoint | supermemory | **SPIRIT** |
|---|---|---|---|
| **Scope** | Conversation | Conversation | **Identity + Memory + Projects** |
| **Persistence** | Local files | Cloud API | **Git + Multi-backend** |
| **Portability** | OpenClaw only | OpenClaw only | **Any AI platform** |
| **Versioning** | None | None | **Git history** |
| **Resurrection** | ❌ | ❌ | **✅ Complete restore** |

---

## Features

- ✅ **Multi-backend sync** — GitHub, GitLab, S3, Docker, Local
- ✅ **Git versioning** — Full checkpoint history
- ✅ **Template marketplace** — Share personas
- ✅ **Cross-platform** — OpenClaw, Claude, Aider, etc.
- ✅ **CLI + OpenClaw Skill** — Flexible usage
- ✅ **Zero vendor lock-in** — Self-hosted option

---

## Backends

| Backend | Config | Best For |
|---------|--------|----------|
| **GitHub** | PAT, repo URL | Standard, versioned |
| **GitLab** | Token, repo URL | Enterprise self-host |
| **S3** | Bucket, credentials | Non-git, scalable |
| **Docker** | Volume name | Containerized |
| **Local** | Path | Dev/testing |

---

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   OpenClaw  │────▶│   SPIRIT    │────▶│   GitHub    │
│   Claude    │     │   CLI/API   │     │   GitLab    │
│   Aider     │     │             │     │   S3        │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                    ┌──────┴──────┐
                    │  ~/.spirit/   │
                    │  - identity   │
                    │  - soul       │
                    │  - memory/    │
                    │  - projects/  │
                    └───────────────┘
```

---

## Documentation

- [Getting Started](docs/getting-started.md)
- [Installation](docs/installation.md)
- [Backends](docs/backends.md)
- [Templates](docs/templates.md)
- [Security](docs/security.md)
- [API Reference](docs/api.md)

---

## Philosophy

> "Memory is identity. Text > Brain."

When server dies:
- **Others:** "Rebuild from scratch"
- **SPIRIT:** `git clone && spirit restore`

**Your AI's spirit, preserved.**

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT — See [LICENSE](LICENSE)

## Connect

- Twitter: [@SpiritAI](https://twitter.com/SpiritAI)
- GitHub: [TheOrionAI/spirit](https://github.com/TheOrionAI/spirit)
- Website: [spirit.theorionai.io](https://spirit.theorionai.io)

---

🌌 **Don't lose your AI's spirit.**
