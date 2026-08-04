---
title: Semantic Versioning
kategorie: Evolution
unterkategorie: Versionierung
quelle: https://microservice-api-patterns.org/patterns/evolution/SemanticVersioning
---

# Semantic Versioning

*a.k.a.* *Version Number Triplet*, *Three-Number Version*

**Kurzform:** Der [Version Identifier](VersionIdentifier.md) wird als hierarchisches Tripel
`major.minor.patch` strukturiert, sodass Stakeholder allein am Versionsvergleich ablesen
können, ob zwei Stände kompatibel sind.

## Kontext

Ein einzelner, unstrukturierter Versionsindikator sagt nichts über die Tragweite einer
Änderung. Jeder Client müsste die Differenz selbst analysieren, um zu entscheiden, ob eine
Migration nötig ist. Provider brauchen umgekehrt eine Regel, um zu bestimmen, ob eine geplante
Änderung die Zusagen aus [Limited Lifetime Guarantee](LimitedLifetimeGuarantee.md) oder
[Two in Production](TwoInProduction.md) verletzt.

## Problem

Wie können Stakeholder API-Versionen vergleichen und sofort erkennen, ob sie kompatibel sind?

## Forces

- **Minimaler Aufwand zur Erkennung von Inkompatibilität**, vor allem auf Client-Seite.
- **Klarheit über die Tragweite** einer Änderung.
- **Saubere Trennung** von Änderungen unterschiedlicher Wirkung in getrennte Stellen der Zahl.
- **Beherrschbarer Governance-Aufwand** — Zahl paralleler Versionen, Branches, Quality Gates.
- **Klarheit über die Evolutionslinie** der API (Vergangenheit, Gegenwart, Zukunft).

## Lösung

Ein dreistelliges Schema `x.y.z` einführen:

| Stelle | Erhöht sich bei | Kompatibilität |
|---|---|---|
| `major` (`x`) | inkompatibler Vertragsänderung | bricht Clients |
| `minor` (`y`) | rückwärtskompatibler Erweiterung | alte Clients laufen weiter |
| `patch` (`z`) | Fehlerkorrektur ohne Vertragsänderung | transparent |

Das Muster ist bewusst technologieunabhängig: es regelt nur, *welche Zahl sich wann ändert*,
nicht Serialisierungsformat oder Transport.

## Beispiel

Das Beispiel der Quellseite, ein Börsendaten-Anbieter:

- `1.0.0` — Suchoperation über Aktiensymbole, Preise in USD.
- `1.1.0` — die Suche akzeptiert optional einen Zeitraum für historische Kurse. Ohne Zeitraum
  gilt das alte Verhalten: rückwärtskompatibel, also `minor`.
- `1.1.1` — Bugfix: die Suche fand bisher nur Präfixtreffer statt Teilstrings.
- `2.0.0` — die Antwort enthält zusätzlich ein `currency`-Attribut, was die Bedeutung des
  Preisfeldes ändert: inkompatibel, also `major`.

## Konsequenzen

**Vorteile:**

- Kompatibilitätsentscheidung ohne Diff-Analyse — ein Zahlenvergleich reicht.
- Die Regel diszipliniert auch den Provider: die Frage „ist das breaking?“ muss vor dem Release
  beantwortet werden.
- Lässt sich mit allen Lifecycle-Mustern kombinieren.

**Nachteile / Kosten:**

- Die Einstufung einer Änderung ist Ermessenssache; „kompatibel“ ist auf semantischer Ebene
  schwerer zu belegen als auf syntaktischer.
- Werden alle drei Stellen nach außen exponiert (z.B. im URI), erzwingt schon ein Patch eine
  Client-Änderung. Deshalb exponieren viele Anwender nur `major` bzw. `major.minor`.
- Häufige `major`-Bumps erzeugen Migrationsdruck und Versionswildwuchs.

## Bekannte Verwendungen

Eine große Schweizer Finanzinstitution und das Terravis-Projekt exponieren `major` und `minor`
im Namespace, während die Fix-Version nur in der Vertragsdokumentation steht; Terravis erreicht
Kompatibilität zwischen Minor-Versionen durch Transformation ein- und ausgehender Nachrichten.
Die eBay-REST-API kombiniert das Tripel mit [Aggressive Obsolescence](AggressiveObsolescence.md).
Außerhalb von Remote-APIs ist SemVer 2.0.0 der De-facto-Standard für Bibliotheken; Snowplow
hat mit „SchemaVer“ eine auf Datenschemata zugeschnittene Variante beschrieben (Model /
Revision / Addition).

## Verwandte Patterns

- [Version Identifier](VersionIdentifier.md) — zwingende Voraussetzung; *Semantic Versioning*
  gibt ihm nur Struktur.
- [Two in Production](TwoInProduction.md) — vollständig kompatible Stände (Patch) dürfen eine
  aktive Version ersetzen, ohne gegen die Zwei-Versionen-Regel zu verstoßen.
- [Limited Lifetime Guarantee](LimitedLifetimeGuarantee.md),
  [Eternal Lifetime Guarantee](EternalLifetimeGuarantee.md),
  [Aggressive Obsolescence](AggressiveObsolescence.md),
  [Experimental Preview](ExperimentalPreview.md) — unterscheiden sich im Commitment-Grad;
  SemVer liefert die gemeinsame Sprache dafür.
