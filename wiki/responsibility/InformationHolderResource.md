---
title: Information Holder Resource
kategorie: Responsibility
unterkategorie: Endpoint Roles
quelle: https://microservice-api-patterns.org/patterns/responsibility/endpointRoles/InformationHolderResource
---

# Information Holder Resource

*a.k.a.* Generic Information Service, Data Entity Resource, Siloed/Isolated Data Holder

**Kurzform:** Ein API-Endpunkt, dessen Identität eine *Datenentität* ist. Er exponiert Create, Read,
Update, Delete und Suche auf dieser Entität und schützt sie gegen konkurrierende Zugriffe, ohne die
Implementierung offenzulegen.

## Kontext

Ein Domänenmodell, ein ER-Diagramm oder ein Glossar der Kernbegriffe liegt vor. Es enthält Entitäten
mit Identität, Lebenszyklus und Attributen, die einander referenzieren. Die Analyse zeigt, dass
strukturierte Daten an mehreren Stellen des verteilten Systems gebraucht werden und daher für
mehrere entfernte Clients zugänglich sein müssen. Ein Verstecken hinter fachlicher Logik ist nicht
möglich oder nicht sinnvoll — die Anwendung hat keinen ausgeprägten Workflow-Charakter.

## Problem

Wie können Domänendaten in einer API exponiert werden, ohne ihre Implementierung offenzulegen? Und
wie können mehrere Clients diese Entitäten nebenläufig lesen und ändern, ohne Datenintegrität und
Datenqualität zu gefährden?

## Forces

- **Modellierungsansatz** und dessen Kopplungswirkung: das API-Schema wird zum geteilten Vokabular
  und damit zum Evolutionsengpass.
- **Qualitätskonflikte**: Nebenläufigkeit, Konsistenz, Datenqualität und -integrität,
  Wiederherstellbarkeit, Verfügbarkeit, Veränderlichkeit vs. Unveränderlichkeit.
- **Sicherheit** — Daten sind ein deutlich attraktiveres Ziel als Aktionen.
- **Datenaktualität vs. Konsistenz** — Caching und Replikation kaufen Latenz mit Staleness.
- **Konformität zu Architekturprinzipien**: lose Kopplung, logische und physische
  Datenunabhängigkeit, unabhängige Deploybarkeit der Microservices.

Die spezifischen Kräfte werden in den fünf Verfeinerungen ausdifferenziert (siehe *Varianten*). Die
Leitentscheidung bleibt: datenorientierte oder aktivitätsorientierte Endpunktsemantik? Letztere
beschreibt [Processing Resource](ProcessingResource.md).

## Lösung

Füge der API einen *Information Holder Resource*-Endpunkt hinzu, der eine datenorientierte Entität
repräsentiert. Exponiere darauf Create-, Read-, Update-, Delete- und Suchoperationen. Koordiniere die
Aufrufe in der API-Implementierung so, dass die Entität geschützt bleibt.

## Varianten

Fünf Verfeinerungen sind eigene Patterns; sie unterscheiden sich in Veränderlichkeit, Lebensdauer,
Referenzrichtung und Eigentümerschaft:

| Variante | Lebensdauer | Änderbarkeit | Referenzen |
|---|---|---|---|
| [Operational Data Holder](OperationalDataHolder.md) | kurz bis mittel | häufig | viele ausgehende |
| [Master Data Holder](MasterDataHolder.md) | lang | selten | viele eingehende |
| [Reference Data Holder](ReferenceDataHolder.md) | sehr lang | für Clients keine | nur eingehende |
| [Data Transfer Resource](DataTransferResource.md) | temporär | durch Clients | — (Daten gehören den Clients) |
| [Link Lookup Resource](LinkLookupResource.md) | — | selten | hält nur Metadaten/Adressen |

## Beispiel

Der `Customer Core`-Microservice aus Lakeside Mutual exponiert Stammdaten. Die Operationen sind
daten- statt aktionsorientiert benannt; `getCustomer` unterstützt mit dem Query-Parameter `fields`
eine [Wish List](../quality/WishList.md).

