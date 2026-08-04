# MAP-Wiki — Microservice API Patterns, gelesen mit Kubernetes-Brille

Zusammenfassung der Pattern-Sprache **Microservice API Patterns (MAP)** von
Zimmermann, Stocker, Lübke, Pautasso und Zdun — Quelle
<https://microservice-api-patterns.org/> plus die dort verlinkten Papers.

Jede Pattern-Seite hat einen Abschnitt **„Bezug zu Kubernetes / KRM"**, der einordnet, ob das
Kubernetes Resource Model das Pattern genauso, anders oder gar nicht umsetzt — und warum.
Die durchgehende Analyse steht in **[KRM vs. MAP](KRM-vs-MAP.md)**.

## Einstieg

| Seite | Wofür |
|---|---|
| [Überblick](meta/overview.md) | Was MAP ist, die fünf Kategorien, POINT-Prinzipien, Autoren |
| [Introduction](meta/introduction.md) · [Primer](meta/primer.md) | Motivation und Hintergrund |
| [Cheat Sheet](meta/cheatsheet.md) | Alle Patterns auf einer Seite |
| [Terminologie](meta/terms.md) | Endpoint, Operation, Message, Representation — mit KRM-Pendants |
| [Tutorials](meta/tutorials.md) | Geführtes Durchspielen an einem Fallbeispiel |
| [Verwandte Pattern-Sprachen](meta/related-pattern-languages.md) | EIP, PoEAA, Refactoring-Katalog |
| [Buch und Ressourcen](meta/book-and-resources.md) | Bibliographie, Slides, MDSL |

**Sichten:** [nach Scope](meta/navigation-byscope.md) ·
[nach Qualität](meta/navigation-byquality.md) ·
[nach Phase](meta/navigation-byphase.md) ·
[nach Rolle](meta/navigation-byrole.md)

## Foundation — [Kategorieseite](meta/category-foundation.md)

Wo liegt die API-Grenze, wer darf sie sehen, wie wird sie beschrieben?

| Pattern | Kurzform |
|---|---|
| [Frontend Integration](foundation/FrontendIntegration.md) | Backend bedient eine UI über eine Remote-API |
| [Backend Integration](foundation/BackendIntegration.md) | Getrennt deployte Systemteile integrieren nachrichtenbasiert |
| [Public API](foundation/PublicAPI.md) | Im offenen Internet exponiert, unbekannte Clients |
| [Community API](foundation/CommunityAPI.md) | Zugangsbeschränkt für einen bekannten Kreis von Organisationen |
| [Solution-Internal API](foundation/SolutionInternalAPI.md) | Nur innerhalb einer Anwendung sichtbar |
| [API Description](foundation/APIDescription.md) | Der explizite, dokumentierte Vertrag |

## Responsibility — [Kategorieseite](meta/category-responsibility.md)

Welche architektonische Rolle hat ein Endpunkt, welche Verantwortung eine Operation?

**Endpunkt-Rollen**

| Pattern | Kurzform |
|---|---|
| [Processing Resource](responsibility/ProcessingResource.md) | Endpunkt-Identität ist eine *Aktivität* |
| [Information Holder Resource](responsibility/InformationHolderResource.md) | Endpunkt-Identität ist eine *Datenentität* |

**Typen von Information Holdern**

| Pattern | Kurzform |
|---|---|
| [Operational Data Holder](responsibility/OperationalDataHolder.md) | Kurzlebig, häufig geändert |
| [Master Data Holder](responsibility/MasterDataHolder.md) | Langlebig, selten geändert, viel referenziert |
| [Reference Data Holder](responsibility/ReferenceDataHolder.md) | Statisch, nur vom Provider änderbar |
| [Data Transfer Resource](responsibility/DataTransferResource.md) | Geteilter Speicher als Blackboard zwischen Clients |
| [Link Lookup Resource](responsibility/LinkLookupResource.md) | Hält Adressen *anderer* Endpunkte |

**Operations-Verantwortlichkeiten**

| Pattern | Kurzform |
|---|---|
| [State Creation Operation](responsibility/StateCreationOperation.md) | Legt Provider-Zustand an |
| [State Transition Operation](responsibility/StateTransitionOperation.md) | Kombiniert Eingabe mit vorhandenem Zustand |
| [Retrieval Operation](responsibility/RetrievalOperation.md) | Liest und projiziert Provider-Zustand |
| [Computation Function](responsibility/ComputationFunction.md) | Rechnet nur aus der Eingabe, ohne Zustand |

## Structure — [Kategorieseite](meta/category-structure.md)

Wie sind die Payloads aufgebaut?

**Repräsentations-Elemente**

| Pattern | Kurzform |
|---|---|
| [Atomic Parameter](structure/AtomicParameter.md) | Ein einzelner unstrukturierter Wert |
| [Atomic Parameter List](structure/AtomicParameterList.md) | Flache Liste zusammengehöriger Werte |
| [Parameter Tree](structure/ParameterTree.md) | Hierarchie mit genau einer Wurzel |
| [Parameter Forest](structure/ParameterForest.md) | Mehrere gleichrangige Bäume nebeneinander |

**Element-Stereotypen**

