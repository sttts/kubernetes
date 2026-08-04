---
title: Retrieval Operation
kategorie: Responsibility
unterkategorie: Operation Responsibilities
quelle: https://microservice-api-patterns.org/patterns/responsibility/operationResponsibilities/RetrievalOperation
---

# Retrieval Operation

*a.k.a.* Read-Only Operation, State Lookup Operation, Query, Data Extractor

**Kurzform:** Eine lesende Operation, die Provider-seitigen Zustand auswertet und als maschinenlesbaren
Ergebnisbericht liefert — mit Such-, Filter- und Formatierungsmöglichkeiten in der Signatur, aber ohne
jede Zustandsänderung.

## Kontext

Ein Endpunkt existiert, aber seine Operationen decken den Integrationsbedarf noch nicht ab: Clients
brauchen lesenden Zugriff auf große Mengen strukturierter, ggf. aggregierter Daten. Diese Daten sind
typischerweise anders strukturiert als das zugrunde liegende Domänenmodell — etwa bezogen auf ein
Zeitintervall oder ein Teilgebiet (Produktkategorie, Kundenprofilgruppe). Der Informationsbedarf tritt
ad hoc oder regelmäßig auf (wöchentlich, monatlich, quartalsweise).

## Problem

Wie können Informationen eines entfernten Providers abgerufen werden, um einen Informationsbedarf zu
decken oder clientseitige Weiterverarbeitung zu ermöglichen? Teilprobleme laut Quelle:

- Wie werden Datenmodellunterschiede überbrückt und Daten aggregiert bzw. mit anderen Quellen kombiniert?
- Wie beeinflusst der Client Umfang und Auswahlkriterien des Ergebnisses?
- Wie wird der Zeitrahmen für Berichte angegeben?

## Forces

- **Veracity, Variety, Velocity, Volume** — vier der fünf Big-Data-Vs. Datenqualität, Heterogenität,
  Änderungsrate und Menge kollidieren mit einem einfachen Kontrakt.
- **Lastmanagement:** Eine offene Query-Schnittstelle ist ein offenes Scheunentor für teure Anfragen.
- **Netzwerkeffizienz gegen Datensparsamkeit:** Wenige große Antworten oder viele kleine?

## Lösung

Eine Operation `ro: (in, S) -> out` wird dem Endpunkt hinzugefügt — meist einer
[Information Holder Resource](InformationHolderResource.md). Sie liefert einen Ergebnisbericht mit
maschinenlesbarer Repräsentation der angefragten Information. Der Signatur werden Such-, Filter- und
Formatierungsparameter beigegeben.

## Beispiel

Aus der Quelle, mit Paginierungsparametern:

```java
// curl http://localhost:8080/claims?limit=10&offset=0
@GET
public ClaimsDTO listClaims(
  @DefaultValue("3") @QueryParam("limit")  Integer limit,
  @DefaultValue("0") @QueryParam("offset") Integer offset,
  @QueryParam("orderBy") String orderBy
) { ... }
```

## Konsequenzen

**Vorteile:**

- Klare Trennung von Lese- und Schreibverantwortung; Endpunkte, die ausschließlich Retrieval Operations
  anbieten, bilden das Query Model in CQRS.
- Cachebar und beliebig wiederholbar, weil seiteneffektfrei.
- Filter- und Wunschparameter verlagern die Selektion auf den Provider und sparen Netzwerkbandbreite.

**Nachteile / Kosten:**

- Ausdrucksstarke Queries (GraphQL, SOQL, restSQL) verlagern Last unvorhersehbar zum Provider und
  erschweren Kapazitätsplanung.
- Aggregierende Berichte koppeln an das Provider-Datenmodell, wenn keine eigene Berichtsstruktur
  definiert wird.
- Große Ergebnismengen erzwingen zusätzliche Mechanik: [Pagination](../quality/Pagination.md),
  Größenlimits, Detailstufen.

## Bekannte Verwendungen

- eBays Traffic-Report-Operation; `files.list` der Slack Web API; Force.com mit SOQL.
- `Cargo find(TrackingId)` und `List findAll()` im Cargo Repository der DDD Sample Application.
- Open Weather Map mit zahlreichen Lookup-Parametern; Zefix (`searchByName`) der Schweizer eGov-Initiative.
- Terravis' Parzellen-Abfrage als Fassade über föderierte Grundbuchdaten — mit hartem Limit von zehn
  Parzellen pro Anfrage zum Lastschutz und wählbaren Detailstufen (full, partial, full history).
- GraphQL-Queries (nicht die Mutations) und restSQL als besonders flexible Instanzen.

## Verwandte Patterns

- [Computation Function](ComputationFunction.md) — ebenfalls zustandsneutral, bezieht aber alle Daten
  aus der Anfrage statt aus Provider-Zustand.
- [State Transition Operation](StateTransitionOperation.md) und
  [State Creation Operation](StateCreationOperation.md) — die schreibenden Schwestern.
- [Pagination](../quality/Pagination.md) — fast immer nötig bei Sammlungsabfragen.
- [Wish List](../quality/WishList.md) und [Wish Template](../quality/WishTemplate.md) — steuern die
  Detailtiefe der Antwort.
