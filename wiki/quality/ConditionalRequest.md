---
title: Conditional Request
kategorie: Quality
unterkategorie: Data Transfer Parsimony
quelle: https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/ConditionalRequest
---

# Conditional Request

*a.k.a.* Conditional Retrieval, Conditional Modifications (Terminologie der Google Calendar API)

**Kurzform:** Der Request trägt Metadaten über den Zustand, den der Client bereits kennt; der
Provider verarbeitet ihn nur, wenn die daraus gebildete Bedingung erfüllt ist, und antwortet sonst
mit einer minimalen Nachricht.

## Kontext

Clients fragen dieselben serverseitigen Daten wiederholt ab, obwohl sich diese zwischen den
Anfragen selten oder gar nicht ändern — Polling-Schleifen, Caches, periodische Synchronisation.

## Problem

Wie lassen sich unnötige serverseitige Verarbeitung und Bandbreitennutzung vermeiden, wenn häufig
Operationen aufgerufen werden, die selten wechselnde Daten liefern?

## Forces

- **Komplexität von Endpoint-, Client- und Payload-Design:** Der Client muss Zustand mitführen und
  eine zweite Antwortform behandeln.
- **Korrektheit von Reporting und Abrechnung:** Zählt ein `304` als Aufruf? Beeinflusst er ein
  [Rate Limit](RateLimit.md)?
- **Nachrichtengröße:** Der eigentliche Hebel bei reiner Bandbreitenersparnis.
- **Client-Last:** Der Client spart Deserialisierung und Verarbeitung.
- **Provider-Last:** Nur wenn die Bedingung *früh* geprüft wird, spart der Provider auch Rechenzeit.
- **Aktualität vs. Korrektheit:** Wie alt darf eine als „unverändert" gemeldete Sicht sein?

## Lösung

Requests konditional machen, indem [Metadata Elements](../structure/MetadataElement.md) in die
Nachrichtenrepräsentation oder in Protokoll-Header aufgenommen werden. Der Provider verarbeitet die
Anfrage nur, wenn die durch diese Metadaten ausgedrückte Bedingung erfüllt ist.

## Varianten

Die Quelle benennt explizit die *Fingerprint-Based Conditional Request*-Variante (ETag über den
Inhalt). Die zeitbasierte Spielart arbeitet stattdessen mit Zeitstempeln (`If-Modified-Since`,
`If-Unmodified-Since`). Orthogonal dazu: Konditionale **Lesezugriffe** (Caching-Optimierung) versus
konditionale **Schreibzugriffe** (Optimistic Concurrency, `If-Match`).

## Beispiel

Antwort mit Fingerprint:

```http
HTTP/1.1 200
ETag: "0c2c09ecd1ed498aa7d07a516a0e56ebc"
Content-Type: application/hal+json;charset=UTF-8
```

Folgeanfrage mit Bedingung, und die Antwort bei unverändertem Zustand:

```http
GET /customers/1c184cf1-... HTTP/1.1
If-None-Match: "0c2c09ecd1ed498aa7d07a516a0e56ebc"

HTTP/1.1 304 Not Modified
ETag: "0c2c09ecd1ed498aa7d07a516a0e56ebc"
```

Die Quelle weist auf eine wichtige Feinheit hin: Wird der ETag — wie bei Springs
`ShallowEtagHeaderFilter` — als Response-Filter berechnet, ist die Antwort bereits vollständig
erzeugt worden und wird nur verworfen. Das spart Bandbreite, aber keine Rechenzeit. Erst eine
Prüfung auf der Request-Seite, idealerweise gegen einen serverseitigen Cache, spart beides.

## Konsequenzen

**Vorteile:**

- Bei stabilen Daten sinkt das übertragene Volumen auf einen Header.
- Der Client kann seinen lokalen Cache mit definierter Semantik gültig halten.
- Bei GitHub zählen `304`-Antworten nicht gegen das [Rate Limit](RateLimit.md) — ein direkter
  Anreiz, das Pattern zu nutzen.

