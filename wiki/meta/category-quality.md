---
title: Kategorie Quality
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/quality
---

# Kategorie: Quality Patterns

## Worum es geht

API-Qualität hat viele Dimensionen — Zuverlässigkeit, Performance, Sicherheit, Skalierbarkeit —,
die untereinander in Konflikt stehen und zusätzlich gegen wirtschaftliche Kräfte (Kosten,
Time-to-Market) abgewogen werden müssen. Die Kategorie behandelt zwei Fragen: Wie erreicht ein
Provider ein bestimmtes Qualitätsniveau bei kosteneffizientem Ressourceneinsatz, und wie werden
die dabei eingegangenen Trade-offs kommuniziert und abgerechnet?

Die Quellseite macht eine nützliche Beobachtung zur *Entscheidungsgranularität*: Die meisten
Entscheidungen dieser Kategorie gelten für Kombinationen aus Client-Gruppe und API (etwa „alle
Freemium-Clients"), nur die Vermeidung unnötiger Datenübertragung wird pro Operation entschieden.

## Inhalt

### Reference Management — eingebettet oder verlinkt?

| Pattern | Einzeiler |
|---|---|
| [Embedded Entity](../quality/EmbeddedEntity.md) | Die Daten am Zielende einer Beziehung werden direkt in die Repräsentation der Quelle eingebettet, um Folgeaufrufe zu vermeiden. |
| [Linked Information Holder](../quality/LinkedInformationHolder.md) | Statt der Daten wird ein [Link Element](../structure/LinkElement.md) auf einen eigenen Endpunkt mitgegeben, damit die Nachricht klein bleibt. |

Diese beiden sind das kanonische Gegensatzpaar von MAP: wenige große Nachrichten gegen viele
kleine. Der empfohlene Übergang läuft in beide Richtungen und wird durch gemessenes Verhalten
ausgelöst, nicht durch Vorabdesign.

### Data Transfer Parsimony — sparsam übertragen

| Pattern | Einzeiler |
|---|---|
| [Pagination](../quality/Pagination.md) | Große Ergebnismengen werden in Chunks ausgeliefert, mit Angabe von Gesamt-/Restmenge und Verweis auf den nächsten Chunk. |
| [Wish List](../quality/WishList.md) | Der Client zählt im Request die gewünschten Datenelemente auf; der Provider liefert nur diese („Response Shaping"). |
| [Wish Template](../quality/WishTemplate.md) | Für geschachtelte Daten: Der Request spiegelt die Hierarchie der Antwort, markierte Zweige werden geliefert. |
| [Conditional Request](../quality/ConditionalRequest.md) | Metadaten machen den Request bedingt; der Provider verarbeitet ihn nur, wenn die Bedingung erfüllt ist (z. B. Daten geändert). |
| [Request Bundle](../quality/RequestBundle.md) | Mehrere unabhängige Requests werden in einer Nachricht gebündelt, samt Identifikatoren und Zähler. |

Die fünf greifen an verschiedenen Stellen an: *Pagination* begrenzt die Anzahl der Elemente,
*Wish List*/*Wish Template* die Breite bzw. Tiefe je Element, *Conditional Request* die Anzahl
der Übertragungen, *Request Bundle* die Anzahl der Round-Trips.

### Quality Management and Governance

| Pattern | Einzeiler |
|---|---|
| [Rate Limit](../quality/RateLimit.md) | Eine durchgesetzte Nutzungsobergrenze schützt den Provider vor übermäßiger Nutzung durch einzelne Clients. |
| [Pricing Plan](../quality/PricingPlan.md) | Nutzungsmetriken je Operation ermöglichen Messung und Abrechnung der API-Nutzung. |
| [Service Level Agreement](../quality/ServiceLevelAgreement.md) | Strukturierte, testbare Service-Level-Objectives samt Konsequenzen bei Nichteinhaltung. |

## Verwandte Wiki-Seiten

- [Kategorie Structure](category-structure.md) — die Nachrichtenstrukturen, die hier optimiert werden
- [Patterns nach Qualitätsattribut](navigation-byquality.md) — Einstieg über Performance, Sicherheit, …
- [Cheat Sheet](cheatsheet.md) — Abschnitt „Continuous API Improvement"
- [Tutorials](tutorials.md) — Tutorial 1 behandelt genau diese Kategorie

## Bezug zu Kubernetes / KRM

Die *Data-Transfer-Parsimony*-Patterns sind in KRM teils vorhanden, teils bewusst durch etwas
anderes ersetzt. [Pagination](../quality/Pagination.md) existiert direkt: `ListOptions.limit` und
`ListOptions.continue`, mit `ListMeta.continue` und `remainingItemCount` in der Antwort — also
Cursor-basiert, nicht Offset-basiert, weil der etcd-Snapshot einen Fortsetzungstoken natürlich
hergibt. Serverseitiges Response Shaping im Sinne einer [Wish List](../quality/WishList.md) gibt
es nur in grober Form: `PartialObjectMetadata` über den `Accept`-Header liefert ausschließlich
`metadata`, und die `Table`-Konvertierung liefert vordefinierte Spalten (bei CRDs über
`additionalPrinterColumns` konfigurierbar). Eine freie Feldauswahl wie bei GraphQL oder eine
[Wish Template](../quality/WishTemplate.md) kennt KRM nicht — `--field-selector` filtert *Objekte*,
nicht *Felder*, und ist auf wenige indizierte Felder beschränkt.

Der [Conditional Request](../quality/ConditionalRequest.md) hat in KRM eine ungewöhnliche
Ausprägung. HTTP-Caching mit `ETag`/`If-None-Match` wird nicht verwendet; stattdessen ist
`metadata.resourceVersion` das universelle Konditionalitäts-Metadatum: als optimistische Sperre
beim Update (Konflikt → HTTP 409), als Startpunkt beim `watch`, und als „nicht neuer als"-Hinweis
beim `list` (`resourceVersion=0` erlaubt eine Antwort aus dem Cache). Vor allem aber ersetzt das
Watch-Modell das wiederholte bedingte Pollen komplett: Der Client konsumiert einen Stream von
Änderungsereignissen statt periodisch zu fragen, ob sich etwas geändert hat. Das ist eine deutlich
weitergehende Lösung derselben Kraft.

Ein [Request Bundle](../quality/RequestBundle.md) fehlt: Es gibt keinen Batch-Endpunkt. Am
nächsten kommen `deletecollection` (ein Verb über viele Objekte) und `kubectl apply -f` mit
mehreren Dokumenten, das aber clientseitig in Einzelaufrufe zerlegt wird. Der Grund liegt im
Modell: Level-getriggerte Reconciliation macht Aufrufbündelung weniger wichtig, weil kein
Aufruf zeitkritisch einzeln quittiert werden muss.

Beim Governance-Teil ist der [Rate Limit](../quality/RateLimit.md) sehr gut vertreten: API
Priority and Fairness (`flowcontrol.apiserver.k8s.io/v1` mit `FlowSchema` und
`PriorityLevelConfiguration`) klassifiziert eingehende Anfragen nach Nutzer und Ressource, weist
sie Prioritätsstufen mit eigenen Concurrency-Anteilen zu und antwortet bei Überlast mit HTTP 429
und `Retry-After`. Zusätzlich drosseln Clients sich selbst (`rest.Config.QPS`/`Burst`). Ein
[Pricing Plan](../quality/PricingPlan.md) existiert in KRM nicht — das Gegenstück ist
`ResourceQuota`/`LimitRange`, das Ressourcenverbrauch begrenzt statt ihn abzurechnen. Ein
[Service Level Agreement](../quality/ServiceLevelAgreement.md) als API-Artefakt gibt es ebenfalls
nicht; Kubernetes definiert SLOs nur projektintern (etwa die Skalierbarkeits-SLOs für
API-Aufruflatenzen), nicht als Vertragsobjekt gegenüber Clients.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality)
