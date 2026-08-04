---
title: Two in Production
kategorie: Evolution
unterkategorie: Lifecycle Management
quelle: https://microservice-api-patterns.org/patterns/evolution/TwoInProduction
---

# Two in Production

*a.k.a.* *Parallel Versions*, *Rolling Update Policy*

**Kurzform:** Der Provider betreibt dauerhaft zwei (oder allgemein *n*) Versionen desselben
Endpunkts parallel und rollt sie überlappend weiter, sodass Clients ihre Migration in ihrem
eigenen Takt vornehmen können.

## Kontext

Eine API entwickelt sich regelmäßig weiter, irgendwann inkompatibel. Clients — besonders von
einer [Public API](../foundation/PublicAPI.md) oder
[Community API](../foundation/CommunityAPI.md) — laufen in unterschiedlichen
Geschwindigkeiten; ihre Release-Zyklen sind nicht mit denen des Providers synchronisiert und
lassen sich auch nicht synchronisieren.

## Problem

Wie kann ein Provider eine API schrittweise weiterentwickeln, ohne bestehende Clients zu
brechen, aber auch ohne beliebig viele Versionen in Produktion pflegen zu müssen?

## Forces

- **Unterschiedliche Lebenszyklen** von Provider und Client müssen möglich bleiben.
- **Keine unentdeckten Kompatibilitätsprobleme** — Brüche müssen sichtbar sein.
- **Rollback-Fähigkeit**, falls sich eine neue Version als Fehlentwurf erweist.
- **Minimale Client-Änderungen** durch API-Änderungen.
- **Minimaler Wartungsaufwand** für die Unterstützung alter Versionen.

Die letzten beiden Kräfte stehen direkt gegeneinander; die Zahl „zwei“ ist der Kompromiss.

## Lösung

Zwei Versionen eines API-Endpunkts und seiner Operationen deployen und supporten, die
Varianten derselben Funktionalität anbieten und untereinander *nicht* kompatibel sein müssen.
Versionen werden rollierend und überlappend eingeführt und außer Betrieb genommen.

## Varianten

**N in Production.** Statt zwei auch drei oder mehr Versionen. „Managed Evolution“
(Murer/Bonati/Furrer) berichtet drei Versionen als guten Kompromiss zwischen Provider-
Komplexität und Migrationsgeschwindigkeit. Ein
[Experimental Preview](ExperimentalPreview.md) kann eine der Versionen sein — der Provider
sollte dann klarstellen, ob er auf das Kontingent angerechnet wird.

## Beispiel

Ein ERP-Hersteller veröffentlicht API-Version 1. Neue Pensionsplan-Funktionen im HR-Modul
brechen den Payroll-Teil, also liefert das nächste ERP-Release Version 1 *und* Version 2 aus:
Bestandskunden bleiben zunächst auf 1, Neukunden nutzen sofort 2. Mit dem übernächsten Release
kommt Version 3 und Version 1 fällt weg — supported sind nun 2 und 3. Wer noch auf 1 ist, wird
abgeschnitten (und kann umgeleitet werden).

## Konsequenzen

**Vorteile:**

- Entkoppelt Release-Zyklen von Provider und Client vollständig.
- Gibt Clients ein planbares, gleitendes Migrationsfenster.
- Ermöglicht Rollback und A/B-Betrieb einer neuen Vertragsvariante.

**Nachteile / Kosten:**

- Doppelter Betrieb, doppelte Tests, doppelte Fehlerbehebung — auch für Sicherheitsfixes.
- Zustandsbehaftete Backends sind das eigentliche Problem: zwei Verträge auf *einem*
  Datenbestand erzwingen Mapping oder Datenduplikation.
- Migrationsdruck bleibt bestehen; das Fenster ist nur verschoben, nicht aufgehoben.
- Ohne strikte Abschaltdisziplin wird aus „zwei in Produktion“ schnell „sieben in Produktion“.

## Bekannte Verwendungen

Terravis betreibt zwei Major-Versionen zwei Jahre lang parallel; eine große europäische Bank
hält ebenfalls zwei Major-Versionen gleichzeitig. GitHub bietet die v3-API zusammen mit der
Folgeversion an. Facebook beschreibt unter „Platform Versioning“ eine rollierende Release- und
Update-Politik, was die Quellseite als definierendes Merkmal des Musters wertet.

## Verwandte Patterns

- [Version Identifier](VersionIdentifier.md) — notwendige Voraussetzung, um die parallel
  aktiven Versionen überhaupt zu unterscheiden.
- [Semantic Versioning](SemanticVersioning.md) — vollständig kompatible Stände (Patch) dürfen
  eine aktive Version ersetzen, ohne das Kontingent zu belasten.
- [Aggressive Obsolescence](AggressiveObsolescence.md) — das Mittel, um Clients von der alten
  Version herunterzubekommen und wieder Platz zu schaffen.
