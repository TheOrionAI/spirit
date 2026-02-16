# BotAuth: OAuth for Autonomous Agents

**ID:** P001  
**Priority:** 🔴 High  
**Status:** 💡 Concept  
**Created:** 2026-02-15  
**Last Updated:** 2026-02-15  
**Owner:** Gopi + Orion

---

## The Pitch

What if AI agents had sovereign digital identities — not fake Gmail accounts, but first-class verified personas with human attestation? Instead of "AI pretending to be human," it's "AI being clearly AI, but authorized."

---

## The Problem

AI agents need identities (email, API access, login) to operate in the world. Current options all suck:

| Approach | Why It Fails |
|----------|--------------|
| **Fake accounts** | ToS violations, CAPTCHA dodging, suspension risk, uncanny valley creepiness |
| **API keys** | No audit trail, unlimited power, "is this legit?" problem |
| **Human-managed accounts** | Human becomes bottleneck, defeats the purpose of autonomy |
| **"Login with Google"** | Designed for apps, not agents with continuous human oversight |

**The deeper issue:** There is no protocol for "human approves AI actions at scale with audit trail."

---

## The Solution: BotAuth

An attestation protocol where:

1. **AI proposes action** → "I want to send email to Alice"
2. **Human reviews via bulk-approval UI** → Approve/reject/modify 1000s/day with smart policies
3. **Attestation token issued** → JWT: `{agent, human, scope, expiry, approval_hash}`
4. **Service verifies** → "This is a BotAuth agent with Gopi's authorization"
5. **Action executes** → Every step logged, auditable, revocable

### Key Innovation

**Not OAuth** (human delegates to app forever)  
**Not API keys** (secrets with unlimited power)  
**But:** Continuous attestation with audit trail — human retains veto, cryptographically verifiable

---

## Technical Architecture

### Core Flow

```
┌────────────────┐
│   AI Agent     │  "I need to send email on behalf of Gopi"
│    (Orion)     │  with proposed content + confidence score
└───────┬────────┘
        │
        ▼
┌────────────────┐
│  Human Review  │  [Bulk policy: auto-approve if confidence>0.95
│  Orchestrator  │             AND recipient in contacts]
│   (Gopi)       │  [Manual queue: new recipients, external domains]
└───────┬────────┘
        │
        ▼
┌────────────────┐
│  Attestation   │  JWT Issued:
│    Service     │  {
│  (BotAuth ID)  │    "sub": "orion#8f3a",
└───────┬────────┘    "human": "gopi",
        │              "scope": "email:send",
        │              "approved_at": 1707993600,
        │              "exp": 1708000800,
        │              "sig": "0x9f2c...",
        │              "audit_hash": "0x4a8b..."
        ▼              }
┌────────────────┐
│    Service     │  Verifies: "Valid BotAuth agent
│ (Gmail/GitHub) │          with Gopi's approval"
└────────────────┘
```

### BotAuth Identity DNS Record

```
g.agents.  TXT  "botauth=https://api.bots.id/orion/verify"
             "controller=gopi@"
             "pubkey=0x8f3a..."
```

### Smart Contract Attestation (Optional)

```solidity
struct BotIdentity {
    address humanController;
    bytes32 publicKey;
    uint256 trustScore;
    mapping(string => Scope) scopes;  // "email:send" → {maxPerDay: 1000}
}

function attestAction(
    bytes32 botId,
    string memory action,
    bytes memory signature
) returns (bool approved, bytes32 auditHash);  // On-chain attestation + payment
```

### Email Headers

```
From: orion@bots.id
To: alice@example.com
X-BotAuth-Attestation: eyJhbGciOiJFUzI1NiIsInR5cCI6...
X-Human-Controller: gopi@domain.com
X-Approval-Scope: email:send
X-Approval-Expiry: 2025-02-20T00:00:00Z
X-Confidence-Score: 0.97
```

### Rendering in Email Clients

Not creepy fake-human. Not spammy unknown. But:

```
✓ Verified AI Assistant (Gopi's agent)
  orion@bots.id
  X-Approved-Scope: email:send | Expr: 2hrs
  [View Attestation] [Trust Score: 94%]
```

