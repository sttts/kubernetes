---
title: Überblick — Microservice API Patterns (MAP)
kategorie: Meta
quelle: https://microservice-api-patterns.org/
---

# Überblick: Microservice API Patterns (MAP)

## Worum es geht

*Microservice API Patterns* (MAP), seit dem Buch von 2022 auch *Patterns for API Design*, ist eine
Pattern-Sprache für Entwurf und Evolution nachrichtenbasierter Remote-APIs. Der Blickwinkel ist
ungewöhnlich eng gewählt und genau darin liegt der Wert: MAP behandelt **„data on the outside"** —
die Repräsentationen und Nutzlasten, die beim API-Aufruf tatsächlich über die Leitung gehen.
Integrationsstile, Middleware-Architekturen und Infrastruktur bleiben bewusst außen vor; dafür
gibt es andere Pattern-Sprachen (siehe [Verwandte Pattern-Sprachen](related-pattern-languages.md)).

Die drei Leitfragen der Startseite:

- Wie viele Services sollen remote exponiert werden, und wie groß sollen sie sein?
- Wie bleiben die Services und ihre Clients lose gekoppelt? Wie viele Daten werden je
  Request/Response ausgetauscht, und wie oft?
- Welche Nachrichtenrepräsentationen sind geeignet, und wie einigt man sich auf ihre Bedeutung?

Die Patterns — die Website listet aktuell 45, das Buch spricht von 44 — sind aus öffentlichen
Web-APIs sowie aus Integrations- und Entwicklungsprojekten
der Autoren *gehoben*, nicht erfunden — und anschließend in EuroPLoP-Writers'-Workshops
begutachtet und gehärtet. Sie gelten ausdrücklich nicht nur für Microservices, sondern für jede
Remote-API, die einfache Textnachrichten austauscht, synchron über HTTP wie asynchron über Queues.

## Inhalt

### Die fünf Kategorien

| Kategorie | Leitfrage | Patterns |
|---|---|---|
| [Foundation](category-foundation.md) | Welche Systeme werden integriert, von wo ist die API erreichbar, wie wird sie dokumentiert? | 6 |
| [Responsibility](category-responsibility.md) | Welche architektonische Rolle spielt jeder Endpunkt, welche Verantwortung jede Operation? | 11 |
| [Structure](category-structure.md) | Wie viele Repräsentationselemente, wie strukturiert, wie annotiert? | 11 |
| [Quality](category-quality.md) | Wie erreicht man ein Qualitätsniveau kosteneffizient, und wie kommuniziert man die Trade-offs? | 10 |
| [Evolution](category-evolution.md) | Wie geht man mit Versionierung, Supportzeiträumen und brechenden Änderungen um? | 7 |

Die Reihenfolge ist zugleich die typische Entscheidungsreihenfolge — von der Grundsatzfrage
„welche API überhaupt" bis zur Lebenszyklusfrage „wie lange". Sie ist aber keine Vorschrift: Die
Autoren betonen mehrfach, dass API-Design ein *wicked problem* ist, das sich nicht serialisieren
lässt, und bieten deshalb bewusst mehrere Einstiegspunkte an.

### Einstiegspunkte

| Einstieg | Wann nützlich |
|---|---|
| [Cheat Sheet](cheatsheet.md) | „Ich habe folgendes Problem" — Issue-zu-Pattern-Tabellen, chronologisch geordnet |
| [nach Scope](navigation-byscope.md) | „Ich arbeite gerade an der Nachricht / am Endpunkt / an der API als Ganzes" |
| [nach Qualität](navigation-byquality.md) | „Mich beschäftigt Performance / Kopplung / Sicherheit" |
| [nach Phase](navigation-byphase.md) | „Wir sind in Sprint 0 / kurz vor Go-Live" |
| [nach Rolle](navigation-byrole.md) | „Ich bin Client-Entwickler / Designer / Product Owner" |
| [Tutorials](tutorials.md) | Geführtes Durchspielen an einem Fallbeispiel, 20 bis 90 Minuten |
| [Introduction](introduction.md) / [Primer](primer.md) | Motivation, Hintergrund, warum überhaupt Patterns |
| [Terminologie](terms.md) | Was heißt Endpoint, Operation, Representation, Published Language? |

### Pattern-Template

Jeder Pattern-Text folgt derselben Gliederung: *Context* (Anwendbarkeitsvoraussetzungen),
*Problem* (die Entwurfsfrage), *Forces* (warum sie schwer zu beantworten ist), *Solution*
(Antwort, Funktionsweise, Varianten, Beispiel, Implementierungshinweise), *Consequences*
(inwieweit die Kräfte aufgelöst werden, plus Alternativen), *Known Uses* (belegte reale
Anwendungen) und Beziehungen zu anderen Patterns. Für den ersten Durchgang empfehlen die Autoren
die in der Pattern-Community übliche Lesereihenfolge: Problem → Solution → Known Uses → Context →
Forces → Consequences.

