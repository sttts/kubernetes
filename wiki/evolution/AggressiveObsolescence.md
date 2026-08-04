---
title: Aggressive Obsolescence
kategorie: Evolution
unterkategorie: Lifecycle Management
quelle: https://microservice-api-patterns.org/patterns/evolution/AggressiveObsolescence
---

# Aggressive Obsolescence

*a.k.a.* *Early Sunset*, *Planned Obsolescence*

**Kurzform:** Der Provider kündigt so früh wie möglich ein Abschaltdatum an, markiert die
betroffenen API-Teile als weiterhin verfügbar, aber nicht mehr empfohlen — und entfernt sie,
sobald die Frist abgelaufen ist.

## Kontext

Eine API ist veröffentlicht und entwickelt sich weiter; Funktionalität kommt hinzu, ändert sich
oder wird überflüssig. Der Provider will bestimmte Teile nicht mehr unterstützen, weil sie kaum
noch genutzt werden oder durch Alternativen abgelöst sind.

## Problem

Wie reduzieren API-Provider den Aufwand für die Pflege einer API oder ihrer Teile — Endpunkte,
Operationen, Nachrichtenrepräsentationen — bei zugesicherten Servicequalitäten?

## Forces

- **Minimierung des Wartungsaufwands**, insbesondere Begrenzung des Supports für alte Clients.
- **Reduktion erzwungener Client-Änderungen** in einem gegebenen Zeitraum.
- **Machtverhältnisse** zwischen Provider und Client anerkennen — wie viel Einfluss haben
  Clients auf Design und Evolution?
- **Kommerzielle Ziele und Randbedingungen**, etwa Auswirkungen auf einen
  [Pricing Plan](../quality/PricingPlan.md).

## Lösung

Ein Abschaltdatum für die gesamte API oder ihre obsoleten Teile so früh wie möglich ankündigen.
Die betroffenen Teile als verfügbar, aber nicht mehr empfohlen deklarieren, damit Clients
gerade genug Zeit für den Umstieg haben. Nach Ablauf der Frist werden die deprecated Teile und
ihr Support entfernt.

Zwei Eigenschaften unterscheiden das Muster von den übrigen Lifecycle-Mustern:

1. **Feinkörnigkeit** — es zielt nicht auf syntaktische Einheiten wie Versionen oder Endpunkte,
   sondern kann einzelne Repräsentationselemente betreffen. Damit sind weniger störende
   Änderungen möglich.
2. **Relative Fristen** — die Deprecation-Periode läuft ab dem Zeitpunkt der Markierung, nicht
   ab dem Release-Datum. [Two in Production](TwoInProduction.md) und
   [Limited Lifetime Guarantee](LimitedLifetimeGuarantee.md) arbeiten mit absoluten Zeitpunkten.

Das Muster ist sowohl proaktiv als auch reaktiv während der Wartung anwendbar.

## Beispiel

Ein Zahlungsdienstleister erlaubt die Identifikation von Konten über die alten,
länderspezifischen Konto- und Banknummern oder über IBAN. Weil IBAN der neue Standard ist und
die alten Nummern kaum noch genutzt werden, markiert der Provider sie in der Dokumentation als
deprecated, veröffentlicht eine Abkündigung auf der API-Doku-Seite und benachrichtigt die
registrierten Clients — mit einem Jahr Frist. Nach einem Jahr wird eine Implementierung ohne
Unterstützung der alten Nummern ausgerollt, die Attribute verschwinden aus der Dokumentation,
und entsprechende Aufrufe schlagen fehl.

## Konsequenzen

**Vorteile:**

- Der Provider bekommt seinen Wartungsaufwand tatsächlich zurück — Code kann gelöscht werden.
- Feingranular anwendbar: ein einzelnes Feld statt einer ganzen Version.
- Die Ankündigung ist eine schwache, günstige Zusage; sie bindet den Provider nur für die Frist.

