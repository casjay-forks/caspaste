## [ ] Create .claude/rules/*.md cheatsheet files (14 files)
Read: AI.md PART 0
PART 0 "Session Initialization" mandates creating
.claude/rules/{ai,project,config,binary,backend,api,frontend,features,
service,makefile,docker,cicd,testing,optional}-rules.md if the .claude/rules/
directory does not exist. It currently does not exist. Each file must follow
the required format (header with PART numbers, NON-NEGOTIABLE warning,
CRITICAL NEVER/ALWAYS sections, key rules summary, reference line) sourced
from the corresponding AI.md PARTs. This requires reading PARTs 7-36, which
was out of scope for this reconciliation pass.

## [ ] Verify README.md section order against PART 1 spec
Read: AI.md PART 1
PART 1 "README.md — Section Order" requires: Title & Badges, About,
Official Site, Features, Production, Client, Configuration, API, Other,
Development (always last). Current README.md order is: Install, Docker,
CLI Client, Server, API, Configuration, Development, License — does not
match prescribed order/section names. Not rewritten in this pass per
"prefer flagging over silently changing large working files"; also blocked
on the naming-conflict item above since README content itself may need to
change depending on that decision.

## [ ] Verify docs/ completeness against PART 3 tree
Read: AI.md PART 3
Required docs/ files per the structure tree: index.md, installation.md,
configuration.md, api.md, cli.md, admin.md, security.md, integrations.md,
development.md, stylesheets/dark.css (+optional light.css), requirements.txt.
Present: admin.md, api.md, cli.md, configuration.md, development.md, index.md,
installation.md, requirements.txt. Missing: security.md, integrations.md,
stylesheets/ (dark.css). PART 30 governs exact content requirements and was
out of scope for this pass (line range cutoff at PART 6).

## [ ] Verify LICENSE.md copyright year
Read: AI.md PART 2
PART 2 requires copyright year = "current year or year of first publication".
Copyright holder updated to "webappsgo" (= project_org) as part of the
caspaste/webappsgo naming-conflict resolution. LICENSE.md currently states
"Copyright (c) 2024 webappsgo". The year 2024 predates the first commit
visible in this git history (2026) — could be correct if the project was
first published elsewhere in 2024, or could be stale. Needs human
confirmation of the actual first-publication year; not changed in this pass.

## [ ] Verify src/mode + src/cli + src/tui runtime-mode dispatch against PART 6
Read: AI.md PART 6
src/mode/mode.go only implements a Production/Development app-mode toggle
(debug flag), not full smart-detection dispatch (server vs CLI vs TUI based
on TTY/args/flags) described in PART 6. Need to check src/cli and src/tui
entry points and src/main.go against the exact PART 6 dispatch rules — not
completed in this pass due to time/scope; flag only, do not implement.