```java
@RestController
@RequestMapping("/customers")
public class CustomerInformationHolder {
    @PutMapping(value = "/{customerId}/address")
    public ResponseEntity changeAddress(
        @PathVariable CustomerId customerId,
        @Valid @RequestBody AddressDto requestDto) { ... }

    @GetMapping(value = "/{ids}")
    public ResponseEntity getCustomer(
        @PathVariable String ids,
        @RequestParam(value = "fields", required = false, defaultValue = "") String fields) { ... }
}
```

## Konsequenzen

**Vorteile:**

- Uniforme, vorhersagbare Semantik: ein Schema und ein CRUD-Verbsatz statt vieler Einzelaktionen.
- Generische Werkzeuge (Clients, Caches, Codegeneratoren, HTTP-Zwischenschichten) tragen viel bei.
- Suche und Teilabruf lassen sich orthogonal ergänzen
  ([Pagination](../quality/Pagination.md), [Wish List](../quality/WishList.md)).

**Nachteile / Kosten:**

- Das exponierte Datenschema koppelt Clients an das Domänenmodell — Evolution wird teuer.
- Fachliche Invarianten liegen beim Client, wenn dieser mehrere Schreibzugriffe orchestrieren muss;
  Nebenläufigkeitsschutz (optimistisches Sperren) wird zur Pflicht.
- Anti-Pattern-Risiko: der Endpunkt degeneriert zum entfernten Datenbank-Wrapper („Siloed Data
  Holder“) und exponiert Implementierungsdetails.

## Bekannte Verwendungen

- Die Star Wars API mit `films`, `people`, `planets`, `species`, `starships`, `vehicles`.
- Dokumentendatenbanken mit HTTP-Schnittstelle, z.B. CouchDB (`GET /{db}/_all_docs`) und MongoDB.
- Master-Data-Management- und Product-Information-Management-Systeme per Definition; ebenso Konto-,
  Abrechnungs- und Währungscode-Endpunkte in Cloud-Provider-APIs.
- Storage-Angebote wie Dropbox, ownCloud und Amazon S3 (Dateisystem- bzw. Bucket-Abstraktionen).
- Enterprise-/Government-SOAs: das „Dynamic Interface“ eines Kernbankensystems, `GetParcelIndex` in
  Terravis (liefert `EGRID`-Grundstücksidentifikatoren), Open-Government-Data-Szenarien.

## Verwandte Patterns

- [Processing Resource](ProcessingResource.md) — komplementäre Semantik, die Alternative.
- [Operational Data Holder](OperationalDataHolder.md), [Master Data Holder](MasterDataHolder.md),
  [Reference Data Holder](ReferenceDataHolder.md) — Verfeinerungen nach Lebensdauer und Änderbarkeit.
- [Data Transfer Resource](DataTransferResource.md) — Spezialisierung für temporäre, clientseitig
  besessene Daten.
- [Link Lookup Resource](LinkLookupResource.md) — Spezialisierung, die nur Adressen hält; ihre
  Trefferliste zeigt oft auf *Information Holder Resources*.
- [State Creation Operation](StateCreationOperation.md) und [Retrieval Operation](RetrievalOperation.md)
  sind typisch; [Computation Function](ComputationFunction.md) und
  [State Transition Operation](StateTransitionOperation.md) sind erlaubt, aber seltener.
- [Embedded Entity](../quality/EmbeddedEntity.md) vs.
  [Linked Information Holder](../quality/LinkedInformationHolder.md) — Referenzen einbetten oder verlinken.
- [Id Element](../structure/IdElement.md), [Link Element](../structure/LinkElement.md),
  [Metadata Element](../structure/MetadataElement.md) — Bausteine der Repräsentation;
  [Conditional Request](../quality/ConditionalRequest.md) für Nebenläufigkeit und Caching.

## Bezug zu Kubernetes / KRM

KRM ist eine **kompromisslose Umsetzung genau dieses Patterns** (Fall a, in extremer Ausprägung):
praktisch jedes API-Objekt ist eine *Information Holder Resource*. Der Verbsatz ist für alle
Ressourcen identisch, und die Repräsentation folgt einem uniformen Rahmen. Das ist stärker
vereinheitlicht, als MAP es fordert: MAP lässt jedem Endpunkt sein eigenes Operationsprofil, KRM
zwingt alle in dasselbe.

