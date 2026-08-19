# Releasing

All five SDKs are versioned and released independently. Each has its own
GitHub Actions workflow, triggered by its own tag prefix:

| SDK | Tag pattern | Workflow | Publishes to |
|---|---|---|---|
| TypeScript | `typescript-vX.Y.Z` | [`release-typescript.yml`](.github/workflows/release-typescript.yml) | [npm](https://www.npmjs.com/package/tachyon-sdk) |
| Python | `python-vX.Y.Z` | [`release-python.yml`](.github/workflows/release-python.yml) | [PyPI](https://pypi.org/project/tachyon-sdk/) |
| Go | `go/vX.Y.Z` (slash, not hyphen — see below) | [`release-go.yml`](.github/workflows/release-go.yml) | nowhere — the git tag *is* the release |
| C# | `csharp-vX.Y.Z` | [`release-csharp.yml`](.github/workflows/release-csharp.yml) | [NuGet](https://www.nuget.org/packages/TachyonSdk) |
| Java | `java-vX.Y.Z` | [`release-java.yml`](.github/workflows/release-java.yml) | [Maven Central](https://central.sonatype.com/artifact/io.github.adikeshri/tachyon-sdk) |

The TypeScript, Python, and C# workflows all follow the same shape: build
once, run the full unit + integration suite (against a real
`adikeshri/tachyon` container) on that build, verify the tag matches the
version in the package manifest, then promote the *exact same* built
artifact to `publish` — gated by an environment approval — which publishes
and cuts a GitHub release. Go and Java depart from that shape for reasons
specific to their ecosystems, covered in their own sections below.

## Cutting a release

1. **Bump the version** in the relevant manifest and open a normal PR
   (skip this step for Go — it has no manifest version, see below):
   - TypeScript: [`typescript/package.json`](typescript/package.json) (`version` field)
   - Python: [`python/pyproject.toml`](python/pyproject.toml) (`[project].version`)
   - C#: [`csharp/Tachyon.Sdk/Tachyon.Sdk.csproj`](csharp/Tachyon.Sdk/Tachyon.Sdk.csproj) (`<Version>`)
   - Java: [`java/pom.xml`](java/pom.xml) (`<version>`)
2. **Merge to `main`.**
3. **Tag the merge commit** and push the tag:

   ```bash
   git checkout main && git pull
   git tag typescript-v0.2.0   # or python-v0.2.0, csharp-v0.2.0, java-v0.2.0, go/v0.2.0
   git push origin typescript-v0.2.0
   ```

4. The matching release workflow picks up the tag, runs the full pipeline,
   and — if everything passes and the environment approval (if configured)
   goes through — publishes and creates a GitHub release with
   auto-generated notes scoped to that SDK's own tag history.

For TypeScript/Python/C#/Java, if the tag's version doesn't match the
manifest's version, the workflow fails fast at the "Determine version and
verify it matches the tag" step before anything is built for publishing.

## Dry runs

Every workflow except Go's also accepts `workflow_dispatch` from the
Actions tab, on any branch, with a `dry_run` input (defaults to `true`).
This runs the entire pipeline — install, typecheck/mypy, unit tests,
integration tests against a live Tachyon container, build, package sanity
checks — without publishing or creating a release. Use it to validate a
workflow change or a not-yet-tagged package before committing to a real
release.

Set `dry_run` to `false` on a manual run only if you deliberately want to
publish without a corresponding tag push (e.g. re-running a release that
failed after tagging). The `publish` job is still gated by its environment
either way.

Go's `workflow_dispatch` has no `dry_run` input — see the Go section below
for why.

## Go: the tag *is* the release

Go modules aren't published to a registry the way the other four are —
`go get` fetches straight from this repo via the Go module proxy
(`proxy.golang.org`), keyed entirely off the git tag. The moment
`go/vX.Y.Z` is pushed, anyone can fetch it; there's no separate publish
action to gate, run approval on, or undo.

Two things follow from that:

- **The tag format is different on purpose.** Go requires
  `<subdirectory>/vX.Y.Z` for a module that lives in a repo subdirectory
  (ours is `github.com/adikeshri/tachyon-sdk/go`) — see [the Go modules
  reference](https://go.dev/ref/mod#vcs-version). `go-v0.2.0` would not
  work; it has to be `go/v0.2.0`.
- **The real gate is upstream of the workflow**: only tag a commit where CI
  has already passed on `main`. `release-go.yml` still runs the full test
  suite against the tagged commit and creates a GitHub release, but that
  happens *after* the release already exists from Go's perspective.

There's also no version field to bump anywhere — skip step 1 above for Go
releases.

## One-time repository setup

None of the following can be configured from this repo's files alone.

### GitHub Environments

Create these under Settings → Environments: `npm-release`, `pypi-release`,
`nuget-release`, `maven-central-release`. Adding required reviewers to each
turns every publish into a manual-approval step; recommended, since none of
these registries allow overwriting a published version. (Go has no
environment — see above.)

### npm and PyPI

Pick one per registry:

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

### NuGet

Uses [NuGet Trusted Publishing](https://learn.microsoft.com/en-us/nuget/nuget-org/trusted-publishing)
(OIDC) — no long-lived API key stored anywhere. `release-csharp.yml`'s
`publish` job exchanges its GitHub Actions identity token for a NuGet API
key that's valid for exactly one hour, via the `NuGet/login` action, right
before pushing.

One-time setup on nuget.org:

1. Sign in → your username → **Trusted Publishing** → add a new policy
2. **Repository Owner**: `adikeshri`, **Repository**: `tachyon-sdk`,
   **Workflow File**: `release-csharp.yml` (just the filename, not the
   `.github/workflows/` path), **Environment**: `nuget-release` (must match
   this repo's `environment: nuget-release` exactly, or the token exchange
   is rejected)
3. If `TachyonSdk` isn't published yet, the policy starts in a 7-day
   provisional state (nuget.org can't lock it to a specific repo ID until a
   real publish happens) — the first tagged release needs to land inside
   that window, or you restart the timer from the same policy page.

Then set one secret, `NUGET_USER`, on the `nuget-release` environment —
your nuget.org **profile name** (not email, not a key) — the `NuGet/login`
step passes it as the identity being authenticated as.

### Maven Central

The heaviest one-time setup of the five, and none of it can be done from
CI:

1. **Claim the namespace.** Sign up at
   [central.sonatype.com](https://central.sonatype.com), then verify the
   `io.github.adikeshri` namespace — for a `io.github.*` groupId this is
   automatic (tied to owning the `adikeshri` GitHub account), no DNS TXT
   record needed.
2. **Generate a user token** (Account → Generate User Token on Central
   Portal — this is a token pair, not your login password). Set the two
   values as `CENTRAL_USERNAME` / `CENTRAL_PASSWORD` secrets on the
   `maven-central-release` environment.
3. **Generate a GPG key pair** — Central requires every release artifact to
   be signed:
   ```bash
   gpg --full-generate-key
   gpg --keyserver keyserver.ubuntu.com --send-keys <KEY_ID>
   gpg --armor --export-secret-keys <KEY_ID>   # paste this whole block into the secret below
   ```
   Set the exported private key block as the `GPG_PRIVATE_KEY` secret, and
   the key's passphrase as `GPG_PASSPHRASE`, both on the
   `maven-central-release` environment.

`release-java.yml` uses `actions/setup-java`'s built-in GPG/Maven-settings
support to wire all four secrets into the build without ever writing them
to disk outside the runner's temp keyring. `java/pom.xml`'s `release`
[Maven profile](java/pom.xml) — source jar, javadoc jar, GPG signing, and
the `central-publishing-maven-plugin` — only activates with `-P release`,
so plain `mvn test`/`mvn verify` (CI, local dev) never needs a GPG key
present.

`autoPublish` is enabled in that plugin config, meaning the GitHub
environment approval is the *only* human gate — once approved, the release
goes live on Central without an extra manual click in the Central Portal
UI (Central's own validation still runs and can still fail the deploy on
malformed metadata).

## CI on every push/PR

[`ci.yml`](.github/workflows/ci.yml) runs on every push to `main` and every
PR: unit tests across a matrix of Node/Python/Go versions (C# and Java run
against whatever's current, since neither has multiple actively-diverging
major versions worth matrixing here), plus a separate integration job per
SDK against a live `adikeshri/tachyon` service container, plus a packaging
sanity check per language (`npm publish --dry-run`,
`python -m build && twine check`, `dotnet pack`, and equivalent build
verification for Go/Java). Release workflows re-run all of this against
the actual release commit rather than trusting a prior CI run, so a release
is never gated on a check that ran against a different commit.
