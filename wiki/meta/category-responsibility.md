---
title: Kategorie Responsibility
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/responsibility
---

# Kategorie: Responsibility Patterns

## Worum es geht

Die *Responsibility Patterns* beantworten zwei Fragen:

- Welche architektonische Rolle spielt ein API-Endpunkt?
- Welche Verantwortung trägt eine einzelne Operation — und wie wirken sich Rolle und Verantwortung
  auf Service-Schnitt und Granularität aus?

Die Begriffe stammen aus Responsibility-Driven Design (RDD): Eine *Rolle* ist ein Satz
zusammengehöriger Verantwortungen (hier: der Endpunkt), eine *Verantwortung* eine Verpflichtung,
eine Aufgabe auszuführen oder eine Information zu kennen (hier: die Operation).

Die Kategorie ist dreigeteilt: zwei generische Endpunkt-Rollen, fünf Spezialisierungen der
datenorientierten Rolle, und vier Operationsverantwortungen. Die Operationsverantwortungen sind
über ihre Zustandssignatur definiert und damit trennscharf — `S` steht für den providerseitigen
Zustand:

| Verantwortung | Signatur | Liest Zustand | Ändert Zustand |
|---|---|---|---|
| [Computation Function](../responsibility/ComputationFunction.md) | `cf: in -> out` | nein | nein |
| [Retrieval Operation](../responsibility/RetrievalOperation.md) | `ro: (in,S) -> out` | ja | nein |
| [State Creation Operation](../responsibility/StateCreationOperation.md) | `sco: in -> (out,S')` | nein | ja |
| [State Transition Operation](../responsibility/StateTransitionOperation.md) | `sto: (in,S) -> (out,S')` | ja | ja |

## Inhalt

### Endpunkt-Rollen

| Pattern | Einzeiler |
|---|---|
| [Processing Resource](../responsibility/ProcessingResource.md) | Aktivitätsorientierter Endpunkt: bündelt und kapselt Kommandos bzw. Aktivitäten auf Anwendungsebene, die der Client auslösen will. |
| [Information Holder Resource](../responsibility/InformationHolderResource.md) | Datenorientierter Endpunkt: exponiert eine Datenentität mit CRUD- und Suchoperationen, hält aber die Implementierung verborgen und schützt die Datenintegrität. |

Die Wahl zwischen beiden ist die grundlegende Weiche: *Processing Resource* koppelt lose (der
Client kennt nur ein Kommando), *Information Holder Resource* koppelt stärker (der Client kennt
das Datenmodell), ist dafür aber flexibler nutzbar.

### Typen von Information Holder Resources

| Pattern | Einzeiler |
|---|---|
| [Operational Data Holder](../responsibility/OperationalDataHolder.md) | Kurzlebige, sich häufig ändernde Transaktionsdaten mit vielen ausgehenden Beziehungen; volles, schnelles CRUD. |
| [Master Data Holder](../responsibility/MasterDataHolder.md) | Langlebige, selten geänderte, von vielen referenzierte Stammdaten; Löschen wird als Sonderform des Änderns behandelt. |
| [Reference Data Holder](../responsibility/ReferenceDataHolder.md) | Für Clients unveränderliche, langlebige, vielfach referenzierte Nachschlagedaten; nur Leseoperationen. |
| [Link Lookup Resource](../responsibility/LinkLookupResource.md) | Verzeichnis-Endpunkt, der [Link Elements](../structure/LinkElement.md) auf die aktuellen Adressen anderer Endpunkte liefert, damit Nachrichten keine Adressen fest verdrahten müssen. |
| [Data Transfer Resource](../responsibility/DataTransferResource.md) | Gemeinsamer Ablageraum unter einer global eindeutigen Adresse; entkoppelt Sender und Empfänger zeitlich und identitätsbezogen. |

Die ersten drei bilden eine Achse über Lebensdauer und Mutabilität (kurzlebig/veränderlich →
langlebig/selten veränderlich → langlebig/unveränderlich). Die letzten beiden sind funktionale
Spezialfälle, keine Punkte auf dieser Achse.

### Operationsverantwortungen

