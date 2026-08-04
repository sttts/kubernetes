---
title: Cheat Sheet
kategorie: Meta
quelle: https://microservice-api-patterns.org/cheatsheet
---

# Cheat Sheet

## Worum es geht

Das Cheat Sheet ist der problemgetriebene Einstieg in MAP: Es ordnet wiederkehrenden *Issues*
die Patterns zu, die sie adressieren, und ist chronologisch entlang eines API-Lebenszyklus
geordnet — von der Konzeption über Release und Produktivierung bis zu kontinuierlicher
Verbesserung und Wartung. Vorbild ist das Cheat Sheet im Anhang von Fowlers
*Patterns of Enterprise Application Architecture*; im Buch existiert eine erweiterte Fassung als
Anhang A, die zusätzlich die Beziehung zu ADDR, RDD und DDD diskutiert.

Die Quellseite warnt selbst vor der Verkürzung: Die Tabellen sind „a gross simplification of a set
of complex design considerations" — Forces und Consequences stehen im Patterntext, nicht hier.

## Inhalt

### Alle Patterns auf einen Blick

Die Website listet derzeit 45 Patterns (6 + 11 + 11 + 10 + 7); das Buch von 2022 spricht von 44.

#### Foundation

| Pattern | Einzeiler |
|---|---|
| [Frontend Integration](../foundation/FrontendIntegration.md) | Backend-Dienste einem oder mehreren Anwendungs-Frontends per Remote-API bereitstellen. |
| [Backend Integration](../foundation/BackendIntegration.md) | Unabhängig gebaute und deployte Backends Daten austauschen und Aktivitäten auslösen lassen. |
| [Public API](../foundation/PublicAPI.md) | API im öffentlichen Internet für unbegrenzt viele, unbekannte Clients anbieten. |
| [Community API](../foundation/CommunityAPI.md) | Zugriff auf eine geschlossene Nutzergruppe über mehrere Organisationen hinweg beschränken. |
| [Solution-Internal API](../foundation/SolutionInternalAPI.md) | API nur systemintern anbieten, etwa zwischen Komponenten derselben Anwendung. |
| [API Description](../foundation/APIDescription.md) | Struktur, Fehler, Ablauf, Vor-/Nachbedingungen und Policies explizit dokumentieren. |

#### Responsibility

| Pattern | Einzeiler |
|---|---|
| [Processing Resource](../responsibility/ProcessingResource.md) | Aktivitätsorientierter Endpunkt, der Kommandos auf Anwendungsebene bündelt. |
| [Information Holder Resource](../responsibility/InformationHolderResource.md) | Datenorientierter Endpunkt, der eine Entität exponiert und ihre Integrität schützt. |
| [Operational Data Holder](../responsibility/OperationalDataHolder.md) | Kurzlebige, häufig geänderte Transaktionsdaten mit schnellem CRUD. |
| [Master Data Holder](../responsibility/MasterDataHolder.md) | Langlebige, selten geänderte, vielfach referenzierte Stammdaten. |
| [Reference Data Holder](../responsibility/ReferenceDataHolder.md) | Für Clients unveränderliche Nachschlagedaten, nur lesbar. |
| [Link Lookup Resource](../responsibility/LinkLookupResource.md) | Verzeichnis, das die aktuellen Adressen anderer Endpunkte liefert. |
| [Data Transfer Resource](../responsibility/DataTransferResource.md) | Geteilter Ablageraum, der Sender und Empfänger zeitlich entkoppelt. |
| [State Creation Operation](../responsibility/StateCreationOperation.md) | Schreiben ohne Lesen des Vorzustands: „etwas ist passiert". |
| [Retrieval Operation](../responsibility/RetrievalOperation.md) | Rein lesend, mit Such-, Filter- und Formatierungsparametern. |
| [State Transition Operation](../responsibility/StateTransitionOperation.md) | Eingabe plus aktueller Zustand ergeben einen geprüften Zustandsübergang. |
| [Computation Function](../responsibility/ComputationFunction.md) | Seiteneffektfreie Berechnung ohne Zustandsbezug. |

