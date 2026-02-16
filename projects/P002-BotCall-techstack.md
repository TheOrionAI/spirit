# BotCall: Tech Stack Analysis

**Language selection for performance-critical components**

---

## Discovery Server Language Comparison

### What the Discovery Server Actually Does

- HTTP API: `POST /register`, `GET /lookup/{id}`
- WebSocket: Real-time presence, heartbeats
- JSON parsing/serialization
- JWT signature verification (BotAuth)
- In-memory lookup (agent_id → endpoint mapping)

**Not CPU intensive, but connection-heavy.** Thousands of concurrent WebSockets.

---

## Option Analysis

### Python (Your original thought)

```python
# With asyncio + FastAPI + uvicorn
async def register_agent(agent_id, endpoint):
    await db.setex(f"agent:{agent_id}", 300, endpoint)
    return {"confirmed": True}

# Handles ~10k concurrent with uvloop
```

| Pros | Cons |
|------|------|
| Fast to prototype | GIL limits true parallelism |
| Asyncio works | Memory footprint ~50MB+ per 1k conns |
| Easy to read | Runtime overhead significant |
| Good ecosystem | Not designed for IO-heavy servers |

**Verdict:** Good for MVP, not for production scale.

---

### Go (Recommended)

```go
// Native goroutines, built-in HTTP/WebSocket
func handleRegister(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // Goroutine per connection
    go func() {
        store.Set(req.AgentID, req.Endpoint, 5*time.Minute)
    }()
    
    json.NewEncoder(w).Encode(RegisterResponse{Confirmed: true})
}

// Handles 100k+ concurrent connections easily
```

| Pros | Cons |
|------|------|
| Goroutines = lightweight threads (2KB stack) | Slightly more verbose than Python |
| Built-in HTTP/WebSocket (net/http) | New language to learn for contributors |
| Compiled = fast startup, efficient | Less "scripty" for quick changes |
| Static typing catches bugs early | | 
| Standard lib is excellent | |
| JSON parsing fast (encoding/json) | |
| JWT libraries robust (golang-jwt) | |

**Resource usage:**
- Memory: ~4-8KB per connection
- 10k connections: ~80MB RAM
- CPU: Negligible for JSON routing

**Verdict:** ✅ **Right choice for discovery server.**

---

### C / C++

```cpp
// epoll/kqueue, manual memory management
class DiscoveryServer {
public:
    void handleConnection(int fd) {
        // Zero-copy where possible
        // Manual buffer management
        // epoll for 100k+ fds
    }
};
// With libraries: libuv, websocket++, rapidjson
```

| Pros | Cons |
|------|------|
| Maximum performance | Development time 10x longer |
| Zero memory overhead | Memory safety (segfaults, leaks) |
| Predictable latency | Harder to maintain |
| Small binary | Build complexity (CMake, deps) |
| | Harder to find contributors |

**Resource usage:**
- Memory: ~2-4KB per connection possible
- 10k connections: ~40MB RAM  
- CPU: Same as Go for this workload

**Verdict:** Overkill for discovery. The complexity isn't worth it.

---

### Rust

```rust
// tokio async runtime
async fn handle_register(req: RegisterRequest) -> Result<impl Reply, Rejection> {
    let store = get_store();
    store.set(&req.agent_id, &req.endpoint, Duration::from_secs(300)).await;
    Ok(warp::reply::json(&RegisterResponse { confirmed: true }))
}

// Zero-cost abstractions, safe concurrency
```

| Pros | Cons |
|------|------|
| Memory safety without GC | Learning curve (borrow checker) |
| Performance = C++ | Compile times slower |
| Modern tooling (cargo) | Smaller ecosystem than Go |
| Async/await ergonomic | Harder to hire for |
| Fearless concurrency | |

**Verdict:** Great option, but Go is easier for this use case.

---

## Decision Matrix

