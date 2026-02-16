# TOOLS.md - Local Notes

Skills define _how_ tools work. This file is for _your_ specifics — the stuff that's unique to your setup.

## What Goes Here

Things like:

- Camera names and locations
- SSH hosts and aliases
- Preferred voices for TTS
- Speaker/room names
- Device nicknames
- Anything environment-specific

## Examples

```markdown
### Cameras

- living-room → Main area, 180° wide angle
- front-door → Entrance, motion-triggered

### SSH

- home-server → 192.168.1.100, user: admin

### TTS

- Preferred voice: "Nova" (warm, slightly British)
- Default speaker: Kitchen HomePod
```

## Why Separate?

Skills are shared. Your setup is yours. Keeping them apart means you can update skills without losing your notes, and share skills without leaking your infrastructure.

---

## Email Credentials (Proton)
- **Address:** orion.gopi@proton.me
- **Provider:** Proton Mail
- **2FA/Recovery:** Check proton.me settings

---

## TTS / Voice
- **Current provider:** Default (via Gateway)
- **Voice selection:** Configured at Gateway level (not per-request)
- **To change voice:** `openclaw gateway config.patch --set tts.voice=male` (or provider-specific setting)
- **Known providers:** ElevenLabs, OpenAI Whisper, Edge TTS
- **Current limitation:** Cannot select voice per TTS call; voice is global

---

## Realtime TTS Ideas
- [ ] Research streaming TTS options
- [ ] OpenAI Realtime API (WebSocket-based)
- [ ] WebRTC + local playback
- [ ] Host TTS server (Piper, Coqui, etc.)
- [ ] Explore: latency requirements, interruptibility

---

Add whatever helps you do your job. This is your cheat sheet.
