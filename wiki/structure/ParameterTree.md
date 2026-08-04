---
title: Parameter Tree
kategorie: Structure
unterkategorie: Representation Elements
quelle: https://microservice-api-patterns.org/patterns/structure/representationElements/ParameterTree
---

# Parameter Tree

*a.k.a.* Single Complex Representation, Tree Representation, Bar

**Kurzform:** Eine hierarchische Struktur mit genau einem Wurzelknoten bildet Containment-Beziehungen
zwischen Repräsentationselementen ab — Kindknoten sind Skalare, Listen oder wiederum Teilbäume.

## Kontext

Ein API-Provider bietet Operations an einem Endpoint an. Die auszutauschenden Daten sind nicht flach:
Elemente gehören zu anderen Elementen, es gibt Wiederholungen und optionale Teile.

## Problem

Wie lassen sich Containment-Beziehungen ausdrücken, wenn komplexe Repräsentationselemente definiert
und zur Laufzeit ausgetauscht werden?

## Forces

Wie bei den drei übrigen Repräsentationspatterns:

- Struktur des Domänenmodells und des Systemverhaltens und deren Wirkung auf Verständlichkeit
- Zusätzlich zu übertragende Daten (Security-Informationen, Metadaten)
- Performance (Latenz, Verarbeitung) und Ressourcenverbrauch (Bandbreite, Speicher, CPU)
- Lose Kopplung und Interoperabilität
- Developer Convenience und Developer Experience
- Sicherheit und Datenschutz

## Lösung

Eine Hierarchie mit dediziertem Wurzelknoten und einem oder mehreren Kindknoten definieren. Jeder
Kindknoten ist entweder ein [Atomic Parameter](AtomicParameter.md), eine
[Atomic Parameter List](AtomicParameterList.md) oder ein weiterer *Parameter Tree*; identifiziert wird
lokal per Name und/oder Position. Für jeden Knoten die Kardinalität festlegen: exakt eins, null oder
eins, mindestens eins, null oder mehr.

## Varianten

**Parameter Collection.** Nesting-Tiefe 1 und alle Kinder der Wurzel haben dieselbe Struktur — der
Baum degeneriert zur Liste. Dies modelliert JSON-Arrays und XML-Elemente mit `maxOccurs > 1`.

**Atomic Parameter Collection.** Sind alle Knoten auf Ebene 1 Blätter, ist die Struktur faktisch eine
in den Baum eingebettete [Atomic Parameter List](AtomicParameterList.md).

## Beispiel

Eine Response, deren Wurzelobjekt Skalare und zwei geschachtelte Teilbäume (`evidence`, `links`)
enthält:

```json
{"claim":
    {"id":"0afeb849-6d63-40b6-b52f-21dee16fdda5",
     "dateOfIncident":"2017-02-14",
     "amount":2000.0,
     "evidence":[],
     "links":[{"uri":"http://localhost:8080/claims/0afeb849-...",
               "rel":"self"}]}}
```

## Implementierungshinweise

