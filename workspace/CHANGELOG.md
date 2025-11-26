## [Unreleased]

### Features

- **api:** Decouple domain models from GORM and sqlc dependencies
- **api:** Add mapper layer to convert between sqlc types and domain models
- **api:** Refactor match handler to use clean domain models in API responses
- **api:** Implement extensible webhook provider pattern (Adapter + Strategy + Registry)
- **api:** Add support for multiple webhook providers (Generic, Opta, StatsBomb)
- **api:** Add batch event processing support for webhook endpoints
- **api:** Implement event type normalization and validation system
- **api:** Add period handling (first half, second half, extra time, penalties)
- **api:** Add second-level precision to match events (0-59 seconds)
- **api:** Add per-provider webhook secret configuration
- **api:** Add database migration for event time fields (second, period)
- **docs:** Add comprehensive event types documentation
- **docs:** Add duration and period handling documentation
- **docs:** Update architecture documentation with provider pattern
- **docs:** Update tech stack presentation with new patterns

### Refactoring

- **api:** Remove GORM dependencies from all domain models (user, team, player, match, match_event, statistics)
- **api:** Replace `uint` with `int32` in domain models to align with sqlc
- **api:** Replace `gorm.DeletedAt` with `*time.Time` for soft deletes
- **api:** Extract async match event publishing to reduce cognitive complexity
- **api:** Refactor webhook handler to use provider abstraction layer
- **api:** Fix all linting warnings (variable shadowing, error handling, comments)

### Bug Fixes

- **api:** Fix pgtype.Date conversion in player mapper
- **api:** Fix pgtype.Numeric conversions in match event handler
- **api:** Improve error handling in position coordinate scanning

## 1.0.0 (2025-11-20)

### Features

- Add markdown linting with markdownlint-cli2 ([58b0b7b](https://github.com/emiliospot/footie/commit/58b0b7b1660f2be568ed63f5f48d2a189082b779))
- Initial professional full-stack monorepo setup ([16bd300](https://github.com/emiliospot/footie/commit/16bd30063190527450ec18c908cf41aade02c0d4))
- **tooling:** Add conventional-changelog for automated CHANGELOG generation ([63564d4](https://github.com/emiliospot/footie/commit/63564d4e141067aaa69983976acd810fe3142b4b))

### Bug Fixes

- **vscode:** Add Docker kill task and unignore .vscode for team consistency ([679a134](https://github.com/emiliospot/footie/commit/679a134187264206172717a14afd6cb63a2f30cc))
