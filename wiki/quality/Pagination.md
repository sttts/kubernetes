---
title: Pagination
kategorie: Quality
unterkategorie: Data Transfer Parsimony
quelle: https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/Pagination
---

# Pagination

*a.k.a.* Query with Partial Result Sets, Response Sequence

**Kurzform:** Große Ergebnismengen werden in transportierbare Teilmengen („Chunks", „Pages")
zerlegt; pro Antwortnachricht geht ein Chunk an den Client, zusammen mit Metadaten, wie er den
nächsten anfordert.

## Kontext

Clients fragen Collections ab, um sie einem Benutzer anzuzeigen oder weiterzuverarbeiten. Die
Treffermenge ist oft deutlich größer als das, was der Client tatsächlich braucht oder verarbeiten
kann. Die Elemente können homogen sein (Zeilen aus einer relationalen Tabelle) oder heterogen
(Fragmente aus einem Dokumentspeicher). Entscheidend für die Variantenwahl: ob sich der Datenbestand
*während* des Abrufs ändert.

## Problem

Wie liefert ein API-Provider große Sequenzen strukturierter Daten aus, ohne Clients, Netz und
sich selbst zu überlasten?

## Forces

- **Performance, Skalierbarkeit, Ressourcenverbrauch:** Serialisierung und Übertragung einer
  vollständigen Collection binden Speicher auf beiden Seiten.
- **Individuelle Informationsbedarfe:** Der Client braucht meist nur die ersten n Elemente.
- **Lose Kopplung und Interoperabilität:** Der Paginierungsmechanismus wird Teil des Vertrags.
- **Developer Experience:** Chunking verlagert eine Schleife in den Client.
- **Sicherheit und Datenschutz:** Paginierung erlaubt es, das Abgreifen ganzer Bestände zu bremsen.
- **Test- und Wartungsaufwand:** Randfälle (leerer Chunk, gelöschte Elemente, abgelaufener Cursor).
- **Session-Bewusstsein und Isolation:** Serverseitiger Cursor-State kostet Ressourcen und bindet
  den Client an eine Instanz.
- **Größe und Zugriffsprofil des Datensatzes:** Statischer Report vs. hochfrequent mutierende Liste.

## Lösung

Die Ergebnismenge in Chunks aufteilen, pro Antwort einen Chunk senden und den Client über die
Gesamt- bzw. Restanzahl informieren. Optional Filter anbieten, um die Menge vorab zu verkleinern,
und einen Verweis (Link, Token) auf den Folge-Chunk mitliefern.

## Varianten

| Variante | Selektion | Verhalten bei Änderungen |
|---|---|---|
| *Page-Based Pagination* | Seitenindex + Seitengröße | Elemente können übersprungen/doppelt gesehen werden |
| *Offset-Based Pagination* | `offset` + `limit` | dito; äquivalent zu Page-Based (offset = page × size) |
| *Cursor-Based Pagination* (a.k.a. *Token-Based*) | Id Element des letzten Elements + Anzahl | stabil, da nicht positionsabhängig |
| *Time-Based Pagination* | Zeitstempel statt ID | stabil, selten genutzt |

Die positionsbasierten Varianten (Page/Offset) sind bei veränderlichen Daten fehleranfällig: Wird
zwischen dem Abruf von Seite 1 und Seite 2 ein Element am Anfang entfernt, rutscht ein Element von
Seite 2 auf Seite 1 und wird nie ausgeliefert.

## Beispiel

Offset-basiert (Lakeside Mutual), mit Kontroll-[Metadata Elements](../structure/MetadataElement.md)
und HATEOAS-Link auf den Folge-Chunk:

```json
{
  "offset": 0,
  "limit": 2,
  "size": 50,
  "customers": [ /* ... */ ],
  "_links": { "next": { "href": "http://localhost:8080/customers?limit=2&offset=2" } }
}
```

