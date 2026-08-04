---
title: Limited Lifetime Guarantee
kategorie: Evolution
unterkategorie: Lifecycle Management
quelle: https://microservice-api-patterns.org/patterns/evolution/LimitedLifetimeGuarantee
---

# Limited Lifetime Guarantee

*a.k.a.* *Fixed Lifetime Guarantee*, *Guaranteed Availability Period*

**Kurzform:** Der Provider sagt zu, eine veröffentlichte API-Version für einen festen Zeitraum
nicht zu brechen, und versieht jede Version mit einem Ablaufdatum.

## Kontext

Eine API ist veröffentlicht und hat mindestens einen Client. Der Provider kann die Roadmaps
seiner Clients nicht steuern, oder der Schaden durch erzwungene Client-Änderungen (finanziell,
reputationsseitig) wäre hoch. Er will deshalb keine Breaking Changes an der publizierten API
vornehmen — aber die API auch nicht für immer einfrieren.

## Problem

Wie teilt ein Provider seinen Clients mit, wie lange sie sich auf die veröffentlichte Version
einer API verlassen können?

## Forces

- **Planbarkeit** der durch API-Änderungen erzwungenen Client-Änderungen.
- **Begrenzung des Wartungsaufwands** für die Unterstützung alter Clients.

Das Muster löst die Spannung, indem es die Unsicherheit auf beiden Seiten in ein Datum
überführt: der Client weiß, bis wann er migrieren muss, der Provider weiß, ab wann er
aufräumen darf.

## Lösung

Als Provider garantieren, die veröffentlichte API für einen festen Zeitraum nicht zu brechen,
und jede API-Version mit einem Ablaufdatum kennzeichnen. Nach Ablauf darf der Provider
beliebige Änderungen vornehmen oder die Version abschalten — die Frist wirkt damit als
eingebauter, impliziter Deprecation-Mechanismus.

## Beispiel

Die Quellseite nennt die IBAN-Einführung in Europa: Eine EU-Parlamentsresolution von 2012 setzte
eine Frist bis 2014, nach der die alten nationalen Kontonummern durch den neuen Standard
ersetzt sein mussten. Systeme, die Konten identifizieren, mussten für ihre alten Operationen
eine *Limited Lifetime Guarantee* aussprechen. Das Beispiel zeigt zugleich, dass Evolutions-
und Versionierungsstrategien nicht immer vom Provider allein bestimmt werden, sondern von
Gesetzgebung oder Branchenkonsortien vorgegeben sein können.

## Konsequenzen

**Vorteile:**

- Migrationen werden budgetierbar und projektierbar — beim Client wie beim Provider.
- Der Provider bekommt ein legitimes Recht zum Abschalten, ohne einzeln verhandeln zu müssen.
- Die Zusage ist prüfbar und eignet sich als Vertragsbestandteil.

**Nachteile / Kosten:**

- Der Provider bindet sich für die gesamte Frist — auch wenn sich die Version als Fehlentwurf
  erweist.
- Innerhalb der Frist bleibt technische Schuld liegen; Sicherheits- und Infrastrukturupgrades
  müssen die alte Version mittragen.
- Der Cutoff ist hart: Clients, die die Frist verpassen, fallen aus.
- Fristen müssen pro Version verwaltet und kommuniziert werden — Governance-Aufwand.

## Bekannte Verwendungen

Facebook gibt für Kern-API und SDK eine Zwei-Jahres-Garantie, auch gegenüber anonymen Clients
(„Platform Versioning“). Google Adwords nannte 10 Monate. eBay arbeitet mit expliziten
Deprecation Policies. Twitter stellte die Unterstützung für HTTP Basic Authentication zum
31.08.2010 ein — angekündigt am 28.04.2010, samt eigener Countdown-Seite. Terravis betreibt
zwei Major-Versionen für jeweils zwei Jahre parallel.

## Verwandte Patterns

- [Eternal Lifetime Guarantee](EternalLifetimeGuarantee.md) — der Grenzfall mit unbegrenzter
  Frist.
- [Aggressive Obsolescence](AggressiveObsolescence.md) — das schwächere Gegenstück; *Limited
  Lifetime Guarantee* mischt Eigenschaften beider Muster.
- [Two in Production](TwoInProduction.md) — gleitendes statt fixes Fenster; beide Muster lassen
  sich kombinieren.
- [Experimental Preview](ExperimentalPreview.md) — die schwächste Zusage der Kategorie.
- [Version Identifier](VersionIdentifier.md) — eine befristete Zusage braucht in der Regel eine
  explizit benannte Version, auf die sie sich bezieht.
- [Semantic Versioning](SemanticVersioning.md) — definiert, was als Bruch der Zusage gilt.
- [API Description](../foundation/APIDescription.md),
  [Service Level Agreement](../quality/ServiceLevelAgreement.md) — tragen das konkrete
  Ablaufdatum.

## Bezug zu Kubernetes / KRM

Kubernetes erfüllt das Muster, ersetzt aber das *absolute Datum* durch eine **relative,
projektweite Policy**. Die Deprecation Policy des Projekts formuliert Mindestfristen, gemessen
in Releases und Monaten, ab dem Zeitpunkt der Deprecation:

