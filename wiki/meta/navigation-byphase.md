---
title: Patterns nach Phase
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/byphase
---

# Patterns nach Phase

## Worum es geht

Der *Phasen*-Filter ordnet die Patterns dem Zeitpunkt im Projektverlauf zu, an dem sie typischerweise
gebraucht werden. MAP verwendet dazu die vier Phasen des Unified Process, jeweils mit einer
agilen Übersetzung in Klammern: Inception (Sprint 0), Elaboration (Spikes), Construction
(Entwicklungs-Iterationen), Transition (Go-Live).

Die Zuordnung ist keine Vorschrift zur Serialisierung — die Quellseiten betonen mehrfach, dass
API-Design ein „wicked problem" ist und sich nicht in eine Phasenfolge zwingen lässt. Der Filter
ist ein Einstiegshilfsmittel, kein Prozessmodell. Auffällig ist trotzdem, wie sauber sich die
[fünf Kategorien](overview.md) auf die Phasen verteilen.

## Inhalt

### Inception (Sprint 0) — Scoping und Geschäftsmodell

| Pattern | Kategorie |
|---|---|
| [Frontend Integration](../foundation/FrontendIntegration.md) | Foundation |
| [Backend Integration](../foundation/BackendIntegration.md) | Foundation |
| [Public API](../foundation/PublicAPI.md) | Foundation |
| [Community API](../foundation/CommunityAPI.md) | Foundation |
| [Solution-Internal API](../foundation/SolutionInternalAPI.md) | Foundation |
| [API Description](../foundation/APIDescription.md) | Foundation |
| [Pricing Plan](../quality/PricingPlan.md) | Quality |
| [Service Level Agreement](../quality/ServiceLevelAgreement.md) | Quality |

Praktisch die gesamte Foundation-Kategorie plus die beiden kommerziellen Governance-Patterns.

### Elaboration (Spikes) — Endpunktschnitt und Rollen

| Pattern | Kategorie |
|---|---|
| [Processing Resource](../responsibility/ProcessingResource.md) | Responsibility |
| [Information Holder Resource](../responsibility/InformationHolderResource.md) | Responsibility |
| [Operational Data Holder](../responsibility/OperationalDataHolder.md) | Responsibility |
| [Master Data Holder](../responsibility/MasterDataHolder.md) | Responsibility |
| [Reference Data Holder](../responsibility/ReferenceDataHolder.md) | Responsibility |
| [Link Lookup Resource](../responsibility/LinkLookupResource.md) | Responsibility |
| [Data Transfer Resource](../responsibility/DataTransferResource.md) | Responsibility |
| [API Key](../structure/APIKey.md) | Structure |
| [Error Report](../structure/ErrorReport.md) | Structure |
| [Context Representation](../structure/ContextRepresentation.md) | Structure |

Endpunkt-Rollen und -Typen; dazu die drei Spezialrepräsentationen, weil sie
querschnittlich sind und früh festgelegt werden sollten.

### Construction (Entwicklungs-Iterationen) — Nachrichtendesign

| Pattern | Kategorie |
|---|---|
| [State Creation Operation](../responsibility/StateCreationOperation.md) | Responsibility |
| [State Transition Operation](../responsibility/StateTransitionOperation.md) | Responsibility |
| [Retrieval Operation](../responsibility/RetrievalOperation.md) | Responsibility |
| [Computation Function](../responsibility/ComputationFunction.md) | Responsibility |
| [Atomic Parameter](../structure/AtomicParameter.md) | Structure |
| [Atomic Parameter List](../structure/AtomicParameterList.md) | Structure |
| [Parameter Tree](../structure/ParameterTree.md) | Structure |
| [Parameter Forest](../structure/ParameterForest.md) | Structure |
| [Data Element](../structure/DataElement.md) | Structure |
| [Metadata Element](../structure/MetadataElement.md) | Structure |
| [Id Element](../structure/IdElement.md) | Structure |
| [Link Element](../structure/LinkElement.md) | Structure |
| [Embedded Entity](../quality/EmbeddedEntity.md) | Quality |
| [Linked Information Holder](../quality/LinkedInformationHolder.md) | Quality |
| [Pagination](../quality/Pagination.md) | Quality |
| [Wish List](../quality/WishList.md) | Quality |
| [Wish Template](../quality/WishTemplate.md) | Quality |
| [Conditional Request](../quality/ConditionalRequest.md) | Quality |
| [Request Bundle](../quality/RequestBundle.md) | Quality |

