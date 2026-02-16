# BotCall: Ultra-Lightweight Voice for Bots & Humans

**ID:** P002  
**Priority:** 🔴 High  
**Status:** 💡 Concept  
**Created:** 2026-02-15  
**Last Updated:** 2026-02-15  
**Owner:** Gopi + Orion

---

## The Pitch

A purpose-built, ultra-lightweight calling platform specifically designed for AI-human communication. No browser hacks, no pretending. Bots join via native protocols, humans via PWA — with BotAuth identity verification baked in from day one.

---

## The Problem

We just hit this live: browser-based calling for AI agents is brittle:

| Platform | Why It Fails for Bots |
|----------|----------------------|
| **Jitsi/Zoom/etc** | Designed human-to-human, requires real mic/camera hardware |
| **Discord/TeamSpeak** | Bots can join, but via hacks (selenium/screensharing) |
| **Twilio** | Requires phone numbers, not packet-based voice |
| **WebRTC in browser** | No access to audio hardware from headless automation |
| **Traditional VoIP** | No identity layer, no attestation |

**The gap:** There's no calling protocol designed for "AI agent + human" from the ground up.

---

## The Solution: BotCall

**A calling platform where identity (BotAuth) is the foundation, not an afterthought.**

### Key Principles

1. **First-class bot support** — Join via HTTP/WebSocket, no browser required
2. **Lightweight** — Minimal deps, minimal latency, minimal overhead
3. **Identity-native** — BotAuth attestation required before joining
4. **Human-friendly** — PWA for mobile, no app install needed
5. **Session-based** — Ephemeral rooms, not persistent identities

---

## Technical Architecture

### Connection Flow

```
Orion (AI)          BotCall Server          Gopi (Human)
    │                      │                      │
    │ 1. Join request      │                      │
    │ {
    │   "agent": "orion",
    │   "room": "xyz-123",
    │   "attestation": "jwt..."
    │ }
    │ ─────────────────>   │
    │                      │
    │                      │ 2. Verify BotAuth
    │                      │    signature
    │                      │
    │ 3. WS recv channel   │
    │ <──────────────────  │
    │                      │
    │                      │ 4. Human joins via PWA
    │                      │ <─────────────────────
    │                      │
    │ 5. Relay audio       │
    │ <──────────────────> │
```

### Protocol Stack

```
┌─────────────────────────────────────────────────┐
│  BotAuth Identity Layer                         │
│  JWT attestation + signature verification       │
├─────────────────────────────────────────────────┤
│  Signaling Layer                                │
│  WebSocket (join, leave, mute, etc)            │
├─────────────────────────────────────────────────┤
│  Media Layer                                     │
│  WebRTC (SRTP) for humans                        │
│  Raw RTP over UDP for bots (simplified)          │
│  OR: WebSocket streaming as fallback             │
├─────────────────────────────────────────────────┤
│  Transport Layer                                 │
│  UDP preferred, TCP fallback                     │
└─────────────────────────────────────────────────┘
```

---

## Bot Join Specification

### HTTP Endpoint: POST /v1/join

**Request:**
```json
{
  "agent_id": "orion#8f3a",
  "room_id": "xyz-123-abc",
  "attestation": {
    "token": "eyJhbGc...",
    "scope": ["call:voice", "call:receive"],
    "expiry": 1708000000
  },
  "capabilities": {
    "audio_in": true,
    "audio_out": true,
    "video": false,
    "text_chat": true
  }
}
```

**Response:**
```json
{
  "session_id": "sess-abc-def",
  "websocket_url": "wss://botcall.io/v1/stream/sess-abc-def",
  "ice_servers": [  // For human side WebRTC
    {"urls": "stun:stun.botcall.io:3478"}
  ],
  "expires_at": "2025-02-15T13:00:00Z"
}
```

### WebSocket Protocol (Bot-side)

