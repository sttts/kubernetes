---
title: Atomic Parameter List
kategorie: Structure
unterkategorie: Representation Elements
quelle: https://microservice-api-patterns.org/patterns/structure/representationElements/AtomicParameterList
---

# Atomic Parameter List

*a.k.a.* Multiple Scalar Representations, Dotted Line

**Kurzform:** Mehrere zusammengehörige [Atomic Parameters](AtomicParameter.md) werden zu einem
kohäsiven Repräsentationselement gruppiert — jedes einzelne bleibt simpel, aber die Zusammengehörigkeit
wird im Kontrakt und zur Laufzeit sichtbar.

## Kontext

Ein API-Provider bietet Operations an einem Endpoint an; Client und Provider müssen sich auf die
Struktur der Request- und Response-Nachrichten einigen. Ein einzelner Skalar reicht nicht, ein voller
Baum ist überdimensioniert.

## Problem

Wie lassen sich mehrere verwandte Atomic Parameters so kombinieren, dass jeder für sich einfach bleibt,
ihre Zusammengehörigkeit aber in der API Description und in den Nachrichten explizit wird?

## Forces

Identisch zu den anderen drei Repräsentationspatterns:

- Struktur des Domänenmodells und des Systemverhaltens und deren Wirkung auf Verständlichkeit
- Zusätzlich zu übertragende Daten (Security-Informationen, Metadaten)
- Performance (Latenz, Verarbeitung) und Ressourcenverbrauch (Bandbreite, Speicher, CPU)
- Lose Kopplung und Interoperabilität
- Developer Convenience und Developer Experience
- Sicherheit und Datenschutz

## Lösung

Zwei oder mehr einfache Datenelemente in einem kohäsiven Repräsentationselement gruppieren. Die
Einträge werden entweder über die Position (Index) oder über einen String-Key identifiziert. Die Liste
als Ganzes bekommt einen eigenen Namen, sofern der Empfänger sie als Einheit verarbeiten muss.
Explizit spezifizieren, wie viele Elemente erforderlich und wie viele zulässig sind.

## Varianten

**Atomic Parameter Collection.** Haben alle Einträge dieselbe Struktur und werden sie über die Position
statt über einen Key identifiziert, wird die *Atomic Parameter List* zu einer homogenen
*Atomic Parameter Collection*.

## Beispiel

Skalare in Query-String und JSON-Body:

```
curl http://localhost:8080/riskreport?firstYear=2017&noOfYears=1

curl http://localhost:8080/claims -H "Content-Type: application/json" \
     -d '{"dateOfIncident":"2017-02-01", "amount": 2000 }'
```

Response mit fünf skalaren Elementen:

```json
[{"year":2017,"policyCount":42,
  "totalClaimValue":2000.0,"claimCount":1,
  "totalInsuredValue":1000000.0}]
```

## Konsequenzen

**Vorteile:**

- Zusammengehörigkeit wird explizit, ohne dass ein Objektschema mit Wurzelknoten nötig wird.
- Weiterhin flach und daher billig zu serialisieren, zu validieren und in URLs abzubilden.
- Gute Zwischenstufe: der Weg von [Atomic Parameter](AtomicParameter.md) hierher ist klein.

**Nachteile / Kosten:**

- Positionsbasierte Identifikation ist brüchig — Umsortieren oder Einfügen bricht Clients.
- Skaliert schlecht: ab einer gewissen Länge wird die Liste unlesbar; dann ist
  [Parameter Tree](ParameterTree.md) oder [Parameter Forest](ParameterForest.md) fällig.
- Verschachtelte Domänenstrukturen lassen sich nicht abbilden, ohne Namen künstlich zu kodieren
  (`address_street`, `address_city` etc.).
- Manche Technologien können mehrere Top-Level-Parameter nicht transportieren; dann ist ein
  [Parameter Tree](ParameterTree.md) als Wrapper nötig.

## Bekannte Verwendungen

- Mehrere URI-Query-Parameter an einem HTTP-GET; ebenso *URI Templates* nach RFC 6570.
- Mehrere einfache Parameter in einem HTTP-POST.
- Facebook Graph API: Response von GET events-per-user.
- Twitter API, z.B. `POST /1.1/statuses/update.json?status=...&lat=...&lon=...`.
- Protocol-Buffers-Nachrichten, die ausschließlich einfache Feldtypen enthalten.
- Swagger/OpenAPI: *parameters definitions* zur Wiederverwendung über Operations hinweg.

## Verwandte Patterns

- [Atomic Parameter](AtomicParameter.md) — sowohl einfachere Alternative als auch Baustein.
- [Parameter Tree](ParameterTree.md) — nächster Evolutionsschritt bei wachsender Liste; dient auch als
  Wrapper, wenn die Technologie keine Mehrfachparameter kennt.
