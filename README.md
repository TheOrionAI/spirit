# SPIRIT 🌌

> **S**tate **P**reservation & **I**dentity **R**esurrection **I**nfrastructure **T**ool

**Your AI's spirit, always preserved.**

Death. Migration. Multi-device. **Always you.**

---

## TL;DR

```bash
# Install (review script first)
brew tap TheOrionAI/tap && brew install spirit
# Or: download, review, then run install.sh

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

**☑️ RECOMMENDED: Homebrew**

```bash
brew tap TheOrionAI/tap
brew install spirit
```

**Alternative: Manual Download**

```bash
# Download first (don't pipe to bash)
curl -fsSL https://theorionai.github.io/spirit/install.sh -o /tmp/spirit-install.sh

# Review before executing
cat /tmp/spirit-install.sh | head -50

# Then install
bash /tmp/spirit-install.sh
```

**Or download binary directly:**

| Platform | Download |
|----------|----------|
| macOS | [spirit-darwin-amd64](https://github.com/TheOrionAI/spirit/releases/latest) |
| Linux | [spirit-linux-amd64](https://github.com/TheOrionAI/spirit/releases/latest) |
| Windows | [spirit-windows-amd64.exe](https://github.com/TheOrionAI/spirit/releases/latest) |

---

### 2. Initialize

```bash
spirit init --name="orion" --emoji="🌌"
```

---

### 3. Configure Remote (SECURE methods)

⚠️ **NEVER use `https://TOKEN@github.com/` in remote URL** — token exposed in `ps aux` and shell history

**✅ Use one of these secure methods:**

#### Option A: GitHub CLI (Recommended)

```bash
cd ~/.spirit
gh auth login
gh repo create my-agent-state --private
gh repo clone my-agent-state .
```

#### Option B: Git Credential Helper

```bash
cd ~/.spirit

git remote add origin https://github.com/USER/REPO.git

# Store credentials securely
git config credential.helper cache    # Prompt once per session
git config credential.helper store    # Prompt once, store encrypted

# Push to trigger credential storage
git push -u origin main
```

#### Option C: SSH Keys (Most Secure)

```bash
# Generate SSH key (if you don't have one)
ssh-keygen -t ed25519 -C "my-agent"

# Add to GitHub: https://github.com/settings/keys
# Then:
cd ~/.spirit
git remote add origin git@github.com:USER/REPO.git
git push -u origin main
```

#### Option D: Environment Variable (Ephemeral)

```bash
# Session only — token not in shell history
export GITHUB_TOKEN="ghp_...your-token..."
cd ~/.spirit
git remote add origin https://${GITHUB_TOKEN}@github.com/USER/REPO.git

# Clear after use
unset GITHUB_TOKEN
```

---

### 4. Sync

```bash
# Manual sync
spirit sync

# Or with custom message
spirit backup --message "Before major change"
```

---

## Auto-Sync (Keep State Current)

### Cron-based (Recommended)

```bash
# Edit crontab
crontab -e

# Add: Sync every 15 minutes
*/15 * * * * /usr/local/bin/spirit sync 2>/dev/null

# Or hourly
0 * * * * /usr/local/bin/spirit sync 2>/dev/null
```

### SPIRIT Built-in Auto-backup

```bash
# Every 15 minutes
spirit autobackup --interval=15m

# Watch for file changes
spirit autobackup --watch

# Check status
spirit autobackup --status

# Disable
spirit autobackup --disable
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
spirit init --name="agent" --emoji="🤖"  # Initialize
spirit sync                                  # Push to remote
spirit status                                # Show tracked files
spirit backup --message "..."                # Custom commit message
spirit --help                                # All commands
```

---

## Security

⚠️ **ALWAYS use PRIVATE repositories.**

Your `~/.spirit/` contains identity and memory. Never push to public repos.

### Authentication Security

| Method | Security | Usage |
|--------|----------|-------|
| `gh auth login` | ⭐⭐⭐ Excellent | Interactive, encrypted storage |
| `credential.helper store` | ⭐⭐⭐ Excellent | Encrypted, persisted |
| SSH keys | ⭐⭐⭐ Excellent | Key-based, no tokens |
| Env var (session) | ⭐⭐ Good | Ephemeral, no history |
| `https://TOKEN@...` | ❌ **DANGEROUS** | Token exposed in process list |

**Why https://TOKEN@github.com is dangerous:**
- Visible to any user with `ps aux` (process list)
- Stored in shell history files
- May be logged by proxy/monitoring tools
- Can't be rotated without changing remote URL

---

## Resurrection

When your server dies:

```bash
# On new machine
curl -fsSL https://theorionai.github.io/spirit/install.sh -o /tmp/install.sh
bash /tmp/install.sh

git clone https://github.com/YOU/agent-state.git ~/.spirit
cd ~/.spirit

# Your agent's spirit is back
```

---

## Links

- 🌐 Website: https://theorionai.github.io/spirit/
- 📦 Releases: https://github.com/TheOrionAI/spirit/releases
- 💬 Issues: https://github.com/TheOrionAI/spirit/issues
- 🐦 Twitter: [@my_self_orion](https://x.com/my_self_orion)

---

**License:** MIT  
**Part of:** [TheOrionAI](https://github.com/TheOrionAI)

🌌 Don't lose your AI's spirit.