Cursor-basiert: Der Client erhält statt `offset` ein undurchsichtiges Token, das er unverändert
zurückschickt — im Beispiel `?page-size=2&cursor=mfn834fj`.

## Konsequenzen

**Vorteile:**

- Nachrichtengrößen und Peak-Speicherverbrauch werden vorhersagbar begrenzt.
- Der Client kann früh mit der Verarbeitung beginnen und vorzeitig abbrechen.
- Ein `size`/`remaining`-Feld erlaubt Fortschrittsanzeigen, ohne selbst notwendig zu sein.

**Nachteile / Kosten:**

- Mehr Roundtrips; die Gesamtlatenz für den vollständigen Bestand steigt.
- Der Provider braucht eine stabile Ordnung; ohne sie ist keine Variante korrekt.
- Positionsbasierte Varianten sind bei Nebenläufigkeit schlicht falsch.
- Cursor-Token sind Zustand: Sie können ablaufen oder durch Konfigurationsänderungen ungültig werden.

## Bekannte Verwendungen

JSON:API spezifiziert Offset- und Page-Varianten. Weit verbreitet in Public Web APIs (GitHub Query
API, Slack, Google-Suche, JIRA Cloud REST API). Die Twitter-REST-API nutzt wegen der ständig
mutierenden Timeline `since_id` — de facto Cursor-Pagination über ein
[Id Element](../structure/IdElement.md). RFC 5005 beschreibt Feed Paging für Atom.

## Verwandte Patterns

- [Request Bundle](RequestBundle.md) — das genaue Gegenteil: bündelt viele kleine Nachrichten, statt eine große zu zerlegen.
- [Wish List](WishList.md) — reduziert die Nachrichtengröße über die Breite (Felder) statt über die Länge (Elemente).
- [Wish Template](WishTemplate.md) — dasselbe, aber strukturiert.
- [Conditional Request](ConditionalRequest.md) — vermeidet die Übertragung ganz, wenn sich nichts geändert hat.
- [Atomic Parameter List](../structure/AtomicParameterList.md) — typische Struktur der Query-Parameter einer paginierten Abfrage.
- [Parameter Tree](../structure/ParameterTree.md) — typische Struktur des einzelnen Chunks.
- [Id Element](../structure/IdElement.md) — Basis der Cursor-Variante.
- [Rate Limit](RateLimit.md) — Paginierung entlastet das Kontingent.
- [Information Holder Resource](../responsibility/InformationHolderResource.md) — der übliche Träger paginierter Abfragen.

## Bezug zu Kubernetes / KRM

KRM erfüllt das Pattern, aber **ausschließlich in der Cursor-Variante** — Offset und Seitennummer
gibt es bewusst nicht. `metav1.ListOptions` trägt `limit` (max. Anzahl Objekte) und `continue`
(opakes Fortsetzungs-Token); `metav1.ListMeta` liefert `continue` zurück, solange weitere Daten
existieren, plus optional `remainingItemCount` als *Schätzung* der Restmenge (nicht gesetzt, wenn
Label- oder Field-Selektoren im Spiel sind). Ein leeres `continue` ist das einzige verlässliche
Ende-Signal — ein Chunk darf auch null Objekte enthalten, wenn alle herausgefiltert wurden.

```http
GET /api/v1/namespaces/default/pods?limit=500

{
  "kind": "PodList",
  "apiVersion": "v1",
  "metadata": {
    "resourceVersion": "8675309",
    "continue": "eyJ2IjoibWV0YS5rOHMuaW8vdjEiLCJydiI6ODY3NTMwOSwic3RhcnQiOiJk...",
    "remainingItemCount": 1483        // <- nur eine Schätzung, fehlt bei aktivem Selektor
  },
  "items": [ /* 500 Pods */ ]
}

GET /api/v1/namespaces/default/pods?limit=500&continue=eyJ2IjoibWV0YS5rOHMuaW8vdjEiLCJydiI6ODY3NTMwOSwic3RhcnQiOiJk...

{
  "metadata": {
    "resourceVersion": "8675309",     // <- unverändert: derselbe Snapshot wie Chunk 1
    "continue": ""                    // <- leer = Ende, unabhängig von len(items)
  },
  "items": [ /* Rest */ ]
}
```

