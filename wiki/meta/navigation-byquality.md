---
title: Patterns nach Qualitätsattribut
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/byquality
---

# Patterns nach Qualitätsattribut

## Worum es geht

Der *Qualitäts*-Filter dreht die Perspektive um: Nicht „woran arbeite ich?", sondern „welche
Eigenschaft macht mir Sorgen?". MAP ordnet jedes Pattern den Qualitätsattributen zu, die es
adressiert — sei es fördernd oder als Trade-off. Der Filter deckt zwölf Attribute ab, darunter
neben den klassischen ISO-25010-Qualitäten (Interoperabilität, Wartbarkeit, Performance,
Zuverlässigkeit, Sicherheit, Portabilität, Usability) auch die architekturspezifischen Größen
*Kopplung*, *Granularität*, *Konsistenz*, *Flexibilität* und *Kosten*.

Die Zuordnung ist bewusst nicht exklusiv: Ein Pattern taucht mehrfach auf, weil es typischerweise
eine Qualität verbessert und eine andere verschlechtert. Der [API Key](../structure/APIKey.md)
etwa erscheint unter Sicherheit, Zuverlässigkeit, Performance *und* Kosten — er ermöglicht
Zurechenbarkeit und damit Abrechnung und Drosselung, kostet aber Verwaltungsaufwand.

## Inhalt

Kurzschreibweise: F = [Foundation](category-foundation.md),
R = [Responsibility](category-responsibility.md), S = [Structure](category-structure.md),
Q = [Quality](category-quality.md), E = [Evolution](category-evolution.md).

### Kopplung

[Frontend Integration](../foundation/FrontendIntegration.md) (F) ·
[Backend Integration](../foundation/BackendIntegration.md) (F) ·
[Processing Resource](../responsibility/ProcessingResource.md) (R) ·
[Data Transfer Resource](../responsibility/DataTransferResource.md) (R) ·
[Link Lookup Resource](../responsibility/LinkLookupResource.md) (R) ·
[Data Element](../structure/DataElement.md) (S) ·
[Metadata Element](../structure/MetadataElement.md) (S) ·
[Link Element](../structure/LinkElement.md) (S) ·
[Context Representation](../structure/ContextRepresentation.md) (S)

### Granularität

[Processing Resource](../responsibility/ProcessingResource.md) (R) ·
[Information Holder Resource](../responsibility/InformationHolderResource.md) (R) ·
[Operational Data Holder](../responsibility/OperationalDataHolder.md) (R) ·
[Master Data Holder](../responsibility/MasterDataHolder.md) (R) ·
[Reference Data Holder](../responsibility/ReferenceDataHolder.md) (R) ·
[Link Lookup Resource](../responsibility/LinkLookupResource.md) (R) ·
[Data Transfer Resource](../responsibility/DataTransferResource.md) (R) ·
[State Creation Operation](../responsibility/StateCreationOperation.md) (R) ·
[State Transition Operation](../responsibility/StateTransitionOperation.md) (R) ·
[Retrieval Operation](../responsibility/RetrievalOperation.md) (R) ·
[Computation Function](../responsibility/ComputationFunction.md) (R)

### Performance

[Computation Function](../responsibility/ComputationFunction.md) (R) ·
[State Creation Operation](../responsibility/StateCreationOperation.md) (R) ·
[State Transition Operation](../responsibility/StateTransitionOperation.md) (R) ·
[Retrieval Operation](../responsibility/RetrievalOperation.md) (R) ·
[API Key](../structure/APIKey.md) (S) ·
[Embedded Entity](../quality/EmbeddedEntity.md) (Q) ·
[Linked Information Holder](../quality/LinkedInformationHolder.md) (Q) ·
[Pagination](../quality/Pagination.md) (Q) ·
[Wish List](../quality/WishList.md) (Q) ·
[Wish Template](../quality/WishTemplate.md) (Q) ·
[Conditional Request](../quality/ConditionalRequest.md) (Q) ·
[Request Bundle](../quality/RequestBundle.md) (Q) ·
[Rate Limit](../quality/RateLimit.md) (Q) ·
[Pricing Plan](../quality/PricingPlan.md) (Q)

### Interoperabilität

[Public API](../foundation/PublicAPI.md) (F) ·
[Community API](../foundation/CommunityAPI.md) (F) ·
[API Description](../foundation/APIDescription.md) (F) ·
[Atomic Parameter](../structure/AtomicParameter.md) (S) ·
[Atomic Parameter List](../structure/AtomicParameterList.md) (S) ·
[Parameter Tree](../structure/ParameterTree.md) (S) ·
[Parameter Forest](../structure/ParameterForest.md) (S) ·
[Data Element](../structure/DataElement.md) (S) ·
[Metadata Element](../structure/MetadataElement.md) (S) ·
[Id Element](../structure/IdElement.md) (S) ·
[Link Element](../structure/LinkElement.md) (S) ·
[Context Representation](../structure/ContextRepresentation.md) (S) ·
[Error Report](../structure/ErrorReport.md) (S) ·
[API Key](../structure/APIKey.md) (S) ·
[Service Level Agreement](../quality/ServiceLevelAgreement.md) (Q) ·
[Semantic Versioning](../evolution/SemanticVersioning.md) (E)