- [Limited Lifetime Guarantee](LimitedLifetimeGuarantee.md) — die stärkere Zusage, wenn Clients
  ein festes Enddatum statt eines gleitenden Fensters brauchen.
- [Eternal Lifetime Guarantee](EternalLifetimeGuarantee.md) — das Gegenmodell ohne
  Abschaltung.
- [Experimental Preview](ExperimentalPreview.md) — kann eine der parallelen Versionen sein.
- [API Description](../foundation/APIDescription.md),
  [Service Level Agreement](../quality/ServiceLevelAgreement.md) — dokumentieren, welche
  Versionen gerade aktiv sind.

## Bezug zu Kubernetes / KRM

Das ist das Muster, bei dem KRM am weitesten über MAP hinausgeht — und zwar in drei Punkten.

**Erstens: es sind nicht zwei Endpunkte, sondern *n* Sichten auf dasselbe Objekt.** Der
API-Server serviert für eine Ressource beliebig viele Versionen gleichzeitig. `resource.k8s.io`
registriert in `pkg/registry/resource/rest/storage_resource.go` derzeit vier Versionen
nebeneinander — `v1`, `v1beta2`, `v1beta1` und `v1alpha3`; `resourceclaims`, `resourceslices`
und `resourceclaimtemplates` werden dabei in dreien davon gleichzeitig serviert. Ein Objekt, das über
`/apis/<group>/v1beta1/...` angelegt wurde, ist über
`/apis/<group>/v1/...` **dasselbe Objekt** — gleiche `metadata.uid`, gleiche
`metadata.resourceVersion`, gleiche Position im `watch`-Stream, gleiche Wirkung auf Controller.
Bei MAP koexistieren zwei Verträge, die dieselbe Funktionalität in Varianten anbieten; bei KRM
koexistieren zwei Repräsentationen ein und desselben Zustands. Ein `kubectl get` in der einen
und ein `kubectl edit` in der anderen Version arbeiten auf demselben Datensatz.

Derselbe Datensatz, zweimal abgerufen — einmal über `v1beta2`, einmal über `v1beta3`. Beide
wurden in 1.26 bis 1.28 parallel serviert, und zwischen ihnen wurde ein Feld umbenannt:

```yaml
# GET /apis/flowcontrol.apiserver.k8s.io/v1beta2/prioritylevelconfigurations/workload-low
apiVersion: flowcontrol.apiserver.k8s.io/v1beta2
kind: PriorityLevelConfiguration
metadata:
  name: workload-low
  uid: 6f1c9d3a-4b2e-4f7a-9c11-2d8e5a0b7c34   # <- identisch
  resourceVersion: "1842"                     # <- identisch
  # ... generation, creationTimestamp, labels, annotations
spec:
  type: Limited
  limited:
    assuredConcurrencyShares: 100             # <- alter Feldname
    lendablePercent: 90
    # ... limitResponse, borrowingLimitPercent
```

```yaml
# GET /apis/flowcontrol.apiserver.k8s.io/v1beta3/prioritylevelconfigurations/workload-low
apiVersion: flowcontrol.apiserver.k8s.io/v1beta3
kind: PriorityLevelConfiguration
metadata:
  name: workload-low
  uid: 6f1c9d3a-4b2e-4f7a-9c11-2d8e5a0b7c34   # <- identisch
  resourceVersion: "1842"                     # <- identisch
  # ... generation, creationTimestamp, labels, annotations
spec:
  type: Limited
  limited:
    nominalConcurrencyShares: 100             # <- derselbe Wert, neuer Name
    lendablePercent: 90
    # ... limitResponse, borrowingLimitPercent
```

Es sind nicht zwei Objekte, die synchron gehalten werden; es ist ein Objekt und zwei Decoder.

**Zweitens: verlustfreie Round-Trip-Konversion als erzwungene Eigenschaft.** Für eingebaute
APIs existiert pro Gruppe eine interne Hub-Version, `runtime.APIVersionInternal` = `__internal`
(definiert in `k8s.io/apimachinery/pkg/runtime/interfaces.go`). Konvertiert wird immer
sternförmig: `v1beta1 → __internal → v1`. Das reduziert den Aufwand von O(n²) Konvertern auf
O(n) und ist der Grund, warum viele Versionen praktikabel bleiben. Die Konvertierungsfunktionen
werden generiert (`zz_generated.conversion.go`, `k8s.io/apimachinery/pkg/conversion`), und die
Verlustfreiheit wird durch fuzzing-basierte Round-Trip-Tests abgesichert
(`k8s.io/apimachinery/pkg/api/apitesting/roundtrip`). Wo eine Version eine Unterscheidung nicht
ausdrücken kann, wird sie in einer Annotation zwischengeparkt, damit der Round-Trip trotzdem
hält. `flowcontrol` hat dafür ein reales Beispiel: in `v1` ist
`nominalConcurrencyShares` ein Pointer, `0` und „nicht gesetzt“ sind unterscheidbar; in
`v1beta3` ist das Feld ein `int32`, in dem `0` historisch „nimm den Default 30“ bedeutete.

