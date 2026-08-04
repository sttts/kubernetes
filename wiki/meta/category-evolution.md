---
title: Kategorie Evolution
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/evolution
---

# Kategorie: Evolution Patterns

## Worum es geht

Die *Evolution Patterns* behandeln Lebenszyklusfragen: Wie lange wird eine API unterstützt, wie
wird versioniert, wie werden brechende Änderungen kommuniziert, und wie hält man die Balance
zwischen Kompatibilität (den Clients zuliebe) und Erweiterbarkeit (dem Provider zuliebe)?

Die Quellseite legt eine Lesereihenfolge nahe, die zugleich die Entscheidungsreihenfolge ist:

1. **Versionieren wir überhaupt — und auf welcher Granularität?**
   → [Version Identifier](../evolution/VersionIdentifier.md)
2. **Reicht eine einfache Nummer, oder müssen wir Änderungsarten unterscheiden?**
   → [Semantic Versioning](../evolution/SemanticVersioning.md)
3. **Wie viele Versionen laufen parallel, und wie lange?**
   → [Two in Production](../evolution/TwoInProduction.md),
   [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md),
   [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md),
   [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md)
4. **Wie vermeiden wir verfrühte Zusagen bei Neuem?**
   → [Experimental Preview](../evolution/ExperimentalPreview.md)

Alle diese Zusagen werden in der [API Description](../foundation/APIDescription.md) festgehalten;
ohne sie hat die Kategorie keinen Ort.

## Inhalt

### Versionierung und Kompatibilität

| Pattern | Einzeiler |
|---|---|
| [Version Identifier](../evolution/VersionIdentifier.md) | Ein expliziter Versionsindikator in Beschreibung *und* Nachrichten — als [Metadata Element](../structure/MetadataElement.md) in Adresse, Protokoll-Header oder Nutzlast. |
| [Semantic Versioning](../evolution/SemanticVersioning.md) | Dreistelliges Schema `x.y.z` (major/minor/patch), damit Stakeholder Kompatibilität am Bezeichner ablesen können. |

### Lebenszyklus-Garantien

| Pattern | Einzeiler |
|---|---|
| [Two in Production](../evolution/TwoInProduction.md) | Genau zwei Versionen eines Endpunkts laufen parallel; Aktualisierung und Abkündigung erfolgen rollierend und überlappend. |
| [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) | Der Provider sagt zu, eine Version für einen festen Zeitraum nicht zu brechen, und versieht jede Version mit einem Ablaufdatum. |
| [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) | Der Provider sagt zu, eine veröffentlichte Version nie zu brechen oder abzuschalten. |
| [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) | Ein möglichst früh angekündigtes Abkündigungsdatum; die Teile bleiben verfügbar, aber nicht empfohlen, und verschwinden zum Stichtag. |
| [Experimental Preview](../evolution/ExperimentalPreview.md) | Zugang auf Best-Effort-Basis ohne Zusagen zu Funktionsumfang, Stabilität und Langlebigkeit — explizit kommuniziert. |

Die vier Garantie-Patterns bilden eine Skala des Provider-Commitments:
*Experimental Preview* (keins) → *Aggressive Obsolescence* (minimal, mit Frist) →
*Limited Lifetime Guarantee* (befristet) → *Eternal Lifetime Guarantee* (unbefristet).
*Two in Production* ist orthogonal dazu und betrifft die Parallelität, nicht die Dauer.

## Verwandte Wiki-Seiten

- [Kategorie Foundation](category-foundation.md) — [API Description](../foundation/APIDescription.md) als Trägerartefakt
- [Cheat Sheet](cheatsheet.md) — Abschnitt „API Support and Maintenance"
- [Patterns nach Phase](navigation-byphase.md) — Evolution-Patterns gehören zu „transition (go live)"
- [Tutorials](tutorials.md) — Tutorial 2, Schritt 5

## Bezug zu Kubernetes / KRM

Kubernetes hat die vollständigste und am weitesten formalisierte Umsetzung dieser Kategorie, die
man in der Praxis findet — sie liegt in der Deprecation Policy und den API-Änderungsregeln fest.

Der [Version Identifier](../evolution/VersionIdentifier.md) ist Teil der Adressierung *und* der
Nutzlast: `/apis/apps/v1/...` im Pfad und `apiVersion: apps/v1` im Objekt. Die Granularität ist
die Group-Version, nicht die einzelne Operation und nicht die gesamte API — ein Mittelweg, den
MAP als Option nennt. Anders als bei [Semantic Versioning](../evolution/SemanticVersioning.md)
gibt es aber keine drei Zahlen: Die Stabilitätsstufe steckt im Namen selbst (`v1alpha1`,
`v1beta2`, `v1`), und die Kompatibilitätsaussage ergibt sich aus der Stufe, nicht aus einem
Zahlenvergleich. Semantic Versioning gilt in Kubernetes für das *Release* (`v1.34.2`), nicht für
die API-Gruppen.

[Two in Production](../evolution/TwoInProduction.md) wird deutlich übererfüllt: Der Server liefert
alle nicht abgekündigten Versionen einer Gruppe gleichzeitig aus, und zwar *verlustfrei
konvertierend* über einen gemeinsamen internen Hub-Typ; `metadata.uid`, `resourceVersion` und der
gespeicherte Inhalt sind versionsübergreifend identisch. MAP verlangt für *Two in Production*
ausdrücklich nicht, dass die Versionen kompatibel sind — KRM verlangt genau das Gegenteil, weil
sonst Controller und `kubectl` bei unterschiedlichen bevorzugten Versionen inkonsistente Sichten
bekämen.

Die Garantien sind an die Stufe gebunden, was die MAP-Patterns exakt reproduziert:
*Experimental Preview* = Alpha (standardmäßig abgeschaltete Feature Gates, Entfernung in jedem
Release möglich, kein Upgrade-Pfad zugesichert). *Limited Lifetime Guarantee* = Beta (nach der
Deprecation Policy mindestens neun Monate oder drei Releases weiter unterstützt).
*Eternal Lifetime Guarantee* ist für GA-APIs faktisch erreicht: Eine stabile API-Version darf
innerhalb einer Major-Version nicht entfernt werden, und ein Kubernetes 2.0 gibt es nicht — `v1`
läuft praktisch unbefristet. *Aggressive Obsolescence* ist als Prozess ausgeführt: Abkündigung im
Release Notes und Deprecation Guide, `Warning`-Header (RFC 7234) auf Antworten auf abgekündigte
Endpunkte, die `kubectl` dem Nutzer anzeigt, dann Entfernung zum angekündigten Release.

Bewertung: Kategorie **genauso**, aber mit deutlich schärferen Regeln und einer ungewöhnlichen
Verschiebung — die Version ist in KRM keine Eigenschaft der Nachricht, sondern eine *Sicht auf
dasselbe persistierte Objekt*.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/evolution)
