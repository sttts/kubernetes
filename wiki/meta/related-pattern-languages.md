---
title: Verwandte Pattern-Sprachen
kategorie: Meta
quelle: https://microservice-api-patterns.org/relatedPatternLanguages
---

# Verwandte Pattern-Sprachen

## Worum es geht

MAP ist bewusst als *ergänzende* Sprache angelegt: Sie deckt die Repräsentation ausgetauschter
Daten ab und setzt für alles andere auf bestehende Pattern-Werke auf. Diese Seite ordnet MAP in
diese Landschaft ein — nützlich vor allem, um zu erkennen, wann man *nicht* bei MAP nachschlagen
sollte.

## Inhalt

### Verteilung und Remoting

| Werk | Fokus | Verhältnis zu MAP |
|---|---|---|
| POSA Vol. 4 (Buschmann, Henney, Schmidt 2007) | Pattern-Sprache, die Patterns verteilter Systeme von Architektur- bis Detailebene verbindet | Rahmen, in den sich MAP einfügt |
| Remoting Patterns (Voelter, Kircher, Zdun 2004) | Broker-Entwurf und Middleware-Interna; Remote Objects, Servants, Lifecycle, asynchrone Invocation | MAP verfeinert die API-Aspekte für nachrichtenbasierte Remote-APIs |
| Enterprise Integration Patterns (Hohpe, Woolf 2003) | Asynchrones Messaging: Routing, Transformation, Guaranteed Delivery | MAP behandelt den *Inhalt* der Nachrichten, EIP ihren *Transport* |
| Patterns of Enterprise Application Architecture (Fowler 2002) | Remote Facade, Data Transfer Object u. a. | Berührt Remote-API-Entwurf, ohne ihn systematisch zu behandeln |
| Domain-Driven Design (Evans 2003) | Bounded Context, Aggregate, Service | Liefert MAP die Endpunktidentifikation (siehe [Cheat Sheet](cheatsheet.md)) |
| Design Patterns (Gamma et al. 1995) | Facade, Proxy u. a. | Auch für Remote-APIs relevant, aber auf Objektebene |

### Datenmodellierung

*Data Model Patterns* (Hay 1996) und die Archetypen von Arlow/Neustadt (2004) behandeln
Datenrepräsentation und Bedeutung — aber im Kontext von *Speicherung und Darstellung*, nicht von
*Transport*. Weil sich Kräfte und Lösungen dadurch unterscheiden, sind sie nicht direkt
übertragbar; genau hier liegt Pat Hellands Unterscheidung von „data on the inside" und „data on
the outside", die MAP als Ausgangspunkt nimmt.

### Serviceorientierung

- **Service Design Patterns** (Daigneau 2011): Patterns auf Plattform-/Technologieebene für REST
  und WSDL/SOAP; behandelt unter anderem Contract Versioning — die nächste Nachbarschaft zur
  [Evolution-Kategorie](category-evolution.md).
- **Process-Driven SOA** (Hentrich, Zdun 2011): Orchestrierung über Prozess-/Workflow-Engines,
  eine Abstraktionsebene über MAP.
- **SOA Patterns** (Rotem-Gal-Oz 2012): weitgehend infrastruktur- und plattformzentriert;
  untersucht Nachrichteninhalt und -struktur nicht in der Tiefe.
- **Service Interaction Patterns** (Barros, Dumas, ter Hofstede 2005) und **Conversation
  Patterns** (Hohpe 2007, sowie Pautasso/Ivanchikj/Schreier 2016 für REST): zustandsbehaftete
  Interaktionen aus mehreren Nachrichtenaustauschen. Die Beziehung ist wechselseitig: Jeder
  Austausch in einer Konversation braucht Repräsentationen, und grobgranulare APIs führen zu
  einfachen Konversationen, feingranulare zu geschwätzigen.
- **Microservices Patterns** (Richardson 2018) und **Cloud Computing Patterns** (Fehling et al.
  2014): eigene Patterns für cloud- und microservice-spezifische Architekturen.
- **Cloud Adoption Patterns** (Kyle Brown et al.): laufende Arbeit, komplementär.

### Weitere Quellen

Ergänzend zu Pattern-Sprachen nennt die Seite plattformspezifische Best-Practice-Sammlungen, etwa
das *RESTful Web Services Cookbook* (Allamaraju 2010) und die Entscheidungssammlung in
*Perspectives on Web Services* (Zimmermann, Tomlinson, Peuser 2003). Erwähnt wird außerdem der
*Interface Refactoring Catalog* als Schwesterprojekt, das MAP-Patterns als Refactoring-Ziele
verwendet.

## Verwandte Wiki-Seiten

- [Überblick MAP](overview.md) · [Primer](primer.md) · [Introduction](introduction.md)
- [Buch und Ressourcen](book-and-resources.md) — das Buch diskutiert dieselbe Landschaft
- [Terminologie](terms.md) — die aus DDD übernommenen Begriffe

## Bezug zu Kubernetes / KRM

Für Kubernetes ist die Abgrenzung der Seite direkt praktisch: Fragen zu Controller-Verhalten,
Retry, Backoff, Ereignisverarbeitung und Idempotenz gehören nicht zu MAP, sondern in die Nähe von
*Enterprise Integration Patterns* und den *Conversation Patterns* — die kanonische
KRM-Konversation aus initialem `list` mit `resourceVersion` und anschließendem `watch` ist ein
Conversation-Pattern, kein MAP-Pattern. Umgekehrt liefern die DDD-Begriffe den brauchbarsten
Rahmen für den Ressourcenschnitt: Ein `Kind` entspricht am ehesten einem Aggregate, und die
Atomaritätsgarantie je Objekt ist genau die Aggregate-Konsistenzgrenze.

Das *Published Language*-Konzept ist die präziseste Beschreibung dessen, was eine versionierte
Kubernetes-API-Gruppe ist: ein veröffentlichtes Austauschvokabular, das weder Provider noch Client
intern bindet — die internen Go-Typen sind ausdrücklich kein Teil davon und werden nie
serialisiert.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/relatedPatternLanguages)
