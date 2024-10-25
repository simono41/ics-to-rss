package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/apognu/gocal"
    "github.com/gorilla/feeds"
)

func convertICStoRSS(w http.ResponseWriter, r *http.Request) {
    // Lese die ICS-URL aus dem Query-Parameter
    icsURL := r.URL.Query().Get("ics")
    if icsURL == "" {
        log.Println("[ERROR] Missing 'ics' query parameter")
        http.Error(w, "Missing 'ics' query parameter", http.StatusBadRequest)
        return
    }

    // Lese den Zeitraum aus dem Query-Parameter
    timeRange := r.URL.Query().Get("range")
    if timeRange == "" {
        timeRange = "all" // Standardwert, wenn kein Zeitraum angegeben ist
    }

    log.Printf("[INFO] Fetching ICS file from URL: %s\n", icsURL)

    // ICS-Datei herunterladen
    resp, err := http.Get(icsURL)
    if err != nil {
        log.Printf("[ERROR] Unable to fetch ICS file: %v\n", err)
        http.Error(w, "Unable to fetch ICS file", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    log.Println("[INFO] Parsing ICS file")

    // ICS-Datei parsen
    parser := gocal.NewParser(resp.Body)
    parser.Parse()

    // RSS-Feed erstellen
    feed := &feeds.Feed{
        Title:       "Converted Calendar Feed",
        Link:        &feeds.Link{Href: icsURL},
        Description: fmt.Sprintf("This is a converted calendar feed for %s", timeRange),
        Created:     time.Now(),
    }

    now := time.Now()
    for _, event := range parser.Events {
        if shouldIncludeEvent(event, timeRange, now) {
            item := &feeds.Item{
                Title:       event.Summary,
                Description: event.Description,
                Link:        &feeds.Link{Href: icsURL},
                Created:     *event.Start,
            }
            feed.Items = append(feed.Items, item)
        }
    }

    log.Println("[INFO] Generating RSS feed")

    // RSS als XML generieren
    rssData, err := feed.ToRss()
    if err != nil {
        log.Printf("[ERROR] Unable to generate RSS feed: %v\n", err)
        http.Error(w, "Unable to generate RSS feed", http.StatusInternalServerError)
        return
    }

    // Header setzen für XML-Anzeige im Browser
    w.Header().Set("Content-Type", "application/xml; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")

    // XML-Daten an den Browser senden
    fmt.Fprint(w, rssData)

    log.Println("[INFO] RSS feed successfully generated and sent as XML")
}

func shouldIncludeEvent(event gocal.Event, timeRange string, now time.Time) bool {
    switch timeRange {
    case "today":
        return event.Start.Year() == now.Year() && event.Start.YearDay() == now.YearDay()
    case "week":
        _, thisWeek := now.ISOWeek()
        _, eventWeek := event.Start.ISOWeek()
        return event.Start.Year() == now.Year() && eventWeek == thisWeek
    case "month":
        return event.Start.Year() == now.Year() && event.Start.Month() == now.Month()
    default:
        return true // "all" oder jeder andere Wert zeigt alle Ereignisse
    }
}

func main() {
    http.HandleFunc("/rss", convertICStoRSS)
    log.Println("[INFO] Server is running at :8080")
    
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("[FATAL] Server failed to start: %v\n", err)
    }
}