- [Parameter Forest](ParameterForest.md) — enthält Listen aus Elementen aller drei anderen Patterns.
- [Pagination](../quality/Pagination.md) — die Query-Parameter der Seitenanfrage bilden typischerweise
  eine *Atomic Parameter List*.
- [Wish List](../quality/WishList.md) — Liste gewünschter Feldnamen als Skalare.
- [Metadata Element](MetadataElement.md), [Data Element](DataElement.md), [Id Element](IdElement.md) —
  Stereotypen, die als Listeneinträge auftreten.

Wie [Atomic Parameter](AtomicParameter.md) verwendbar in *Command Message*, *Document Message* und
*Event Message* (Hohpe/Woolf 2003), sofern der Inhalt als Liste darstellbar ist.

## Bezug zu Kubernetes / KRM

Die deutlichste Instanz im KRM ist der Query-String der List-/Watch-Operationen:
`?labelSelector=app%3Dnginx&limit=500&continue=...&resourceVersion=0&watch=true` ist eine Sammlung
verwandter Skalare, die gemeinsam eine Abfrage beschreiben. Kubernetes lässt sie aber nicht als lose
Liste stehen, sondern bindet sie an den versionierten Typ `metav1.ListOptions`
(`staging/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go`) — die Liste ist damit ein benanntes,
schemabehaftetes Objekt und nicht nur eine Konvention über Positionen oder Keys. Gleiches gilt für
`metav1.DeleteOptions` (`gracePeriodSeconds`, `propagationPolicy`, `dryRun`) und
`metav1.PatchOptions` (`dryRun`, `force`, `fieldManager`, `fieldValidation`). Das ist Variante (b):
das Pattern wird erfüllt, aber über einen einheitlichen Objekttyp statt über Parametersignaturen.

```http
GET /api/v1/namespaces/default/pods?labelSelector=app%3Dnginx&limit=500&resourceVersion=0&watch=true
```

```yaml
# Derselbe Aufruf als typisiertes Objekt — ListOptions trägt TypeMeta,
# die Parameterliste ist also selbst benannt und versioniert:
apiVersion: v1
kind: ListOptions          # <- die "Liste" hat einen Namen im Schema
labelSelector: app=nginx
limit: 500
resourceVersion: "0"
watch: true
# ... fieldSelector, continue, resourceVersionMatch, allowWatchBookmarks, timeoutSeconds
```

Innerhalb von Payloads verbieten die Kubernetes-API-Konventionen genau die naive Form dieses Patterns:
ein Top-Level-JSON-Array als Request- oder Response-Body ist nicht zulässig, weil es weder `apiVersion`
und `kind` noch spätere additive Erweiterung erlaubt. Sammlungen erscheinen deshalb immer als Feld
innerhalb eines Objekts — `items` in `metav1.List`, `rows` und `columnDefinitions` in `metav1.Table`.
Das entspricht exakt dem in der Quelle beschriebenen Fall, dass eine *Atomic Parameter List* einen
[Parameter Tree](ParameterTree.md) als Transport-Wrapper braucht.

```yaml
# So nicht — ein Top-Level-Array als Request- oder Response-Body ist im KRM unzulässig:
- metadata: { name: web-1 }
- metadata: { name: web-2 }
```

```yaml
# Sondern: die Sammlung ist ein Feld unter einer Wurzel mit apiVersion/kind.
apiVersion: v1
kind: PodList
metadata:                     # <- ListMeta, nicht ObjectMeta
  resourceVersion: "12345"
  continue: eyJ2IjoibWV0YS5rOHMuaW8vdjEi...
items:                        # <- der Listen-Wrapper, den die Quelle beschreibt
  - metadata:
      name: web-1
    # ... spec, status
  - metadata:
      name: web-2
    # ... spec, status
```

Für flache Skalarlisten *innerhalb* des Payloads geht KRM über MAP hinaus: Felder werden mit
`+listType=atomic` oder `+listType=set` annotiert (in OpenAPI/CRDs `x-kubernetes-list-type`).
`atomic` heißt, die gesamte Liste wird beim Update als eine Einheit ersetzt — der Skalar-Charakter
der Liste wird also maschinenlesbar erklärt. `set` verlangt Eindeutigkeit der Werte und macht die
Reihenfolge irrelevant.

```yaml
apiVersion: v1
kind: Pod
metadata:
  finalizers:                 # <- +listType=set: eindeutig, Reihenfolge bedeutungslos,
    - foregroundDeletion      #    mehrere Akteure dürfen unabhängig Einträge halten
    - example.com/drain-connections
spec:
  containers:
    - name: app
      args:                   # <- +listType=atomic: nur als Ganzes ersetzbar,
        - --port=8080         #    Reihenfolge ist bedeutungstragend
        - --v=2
```

Ohne diese Annotation ist für Server-Side Apply nicht entscheidbar, ob zwei Feldmanager dieselbe
Liste unabhängig befüllen dürfen; MAP kennt diese Merge-Dimension nicht.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/representationElements/AtomicParameterList)