Das Token ist nicht wirklich opak, sondern base64url-kodiertes JSON (siehe
`k8s.io/apiserver/pkg/storage/continue.go`): der etcd-Startschlüssel des nächsten Bereichs
plus die `resourceVersion` des ersten Requests.

```json
// base64url-dekodiertes continue-Token — Clients dürfen das nicht auswerten:
{"v":"meta.k8s.io/v1","rv":8675309,"start":"default/web-7d9f8c-nkm2p"}
```

Genau daraus folgt die zentrale Garantie: Alle Chunks werden aus demselben MVCC-Snapshot bedient,
sodass die Sequenz identisch zu einem einzigen
`list` ohne `limit` ist — kein Objekt wird übersprungen oder doppelt geliefert. Diese Kombination
aus geordneten Key-Ranges und Snapshot-Read ist der Grund, warum ein Offset-Modell für KRM nie
in Frage kam: Es könnte diese Konsistenz nicht liefern.

Der Preis ist Vergänglichkeit. Snapshots werden nach etwa fünf bis fünfzehn Minuten kompaktiert;
danach antwortet der Server mit `410 Gone` und `Status.Reason: "Expired"`
(`metav1.StatusReasonExpired`). Die Fehlerantwort enthält ihrerseits ein neues `continue`-Token:
Wer Konsistenz braucht, muss von vorn listen; wer nur Vollständigkeit im weiteren Sinne braucht,
kann mit dem Token aus dem Fehler weiterlaufen — dann allerdings aus einem neueren Snapshot und
damit inkonsistent zu den bereits gelesenen Chunks. Diese explizite Wahl zwischen Konsistenz und
Fortsetzbarkeit ist in MAP nicht vorgesehen. Der Fehler ist ein gewöhnliches `metav1.Status`-Objekt,
und ausgerechnet dessen `metadata` trägt das Ersatz-Token:

```http
HTTP/1.1 410 Gone

{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {
    "continue": "eyJ2IjoibWV0YS5rOHMuaW8vdjEiLCJydiI6LTEsInN0YXJ0Ijoi..."  // <- rv:-1 = ab jetzt
  },
  "status": "Failure",
  "message": "The provided continue parameter is too old to display a consistent list result. You can start a new list without the continue parameter, or use the continue token in this response to retrieve the remainder of the results. Continuing with the provided token results in an inconsistent list ...",
  "reason": "Expired",          // <- metav1.StatusReasonExpired
  "code": 410
}
```

Clientseitig ist Chunking Standard: `k8s.io/client-go/tools/pager` verwendet `defaultPageSize = 500`,
`kubectl` exponiert es als `--chunk-size` (Default ebenfalls 500), und `Reflector.WatchListPageSize`
steuert das Chunking initialer Listen für Informer. Umgekehrt ist Paginierung bei `watch=true`
nicht unterstützt — ein Watch ist eine unbegrenzte Ereignissequenz, kein Chunk-Verfahren; sein
Volumenproblem wird über [Conditional Request](ConditionalRequest.md)-artige Mechanismen
(`resourceVersion`, Bookmarks) gelöst. Als orthogonale Ergänzung gibt es seit KEP-5866 hinter dem
Alpha-Gate `ShardedListAndWatch` einen `shardSelector` in `ListOptions` und ein `shardInfo` in
`ListMeta`: Damit wird die Collection nicht sequenziell in Seiten, sondern über einen Hash-Bereich
von `metadata.uid` in parallel abrufbare Shards zerlegt — eine Partitionierung, die MAP unter
*Pagination* nicht kennt.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/Pagination)