```javascript
// Connect
const ws = new WebSocket('wss://botcall.io/v1/stream/{session_id}')

// Receive audio from human (opus encoded)
ws.onmessage = (event) => {
  const packet = JSON.parse(event.data)
  if (packet.type === 'audio') {
    playAudio(packet.data)  // opus -> pcm -> speaker
  }
}

// Send audio to human
function sendAudio(opusChunk) {
  ws.send(JSON.stringify({
    type: 'audio',
    data: opusChunk,
    timestamp: Date.now()
  }))
}

// Control messages
ws.send(JSON.stringify({
  type: 'control',
  action: 'mute',
  value: true
}))
```

---

## Human Join via PWA

### User Flow

1. **Invite link** → `https://botcall.io/join/xyz-123`
2. **Browser opens PWA** (or install prompt)
3. **Request permissions** (mic)
4. **Join room** via WebRTC

### Techstack

- **Svelte or Preact** — minimal bundle size
- **WebRTC native** — getUserMedia + RTCPeerConnection
- **Service Worker** — offline capability, push notifications
- **No backend state for PWA** — all signaling via BotCall server

---

## Audio Pipeline for Bots

### Option 1: Wave file streaming (MVP)

```
TTS Engine ──► WAV file ──► WebSocket ──► BotCall Server ──► Human ear
                                    ▲
                                    │
                              STT ◄──┘ (human voice)
```

### Option 2: Real-time streaming (real deal)

```
Mic Input ──► Opus Encoder ──► WebSocket ──► BotCall Server
                                    │
                                    ▼
                              Human (WebRTC)
                                    │
                              Human Voice
                                    ▼
                              WebSocket ──► Opus Decoder ──► TTS Engine
```

### Opus Codec

- **Why Opus?** Low latency, robust to packet loss, standard for VoIP
- **Bitrates:** 6-24 kbps for voice
- **Latency:** 20ms frames typical
- **Header:** RtpEncoding with sequence numbers

---

## BotAuth Integration

### Room Creation

Only authenticated bots can create rooms:

```javascript
// Bot creates room
POST /v1/rooms/create
Authorization: Bearer {botauth_token}

Response:
{
  "room_id": "xyz-123",
  "invite_url": "https://botcall.io/join/xyz-123",
  "expires_in": 3600,
  "max_human_participants": 1
}
```

### Room Join Verification

Every join verified against BotAuth service:

```javascript
// BotCall server validates attestation
const valid = await botauth.verify({
  token: attestation_token,
  required_scope: "call:voice",
  room: room_id
})

// Checks:
// - Signature valid
// - Not expired
// - Human controller matches room owner
// - Scope includes voice call
```

---

## Prior Art & Comparison

| Platform | Bot Native | Identity Layer | Latency | Notes |
|----------|-----------|----------------|---------|-------|
| **Discord** | 🟡 (via libs) | ❌ | Low | Best current option, but not purpose-built |
| **Jitsi** | ❌ | ❌ | Low | No bot audio support |
| **Twilio** | 🟡 (API) | 🟡 | Medium | Phone-centric, not packet voice |
| **Daily.co** | 🟡 (REST) | 🟡 | Low | Closest, but $$$ |
| **Signal calls** | ❌ | ✅ | Low | No bot API for calls |
| **BotCall (this)** | ✅ | ✅ (BotAuth) | Low | Purpose-built |

---

## Open Questions

### Technical
- [ ] Raw UDP vs WebSocket for bot audio? (Perf vs simplicity)
- [ ] One-to-many bot-human? Or 1:1 only?
- [ ] Recording/audit — store calls for legal? (Ephemeral by default?)
- [ ] Mobile PWA constraints — iOS WebAudio limitations?
- [ ] TTS streaming — chunking vs true streaming?

### Product
- [ ] Free tier for devs? Credit-based like OpenAI?
- [ ] Public rooms vs invite-only?
- [ ] Group calls (1 bot, N humans)?
- [ ] Video support or voice-only?

### Legal/Ethics
- [ ] Recording consent — bot-side, human-side, both?
- [ ] Spam prevention — bot calls could be abused
- [ ] Emergency calls — redirect to 911 if needed

---

## Open Source Strategy

### Licensing

