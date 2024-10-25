# ICS to RSS Converter

Dieses Projekt ist eine Go-Anwendung, die ICS-Kalenderdateien in RSS-Feeds konvertiert. Die Anwendung wird in einem Docker-Container ausgeführt und kann mit Docker Compose einfach bereitgestellt werden.

## Voraussetzungen

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Installation

1. **Repository klonen:**

   ```bash
   git clone https://github.com/dein-benutzername/ics-to-rss.git
   cd ics-to-rss
   ```

2. **Docker Compose starten:**

   Stelle sicher, dass Docker und Docker Compose installiert sind und führe dann den folgenden Befehl aus:

   ```bash
   docker-compose up -d --build
   ```

   Dieser Befehl baut das Docker-Image und startet den Container.

## Verwendung

Die Anwendung läuft auf `http://localhost:8080`. Um eine ICS-Datei in einen RSS-Feed zu konvertieren, verwende die folgende URL-Struktur:

```
http://localhost:8080/rss?ics=<ICS_URL>
```

Ersetze `<ICS_URL>` durch die URL deiner ICS-Datei.

Und Optional wenn nur die Events des heutigen Tages, Woche, oder Monats ausgegeben werden soll:

- Für alle Events: `/rss?ics=<ICS_URL>`
- Für heutige Events: `/rss?ics=<ICS_URL>&range=today`
- Für Events dieser Woche: `/rss?ics=<ICS_URL>&range=week`
- Für Events dieses Monats: `/rss?ics=<ICS_URL>&range=month`

## Docker Compose Datei

Hier ist der Inhalt der `docker-compose.yml` Datei:

```yaml
version: '3.8'

services:
  ics-to-rss:
    build: .
    ports:
      - "8080:8080"
```

## Dockerfile

Stelle sicher, dass du auch eine `Dockerfile` im selben Verzeichnis hast:

```dockerfile
# Verwende ein offizielles Golang-Image als Build-Umgebung
FROM golang:1.20 as builder

WORKDIR /app

# Kopiere den Go-Modul-Dateien und installiere Abhängigkeiten
COPY go.mod go.sum ./
RUN go mod download

# Kopiere den Rest des Codes
COPY . .

# Baue die Anwendung
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Verwende ein schlankes Image für die Produktion
FROM alpine:3.18

WORKDIR /root/

# Kopiere das gebaute Go-Binary aus der vorherigen Stufe
COPY --from=builder /app/main .

# Exponiere Port 8080 und starte die Anwendung
EXPOSE 8080
CMD ["./main"]
```

## Lizenz

Dieses Projekt steht unter der MIT-Lizenz. Weitere Informationen findest du in der `LICENSE` Datei.
