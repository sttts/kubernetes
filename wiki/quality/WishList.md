---
title: Wish List
kategorie: Quality
unterkategorie: Data Transfer Parsimony
quelle: https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/WishList
---

# Wish List

*a.k.a.* Data Wish Enumeration, Partial Response Representation Request, Data Selection Profile

**Kurzform:** Der Client zählt im Request flach auf, welche Datenelemente er haben will; der
Provider liefert genau diese und nichts sonst („Response Shaping").

## Kontext

Ein Endpoint bedient mehrere Clients, die dieselben Operationen aufrufen, aber unterschiedliche
Informationsbedarfe haben: Ein Übersichtsscreen braucht drei Attribute, ein Reporting-Job den
vollständigen Datensatz. Ein einziges, maximales Response-Schema bedient beide schlecht.

## Problem

Wie teilt ein API-Client dem Provider **zur Laufzeit** mit, an welchen Daten er interessiert ist?

## Forces

- **Performance, Skalierbarkeit, Ressourcenverbrauch:** Nicht angeforderte Felder kosten
  Datenbankzugriffe, Serialisierung und Bandbreite.
- **Individuelle Informationsbedarfe:** Sonst müsste der Provider pro Client eine eigene Operation
  oder Repräsentation pflegen.
- **Lose Kopplung und Interoperabilität:** Die Feldnamen werden Teil der Aufrufsyntax; der Client
  koppelt sich an die interne Struktur.
- **Developer Experience:** Ein `fields=`-Parameter ist trivial zu benutzen, aber schwer zu
  dokumentieren und zu versionieren.
- **Sicherheit und Datenschutz:** Datensparsamkeit als Feature — heikle Felder werden gar nicht
  erst übertragen; umgekehrt braucht die Selektion eine Autorisierungsprüfung pro Feld.
- **Test- und Wartungsaufwand:** Die Zahl möglicher Antwortformen wächst kombinatorisch.

## Lösung

Der Client übergibt im Request eine *Wish List*: eine flache Aufzählung der gewünschten
Datenelemente, typischerweise als [Atomic Parameter List](../structure/AtomicParameterList.md) im
Query String. Der Provider liefert ausschließlich die genannten Elemente. Fehlt die Liste, gilt
üblicherweise „alles"; Wildcards sind eine gängige Erweiterung.

## Beispiel

Ohne Wish List liefert `GET /customers/gktlipwhjr` sämtliche Attribute inklusive
`moveHistory` und `customerInteractionLog`. Mit Wish List:

```
curl 'http://localhost:8080/customers/gktlipwhjr?fields=customerId,birthday,postalCode'
```

```json
{
  "customerId": "gktlipwhjr",
  "birthday": "1989-12-31T23:00:00.000+0000",
  "postalCode": "8640"
}
```

## Konsequenzen

**Vorteile:**

- Deutlich kleinere Antwortnachrichten ohne Vervielfachung der Endpoints.
- Ein Endpoint kann breite und schmale Clients gleichzeitig bedienen.
- Wirkt positiv auf ein [Rate Limit](RateLimit.md), das Datenvolumen einbezieht.

**Nachteile / Kosten:**

- Provider-seitig ist Response Shaping Aufwand: Projektion, Validierung unbekannter Feldnamen,
  Fehlerbehandlung.
- Caching wird schwerer — jede Feldkombination ist eine eigene Repräsentation.
- Statisch typisierte Clients bekommen Objekte mit Löchern; Pflichtfelder werden optional.
- Flach: Verschachtelte Strukturen lassen sich nur mit einer Pfadsyntax ausdrücken — dann ist man
  faktisch bei [Wish Template](WishTemplate.md).

## Bekannte Verwendungen

Sparse Fieldsets in JSON:API; `fields`-Parameter der Google Calendar API (mit Wildcard-Syntax);
Field Expansion in der Facebook Graph API; `expand` in Atlassian JIRA und Confluence; `$expand` in
der Microsoft Graph API; „Response Decoration" in den LinkedIn-APIs; Fields Masks in Flask-RESTPlus;
Attributselektion in den TMForum-REST-APIs. Der Netflix Technology Blog beschreibt Protobuf
`FieldMask` als gRPC-Umsetzung und schlägt vor, häufige Feldkombinationen als vorgefertigte Masken
in Client-Libraries auszuliefern.

## Verwandte Patterns

- [Wish Template](WishTemplate.md) — löst dasselbe Problem mit einer Struktur statt einer flachen Liste.
- [Pagination](Pagination.md) — verkleinert Antworten über die Länge statt über die Breite.
- [Conditional Request](ConditionalRequest.md) — kombinierbar: Feldauswahl greift nur, wenn überhaupt gesendet wird.
- [Request Bundle](RequestBundle.md) — Alternative bzw. Ergänzung, wenn viele Entitäten auf einmal gebraucht werden.
- [Atomic Parameter List](../structure/AtomicParameterList.md) — übliche Repräsentation der Wunschliste.
- [Parameter Tree](../structure/ParameterTree.md) / [Parameter Forest](../structure/ParameterForest.md) — die Strukturen, auf die sich die Wünsche beziehen.
- [Rate Limit](RateLimit.md) — profitiert von reduziertem Volumen.
- [Embedded Entity](EmbeddedEntity.md) / [Linked Information Holder](LinkedInformationHolder.md) — statische Antwort auf dieselbe Frage, zur Designzeit statt zur Laufzeit.

## Bezug zu Kubernetes / KRM

Hier hat KRM eine **echte Lücke**: Es gibt **kein generisches Feldselektionsprotokoll**. Weder
`GET /api/v1/namespaces/default/pods/x` noch ein `list` akzeptiert einen `fields`-Parameter im
Sinne von Sparse Fieldsets; die Antwort ist immer das vollständige Objekt gemäß Schema der
angeforderten Group/Version. Der Grund ist strukturell: KRM setzt auf ein *uniformes* Schema über
alle Ressourcen, auf dem generische Maschinerie (Informer-Cache, Server-Side Apply, Konvertierung
zwischen Versionen, Deep-Copy, Protobuf-Codecs) aufsetzt. Ein pro Request variables Teilobjekt
wäre in diesem Modell kein gültiges Objekt mehr — es ließe sich nicht von einem Objekt
unterscheiden, bei dem die Felder tatsächlich leer sind, und wäre damit für Level-Triggered
Reconciliation unbrauchbar.

```http
# So nicht — es gibt keinen Feldselektor für Antwortinhalte:
GET /api/v1/namespaces/default/pods?fields=metadata.name,status.phase
# Der Parameter wird ignoriert; die Antwort ist das vollständige PodList.
```

Statt einer freien Feldauswahl bietet KRM einige **vordefinierte Projektionen**, die per
Content-Negotiation über Accept-Parameter (`as`, `g`, `v`, ausgewertet in
`apiserver/pkg/endpoints/handlers/negotiation`) ausgehandelt werden:

- **Metadata-only:** liefert `TypeMeta` + `ObjectMeta` ohne `spec` und `status`. Das ist genau
  eine — die wichtigste — Wish List, fest verdrahtet. `client-go/metadata` setzt diesen Header, und
  der Garbage Collector betreibt darüber Metadata-only-Informer über sämtliche Ressourcentypen,
  ohne je ein vollständiges Objekt zu laden.

  ```http
  GET /api/v1/namespaces/default/pods
  Accept: application/json;as=PartialObjectMetadataList;g=meta.k8s.io;v=v1

  {
    "kind": "PartialObjectMetadataList",
    "apiVersion": "meta.k8s.io/v1",
    "metadata": { "resourceVersion": "8675309" },
    "items": [{
      "kind": "PartialObjectMetadata",          // <- jedes Item bekommt eigene TypeMeta
      "apiVersion": "meta.k8s.io/v1",
      "metadata": {
        "name": "web-7d9f8c-nkm2p",
        "namespace": "default",
        "uid": "3b1c...",
        "resourceVersion": "8675001",
        "ownerReferences": [ /* ... */ ]        // <- alles, was der GC braucht
      }
    }]
    // kein "spec", kein "status"
  }
  ```

- **Table:** liefert die Spaltendarstellung, die `kubectl get` anzeigt. `includeObject`
  (`None` | `Metadata` | `Object`, aus `metav1.TableOptions`) steuert dabei, wie viel des
  Originalobjekts pro Zeile mitkommt — der einzige Stellhebel in KRM, der einer Wish List in Form
  einer expliziten Client-Angabe nahekommt. CRDs definieren ihre Spalten über
  `spec.versions[].additionalPrinterColumns` mit JSONPath, also **provider-seitig**, nicht vom
  Client gewünscht.

  ```http
  GET /api/v1/namespaces/default/pods?includeObject=Metadata
  Accept: application/json;as=Table;g=meta.k8s.io;v=v1

  {
    "kind": "Table",
    "apiVersion": "meta.k8s.io/v1",
    "columnDefinitions": [
      { "name": "Name",     "type": "string", "format": "name", "description": "..." },
      { "name": "Ready",    "type": "string", "description": "..." },
      { "name": "Status",   "type": "string", "description": "..." },
      { "name": "Restarts", "type": "string", "description": "..." },
      { "name": "Age",      "type": "string", "description": "..." },
      { "name": "IP",       "type": "string", "priority": 1, "description": "..." }
      // priority > 0: nur bei kubectl get -o wide angezeigt
    ],
    "rows": [{
      "cells": ["web-7d9f8c-nkm2p", "1/1", "Running", "0", "4d", "10.244.1.7"],
      "object": { /* wegen includeObject=Metadata nur PartialObjectMetadata */ }
    }]
  }
  ```

- **Zeilenselektion statt Feldselektion:** `labelSelector` und `fieldSelector` in `ListOptions`
  reduzieren, *welche* Objekte kommen, nicht *welche Teile*. `fieldSelector` ist zudem nicht
  generisch: Jede Ressource registriert eine feste Menge selektierbarer Felder (für Pods etwa
  `spec.nodeName`, `spec.schedulerName`, `status.phase`, `status.podIP` — siehe
  `ToSelectableFields` in `pkg/registry/core/pod/strategy.go`); CRDs können seit
  `spec.versions[].selectableFields` eigene deklarieren.

Wichtig für die ehrliche Bewertung: `kubectl get -o jsonpath=...` und `-o custom-columns=...` sind
**client-seitig** und sparen kein einziges Byte auf der Leitung. Die *Data Transfer Parsimony*
entsteht nur bei `as=Table` und `as=PartialObjectMetadata`. Verglichen mit MAP fehlt KRM damit die
Laufzeit-Flexibilität des Patterns vollständig; es hat sie gegen eine kleine Zahl generischer,
cachebarer und über alle Ressourcen einheitlicher Projektionen eingetauscht.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/WishList)