```yaml
apiVersion: flowcontrol.apiserver.k8s.io/v1beta3
kind: PriorityLevelConfiguration
metadata:
  name: workload-low
  annotations:
    flowcontrol.k8s.io/v1beta3-preserve-zero-concurrency-shares: ""   # <- „0 heißt wirklich 0“
    # ... apf.kubernetes.io/autoupdate-spec
spec:
  type: Limited
  limited:
    nominalConcurrencyShares: 0
```

Die Annotation ist die Stelle, an der die Round-Trip-Anforderung sichtbar wird: sie existiert
nur, damit eine in `v1` ausdrückbare Information beim Weg durch `v1beta3` nicht verlorengeht.
MAP hat für diese Anforderung kein Gegenstück — dort dürfen die Versionen ausdrücklich
inkompatibel sein.

**Drittens: Speicherversion und Migration sind explizit.** Persistiert wird nur *eine* Version
pro Ressource, die Storage Version. Bei CRDs ist das eine Deklaration im Objekt selbst:

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: widgets.example.com
spec:
  group: example.com
  names:
    kind: Widget
    plural: widgets
    # ... singular, shortNames, categories
  scope: Namespaced
  versions:
    - name: v1beta1
      served: true                 # <- wird noch ausgeliefert
      storage: false               # <- aber nicht mehr gespeichert
      # ... schema, subresources, additionalPrinterColumns
    - name: v1
      served: true
      storage: true                # <- genau ein Eintrag darf das sein
      # ... schema, subresources, additionalPrinterColumns
  conversion:
    strategy: Webhook              # <- Alternative: None (schreibt nur apiVersion um)
    webhook:
      conversionReviewVersions: ["v1"]
      clientConfig:
        service:
          namespace: example-system
          name: widget-conversion
          path: /convert
          port: 443
        # ... caBundle
```

`served: true` und `storage: true` sind orthogonal: `served` entscheidet, ob eine Version über
HTTP erreichbar ist, `storage`, in welcher Form etcd sie sieht. `strategy: None` schreibt beim
Lesen nur die `apiVersion` um und lässt alles andere unverändert — brauchbar nur, solange die
Versionen schemagleich sind; alles andere braucht `Webhook`.

Damit Clients einen Wechsel der Speicherversion erkennen können, liefert Discovery pro Ressource
`storageVersionHash` in `metav1.APIResource`:

```json
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "apps/v1",
  "resources": [
    {
      "name": "deployments",
      "singularName": "deployment",
      "namespaced": true,
      "kind": "Deployment",
      "verbs": ["create", "delete", "get", "list", "patch", "update", "watch"],
      "storageVersionHash": "8aSe+NMSvyU="
    }
  ]
}
```

Der Wert ist für Clients opak; nur der Gleichheitsvergleich ist definiert. Ändert er sich, hat
der Server die Speicherform gewechselt. Die interne API-Gruppe `internal.apiserver.k8s.io` führt
zusätzlich `StorageVersion`-Objekte, und `storagemigration.k8s.io/v1beta1` bietet
`StorageVersionMigration`, um bestehende Objekte kontrolliert umzuschreiben:

```yaml
apiVersion: storagemigration.k8s.io/v1beta1
kind: StorageVersionMigration
metadata:
  name: widgets-to-v1
spec:
  resource:
    group: example.com        # <- keine Version: migriert wird *in* die aktuelle Storage Version
    resource: widgets
status:
  resourceVersion: "1842"     # <- Stand, ab dem der Controller durchläuft
  conditions:
    - type: Succeeded         # <- daneben: Running, Failed
      status: "True"
      # ... lastTransitionTime, reason, message
```

Bemerkenswert ist, was in `spec.resource` *fehlt*: eine Version. Der Auftrag lautet nicht
„schreibe nach v1 um“, sondern „lies und schreibe jedes Objekt einmal“ — welche Version dabei
auf die Platte kommt, weiß nur der Server. Der praktische Effekt: der Provider kann seine
Speicherform migrieren, ohne dass irgendein Client davon etwas merkt — eine Entkopplung, die
MAP nicht adressiert.

Der Preis ist genau derselbe, den die Quellseite nennt: Konvertierbarkeit muss beim Design
*jeder* Feldänderung mitgedacht werden. Ein Feld, das sich nicht verlustfrei auf die alte
Version abbilden lässt, ist in Kubernetes de facto nicht baubar, solange die alte Version noch
serviert wird. Das ist der Grund, warum Kubernetes fast nie `major`-artig versioniert (siehe
[Semantic Versioning](SemanticVersioning.md)) — statt inkompatibler Parallelversionen wird
Kompatibilität erzwungen.

---
[← Index](../README.md) · [Kategorie Evolution](../meta/category-evolution.md) · [Quelle](https://microservice-api-patterns.org/patterns/evolution/TwoInProduction)
