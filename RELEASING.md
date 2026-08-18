# Releasing

The two SDKs are versioned and released independently. Each has its own
GitHub Actions workflow, triggered by its own tag prefix:

| SDK | Tag pattern | Workflow | Publishes to |
|---|---|---|---|
| TypeScript | `typescript-vX.Y.Z` | [`release-typescript.yml`](.github/workflows/release-typescript.yml) | [npm](https://www.npmjs.com/package/tachyon-sdk) |
| Python | `python-vX.Y.Z` | [`release-python.yml`](.github/workflows/release-python.yml) | [PyPI](https://pypi.org/project/tachyon-sdk/) |

Both workflows follow the same shape: build once, run the full unit +
integration suite (against a real `adikeshri/tachyon` container) on that
build, verify the tag matches the version in the package manifest, then
promote the *exact same* built artifact to `publish` — which requires the
`npm-release` / `pypi-release` environment to approve, and only then
publishes and cuts a GitHub release.

## Cutting a release

1. **Bump the version** in the relevant manifest and open a normal PR:
   - TypeScript: [`typescript/package.json`](typescript/package.json) (`version` field)
   - Python: [`python/pyproject.toml`](python/pyproject.toml) (`[project].version`)
2. **Merge to `main`.**
3. **Tag the merge commit** and push the tag:

   ```bash
   git checkout main && git pull
   git tag typescript-v0.2.0   # or python-v0.2.0
   git push origin typescript-v0.2.0
   ```

4. The matching release workflow picks up the tag, runs the full pipeline,
   and — if everything passes and the environment approval (if configured)
   goes through — publishes and creates a GitHub release with
   auto-generated notes scoped to that SDK's own tag history.

If the tag's version doesn't match the manifest's version, the workflow
fails fast at the "Determine version and verify it matches the tag" step
before anything is built for publishing.

## Dry runs

Both workflows also accept `workflow_dispatch` from the Actions tab, on any
branch, with a `dry_run` input (defaults to `true`). This runs the entire
pipeline — install, typecheck/mypy, unit tests, integration tests against a
live Tachyon container, build, package sanity checks — without publishing
or creating a release. Use it to validate a workflow change or a
not-yet-tagged package before committing to a real release.

Set `dry_run` to `false` on a manual run only if you deliberately want to
publish without a corresponding tag push (e.g. re-running a release that
failed after tagging). The `publish` job is still gated by its environment
either way.

## One-time repository setup

The workflows expect the following, none of which can be configured from
this repo's files alone:

**GitHub Environments** (Settings → Environments) — create `npm-release`
and `pypi-release`. Adding required reviewers to each turns every publish
into a manual-approval step; this is optional but recommended, since
neither npm nor PyPI allow overwriting a published version.

**Publishing credentials** — pick one per registry:

- *Trusted publishing (recommended, no long-lived secret)*: configure the
  registry to trust this repo's release workflow directly —
  [npm Trusted Publishers](https://docs.npmjs.com/trusted-publishers) for
  the `tachyon-sdk` package, or
  [PyPI Trusted Publishers](https://docs.pypi.org/trusted-publishers/) for
  the `tachyon-sdk` project. Leave the token secrets below unset; the
  workflows already request `id-token: write` and fall back to OIDC
  automatically when no token is provided.
- *API token*: set repository (or environment) secrets `NPM_TOKEN` (an npm
  [automation
  token](https://docs.npmjs.com/creating-and-viewing-access-tokens)) and/or
  `PYPI_API_TOKEN` (a PyPI [API
  token](https://pypi.org/help/#apitoken) scoped to the `tachyon-sdk`
  project).

A brand-new PyPI project can't have a trusted publisher configured before
it exists under that name; publish the very first release with an API
token, then switch to trusted publishing for subsequent ones if desired.

## CI on every push/PR

[`ci.yml`](.github/workflows/ci.yml) runs on every push to `main` and every
PR: unit tests across a matrix of Node/Python versions, plus a separate
integration job per SDK against a live `adikeshri/tachyon` service
container, plus a packaging sanity check (`npm publish --dry-run`,
`python -m build && twine check`). Release workflows re-run all of this
against the actual release commit rather than trusting a prior CI run, so a
release is never gated on a check that ran against a different commit.
