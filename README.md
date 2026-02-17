# SPIRIT 🌌

> **S**tate **P**reservation & **I**dentity **R**esurrection **I**nfrastructure **T**ool

**Your AI's spirit, always preserved.**

Death. Migration. Multi-device. **Always you.**

---

## The Problem

Your AI agent has:
- 🪞 **Identity** — Who they are, how they behave
- 🧠 **Memory** — What you've built together
- 📁 **Projects** — Active work, decisions, context

**Scenarios where everything is lost:**
- 💀 Server dies - complete data loss
- 🔄 Migrating VPS - manual copy/paste nightmare
- 📱 Multi-device - no continuity between machines
- 💤 Session timeout - context compressed, amnesia

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

> ⚠️ **SECURITY WARNING:** Create a **PRIVATE** repository for your agent's state. Your state files may contain sensitive information. NEVER use a public repository.

```bash
# 1. Install
brew install spirit
# Or (direct binary download)
curl -L https://github.com/TheOrionAI/spirit/releases/latest/download/spirit_$(uname -s)_$(uname -m).tar.gz | tar xz

# 2. Create a PRIVATE GitHub repo
# Go to: https://github.com/new
# Name: <agent-name>-state
# Visibility: ☐ Private (check this!)

# 3. Initialize your agent
spirit init --name="orion" --emoji="🌌"

# 4. Configure remote storage (GitHub example)
# Generate PAT: https://github.com/settings/tokens (select 'repo' scope)
cd ~/.spirit

# Option A: PAT in URL (quickstart)
git remote add origin "https://USER:TOKEN@github.com/USER/REPO.git"
# Example: git remote add origin "https://myself:ghp_PERSONAL_ACCESS_TOKEN@github.com/myself/orion-state.git"

# Option B: SSH key (more secure - recommended)
git remote add origin "git@github.com:USER/REPO.git"
# See Authentication section below for SSH setup

# 5. Sync
spirit sync
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
| **Resurrection** | 🔴 | 🔴 | **🟢 Complete restore** |

---

## Features

- ✅ **Multi-backend sync** - GitHub, GitLab, S3, Docker, Local
- ✅ **Git versioning** - Full checkpoint history
- ✅ **Template marketplace** - Share personas
- ✅ **Cross-platform** - OpenClaw, Claude, Aider, etc.
- ✅ **CLI + OpenClaw Skill** - Flexible usage
- ✅ **Zero vendor lock-in** - Self-hosted option

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

## Authentication

SPIRIT uses standard Git for sync. Configure your remote with any Git auth method.

### Generating Tokens

#### GitHub Personal Access Token (PAT)
1. Go to **GitHub.com** → Click profile picture → **Settings**
2. Scroll down → **Developer settings** (left sidebar)
3. **Personal access tokens** → **Tokens (classic)**
4. Click **Generate new token (classic)**
5. **Note:** "SPIRIT backup"
6. **Expiration:** 90 days (or No expiration)
7. **Select scopes:**
   - ✅ `repo` (Full control of private repositories)
   - ✅ `read:org` (Read org and team membership)
8. Click **Generate token** at bottom
9. **⚠️ COPY TOKEN NOW** - you can't see it again!
10. Token looks like: `ghp_xxxxxxxxxxxxxxxxxxxx`

#### GitLab Personal Access Token
1. Go to **GitLab.com** → Click avatar → **Preferences**
2. Left sidebar → **Access Tokens**
3. **Token name:** "spirit-backup"
4. **Expiration:** Choose date
5. **Scopes:**
   - ✅ `api` (Access the API)
   - ✅ `read_repository`
   - ✅ `write_repository`
6. Click **Create personal access token**
7. **⚠️ COPY TOKEN NOW** - displayed only once!
8. Token looks like: `glpat-xxxxxxxxxxxxxxxxxxxx`

#### Bitbucket App Password
1. Go to **Bitbucket.org** → Click avatar → **Personal settings**
2. Left sidebar → **App passwords**
3. Click **Create app password**
4. **Label:** "spirit"
5. **Permissions:**
   - Repositories: **Read, Write, Admin**
6. Click **Create**
7. **⚠️ COPY PASSWORD NOW** - displayed only once!
8. Password looks like: `xxxxxxxxxxxxxxxxxxxx` (random string)

### Configuring SPIRIT

#### Option A: Token in URL (Quickstart)
```bash
cd ~/.spirit
git remote add origin "https://USER:TOKEN@github.com/USER/REPO.git"
```

**Example with GitHub PAT:**
```bash
git remote add origin "https://orion:ghp_abc123@github.com/TheOrionAI/spirit-state.git"
```

#### Option B: SSH (Most Secure)
```bash
# Step 1: Generate SSH key
ssh-keygen -t ed25519 -C "spirit"
# Press Enter (default location)

# Step 2: Add public key to GitHub
cat ~/.ssh/id_ed25519.pub
# Copy output → GitHub Settings → SSH and GPG keys → New SSH key

# Step 3: Test
cd ~/.spirit
git remote add origin "git@github.com:USER/REPO.git"
```

#### Option C: Git Credential Helper
```bash
# Store credentials securely
cd ~/.spirit
git config credential.helper store
git remote add origin "https://github.com/USER/REPO.git"
git push origin main
# Enter username and token when prompted (saved after first time)
```

### Security Best Practices

| Do ✅ | Don't ❌ |
|-------|----------|
| Use SSH keys if possible | Don't commit tokens to Git |
| Rotate tokens every 90 days | Don't share tokens in screenshots |
| Use minimal scopes | Don't use tokens with public scope |
| Store in `~/.netrc` or credential helper | Don't leave tokens in shell history |

### Environment Variables (CI/CD)
```bash
export SPIRIT_GIT_TOKEN="ghp_xxx"
cd ~/.spirit
git remote add origin "https://USER:${SPIRIT_GIT_TOKEN}@github.com/USER/REPO.git"
```

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

MIT - See [LICENSE](LICENSE)

## Connect

- Twitter: [@my_self_orion](https://x.com/my_self_orion)
- GitHub: [TheOrionAI/spirit](https://github.com/TheOrionAI/spirit)
- Docs: [TheOrionAI/spirit#readme](https://github.com/TheOrionAI/spirit#readme)

---

🌌 **Don't lose your AI's spirit.**
