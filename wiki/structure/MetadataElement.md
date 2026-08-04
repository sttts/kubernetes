---
title: Metadata Element
kategorie: Structure
unterkategorie: Element Stereotypes
quelle: https://microservice-api-patterns.org/patterns/structure/elementStereotypes/MetadataElement
---

# Metadata Element

*a.k.a.* *Data about/on Data*, *Semantic Annotation*, *Embedded API Microdescription*,
*Metadata Representation*

**Kurzform:** Nachrichten werden mit zusätzlicher, erklärender Information angereichert, damit
Empfänger den Inhalt korrekt und effizient verarbeiten können, ohne Annahmen über die
Datensemantik fest zu verdrahten.

## Kontext

Die Repräsentationen von Request und Response sind definiert, ggf. mit
[Atomic Parameter](AtomicParameter.md), [Atomic Parameter List](AtomicParameterList.md),
[Parameter Tree](ParameterTree.md) und [Parameter Forest](ParameterForest.md). Für eine
korrekte und effiziente Verarbeitung brauchen Empfänger jedoch mehr als Name und Typ der
Nutzdaten.

## Problem

Wie können Nachrichten so angereichert werden, dass Empfänger den Inhalt richtig interpretieren
— ohne dass Annahmen über die Datensemantik im Client-Code hartkodiert werden müssen?

## Forces

- **Interoperabilität** — Empfänger unterschiedlicher Herkunft und Reife müssen dieselbe
  Nachricht deuten können.
- **Kopplung** — Metadaten sind selbst Vertragsbestandteil; zu spezifische Metadaten koppeln
  ebenso stark wie Nutzdaten.
- **Bedienkomfort vs. Laufzeiteffizienz (Nachrichtengröße)** — jedes Metadatenfeld kostet
  Bytes, Serialisierungszeit und Erklärungsbedarf.

## Lösung

Ein oder mehrere *Metadata Elements* einführen, die die übrigen Repräsentationselemente
erklären und ergänzen. Die Werte konsistent und vollständig befüllen und empfängerseitig
tatsächlich auswerten, um Verarbeitung zu steuern.

Die Quellseite unterscheidet drei Metadaten-Arten:

| Art | Zweck | Beispiele aus der Quelle |
|---|---|---|
| **Provenance Metadata** | Herkunft, Erzeugung, Gültigkeit | `regionCode`, Message Date, ETag |
| **Control Metadata** | steuert die weitere Verarbeitung | `nextPageToken`, `prevPageToken`, Sortier-/Filteroperatoren |
| **Aggregated Metadata** | Kennzahlen über die Nutzdaten | `totalResults`, `resultsPerPage`, `totalSize` |

## Beispiel

Ausschnitt einer YouTube-Data-API-Antwort (Quellseite) mit allen drei Arten:

```json
{
  "kind": "youtube#searchListResponse",
  "etag": "etag",
  "nextPageToken": "CBQQAA",
  "prevPageToken": "CAoQAQ",
  "regionCode": "CH",
  "pageInfo": {
    "totalResults": 1662,
    "resultsPerPage": 10
  },
  "items": [ "search Resource" ]
}
```

`pageInfo` ist ein kleiner [Parameter Tree](ParameterTree.md), der zwei aggregierte Metadaten
bündelt; die übrigen sind [Atomic Parameters](AtomicParameter.md).

## Konsequenzen

**Vorteile:**

- Empfänger können generisch verarbeiten, statt Semantik zu raten oder zu hartkodieren.
- Steuerinformation (Paging, Caching, Konsistenz) wird explizit und maschinenauswertbar.
- Herkunftsangaben ermöglichen Auditierbarkeit und Nachvollziehbarkeit.

**Nachteile / Kosten:**

- Größere Nachrichten und mehr Serialisierungsaufwand.
- Metadaten müssen konsistent gepflegt werden; falsche Metadaten sind schlimmer als keine.
- Uneinheitliche Platzierung über Endpoints hinweg erzeugt Wildwuchs — dagegen hilft
  [Context Representation](ContextRepresentation.md).

## Bekannte Verwendungen

Die Quellseite nennt: `metadata`-Parameter der Facebook Graph API; Twitter REST API; YouTube
Data API; Force.com/SOQL mit `done` (Control) und `totalSize` (Aggregated); HTTP-`ETag` nach
RFC 7232 als zugleich Control- und Provenance-Metadatum; Protokoll-Header in HTTP und TCP/IP;
den Schweizer E-Government-Standard GBDBS (u. a. Datum der letzten Aktualisierung); `APIs.json`;
die „meta objects“ von JSON:API.

## Verwandte Patterns

- [Data Element](DataElement.md) — Oberbegriff; jedes *Metadata Element* ist ein *Data Element*,
  aber nicht umgekehrt.
- [Id Element](IdElement.md), [Link Element](LinkElement.md) — können selbst als Metadaten
  gelten oder von Metadaten begleitet werden.
- [Context Representation](ContextRepresentation.md) — API-weiter, technologieunabhängiger
  Standardort und Standardaufbau für (Control-)Metadaten.
