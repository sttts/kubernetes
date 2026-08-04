---
title: Buch, Publikationen und Ressourcen
kategorie: Meta
quelle: https://microservice-api-patterns.org/book
---

# Buch, Publikationen und Ressourcen

## Worum es geht

Diese Seite bündelt die Quellen hinter MAP: das Buch von 2022, die EuroPLoP- und
Konferenzpapiere, aus denen die Patterns hervorgegangen sind, ergänzende Ressourcen (Vorträge,
Beispielimplementierungen, DSL) und die Autoren. Fasst die MAP-Seiten `/book`, `/publications`,
`/quickstart-resources`, `/slides` und `/about` zusammen.

## Inhalt

### Das Buch

**Patterns for API Design: Simplifying Integration with Loosely Coupled Message Exchanges.**
Olaf Zimmermann, Mirko Stocker, Daniel Lübke, Uwe Zdun, Cesare Pautasso.
Addison-Wesley Professional, Addison-Wesley Signature Series (Vernon), November 2022,
ISBN 9780137670109. Übersetzungen: Chinesisch (traditionell), Koreanisch, Portugiesisch.

Aufbau in drei Teilen:

| Teil | Inhalt |
|---|---|
| 1 — Foundations and Narratives | Kap. 1: API-Grundlagen und ein Domänenmodell für Remote-APIs. Kap. 2: Fallstudie *Lakeside Mutual*. Kap. 3: sechs Entscheidungsnarrative mit 29 wiederkehrenden Entscheidungen samt Optionen und Kriterien. |
| 2 — The Patterns | Kap. 4–9: die 44 Patterns. Kap. 4 Einführung und Foundations, Kap. 5 Endpunkttypen und Operationen, Kap. 6 Nachrichtenrepräsentationen, Kap. 7 Qualitätsverfeinerung, Kap. 8 Evolution, Kap. 9 Dokumentation. |
| 3 — Our Patterns in Action | Kap. 10: zwei reale Systeme (Schweizer Hypothekargeschäft, Angebots- und Bestellprozesse im Bauwesen). Kap. 11: Retrospektive und Ausblick. |

Anhänge: A — Endpoint Identification and Pattern Selection Guides, die erweiterte Fassung des
[Cheat Sheets](cheatsheet.md), inklusive Bezug zu ADDR, RDD und DDD. B — Implementierung der
Lakeside-Mutual-Fallstudie. C — Einführung in MDSL, die Microservice Domain Specific Language.

BibTeX-Schlüssel der Autoren: `PatternsForAPIDesign:2022`.

### Publikationen (Auswahl)

Die Patterns wurden zuerst als Konferenzpapiere begutachtet; die Papierfassungen sind teils
ausführlicher als die Website und enthalten Varianten, die dort fehlen.

| Jahr | Publikation | Deckt ab |
|---|---|---|
| 2017 | Interface Representation Patterns — Crafting and Consuming Message-Based Remote APIs (EuroPLoP) | frühe Fassungen der [Structure](category-structure.md)-Patterns |
| 2018 | Interface Quality Patterns — Communicating and Improving the Quality of Microservices APIs (EuroPLoP) | fünf [Quality](category-quality.md)-Patterns |
| 2018 | Guiding Architectural Decision Making on Quality Aspects in Microservice APIs (ICSOC) | Entscheidungsmodell zu Qualität |
| 2019 | Interface Evolution Patterns — Balancing Compatibility and Extensibility across Service Life Cycles (EuroPLoP) | die [Evolution](category-evolution.md)-Kategorie |
| 2020 | Interface Responsibility Patterns: Processing Resources and Operation Responsibilities (EuroPLoP) | Endpunktrollen und Operationsverantwortungen |
| 2020 | Data-Oriented Interface Responsibility Patterns: Types of Information Holder Resources (EuroPLoP) | die Information-Holder-Typen |
| 2017/2019 | Introduction to Microservice API Patterns (MAP), Microservices-Konferenz-Postproceedings | Gesamtüberblick, Domänenmodell, Pattern-Template |

### Weitere Ressourcen

- **Vorträge (PDF):** „API First with Patterns for API Design" (Keynote, niederländische
  API-Community, November 2025), „APIs as Service Activators" (ZEUS 2023), „Microservice API
  Patterns – Step by Step" (APICon 2022), eine längere JAX-2021-Präsentation, eine Präsentation
  zu DDD/MAP/MDSL (2020) und die Keynote von der Microservices-2019-Konferenz, mit der die
  Website angekündigt wurde. Alle Patterns sind zusätzlich als Web-Slideshow verfügbar.
- **Video und Podcast:** ein 24-Minuten-Video mit Daniel Lübke auf Erik Wildes YouTube-Kanal
  (Februar 2022) und eine einstündige Folge in Vaughn Vernons „Add Dot"-Reihe.
- **Beispielimplementierung:** *Lakeside Mutual*, eine fiktive Versicherung, deren Kern-Capabilities
  (Kunden-, Vertrags- und Risikomanagement) als Microservices mit zugehörigen Frontends
  implementiert sind. Öffentliches GitHub-Repository; viele Patterns sind dort belegt.
- **MDSL:** die Microservice Domain Specific Language mit Werkzeugen; unterstützt alle Patterns
  und wird in den [Tutorials](tutorials.md) für Vertragsausschnitte verwendet.
- **Design Practice Repository (DPR):** Methodensammlung auf GitHub und LeanPub, deren
  „Stepwise Service Design" viele MAP-Patterns einbindet.
- **Interface Refactoring Catalog:** Schwesterprojekt, seit Juli 2025 abgeschlossen; elf
  MAP-Patterns dienen dort als Refactoring-Ziele.
- **„A Pattern a Day":** seit Juni 2026 ein kostenloses, E-Mail-basiertes Lernangebot von
  Co-Autor Daniel Lübke, das alle 44 Patterns des Buchs durchläuft.

### Die Autoren

| Autor | Zugehörigkeit |
|---|---|
| Olaf Zimmermann | beratender Softwarearchitekt und Senior Researcher |
| Mirko Stocker | Professor an der OST Rapperswil (Schweiz); Managing Director, LegalGo GmbH |
| Daniel Lübke | Digital Solution Architecture, Hannover |
| Uwe Zdun | Universität Wien, Software Architecture Research Group |
| Cesare Pautasso | Software Institute, USI Lugano (Schweiz) |

Die Forewords stammen von Vaughn Vernon (Serienherausgeber) und Frank Leymann. Die Autoren danken
den Shepherds und Writers'-Workshop-Teilnehmern der EuroPLoP 2017–2020 sowie den Studierenden der
HSR/OST-Vorlesungen seit 2018 — ein Hinweis darauf, wie stark die Sprache aus Peer-Review
entstanden ist.

## Verwandte Wiki-Seiten

- [Überblick MAP](overview.md) · [Introduction](introduction.md) · [Primer](primer.md)
- [Cheat Sheet](cheatsheet.md) — die Online-Vorgängerfassung von Anhang A
- [Tutorials](tutorials.md) — Kapitel 3 des Buchs ist das großformatige Tutorial
- [Verwandte Pattern-Sprachen](related-pattern-languages.md)

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/book)