**Nachteile / Kosten:**

- Zwei Antwortformen pro Operation; Fehlerbehandlung und Tests verdoppeln sich.
- Ein starker ETag verlangt eine deterministische Serialisierung — Map-Iterationsreihenfolge oder
  wechselnde Feldreihenfolgen zerstören ihn.
- Zeitstempelbasierte Varianten haben Auflösungsprobleme (Sekundengranularität) und Uhrenfragen.
- Ohne request-seitige Auswertung entfällt der Löwenanteil der Ersparnis.

## Bekannte Verwendungen

GitHub API v3 (ETags in den meisten Antworten, `304` ohne Rate-Limit-Anrechnung); Google Calendar
API (ETags beim Lesen *und* Ändern); Salesforce REST API (ETags oder `If-Modified-Since` /
`If-Unmodified-Since`). Frameworks unterstützen das Pattern direkt: Spring mit
`ShallowEtagHeaderFilter`, die `CacheControl`-Bibliothek des Play Frameworks. RFC 7232 ist die
maßgebliche HTTP-Spezifikation.

## Verwandte Patterns

- [Pagination](Pagination.md) — reduziert Volumen durch Zerlegung statt durch Auslassung.
- [Wish List](WishList.md) / [Wish Template](WishTemplate.md) — kombinierbar: Sie bestimmen, *was*
  gesendet wird, wenn die Bedingung erfüllt ist.
- [Request Bundle](RequestBundle.md) — kombinierbar, aber laut Quelle nur bei nachweisbarem Gewinn.
- [Metadata Element](../structure/MetadataElement.md) — die Bedingung reist als Kontroll-Metadatum.
- [Rate Limit](RateLimit.md) — profitiert unmittelbar.
- [Information Holder Resource](../responsibility/InformationHolderResource.md) — der typische Adressat.
- [Retrieval Operation](../responsibility/RetrievalOperation.md) — die typische konditionale Operation.

## Bezug zu Kubernetes / KRM

KRM erfüllt dieses Pattern **sehr stark, aber auf einem eigenen Weg**. HTTP-ETags spielen fast keine
Rolle: `ETag`/`If-None-Match` und `304 Not Modified` gibt es im Kubernetes-API-Server nur an zwei
Stellen — bei der Aggregated Discovery
(`apiserver/pkg/endpoints/discovery/aggregated/etag.go`, SHA-512-Hash über das
`APIGroupDiscoveryList`) und in den OpenAPI-v2/v3-Handlern von `kube-openapi`. Für Ressourcenpfade
wie `/api/v1/pods` existiert kein ETag. An seine Stelle tritt durchgängig die
**`resourceVersion`** — ein Metadata Element, das an jedem Objekt (`metadata.resourceVersion`),
an jeder Liste (`ListMeta.resourceVersion`) und in `ListOptions` auftaucht und intern der
etcd-MVCC-Revision entspricht. Sie ist strikt opak: Clients dürfen sie nur vergleichen bzw.
unverändert zurückschicken, nicht interpretieren oder ordnen.

Auf der **Schreibseite** ist das Pattern zwingend, nicht optional. Ein `update` (`PUT`) mit gesetzter
`metadata.resourceVersion` ist ein `If-Match`: Stimmt sie nicht mit der aktuellen überein, antwortet
der Server mit `409 Conflict` und `Status.Reason: "Conflict"` (`metav1.StatusReasonConflict`) — das
Standardmuster für Optimistic Concurrency, auf dem sämtliche Controller-Retry-Schleifen
(`client-go/util/retry.RetryOnConflict`) beruhen.

```http
PUT /apis/apps/v1/namespaces/default/deployments/web

{
  "apiVersion": "apps/v1", "kind": "Deployment",
  "metadata": { "name": "web", "resourceVersion": "8675309" },   // <- das If-Match von KRM
  "spec": { "replicas": 5 }
}

HTTP/1.1 409 Conflict

{
  "kind": "Status",
  "apiVersion": "v1",
  "status": "Failure",
  "message": "Operation cannot be fulfilled on deployments.apps \"web\": the object has been modified; please apply your changes to the latest version and try again",
  "reason": "Conflict",
  "details": { "group": "apps", "kind": "deployments", "name": "web" },
  "code": 409
}
```