| Pattern | Kurzform |
|---|---|
| [Data Element](structure/DataElement.md) | Fachliche Nutzdaten in eigenem Vokabular |
| [Metadata Element](structure/MetadataElement.md) | Erklärende Zusatzinformation über die Nachricht |
| [Id Element](structure/IdElement.md) | Macht Elemente unterscheidbar |
| [Link Element](structure/LinkElement.md) | Id, die zugleich adressierbar ist |

**Spezielle Repräsentationen**

| Pattern | Kurzform |
|---|---|
| [API Key](structure/APIKey.md) | Providerseitig vergebenes Client-Token |
| [Error Report](structure/ErrorReport.md) | Fehler als maschinenlesbare Repräsentation |
| [Context Representation](structure/ContextRepresentation.md) | Aufrufkontext protokollneutral gebündelt |

## Quality — [Kategorieseite](meta/category-quality.md)

Wie werden Datenmenge, Referenzen und Zusagen bewirtschaftet?

**Referenz-Management**

| Pattern | Kurzform |
|---|---|
| [Embedded Entity](quality/EmbeddedEntity.md) | Verwandte Daten mitliefern |
| [Linked Information Holder](quality/LinkedInformationHolder.md) | Verwandte Daten nur referenzieren |

**Datensparsamkeit**

| Pattern | Kurzform |
|---|---|
| [Pagination](quality/Pagination.md) | Große Ergebnismengen in Chunks |
| [Wish List](quality/WishList.md) | Client zählt gewünschte Felder flach auf |
| [Wish Template](quality/WishTemplate.md) | Client schickt ein Skelett der Wunschstruktur |
| [Conditional Request](quality/ConditionalRequest.md) | Antwort nur, wenn sich etwas geändert hat |
| [Request Bundle](quality/RequestBundle.md) | Mehrere Requests in einer Nachricht |

**Governance**

| Pattern | Kurzform |
|---|---|
| [Rate Limit](quality/RateLimit.md) | Aufrufe pro Client und Zeitfenster begrenzen |
| [Pricing Plan](quality/PricingPlan.md) | Nutzung messen und abrechnen |
| [Service Level Agreement](quality/ServiceLevelAgreement.md) | Messbare QoS-Zusagen neben der Description |

## Evolution — [Kategorieseite](meta/category-evolution.md)

Wie ändert sich eine API, ohne ihre Clients zu zerstören?

| Pattern | Kurzform |
|---|---|
| [Version Identifier](evolution/VersionIdentifier.md) | Expliziter Versionsindikator |
| [Semantic Versioning](evolution/SemanticVersioning.md) | Major/Minor/Patch mit Kompatibilitätssemantik |
| [Two in Production](evolution/TwoInProduction.md) | Mehrere Versionen laufen parallel |
| [Limited Lifetime Guarantee](evolution/LimitedLifetimeGuarantee.md) | Zusage über einen festen Zeitraum |
| [Eternal Lifetime Guarantee](evolution/EternalLifetimeGuarantee.md) | Zusage, nie zu brechen |
| [Experimental Preview](evolution/ExperimentalPreview.md) | Zugriff ohne Stabilitätszusage |
| [Aggressive Obsolescence](evolution/AggressiveObsolescence.md) | Früh ankündigen, hart abschalten |

## Papers

Die Website kürzt seit Erscheinen des Buchs die Abschnitte *Forces* und *Consequences* auf
Stichworte. Die Original-Papers enthalten sie vollständig — diese Seiten arbeiten genau das heraus.

| Paper | Venue | Deckt ab |
|---|---|---|
| [Introduction to MAP](papers/2020-introduction-to-map.md) | Microservices 2017/2019, OASIcs | Domänenmodell, Kategorien, Gesamtbild |
| [Interface Representation Patterns](papers/2017-interface-representation-patterns.md) | EuroPLoP 2017 | Structure-Kategorie |
| [Interface Quality Patterns](papers/2018-interface-quality-patterns.md) | EuroPLoP 2018 | Quality-Kategorie |
| [Guiding Architectural Decision Making](papers/2018-icsoc-decision-guidance.md) | ICSOC 2018 | Entscheidungsmodell hinter Quality |
| [Interface Evolution Patterns](papers/2019-interface-evolution-patterns.md) | EuroPLoP 2019 | Evolution-Kategorie |
| [Interface Responsibility Patterns](papers/2020-interface-responsibility-patterns.md) | EuroPLoP 2020 | Endpunkt-Rollen, Operationen |
| [Data-Oriented Responsibility Patterns](papers/2020-data-oriented-responsibility-patterns.md) | EuroPLoP 2020 | Information-Holder-Typen |

## Die Analyse

- **[KRM vs. MAP](KRM-vs-MAP.md)** — die vollständige Gegenüberstellung: was Kubernetes anders
  macht, wo es besser ist und was fehlt. Bilanz über alle 45 Patterns.
- **[One Schema to Rule Them All](blog/krm-through-the-lens-of-map.md)** (englisch) — dieselbe
  Analyse als Blogpost-Deep-Dive, mit Beispielen.

---

*Quelle: <https://microservice-api-patterns.org/>, Zimmermann/Stocker/Lübke/Pautasso/Zdun.
Diese Zusammenfassungen sind eigene Formulierungen, keine Kopien; die Pattern-Sprache selbst
gehört ihren Autoren.*
