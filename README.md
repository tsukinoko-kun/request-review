# request-review

Forgeless code reviews

## Usage

### Install on your system

Install binary from Releases: https://github.com/tsukinoko-kun/request-review/releases/latest

Install via Homebrew:

```shell
brew install tsukinoko-kun/tap/request-review
```

Install via Scoop:

```shell
scoop bucket add tsukinoko-kun https://github.com/tsukinoko-kun/scoop-bucket
scoop install tsukinoko-kun/request-review
```

Install via Go install:

```shell
go install github.com/tsukinoko-kun/request-review/cmd/request-review@latest
```

### Run Docker image

Example for Windows PWSH (I haven't tried this)

```shell
docker run --rm -w /wd -v "${PWD}:/wd" -v ~/.git-credentials:/root/.git-credentials:ro ghcr.io/tsukinoko-kun/request-review:latest
```