### Wartbarkeit

[Parameter Tree](../structure/ParameterTree.md) (S) ·
[Parameter Forest](../structure/ParameterForest.md) (S) ·
[Data Element](../structure/DataElement.md) (S) ·
[Context Representation](../structure/ContextRepresentation.md) (S) ·
[Version Identifier](../evolution/VersionIdentifier.md) (E) ·
[Semantic Versioning](../evolution/SemanticVersioning.md) (E) ·
[Two in Production](../evolution/TwoInProduction.md) (E) ·
[Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) (E) ·
[Experimental Preview](../evolution/ExperimentalPreview.md) (E) ·
[Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) (E) ·
[Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) (E)

### Sicherheit

[Frontend Integration](../foundation/FrontendIntegration.md) (F) ·
[Public API](../foundation/PublicAPI.md) (F) ·
[Community API](../foundation/CommunityAPI.md) (F) ·
[Solution-Internal API](../foundation/SolutionInternalAPI.md) (F) ·
[API Key](../structure/APIKey.md) (S) ·
[Data Element](../structure/DataElement.md) (S) ·
[Id Element](../structure/IdElement.md) (S) ·
[Context Representation](../structure/ContextRepresentation.md) (S) ·
[Error Report](../structure/ErrorReport.md) (S) ·
[Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) (E)

### Zuverlässigkeit

[Backend Integration](../foundation/BackendIntegration.md) (F) ·
[Public API](../foundation/PublicAPI.md) (F) ·
[API Key](../structure/APIKey.md) (S) ·
[Error Report](../structure/ErrorReport.md) (S) ·
[Request Bundle](../quality/RequestBundle.md) (Q) ·
[Rate Limit](../quality/RateLimit.md) (Q) ·
[Pricing Plan](../quality/PricingPlan.md) (Q) ·
[Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) (E)

### Konsistenz

[Backend Integration](../foundation/BackendIntegration.md) (F) ·
[Information Holder Resource](../responsibility/InformationHolderResource.md) (R) ·
[Operational Data Holder](../responsibility/OperationalDataHolder.md) (R) ·
[Master Data Holder](../responsibility/MasterDataHolder.md) (R)

### Flexibilität

[Frontend Integration](../foundation/FrontendIntegration.md) (F) ·
[Public API](../foundation/PublicAPI.md) (F) ·
[Solution-Internal API](../foundation/SolutionInternalAPI.md) (F) ·
[Data Transfer Resource](../responsibility/DataTransferResource.md) (R) ·
[Link Lookup Resource](../responsibility/LinkLookupResource.md) (R) ·
[API Key](../structure/APIKey.md) (S) ·
[Linked Information Holder](../quality/LinkedInformationHolder.md) (Q) ·
[Two in Production](../evolution/TwoInProduction.md) (E) ·
[Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) (E) ·
[Experimental Preview](../evolution/ExperimentalPreview.md) (E)

### Usability

[Frontend Integration](../foundation/FrontendIntegration.md) (F) ·
[Processing Resource](../responsibility/ProcessingResource.md) (R) ·
[Embedded Entity](../quality/EmbeddedEntity.md) (Q) ·
[Linked Information Holder](../quality/LinkedInformationHolder.md) (Q) ·
[Wish List](../quality/WishList.md) (Q) ·
[Wish Template](../quality/WishTemplate.md) (Q) ·
[Experimental Preview](../evolution/ExperimentalPreview.md) (E)

### Kosten / Affordability

[Public API](../foundation/PublicAPI.md) (F) ·
[Community API](../foundation/CommunityAPI.md) (F) ·
[Solution-Internal API](../foundation/SolutionInternalAPI.md) (F) ·
[API Key](../structure/APIKey.md) (S) ·
[Pricing Plan](../quality/PricingPlan.md) (Q) ·
[Rate Limit](../quality/RateLimit.md) (Q) ·
[Version Identifier](../evolution/VersionIdentifier.md) (E) ·
[Semantic Versioning](../evolution/SemanticVersioning.md) (E) ·
[Two in Production](../evolution/TwoInProduction.md) (E) ·
[Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) (E) ·
[Experimental Preview](../evolution/ExperimentalPreview.md) (E) ·
[Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) (E) ·
[Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) (E)

### Portabilität

[API Key](../structure/APIKey.md) (S) ·
[Error Report](../structure/ErrorReport.md) (S)

Die kürzeste Liste — beide Patterns lösen ihr Problem bewusst *unabhängig* von der
Protokolltechnologie und sind deshalb die einzigen, die Portabilität direkt adressieren.

## Verwandte Wiki-Seiten

- [Patterns nach Scope](navigation-byscope.md) · [nach Phase](navigation-byphase.md) · [nach Rolle](navigation-byrole.md)
- [Cheat Sheet](cheatsheet.md)
- [Kategorie Quality](category-quality.md)

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/byquality)