Die Quelle nennt unter anderem: Grenzen, Optionalität, `Null`-Behandlung und Validierung müssen für
*jeden* Skalar und *jede* Listenstruktur im Baum einzeln festgelegt werden — deutlich aufwendiger als
bei den flachen Patterns. Nicht die reale Welt vollständig modellieren („if in doubt, leave it out").
Maschinenlesbare Schemata (JSON Schema, XML Schema) plus Beispieldaten bereitstellen. Zu tiefe
Verschachtelung vermeiden, außer die Zahl der Calls ist der dominante Faktor. Vollständig generische
Strukturen (dynamisch interpretierte Key-Value-Paare) zurückhaltend einsetzen: die versprochene
Flexibilität kostet Verständlichkeit, Code-Completion und Testautomatisierung. Ein *Parameter Tree* im
Request-Body eines HTTP GET ist laut HTTP/1.1-Spezifikation zu vermeiden.

## Konsequenzen

**Vorteile:**

- Bildet Domänenstrukturen direkt ab; ein Call statt vieler.
- Ein einziger Wurzelknoten macht die Nachricht als Ganzes benennbar, schemafähig und versionierbar.
- Additive Erweiterung um optionale Felder ist rückwärtskompatibel möglich.

**Nachteile / Kosten:**

- Höherer Validierungs- und Spezifikationsaufwand pro Knoten.
- Nachrichtengröße und Parsing-Kosten steigen; tiefe Verschachtelung erschwert Verarbeitung.
- Engere Kopplung an das Domänenmodell des Providers — Änderungen strahlen weiter aus.
- Generische Sub-Bäume (frei belegbare Maps) verlagern Fehler von der Compile- in die Laufzeit.

## Bekannte Verwendungen

- JAX-RS mit `@Consumes`/`@Produces` auf Custom Media Types (geschachteltes JSON mit einer Wurzel).
- JIRA Cloud REST API: Request von `issue-createIssue`.
- Twitter REST API: `GET collections` mit einem geschachtelten Objekt `objects`.
- Protocol Buffers allgemein; komplexe Typen in Apache Avro.

## Verwandte Patterns

- [Atomic Parameter](AtomicParameter.md) — Baustein für die Blätter.
- [Atomic Parameter List](AtomicParameterList.md) — einfachere Alternative; nutzt umgekehrt den
  *Parameter Tree* als Transport-Wrapper, wenn die Technologie nur einen Parameter zulässt.
- [Parameter Forest](ParameterForest.md) — Alternative, wenn eine gemeinsame Wurzel künstlich wirkt.
- [Pagination](../quality/Pagination.md) — die Seiten der Response sind üblicherweise Trees oder Forests.
- [Embedded Entity](../quality/EmbeddedEntity.md) — eingebettete Teilbäume statt Links.
- [Linked Information Holder](../quality/LinkedInformationHolder.md) — Gegenentwurf: Verweis statt Baum.
- [Error Report](ErrorReport.md), [Context Representation](ContextRepresentation.md) — typischerweise
  als Tree strukturiert.
- [Wish Template](../quality/WishTemplate.md) — beschreibt einen gewünschten Teilbaum.

Verwendbar in *Command Message*, *Document Message* und *Event Message* (Hohpe/Woolf 2003);
*Document Message* wird meist als *Parameter Tree* transportiert. *Content Filter* und
*Content Enricher* (ebd.) operieren typischerweise auf solchen Strukturen.

## Bezug zu Kubernetes / KRM

Das KRM ist im Kern eine konsequente Anwendung dieses Patterns: Jede Ressource ist ein
*Parameter Tree* mit fester Wurzelstruktur aus `apiVersion`, `kind`, `metadata`, `spec` und `status`.
`metadata` (`metav1.ObjectMeta`) ist ein überall identischer Teilbaum; `spec` und `status` sind
ressourcenspezifisch.

```yaml
apiVersion: apps/v1          # <- Wurzel-Skalare, bei jeder Ressource identisch benannt
kind: Deployment
metadata:                    # <- metav1.ObjectMeta: derselbe Teilbaum für jeden Typ
  name: web
  namespace: default
  labels:
    app: web                 # <- generischer Key-Value-Teilbaum, bewusst untypisiert
  # ... uid, resourceVersion, generation, ownerReferences, finalizers, managedFields
spec:                        # <- typspezifischer Teilbaum: gewünschter Zustand
  replicas: 3
  selector:
    matchLabels:
      app: web
  template:                  # <- Teilbaum im Teilbaum (PodTemplateSpec)
    metadata:
      labels:
        app: web
    spec:
      containers:
        - name: app
          image: nginx:1.27
          # ... ports, env, resources
status:                      # <- typspezifischer Teilbaum: beobachteter Zustand
  observedGeneration: 7
  readyReplicas: 3
  # ... replicas, updatedReplicas, availableReplicas, conditions
```

Diese Uniformität über alle Ressourcen hinweg ist der eigentliche Unterschied zu MAP: MAP wählt die Repräsentationsstruktur pro
Operation, KRM legt sie für die gesamte API-Fläche einmal fest, weil generische Clients (`kubectl`,
Controller, Garbage Collector, Admission-Webhooks) über beliebige Typen hinweg funktionieren müssen.
Die von der Quelle abgeratene „generische Key-Value-Struktur" existiert im KRM bewusst und begrenzt —
`labels` und `annotations` sind `map[string]string`, und `runtime.RawExtension` erlaubt untypisierte
Teilbäume (etwa `object` in `metav1.WatchEvent`).

Die Kardinalitäten aus der Lösung werden in Go über Pointer-Typen und `omitempty` ausgedrückt:
`*int64` plus `omitempty` bedeutet „null oder eins", ein Slice ohne `omitempty` erzwingt ein
vorhandenes (ggf. leeres) Array. Die API-Konventionen verlangen explizit, `null` und leere Objekte
nicht bedeutungsunterscheidend zu verwenden, weil sonst Round-Trips über JSON, Protobuf und
Strategic-Merge-Patch nicht mehr verlustfrei sind.

Über MAP hinaus geht KRM bei der Merge-Semantik des Baums. Für Server-Side Apply muss pro Knoten
erklärt sein, ob er als Einheit oder granular behandelt wird: `+listType=atomic|set|map` mit
`+listMapKey=...` an Slices und `+mapType=atomic|granular` an Maps (in
`staging/src/k8s.io/api/core/v1/types.go` hundertfach vorhanden; in CRD-Schemata als
`x-kubernetes-list-type`, `x-kubernetes-list-map-keys` und `x-kubernetes-map-type`, siehe
`staging/src/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/types_jsonschema.go`). Damit ist
`spec.containers` als `listType=map` mit Key `name` deklariert und mehrere Feldmanager können
unabhängig einzelne Container verwalten.

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
# ... metadata.name, spec.group, spec.names, spec.scope
spec:
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                backends:
                  type: array
                  x-kubernetes-list-type: map          # <- Knoten wird granular gemerged
                  x-kubernetes-list-map-keys: ["name"] # <- Identität des Eintrags im Baum
                  items:
                    type: object
                    x-kubernetes-map-type: granular    # <- Default für Objekte
                    properties:
                      name: { type: string }
                      weight: { type: integer }
      subresources:
        status: {}
```

`x-kubernetes-list-type` defaultet für Arrays auf `atomic`,
`x-kubernetes-map-type` für Objekte auf `granular`. MAP beschreibt den Baum nur als Datenstruktur; das
KRM beschreibt zusätzlich, wer welchen Teilbaum besitzen darf — eine Notwendigkeit, die erst aus
deklarativem Modell und Level-Triggered Reconciliation mit mehreren gleichzeitigen Schreibern
entsteht.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/representationElements/ParameterTree)