Analog erlaubt `metav1.Preconditions` in `DeleteOptions` eine Bedingung über `uid` und/oder
`resourceVersion`, sodass sich ein Delete nicht versehentlich auf ein neu erzeugtes, gleichnamiges
Objekt bezieht. Server-Side Apply ersetzt den
Vergleich durch Feld-Eigentümerschaft und meldet Konflikte ebenfalls als `409`, überwindbar per
`force`.

Auf der **Leseseite** steuert `ListOptions.resourceVersion` zusammen mit `resourceVersionMatch`
(`NotOlderThan` | `Exact`) die zulässige Aktualität. `resourceVersion=0` bzw.
`resourceVersionMatch=NotOlderThan` erlaubt dem API-Server, aus dem Watch-Cache statt aus etcd zu
antworten — das ist keine Bandbreitenersparnis, sondern eine Ersparnis an teuren Quorum-Reads, also
genau die „Provider-Last"-Force. Wird eine `resourceVersion` angefordert, die der Watch-Cache nicht
mehr vorhält, kommt `410 Gone` mit `Status.Reason: "Expired"` und der Meldung
`too old resource version` (`storage/cacher/watch_cache.go`) — die KRM-Entsprechung dazu, dass ein
Validator ungültig geworden ist.

```http
# "irgendein Stand, aber nicht älter als der, den ich schon kenne":
GET /api/v1/namespaces/default/pods?resourceVersion=8675309&resourceVersionMatch=NotOlderThan

# Sonderfall resourceVersion=0: "beliebig alt, Hauptsache aus dem Cache" — kein Quorum-Read.
GET /api/v1/namespaces/default/pods?resourceVersion=0
```

Der eigentliche Unterschied zu MAP ist jedoch, dass Kubernetes das Problem eine Ebene höher löst:
**`watch` ersetzt Polling überhaupt**. Statt periodisch konditional zu fragen „hat sich etwas
geändert?", öffnet der Client einen Stream ab einer bekannten `resourceVersion` und bekommt nur noch
Deltas. `allowWatchBookmarks=true` liefert dazu `BOOKMARK`-Events, die die aktuelle
`resourceVersion` fortschreiben, ohne Objektdaten zu übertragen — die extremste Form eines
„nichts Neues"-Signals, sparsamer als jedes `304`, weil sie unaufgefordert kommt und keinen Request
kostet. `sendInitialEvents=true` (WatchList, `client-go`-Feature `WatchListClient`) verschmilzt
schließlich den initialen Vollabgleich mit dem Stream und markiert dessen Ende über ein Bookmark
mit der Annotation `k8s.io/initial-events-end`.

```http
GET /api/v1/namespaces/default/serviceaccounts?watch=true&allowWatchBookmarks=true
    &sendInitialEvents=true&resourceVersionMatch=NotOlderThan

{"type":"ADDED","object":{"kind":"ServiceAccount","apiVersion":"v1","metadata":{"name":"foobar","namespace":"default","resourceVersion":"217","uid":"a1453396-..."}}}
{"type":"BOOKMARK","object":{"kind":"ServiceAccount","apiVersion":"v1","metadata":{"resourceVersion":"964","annotations":{"k8s.io/initial-events-end":"true"}}}}
```

Das Bookmark-Objekt trägt nur `resourceVersion` und die Annotation — kein `spec`, kein `status`,
kein Name. Es sagt „bis hierher bist du synchron" und kostet dafür keinen Request. Das Ergebnis ist
ein Modell, in dem ein Controller nach einem einmaligen Sync dauerhaft ohne wiederholte
Vollabfragen auskommt — eine Effizienzstufe,
die MAP unter *Conditional Request* nicht abbildet, weil das Pattern das Request-Response-Modell
nicht verlässt.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/ConditionalRequest)