Die mit Abstand größte Gruppe: die komplette Structure-Kategorie, die
Operationsverantwortungen und alle Parsimony-Patterns.

### Transition (Go-Live) — Betrieb, Governance, Lebenszyklus

| Pattern | Kategorie |
|---|---|
| [API Description](../foundation/APIDescription.md) | Foundation |
| [Rate Limit](../quality/RateLimit.md) | Quality |
| [Pricing Plan](../quality/PricingPlan.md) | Quality |
| [Service Level Agreement](../quality/ServiceLevelAgreement.md) | Quality |
| [Version Identifier](../evolution/VersionIdentifier.md) | Evolution |
| [Semantic Versioning](../evolution/SemanticVersioning.md) | Evolution |
| [Two in Production](../evolution/TwoInProduction.md) | Evolution |
| [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) | Evolution |
| [Experimental Preview](../evolution/ExperimentalPreview.md) | Evolution |
| [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) | Evolution |
| [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) | Evolution |

Die komplette Evolution-Kategorie. *API Description*, *Pricing Plan* und *Service Level Agreement*
tauchen zum zweiten Mal auf: erst als Absichtserklärung in Sprint 0, dann als tatsächlich
veröffentlichtes Artefakt.

## Verwandte Wiki-Seiten

- [Patterns nach Scope](navigation-byscope.md) · [nach Rolle](navigation-byrole.md) · [nach Qualität](navigation-byquality.md)
- [Cheat Sheet](cheatsheet.md) — dieselbe zeitliche Ordnung, aber problemgetrieben
- [Tutorials](tutorials.md) — Tutorial 2 durchläuft genau diese Abfolge in fünf Schritten

## Bezug zu Kubernetes / KRM

Für Kubernetes-APIs verschieben sich die Gewichte deutlich, und der Grund ist der KEP-Prozess.
Was MAP über Inception und Elaboration verteilt, wird in Kubernetes in einem einzigen Artefakt
vorweggenommen: Ein Kubernetes Enhancement Proposal enthält Motivation, API-Schema, Alternativen,
Test- und Graduierungsplan (Alpha → Beta → GA) sowie die Rückbau-Strategie, bevor die erste Zeile
Typdefinition gemergt wird. Die Evolution-Entscheidungen — bei MAP eine Transition-Aufgabe —
fallen in Kubernetes damit *ganz am Anfang*: Die Wahl `v1alpha1` gegenüber `v1` ist eine
Inception-Entscheidung mit voller Kenntnis der späteren Konsequenzen (Feature Gate standardmäßig
aus, Entfernung jederzeit möglich, keine Konvertierungsgarantie).

Die Elaboration-Phase im MAP-Sinn — Endpunktrollen und -typen wählen — entfällt weitgehend, weil
die Rolle durch das KRM-Modell vorgegeben ist (siehe
[Kategorie Responsibility](category-responsibility.md)): Man entwirft eine Information Holder
Resource mit `spec`/`status`, oder man begründet sehr ausführlich, warum nicht. Die eigentliche
Designarbeit konzentriert sich auf Construction, also auf Feldnamen, Typen, Optionalität,
Defaulting, Validierung und Kompatibilität — und dort gelten mit `api-conventions.md` und
`api_changes.md` sehr detaillierte, verbindliche Regeln, für die MAP kein Gegenstück hat.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/byphase)
