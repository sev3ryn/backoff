# Backoff

[![tests](https://github.com/sev3ryn/backoff/workflows/Test/badge.svg)](https://github.com/sev3ryn/backoff/actions)
[![GitHub license](https://img.shields.io/github/license/sev3ryn/backoff.svg)](https://github.com/sev3ryn/backoff/blob/master/LICENSE)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/sev3ryn/backoff)](https://pkg.go.dev/github.com/sev3ryn/backoff)



Boring-simple, goroutine-safe Go(golang) library for exponential backoff with jitter(optional)

## Installationśś

```
go get -u github.com/sev3ryn/backoff
```


## Usage

```go

b := backoff.Exponential(backoff.Config{
	Min:       100 * time.Millisecond,
	Max:       4 * time.Second,
	Factor:    2,
	JitterPct: 40,
})

for i := 0; ; i++ {
	err := doSmth()
	if err == nil {
		break
	}

	time.Sleep(b.Attempt(i))
}

```