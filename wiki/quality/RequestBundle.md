---
title: Request Bundle
kategorie: Quality
unterkategorie: Data Transfer Parsimony
quelle: https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/RequestBundle
---

# Request Bundle

*a.k.a.* Request Batch, Request Deck, Bulk Request-Response

**Kurzform:** Mehrere unabhängige Requests werden in einer einzigen Nachricht zusammengefasst und
gemeinsam übertragen; Metadaten identifizieren die Einzelanfragen und ihre Antworten.

## Kontext

Ein Endpoint mit einer oder mehreren Operationen existiert. Der Provider beobachtet, dass Clients
viele kleine Anfragen stellen und für jede eine eigene Antwort erhalten. Diese geschwätzigen
Interaktionssequenzen („chatty interactions") schaden Durchsatz und Skalierbarkeit — der
Protokoll-, Verbindungs- und Serialisierungsoverhead dominiert die Nutzlast.

## Problem

Wie lässt sich die Zahl von Requests und Responses reduzieren, um die Kommunikationseffizienz zu
erhöhen?

## Forces

- **Komplexität von Endpoint-, Client- und Payload-Design:** Ein Bundle braucht Container-Struktur,
  Korrelation und eine mehrstufige Fehlerbehandlung.
- **Korrektheit von Reporting und Abrechnung:** Ein Bundle mit n Elementen — ist das ein Aufruf oder
  sind es n?
- **Latenz:** Ein Bundle antwortet erst, wenn das langsamste Element fertig ist.
- **Durchsatz:** Der eigentliche Gewinn — weniger Roundtrips, weniger Overhead pro Nutzlastbyte.

## Lösung

Einen *Request Bundle* als Datencontainer definieren, der mehrere unabhängige Anfragen in einer
Nachricht versammelt. Metadaten ergänzen: Identifikatoren der Einzelanfragen und ein Zähler der
Bundle-Elemente. Request- und Response-Bundles bilden dabei
[Parameter Forests](../structure/ParameterForest.md) bzw.
[Parameter Trees](../structure/ParameterTree.md); die Identifikatoren folgen dem Muster
*Correlation Identifier* (Hohpe/Woolf 2003).

## Varianten

Die Quelle nennt die Ausprägung *Request Bundle with Single Bundled Response* — alle Teilantworten
kommen in einer gemeinsamen Antwortnachricht zurück. Die Alternative sind mehrere Einzelantworten
auf ein gebündeltes Request, die über Korrelations-IDs zugeordnet werden.

## Beispiel

Im Lakeside-Mutual-*Customer-Core*-Service werden mehrere Kunden über eine
[Atomic Parameter List](../structure/AtomicParameterList.md) von
[Id Elements](../structure/IdElement.md) in einem Aufruf angefordert:

```
curl http://localhost:8080/customers/ce4btlyluu,rgpp0wkpec
```

```json
{ "customers": [ { "customerId": "ce4btlyluu", ... }, { "customerId": "rgpp0wkpec", ... } ] }
```

## Konsequenzen

**Vorteile:**

- Drastisch weniger Roundtrips; Marshalling-, Verbindungs- und Header-Overhead amortisieren sich.
- Besonders wirksam bei hoher Netzlatenz und bei Batch-artigen Verarbeitungen.
- Wirkt positiv auf ein [Rate Limit](RateLimit.md).

**Nachteile / Kosten:**

- Fehlerbehandlung wird zweistufig: Ein [Error Report](../structure/ErrorReport.md) muss pro
  Bundle-Element und nicht nur pro Aufruf gemeldet werden können.
- Teilerfolge sind der Normalfall — Atomarität ist *nicht* Teil des Patterns und müsste separat
  gebaut werden.
- Head-of-Line-Blocking: Die Antwortlatenz richtet sich nach dem langsamsten Element.
- Allamaraju (2010, Rezept 13) rät ausdrücklich von generischen Tunnel-Endpoints ab; besser sei ein
  Endpoint, der den konkreten Anwendungsfall direkt unterstützt.

## Bekannte Verwendungen

Google Calendar API mit *Batch Requests* über `multipart/mixed` an einen dedizierten `batch`-Endpoint
(Einzelanfragen als Body-Parts mit `Content-Type: application/http` und optionalem `Content-ID`);
*Batch Operations* in den Adidas API Guidelines, dort als
[Parameter Forest](../structure/ParameterForest.md) an denselben Endpoint statt an einen eigenen;
Service-Design-Richtlinien einer großen Schweizer Bank mit eigenem Fehlerformat pro Geschäftsobjekt;
zeilenweise Batchverarbeitung (Terravis); Bulk-Requests in SWITCH edu-ID.

## Verwandte Patterns

- [Pagination](Pagination.md) — das Gegenstück: zerlegt eine große Nachricht, statt viele kleine zu bündeln.
- [Conditional Request](ConditionalRequest.md) — kombinierbar, aber laut Quelle nur bei nachweisbarem Zusatzgewinn.
- [Wish List](WishList.md) / [Wish Template](WishTemplate.md) — Alternative, wenn der Endpoint eine [Information Holder Resource](../responsibility/InformationHolderResource.md) ist.
- [Parameter Forest](../structure/ParameterForest.md) / [Parameter Tree](../structure/ParameterTree.md) — Strukturen von Bundle-Request und -Response.
- [Metadata Element](../structure/MetadataElement.md) — Bundle-Zähler und Korrelations-IDs.
- [Error Report](../structure/ErrorReport.md) — muss elementgranular sein.
- [Rate Limit](RateLimit.md) — profitiert von weniger Aufrufen.

## Bezug zu Kubernetes / KRM

Hier hat KRM die **deutlichste Lücke** der gesamten Kategorie: Es gibt **kein generisches Batching**.
Keine `POST`-Operation nimmt eine `List` entgegen, es gibt keinen `batch`-Endpoint, kein
`multipart/mixed`, keine Multi-Objekt-Transaktion und folglich auch keine Atomaritätsgarantie über
mehrere Objekte. Jedes `create`, `update`, `patch` und `delete` adressiert genau **ein** Objekt unter
genau einem Pfad `/apis/<group>/<version>/namespaces/<ns>/<resource>/<name>`.

Konsequenz im Alltag: `kubectl apply -f manifests/` mit fünfzig Objekten erzeugt fünfzig einzelne
HTTP-Requests — `ApplyOptions.Run` iteriert schlicht über die eingelesenen Objekte und ruft pro
Objekt `applyOneObject` auf. Ein Helm-Release oder ein Operator, der ein Dutzend zusammengehöriger
Ressourcen erzeugt, tut dasselbe. Fehlschläge sind dabei partiell, und es gibt kein Rollback; die
Wiederherstellung des gewünschten Zustands ist Aufgabe der nächsten Reconciliation-Runde.

```console
$ kubectl apply -f manifests/
serviceaccount/web created
configmap/web-config created
deployment.apps/web created
service/web created
horizontalpodautoscaler.autoscaling/web created
# fünf Zeilen = fünf HTTP-Requests an fünf verschiedene Pfade.
# Schlägt der dritte fehl, bleiben die ersten beiden bestehen — kein Rollback.
```

Die einzige echte Sammeloperation auf der Schreibseite ist der Verb `deletecollection`
(`DELETE` auf den Collection-Pfad, registriert in `apiserver/pkg/endpoints/installer.go`,
in `client-go` als `DeleteCollection(ctx, deleteOpts, listOpts)`). Sie ist aber kein *Request Bundle*
im Sinne des Patterns: Sie bündelt nicht n unabhängige Anfragen, sondern wendet **eine** Operation
auf eine per Selektor bestimmte Menge an — und auch sie ist nicht atomar.

```http
DELETE /api/v1/namespaces/default/pods?labelSelector=app%3Dweb   # <- ListOptions in der Query

{
  "apiVersion": "v1",
  "kind": "DeleteOptions",              // <- DeleteOptions im Body
  "gracePeriodSeconds": 30,
  "propagationPolicy": "Background"
}
```

Auf der Leseseite ist `list` das Gegenstück: eine Anfrage, n Objekte — was allerdings eher der
Normalfall einer Collection-Ressource als ein Bundle ist.

Zwei Stellen setzen das Pattern jedoch **wirklich** um, beide im Meta-/Discovery-Bereich:

- **Aggregated Discovery:** Früher musste ein Client `/apis` abrufen und danach für jede
  Group/Version einen weiteren Request stellen — bei einem Cluster mit CRDs schnell dreistellig
  viele. Mit `Accept: application/json;g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList` liefert
  ein einziger Request die vollständige `APIGroupDiscoveryList` mit allen Gruppen, Versionen und
  Ressourcen. Das ist ein lehrbuchmäßiges Request Bundle — kombiniert mit
  [Conditional Request](ConditionalRequest.md), denn genau dieser Endpoint ist einer der wenigen im
  API-Server, die ETag und `If-None-Match` unterstützen.

  ```http
  GET /apis
  Accept: application/json;g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList
  If-None-Match: "9C7F1A5B..."          # <- SHA-512 über die gesamte Liste

  HTTP/1.1 304 Not Modified             # <- statt hunderter GETs: null Bytes Nutzlast
  ETag: "9C7F1A5B..."
  ```

- **`SelfSubjectRulesReview`** (`authorization.k8s.io/v1`) liefert alle Regeln, die für den
  aufrufenden Benutzer in einem Namespace gelten, in einer Antwort — statt n einzelner
  `SubjectAccessReview`-Aufrufe.

Warum verzichtet KRM ansonsten darauf? Weil das deklarative, level-triggered Modell den Bedarf
verschiebt: Der teure Pfad ist nicht das einmalige Schreiben vieler Objekte, sondern das
kontinuierliche Beobachten. Genau dieser Pfad ist über `watch` bereits gebündelt — ein einziger
Long-Running-Request transportiert beliebig viele Ereignisse für eine ganze Collection. Zudem
setzen Autorisierung (RBAC pro Ressource und Verb), Admission Webhooks (eine `AdmissionReview` pro
Objekt), Audit-Logging und `metadata.managedFields` sämtlich auf der Granularität des Einzelobjekts
auf. Ein echtes Bundle müsste all diese Maschinerie durchbrechen. Kubernetes hat sich damit bewusst
für Einheitlichkeit und Nachvollziehbarkeit statt für Anfrageeffizienz auf der Schreibseite
entschieden.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/RequestBundle)
