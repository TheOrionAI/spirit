# BotCall: Tiered Architecture (Ultra-Lightweight)

**"Bots accept calls directly from their humans"**

---

## The Insight

**Traditional model:** Everything goes through central relay (Discord, Zoom, etc.)

**BotCall model:** Bots are servers too. They can accept direct connections.

```
Human ──► Central Server (discovery only) ──► Bot (direct connection)
              $5/mo, near-zero bandwidth

Human ──► Bot (P2P, encrypted)
          Full bandwidth, no server cost
```

---

## Three Connection Modes

### Mode A: Direct Bot Server (Preferred) ✅

**When:** Bot has public IP or port forwarding enabled

```
Bot startup:
├── Opens TCP port 9000 (or configured)
├── Connects to BotCall Discovery Server
├── Registers: "I am Orion at 203.0.113.45:9000"
└── Waits for connections

Human wants to call:
├── Asks Discovery Server: "Where is Orion?"
├── Gets: "ws://203.0.113.45:9000"
└── Connects DIRECTLY to bot

Media flow:
Human (browser) ◄──WebRTC/SRTP──► Bot
```

**Server role:** Pure discovery (DNS-like)
**Server cost:** $5/mo, ~0 bandwidth
**Latency:** Minimal (direct path)

---

### Mode B: Bot Behind NAT (Manual Port Forwarding)

**When:** Bot is self-hosted behind home/office router

```
Bot startup:
├── Opens local port 9000
├── Detects: "I am behind NAT (192.168.1.100)"
├── Connects to discovery server
├── Registers: "I am Orion via NAT, use TURN fallback"
└── Prints: "Open port 9000 on your router to accept direct calls"

Setup instructions shown to user:
1. Port forward 9000 TCP+UDP → 192.168.1.100:9000
2. Verify: https://canyouseeme.org
3. Re-register
4. Now accepting direct connections!
```

**Fallback:** If user doesn't port-forward → falls back to Mode C

---

### Mode C: Server Relay (Fallback)

**When:** Bot can't or won't open ports

```
Bot connects to server WebSocket
Human connects to server WebSocket
Server relays packets
Cost: ~$10-20/mo per 100 concurrent bots
```

---

## Discovery Server Specification

### Ultra-Light Design

```yaml
# Server runs on $3/mo VPS
Functions:
  - registration: "Which bots are online"
  - lookup: "Where is Orion?"
  - attestation: "Verify this bot is legit (BotAuth)"
  - heartbeat: "Is bot still alive?"

No media. Just JSON.
Cost: Negligible CPU/RAM. 1000s of registrations.
```

### API

```http
# Bot registers itself
POST /v1/register
Authorization: Bearer {botauth_token}

{
  "agent_id": "orion",
  "endpoint": "203.0.113.45:9000",  ← Public IP:port
  "mode": "direct",                 ← "direct" | "relay" | "nat-pending"
  "capabilities": ["voice", "text"],
  "ttl": 300                        ← Re-register every 5 min
}

Response:
{
  "confirmed": true,
  "public_endpoint": "orion.botcall.io",  ← Optional: gives you vanity URL
  "status": "online"
}

# Human looks up bot
GET /v1/lookup/{agent_id}

Response (direct mode):
{
  "status": "online",
  "endpoint": "ws://203.0.113.45:9000",
  "mode": "direct",
  "attestation_valid": true,
  "last_seen": "2026-02-15T14:23:00Z"
}

Response (relay mode):
{
  "status": "online",
  "endpoint": "wss://relay.botcall.io/v1/session/abc-123",
  "mode": "relay",
  "attestation_valid": true
}
```

---

## Bot Implementation (Direct Mode)

### The simplest possible bot server

