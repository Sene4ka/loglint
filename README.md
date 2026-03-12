# loglint

Static analysis tool for detecting style violations and sensitive data leaks in Go log messages.

## Rules

| Rule | Description                                               | Defaults                                                                    |
|------|-----------------------------------------------------------|-----------------------------------------------------------------------------|
| `shouldStartWithLowercase` | Log message must start with lowercase letter              | Enabled: `true`                                                             |
| `shouldContainOnlyEnglish` | Only English letters allowed in messages                  | Enabled: `true`                                                             |
| `shouldNotContainSpecialSymbols` | No emojis or special characters (configurable allow-list) | Enabled: `true`, Allowed: `[':', '_', '-', '=', '%']`                       |
| `shouldNotContainSensitiveInformation` | Detects keywords in strings and varible names             | Enabled: `true`, Keywords: `['key', 'password', 'secret', 'auth', 'token']` |

## Installation

### Standalone

#### Variant 1
```bash
go install github.com/Sene4ka/loglint/cmd/loglint@latest
```
#### Variant 2
```bash
git clone https://github.com/Sene4ka/loglint.git
cd log-linter

go build -o loglint ./cmd/loglint
# or
go install ./cmd/loglint
```

#### Usage
```bash
# Run with auto-detected .loglint.yml
loglint ./...

# Run with explicit config path
loglint --config ./configs/loglint.yml ./...

# Run on specific files
loglint ./pkg/... ./cmd/...
```

#### See [./loglint.example.yml](.loglint.example.yml) for config options

### Golangci-lint plugin (Doesn't work that good on Windows)

#### Create .custom-gcl.yml [Example](.custom-gcl.example.yml)

#### Run

```bash
# May crash on Windows because of git-specifics
golangci-lint custom
```

#### Enable plugin in .golangci.yml [Example with configuration settings](.golangci.example.yml)