**Nachteile / Kosten:**

- Clients tragen die Migrationslast und müssen Ankündigungen aktiv verfolgen.
- Kurze oder schlecht kommunizierte Fristen beschädigen die Beziehung — das Muster setzt
  Verhandlungsmacht des Providers voraus.
- Feingranulare Deprecations sind schwer zu verfolgen: ein Client muss wissen, welche *Felder*
  er nutzt, nicht nur welche Version.
- Ohne Nutzungsmessung weiß der Provider nicht, wen er abschneidet.

## Bekannte Verwendungen

Microsoft Graph arbeitet mit einer 24-Monats-Deprecation-Frist. Riot Games nutzt(e) das Muster
für die Riot Games API, eBay kombiniert es mit
[Semantic Versioning](SemanticVersioning.md): eine neue Major-Version bleibt zunächst
kompatibel, kündigt aber Funktionalität ab, und Clients sollen binnen sechs Monaten die
Abhängigkeit auflösen. Google wird im Zusammenhang mit eingestellten Online-Diensten (etwa
Google Wave) genannt. Auch lokale Programmier-APIs wie das Zend Framework verwenden das Muster.

## Verwandte Patterns

- [Two in Production](TwoInProduction.md) — das gleitende Fenster, das Aggressive Obsolescence
  offenhält; die beiden greifen ineinander.
- [Limited Lifetime Guarantee](LimitedLifetimeGuarantee.md) — die stärkere, absolute Zusage.
- [Eternal Lifetime Guarantee](EternalLifetimeGuarantee.md) — das Gegenmodell.
- [Experimental Preview](ExperimentalPreview.md) — noch schwächeres Commitment; Aggressive
  Obsolescence ist die nächste Stufe darüber.
- [Version Identifier](VersionIdentifier.md) — optional; das Muster funktioniert auch
  unterhalb der Versionsgranularität.
- [API Description](../foundation/APIDescription.md),
  [Service Level Agreement](../quality/ServiceLevelAgreement.md) — tragen die
  Deprecation-Metadaten.
- [Pricing Plan](../quality/PricingPlan.md) — Abkündigungen haben kommerzielle Wirkung.

## Bezug zu Kubernetes / KRM

Kubernetes wendet das Muster an, aber **nur unterhalb von GA** — und ersetzt die einzelne
Provider-Ankündigung durch eine projektweite, maschinenlesbare Policy. Die Fristen sind relativ
zum Deprecation-Zeitpunkt, genau wie das Muster es vorsieht:

| Stufe | Frist ab Deprecation | Tatsächlich entfernt? |
|---|---|---|
| Alpha | keine | ja, regelmäßig |
| Beta | 9 Monate oder 3 Releases, das Längere | ja — z.B. `flowcontrol.apiserver.k8s.io/v1beta1`…`v1beta3`, `extensions/v1beta1`, `apps/v1beta1` |
| GA | 12 Monate oder 3 Releases, das Längere | praktisch nie, siehe [Eternal Lifetime Guarantee](EternalLifetimeGuarantee.md) |

Die Frist muss dabei nicht einmal von Hand gesetzt werden: fehlen die Marker `:deprecated=` und
`:removed=`, leitet der Generator sie aus `:introduced=` ab — jeweils „plus drei Minor“. Für
`flowcontrol.apiserver.k8s.io/v1beta3` ergibt `introduced=1.26` damit automatisch
`deprecated=1.29` und `removed=1.32`.

Die Ankündigung ist kein Blogpost, sondern Teil der API. Für eingebaute Typen werden aus den
Markern `+k8s:prerelease-lifecycle-gen:deprecated=`, `:removed=` und `:replacement=` die
Methoden `APILifecycleDeprecated()`, `APILifecycleRemoved()` und `APILifecycleReplacement()`
generiert. `k8s.io/apiserver/pkg/endpoints/deprecation` wertet sie gegen die laufende
Serverversion aus, und `installer.go` hängt an jede betroffene Route einen Warning-Handler. Der
Client erhält daraufhin einen `Warning`-Header:

