---
title: Patterns nach Rolle
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/byrole
---

# Patterns nach Rolle

## Worum es geht

Der *Rollen*-Filter ordnet die Patterns danach, wer sie liest. MAP unterscheidet vier
Stakeholder-Rollen — API Client Developer, API Designer, API Developer und API Product Owner.

Die Aufteilung folgt der Blickrichtung: Der *Client Developer* konsumiert und will wissen, womit
er rechnen muss (Fehler, Limits, Versionen, Nachrichtengrößen). Der *Designer* entscheidet über
Schnitt und Verantwortung. Der *Developer* implementiert die Nachrichtenstrukturen. Der
*Product Owner* verantwortet Sichtbarkeit, Kosten und Lebenszyklus. Client Developer und
Product Owner sitzen dabei an den entgegengesetzten Enden derselben Verträge — sie lesen zu großen
Teilen dieselben Patterns aus verschiedenen Perspektiven.

## Inhalt

### API Client Developer

| Pattern | Kategorie |
|---|---|
| [Frontend Integration](../foundation/FrontendIntegration.md) | Foundation |
| [API Description](../foundation/APIDescription.md) | Foundation |
| [API Key](../structure/APIKey.md) | Structure |
| [Error Report](../structure/ErrorReport.md) | Structure |
| [Embedded Entity](../quality/EmbeddedEntity.md) | Quality |
| [Linked Information Holder](../quality/LinkedInformationHolder.md) | Quality |
| [Pagination](../quality/Pagination.md) | Quality |
| [Wish List](../quality/WishList.md) | Quality |
| [Wish Template](../quality/WishTemplate.md) | Quality |
| [Request Bundle](../quality/RequestBundle.md) | Quality |
| [Rate Limit](../quality/RateLimit.md) | Quality |
| [Version Identifier](../evolution/VersionIdentifier.md) | Evolution |

### API Designer

| Pattern | Kategorie |
|---|---|
| [Frontend Integration](../foundation/FrontendIntegration.md) | Foundation |
| [Backend Integration](../foundation/BackendIntegration.md) | Foundation |
| [Public API](../foundation/PublicAPI.md) | Foundation |
| [Community API](../foundation/CommunityAPI.md) | Foundation |
| [Solution-Internal API](../foundation/SolutionInternalAPI.md) | Foundation |
| [Processing Resource](../responsibility/ProcessingResource.md) | Responsibility |
| [Information Holder Resource](../responsibility/InformationHolderResource.md) | Responsibility |
| [Operational Data Holder](../responsibility/OperationalDataHolder.md) | Responsibility |
| [Master Data Holder](../responsibility/MasterDataHolder.md) | Responsibility |
| [Reference Data Holder](../responsibility/ReferenceDataHolder.md) | Responsibility |
| [Link Lookup Resource](../responsibility/LinkLookupResource.md) | Responsibility |
| [Data Transfer Resource](../responsibility/DataTransferResource.md) | Responsibility |
| [State Creation Operation](../responsibility/StateCreationOperation.md) | Responsibility |
| [State Transition Operation](../responsibility/StateTransitionOperation.md) | Responsibility |
| [Retrieval Operation](../responsibility/RetrievalOperation.md) | Responsibility |
| [Computation Function](../responsibility/ComputationFunction.md) | Responsibility |
| [Context Representation](../structure/ContextRepresentation.md) | Structure |

Die Designer-Liste ist praktisch deckungsgleich mit Foundation + Responsibility: Wer schneidet,
entscheidet über Rollen und Verantwortungen, nicht über Feldnamen.

### API Developer

| Pattern | Kategorie |
|---|---|
| [Atomic Parameter](../structure/AtomicParameter.md) | Structure |
| [Atomic Parameter List](../structure/AtomicParameterList.md) | Structure |
| [Parameter Tree](../structure/ParameterTree.md) | Structure |
| [Parameter Forest](../structure/ParameterForest.md) | Structure |
| [Data Element](../structure/DataElement.md) | Structure |
| [Metadata Element](../structure/MetadataElement.md) | Structure |
| [Id Element](../structure/IdElement.md) | Structure |
| [Link Element](../structure/LinkElement.md) | Structure |
| [API Key](../structure/APIKey.md) | Structure |
| [Error Report](../structure/ErrorReport.md) | Structure |
| [Embedded Entity](../quality/EmbeddedEntity.md) | Quality |
| [Linked Information Holder](../quality/LinkedInformationHolder.md) | Quality |
| [Pagination](../quality/Pagination.md) | Quality |
| [Wish List](../quality/WishList.md) | Quality |
| [Wish Template](../quality/WishTemplate.md) | Quality |
| [Conditional Request](../quality/ConditionalRequest.md) | Quality |
| [Request Bundle](../quality/RequestBundle.md) | Quality |
| [Rate Limit](../quality/RateLimit.md) | Quality |

### API Product Owner

| Pattern | Kategorie |
|---|---|
| [Frontend Integration](../foundation/FrontendIntegration.md) | Foundation |
| [Backend Integration](../foundation/BackendIntegration.md) | Foundation |
| [Public API](../foundation/PublicAPI.md) | Foundation |
| [Community API](../foundation/CommunityAPI.md) | Foundation |
| [Solution-Internal API](../foundation/SolutionInternalAPI.md) | Foundation |
| [API Description](../foundation/APIDescription.md) | Foundation |
| [Pricing Plan](../quality/PricingPlan.md) | Quality |
| [Rate Limit](../quality/RateLimit.md) | Quality |
| [Service Level Agreement](../quality/ServiceLevelAgreement.md) | Quality |
| [Version Identifier](../evolution/VersionIdentifier.md) | Evolution |
| [Semantic Versioning](../evolution/SemanticVersioning.md) | Evolution |
| [Two in Production](../evolution/TwoInProduction.md) | Evolution |
| [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) | Evolution |
| [Experimental Preview](../evolution/ExperimentalPreview.md) | Evolution |
| [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) | Evolution |
| [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) | Evolution |

## Verwandte Wiki-Seiten

- [Patterns nach Scope](navigation-byscope.md) · [nach Phase](navigation-byphase.md) · [nach Qualität](navigation-byquality.md)
- [Cheat Sheet](cheatsheet.md)
- [Überblick MAP](overview.md)

## Bezug zu Kubernetes / KRM

Die Rollentrennung existiert im Kubernetes-Projekt ähnlich, allerdings mit anderen Namen und einer
zusätzlichen Rolle, die MAP nicht kennt. „API Designer" und „API Developer" fallen im KEP-Prozess
weitgehend zusammen und werden von den zuständigen SIGs (SIG API Machinery für Querschnitt, die
Fach-SIGs für einzelne Gruppen) sowie den API Reviewers wahrgenommen. Die Rolle des „API Product
Owner" mit Preis- und SLA-Verantwortung fehlt in einem Open-Source-Projekt vollständig; ihre
Lebenszyklus-Aufgaben übernimmt die Deprecation Policy als Projektregel, nicht als
Produktentscheidung.

Die zusätzliche Rolle ist der *Controller-Autor*, der weder reiner Client noch reiner Provider
ist: Er konsumiert über Informer/Watch (Client-Sicht) und schreibt über die
`/status`-Subresource zurück (Provider-Sicht). Für ihn ist die relevanteste
Patternmenge eine andere als für alle vier MAP-Rollen — [Pagination](../quality/Pagination.md)
und [Conditional Request](../quality/ConditionalRequest.md) treffen ihn über `resourceVersion`
und Watch-Bookmarks, und [Version Identifier](../evolution/VersionIdentifier.md) über die
Konvertierungssemantik.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/byrole)
