# Cores Common

Gemeinsame Go-Pakete für die Cores-Dienste. Das Modul bündelt wiederverwendbare
Bausteine für Branding, Konfiguration, Healthchecks, JWT-Verarbeitung, CORS und
JSON-Antworten.

## Verwendung

```go
require github.com/nbt4/cores-common <version>
```

Das Modul benötigt Go 1.25. Neue Pakete müssen serviceunabhängig bleiben und mit
`go test ./...` geprüft werden.

## Quellcode

[GitHub: nbt4/cores-common](https://github.com/nbt4/cores-common) ·
[Cores-Monorepo](https://github.com/nbt4/cores)
