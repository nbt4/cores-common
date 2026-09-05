# Cores Common

## Suite-Vertrag

`cores-common` besitzt keine Benutzeroberfläche. Sobald dieses Modul UI-nahe Verträge oder auslieferbare Webressourcen erhält, gelten ohne Ausnahme das [Cores Suite Designsystem](https://github.com/nbt4/cores/blob/main/docs/DESIGN_SYSTEM.md) und dessen Sync-/Prüfworkflow.

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

## Authentifizierte Sitzungen

`pkg/jwt.ValidateSession` prüft HS256-Suite-Tokens mit verpflichtendem Ablaufdatum
und lädt über `DatabaseUserLookup` den aktuellen Benutzerstatus sowie die aktuelle
Administratorrolle. Kontosperren, gelöschte Benutzer und Rollenänderungen gelten
ab der nächsten Anfrage, auch bei bereits ausgestellten Tokens. Datenbankfehler
verweigern den Zugriff; Kontoabfragen besitzen ein Drei-Sekunden-Limit und keinen
Cache. Fachliche Rollen und Planner-Mitgliedschaften bleiben zusätzlich erforderlich.

## Branding-Modell

`pkg/branding` definiert die zentrale Konfiguration und semantische Asset-Sätze
für Produkt- und Unternehmensmarken. Der Service liest
`branding_config.assets_json`, hält die bisherigen Sidebar-/Login-/Favicon-
Spalten als Fallback kompatibel und erzeugt dynamische PWA-Manifeste. Ein
Favicon wird bewusst nicht automatisch als App-Icon verwendet.

## Quellcode

[GitHub: nbt4/cores-common](https://github.com/nbt4/cores-common) ·
[Cores-Monorepo](https://github.com/nbt4/cores)
