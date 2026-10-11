# Job Application Answer Bank — Ari Quinones

Reference doc of reusable stories, facts, and drafted answers from job application work. Use this as context to draft new answers by tweaking tone/emphasis per company rather than starting from scratch.

---

## Core background facts (reuse across any application)

- **9 years of experience** as a software engineer, currently a de facto lead engineer.
- **GitHub:** https://github.com/ariquinones
- **Stack:** Go (core backend/systems work), React (UI work), Node.js (microservices), PostgreSQL (previous role), AWS (EventBridge, SQS, Lambda, DynamoDB, CloudWatch).
- **Python:** personal-project level only, not production experience — be upfront about this if a role requires it as a hard minimum.
- **AI tools used daily:** Claude Code (terminal, agentic/task-driven flows) and Cursor (GUI, GPT-5.6 and Grok as underlying models).
- **Signature agentic workflow:** write a detailed PRD → Claude Code breaks it into a discrete task list → run a **"ralph loop"** where a sub-agent completes each task and opens its own PR → every PR gets an automated review pass via **Greptile** (1-5 score) → personal manual review → end-to-end functional testing → CI must pass full test suite (hard regression gate) before merge.
- **Mentoring style:** whenever chosen as lead on a project, pairs with 1-2 other engineers — breaks the project into milestones, writes the technical design doc up front so everyone is aligned before code is written, assigns tasks, and walks through architecture/reasoning as needed.
- **Personal/career-origin story:** started career in healthcare — marketing/design intern at a physical therapy clinic, and also worked in a biotech lab doing pharmacogenomics — before moving into software.
- **Personal health story:** has worn glasses since age 4; eyesight bad enough to not be a LASIK candidate — lifelong condition, not something ever "solved."

### Flagship project #1 — Automations engine (event-driven workflow engine)
- Built and owned end-to-end: trigger/condition/action data model, scheduling + batch-execution engine (AWS EventBridge, SQS, Lambda), run tracking, extensibility layer that let third-party marketplace apps plug into the same system.
- **Metrics:** ~95% of customers actively use it, <1% of runs end in error, 5M+ automation runs processed in production.
- **Hardest/most technical problem:** exactly-once execution under concurrency — a recurring scheduled trigger, a fan-out of thousands of SQS-queued sub-runs, and retries could all be in-flight simultaneously, risking double-firing the same action (e.g., sending a file twice) for the same recipient. Solved by checking each sub-run against an **idempotency record maintained by another service** before marking it successful.
- **Key trade-offs made:**
  - Traded latency/an extra dependency for correctness (idempotency check) — worth it since a duplicate customer-facing action is worse than a slightly slower one.
  - Introduced a **partial-success run status** instead of binary success/failure, since reporting a batch as "failed" when 4,950 of 5,000 sub-runs succeeded was misleading — cost: more states to track/test.
  - Runs retry a bounded number of times before transitioning to **"skipped"** rather than retrying forever or failing silently.
- Zero-to-one build; this is the strongest "built something from nothing" example available.

### Flagship project #2 — Claude Code skill for staging environments ("manual work eliminated" story)
- Built a Claude Code skill automating a manual, multi-step staging test-workspace setup (test signup, plan upgrade, test payment method, credit grant overrides, ordered creation of two role-differentiated users).
- **Before:** ~25 minutes of manual setup, easy to get wrong.
- **After:** one command, ready-to-use credentials. Adopted by a team of 9 engineers, in daily use for 1+ month.
- Good "eliminated manual work" story when a role doesn't want the bigger Automations story.

### Flagship project #3 — Marketplace Apps API technical design document
- Wrote a design doc solving: how to extend a platform's first-class API model (already available to in-house apps) to apps built by external/third-party partners, without hand-building bespoke integrations per partner or giving external code unchecked access.
- Designed a schema-driven data model letting third-party apps register their own API surface (endpoints, input/output schemas, validation) so the platform could validate and proxy requests generically.
- Made explicit calls on error handling (wrapping third-party errors so nothing unexpected leaks to users) and named versioning/backward-compatibility as an open, unsolved risk.
- Broke the project into shippable phases due to a hard external dependency on a partner delivering their own API first.
- Note: original document is confidential/internal — only a sanitized/generalized description should ever be shared externally, never the raw doc.

