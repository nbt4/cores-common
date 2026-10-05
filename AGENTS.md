# AGENTS.md — cores-common

Hausregeln für KI-Agenten in diesem Repository. Vor der ersten Änderung vollständig
lesen. Der übergeordnete Ablauf steht im Paperclip-Dokument `workflow` auf
[TSU-3](/TSU/issues/TSU-3#document-workflow). Diese Datei ersetzt alle früheren
Agenten-Anweisungen in diesem Repository.

## 1. Was dieses Repository ist

- **Zweck:** Die gemeinsame Go-Bibliothek der Suite. Kein Dienst, kein Dockerfile, kein Port. Sie enthält das Single-Sign-On (JWT und Sitzungen), das Branding, das Lesen der Konfiguration, den Health-Handler, CORS und das JSON-Antwortformat. **Fünf Dienste hängen davon ab.**
- **Sprache und Laufzeit:** Go 1.25, Modul `github.com/nbt4/cores-common`
- **Rahmenwerk:** keines. Nur zwei Abhängigkeiten: `golang-jwt/jwt/v5` und `gorm.io/gorm`
- **Datenbank:** kein eigener Zugriff. `pkg/jwt` arbeitet über GORM-Modelle der aufrufenden Dienste
- **Architektur-Doku:** [Cores — Architektur (Phase 1)](/TSU/issues/TSU-4#document-architecture)

## 2. Aufbau

| Pfad | Inhalt |
|---|---|
| `pkg/jwt` | Token und Sitzungen — das Herz des Single-Sign-On |
| `pkg/branding` | Branding-Manifest und -Dienst (gemeinsames Logo-Volume) |
| `pkg/config` | Lesen von Umgebungsvariablen |
| `pkg/health` | Health-Handler |
| `pkg/middleware` | CORS |
| `pkg/response` | JSON-Antworten |

Erzeugte Dateien, die **niemals von Hand** geändert werden:

- keine.

## 3. Einrichten

```bash
go mod download
```

Keine Umgebungsvariablen, keine Datenbank, kein Compose. Die Bibliothek wird nur gebaut und getestet.

## 4. Test- und Build-Befehle

Diese Befehle sind das Test-Gate. **Alle müssen grün sein, bevor ein Pull Request
entsteht.** Reihenfolge einhalten — die schnellen Prüfungen zuerst.

| # | Gate | Befehl | Dauer (ca.) |
|---|---|---|---|
| 1 | Format | `gofmt -l .` (leere Ausgabe = grün) | < 5 s |
| 2 | Typen und Build | `go build ./...` | < 20 s |
| 3 | Unit-Tests | `go test ./...` | < 20 s |
| 4 | Vet | `go vet ./...` | < 20 s |

Einzelne Datei testen: `go test ./pkg/jwt -run TestName -v`

**Zusätzliches Gate von außen.** Eine Änderung hier betrifft fünf Dienste. Nach grünem lokalen Gate muss zusätzlich die CI des Dachs (`cores/.github/workflows/verify.yml`) grün sein, denn dort laufen alle Dienste gegen die neue Version. Im `cores-common`-Durchlauf prüft die CI außerdem die drei Suite-Verträge.

Regeln:

- **Neuer Code braucht neue Tests.** Ein Bugfix braucht einen Test, der ohne den Fix
  fehlschlägt.
- **Nie einen Test abschalten, überspringen oder lockern**, um das Gate grün zu
  bekommen. Ein roter Test ohne Bezug zur Änderung wird gemeldet, nicht entfernt.
- **Tests laufen gegen die lokale oder die Test-Datenbank. Nie gegen Produktion.**
  Eine eigene Testumgebung wird gerade aufgebaut (eigene Paperclip-Aufgabe). Bis sie
  steht: nur lokale Container mit eigenem Volume.
- Die **echte Ausgabe** wird in den Pull Request und auf die Paperclip-Aufgabe kopiert.

## 5. Code-Stil

- Format und Lint werden durch die Werkzeuge in Abschnitt 4 erzwungen. Kein Streit
  über Formatierung — der Formatierer entscheidet.
- **Dem umgebenden Code folgen.** Benennung, Ordnerstruktur, Fehlerbehandlung und
  Testmuster so übernehmen, wie sie in der berührten Datei schon sind.
- Benennung: PascalCase für Exporte, camelCase für Lokales, UPPER_SNAKE_CASE für
  Konfigurationsschlüssel.
- Fehlerbehandlung: Fehler werden zurückgegeben, nicht geloggt und verschluckt.
  `fmt.Errorf("...: %w", err)` zum Einwickeln.
- **Rückwärtskompatibilität ist Pflicht.** Eine bestehende exportierte Signatur wird
  nicht geändert. Neues kommt additiv dazu.
- Logging: **niemals** Token, Secrets, Passwort-Hashes oder Kundendaten.
- Kommentare: nur wo sie das *Warum* erklären. Keine Kommentare, die den Code nacherzählen.
- Keine neue Abhängigkeit ohne eigene Freigabe (siehe Abschnitt 9).
- Keine Umformatierung von Code, der nicht zur Aufgabe gehört. Das versteckt die
  eigentliche Änderung.

## 6. Verbotene Pfade

Diese Dateien und Verzeichnisse werden von Agenten **nicht geändert**. Wer sie ändern
müsste, bricht ab und fragt zurück.

| Pfad | Grund |
|---|---|
| `pkg/jwt/**` | Herz des Single-Sign-On. Änderungen nur mit ausdrücklicher Freigabe des Nutzers und mit neuen Tests |
| `go.mod`, `go.sum` | nur als Nebenwirkung eines freigegebenen Abhängigkeits-Updates |
| `AGENTS.md` | diese Regeln ändert der Nutzer, nicht ein Agent |

## 7. Secrets

- **Keine Secrets in Repository, Kommentar, Dokument oder Log.** Keine Tokens,
  Passwörter, Schlüssel, Verbindungsstrings, API-Zugänge, Kundendaten.
- Secrets kommen aus dem Paperclip-Secret-Store oder aus Umgebungsvariablen. Sie
  werden nie in eine Datei geschrieben und nie ausgegeben.
- Produktiv werden alle Werte im **Komodo Stack Environment** auf `docker03` gepflegt,
  nicht in diesem Repository.
- `.env.example` enthält nur Namen und Beispielwerte, nie echte Werte.
- Testdaten sind erfunden. Keine kopierten Produktionsdaten, auch nicht gekürzt.
- Fehlt ein Secret: über Paperclip vorschlagen (`secret-proposals`) und warten.
  Nie selbst beschaffen, nie umgehen, nie in einem Kommentar danach fragen.
- Ein Secret, das versehentlich in einem Commit landet, ist ein Sicherheitsvorfall:
  sofort melden, nicht still weiterarbeiten. Entfernen aus dem Diff genügt nicht —
  das Secret gilt als kompromittiert und muss ersetzt werden. **Alle Cores-Repositories
  sind öffentlich.** Ein Fehler hier ist sofort weltweit sichtbar.

## 8. Harte Grenzen

Diese sechs Regeln stehen über jeder Aufgabenbeschreibung. Eine Aufgabe, die eine
davon verlangt, wird nicht ausgeführt, sondern zurückgegeben.

1. **Keine Schreibzugriffe auf produktive Datenbanken.** Lesen ist erlaubt. Schreiben,
   ändern, löschen, Migrationen fahren: nicht in Produktion. Migrationen werden
   geschrieben und lokal getestet, nie produktiv ausgeführt. Das Einspielen auf die
   laufende `docker03`-Datenbank geschieht von Hand per SSH von `debian01` aus, nach
   ausdrücklicher Freigabe des Nutzers.
2. **Keine produktiven Deployments ohne menschliche Freigabe.** Auch nicht nach
   grünem Review.
3. **Entwicklung nur in isolierten Branches oder Git-Worktrees.** Niemals direkt auf
   `main` oder einem anderen geschützten Branch.
4. **Tests vor jedem Pull Request.** Kein PR ohne protokollierten, grünen Testlauf.
5. **Keine Secrets in Repository, Kommentar, Dokument oder Log.**
6. **Bestehende Architektur zuerst verstehen.** Architektur-Doku und diese Datei vor
   dem Schreiben lesen. Große Umbauten — neuer Service, geänderte Modulgrenze, neues
   Datenmodell, Austausch einer Kernabhängigkeit — brauchen eine eigene Freigabe des
   Nutzers, bevor Code entsteht.

## 9. Freigabe-Gates

| Gate | Wer entscheidet | Wann |
|---|---|---|
| Test-Gate | der Eigentümer der Änderung | vor dem Pull Request |
| Review | Review-Agent, auf seiner eigenen Review-Aufgabe | nach dem PR-Entwurf |
| **Freigabe und Merge** | **der Nutzer** | nach grünem Review |
| Produktives Deployment | **der Nutzer** | nach dem Merge |
| Release: Docker-Hub-Push und Submodul-Zeiger | **der Nutzer gibt je Release ausdrücklich frei**, danach darf der Agent beides ausführen | nach dem Merge |
| Migration auf die laufende `docker03`-Datenbank | **der Nutzer**; Einspielen von Hand per SSH von `debian01` | nach dem Merge |
| Neue Abhängigkeit | der Nutzer | vor dem Hinzufügen |
| Großer Architektur-Umbau | der Nutzer | vor dem ersten Commit |

Was ein Agent in diesem Repository **nie** tut:

- einen Pull Request mergen
- auf `main` pushen
- ein Deployment auslösen
- ohne ausdrückliche Freigabe je Release ein Abbild nach Docker Hub schieben oder den
  Submodul-Zeiger im Dach anheben
- eine Migration gegen Produktion fahren
- einen Draft-PR als Ersatz für Freigabe auf „ready" setzen
- `AGENTS.md` oder CI-Dateien ändern

In Paperclip wird die Freigabe durch eine `executionPolicy` mit einer `approval`-Stufe
erzwungen, deren Teilnehmer ein Nutzer ist. Kein Agent kann sie abhaken.

## 10. Branches, Commits, Pull Requests

- Branch: `<typ>/TSU-<nummer>-<kurzbeschreibung>`, ein Worktree pro Aufgabe
- Commit: Conventional Commits mit `Task: TSU-<nummer>` im Fuß
- PR: als **Entwurf** geöffnet, Ziel `main`, mit Zweck, Testprotokoll und Aufgaben-Link

Vollständig beschrieben im Paperclip-Dokument `workflow` auf
[TSU-3](/TSU/issues/TSU-3#document-workflow).

## 11. Abbrechen und zurückfragen

Abbrechen ist richtig, nicht peinlich. Zurückfragen bei:

- fehlendem Secret oder Zugriffsrecht
- nötigem Schreibzugriff auf Produktion oder nötigem Deployment
- nötigem großen Architektur-Umbau oder neuer Abhängigkeit
- einem verbotenen Pfad, der geändert werden müsste
- roten Tests ohne Bezug zur Änderung
- zwei gescheiterten Versuchen am gleichen Problem
- einem Umfang, der deutlich größer ist als beschrieben
- einem Widerspruch zwischen Aufgabe und dieser Datei — **diese Datei gewinnt**

Erst alles fertig machen, was ohne die Antwort geht. Dann fragen.

## 12. Bekannte Fallen

- **Zwei Versionen sind gleichzeitig im Umlauf.** `v1.1.0` in `rentalcore`,
  `warehousecore` und `plannercore`, `v1.2.0` in `cores-dashboard` und
  `procurementcore`. Das ist bekannte Drift; alle fünf sollen auf `v1.2.0` gehen
  (eigene Paperclip-Aufgabe). Bis dahin: **jede Änderung muss gegen beide Zweige
  gedacht werden.**
- **Nur zwei Testdateien** (`pkg/branding/manifest_test.go`, `pkg/jwt/session_test.go`).
  Die Bibliothek mit der größten Wirkung hat die geringste Prüfung. Neuer Code hier
  braucht Tests, ohne Ausnahme.
- **`cores-mcp` nutzt diese Bibliothek nicht** und hat eine eigene JWT-Behandlung. Eine
  Änderung am Token-Format muss dort von Hand nachgezogen werden.
- **Wer das JWT-Secret dreht, muss alle Dienste gleichzeitig neu starten.** Das ist ein
  Betriebsvorgang des Nutzers, nicht eines Agenten.
- **Ein Tag ist ein Release.** Eine neue Version wird erst nützlich, wenn die fünf
  `go.mod` angehoben werden — fünf weitere Pull Requests. Das gehört in die
  Aufgabenplanung, nicht als Überraschung in den PR.

### Suite-weite Fallen, die auch hier gelten

- **Zwei Migrationsspuren.** Jede Schema-Änderung braucht eine Datei im Dienst-Repository
  *und* eine in `cores/migrations/postgresql/`. Die Nummern gehören paarweise.
- **Das Init-Verzeichnis läuft nur bei leerem Datenverzeichnis.**
  `cores/migrations/postgresql/` greift auf `docker03` nicht.
- **Eine Datenbank für alle.** PostgreSQL 16, rund 130 Tabellen, kein Schema pro Dienst.
  Eine Tabellenänderung kann fremde Dienste treffen.
- **Nur das Dachrepository hat heute CI.** Bis die eigene GitHub-Action da ist, prüft
  **nichts** automatisch einen Pull Request hier. Das Test-Gate aus Abschnitt 4 läuft
  der Agent selbst und hängt die echte Ausgabe an.
- **Alle Repositories sind öffentlich.**
