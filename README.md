# Release tests

This branch contains the changes applied to an upstream checkout. It does not
contain another maintained copy of actionlint. Edit these patches; regenerate
`tests` instead of editing that branch.

| Branch                      | Purpose                                                               |
| --------------------------- | --------------------------------------------------------------------- |
| `master`                    | Exact upstream `master` snapshot.                                     |
| Original stack branch names | Exact copies of the six upstream PR heads.                            |
| `patches`                   | Ordered patches, reusable workflow, historical capture, instructions. |
| `tests`                     | Generated source; default branch for running the workflows.           |
| `test-release`              | Preserved execution history from before this restructuring.           |

## Prepare source

Run **Prepare tests** on `tests`, supplying the source repository, source ref, and
an unused stable version. The source ref can be a branch, tag, or commit; no stack
tip is hardcoded. The patches must apply to the selected source.

```sh
gh workflow run prepare-tests.yml --repo kajaux/actionlint --ref tests \
  -f source-repository=kjanat/actionlint -f source-ref=YOUR_REF -f version=YOUR_VERSION
```

The caller invokes this branch's reusable workflow. The workflow retrieves its
actual patch SHA from GitHub's **run-attempt** metadata, checks out that revision,
verifies `series.sha256`, and applies `series.txt` in order. Conflicts stop the run
before the remote `tests` branch changes. There is no fallback or three-way merge.

The patched version helper updates all declared version references. It uses the
source commit's calendar date for the changelog, not the runner's current date.
Git author/committer dates also come from the source commit. The generated commits
are rooted in that source SHA and name the patch revision and selected version.
The `tests-source` artifact records inputs, both resolved SHAs, date, commit, and
resulting tree. Identical inputs and patch revision produce identical trees.

The final push updates only `tests`, with a lease against the head observed at
the start. It does not create a tag, draft, release, or package. No standalone
preparation script is required.

GitHub can reject `GITHUB_TOKEN` pushes that introduce changed workflow files.
That failure leaves `tests` unchanged; it requires a workflow-capable credential
before the run can publish those changes. A successful push with unchanged
workflow content does not establish that permission.

[Verification](verification.json) records the exact source and patch revisions,
matching branch SHAs, historical trees, successful generation, and an intentional
conflict run that left `tests` unchanged.

## Review and maintain patches

- `001`: publication isolation targeting `kajaux/actionlint` directly; no npmjs,
  Docker Hub, Homebrew, Scoop, AUR, WinGet, or Nix execution.
- `002`: existing additional published Action tests and GHCR destination fix.
- `003`: `tests` branch assumptions, deterministic preparation-only version mode,
  and matching helper tests.
- `004`: thin workflow-dispatch caller on the generated branch.

Publication is restricted to `kajaux/actionlint`, GitHub Packages (`@kajaux`),
and `ghcr.io/kajaux/actionlint`. The existing facade package name remains
`@kajaux/actionlint-rehearsal`; this restructuring does not rename published packages.

Update the patch files and their hashes together. To produce an amended patch,
edit an expendable checkout at that patch's input tree and export
`git diff --binary --full-index`; never make unrecorded edits on `tests`.
Upstream changes absorbed into a newer source may require removing or revising a
patch. That is an explicit reviewed edit, not a generator heuristic.

## Release execution is separate

After inspecting generated `tests`, explicitly dispatch `release-prepare.yml`
on that branch with the version and expected generated commit. Promotion still
uses `scripts/release-candidate.mjs promote` from a clean `tests` checkout with
`--repo kajaux/actionlint`. It verifies the prepared source, signed tag, assets,
and exact workflow attempt. Existing immutable tags are never replaced.

The preparation workflow does not invoke either operation. The generated
publication workflows retain their candidate and published-consumer checks.

## Historical capture

`history/manifest.json` identifies the original upstream source and each modified
tree. Apply its two patches in order to that source, using `git apply --index`,
then compare `git write-tree` after each patch with the recorded tree. This
reconstructs the old setup exactly, including its old repository name and fixed
version bump. Those historical files are **not** the active series.

The old execution commits remain on `test-release`, and existing release tags
remain untouched. This repository contains no changes to `kjanat/actionlint`.