| Pattern | Einzeiler |
|---|---|
| [State Creation Operation](../responsibility/StateCreationOperation.md) | Schreibende Operation ohne Lesen des Vorzustands: Der Client meldet, dass etwas passiert ist, was der Provider wissen muss. |
| [Retrieval Operation](../responsibility/RetrievalOperation.md) | Rein lesende Operation mit Such-, Filter- und Formatierungsmöglichkeiten in der Signatur. |
| [State Transition Operation](../responsibility/StateTransitionOperation.md) | Kombiniert Clienteingabe und aktuellen Zustand zu einem Zustandsübergang; gültige Übergänge werden im Endpunkt modelliert und zur Laufzeit geprüft. |
| [Computation Function](../responsibility/ComputationFunction.md) | Seiteneffektfreie Berechnung: gleiche Eingabe, gleiches Ergebnis, kein Zustand. |

## Verwandte Wiki-Seiten

- [Kategorie Foundation](category-foundation.md) — was vorher entschieden wird
- [Kategorie Structure](category-structure.md) — was nachher entschieden wird
- [Cheat Sheet](cheatsheet.md) — Issue-zu-Pattern-Zuordnung für Endpunkt- und Operationsdesign
- [Tutorials](tutorials.md) — Tutorial 2, Schritt 2, wendet die Kategorie auf einen Fall an
- [Patterns nach Granularität und Rolle](navigation-byquality.md)

## Bezug zu Kubernetes / KRM

KRM ist nahezu vollständig datenorientiert: Fast jede Ressource — `Pod`, `Deployment`,
`ConfigMap`, `Node` — ist eine *Information Holder Resource*, adressiert über eine
Group-Version-Resource und die Verben `create`, `get`, `list`, `watch`, `update`, `patch`,
`delete`, `deletecollection`. Die MAP-Typenachse lässt sich sauber wiederfinden: `Event`
und `Lease` verhalten sich wie *Operational Data Holder* (kurzlebig, TTL-behaftet, hohe
Änderungsrate), `Namespace`, `CustomResourceDefinition` oder `StorageClass` wie *Master Data
Holder* (langlebig, vielfach referenziert), und die Discovery-Antworten unter `/apis` sowie
`APIService`-Objekte wie eine *Link Lookup Resource*. Ein echter *Reference Data Holder* ist die
OpenAPI-Auslieferung unter `/openapi/v3`: nur lesbar, für Clients unveränderlich.

*Processing Resources* gibt es in KRM nur als bewusst begrenzte Ausnahme, und zwar als
Subresources: `pods/exec`, `pods/portforward`, `pods/eviction`, `deployments/scale`,
`serviceaccounts/token`. Auch die Review-APIs (`SubjectAccessReview`, `TokenReview`,
`SelfSubjectRulesReview`) gehören hierher — sie sind *Computation Functions* im MAP-Sinn: Man
`create`t ein Objekt, bekommt es mit gefülltem `status` zurück, und es wird nichts gespeichert.
Das ist eine KRM-Eigenwilligkeit: Weil das uniforme Schema nur Ressourcen und Verben kennt, wird
ein RPC-artiger Aufruf als „Erzeugen eines Objekts, das nie persistiert wird" modelliert.

Am stärksten weicht KRM bei den *State Transition Operations* ab. Das deklarative, level-getriggerte
Modell verbietet dem Client, einen Zustandsübergang zu *befehlen*: Er schreibt nur den Wunschzustand
nach `spec`, und ein Controller konvergiert dorthin und berichtet über `status`. Die
`update`/`patch`-Aufrufe sind damit formal *State Transition Operations* (`(in,S) -> (out,S')`,
inklusive Prüfung über Admission und `metadata.resourceVersion`-Precondition), semantisch aber
keine Geschäftsübergänge — die Übergangslogik liegt im Controller, nicht im Endpunkt. Wo MAP
explizite Übergangsoperationen empfiehlt, hat KRM die Statusmaschine aus der API herausgezogen.
Bewertung: Endpunktrollen und Retrieval Operations **genauso**, Processing Resource und
State Transition Operation **anders**, weil das uniforme Schema und die Reconciliation-Schleife
die Verantwortung verschieben.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility)