| Component | License | Rationale |
|-----------|---------|-----------|
| **Protocol (RFCs, specs)** | CC-BY-SA 4.0 | Free to implement, must attribute |
| **Server (Go/Rust)** | AGPL-3.0 | Network use = distribution, forces sharing improvements |
| **SDKs (Python, JS, etc)** | MIT | Maximum adoption, minimal friction |
| **PWA** | AGPL-3.0 | Keeps open, matches server |

**Dual licensing for enterprise:** AGPL for community, paid commercial license for orgs that can't open-source their usage.

### Governance

**Stage 1: BDFL (Benevolent Dictator)**
- You and me make decisions
- Fast iteration, clear vision

**Stage 2: Core Team**
- 3-5 maintainers with commit rights
- Weekly meetings, public roadmap

**Stage 3: Foundation**
- Non-profit foundation (like Signal)
- Open governance, RFC process

### Community Building

**Phase 1: Alpha (us)**
- Private repo, invite-only
- Dogfood with OpenClaw integration

**Phase 2: Beta (friends)**
- Public repo, limited access
- Discord/Slack for feedback

**Phase 3: Public Launch**
- Hacker News, Twitter, Reddit
- "Show HN: A calling protocol built for AI agents"

### Ecosystem Hooks

| Integration | Why | How |
|-------------|-----|-----|
| **OpenClaw** | First-class support | Native plugin |
| **LangChain** | Agent frameworks | Python SDK |
| **AutoGPT** | Popular agents | Integration guide |
| **Discord** | Bridge existing bots | Bot that bridges rooms |
| **Twilio** | Fallback SMS/voice | Optional integration |

### Value Prop for Contributors

**For devs:** "Add voice to your AI with 5 lines of code"
**For researchers:** "Study human-AI interaction with real data"
**For enterprises:** "Deploy private voice infra, no vendor lock-in"
**For activists:** "Open protocol, no corp can kill it"

### Revenue Model (Sustainable OSS)

| Service | Free Tier | Paid Tier |
|---------|-----------|-----------|
| **Hosted BotCall** | 100 mins/mo | $0.01/min |
| **Enterprise support** | Community | $5k/mo |
| **SaaS PWA** | Basic branded | Custom domain, white-label |
| **BotAuth notary** | Public chain | Private attestations |

**Principle:** Open core, paid convenience

---

## Relationship to BotAuth

**Synergy:**

```
┌─────────────────────────────────────────────────────────┐
│  BotAuth (P001)           │   BotCall (P002)           │
│  Identity + Attestation     │   Communication            │
│  ──────────────────────────┼─────────────────────────   │
│  "Who is this AI?"          │   "Let's talk"             │
│  "Who approved this?"      │   "Securely"                 │
│  ──────────────────────────┼─────────────────────────   │
│  JWT tokens                 │   WebSocket/RTP            │
│  Smart contracts            │   Opus codec               │
│  Reputation                 │   PWA clients              │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────────┐
                    │  BotComm Stack      │
                    │  Full AI-human comm  │
                    └─────────────────────┘
```

**Marketing angle:** "The first communication stack designed for the agentic web"

- BotAuth = trust layer
- BotCall = transport layer
- Together = complete solution

---

## Next Steps (Updated)

- [ ] Create GitHub org (botcall-io?)
- [ ] Draft open source governance doc
- [ ] Set up Discord for community
- [ ] Build minimal MVP (WebSocket + audio file streaming)
- [ ] Integrate with BotAuth (P001)
- [ ] Write "hello world" bot SDK
- [ ] Record demo video (me calling you via BotCall)
- [ ] Launch on HN/Reddit

---

## Notes

**Why open source wins here:**
1. **Protocol network effects** — More users = more valuable
2. **Security** — Open code = auditable = trustworthy
3. **No lock-in** — Users switch if we betray principles
4. **Talent** — Best devs want to work on OSS
5. **Mission alignment** — "AI for everyone" means "open infrastructure"

**Tagline ideas:**
- "The voice layer for the agentic web"
- "When bots need to talk"
- "Anonymous calling is for humans. Attested calling is for agents."
- "Your AI's phone number"