# Agent skills

This repository ships on-demand agent skills for the stack. They are
committed under `.agents/skills/`. `.claude/skills/` and `.windsurf/skills/`
contain relative symlinks to the same directories.

`skills-lock.json` records the source repository, skill path and content hash
for the vendor skills. The project skills are local and are not in the lock
file.

Skills are instructions, not code. They do not run by themselves. The agent
loads one when its description matches the task. Review a skill before you
trust it, because a skill runs with the same permissions as the agent.

## Inventory

| Stack            | Source                     | Skills                                                                                                                                                    |
| ---------------- | -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Go               | `samber/cc-skills-golang`  | `golang-code-style`, `golang-context`, `golang-database`, `golang-design-patterns`, `golang-documentation`, `golang-error-handling`, `golang-lint`, `golang-naming`, `golang-project-layout`, `golang-security`, `golang-testing` |
| shadcn/ui        | `shadcn/ui`                | `shadcn`                                                                                                                                                  |
| Docker           | `docker/skills`            | `docker-project-foundations`, `docker-build-strategies`, `docker-compose-patterns`, `docker-docs`                                                         |
| Playwright       | `microsoft/playwright-cli` | `playwright-cli`                                                                                                                                          |
| This project     | local                      | `grb-architecture`, `grb-auth-rbac`, `grb-api-contract`, `grb-postgres-migrations`, `grb-release-pipeline`, `grb-i18n-themes`, `grb-testing`               |

## Project skills

| Skill                      | Use it when                                                     |
| -------------------------- | --------------------------------------------------------------- |
| `grb-architecture`         | Adding a feature or deciding where code belongs                 |
| `grb-auth-rbac`            | Changing sessions, 2FA, approval, users, roles or permissions   |
| `grb-api-contract`         | Adding an endpoint or a React data fetch                        |
| `grb-postgres-migrations`  | Changing the schema, seed, gorm store or migration             |
| `grb-release-pipeline`     | Changing Vite, prerendering, go:embed, Docker or Compose        |
| `grb-i18n-themes`          | Changing messages, locale behavior or color palettes           |
| `grb-testing`              | Adding or running Go, Vitest or Playwright tests                |

These skills point at the project files. They do not replace `AGENTS.md`.
They go deeper on one workflow.

## Deliberately not shipped

The Next.js sibling also ships `vercel/next.js`, `prisma/skills` and
`better-auth/skills`. Those do not apply to this stack:

- The server is Go, not Next.js.
- The data layer is gorm and SQL migrations, not Prisma.
- Auth is app-owned, not Better Auth.

Do not add those skills to this repository. Use the Go, shadcn, Docker and
Playwright sets, plus the project skills.

## Stack items without an official vendor skill

- Tailwind CSS — use the `shadcn` skill for Tailwind 4 and semantic tokens.
- PostgreSQL and gorm — use `grb-postgres-migrations` and `golang-database`.
- Vite, React Router, TanStack Query, Zod, React Hook Form, Vitest,
  Sonner and Recharts — use the project skills and the codebase itself.
- Go testing and linting — use the `golang-testing` and `golang-lint` skills.

Do not add a community skill only to fill a gap. If a community skill is
useful, review it, pin it in `skills-lock.json` and document why it is here.

## Commands

```bash
npx skills list                 # installed skills
npx skills update               # update vendor skills to their latest versions
npx skills experimental_install # restore vendor skills from skills-lock.json
```

Install a specific skill into the canonical `.agents` layout:

```bash
npx skills add shadcn/ui --skill shadcn --copy -a opencode -y

npx skills add docker/skills \
  --skill docker-project-foundations \
  --skill docker-build-strategies \
  --skill docker-compose-patterns \
  --skill docker-docs \
  --copy -a opencode -y

npx skills add microsoft/playwright-cli \
  --skill playwright-cli \
  --copy -a opencode -y

npx skills add samber/cc-skills-golang \
  --skill golang-code-style \
  --skill golang-context \
  --skill golang-database \
  --skill golang-design-patterns \
  --skill golang-documentation \
  --skill golang-error-handling \
  --skill golang-lint \
  --skill golang-naming \
  --skill golang-project-layout \
  --skill golang-security \
  --skill golang-testing \
  --copy -a opencode -y
```

After a vendor install or update, refresh the agent links:

```bash
mkdir -p .claude/skills .windsurf/skills
for d in .agents/skills/*/; do
  name=$(basename "$d")
  ln -sfn "../../.agents/skills/$name" ".claude/skills/$name"
  ln -sfn "../../.agents/skills/$name" ".windsurf/skills/$name"
done
```

## Add a project skill

1. Create `.agents/skills/grb-<name>/SKILL.md`.
2. Add the YAML frontmatter with `name` and `description`.
3. Keep it short and action-oriented.
4. Point to the project files it protects.
5. Add a relative symlink in `.claude/skills` and `.windsurf/skills`.
6. Add it to the inventory in this file.

Do not duplicate the rules in `AGENTS.md`. Go deeper on one workflow.

## Security

The installer printed a risk summary for each vendor skill. The Go set, the
Docker set and Playwright are widely used. The shadcn skill is maintained by
the shadcn project. Review any skill before you install it. A skill can ask
the agent to run commands.

The app also has no skill that can bypass the server-side auth and RBAC
checks. The checks live in Go, not in the agent instructions.

## Licenses

| Source                    | License    |
| ------------------------- | ---------- |
| `samber/cc-skills-golang` | MIT        |
| `shadcn/ui`               | MIT        |
| `docker/skills`           | Apache-2.0 |
| `microsoft/playwright-cli`| Apache-2.0 |
| Project skills            | MIT        |

The license of each vendor skill is also in its `SKILL.md` frontmatter where
the upstream project provides one.