```http
HTTP/1.1 200 OK
Warning: 299 - "flowcontrol.apiserver.k8s.io/v1beta3 FlowSchema is deprecated in v1.29+, unavailable in v1.32+; use flowcontrol.apiserver.k8s.io/v1 FlowSchema"
```

Damit sind alle drei Informationen — deprecated seit, entfällt ab, Ersatz — in jeder einzelnen
Antwort enthalten. `kubectl` registriert dafür in `cmd.go` einen
`rest.NewWarningWriter(o.IOStreams.ErrOut, ...)` mit Deduplizierung, sodass die Abkündigung im
Terminal und in CI-Logs auftaucht statt nur in der Doku:

```console
$ kubectl get flowschemas.v1beta3.flowcontrol.apiserver.k8s.io
Warning: flowcontrol.apiserver.k8s.io/v1beta3 FlowSchema is deprecated in v1.29+, unavailable in v1.32+; use flowcontrol.apiserver.k8s.io/v1 FlowSchema
NAME     PRIORITYLEVEL   MATCHINGPRECEDENCE   DISTINGUISHERMETHOD   AGE   MISSINGPL
exempt   exempt          1                    <none>                12d   False
# ^ stdout, die Warnung darüber geht auf stderr — Pipes bleiben sauber
```

CRD-Autoren bekommen denselben Kanal deklarativ: pro Eintrag in `spec.versions[]` gibt es
`deprecated: true` und einen freien `deprecationWarning`-Text.

Die von der Quellseite geforderte Nutzungsmessung liefert der API-Server selbst. Die stabile
Metrik `apiserver_requested_deprecated_apis` ist ein Gauge über tatsächlich angefragte
deprecated APIs mit den Labels `group`, `version`, `resource`, `subresource` und
`removed_release`. Ein Betreiber kann vor einem Upgrade exakt bestimmen, welche Workloads noch
auf einer im Zielrelease verschwindenden Version arbeiten — der Provider weiß also, wen er
abschneidet, was in der MAP-Beschreibung als offenes Problem stehen bleibt.

Die **Feinkörnigkeit** löst KRM anders als MAP. Einzelne Felder einer GA-Ressource werden nur
in der Dokumentation als deprecated markiert und faktisch nie entfernt — die Round-Trip-
Anforderung aus [Two in Production](TwoInProduction.md) verbietet es. Feingranulare Abschaltung
findet stattdessen über Feature Gates statt (`--feature-gates`, Stufe `Deprecated` in
`k8s.io/component-base/featuregate`). Der Ablauf steht als versionierte Leiter im Code, samt
Enddatum im Kommentar:

```go
// pkg/features/kube_features.go
AllowOverwriteTerminationGracePeriodSeconds: {
	{Version: version.MustParse("1.0"),  Default: true,  PreRelease: featuregate.GA},
	{Version: version.MustParse("1.32"), Default: false, PreRelease: featuregate.Deprecated},                       // <- default-off, noch einschaltbar
	{Version: version.MustParse("1.35"), Default: false, PreRelease: featuregate.Deprecated, LockToDefault: true},  // remove in 1.38
},
```

Drei Übergänge im Abstand von je drei Releases: in 1.32 kippt der Default auf `false`, das Gate
bleibt aber einschaltbar; in 1.35 fixiert `LockToDefault` es endgültig; in 1.38 verschwindet es.
Das ist exakt die relative Frist des Musters — nur eben im Quelltext deklariert statt in einer
Ankündigung. Das Muster ist also erfüllt, aber der Hebel ist Verhaltens- statt Schemaänderung.

---
[← Index](../README.md) · [Kategorie Evolution](../meta/category-evolution.md) · [Quelle](https://microservice-api-patterns.org/patterns/evolution/AggressiveObsolescence)