- [API Description](../foundation/APIDescription.md),
  [Service Level Agreement](../quality/ServiceLevelAgreement.md) — tragen die
  Versionierungsinformation zum Client.

## Bezug zu Kubernetes / KRM

Hier weicht KRM bewusst ab: **API-Gruppen in Kubernetes verwenden kein SemVer.** Statt
`x.y.z` gibt es eine kurze, ordinale Reifegradskala, deren Schema
`k8s.io/apimachinery/pkg/version` als Regex kodiert: `v<n>`, `v<n>beta<m>`, `v<n>alpha<m>`.
Es gibt keine Minor- und keine Patch-Stelle, weil kompatible Erweiterungen *innerhalb*
derselben Version stattfinden dürfen und nicht angezeigt werden müssen — additive, optionale
Felder sind laut API-Konventionen erlaubt, ohne die Version zu ändern. Damit entfällt der
Anwendungsfall, für den `minor`/`patch` gedacht sind.

Die Skala kodiert dafür etwas, das SemVer nicht kennt: das **Stabilitätsversprechen**.

| Stufe | Beispiel aus dem Repo | Zusage (Kubernetes Deprecation Policy) |
|---|---|---|
| Alpha | `resource.k8s.io/v1alpha3` | keine; jederzeit änder- und entfernbar, per Default nicht serviert |
| Beta | `resource.k8s.io/v1beta1`, `v1beta2` | Support für mindestens 9 Monate oder 3 Releases nach Deprecation (das Längere) |
| GA / stable | `resource.k8s.io/v1` | Support für mindestens 12 Monate oder 3 Releases nach Deprecation (das Längere); praktisch nie entfernt, siehe [Eternal Lifetime Guarantee](EternalLifetimeGuarantee.md) |

Die Gruppe `resource.k8s.io` durchläuft die Skala gerade sichtbar: `staging/src/k8s.io/api/resource/`
enthält alle vier Verzeichnisse nebeneinander. Innerhalb einer Stufe zählt die Ziffer hoch
(`v1beta1` → `v1beta2`), und der Übergang `v1beta2` → `v1` ist die einzige „Major-artige“
Bewegung, die es üblicherweise gibt. Ein Sprung auf `v2` einer bestehenden Gruppe ist extrem
selten — statt inkompatibel zu versionieren, wird eher eine neue Ressource oder ein neues Feld
eingeführt.

`CompareKubeAwareVersionStrings` implementiert die Ordnung explizit: GA schlägt Beta schlägt
Alpha, innerhalb einer Stufe gewinnt die höhere Zahl. Diese Ordnung ist funktional relevant,
denn sie bestimmt unter anderem, welche Version Discovery-Clients und `kubectl` bevorzugen. Was
ein Client tatsächlich zu sehen bekommt, ist deshalb eine flache Liste von Gruppenversionen
ohne jede Zahlenhierarchie — hier die Default-Konfiguration aus
`pkg/controlplane/instance.go`:

```console
$ kubectl api-versions
admissionregistration.k8s.io/v1
apps/v1
autoscaling/v1
autoscaling/v2                    # <- keine Minor-/Patch-Stelle, nur v1 neben v2
batch/v1
certificates.k8s.io/v1
coordination.k8s.io/v1
discovery.k8s.io/v1
events.k8s.io/v1
flowcontrol.apiserver.k8s.io/v1
networking.k8s.io/v1
node.k8s.io/v1
policy/v1
rbac.authorization.k8s.io/v1
resource.k8s.io/v1
scheduling.k8s.io/v1
storage.k8s.io/v1
v1                                # <- Core-Gruppe, leerer Gruppenname
```

`autoscaling/v1` und `autoscaling/v2` sind dabei der seltene Fall einer echten
Major-Nachfolge — beide werden parallel serviert (siehe
[Two in Production](TwoInProduction.md)).

Semver-artig ist dagegen die **Release-Version von Kubernetes selbst** (`1.34.2`,
`v1.35.0-alpha.1`): `minor` steht für den quartalsweisen Feature-Release, `patch` für
Bugfix-Releases. Beide Zahlensysteme stehen nebeneinander und bedeuten Verschiedenes:

```console
$ kubectl version -o yaml
serverVersion:
  major: "1"
  minor: "34"                 # <- Release-Version: semver-artig, mit Patch-Stelle
  gitVersion: v1.34.2
  # ... gitCommit, buildDate, goVersion, platform
```

```yaml
apiVersion: apps/v1           # <- API-Version: ordinal, ohne Patch-Stelle
kind: Deployment
```

Auch hier ist die Analogie unvollständig — die Major-Stelle steht seit
2015 auf `1` und es gibt keine Absicht, sie zu erhöhen, weil ein Major-Bump genau das
Kompatibilitätsversprechen brechen würde, das die Deprecation Policy schützt. Die Verbindung
zwischen Release- und API-Version stellen die Generator-Marker
`+k8s:prerelease-lifecycle-gen:introduced=1.26`, `:deprecated=`, `:removed=` und
`:replacement=` her, aus denen die Methoden `APILifecycleIntroduced()`,
`APILifecycleDeprecated()`, `APILifecycleRemoved()` und `APILifecycleReplacement()` generiert
werden — siehe [Aggressive Obsolescence](AggressiveObsolescence.md).

---
[← Index](../README.md) · [Kategorie Evolution](../meta/category-evolution.md) · [Quelle](https://microservice-api-patterns.org/patterns/evolution/SemanticVersioning)
