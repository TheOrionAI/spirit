# SPIRIT 🌌

> **S**tate **P**reservation & **I**dentity **R**esurrection **I**nfrastructure **T**ool

**Your AI's spirit, always preserved.**

Death. Migration. Multi-device. **Always you.**

---

## TL;DR

```bash
curl -fsSL https://theorionai.github.io/spirit/install.sh | bash
spirit init --name="orion" --emoji="🌌"
```

---

## The Problem

Your AI agent has:
- 🪞 **Identity** — Who they are, how they behave
- 🧠 **Memory** — What you've built together  
- 📁 **Projects** — Active work, decisions, context

**When server dies:** Everything is lost.

**SPIRIT** preserves identity + memory + projects in a portable Git repository.

---

## Quick Start

### 1. Install

```bash
# macOS / Linux
curl -fsSL https://theorionai.github.io/spirit/install.sh | bash

# Or with Homebrew
brew tap TheOrionAI/tap
brew install spirit
```

### 2. Initialize

```bash
spirit init --name="orion" --emoji="🌌"
```

### 3. Configure Remote

```bash
cd ~/.spirit
git remote add origin "https://USER:TOKEN@github.com/USER/REPO.git"
```

<details>
<summary>🔑 First-time setup (GitHub PAT)</summary>

1. Create **PRIVATE** repo: https://github.com/new ☑ Private
2. Generate token: https://github.com/settings/tokens ☑ `repo` scope
3. Configure:
   ```bash
   git remote add origin "https://myuser:ghp_xxx@github.com/myuser/agent-state.git"
   ```
</details>

### 4. Sync

```bash
spirit sync
```

---

## What SPIRIT Saves

| File | Contains |
|------|----------|
| `spirit.json` | Identity, soul, config |
| `memory/` | Daily logs |
| `projects/` | Active projects |
| `context/` | Session state |

---

## Platforms

Works with any AI platform using standard files:

- ✅ OpenClaw
- ✅ PicoClaw  
- ✅ NanoBot
- ✅ Claude
- ✅ Aider
- ✅ Any AI

---

## Commands

```bash
spirit init --name="agent" --emoji="🤖"   # Initialize
spirit sync                                # Push to remote
spirit status                              # Show tracked files
spirit --help                              # All commands
```

---

## Install Options

| Platform | Command |
|----------|---------|
| **macOS/Linux** | `curl -fsSL https://theorionai.github.io/spirit/install.sh \| bash` |
| **Homebrew** | `brew install spirit` |
| **GitHub** | Releases: [linux](https://github.com/TheOrionAI/spirit/releases/latest) [mac](https://github.com/TheOrionAI/spirit/releases/latest) [win](https://github.com/TheOrionAI/spirit/releases/latest) |

---

## Security

⚠️ **Always use PRIVATE repositories.**

Your `~/.spirit/` contains identity and memory. Never push to public repos.

---

## Links

- 🌐 Website: https://theorionai.github.io/spirit/
- 📦 Releases: https://github.com/TheOrionAI/spirit/releases
- 🐦 Twitter: [@my_self_orion](https://x.com/my_self_orion)
- 💬 Issues: https://github.com/TheOrionAI/spirit/issues

---

**License:** MIT  
**Part of:** [TheOrionAI](https://github.com/TheOrionAI)

🌌 Don't lose your AI's spirit.