### Flagship project #4 — Production LLM chatbot (CRM chatbot)
- Built and deployed a CRM chatbot using OpenAI's API.
- Processed incoming user messages through an **async queue** rather than synchronously, so a slow/delayed model response wouldn't block or freeze the UI.
- Monitored with **AWS CloudWatch**, tracking error rates to catch API failures/degraded behavior early.

### Flagship project #5 — Audit log system
- Built a system tracking every customer-facing action, exposed via a queryable API and UI.
- Useful as an adjacent proof point for security/compliance/traceability-minded roles (e.g., regulated industries) even without direct regulated-industry experience.

---

## Application-by-application record

### 1. Biobase — Senior Technical Product Engineer
**Job (1-2 sentences):** AI-native platform/database company. Application skips the résumé/cover-letter format entirely and asks for two specific things by email: something built that eliminated manual work (and what happened after), and something written (spec/RFC/design doc/postmortem) that made a messy problem precise.

- **Q: Something you built that eliminated manual work** → A: Claude Code skill for staging environments (see Flagship #2 above), framed with the context clause "My team runs pricing experiments that require spinning up test workspaces shaped like real customers..."
- **Q: Something you wrote (design doc)** → A: Marketplace Apps API TDD, sanitized (see Flagship #3), framed with "When I was brought on as lead engineer for a project extending our platform's API model to third-party marketplace apps..."
- Full email combined both under subject "Senior Technical Product Engineer — application," opening with a one-line personal intro.

---

### 2. Unnamed SWE application (AI tooling short-answer questions)
**Job (1-2 sentences):** General SWE role; no company link was used for this set of questions. Series of short (3-5 sentence) written answers about AI tooling and workflow.

- **Q: Compare & contrast your experience with various AI vendors' models and toolchains** → A: Claude Code (terminal, agent-first, autonomous task-driven work) vs. Cursor (GUI, model-agnostic — GPT-5.6 and Grok), framed around toolchain philosophy differences and when each is used.
- **Q: Describe your SDLC for dealing with the volume of code LLMs generate** → A: Greptile automated triage (1-5 score) → manual review → end-to-end testing → CI regression gate.
- **Q: Describe a project or feature you developed agentically** → A: Automations module, PRD → Claude Code task breakdown → ralph loop → per-task PR, staying focused on architecture calls (partial-success status, retry cap) while the agent handled implementation.
- **Q: Briefly describe your experience deploying LLM/GenAI applications to production** → A: CRM chatbot story (Flagship #4).

---

### 3. Deya Health — Software Engineer, AI x Healthcare (NYC/Hybrid/Remote)
**Job (1-2 sentences):** Eye-health startup building an AI-driven longitudinal eye-care data platform. Their thesis: the moat is a continuous patient data layer between visits, not another diagnostic/scribe tool.

- **Q: Why note (why interested + why good fit)** → A: personal LASIK-ineligibility story tied directly to their "continuous data layer" thesis, plus data-infra/agentic-workflow background as the fit argument.
- Separate discussion (not a written answer): whether a Loom video demo was safe given proprietary/credentialed product — recommendation was to use a staging/synthetic-data environment or a diagrams/code-walkthrough approach instead of production data, and to check with a manager first.

---

### 4. Ohr.xyz (role not restated, materials/chemistry-adjacent AI platform)
**Job (1-2 sentences):** AI platform company (molecule discovery / digital twins), values human-in-the-loop tools that augment domain experts rather than replace them; emphasizes commercial judgment and starting narrow.

- **Q: How would you use AI to reinvent a legacy industry? Where would you start?** → A: Physical therapy clinic story — the real recovery happens at home, but there's no continuous record between visits; wearables/pose-estimation now exist, but the real unsolved problem is closing the loop back into the therapist's next decision. Proposed starting narrow: one exercise category, phone-camera/sensor tracking, a therapist-facing summary before the next visit.
- **Q: What's the hardest problem you've solved, and what made it interesting?** → A: Idempotency/double-firing problem in the Automations engine (see Flagship #1's "hardest problem" section) — framed around getting the failure model right, not just the mechanics.

---

### 5. Relevant Healthcare — Lead Engineer
**Job (1-2 sentences):** Healthcare software company; playful/quirky posting (TypeScript joke, hidden "giraffe thermos" instruction to test attention to detail) looking for a lead engineer strong on mentorship, architecture ownership, and AI strategy, with an explicit interest in candidates who care about healthcare equity.

- **Full application email** (not a single Q&A) — required elements: subject line "Lead Engineer," resume attached, GitHub link, state where the posting was found (Hacker News "Who's Hiring"), and the exact phrase "giraffe thermos" worked in naturally. Content: personal healthcare-origin story (PT clinic + biotech pharmacogenomics lab), lead-engineer ownership examples, TDD-writing habit as evidence of mentorship-through-writing, agentic AI workflow experience.

---

### 6. Unnamed application B (technical depth / impact-focused written questions)
**Job (1-2 sentences):** General technical engineering role; no company link was used. Series of short written questions focused on technical depth, trade-offs, and impact — includes an "evidence of exceptional ability" style prompt aimed at filtering low-effort applications.

- **Q: Describe a technically complex problem you worked on... focus on what made it hard and the tradeoffs (3-5 sentences)** → A: Idempotency/double-firing problem, reframed explicitly around two named trade-offs: (1) latency/dependency vs. correctness, (2) simplicity vs. accurate partial-success failure reporting.
- **Q: What's the most meaningful impact you've had as an engineer? (3-5 sentences)** → A: Automations engine ownership story with full metrics (95% adoption, <1% error, 5M+ runs).
- **Q: Please provide evidence of exceptional ability (single paragraph)** → A: Broader paragraph anchored on Automations metrics, plus supporting mentions of lead-engineer ownership/TDD-writing and agentic AI workflow building, closing with an explicit "I take ambiguous, foundational problems and turn them into something the customer base depends on" framing.
- **Q: Explain your process for using Claude/ChatGPT to its full potential** → A: Full pipeline answer — PRD discipline → Claude Code task breakdown → ralph loop → Greptile → manual review → e2e testing → CI gate, framed as "the agent executes a well-specified plan; I stay focused on judgment."

---

### 7. AuthZed — (role building SpiceDB, a Zanzibar-inspired permissions database)
**Job (1-2 sentences):** Company building an open-source, distributed authorization/permissions database (SpiceDB). Role is distributed-systems- and database-heavy.

- **Q: Why AuthZed?** → A: Framed around genuine overlap in problem category — correctness-under-concurrency work (the idempotency/double-firing story) as the honest parallel to database correctness, plus real but "consumer-level" database experience (DynamoDB single-table design, index/access-pattern design, sharding awareness) — explicitly acknowledged as data-modeling-level experience, not database-internals/engine-building experience, framed as genuine motivation to go deeper.

---

### 8. Kamiwaza — Senior Software Engineer (MTS)
**Job (1-2 sentences):** Small, well-funded startup building GenAI deployment infrastructure for enterprises (model serving, embeddings/retrieval, data catalogs, evaluation pipelines, cluster orchestration). "MTS" = Member of Technical Staff, a flat-title convention. Notably, their posting explicitly invites disclosure of AI-assisted applications rather than treating it as a red flag.

- **Fit analysis:** Strong overlap on distributed systems, cloud infra, and API design; genuine overlap on GenAI-tool fluency (production chatbot + agentic dev workflow). **Real gap:** Python is a stated minimum requirement (production-grade services) but experience is personal-project level only — handled by disclosing this honestly and confidently rather than omitting it.
- **Application note** (not a single Q&A — covers their "how to apply" ask): GitHub link, a system design they're proud of (Automations engine, full story with metrics and the idempotency problem), an honest paragraph naming the Python gap as a fast-closing gap given systems fundamentals transfer, and an honest note that there's no ML-specific evaluation/benchmark work to share, offset by CI/testing rigor as the same underlying discipline.

---

### 9. Bicycle Health — (role not restated; healthcare/telemedicine, opioid use disorder treatment)
**Job (1-2 sentences):** Telemedicine company treating Opioid Use Disorder (OUD) via medication-assisted treatment, focused on affordability, confidentiality, and reducing stigma to improve care access.

- **Q: What about Bicycle Health's mission resonates most with you?** → A: Values-based answer — healthcare access for underserved populations, using tech to remove real barriers (cost, stigma, convenience) rather than optimizing for people who already have easy access.
- **Q: Do you regularly use AI coding tools? How do you maintain a verification process for catching subtle logic errors or security risks in generated output?** → A: Full verification pipeline (Greptile → manual review → e2e testing → CI gate), with an added closing line connecting that discipline specifically to healthcare stakes (patient data/PHI risk).

---

### 10. Truemed — (Labs team role; zero-to-one, high-autonomy environment)
**Job (1-2 sentences):** Company enabling people to use HSA/FSA funds on health-improving purchases; the specific team ("Labs") is a small, high-autonomy group doing hypothesis-driven prototyping, expected to use agentic AI workflows heavily, full-stack (Next.js/TypeScript frontend, Django/Python backend mentioned as desired stack).

- **Q: Briefly describe how your experience is a match for our position (1-3 sentences)** → A: Led with production LLM feature experience + daily agentic AI workflow use + track record of owning ambiguous/foundational problems end-to-end.
- **Q: Briefly describe why you are interested in Truemed (1-3 sentences)** → A: Genuine draw to the Labs operating model (small team, autonomy, agentic-workflow speed) plus appreciation for a prevention-focused (not just treatment-focused) mission.
- Note: Django/Python backend is listed as a desired skill here too — same honest caveat as Kamiwaza applies if it comes up live.

---

### 11. Prudentia Sciences — (role not restated; AI/LLM-orchestration platform for biopharma)
**Job (1-2 sentences):** Company building an AI-powered platform/LLM orchestration tooling for the biopharma industry.

- **Q: What is the most interesting way you have used AI professionally recently?** → A: The full agentic build loop (PRD → Claude Code task breakdown → ralph loop → per-task PR) used to build the Automations engine, with emphasis on where the person had to stay in the loop (the idempotency bug, the failure-reporting redesign) — framed as "AI as an execution layer you direct, not a code-writing tool."

---

### 12. Benchmark Analytics — Senior Software Engineer (2026)
**Job (1-2 sentences):** Public-safety/law-enforcement-adjacent analytics company. Posting asks for a direct email application (not a form) to a specific person — Noah Gartner, VP of Engineering — emphasizing full-stack ownership (Node.js/React/PostgreSQL), an AI-augmented SDLC, a "trio" mentoring/pairing model, and 7+ years experience.

- **Full personalized email to Noah Gartner** (final version): opens plainly and directly (no flattery), states interest specifically in the trio model and AI-augmented SDLC (with the real agentic-workflow pipeline as proof), then states fit at an overview level — full-stack ownership across React/Node.js/Postgres — anchored by the mentoring paragraph ("Whenever I've led a project, I've paired with one or two other engineers...").
- Lesson learned during drafting: when emailing a **named individual**, avoid a labeled/checklist structure matching requirements one-by-one (reads robotic) — write flowing prose, lead with genuine reason for interest before fit, and don't over-load it with project-level technical detail; save that for the interview.

---

## Other correspondence (not job-application answers, but related outreach)

- **Recruiter introduction email** — short intro email to a recruiter (via a warm intro from a former coworker), no specific role attached, resume mentioned as attached. Reusable template: name the mutual contact immediately, state current situation plainly (open to opportunities, nothing specific yet), keep it to 3-4 sentences.

---

## Patterns worth carrying into new applications

1. **Lead with real numbers when available** — 95% adoption / <1% error rate / 5M+ runs is the strongest single proof point across nearly every answer above. Reuse it whenever "impact" or "accomplishments" is asked.
2. **Name honest gaps rather than hide them** (Python level, no regulated-industry experience, database "consumer" vs. "builder" experience) — this reads as more credible than overclaiming, especially at senior level.
3. **The idempotency/double-firing story is the strongest "hard problem" story** — but it's been reused across multiple applications already (Ohr, AuthZed's parallel, Unnamed App B, Prudentia). If a single company asks multiple open-ended technical questions, have a second story ready so it doesn't repeat within one process.
4. **Personal stories are most powerful when concrete and specific**, not just stated as a fact — e.g., "not a LASIK candidate" (Deya) and "marketing/design intern at a PT clinic" (Ohr) both work because they're specific, not generic passion statements.
5. **When emailing a named individual** (vs. a generic inbox), drop the checklist/requirements-matching structure — write like a person addressing a person.
6. **Confidential/internal material** (TDDs, production systems with real customer data) should always be described/sanitized rather than shared directly — this came up multiple times (Biobase TDD, Deya Loom video).