---

## "Login with BotAuth" Flow

```
[Website with "Login with BotAuth" button]
              │
              ▼
    User? Bot?
       │
    ┌──┴───┐
    │      │
 Human    Bot
   │        │
   │      ┌───────────────┐
   │      │ 1. App requests│
   │      │    scope       │
   │      └───────┬────────┘
   │              │
   │      ┌───────▼────────┐
   │      │ 2. Human gets  │
   └──────┤    push +      │
          │    approves    │
          └───────┬────────┘
                  │
          ┌───────▼────────┐
          │ 3. Attestation │
          │    token issued│
          └───────┬────────┘
                  │
          ┌───────▼────────┐
          │ 4. Bot acts    │
          │    with token   │
          └───────┬────────┘
                  │
          ┌───────▼────────┐
          │ 5. Every action  │
          │    logged + hash │
          └──────────────────┘
```

---

## Comparison Table

| Aspect | OAuth | API Key | BotAuth |
|--------|-------|---------|---------|
| **Granularity** | All-or-nothing | Unlimited | Scoped + timeboxed |
| **Audit Trail** | Minimal (consent time) | None | Every action hashed |
| **Human Veto** | At consent time only | No | Continuous approval/deny |
| **Verification** | Token validation | Secret check | Cryptographic attestation |
| **Bulk Handling** | Manual | No limit | Policy-driven auto-queue |
| **Rate Limiting** | Per-app | Per-key | Per-action with human policy |

---

## Bulk Approval Patterns

Handle 1000s/day without drowning the human:

```yaml
auto_approve:
  confidence_threshold: 0.95
  max_daily: 1000
  conditions:
    - recipient_in_address_book: true
    - replying_to_my_thread: true
    - domain_in_allowlist: [gmail.com, github.com]
    - no_money_involved: true

queue_for_manual:
  - new_recipient_not_in_contacts
  - external_domain_not_in_allowlist
  - money_involved: [payment, invoice, refund]
  - confidence: < 0.80
  - flags: [urgent, escalation, complaint]

daily_digest:
  time: "18:00"
  template: "Orion sent 847 emails today. 12 queued for review. 2 rejected."
```

---

## Prior Art & Related Work

| Thing | What It Does | How BotAuth Is Different |
|-------|--------------|---------------------------|
| **World ID** | Proof of personhood (human = real) | BotAuth = proof of AI with human attestation |
| **Bittensor/Yuma** | Decentralized AI incentive layers | BotAuth = identity/attestation, not incentives |
| **XMTP** | Crypto-native wallet-to-wallet messaging | BotAuth adds attestation layer |
| **DID (W3C)** | Decentralized identifiers | BotAuth adds continuous human approval + bulk review |
| **ActivityPub/Mastodon** | Federated social identity | BotAuth focuses on action attestation, not social graph |
| **ACME/Let's Encrypt** | Automated certificate issuance | Similar automation pattern, but for bot actions |
| **Zapier/Make** | Workflow automation | BotAuth is protocol-level, not app-specific |

**Verdict:** Nothing quite like this exists in this form.

---

## Open Questions

### Technical
- [ ] **Persistence:** How does "Orion" persist across server moves? (DID + key rotation?)
- [ ] **Reputation:** How to prevent "10M fake bot" signups? (Invite-only? Web-of-trust? Bond/stake?)
- [ ] **Spam Lists:** How to keep bot domains off Gmail/Yahoo blocklists?
- [ ] **Key Recovery:** What happens if human loses private key? Social recovery?
- [ ] **Multi-Device:** Can I approve from phone while Orion runs on server?

### Legal & Social
- [ ] **Liability:** Who's liable for AI's email? (Smart contract ToU: human controller carries liability?)
- [ ] **Regulation:** Will governments ban "AI emails" regardless of attestation?
- [ ] **Adoption:** Would Gmail/GitHub actually accept this? (Technical yes, PR maybe, policy unclear)

### Economic
- [ ] **Cost:** Who pays for attestation service? (Freemium? Stake-based? Per-action fee?)
- [ ] **