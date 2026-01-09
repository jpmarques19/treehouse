# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.4](https://github.com/jpmarques19/treehouse/compare/v0.4.3...v0.4.4) (2026-01-09)


### Features

* **workflows:** enhance crew-add with progressive persona development ([390ab55](https://github.com/jpmarques19/treehouse/commit/390ab55e08ce08e65b47fd50ff96d49db82486f7))
* **workflows:** integrate crew-add with th crew --add CLI ([acfba13](https://github.com/jpmarques19/treehouse/commit/acfba13334dddc769873e984151da7378c65f680))


### Code Refactoring

* **workflows:** simplify crew-add to generate-first pattern ([b649100](https://github.com/jpmarques19/treehouse/commit/b6491001b5e4098ff6cc9019cca87226929fdeeb))
* **workflows:** transform crew-add to intent-driven discovery ([56e8afc](https://github.com/jpmarques19/treehouse/commit/56e8afcab343d950830cc267646662d957844cc0))

## [0.4.3](https://github.com/jpmarques19/treehouse/compare/v0.4.2...v0.4.3) (2026-01-09)


### Features

* **cli:** add th crew --add command for agent creation ([0bb7a62](https://github.com/jpmarques19/treehouse/commit/0bb7a6255a77dffd95a44ce95e652a75ca938b4f))

## [0.4.2](https://github.com/jpmarques19/treehouse/compare/v0.4.1...v0.4.2) (2026-01-08)


### Features

* **prune:** implement th prune command for orphan cleanup ([a2676f0](https://github.com/jpmarques19/treehouse/commit/a2676f03ceb1ee726aa5067292766d1afcd3d7b9))
* **remove:** implement th remove command for nook cleanup ([e8f7180](https://github.com/jpmarques19/treehouse/commit/e8f7180ac4c36d5671246c9a267307cd65e43a8a))
* **workflow:** add delete and prune menu actions to treehouse-list ([24d8b73](https://github.com/jpmarques19/treehouse/commit/24d8b73ed31637b96c590f2a91a8d574c4d28fb1))
* **workflows:** add generic agent loader stub ([845b221](https://github.com/jpmarques19/treehouse/commit/845b221fad9fb24d74614bbaa15ac6d4ff0496c9))
* **workflows:** add handoff and crew-add workflows for Epic 5 & 6 ([4cf2da8](https://github.com/jpmarques19/treehouse/commit/4cf2da8ef67b20f84a26106386c2ac34fcf71411))


### Bug Fixes

* address PR review security and code quality issues ([b88114b](https://github.com/jpmarques19/treehouse/commit/b88114ba8d4dba8a1847d725a5852e75063fdd99))


### Code Refactoring

* **workflows:** rename agent stub to agent-loader ([3996099](https://github.com/jpmarques19/treehouse/commit/3996099fc12ab27274b3095a25c79edddf23189e))
* **workflows:** rename handoff workflow to checkpoint ([b3fe163](https://github.com/jpmarques19/treehouse/commit/b3fe163fdebe45c1064ece9f6249a305f234ef4e))
* **workflows:** simplify Claude stub loader pattern ([0476c05](https://github.com/jpmarques19/treehouse/commit/0476c05d34dd3703a48f93a80c991dbacbefe4e2))

## [0.4.1](https://github.com/jpmarques19/treehouse/compare/v0.4.0...v0.4.1) (2026-01-08)


### Features

* add CLI commands for version and workspace init ([c832522](https://github.com/jpmarques19/treehouse/commit/c83252295352efd42b048c549b24486234deb444))
* add CLI entry point ([5656dcf](https://github.com/jpmarques19/treehouse/commit/5656dcfb3b16ef72811b8d00aec47bebde058e35))
* add git repository detection ([bef1917](https://github.com/jpmarques19/treehouse/commit/bef191735a13f2bb6a2f58362afa9afdbce5432d))
* add GitHub Actions release workflow ([7d90f9d](https://github.com/jpmarques19/treehouse/commit/7d90f9d97baeb578eff97cba5af5b44cc2748f43))
* add GoReleaser configuration for cross-compilation ([bec2165](https://github.com/jpmarques19/treehouse/commit/bec2165840ace19c7c4957fd209e1d13b3ef7d60))
* add initial changelog ([c357a45](https://github.com/jpmarques19/treehouse/commit/c357a4582d222bd45009961f318e672500c425d2))
* add install script for binary distribution ([a80a5fe](https://github.com/jpmarques19/treehouse/commit/a80a5fe7bfdc02923849deedbdb538856452d7f0))
* add JSON output infrastructure ([1915a6a](https://github.com/jpmarques19/treehouse/commit/1915a6a697c944dec85b2e1eae15ffdf6212ebbe))
* add nook-context-analyst agent for contextual lineage analysis in BMAD tracking system ([64b5e21](https://github.com/jpmarques19/treehouse/commit/64b5e2173e19f04886717d906930a425c20a08a0))
* add Release Please for semantic versioning ([81f4616](https://github.com/jpmarques19/treehouse/commit/81f46166f033392ff40e28dbb22e8944399f8048))
* add shared test utilities for git operations ([333c309](https://github.com/jpmarques19/treehouse/commit/333c309369fefd615239a3863519109481c3b298))
* **fork:** implement th fork command for nook creation ([fe540d4](https://github.com/jpmarques19/treehouse/commit/fe540d456780b1ac967ba8179cf487d685a11819))
* **init:** embed and install workflows during th init ([64c0e07](https://github.com/jpmarques19/treehouse/commit/64c0e07ab06f792c9b1a4cf4e7457993ceaabeae))
* initialize Go module with Cobra dependency ([36040f1](https://github.com/jpmarques19/treehouse/commit/36040f149b69c61cecaeb87e62748b45d93887da))
* **list:** implement th list command for workspace visibility ([a33563d](https://github.com/jpmarques19/treehouse/commit/a33563d527131ca8347eee6d5e750744b102f4ef))
* **nook-fork:** add agent validation step before deployment ([b137494](https://github.com/jpmarques19/treehouse/commit/b137494de42771c0d57684957f6389eb3f04e58d))
* **nook-fork:** add Agent Wizard with lineage context integration ([6ea8ae4](https://github.com/jpmarques19/treehouse/commit/6ea8ae46e49a1335234c37c89c227d085aef55dd))
* **nook-fork:** refactor to step-file architecture v3.0 ([8c4e452](https://github.com/jpmarques19/treehouse/commit/8c4e4523c2fe67fa485bb4a42be944a2e2d9157a))
* **nook:** implement nook ID generation and deck lineage model ([cb7d17d](https://github.com/jpmarques19/treehouse/commit/cb7d17dae38011cbdebfceb59ea57ffb0578737d))
* **prune:** implement th prune command for orphan cleanup ([a2676f0](https://github.com/jpmarques19/treehouse/commit/a2676f03ceb1ee726aa5067292766d1afcd3d7b9))
* **remove:** implement th remove command for nook cleanup ([e8f7180](https://github.com/jpmarques19/treehouse/commit/e8f7180ac4c36d5671246c9a267307cd65e43a8a))
* **th:** add skip_worktree_paths config option ([2af4557](https://github.com/jpmarques19/treehouse/commit/2af45576d8fa8508fd3b1239f1fafd5aef441e2a))
* **th:** migrate all workflows to step-engine architecture ([35ee431](https://github.com/jpmarques19/treehouse/commit/35ee4311b9c4535847ab374873e1d602ff2d5d03))
* **workflow:** add delete and prune menu actions to treehouse-list ([24d8b73](https://github.com/jpmarques19/treehouse/commit/24d8b73ed31637b96c590f2a91a8d574c4d28fb1))
* **workflows:** add generic agent loader stub ([845b221](https://github.com/jpmarques19/treehouse/commit/845b221fad9fb24d74614bbaa15ac6d4ff0496c9))
* **workflows:** add handoff and crew-add workflows for Epic 5 & 6 ([4cf2da8](https://github.com/jpmarques19/treehouse/commit/4cf2da8ef67b20f84a26106386c2ac34fcf71411))


### Bug Fixes

* address PR review security and code quality issues ([b88114b](https://github.com/jpmarques19/treehouse/commit/b88114ba8d4dba8a1847d725a5852e75063fdd99))
* **config:** simplify sync paths to use bmad/ directory ([bf99f7e](https://github.com/jpmarques19/treehouse/commit/bf99f7e5603c10775d53c0f5ef6a5db7a2c9785d))
* correct gitignore pattern for th binary ([63ec919](https://github.com/jpmarques19/treehouse/commit/63ec9196e5de2a72bfee55a763cf293742b355cf))
* **fork:** improve error handling and optimize deck loading ([f89fe5d](https://github.com/jpmarques19/treehouse/commit/f89fe5d1d43573c8037edc28bae249f1c93f9567))
* **fork:** store nook worktrees inside .treehouse/nooks/ ([1955874](https://github.com/jpmarques19/treehouse/commit/19558747de78850c4f4e1afb420fe0b260dba07a))
* **install:** use GitHub raw URL instead of non-existent domain ([14f932f](https://github.com/jpmarques19/treehouse/commit/14f932feed773f4270abc02749d2b4ddf96e0954))
* **list:** detect base repo when running from nook worktree ([66423bf](https://github.com/jpmarques19/treehouse/commit/66423bfb8289b994a216aaf5ca1f189053fdc5d0))
* **nook-fork:** enforce correct agent output path in YOLO mode ([160808e](https://github.com/jpmarques19/treehouse/commit/160808e3679be5e4ebc93981c91e3df3adf22f6c))
* **nook-fork:** update custom agent path to bmad/agents/ ([9301d16](https://github.com/jpmarques19/treehouse/commit/9301d166ffdf645163464a6486778e66dbd51235))
* resolve .gitignore conflict with main branch ([bed070c](https://github.com/jpmarques19/treehouse/commit/bed070cad8054feaf385c9513253fc4af4719580))
* update module path to match repository ([1560bf1](https://github.com/jpmarques19/treehouse/commit/1560bf19777c67270e5b43ec4ccee336b255f5fa))


### Code Refactoring

* decouple treehouse workflows from BMAD ([c85573f](https://github.com/jpmarques19/treehouse/commit/c85573fb45ec7a0d337559043810c45883c78995))
* migrate workflows to standalone format ([0294013](https://github.com/jpmarques19/treehouse/commit/0294013ab2518f7ffe0768d1b724787c50a6b165))
* **workflows:** rename agent stub to agent-loader ([3996099](https://github.com/jpmarques19/treehouse/commit/3996099fc12ab27274b3095a25c79edddf23189e))
* **workflows:** rename handoff workflow to checkpoint ([b3fe163](https://github.com/jpmarques19/treehouse/commit/b3fe163fdebe45c1064ece9f6249a305f234ef4e))
* **workflows:** simplify Claude stub loader pattern ([0476c05](https://github.com/jpmarques19/treehouse/commit/0476c05d34dd3703a48f93a80c991dbacbefe4e2))

## [Unreleased]

### Features

- Initial release of Treehouse CLI (`th`)
- `th init` command to initialize treehouse workspace
- `/treehouse-init` workflow for Claude Code integration
- JSON-only output for all CLI commands
- Git repository detection and validation
- Git version check (requires 2.5+ for worktree support)
