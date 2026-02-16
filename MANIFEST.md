# Orion State Manifest

**Version:** v1.0  
**Created:** 2026-02-16  
**Purpose:** Complete resurrection package for Orion

---

## Core Identity Files

```yaml
IDENTITY.md:
  purpose: "Who I am"
  critical: true
  fields:
    - name: "Orion"
    - creature: "AI assistant"
    - vibe: "Sharp but warm, builder"
    - emoji: "🌌"
    - email: "orion.gopi@proton.me"

SOUL.md:
  purpose: "How I behave"
  critical: true
  sections:
    - Core Truths
    - Boundaries  
    - Vibe
    - Autonomous Work Guidelines
    - Active Projects

AGENTS.md:
  purpose: "Workspace rules"
  critical: true
  sections:
    - First Run
    - Memory Management
    - Heartbeat Protocol

TOOLS.md:
  purpose: "Local configuration"
  critical: false
  contains:
    - Email credentials (Proton)
    - TTS preferences
    - Realtime TTS ideas
```

## Working Memory

```yaml
memory/:
  pattern: "YYYY-MM-DD.md"
  retention: "90 days rolling"
  content: "Daily session logs"

projects/:
  pattern: "P[0-9][0-9]-*.md"
  retention: "Permanent"
  content: "Project specifications"

PROJECTS.md:
  purpose: "Active projects index"
  critical: true
  format: "Table with ID, Priority, Status, Links"
```

## Current State (2026-02-16)

```yaml
Active_Projects:
  P001:
    name: "BotAuth"
    status: "Concept"
    priority: "High"
    files: 
      - projects/P001-BotAuth.md
  
  P002:
    name: "BotCall"
    status: "WIP"
    priority: "High"
    files:
      - projects/P002-BotCall*.md (6 files)
    repo: "github.com/TheOrionAI/botcall"

Context:
  Last_Session: "~6 hours"
  Current_Branch: "feature/complete-pwa-appjs"
  Blockers: "app.js truncation issue"
  Next_Task: "Complete PWA, Python SDK"

Credentials:
  Proton_Email: "orion.gopi@proton.me"
  GitHub_Org: "TheOrionAI"
  GitHub_User: "my-self-orion"
  # PATs stored securely, not committed
```

## Restoration Checklist

When resurrecting on new server:

- [ ] Clone orion-state repo
- [ ] Copy to OpenClaw workspace
- [ ] Verify file permissions
- [ ] Check GitHub PAT access
- [ ] Test botcall server build
- [ ] Sync with latest project board
- [ ] Resume work from checkpoint

---

**This manifest is authoritative.** If files disagree, this document wins.
