---
title: Primer — Microservices, Patterns, MAP
kategorie: Meta
quelle: https://microservice-api-patterns.org/primer
---

# Primer: Was sind Microservices? Warum Patterns?

## Worum es geht

Der Primer ist die Kurzfassung der [Introduction](introduction.md) in drei Abschnitten: Was sind
Microservices, warum wählt MAP das Pattern-Format, und was genau ist MAP. Er ist der geeignete
Einstieg für Leser, die den Kontext einordnen wollen, ohne die lange Motivationsdiskussion zu
lesen.

## Inhalt

### Was sind Microservices?

Microservices exponieren Geschäftsfähigkeiten lose gekoppelt. Sie sind so entworfen, dass sie
unabhängig deploybar sind — was sie zugleich unabhängig skalierbar und austauschbar macht. Weil
sie ihren eigenen Zustand kapseln, müssen sie über nachrichtenbasierte Remote-APIs kommunizieren.
Historisch sind sie aus SOA hervorgegangen; typischerweise laufen sie in Virtualisierungscontainern
(Docker), die geclustert werden — die Quellseite nennt hier ausdrücklich Kubernetes als Beispiel.
Weitere Tenets: polyglotte Programmierung und Persistenz, dezentrale Continuous Delivery,
Ende-zu-Ende-Monitoring als DevOps-Praxis.

**Nutzen:**

- Agile Entwicklung mit Continuous Delivery — jeder Service gehört genau einem Team, das ihn
  unabhängig entwickelt, deployt und betreibt.
- Eignung für „IDEAL" cloud-native Anwendungen (isolierter Zustand, Verteilung, Elastizität,
  automatisiertes Management, lose Kopplung), inklusive horizontaler On-Demand-Skalierung.
- Inkrementelle Migration von Monolithen — geringeres Risiko bei Modernisierungsvorhaben.

**Preis:**

- Kommunikationsoverhead im verteilten System plus schlechte API-Entwürfe schlagen direkt auf die
  Performance durch.
- Eine große Zahl von Services erfordert disziplinierte Lebenszyklusverwaltung, Monitoring und
  Debugging.
- Single Points of Failure und Kaskadenfehler müssen vermieden werden (Redundanz, Circuit Breaker),
  damit ausfallende Downstream-Instanzen nicht das Gesamtsystem mitreißen.
- Datenkonsistenz und Zustandsverwaltung werden schwierig, wenn monolithische, zustandsbehaftete
  Anwendungen zerlegt werden.
- Autonomie und Konsistenz sind bei klassischer Backup-/Disaster-Recovery-Strategie nicht
  gleichzeitig garantierbar (das *BAC-Theorem* von Pardon, Pautasso und Zimmermann).

### Warum Patterns?

Patterns sind generelle, wiederverwendbare Lösungen für wiederkehrende Probleme — keine fertigen
Entwürfe, sondern Umrisse. Die Hillside-Group-Definition (nach Christopher Alexander) beschreibt
ein Pattern als dreiteilige Regel aus Kontext, einem darin wiederholt auftretenden Kräftesystem
und einer Konfiguration, die diese Kräfte auflöst.

Die Gründe der Autoren für dieses Format:

- Pattern-Namen bilden ein Fachvokabular, eine Ubiquitous Language.
- Patterns folgen einem strukturierten, wiedererkennbaren Textschema.
- Sie sind an den Rändern weich — sie skizzieren, statt Blaupausen vorzugeben.
- Sie werden aus der Praxis gehoben, nicht erfunden, und in Writers' Workshops gehärtet.
- Lösungsskizzen visualisieren den Entwurf (UML, DSLs, Icon-Sprachen, informelle Bilder).
- Pattern-*Sprachen* definieren Beziehungen zwischen Patterns und stützen kriterienbasierte
  Entscheidungen.

Die prägnanteste Formulierung der Seite: Ein Pattern-Autor ist „a journalist/reporter, not an
inventor/designer, first and foremost".

### Was ist MAP?

MAP behandelt Entwurf und Evolution von Web-APIs mit Fokus auf die *Payload-Repräsentationen* —
den Nachrichteninhalt zwischen Provider und Consumer. Vier Aspekte:

1. Struktur der Nachrichten und der darin kritischen Elemente.
2. Rollen und Verantwortungen der API-Aufrufe.
3. Wirkung des Nachrichteninhalts auf die API-Qualität.
4. API-Beschreibungen als Mittel für Governance und Evolution.

Die Patterns gelten für jede Remote-API mit einfachen Dokumentnachrichten — nicht nur
Microservices, nicht nur synchrones HTTP, auch queue-basiert asynchron. Nicht abgedeckt sind
Integrationsstile und Infrastrukturarchitekturen; dafür verweist die Seite auf
*Enterprise Integration Patterns* und *Cloud Computing Patterns* (siehe
[Verwandte Pattern-Sprachen](related-pattern-languages.md)).

## Verwandte Wiki-Seiten

- [Überblick MAP](overview.md) · [Introduction](introduction.md) — die ausführliche Fassung
- [Terminologie](terms.md) · [Cheat Sheet](cheatsheet.md) · [Tutorials](tutorials.md)
- [Verwandte Pattern-Sprachen](related-pattern-languages.md)
- [Buch und Ressourcen](book-and-resources.md)

## Bezug zu Kubernetes / KRM

Der Primer nennt Kubernetes beiläufig als Clustering-Technologie für Container — also in der
Rolle der *Infrastruktur*, die MAP ausdrücklich nicht behandelt. Die interessante Beobachtung
dieses Wikis ist, dass Kubernetes selbst zugleich ein Untersuchungsgegenstand *für* MAP ist: Die
Kubernetes-API ist eine der am gründlichsten durchdachten nachrichtenbasierten Remote-APIs
überhaupt, mit einer expliziten API-Description-Kette (Discovery, OpenAPI), einem strikten
Evolutionsregime und einem einheitlichen Repräsentationsschema.

Die im Primer genannten Microservice-Kosten treffen dabei auf konkrete KRM-Mechanismen: Gegen den
Kommunikationsoverhead steht das Watch-/Informer-Modell (ein Stream statt Polling), gegen
Kaskadenfehler stehen API Priority and Fairness und clientseitige Rate Limiter, gegen die
Konsistenzproblematik die bewusste Beschränkung auf Atomarität *je Objekt* plus
level-getriggerte Reconciliation, die inkonsistente Zwischenzustände als Normalfall behandelt
statt sie zu verhindern.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/primer)
