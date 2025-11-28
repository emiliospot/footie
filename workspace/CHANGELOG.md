## 1.0.0 (2025-11-28)

### Features

* add comprehensive typecheck commands and fix type inconsistencies ([df15515](https://github.com/emiliospot/footie/commit/df15515e1120f435caffdaaafe3db17fa2058aff))
* Add database migration scripts to root package.json ([7472828](https://github.com/emiliospot/footie/commit/7472828b6263d16d5aeacbab306edbfb6875c6c2))
* Add markdown linting with markdownlint-cli2 ([58b0b7b](https://github.com/emiliospot/footie/commit/58b0b7b1660f2be568ed63f5f48d2a189082b779))
* Add real-time WebSocket + Redis Streams architecture ([08f4808](https://github.com/emiliospot/footie/commit/08f48087d8c2da8c237c0578f31fb6ffe9d371ed))
* Change API port from 8081 to 8088 and fix Angular dev server security warning ([c8295ad](https://github.com/emiliospot/footie/commit/c8295ada086e9f4619224470787a06753dfdc1cf))
* implement Competition Rankings feature ([25dbacc](https://github.com/emiliospot/footie/commit/25dbacc5fecac00c6fb292522e77c938c0e2bfa6))
* implement extensible webhook provider pattern and event type system ([e63ee8a](https://github.com/emiliospot/footie/commit/e63ee8a644fb6f7a1acc1917fb7f6419a765fd36))
* Implement sqlc-based handlers and Redis Commander ([32915f0](https://github.com/emiliospot/footie/commit/32915f0a6c83b1cc9caf5e2a318a80ecac330142))
* Initial professional full-stack monorepo setup ([16bd300](https://github.com/emiliospot/footie/commit/16bd30063190527450ec18c908cf41aade02c0d4))
* Migrate to sqlc + pgx + golang-migrate (pro analytics stack) ([61012d0](https://github.com/emiliospot/footie/commit/61012d01ce1dcdee691f6d5044ff501e78f01c4d))
* **tooling:** Add conventional-changelog for automated CHANGELOG generation ([63564d4](https://github.com/emiliospot/footie/commit/63564d4e141067aaa69983976acd810fe3142b4b))
* update Competition Rankings component with layout improvements ([89ab112](https://github.com/emiliospot/footie/commit/89ab1126c60da424fefcc40c14b933fe92169745))

### Bug Fixes

* add ExtractEvents method to OptaProvider and update SQL schema ([704083f](https://github.com/emiliospot/footie/commit/704083fd8c608716619ac99c543f5e64b44ee229))
* **ci:** Skip E2E tests until Phase 1 handlers are implemented ([8ff057b](https://github.com/emiliospot/footie/commit/8ff057b1e14ecc239924a5e23a356fbe63096654))
* Clean up GORM dependencies and fix type errors after sqlc migration ([e6bc1dc](https://github.com/emiliospot/footie/commit/e6bc1dc61d5138758055023bb312b036268f2621))
* Correct database URL format for golang-migrate ([f257ae7](https://github.com/emiliospot/footie/commit/f257ae74819a5897353ab4c897077556db57c899))
* Correct migrations path to ./migrations ([3e6fd9b](https://github.com/emiliospot/footie/commit/3e6fd9b431fb26e3344c5aef15e755cab93476cd))
* **vscode:** Add Docker kill task and unignore .vscode for team consistency ([679a134](https://github.com/emiliospot/footie/commit/679a134187264206172717a14afd6cb63a2f30cc))
## 1.0.0 (2025-11-28)

### Features

* Add database migration scripts to root package.json ([7472828](https://github.com/emiliospot/footie/commit/7472828b6263d16d5aeacbab306edbfb6875c6c2))
* Add markdown linting with markdownlint-cli2 ([58b0b7b](https://github.com/emiliospot/footie/commit/58b0b7b1660f2be568ed63f5f48d2a189082b779))
* Add real-time WebSocket + Redis Streams architecture ([08f4808](https://github.com/emiliospot/footie/commit/08f48087d8c2da8c237c0578f31fb6ffe9d371ed))
* Change API port from 8081 to 8088 and fix Angular dev server security warning ([c8295ad](https://github.com/emiliospot/footie/commit/c8295ada086e9f4619224470787a06753dfdc1cf))
* implement Competition Rankings feature ([25dbacc](https://github.com/emiliospot/footie/commit/25dbacc5fecac00c6fb292522e77c938c0e2bfa6))
* implement extensible webhook provider pattern and event type system ([e63ee8a](https://github.com/emiliospot/footie/commit/e63ee8a644fb6f7a1acc1917fb7f6419a765fd36))
* Implement sqlc-based handlers and Redis Commander ([32915f0](https://github.com/emiliospot/footie/commit/32915f0a6c83b1cc9caf5e2a318a80ecac330142))
* Initial professional full-stack monorepo setup ([16bd300](https://github.com/emiliospot/footie/commit/16bd30063190527450ec18c908cf41aade02c0d4))
* Migrate to sqlc + pgx + golang-migrate (pro analytics stack) ([61012d0](https://github.com/emiliospot/footie/commit/61012d01ce1dcdee691f6d5044ff501e78f01c4d))
* **tooling:** Add conventional-changelog for automated CHANGELOG generation ([63564d4](https://github.com/emiliospot/footie/commit/63564d4e141067aaa69983976acd810fe3142b4b))
* update Competition Rankings component with layout improvements ([89ab112](https://github.com/emiliospot/footie/commit/89ab1126c60da424fefcc40c14b933fe92169745))

### Bug Fixes

* add ExtractEvents method to OptaProvider and update SQL schema ([704083f](https://github.com/emiliospot/footie/commit/704083fd8c608716619ac99c543f5e64b44ee229))
* **ci:** Skip E2E tests until Phase 1 handlers are implemented ([8ff057b](https://github.com/emiliospot/footie/commit/8ff057b1e14ecc239924a5e23a356fbe63096654))
* Clean up GORM dependencies and fix type errors after sqlc migration ([e6bc1dc](https://github.com/emiliospot/footie/commit/e6bc1dc61d5138758055023bb312b036268f2621))
* Correct database URL format for golang-migrate ([f257ae7](https://github.com/emiliospot/footie/commit/f257ae74819a5897353ab4c897077556db57c899))
* Correct migrations path to ./migrations ([3e6fd9b](https://github.com/emiliospot/footie/commit/3e6fd9b431fb26e3344c5aef15e755cab93476cd))
* **vscode:** Add Docker kill task and unignore .vscode for team consistency ([679a134](https://github.com/emiliospot/footie/commit/679a134187264206172717a14afd6cb63a2f30cc))
## 1.0.0 (2025-11-28)

### Features

* Add database migration scripts to root package.json ([7472828](https://github.com/emiliospot/footie/commit/7472828b6263d16d5aeacbab306edbfb6875c6c2))
* Add markdown linting with markdownlint-cli2 ([58b0b7b](https://github.com/emiliospot/footie/commit/58b0b7b1660f2be568ed63f5f48d2a189082b779))
* Add real-time WebSocket + Redis Streams architecture ([08f4808](https://github.com/emiliospot/footie/commit/08f48087d8c2da8c237c0578f31fb6ffe9d371ed))
* Change API port from 8081 to 8088 and fix Angular dev server security warning ([c8295ad](https://github.com/emiliospot/footie/commit/c8295ada086e9f4619224470787a06753dfdc1cf))
* implement Competition Rankings feature ([25dbacc](https://github.com/emiliospot/footie/commit/25dbacc5fecac00c6fb292522e77c938c0e2bfa6))
* implement extensible webhook provider pattern and event type system ([e63ee8a](https://github.com/emiliospot/footie/commit/e63ee8a644fb6f7a1acc1917fb7f6419a765fd36))
* Implement sqlc-based handlers and Redis Commander ([32915f0](https://github.com/emiliospot/footie/commit/32915f0a6c83b1cc9caf5e2a318a80ecac330142))
* Initial professional full-stack monorepo setup ([16bd300](https://github.com/emiliospot/footie/commit/16bd30063190527450ec18c908cf41aade02c0d4))
* Migrate to sqlc + pgx + golang-migrate (pro analytics stack) ([61012d0](https://github.com/emiliospot/footie/commit/61012d01ce1dcdee691f6d5044ff501e78f01c4d))
* **tooling:** Add conventional-changelog for automated CHANGELOG generation ([63564d4](https://github.com/emiliospot/footie/commit/63564d4e141067aaa69983976acd810fe3142b4b))

### Bug Fixes

* add ExtractEvents method to OptaProvider and update SQL schema ([704083f](https://github.com/emiliospot/footie/commit/704083fd8c608716619ac99c543f5e64b44ee229))
* **ci:** Skip E2E tests until Phase 1 handlers are implemented ([8ff057b](https://github.com/emiliospot/footie/commit/8ff057b1e14ecc239924a5e23a356fbe63096654))
* Clean up GORM dependencies and fix type errors after sqlc migration ([e6bc1dc](https://github.com/emiliospot/footie/commit/e6bc1dc61d5138758055023bb312b036268f2621))
* Correct database URL format for golang-migrate ([f257ae7](https://github.com/emiliospot/footie/commit/f257ae74819a5897353ab4c897077556db57c899))
* Correct migrations path to ./migrations ([3e6fd9b](https://github.com/emiliospot/footie/commit/3e6fd9b431fb26e3344c5aef15e755cab93476cd))
* **vscode:** Add Docker kill task and unignore .vscode for team consistency ([679a134](https://github.com/emiliospot/footie/commit/679a134187264206172717a14afd6cb63a2f30cc))
## 1.0.0 (2025-11-28)

### Features

* Add database migration scripts to root package.json ([7472828](https://github.com/emiliospot/footie/commit/7472828b6263d16d5aeacbab306edbfb6875c6c2))
* Add markdown linting with markdownlint-cli2 ([58b0b7b](https://github.com/emiliospot/footie/commit/58b0b7b1660f2be568ed63f5f48d2a189082b779))
* Add real-time WebSocket + Redis Streams architecture ([08f4808](https://github.com/emiliospot/footie/commit/08f48087d8c2da8c237c0578f31fb6ffe9d371ed))
* Change API port from 8081 to 8088 and fix Angular dev server security warning ([c8295ad](https://github.com/emiliospot/footie/commit/c8295ada086e9f4619224470787a06753dfdc1cf))
* implement Competition Rankings feature ([25dbacc](https://github.com/emiliospot/footie/commit/25dbacc5fecac00c6fb292522e77c938c0e2bfa6))
* implement extensible webhook provider pattern and event type system ([e63ee8a](https://github.com/emiliospot/footie/commit/e63ee8a644fb6f7a1acc1917fb7f6419a765fd36))
* Implement sqlc-based handlers and Redis Commander ([32915f0](https://github.com/emiliospot/footie/commit/32915f0a6c83b1cc9caf5e2a318a80ecac330142))
* Initial professional full-stack monorepo setup ([16bd300](https://github.com/emiliospot/footie/commit/16bd30063190527450ec18c908cf41aade02c0d4))
* Migrate to sqlc + pgx + golang-migrate (pro analytics stack) ([61012d0](https://github.com/emiliospot/footie/commit/61012d01ce1dcdee691f6d5044ff501e78f01c4d))
* **tooling:** Add conventional-changelog for automated CHANGELOG generation ([63564d4](https://github.com/emiliospot/footie/commit/63564d4e141067aaa69983976acd810fe3142b4b))

### Bug Fixes

* add ExtractEvents method to OptaProvider and update SQL schema ([704083f](https://github.com/emiliospot/footie/commit/704083fd8c608716619ac99c543f5e64b44ee229))
* **ci:** Skip E2E tests until Phase 1 handlers are implemented ([8ff057b](https://github.com/emiliospot/footie/commit/8ff057b1e14ecc239924a5e23a356fbe63096654))
* Clean up GORM dependencies and fix type errors after sqlc migration ([e6bc1dc](https://github.com/emiliospot/footie/commit/e6bc1dc61d5138758055023bb312b036268f2621))
* Correct database URL format for golang-migrate ([f257ae7](https://github.com/emiliospot/footie/commit/f257ae74819a5897353ab4c897077556db57c899))
* Correct migrations path to ./migrations ([3e6fd9b](https://github.com/emiliospot/footie/commit/3e6fd9b431fb26e3344c5aef15e755cab93476cd))
* **vscode:** Add Docker kill task and unignore .vscode for team consistency ([679a134](https://github.com/emiliospot/footie/commit/679a134187264206172717a14afd6cb63a2f30cc))
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
- **api:** Add Competition Rankings API endpoint with mock data
- **api:** Make database and Redis optional in development mode
- **web:** Add Competition Rankings component (Angular)
- **web:** Embed Competition Rankings in dashboard
- **web:** Add rankings service and models
- **web:** Update auth guard to redirect to dashboard instead of login
- **docs:** Add comprehensive event types documentation
- **docs:** Add duration and period handling documentation
- **docs:** Add Competition Rankings feature README
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