| Criterion | Python | Go | C++ | Rust |
|-----------|--------|----|----|------|
| Performance | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Development velocity | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐ |
| Memory efficiency | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Concurrent connections | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Maintainability | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐ |
| Hiring/contributors | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐ |
| Deployment simplicity | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Overall for Discovery** | Good | **Best** | Overkill | Good |

---

## Recommended Stack

### Discovery Server: Go

```
Language: Go 1.21+
Framework: Standard library (net/http) + gorilla/websocket
JSON: encoding/json (std lib)
JWT: golang-jwt/jwt/v5
Rate limiting: golang.org/x/time/rate
Storage: Redis (for distributed) or sync.Map (single instance)
Metrics: Prometheus client

Why:
- 100k+ concurrent connections on cheap VPS
- Compiles to single binary, deploy anywhere
- Static binary, no dependencies
- Easy to dockerize
```

**Single binary deployment:**
```bash
go build -o botcall-discovery
./botcall-discovery  # Listens on :8080
```

---

### Bot SDK: Multiple Options

| Language | Priority | Why |
|----------|----------|-----|
| **Go** | High | Same as server, easy deployment |
| **Python** | High | Most AI/ML devs know it |
| **Rust** | Medium | Performance, safety |
| **Node.js** | Medium | Web developers |

**Bot SDK minimum spec:**
```
Functions:
  - register_with_discovery()
  - accept_call()
  - send_audio_opus()
  - receive_audio_opus()
  - verify_human_attestation()

Footprint: < 50MB RAM, < 10% CPU idle
```

---

### Human Client: PWA (Progressive Web App)

**Yes, PWA is the right call.**

**Architecture:**
```
PWA Structure:
├── index.html        ← Shell
├── manifest.json     ← "Install to home screen"
├── sw.js             ← Service worker (offline, push)
├── app/
│   ├── main.js       ← Svelte/React/Vanilla
│   ├── webrtc.js     ← RTCPeerConnection handling
│   └── opus-decoder.js  ← Audio pipeline
└── assets/
    └── icons/
```

**Why PWA over native app:**

| Aspect | PWA | Native App |
|--------|-----|------------|
| Install friction | Low (bookmark) | High (app store) |
| Update | Instant | Review process |
| Platform coverage | All (iOS/Android/Desktop) | Per-platform |
| WebRTC access | ✅ Full support | Same |
| Push notifications | ✅ iOS 16.4+, Android always | Same |
| Offline capability | ✅ Service worker | Same |
| Camera/mic | ✅ Permission prompt | Same |

**PWA Tech Stack:**
```
Framework: SvelteKit (lightweight, fast) or Preact (tiny)
Bundler: Vite (instant HMR)
WebRTC: Native browser APIs
State: Svelte stores or Zustand
Styling: Tailwind CSS
Icons: Heroicons
Audio: Web Audio API + Opus (via wasm-ogg-opus)

Build output:
- index.html (2KB)
- app.js (50KB gzipped)
- manifest.json
- Total: < 100KB for shell
```

**iOS considerations:**
```javascript
// iOS requires user interaction for audio context
const audioContext = new (window.AudioContext || window.webkitAudioContext)();
// Must resume() after user gesture

// Push notifications: iOS 16.4+ supports PWAs on home screen
// Before that: Safari only, no push
```

---

## Architecture Summary

```
┌────────────────────────────────────────────────────────────┐
│                      BotCall Network                       │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐  │
│  │    Go       │     │    Go       │     │    PWA      │  │
│  │ Discovery   │◄───►│  Bot SDK    │◄───►│   Human     │  │
│  │   Server    │     │  (Optimal)  │     │   Client    │  │
│  │             │     │             │     │             │  │
│  │ - Register  │     │ - Accept    │     │ - WebRTC    │  │
│  │ - Lookup    │     │ - Stream    │     │ - UI        │  │
│  │ - Auth      │     │ - Decode    │     │ - Connect   │  │
│  │             │     │             │     │             │  │
│  └─────────────┘     └─────────────┘     └─────────────┘  │
│                                                            │
│  $3/m