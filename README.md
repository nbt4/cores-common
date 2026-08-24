# Cores Common

Gemeinsame Go-Pakete für RentalCore, WarehouseCore, PlannerCore,
ProcurementCore und Cores Dashboard. Das Modul bündelt wiederverwendbare
Bausteine für Branding, Konfiguration, Healthchecks, JWT-Verarbeitung, CORS und
JSON-Antworten.

## Verwendung

```go
require github.com/nbt4/cores-common <version>
```

Das Modul benötigt Go 1.25. Neue Pakete müssen serviceunabhängig bleiben und mit
`go test ./...` geprüft werden.

## Branding

`pkg/branding` definiert die zentrale Konfiguration und semantische Asset-Sätze
für Produkt- und Unternehmensmarken. Der Service liest
`branding_config.assets_json`, hält die bisherigen Sidebar-/Login-/Favicon-
Spalten als Fallback kompatibel und erzeugt dynamische PWA-Manifeste. Ein
Favicon wird bewusst nicht automatisch als App-Icon verwendet.

## Quellcode

[GitHub: nbt4/cores-common](https://github.com/nbt4/cores-common) ·
[Cores-Monorepo](https://github.com/nbt4/cores)