- [Embedded Entity](../quality/EmbeddedEntity.md) vs.
  [Linked Information Holder](../quality/LinkedInformationHolder.md) — die zwei Antworten auf die Frage,
  wie viel Nachbardaten mitkommen.
- [Conditional Request](../quality/ConditionalRequest.md) — vermeidet Übertragung unveränderter Daten.
- [Master Data Holder](MasterDataHolder.md), [Reference Data Holder](ReferenceDataHolder.md),
  [Operational Data Holder](OperationalDataHolder.md) — die Endpunkte, die typischerweise gelesen werden.

## Bezug zu Kubernetes / KRM

Von den vier Operationsverantwortlichkeiten ist *Retrieval Operation* diejenige, die KRM am
vollständigsten und wörtlichsten umsetzt. Drei der acht Standardverben sind rein lesend: `get` (ein
Objekt über `metadata.namespace` und `metadata.name`), `list` (eine Collection) und `watch` (ein Strom
von Änderungen). Alle Ressourcen — Core-Typen wie CRDs — bieten sie mit identischer Semantik, was
generische Clients wie `kubectl get` und die Informer-Maschinerie in `client-go` überhaupt erst möglich
macht.

Die vom Pattern geforderten Such- und Filterfähigkeiten stecken in `metav1.ListOptions`
(`staging/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go`): `labelSelector` filtert über
`metadata.labels`, `fieldSelector` über einen pro Ressource fest definierten, kleinen Satz indizierter
Felder (etwa `spec.nodeName` bei Pods). Bewusst gibt es *keine* allgemeine Query-Sprache — das ist die
KRM-Antwort auf die Force „Lastmanagement“: Der Provider gibt nur Filter frei, die er effizient bedienen
kann, statt beliebige Ausdrücke zuzulassen. Aggregierende oder umstrukturierende Berichte im Sinne des
Patterns gibt es folglich nicht; wer sie braucht, baut sie clientseitig auf Informer-Caches oder als
eigene aggregierte API.

Für [Pagination](../quality/Pagination.md) definiert `ListOptions` `limit` und `continue`; die Antwort
trägt `metadata.continue` und `metadata.remainingItemCount`. Konsistenz steuern `resourceVersion` und
`resourceVersionMatch` (`NotOlderThan`, `Exact`) — damit lässt sich zwischen billigem Cache-Read und
teurem quorum-konsistentem Read wählen, eine explizite Ausformulierung des Trade-offs, den das Pattern
nur andeutet. Alle diese Parameter sind Query-Parameter derselben, uniformen Collection-URL:

```http
GET /api/v1/namespaces/default/pods?labelSelector=app%3Dweb&fieldSelector=spec.nodeName%3Dnode-1&limit=500&resourceVersionMatch=NotOlderThan&resourceVersion=0
Accept: application/json

HTTP/1.1 200 OK

{
  "kind": "PodList",
  "apiVersion": "v1",
  "metadata": {
    "resourceVersion": "184203",
    "continue": "eyJ2IjoibWV0YS5rOHMuaW8vdjEiLCJydiI6MTg0MjAzLCJzdGFydCI6...",
    "remainingItemCount": 1240
  },
  "items": []
}
```

`resourceVersion=0` zusammen mit `resourceVersionMatch=NotOlderThan` ist dabei die Bitte „nimm, was
im Cache liegt“; `continue` ist ein opakes Token, kein Offset — der Client darf es nur zurückschicken,
nicht interpretieren.

Die eigentliche Besonderheit ist `watch`. Statt periodischem Polling liefert der API-Server einen
Ereignisstrom (`ADDED`, `MODIFIED`, `DELETED`, optional `BOOKMARK` via `allowWatchBookmarks`), und mit
`sendInitialEvents` lässt sich ein vollständiger Initialbestand gefolgt von einem synthetischen
Bookmark anfordern.

```http
GET /api/v1/namespaces/default/pods?watch=true&allowWatchBookmarks=true&sendInitialEvents=true&resourceVersionMatch=NotOlderThan&resourceVersion=

HTTP/1.1 200 OK
Transfer-Encoding: chunked

{"type":"ADDED","object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"web-7d9f2","resourceVersion":"184190"}}}
{"type":"BOOKMARK","object":{"kind":"Pod","apiVersion":"v1","metadata":{"resourceVersion":"184203","annotations":{"k8s.io/initial-events-end":"true"}}}}
{"type":"MODIFIED","object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"web-7d9f2","resourceVersion":"184287"}}}
```

Das `BOOKMARK`-Objekt trägt keine Nutzdaten, nur eine `resourceVersion`: es sagt „ab hier bist du
auf Stand“ und erlaubt einem abgerissenen Watch, ohne erneutes `list` wieder anzudocken. Damit wird
aus der Pull-Operation des Patterns ein langlebiges Abonnement, das
dieselben Selektoren nutzt. Für Level-Triggered Reconciliation ist das konstitutiv: Ein Controller
listet einmal, beobachtet danach nur noch Deltas und hält lokal einen konsistenten Spiegel. KRM
erfüllt das Pattern also (a) im Kern genauso, ergänzt es aber um eine Streaming-Dimension und
beschneidet bewusst die Ausdrucksstärke der Filter.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/operationResponsibilities/RetrievalOperation)
