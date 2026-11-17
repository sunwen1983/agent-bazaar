# Demo LLM key (D4-b) — placeholder only, no real key

Decision (D4-b): the live site's default path seeds a **demo DeepSeek key**
into the visitor's browser (BYOK profile in localStorage), so judges can
open the site and run it with zero friction. BYOK stays available for
visitors who prefer their own key.

## How it is funded (to be done at deploy time, by Wen)

- Wen creates a **separate DeepSeek account/key for the demo only**, tops
  it up with a small prepaid balance (**~$5 cap**), and rotates the key
  after the contest ends.
- The key is public by design (anyone can read it from the page source or
  localStorage). The only hard limit is the prepaid balance: a leaked key
  can burn at most the balance on that account. Never use a production key
  or an account with auto-recharge.
- Cost estimate: one full team run is ~10–20 short calls; at DeepSeek-V3
  prices (~$0.3–$1 per million tokens) $5 covers hundreds of demos.

## How it is wired (M2 work, not in this skeleton)

- The bazaar `site.js` seeds a `Demo` BYOK profile (localStorage key
  `nsweb-byok-profiles-v1`) with:
  - `base_url`: `https://api.deepseek.com/v1`
  - `model`: `deepseek-chat`
  - `api_key`: the demo key (placeholder below until Wen fills it in at
    deploy time — **never commit the real key**).
- The key never goes into the manifest or any network definition (the spec
  forbids it there). It only lives in the visitor's browser localStorage,
  exactly like a BYOK key.
- The demo video uses Wen's own key via BYOK (typed off-camera), not the
  demo key, so the video shows real model behavior.

## Placeholder

```
DEMO_DEEPSEEK_API_KEY=<fill in at deploy time, never commit>
DEMO_DEEPSEEK_MODEL=deepseek-chat
DEMO_DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
```

## Fallback

If Wen decides against any key exposure, delete the seeding code and the
site is pure BYOK (judges paste their own key). That is a worse judge
experience and Devpost's Testing clause expects the project to be testable
"free of charge and without any restriction", which is why (b) is the
default.
