---
title: Patterns nach Scope
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/byscope
---

# Patterns nach Scope

## Worum es geht

Der *Scope*-Filter beantwortet die Frage: „Woran arbeite ich gerade?" MAP unterscheidet fünf
Artefaktebenen — die API als Ganzes, den einzelnen Endpunkt, die einzelne Operation, die
Nachrichtenrepräsentation und die API-Dokumentation. Ein Pattern kann in mehreren Scopes
auftauchen; die Evolution-Patterns etwa gelten sowohl für die gesamte API als auch für einzelne
Endpunkte, weil beide Granularitäten sinnvoll sind.

Die Ebenen entsprechen dem API-Domänenmodell (siehe [Terminologie](terms.md)) und bilden eine
Enthaltenseinshierarchie: API ⊃ Endpunkt ⊃ Operation ⊃ Nachricht ⊃ Repräsentationselement.

## Inhalt

### API als Ganzes

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
| [Version Identifier](../evolution/VersionIdentifier.md) | Evolution |
| [Semantic Versioning](../evolution/SemanticVersioning.md) | Evolution |
| [Two in Production](../evolution/TwoInProduction.md) | Evolution |
| [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) | Evolution |
| [Experimental Preview](../evolution/ExperimentalPreview.md) | Evolution |
| [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) | Evolution |
| [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) | Evolution |

### API-Endpunkt

| Pattern | Kategorie |
|---|---|
| [Processing Resource](../responsibility/ProcessingResource.md) | Responsibility |
| [Information Holder Resource](../responsibility/InformationHolderResource.md) | Responsibility |
| [Operational Data Holder](../responsibility/OperationalDataHolder.md) | Responsibility |
| [Master Data Holder](../responsibility/MasterDataHolder.md) | Responsibility |
| [Reference Data Holder](../responsibility/ReferenceDataHolder.md) | Responsibility |
| [Link Lookup Resource](../responsibility/LinkLookupResource.md) | Responsibility |
| [Data Transfer Resource](../responsibility/DataTransferResource.md) | Responsibility |
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

### Operation

| Pattern | Kategorie |
|---|---|
| [State Creation Operation](../responsibility/StateCreationOperation.md) | Responsibility |
| [State Transition Operation](../responsibility/StateTransitionOperation.md) | Responsibility |
| [Retrieval Operation](../responsibility/RetrievalOperation.md) | Responsibility |
| [Computation Function](../responsibility/ComputationFunction.md) | Responsibility |
| [API Key](../structure/APIKey.md) | Structure |
| [Pagination](../quality/Pagination.md) | Quality |
| [Wish List](../quality/WishList.md) | Quality |
| [Wish Template](../quality/WishTemplate.md) | Quality |
| [Conditional Request](../quality/ConditionalRequest.md) | Quality |
| [Request Bundle](../quality/RequestBundle.md) | Quality |
| [Rate Limit](../quality/RateLimit.md) | Quality |

### Nachrichtenrepräsentation

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
| [Error Report](../structure/ErrorReport.md) | Structure |
| [Context Representation](../structure/ContextRepresentation.md) | Structure |
| [Embedded Entity](../quality/EmbeddedEntity.md) | Quality |
| [Linked Information Holder](../quality/LinkedInformationHolder.md) | Quality |

### API-Dokumentation

| Pattern | Kategorie |
|---|---|
| [API Description](../foundation/APIDescription.md) | Foundation |
| [API Key](../structure/APIKey.md) | Structure |
| [Service Level Agreement](../quality/ServiceLevelAgreement.md) | Quality |

## Verwandte Wiki-Seiten

- [Patterns nach Phase](navigation-byphase.md) · [nach Rolle](navigation-byrole.md) · [nach Qualität](navigation-byquality.md)
- [Cheat Sheet](cheatsheet.md) — Einstieg über konkrete Probleme statt über Artefaktebenen
- [Terminologie](terms.md) — Definition von Endpoint, Operation, Message, Representation

## Bezug zu Kubernetes / KRM

Die Scope-Hierarchie lässt sich fast eins zu eins auf KRM abbilden, allerdings mit einer
Verschiebung, die viele Missverständnisse erklärt:

| MAP-Scope | KRM-Entsprechung |
|---|---|
| API | Cluster-Endpunkt bzw. eine API-Gruppe (`apps`, `batch`, `rbac.authorization.k8s.io`) |
| API-Endpunkt | Group-Version-Resource, z. B. `apps/v1, Resource=deployments`; Subresources (`/status`, `/scale`) sind eigene Endpunkte |
| Operation | Verb: `get`, `list`, `watch`, `create`, `update`, `patch`, `delete`, `deletecollection` |
| Nachrichtenrepräsentation | Das Objekt selbst (`apiVersion`/`kind`/`metadata`/`spec`/`status`) bzw. `metav1.Status` bei Fehlern |
| API-Dokumentation | Discovery (`/apis`), OpenAPI v2/v3, `kubectl explain` |

Der wesentliche Unterschied: In MAP definiert der Designer je Endpunkt eine eigene Menge von
Operationen mit eigenen Signaturen. In KRM ist die Operationsmenge *fix* und über alle Ressourcen
identisch; der Designer wählt nur noch aus, welche Verben eine Ressource unterstützt (bei CRDs
implizit über `spec.versions[].served` und die vorhandenen Subresources). Die gestalterische
Freiheit wandert damit vollständig vom Scope „Operation" in den Scope „Nachrichtenrepräsentation"
— das Schema ist bei KRM praktisch der gesamte API-Entwurf. Wer eine Kubernetes-API entwirft,
verbringt seine Zeit fast ausschließlich in der Kategorie
[Structure](category-structure.md), nicht in [Responsibility](category-responsibility.md).

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/byscope)
