package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/apognu/gocal"
    "github.com/gorilla/feeds"
)

func convertICStoRSS(w http.ResponseWriter, r *http.Request) {
    // Lese die ICS-URL aus dem Query-Parameter
    icsURL := r.URL.Query().Get("ics")
    if icsURL == "" {
        http.Error(w, "Missing 'ics' query parameter", http.StatusBadRequest)
        return
    }

    // ICS-Datei herunterladen
    resp, err := http.Get(icsURL)
    if err != nil {
        http.Error(w, "Unable to fetch ICS file", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // ICS-Datei parsen
    parser := gocal.NewParser(resp.Body)
    parser.Parse()

    // RSS-Feed erstellen
    feed := &feeds.Feed{
        Title:       "Converted Calendar Feed",
        Link:        &feeds.Link{Href: icsURL},
        Description: "This is a converted calendar feed",
        Created:     time.Now(),
    }

    for _, event := range parser.Events {
        // Dereferenziere event.Start, um den time.Time Wert zu erhalten
        item := &feeds.Item{
            Title:       event.Summary,
            Description: event.Description,
            Link:        &feeds.Link{Href: icsURL},
            Created:     *event.Start,
        }
        feed.Items = append(feed.Items, item)
    }

    // RSS als HTTP-Antwort senden
    rssData, err := feed.ToRss()
    if err != nil {
        http.Error(w, "Unable to generate RSS feed", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/rss+xml")
    w.Write([]byte(rssData))
}

func main() {
    http.HandleFunc("/rss", convertICStoRSS)
    fmt.Println("Server is running at :8080")
    http.ListenAndServe(":8080", nil)
}