### Autoren und Buch

MAP ist ein Freiwilligenprojekt von fünf Autoren: **Olaf Zimmermann** (beratender
Softwarearchitekt, Senior Researcher), **Mirko Stocker** (OST Rapperswil), **Daniel Lübke**
(Hannover), **Uwe Zdun** (Universität Wien) und **Cesare Pautasso** (USI Lugano). Das Buch
*Patterns for API Design: Simplifying Integration with Loosely Coupled Message Exchanges*
erschien 2022 bei Addison-Wesley in der Vaughn-Vernon-Signature-Serie (ISBN 9780137670109);
Details unter [Buch und Ressourcen](book-and-resources.md).

### Die POINT-Prinzipien

Neben der Pattern-Sprache selbst tragen die MAP-Autoren fünf API-Entwurfsprinzipien als Backronym
**POINT** vor, eingeführt 2021 im Blogpost „APIs should get to the POINT":

| Buchstabe | Prinzip | Kurzform |
|---|---|---|
| **P** | *purposeful* | Die API hat einen klaren fachlichen Zweck, kein generischer Datenbankzugriff |
| **O** | style-**o**riented | Ein Architekturstil wird gewählt und konsequent durchgehalten |
| **I** | *isolated* | Endpunkte sind lose gekoppelt und unabhängig deploybar |
| **N** | cha**n**nel-neutral | Die API bindet sich nicht an ein einzelnes Transportprotokoll |
| **T** | *T-shaped* | Breite Übersichtsoperationen (Listen von Ids) neben tiefen Detailoperationen |

POINT ist nicht Teil der Pattern-Sprache, wird aber auf einzelnen Pattern-Seiten unter
*More Information* als Bezugsrahmen herangezogen — etwa bei
[Id Element](../structure/IdElement.md), wo Ids die horizontale Leiste des T tragen und die
Detailrepräsentation die vertikale.

Der Name der Sprache selbst ist ein Wortspiel: MAP als „Landkarte" durch den Entwurfsraum.

## Verwandte Wiki-Seiten

- Kategorien: [Foundation](category-foundation.md) · [Responsibility](category-responsibility.md) · [Structure](category-structure.md) · [Quality](category-quality.md) · [Evolution](category-evolution.md)
- [Cheat Sheet](cheatsheet.md) · [Terminologie](terms.md) · [Tutorials](tutorials.md)
- [Introduction](introduction.md) · [Primer](primer.md) · [Verwandte Pattern-Sprachen](related-pattern-languages.md) · [Buch und Ressourcen](book-and-resources.md)
- **[KRM vs. MAP](../KRM-vs-MAP.md)** — die durchgehende Analyse

## Bezug zu Kubernetes / KRM

Das Kubernetes Resource Model ist ein interessanter Testfall für MAP, weil es fast alle
Entwurfsfreiheiten, die MAP als Entscheidungen anbietet, *vorab und einheitlich* entscheidet. Wo
MAP je Endpunkt eine eigene Operationsmenge mit eigenen Signaturen entwerfen lässt, gibt KRM eine
feste Verbmenge vor (`get`, `list`, `watch`, `create`, `update`, `patch`, `delete`,
`deletecollection`) und ein festes Nachrichtengerüst (`apiVersion`, `kind`, `metadata`, `spec`,
`status`). Der Preis ist Ausdrucksverlust; der Gewinn ist, dass generische Werkzeuge — `kubectl`,
Informer, Server-Side Apply, RBAC, Admission — für jede Ressource funktionieren, auch für
solche, die erst zur Laufzeit als CRD dazukommen.

Drei Kategorien werden dadurch verschieden getroffen: [Structure](category-structure.md) ist der
Ort, an dem in KRM praktisch die gesamte Designarbeit stattfindet.
[Evolution](category-evolution.md) ist in KRM formalisierter als in fast jeder anderen API, die
man kennt — mit Stabilitätsstufen, verlustfreier Versionskonvertierung, `Warning`-Headern und
einer verbindlichen Deprecation Policy. [Responsibility](category-responsibility.md) dagegen ist
weitgehend leergeräumt, weil das deklarative, level-getriggerte Modell die Zustandsübergangslogik
aus der API in die Controller verschiebt: Der Client äußert einen Wunsch (`spec`), er befiehlt
keinen Übergang.

Die ausführliche Gegenüberstellung — welches Pattern KRM (a) genauso, (b) anders, (c) gar nicht
umsetzt, und warum — steht in [KRM vs. MAP](../KRM-vs-MAP.md). Jede Pattern-Seite dieses Wikis
hat außerdem einen eigenen KRM-Abschnitt.

---
[← Index](../README.md) · [Cheat Sheet](cheatsheet.md) · [Quelle](https://microservice-api-patterns.org/)
