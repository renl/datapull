# Project Agent Directory

> Master index of all active documentation artifacts. Maintained by `@chief-of-staff`.

## Active Artifacts

| Artifact | Location | Owner | Status |
|----------|----------|-------|--------|
| Requirements | `docs/requirements/` | `@product-manager` | Not created |
| UX/UI Specification | `docs/ux-ui-spec/` | `@ux-ui-designer` | Not created |
| Architecture Contract | `docs/architecture/` | `@software-architect` | Drafted |
| Test Plan | `docs/test-plan/` | `@sdet-engineer` | Drafted; app-level verification blocked |
| Business Plan | `docs/business-plan.md` | `@business-strategist` | Not created |
| Pitch Deck | `docs/pitch-deck.md` | `@business-strategist` | Not created |
| Messaging Playbook | `docs/messaging-playbook.md` | `@growth-marketer` | Not created |
| Sales Deck | `docs/sales-deck.md` | `@growth-marketer` | Not created |
| Social Media Playbook | `docs/social-media-playbook.md` | `@growth-marketer` | Not created |
| Operations Directory | `docs/operations/` | `@devops-engineer` | Not created |
| Planning Directory | `docs/planning/` | `@planner` | Not created |
| Status Log | `docs/status_log.md` | `@chief-of-staff` | Active |

## Notes

- This repo has begun bootstrap through the URA quick-restore slice; architecture and test-plan artifacts currently exist in this worktree.
- Standalone URA fetch verification passed in this worktree, but direct app-level Refresh verification is currently blocked in automation because the Fyne `GLFW30` window exposes no UI Automation descendants.
- For directory-based artifacts, always read the directory's `README.md` first once it exists.