| Stufe | Mindestfrist nach Deprecation |
|---|---|
| GA (`v1`) | 12 Monate oder 3 Releases, je nachdem was länger ist |
| Beta (`v1beta1`) | 9 Monate oder 3 Releases, je nachdem was länger ist |
| Alpha (`v1alpha1`) | keine; Entfernung jederzeit möglich |

Zusätzlich gilt: Eine API-Version darf erst dann als deprecated markiert werden, wenn eine
Nachfolgeversion serviert wird, und ein Objekt muss in allen servierten Versionen ohne
Informationsverlust les- und schreibbar bleiben (siehe
[Two in Production](TwoInProduction.md)).

Die Zusage ist maschinenlesbar, nicht nur Prosa. Für eingebaute Typen tragen die Go-Structs
Generator-Marker, aus denen `APILifecycleIntroduced()`, `APILifecycleDeprecated()`,
`APILifecycleRemoved()` und `APILifecycleReplacement()` erzeugt werden:

```go
// staging/src/k8s.io/api/flowcontrol/v1beta3/types.go
// +k8s:prerelease-lifecycle-gen:introduced=1.26
// +k8s:prerelease-lifecycle-gen:replacement=flowcontrol.apiserver.k8s.io,v1,FlowSchema
type FlowSchema struct { /* ... */ }
```

Bemerkenswert ist, was hier *nicht* steht: weder `:deprecated=` noch `:removed=`. Beide werden
aus `introduced` abgeleitet — jeweils „plus drei Minor“. Die Frist ist damit nicht das Ergebnis
einer Einzelfallentscheidung, sondern der Default:

```go
// staging/src/k8s.io/api/flowcontrol/v1beta3/zz_generated.prerelease-lifecycle.go
func (in *FlowSchema) APILifecycleIntroduced() (major, minor int) { return 1, 26 }
func (in *FlowSchema) APILifecycleDeprecated() (major, minor int) { return 1, 29 }  // <- abgeleitet
func (in *FlowSchema) APILifecycleRemoved() (major, minor int)    { return 1, 32 }  // <- abgeleitet
```

`k8s.io/apiserver/pkg/endpoints/deprecation` wertet das gegen die laufende Serverversion aus,
und der API-Server hängt beim Registrieren der Routen einen Warning-Handler ein
(`AddWarningsHandler` in `installer.go`), der einen `Warning`-Header nach RFC 7234 zurückgibt:

```http
HTTP/1.1 200 OK
Content-Type: application/json
Warning: 299 - "flowcontrol.apiserver.k8s.io/v1beta3 FlowSchema is deprecated in v1.29+, unavailable in v1.32+; use flowcontrol.apiserver.k8s.io/v1 FlowSchema"
```

Der Statuscode bleibt `200` — die Antwort ist gültig, nur die Frist läuft. `kubectl` und alle
client-go-basierten Clients geben diese Header aus (`k8s.io/client-go/rest/warnings.go`). Der
Ablaufzeitpunkt reist damit *in jeder Antwort mit*, statt nur in der Dokumentation zu stehen —
das geht deutlich über die MAP-Beschreibung hinaus, die das Datum in
[API Description](../foundation/APIDescription.md) oder
[SLA](../quality/ServiceLevelAgreement.md) verortet.

CRDs bekommen denselben Mechanismus als Konfiguration statt als Codegenerierung — pro Eintrag in
`spec.versions[]`, mit frei formuliertem Text auf demselben Warning-Kanal:

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: widgets.example.com
spec:
  # ... group, names, scope
  versions:
    - name: v1beta1
      served: true                 # <- weiterhin erreichbar …
      storage: false
      deprecated: true             # <- … aber abgekündigt
      deprecationWarning: "example.com/v1beta1 Widget is deprecated; use example.com/v1 Widget"
      # ... schema
    - name: v1
      served: true
      storage: true
      # ... schema
```

Betriebsseitig ist die Restlaufzeit messbar. Die Metrik `apiserver_requested_deprecated_apis`
ist als `STABLE` deklariert — sie unterliegt selbst der Deprecation Policy — und ist ein Gauge
über die tatsächlich angefragten deprecated APIs:

```prometheus
apiserver_requested_deprecated_apis{group="flowcontrol.apiserver.k8s.io",version="v1beta3",resource="flowschemas",subresource="",removed_release="1.32"} 1
```

Das Label `removed_release` ist der entscheidende Teil: es trägt genau den Wert, den
`APILifecycleRemoved()` liefert. Ein Cluster-Betreiber kann damit vor einem Upgrade prüfen, ob
noch etwas auf eine Version zugreift, die im Zielrelease verschwindet — Governance-Traceability,
wie sie [Version Identifier](VersionIdentifier.md) als Force nennt, hier als Zeitreihe.

---
[← Index](../README.md) · [Kategorie Evolution](../meta/category-evolution.md) · [Quelle](https://microservice-api-patterns.org/patterns/evolution/LimitedLifetimeGuarantee)
