---
title: Introduction — Microservices und APIs
kategorie: Meta
quelle: https://microservice-api-patterns.org/introduction
---

# Introduction: Microservices und APIs

## Worum es geht

Die Introduction-Seite ist der lange Motivationstext von MAP: Sie begründet, warum
Service-API-Design ein eigenständiges, schwieriges Problem ist, grenzt MAP gegen benachbarte
Wissensbestände ab und erklärt Aufbau und Lesereihenfolge der Pattern-Sprache. Die Autoren merken
selbst an, dass das Buch weniger microservice-zentriert ist als dieser Text.

## Inhalt

### Warum API-Design ein eigenes Problem ist

Microservice-Architekturen sind aus SOA hervorgegangen: unabhängig deploybare, skalierbare und
austauschbare Services mit je einer Verantwortung, die eine Geschäftsfähigkeit modelliert, eigenen
Zustand kapseln und über nachrichtenbasierte Remote-APIs lose gekoppelt kommunizieren — meist in
Containern, oft mit polyglotter Programmierung und Persistenz, DevOps-Praktiken und dezentraler
Continuous Delivery. Über Infrastruktur und Betrieb ist dabei viel geschrieben worden, über den
API-Schnitt selbst wenig. Genau diese Lücke adressiert MAP.

Drei Problemklassen machen den Entwurf schwer:

- **Anforderungsvielfalt:** Die Bedürfnisse der Clients unterscheiden sich und ändern sich
  fortlaufend. Der Provider muss entscheiden, ob er einen guten Kompromiss in einer einheitlichen
  API anbietet oder jeden Client einzeln bedient.
- **Design-Mismatches:** Was Backendsysteme können und wie sie strukturiert sind, deckt sich nicht
  mit dem, was Clients erwarten. Die Differenz muss irgendwo überbrückt werden.
- **Zielkonflikte:** Provider wollen innovieren und schnell ändern, Clients wollen Stabilität.
  Eine API zu veröffentlichen heißt, Kontrolle abzugeben: Jedes exponierte Datum kann genutzt
  werden, mitunter unerwartet — und ein Informationsvorsprung gegenüber dem Wettbewerb kann
  verlorengehen.

Daraus folgen die klassischen Trade-offs: wenige große Aufrufe gegen viele feingranulare;
stabile, standardisierte, breite Schnittstellen gegen schnell wechselnde, spezialisierte;
Datenkonsistenz gegen Verfügbarkeit und Antwortzeit.

### Abgrenzung: was MAP nicht ist

Für RESTful HTTP im Detail gibt es Kochbücher (Allamaraju), für asynchrones Messaging die
*Enterprise Integration Patterns* (Hohpe/Woolf), für Service-Identifikation strategisches
Domain-Driven Design (Evans, Vernon), für Infrastruktur die SOA-, Cloud- und
Microservices-Patternsammlungen. MAP setzt darauf auf und behandelt das, was dort fehlt: die
Strukturierung des Datenaustauschs. Der Bezugspunkt ist Pat Hellands Unterscheidung von *data on
the inside* und *data on the outside* — beide unterscheiden sich erheblich in Veränderlichkeit,
Lebensdauer, Genauigkeit, Konsistenz- und Schutzbedarf, weshalb Datenmodellierungs-Patterns für
Speicher nicht einfach übertragbar sind.

### Warum das Pattern-Format

- Pattern-Namen bilden ein Fachvokabular, eine *Ubiquitous Language*. Vorbild ist EIP, dessen
  Namen zur Lingua franca queue-basierten Messagings wurden und in Frameworks implementiert sind.
- *Forces* und *Consequences* unterstützen informierte Entscheidungen — auch über die Nachteile
  einer Lösung.
- Patterns sind „soft around their edges": Sie skizzieren Lösungen, statt Blaupausen zu liefern,
  die blind zu befolgen wären.

Der Autor eines Patterns ist, wie es der [Primer](primer.md) formuliert, in erster Linie
Journalist und Reporter, nicht Erfinder.

### Kategorien und Pattern-Template