```yaml
apiVersion: apps/v1                 # <- Rahmen: Gruppe/Version …
kind: Deployment                    # <- … und Typ, in jedem Objekt jeder Ressource
metadata:
  name: web
  namespace: default
  uid: a3f7c2e1-9b0d-4f2a-8c31-6d5e0b7a1f42
  generation: 7
  resourceVersion: "184203"         # <- optimistisches Sperren, Konflikt → HTTP 409
  # ... labels, annotations, creationTimestamp, ownerReferences, managedFields
spec:                               # gewünschter Zustand, geschrieben vom Client
  replicas: 3
  # ... selector, strategy, template
status:                             # beobachteter Zustand, geschrieben vom Controller
  observedGeneration: 7
  readyReplicas: 3
  # ... availableReplicas, updatedReplicas, conditions
```

Dass der Verbsatz wirklich uniform ist, sagt der Endpunkt selbst — die Discovery liefert ihn als
Datenfeld, nicht als Prosadokumentation:

```console
$ kubectl get --raw /apis/apps/v1 | jq '.resources[] | select(.name=="deployments" or .name=="deployments/status")'
{
  "name": "deployments",
  "singularName": "deployment",
  "namespaced": true,
  "kind": "Deployment",
  "verbs": ["create","delete","deletecollection","get","list","patch","update","watch"],
  "shortNames": ["deploy"],
  "categories": ["all"]
}
{
  "name": "deployments/status",
  "namespaced": true,
  "kind": "Deployment",
  "verbs": ["get","patch","update"]
}
```

Der zweite Eintrag ist die Pointe: `status` ist ein eigener Endpunkt mit verkürztem Verbsatz und
eigener RBAC-Ressource (`deployments/status`) — die Entität wird nicht als Ganzes geschützt, sondern
hälftig.

Zwei Ergänzungen gehen über das Pattern hinaus (Fall b):

- **`watch` als fünftes Verb.** Statt nur „lesen“ gibt es „lesen und ab jetzt Änderungen mitgeteilt
  bekommen“, gesteuert über `metadata.resourceVersion` und (mit Bookmark-Events) über
  `resourceVersion`-Fortschreibung. Damit löst KRM das *Data freshness versus consistency*-Force
  anders als klassische APIs: nicht durch Polling mit Caching, sondern durch einen dauerhaften
  Änderungsstrom mit einem Informer-Cache im Client.
- **`spec`/`status`-Trennung.** Der gewünschte Zustand gehört dem Client, der beobachtete Zustand
  dem Controller; `status` ist meist ein eigenes Subresource mit eigenem RBAC-Verb. Ein
  Kubernetes-Objekt ist deshalb nie nur ein Datensatz, sondern immer auch ein Auftrag an einen
  Controller — genau deshalb kommt KRM mit so wenigen
  [Processing Resources](ProcessingResource.md) aus.

Die Nebenläufigkeitskräfte des Patterns löst KRM explizit: optimistisches Sperren über
`metadata.resourceVersion` (Konflikt → HTTP 409), deklaratives Merging über Server-Side Apply mit
`metadata.managedFields` und Feldbesitzern, Lebenszyklus über `metadata.ownerReferences`,
`metadata.finalizers` und `metadata.deletionTimestamp`. Suche gibt es bewusst schwach: nur
Label- und Feldselektoren (`labelSelector`, `fieldSelector`) plus `limit`/`continue`-Pagination —
keine allgemeine Query-Sprache, weil jede Anfrage aus dem Watch-Cache bedienbar bleiben soll.

Die Schattenseite des Patterns ist in Kubernetes ebenfalls gut zu beobachten: das exponierte Schema
*ist* der Kontrakt, weshalb Evolution über Versionierung, Konvertierungs-Webhooks und die
Deprecation Policy teuer erkauft wird — vgl.
[Two in Production](../evolution/TwoInProduction.md) und
[Version Identifier](../evolution/VersionIdentifier.md).

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/endpointRoles/InformationHolderResource)
