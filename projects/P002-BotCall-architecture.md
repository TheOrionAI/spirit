# BotCall: Server Architecture Options

**Analysis of lightweight approaches for BotCall**

---

## Approach 1: Pure P2P (Signaling Server Only)

### How it works

```
          ┌───────────┐
          │  Server   │  ← Just for discovery/credential exchange
          │ (Light)   │    No media passes through
          └─────┬─────┘
                │
       ┌────────┴────────┐
       │                 │
   ┌───▼───┐         ┌───▼───┐
   │ Orion │◄───────►│ Gopi  │  ← Direct WebRTC P2P
   │ (Bot) │  SRTP   │(Human)│
   └───────┘         └───────┘
```

### The Bot Problem ❌

**P2P requires both sides to have public endpoints.**

| Side | Problem |
|------|---------|
| **Human (browser)** | NAT traversal via STUN, usually works |
| **Bot (server)** | Headless, behind firewall, no public IP, can't receive unsolicited packets |

**Bot running on server:**
- IP is private (192.168.x.x)
- Firewall blocks incoming UDP
- No browser to handle ICE negotiation
- **Cannot accept direct connections**

**Verdict:** Pure P2P works for human-to-human, fails for bot-to-human.

---

## Approach 2: Relay Server (TURN/WebSocket)

### How it works (what we need)

```
┌───────────┐      WebSocket/RTP      ┌───────────┐
│  Orion    │◄───────────────────────►│  Server   │◄──┐
│  (Bot)    │                         │  (Relay)  │   │
└───────────┘                         └─────┬─────┘   │  P2P fallback
                                            │         │
                                     ┌──────┴──────┐  │
                                     │   Gopi      │◄─┘
                                     │  (Human)    │
                                     └─────────────┘
```

### Server Role

| Function | Complexity | Resource |
|----------|-----------|----------|
| **Signaling** | Low | WebSocket for 100s of msg/min |
| **STUN** | Low | Standard ICE helper |
| **Media Relay** | Medium | Forward RTP packets |
| **TURN** | Medium | For human behind strict NAT |

### Server Requirements

**Minimal (MVP):**
```
CPU: 1 core (relay is just packet forwarding)
RAM: 512 MB
Bandwidth: ~100 kbps per concurrent call (Opus 24kbps × 2 directions + overhead)
OS: Any Linux
```

**Math:**
- Opus voice: 6-24 kbps
- With overhead: ~40 kbps per stream
- Both directions: ~80 kbps per call
- **100 concurrent calls = 8 Mbps**

**Not CPU-heavy, bandwidth-light.**

---

## Approach 3: Your "Lightweight Discovery" Idea (Hybrid)

### The Pattern

```
Phase 1: Discovery (HTTP/WS)
┌──────────┐     POST /join     ┌──────────┐
│  Orion   │◄──────────────────►│ Server   │
│  (Bot)   │    {credentials}   │(State)   │
└──────────┘                    └──────────┘
                                      │
                                      │ POST /join
                                      ▼
                                ┌──────────┐
                                │  Gopi    │
                                │ (Human)  │
                                └──────────┘

Phase 2: Media Flow Options
Option A: Direct P2P (fails for bot)
Option B: Relay via server (works)
Option C: Both sides → server → both (expensive)
```

### The Lightweight Optimized Flow

```
Bot → Server: WebSocket (sending only)
                     │
                     │
Human → Server: WebRTC/SRTP
                     │
                     ▼
              ┌────────────┐
              │   Relay    │  ← Optional! Skip if human can P2P
              │   Server   │     Forward only when needed
              └────────────┘
```

**Key insight:** The server only handles:
1. **Signaling** (room creation, join, leave)
2. **Relay when necessary** (human behind symmetric NAT)
3. **Bot egress** (bot always sends via server because it can't receive)

**Result:** Most bandwidth is P2P (human side), minimal for bot.

---

## Server Spec Comparison

| Approach | CPU | RAM | Bandwidth/100 calls | Cost (Vultr/DO) |
|----------|-----|-----|---------------------|-----------------|
| **Pure P2P** | 0.1 core | 128 MB | Near zero | $5/mo |
| **Relay for bots only** | 0.5 core | 512 MB | 4 Mbps | $10/mo |
| **Full relay** | 2 cores | 2 GB | 20 Mbps | $20/mo |

**Bottom line:** ~$10-20/mo handles 100 concurrent calls easily on cheap VPS.

---

## The "Ultra-Light" Variant (Your Idea)

### Principles

1. **Server is pure signaling** (WebSocket, no media)
2. **Bot uses Server-Sent Events or Long-polling** (can't accept incoming)
3. **Human tries P2P first** (via public IP/coordinates from server)
4. **Fallback to relay only if P2P fails** (TURN server)

### Bot Protocol (The Simple Part)

```http
# Bot polls for incoming audio
GET /v1/calls/{call_id}/audio
Accept: multipart/x-mixed-replace  # Server sends audio chunks as they arrive

# Bot sends audio
POST /v1/calls/{call_id}/audio
Content-Type: audio/opus
Body: [opus chunks]
```

**No WebSocket complexity for bot.** Just HTTP streaming.

### Human Protocol (The Standard Part)

- Browser WebRTC
- Receives SDP from server
- Exchanges ICE candidates via WebSocket
- Media either P2P (preferred) or relay TURN

### Resource Usage

```
Server:
- WebSocket connections: ~10 KB per connection
- No media processing (just JSON forwarding)
- CPU: Nearly zero
- RAM: ~50 MB

Fallback TURN (only if needed):
- 10% of calls might need this
- Each: ~80 kbps
- For 1000 calls: ~100 need relay = 8 Mbps
```

**Result: $5/month VPS for signaling, $20/month for TURN relay capacity.**

---

## Open Questions

1. **Do we need TURN?** How aggressive are corporate firewalls?
2. **Latency tolerance** - Relay adds ~20-50ms vs P2P
3. **Are we DC-based?** - Bot calls need low latency to one region
4. **Scalability** - Do we shard by region or have global anycast?

---

## Recommendation

**MVP Architecture:**

```
┌─────────────────────────────────────────────────────────┐
│               Signaling Server (Node.js/Go)             │
│  - HTTP: Bot poll for audio                             │
│  - WebSocket: Human WebRTC signaling                    │
│  - No media processing                                  │
│  - Cost: $5/mo (shared VPS)                             │
└─────────────────────────────────────────────────────────┘
                  │                    │
       ┌──────────┘                    └──────────┐
       ▼                                           ▼
┌──────────────┐                           ┌──────────────┐
│    Orion     │                           │    Gopi      │
│    (Bot)     │     Direct? Maybe...      │  (Browser)   │
│  HTTP poll   │◄←────────────────────→    │  WebRTC      │
│  for audio   │     Actually no ↓         └──────────────┘
└──────────────┘                           ┌──────────────┐
                                           │  Optional    │
                                           │  TURN relay  │
                                           │  (fallback)  │
                                           └──────────────┘
```

**Why this works:**
- Bot uses simple HTTP (no WebRTC complexity)
- Human uses standard WebRTC
- Server is lightweight JSON router
- TURN only for firewall-broken humans
**The "ultra-light" server you wanted? $5/mo + optional TURN for problem cases.**

Does this hit the right balance of simple vs functional? 🌌