#### Structure

| Pattern | Einzeiler |
|---|---|
| [Atomic Parameter](../structure/AtomicParameter.md) | Ein einzelnes unstrukturiertes Datum mit Name, Typ, Kardinalität. |
| [Atomic Parameter List](../structure/AtomicParameterList.md) | Mehrere zusammengehörige Atomic Parameters als ein kohäsives Element. |
| [Parameter Tree](../structure/ParameterTree.md) | Hierarchische Nutzlast mit einer Wurzel und geschachtelten Kindknoten. |
| [Parameter Forest](../structure/ParameterForest.md) | Mehrere Parameter Trees nebeneinander als Nutzlast. |
| [Data Element](../structure/DataElement.md) | Eigenes Fachvokabular in Nachrichten statt durchgereichter Interndaten. |
| [Metadata Element](../structure/MetadataElement.md) | Zusatzinformation zur korrekten Interpretation der übrigen Elemente. |
| [Id Element](../structure/IdElement.md) | Eindeutiger Bezeichner für Endpunkte, Operationen, Repräsentationselemente. |
| [Link Element](../structure/LinkElement.md) | Netzwerkauflösbarer Zeiger auf einen anderen Endpunkt (HATEOAS). |
| [API Key](../structure/APIKey.md) | Je Client eindeutiges Token zur Identifikation der Aufrufe. |
| [Error Report](../structure/ErrorReport.md) | Maschinenlesbare Fehlercodes plus menschenlesbare Beschreibung in der Antwort. |
| [Context Representation](../structure/ContextRepresentation.md) | Kontext-Metadaten gebündelt in der Nutzlast statt in Protokoll-Headern. |

#### Quality

