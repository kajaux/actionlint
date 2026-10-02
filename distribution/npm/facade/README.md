# @kajaux/actionlint

Static checker for GitHub Actions workflow files, distributed as a prebuilt binary.

This GitHub Packages build comes from [`kajaux/actionlint`](https://github.com/kajaux/actionlint),
the isolated release-test repository for [`kjanat/actionlint`](https://github.com/kjanat/actionlint).
Installing it puts an `actionlint` executable on your `PATH`; no Go toolchain is needed.

## Install

Authenticate to GitHub Packages with a token that can read these packages:

```sh
npm config set @kajaux:registry https://npm.pkg.github.com
npm login --scope=@kajaux --auth-type=legacy --registry=https://npm.pkg.github.com
```

Then install:

```sh
npm install --save-dev @kajaux/actionlint
```

Or run it without adding it to the project:

```sh
npx @kajaux/actionlint
```

## Usage

After installing it as a development dependency, run it in the project and it finds the workflows itself:

```sh
npx actionlint
```

As a package script:

```json
{
  "scripts": {
    "lint:workflows": "actionlint"
  }
}
```

`actionlint` exits `0` when it finds nothing, `1` when it reports problems, and `2` or `3` on a usage error or a fatal
error. Those statuses are forwarded verbatim, so it drops into CI unchanged.

See the [usage documentation](https://github.com/kjanat/actionlint/blob/master/docs/usage.md) for the full command line,
and [the checks list](https://github.com/kjanat/actionlint/blob/master/docs/checks.md) for what it looks for.

The manual page ships in the package as `man/actionlint.1`. npm registered man pages with the system `man` program up to
v11; from v12 it no longer does, so on a current npm read it directly:

```sh
man ./node_modules/@kajaux/actionlint/man/actionlint.1
```

### Configuration schema

The package includes `actionlint.schema.json` for completion and validation in configuration editors. For a config at
`.github/actionlint.yaml`, use the installed copy with:

```yaml
# yaml-language-server: $schema=../node_modules/@kajaux/actionlint/actionlint.schema.json
```

ShellCheck directives reference `schemas/shellcheck/0.11.0.schema.json` relative
to the main schema. The package includes that file, so local schema validation
can run offline. Use the installed schema: this test package is published to GitHub Packages, not npmjs CDNs.

### ShellCheck and Pyflakes

`actionlint` also checks the shell scripts inside `run:` steps with [ShellCheck][shellcheck], and Python scripts with
[Pyflakes][pyflakes], when those are on your `PATH`. Neither is bundled here; install them separately to enable those
checks.

## How this package is put together

This package contains no binary itself. It declares one `optionalDependencies` entry per platform, each published
under the same `@kajaux` scope:

| Package                           | Runs on               |
| --------------------------------- | --------------------- |
| `@kajaux/actionlint-darwin-arm64` | macOS Apple silicon   |
| `@kajaux/actionlint-darwin-x64`   | macOS Intel           |
| `@kajaux/actionlint-freebsd-ia32` | FreeBSD 32-bit x86    |
| `@kajaux/actionlint-freebsd-x64`  | FreeBSD x86-64        |
| `@kajaux/actionlint-linux-arm64`  | Linux ARM64           |
| `@kajaux/actionlint-linux-arm`    | Linux ARMv6 and ARMv7 |
| `@kajaux/actionlint-linux-ia32`   | Linux 32-bit x86      |
| `@kajaux/actionlint-linux-x64`    | Linux x86-64          |
| `@kajaux/actionlint-win32-arm64`  | Windows ARM64         |
| `@kajaux/actionlint-win32-ia32`   | Windows 32-bit x86    |
| `@kajaux/actionlint-win32-x64`    | Windows x86-64        |

Each declares `os` and `cpu`, so your package manager downloads only the one matching your machine and skips the rest.
The `actionlint` command here is a small launcher that resolves that package and execs the binary inside it.

There is no separate musl build: the binaries are statically linked, so the `linux-*` packages run on Alpine and on
glibc distributions alike.

The binaries are the same ones attached to the [GitHub release][releases]; each archive is verified against the
release's published checksums before being repackaged.

### If the binary is not found

The launcher fails with an explanation, but the usual cause is a package manager that skipped optional dependencies.
Reinstall without `--no-optional` or `--omit=optional`.

Using Bun with `minimumReleaseAge`? Add `@kajaux/*` to `minimumReleaseAgeExcludes` alongside
`@kajaux/actionlint`. A fresh release otherwise installs the facade while its binaries are still age-gated.

## Analysis results

CLI `check --json` and the GitHub Action emit the same versioned result contract.
This package ships the schema and types supported by its matching binary:

```typescript
import type { CheckResult, DiagnosticRecord } from '@kajaux/actionlint/result';
```

`@kajaux/actionlint/result.schema.json` exports the result schema;
`@kajaux/actionlint/schemas/results/v1.schema.json` selects its contract version.
`CheckResult` describes a complete result, including failed analysis;
`DiagnosticRecord` describes one JSONL record. Types do not validate untrusted JSON.
See [the result contract](https://github.com/kjanat/actionlint/blob/master/docs/results.md)
for compatibility and migration details.

## Other ways to install

Homebrew, Arch (AUR), Scoop, Docker, a download script, and `go install` are all covered in
[the installation documentation](https://github.com/kjanat/actionlint/blob/master/docs/install.md).

## License

MIT. See [LICENSE.txt](./LICENSE.txt).

[pyflakes]: https://pypi.org/project/pyflakes/
[releases]: https://github.com/kajaux/actionlint/releases
[shellcheck]: https://www.shellcheck.net/
