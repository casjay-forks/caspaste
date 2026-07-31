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

## [ ] Fix go-lint findings (14 issues, pre-existing, found during naming rename)
Read: AI.md PART 26 (Makefile), PART 8 (Server binary CLI)
go-lint agent found 14 violations unrelated to the caspaste/webappsgo rename
that landed in the same working tree; logged rather than fixed to keep the
rename commit scoped:
- Makefile line 37 (GO_DOCKER): missing `-e GOFLAGS=-buildvcs=false` for
  mounted .git directory safety
- Makefile lines 63, 66, 75, 81, 97, 100, 175, 178: `go build` calls missing
  `-buildvcs=false` flag
- Makefile line 25: LDFLAGS missing `-trimpath` (required for build, release,
  docker targets)
- Makefile: missing required `clean` target
- src/cli/cli.go: missing `-h`/`-v` short flag forms (only `--help`/
  `--version` work; both forms required)
- src/server/caspaste.go line ~1253: `--color` flag accepts wrong values
  (`always`/`never`/`auto`; spec requires `auto`/`yes`/`no`)
- src/server/caspaste.go line ~1253: `--color` flag defaults to empty string;
  must default to `auto`

## [ ] Rebuild stale CI `:build` toolchain image (go1.26.4 -> go1.26.5)
Read: AI.md PART 28
CI's `vuln-check` job runs `ghcr.io/{owner}/{repo}:build`, pinned to a Go
toolchain built before go1.26.5 (which fixes GO-2026-5856, the crypto/tls
Encrypted Client Hello privacy leak). `casjaysdev/go:latest` already ships
go1.26.5. Fixing requires triggering the "Build Toolchain Image"
workflow_dispatch to rebuild and push a fresh `:build` tag — not done here
since it changes a shared CI artifact outside this repo's normal push flow;
needs a human OK.