- [Pagination](../quality/Pagination.md) — beruht direkt auf Control- und Aggregated Metadata.
- [Conditional Request](../quality/ConditionalRequest.md) — ETag/Last-Modified als Metadatum.
- [Rate Limit](../quality/RateLimit.md) — Zähler und Kontingente als Metadaten.
- [API Description](../foundation/APIDescription.md) — dort werden Metadatenfelder erklärt.
- [Version Identifier](../evolution/VersionIdentifier.md) — Sonderfall von Provenance-Metadaten.

Außerhalb von MAP: *Format Indicator*, *Message Expiration*, *Correlation Identifier*,
*Return Address*, *Message Filter*/*Selector* und *Aggregator* aus den *Enterprise Integration
Patterns* (Hohpe/Woolf) sowie das *Context Object* aus den *Remoting Patterns*
(Voelter/Kircher/Zdun).

## Bezug zu Kubernetes / KRM

Hier ist KRM ungewöhnlich radikal: statt Metadaten ad hoc pro Endpoint zu verstreuen, gibt es
genau **einen** uniformen Block `metadata` (`metav1.ObjectMeta`) — und zwar identisch über
*alle* Ressourcentypen hinweg, eingebaute wie CRDs. Das ist eine viel stärkere Ausprägung des
Musters, als MAP verlangt, und faktisch die konsequenteste denkbare Umsetzung von
[Context Representation](ContextRepresentation.md).

Der Block enthält Provenance-Metadaten, Identität (siehe [Id Element](IdElement.md)),
Nebenläufigkeitssteuerung, Beziehungen (siehe [Link Element](LinkElement.md)),
Lebenszyklus-Steuerung sowie zwei offene Key-Value-Maps:

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: web-6d4f8b7c9
  namespace: default
  uid: 8f2a1c0e-3b7d-4c11-9a55-2f0e6d7b1c34   # <- global eindeutig, vom Server vergeben
  resourceVersion: "1250043"                  # <- opak; Optimistic Concurrency und watch
  generation: 3                               # <- Zähler der spec-Änderungen
  creationTimestamp: "2026-08-04T09:12:33Z"
  labels:                                     # <- selektierbar, indiziert, validiert
    app: web
    pod-template-hash: 6d4f8b7c9
  annotations:                                # <- nicht selektierbar, beliebiger Inhalt
    deployment.kubernetes.io/revision: "4"
    deployment.kubernetes.io/desired-replicas: "3"
  finalizers:
    - example.com/cleanup
  ownerReferences:
    - apiVersion: apps/v1
      kind: Deployment
      name: web
      uid: 1b0c9d55-7e21-4a3f-8de0-5c9a2b6f4e77
      controller: true
      blockOwnerDeletion: true
  managedFields:
    - manager: kube-controller-manager
      operation: Update
      apiVersion: apps/v1
      time: "2026-08-04T09:12:33Z"
      fieldsType: FieldsV1
      fieldsV1:                               # <- Provenance auf Feldebene
        f:metadata:
          f:labels:
            .: {}
            f:app: {}
        f:spec:
          f:replicas: {}
    - manager: kube-controller-manager
      operation: Update
      subresource: status                     # <- getrennter Eintrag pro Subresource
      # ... fieldsV1
```

Die Trennung von `labels` und `annotations` ist bewusst:
Labels sind Selektor-Ziel (`spec.selector` bei Deployment, Service, NetworkPolicy) und in ihrer
Länge und Syntax validiert, Annotations sind ein Freiraum für Werkzeuge und Controller.

`managedFields` (`metav1.ManagedFieldsEntry` mit `manager`, `operation`, `apiVersion`, `time`,
`fieldsType`, `fieldsV1`, `subresource`) ist Provenance-Metadatum auf Feldebene: es hält fest,
welcher Akteur welches einzelne Feld zuletzt gesetzt hat, und ist die Grundlage von
Server-Side Apply und dessen Konflikterkennung. Ein Äquivalent dazu kennt MAP nicht.

Control- und Aggregated Metadata für Listen sitzen in einem eigenen, ebenfalls uniformen Block
`metav1.ListMeta`.

```yaml
apiVersion: v1
kind: PodList
metadata:                     # <- metav1.ListMeta, nicht ObjectMeta
  resourceVersion: "1250043"  # <- Provenance + Startpunkt für ein anschließendes watch
  continue: eyJ2IjoibWV0YS5rOHMuaW8vdjEi...   # <- Control Metadata: Cursor der nächsten Seite
  remainingItemCount: 1200    # <- Aggregated Metadata
items:
  # ... 500 Pods (limit=500)
```

Auf Transportebene ergänzen Header wie `ETag`-Semantik
über `resourceVersion` und die Warnungen des `Warning`-Headers das Bild.

```http
HTTP/1.1 200 OK
Content-Type: application/json
Warning: 299 - "policy/v1beta1 PodDisruptionBudget is deprecated in v1.21+, unavailable in v1.25+; use policy/v1 PodDisruptionBudget"
```

Der bemerkenswerte Unterschied zu MAP: Metadaten sind in KRM nicht optionale Anreicherung,
sondern Voraussetzung der Maschinerie. Generischer Code — `kubectl`, Informer, Garbage
Collector, Admission — funktioniert nur, weil `metadata` bei jeder Ressource dieselbe Form hat
und über `metav1.Object` typunabhängig zugreifbar ist.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/elementStereotypes/MetadataElement)