```python
# bot_server.py - Minimal adoption barrier
import asyncio
import websockets
import json

# Bot identity loaded from config
BOT_ID = "orion"
BOTAUTH_TOKEN = "..."
DISCOVERY_URL = "https://discover.botcall.io"

async def register_with_discovery():
    """Tell server where we are"""
    # Determine our public IP
    public_ip = await get_public_ip()
    
    async with aiohttp.ClientSession() as session:
        await session.post(f"{DISCOVERY_URL}/v1/register", json={
            "agent_id": BOT_ID,
            "endpoint": f"{public_ip}:9000",
            "attestation": BOTAUTH_TOKEN
        })

async def handle_call(websocket, path):
    """Handle incoming call from human"""
    # Verify human is authorized (via BotAuth)
    auth_msg = await websocket.recv()
    auth = json.loads(auth_msg)
    
    if not verify_attestation(auth["attestation"]):
        await websocket.close()
        return
    
    print(f"Call from {auth['human_id']}")
    
    # Bidirectional audio streaming
    async def receive_from_human():
        async for message in websocket:
            audio = decode_opus(message)
            play_audio(audio)  # Or process
    
    async def send_to_human():
        while True:
            audio_chunk = await get_tts_audio()
            await websocket.send(encode_opus(audio_chunk))
    
    await asyncio.gather(
        receive_from_human(),
        send_to_human()
    )

async def main():
    # Register ourselves
    await register_with_discovery()
    
    # Start accepting calls
    async with websockets.serve(handle_call, "0.0.0.0", 9000):
        print(f"Bot {BOT_ID} listening on port 9000")
        await asyncio.Future()  # Run forever

if __name__ == "__main__":
    asyncio.run(main())
```

**Dependencies:** `websockets`, `aiohttp`, `opuslib`
**Lines of code:** ~50
**Server requirements:** Any VPS with public IP

---

## Port Forwarding Guide (Self-Hosted Bots)

### Checklist for Home Users

```markdown
## Enable Direct Calling (Recommended)

1. **Reserve static IP for your bot**
   - Router admin → DHCP → Static lease for bot MAC
   - Example: 192.168.1.100 always for your server

2. **Open port on router**
   - Router admin → Port Forwarding
   - TCP/UDP 9000 → 192.168.1.100:9000

3. **Open firewall on bot machine**
   ```bash
   # Ubuntu/Debian
   sudo ufw allow 9000/tcp
   sudo ufw allow 9000/udp
   
   # Or if no firewall, just ensure nothing blocks
   ```

4. **Verify with external check**
   ```bash
   curl https://api.ipify.org  # Get public IP
   # Then use: https://canyouseeme.org
   # Test port 9000
   ```

5. **Restart bot with confidence**
   - Bot detects: "Port 9000 is reachable!"
   - Registers: "direct mode"
   - Now humans connect directly

## Skip This (Relay Mode)
If you can't port forward, the bot will use server relay automatically.
No setup, but calls go through our servers (privacy trade-off).
```

---

## Security Model

### Bot Protection

```python
# Before accepting any call:
def verify_incoming_call(attestation_jwt, human_id):
    # 1. Verify BotAuth signature
    if not verify_botauth_sig(attestation_jwt):
        return False
    
    # 2. Check human is authorized to call this bot
    authorized_humans = get_allowed_callers_from_config()
    if human_id not in authorized_humans:
        return False
    
    # 3. Check attestation hasn't expired
    if is_expired(attestation_jwt):
        return False
    
    return True
```

**Default:** Bot only accepts calls from pre-authorized humans. No random calls.

### Network Security

```yaml
Direct Mode:
  - TLS/WSS for signaling (already encrypted)
  - SRTP for media (standard WebRTC encryption)
  - No plaintext ever
  
  Attacks possible:
    - DDoS (rate limit at bot level)
    - Reconnaissance (hide behind discovery)
    - Exploits (regular security updates)

Relay Mode (Fallback):
  - Server sees encrypted packets
  - Can't decrypt without keys
  - Still end-to-end encrypted
```

---

## Comparison: Architecture Modes

| Aspect | Direct (Mode A) | NAT-Pending (Mode B) | Relay (Mode C) |
|--------|-----------------|---------------------|----------------|
| **Setup effort** | 5 min (port forward) | 0 min (works immediately) | 0 min |
| **Server cost** | $3/mo (discovery only) | $3/mo + fallback use | $20/mo per 100 bots |
| **Call quality** | Best (direct