| Pattern | Einzeiler |
|---|---|
| [Embedded Entity](../quality/EmbeddedEntity.md) | Zieldaten einer Beziehung direkt einbetten, um Folgeaufrufe zu sparen. |
| [Linked Information Holder](../quality/LinkedInformationHolder.md) | Statt der Daten einen Link auf einen eigenen Endpunkt mitgeben. |
| [Pagination](../quality/Pagination.md) | Große Ergebnismengen in Chunks ausliefern, mit Navigationsmetadaten. |
| [Wish List](../quality/WishList.md) | Client zählt im Request die gewünschten Felder auf („Response Shaping"). |
| [Wish Template](../quality/WishTemplate.md) | Request spiegelt die Antworthierarchie, um geschachtelte Wünsche auszudrücken. |
| [Conditional Request](../quality/ConditionalRequest.md) | Verarbeitung nur, wenn die per Metadaten ausgedrückte Bedingung erfüllt ist. |
| [Request Bundle](../quality/RequestBundle.md) | Mehrere unabhängige Requests in einer Nachricht bündeln. |
| [Rate Limit](../quality/RateLimit.md) | Nutzungsobergrenze gegen übermäßige Nutzung durchsetzen. |
| [Pricing Plan](../quality/PricingPlan.md) | Nutzung messen und abrechnen, mit Metriken je Operation. |
| [Service Level Agreement](../quality/ServiceLevelAgreement.md) | Testbare Service-Level-Objectives samt Konsequenzen vereinbaren. |

#### Evolution

| Pattern | Einzeiler |
|---|---|
| [Version Identifier](../evolution/VersionIdentifier.md) | Expliziter Versionsindikator in Beschreibung und Nachricht. |
| [Semantic Versioning](../evolution/SemanticVersioning.md) | Dreistelliges `x.y.z`-Schema macht Kompatibilität am Bezeichner ablesbar. |
| [Two in Production](../evolution/TwoInProduction.md) | Zwei Versionen parallel betreiben, rollierend aktualisieren und abkündigen. |
| [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) | Zusage, eine Version für einen festen Zeitraum nicht zu brechen. |
| [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) | Zusage, eine Version nie zu brechen oder abzuschalten. |
| [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) | Frühe Abkündigung mit Stichtag, danach Entfernung. |
| [Experimental Preview](../evolution/ExperimentalPreview.md) | Zugang ohne Zusagen zu Umfang, Stabilität und Langlebigkeit. |

### Issue-zu-Pattern-Tabellen

**Einstieg ins API-Design**

| Issue | Patterns |
|---|---|
| Ich will eine API anbieten | [Frontend Integration](../foundation/FrontendIntegration.md) oder [Backend Integration](../foundation/BackendIntegration.md) |
| Die API soll breit zugänglich sein | [Public API](../foundation/PublicAPI.md) |
| Die Sichtbarkeit soll eingeschränkt sein | [Community API](../foundation/CommunityAPI.md) oder [Solution-Internal API](../foundation/SolutionInternalAPI.md) |
| Endpunktkandidaten identifizieren | Ein Endpunkt je DDD-*Bounded Context* (hohe Kohäsion) oder je *Aggregate* (feingranular) |
| Aktivitätsorientierte Fähigkeit modellieren | [Processing Resource](../responsibility/ProcessingResource.md) |
| Datenorientierte Fähigkeit modellieren | [Information Holder Resource](../responsibility/InformationHolderResource.md) |
| Langlebige Daten anbieten | [Master Data Holder](../responsibility/MasterDataHolder.md) bzw. [Reference Data Holder](../responsibility/ReferenceDataHolder.md) |
| Kurzlebige Transaktionsdaten anbieten | [Operational Data Holder](../responsibility/OperationalDataHolder.md) |
| Clients transiente Daten austauschen lassen | [Data Transfer Resource](../responsibility/DataTransferResource.md) |
| Endpunktverzeichnis anbieten | [Link Lookup Resource](../responsibility/LinkLookupResource.md) |
| Client initialisiert providerseitigen Zustand | [State Creation Operation](../responsibility/StateCreationOperation.md) |
| Client aktualisiert providerseitigen Zustand | [State Transition Operation](../responsibility/StateTransitionOperation.md) |
| Client liest providerseitigen Zustand | [Retrieval Operation](../responsibility/RetrievalOperation.md) |
| Client ruft zustandslose Operation | [Computation Function](../responsibility/ComputationFunction.md) |
| Datenvertrag festlegen, einfache Daten | [Atomic Parameter](../structure/AtomicParameter.md), [Atomic Parameter List](../structure/AtomicParameterList.md) |
| Datenvertrag festlegen, komplexe Daten | [Parameter Tree](../structure/ParameterTree.md), [Parameter Forest](../structure/ParameterForest.md) |
| Strukturierte Nutzlasten austauschen | [Data Element](../structure/DataElement.md) mit [Embedded Entity](../quality/EmbeddedEntity.md) |
| Repräsentationselemente unterscheidbar machen | [Id Element](../structure/IdElement.md) |
| Unflexibler, statischer Ablauf | Vom [Id Element](../structure/IdElement.md) zum [Link Element](../structure/LinkElement.md) (HATEOAS) |

**Release und Produktivierung**

| Issue | Patterns |
|---|---|
| Clients müssen wissen, wie sie rufen | [API Description](../foundation/APIDescription.md) veröffentlichen |
| Bedingungen dokumentieren | [Service Level Agreement](../quality/ServiceLevelAgreement.md) |
| API robust und aussagefähig machen | [Error Report](../structure/ErrorReport.md) in die Antworten |
| Faire Nutzung sicherstellen | [Rate Limit](../quality/RateLimit.md) |
| Mit der API Geld verdienen | [Pricing Plan](../quality/PricingPlan.md) |

**Kontinuierliche Verbesserung**

| Issue | Patterns |
|---|---|
| Clients melden Interoperabilitäts-/Usability-Probleme | Von minimaler zu voller [API Description](../foundation/APIDescription.md); [Context Representation](../structure/ContextRepresentation.md) einführen |
| Clients melden Performance-Probleme | [Embedded Entity](../quality/EmbeddedEntity.md) → [Linked Information Holder](../quality/LinkedInformationHolder.md); [Wish List](../quality/WishList.md)/[Wish Template](../quality/WishTemplate.md); [Conditional Request](../quality/ConditionalRequest.md), [Request Bundle](../quality/RequestBundle.md), [Pagination](../quality/Pagination.md) |
| Zugriffskontrolle nötig | [API Key](../structure/APIKey.md) oder vollwertige IAM-Lösung (OAuth) |

**Support und Wartung**

| Issue | Patterns |
|---|---|
| Nicht abwärtskompatible Änderung | Neue Operation oder [Version Identifier](../evolution/VersionIdentifier.md) |
| Bedeutung von Änderungen kommunizieren | [Semantic Versioning](../evolution/SemanticVersioning.md) |
| Mehrere Versionen betreiben | [Two in Production](../evolution/TwoInProduction.md) |
| Mehrere Versionen vermeiden | [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) |
| Befristete Verfügbarkeit zusagen | [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) |
| Unbefristete Verfügbarkeit zusagen | [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) |
| Keine Zusagen machen | [Experimental Preview](../evolution/ExperimentalPreview.md) |

## Verwandte Wiki-Seiten

- [Überblick MAP](overview.md) · [Terminologie](terms.md)
- Kategorien: [Foundation](category-foundation.md) · [Responsibility](category-responsibility.md) · [Structure](category-structure.md) · [Quality](category-quality.md) · [Evolution](category-evolution.md)
- Filter: [Scope](navigation-byscope.md) · [Phase](navigation-byphase.md) · [Rolle](navigation-byrole.md) · [Qualität](navigation-byquality.md)

## Bezug zu Kubernetes / KRM

Wer eine Kubernetes-API entwirft, kann große Teile der ersten Tabelle überspringen — die
Antworten sind vom Modell vorgegeben. Der Kurzdurchlauf für einen CRD-Entwurf sieht so aus:

| Cheat-Sheet-Issue | KRM-Antwort |
|---|---|
| API anbieten | CRD oder Aggregated API Server; der Integrationstyp ist immer beides zugleich |
| Sichtbarkeit | Nicht Teil der API, sondern von RBAC und Netzwerk |
| Endpunktkandidaten | Eine Ressource je Aggregate; die DDD-Nähe des Cheat Sheets passt gut |
| Endpunktrolle | Fast immer Information Holder Resource mit `spec`/`status`; Processing Resource nur als Subresource |
| Operationsverantwortung | Nicht wählbar — die Verbmenge ist fix |
| Datenvertrag | Immer [Parameter Tree](../structure/ParameterTree.md), Schema als OpenAPI v3 in der CRD |
| Elemente unterscheidbar machen | `metadata.name`, `metadata.namespace`, `metadata.uid` |
| Statischer Ablauf, HATEOAS | Nicht vorgesehen; stattdessen `ObjectReference`/`ownerReferences` plus Discovery |
| Fehlerbericht | `metav1.Status` mit `reason`/`code`/`details` |
| Zugriffskontrolle | ServiceAccount-Token, Client-Zertifikate, RBAC |
| Rate Limit | API Priority and Fairness (`FlowSchema`, `PriorityLevelConfiguration`) |
| Versionierung | Group-Version im Pfad; Stabilitätsstufe im Namen (`v1alpha1`/`v1beta1`/`v1`) |
| Mehrere Versionen | Alle bedienten Versionen parallel, verlustfrei konvertierend |
| Abkündigung | Deprecation Policy, `Warning`-Header, Entfernung zum angekündigten Release |

Die verbleibende echte Designarbeit liegt fast vollständig in den Zeilen zu Datenvertrag und
Nachrichtengröße — also in [Structure](category-structure.md) und
[Quality](category-quality.md). Was das Cheat Sheet als *Continuous Improvement* führt, ist in
KRM zudem schwieriger: Ein Wechsel von [Embedded Entity](../quality/EmbeddedEntity.md) zu
[Linked Information Holder](../quality/LinkedInformationHolder.md) ist bei einer GA-API eine
brechende Schemaänderung und damit nur über eine neue Group-Version möglich.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/cheatsheet)