Die vier Fragen, aus denen die [Kategorien](overview.md) entstehen, betreffen *Structure* der
Nachrichten, *Quality*-Wirkung des Nachrichteninhalts, Rollen und *Responsibility* der Operationen
sowie API-Beschreibungen als Mittel für Governance und *Evolution*.

Das Template: *Context* setzt Voraussetzungen, *Problem* stellt die Entwurfsfrage, *Forces*
erklären, warum sie schwer ist (oft unter Verweis auf konfligierende Qualitätsattribute, teils mit
einer Nicht-Lösung), *Solution* antwortet und beschreibt Funktionsweise, Varianten, Beispiel und
Implementierungshinweise, *Consequences* diskutieren die Auflösung der Kräfte und Alternativen,
*Known Uses* belegen reale Anwendungen, und Beziehungen zu anderen Patterns zeigen, was als
Nächstes interessant wird. Im Buch heißen die Abschnitte anders: *When and Why to Apply*,
*How it Works*, *Discussion*, *Related Patterns*, *More Information*.

Empfohlene Lesereihenfolge beim ersten Durchgang:
**1. Problem → 2. Solution → 3. Known Uses → 4. Context → 5. Forces → 6. Consequences.**

Ein Patterntext ist als eigenständiger Fachartikel von typisch 2000 bis 3000 Wörtern angelegt.
Wem das zu viel ist, dem empfehlen die Autoren, die Texte als Nachschlagewerk zu behandeln.

### Schlussbemerkung der Autoren

MAP ist ein Freiwilligenprojekt. Die meisten Patterns wurden von der Pattern-Community und
erfahrenen Reviewern in Workshops durchgearbeitet. Die Autoren betonen, dass die Patterns auch
dann nützlich sind, wenn man sich gegen Microservices entschieden hat — die Entwurfsprobleme beim
Exponieren von Remote-APIs bleiben, unabhängig davon, ob der Begriff in Mode ist.

## Verwandte Wiki-Seiten

- [Überblick MAP](overview.md) — die Kurzfassung
- [Primer](primer.md) — knappere Variante desselben Stoffs
- [Terminologie](terms.md) — das API-Domänenmodell
- [Cheat Sheet](cheatsheet.md) · [Tutorials](tutorials.md)
- [Verwandte Pattern-Sprachen](related-pattern-languages.md)

## Bezug zu Kubernetes / KRM

Die drei Problemklassen der Introduction lassen sich am Kubernetes-Projekt unmittelbar
wiedererkennen — und Kubernetes hat für jede eine charakteristische Antwort gewählt.

Gegen die **Anforderungsvielfalt** setzt KRM konsequent den unifizierten Kompromiss: eine API für
alle Clients, kein Backend-for-Frontend, keine clientspezifischen Varianten. Angepasst wird über
Content Negotiation (`Table`- und `PartialObjectMetadata`-Repräsentationen) statt über eigene
Endpunkte. Gegen den **Design-Mismatch** steht die strikte Trennung von externen versionierten
Typen und internen Typen mit generierter Konvertierung — Clients sehen nie die interne Struktur.

Am interessantesten ist der **Zielkonflikt** zwischen Änderungswunsch und Stabilitätsbedürfnis. Die
Introduction beschreibt ihn als Verlust an Kontrolle: Alles, was exponiert wird, wird genutzt,
mitunter unerwartet. Genau das ist die Erfahrungsgrundlage der Kubernetes-Deprecation-Policy und
der Regel, dass GA-APIs innerhalb einer Major-Version nicht entfernt werden. Kubernetes hat den
Konflikt nicht aufgelöst, sondern in ein Stufenmodell überführt: Wer maximale Änderungsfreiheit
will, veröffentlicht Alpha (siehe [Experimental Preview](../evolution/ExperimentalPreview.md));
wer Verlässlichkeit zusagt, geht auf GA und akzeptiert, dass die Struktur eingefroren ist. Die
Entscheidung wird damit explizit und für Clients ablesbar gemacht — das ist präzise die Wirkung,
die MAP von Patterns wie [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md)
erwartet.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/introduction)
