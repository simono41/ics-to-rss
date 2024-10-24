package main

import (
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

    log.Println("[INFO] Generating RSS feed")

    // RSS als HTTP-Antwort senden
    rssData, err := feed.ToRss()
    if err != nil {
        log.Printf("[ERROR] Unable to generate RSS feed: %v\n", err)
        http.Error(w, "Unable to generate RSS feed", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/rss+xml")
    w.Write([]byte(rssData))

    log.Println("[INFO] RSS feed successfully generated and sent")
}

func main() {
    http.HandleFunc("/rss", convertICStoRSS)
    log.Println("[INFO] Server is running at :8080")
    
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("[FATAL] Server failed to start: %v\n", err)
    }
}
