---
title: "Interface Representation Patterns (EuroPLoP 2017)"
kategorie: Paper
autoren: Olaf Zimmermann, Mirko Stocker, Daniel Lübke, Uwe Zdun
venue: EuroPLoP 2017
quelle: http://eprints.cs.univie.ac.at/5161/1/WADE-EuroPlop2017Paper-FinalSubmissionOct19.pdf
---

# Interface Representation Patterns (EuroPLoP 2017)

## Bibliografische Angaben

Olaf Zimmermann, Mirko Stocker, Daniel Lübke, Uwe Zdun: *Interface Representation Patterns — Crafting
and Consuming Message-Based Remote APIs.* Proceedings of the 22nd European Conference on Pattern
Languages of Programs (EuroPLoP '17), Irsee, 12.–16. Juli 2017, 36 Seiten. ACM, DOI
[10.1145/3147704.3147734](https://doi.org/10.1145/3147704.3147734). HSR Rapperswil, innoQ Schweiz,
Universität Wien. Erste Publikation der Pattern-Sprache, aus der später
[microservice-api-patterns.org](https://microservice-api-patterns.org) und das Buch *Patterns for
API Design* (2022) wurden.

## Worum es geht

Die etablierten Pattern-Sprachen für verteilte Systeme — POSA 4, *Remoting Patterns*, *Enterprise
Integration Patterns*, Daigneaus *Service Design Patterns*, SOA- und Cloud-Kataloge — behandeln laut
den Autoren die Schnittstelle selbst nicht: Sie beschreiben Middleware, Konversationen und
Infrastruktur, aber nicht **Struktur und Bedeutung des Nachrichteninhalts**. Genau diese Lücke füllt
das Paper, für *message-based remote APIs*: RESTful HTTP, WSDL/SOAP, WebSockets, ausdrücklich auch
gRPC „trotz seines Namens" (S. 2). Als Vokabular dient die *Data Transfer Representation* (DTR), das
Wire-Level-Pendant zum DTO ohne dessen Annahme über das Programmierparadigma der Endpunkte;
Top-Level-Elemente heißen *in* bzw. *out parameters* und bilden die *message signature*.

**Methodik.** Pattern Mining aus eigenen Integrationsprojekten, aus 18 öffentlichen Web-APIs (u.a.
Facebook, GitHub, JIRA, Microsoft Graph, PayPal, Twitter) und aus der Literatur, in einer explizit
dokumentierten Serie von *problem jams*, *forces jams* und *known uses jams* — erst von den Autoren
selbst beantwortet, dann mit ihrem Netzwerk. Beitrag sind fünf ausformulierte *Structure*-Patterns,
die Kategorienstruktur der Sprache und ein Backlog von rund 40 Kandidaten (Tabelle I).

## Behandelte Patterns

| Pattern | Paper-Abschnitt | Entsprechung in diesem Wiki |
|---|---|---|
| *Atomic Parameter* (a.k.a. Single Scalar Representation, Dot) | 4.1, S. 9–13 | [Atomic Parameter](../structure/AtomicParameter.md) |
| *Atomic Parameter List* (a.k.a. Multiple Scalar Representations) | 4.2, S. 13–16 | [Atomic Parameter List](../structure/AtomicParameterList.md) |
| *Parameter Tree* (a.k.a. Single Complex Representation, Bar) | 4.3, S. 16–21 | [Parameter Tree](../structure/ParameterTree.md) |
| *Parameter Forest* (a.k.a. Parameter Comb, Hybrid Parameter List) | 4.4, S. 21–25 | [Parameter Forest](../structure/ParameterForest.md) |
| *Pagination* (a.k.a. Query with Partial Result Sets, Response Sequence) | 4.5, S. 25–31 | [Pagination](../quality/Pagination.md) |

Nur als Kandidaten skizziert, inzwischen ausformuliert: *Id/Link/Entity/Metadata Parameter* →
[Id](../structure/IdElement.md)/[Link](../structure/LinkElement.md)/[Data](../structure/DataElement.md)/[Metadata Element](../structure/MetadataElement.md),
*API Key* → [API Key](../structure/APIKey.md), *Error Reporting* → [Error Report](../structure/ErrorReport.md),
*Wish List/Template* → [Wish List](../quality/WishList.md), *Embedded/Linked Reference Data* →
[Embedded Entity](../quality/EmbeddedEntity.md)/[Linked Information Holder](../quality/LinkedInformationHolder.md).

## Was das Paper gegenüber der Website ergänzt

Seit dem Buch kürzt die Website *Forces* auf Stichwortlisten und ersetzt *Consequences* durch einen
Buchverweis. Das Paper enthält beides ausformuliert; sein Abschnitt „Discussion" entspricht den
heutigen *Consequences*. Alle vier Struktur-Patterns teilen denselben Kräftesatz oberster Ebene —
Interoperabilität, Performance, Verarbeitungsaufwand, Wartbarkeit, Sicherheit —, dessen *Gewicht* mit
der Komplexität der Struktur wächst.

### Atomic Parameter

**Forces.** Der Kontrakt ist geteiltes Wissen und damit unmittelbar Kopplung; die Autoren verorten ihn
als *format autonomy* im Kopplungsmodell von Leymann. Der Zielkonflikt ist zweiseitig: Ein
**unterspezifizierter** Kontrakt erzeugt Interoperabilitätsprobleme, sobald Optionalität ins Spiel
kommt (Abwesenheit eines Werts vs. expliziter NULL-Wert) oder zwischen Repräsentationsvarianten
gewählt werden kann; ein **überspezifizierter** wird unflexibel und macht Rückwärtskompatibilität
schwer. Zweitens: Einfache Datenstrukturen führen zu feingranularen Kontrakten, komplexe zu
grobgranularen — die Repräsentationsentscheidung *ist* eine Granularitätsentscheidung. Drittens: Man
könnte alles als String oder Key-Value-Paar übertragen, aber das vergrößert das *implizit* geteilte
Wissen, erschwert Test und Wartung und bläht Nachrichten auf.

**Consequences.** Skalare sind in jeder Sprache leicht zu lesen, schreiben und verarbeiten und dank
ihrer Einfachheit und intrinsischen Kohäsion leicht abzusichern. Die Autoren relativieren das
ausdrücklich: Der Einfluss auf Interoperabilität und Bandbreite ist auf plattformunabhängiger
Pattern-Ebene gerade *nicht* offensichtlich. Negativ: Die Ausdrucksstärke ist begrenzt; eine
reichhaltige *Published Language* ließe sich nur hineincodiert übertragen. Nutzen viele API-Calls das
Pattern, entstehen geschwätzige, zustandsbehaftete Konversationen — dann ist der Wechsel zu *Atomic
Parameter List* oder *Parameter Tree* fällig. Aus den Hinweisen: Skalare sollen **nicht** binäre Daten
oder codierte Strukturen tunneln — für die Autoren ein Integrations-Antipattern.

### Atomic Parameter List

**Forces.** Die Kräfte verschärfen sich mit der Elementzahl nichtlinear: Mehrere Parameter eröffnen
Sicherheitsrisiken sowohl einzeln als auch *in Kombination*. Zwei Alternativen werden explizit
verworfen. Erstens mehrere Nachrichten mit je einem Skalar — geschwätzig, verschwendet Netzkapazität
und, bei zustandslosen Servern, Serverressourcen, weil der Sitzungszustand je Request aus Nachricht
und Headern rekonstruiert werden muss. Zweitens immer eine komplexe Struktur senden und Unbenötigtes
leer lassen — erzeugt Interoperabilitäts- und Evolutionsprobleme, höheren Testaufwand und
Verarbeitungsverschwendung, weil leere Strukturen trotzdem marshalled werden.

**Consequences.** Die Lösung balanciert Informationsbedarf gegen Verarbeitungs- und
Kommunikationsaufwand; da nur Basistypen verwendet werden, ist die Interoperabilität gut. Zwei
Einschränkungen, die die Website nicht mehr führt: Das Pattern ist in manchen Plattformen **nicht in
Reinform realisierbar** — viele Sprachen erlauben nur einen Rückgabewert, und die Standard-Mappings
nach JSON/XML Schema (JAX-RS, JAX-WS) folgen dem. Dafür führt das Paper die Variante *Atomic Parameter
Tree* ein: eine nachrichtenspezifische Sequenz von Skalaren unter einer Wurzel. In Technologien ohne
Typsystem (WebSockets) müssen sich beide Seiten separat auf das Marshalling einigen. Sicherheitsseitig
verursacht das Pattern *mehr* Aufwand als *Atomic Parameter*, abhängig von Schutzklasse und
semantischer Nähe der Parameter (Service Cutter, Gysel et al. 2016).

### Parameter Tree

**Forces.** Zusätzlich zum gemeinsamen Satz: Bei repetitiven oder verschachtelten Daten sind **Anzahl
der Datenelemente und Verschachtelungstiefe** die relevanten Größen. Die verworfene Alternative —
mehrere Nachrichten mit einfachen Parametern — ist hier nicht nur geschwätzig, sondern riskiert eine
Verletzung der losen Kopplung, weil semantische Abhängigkeiten zwischen Calls und Parametern entstehen.

**Consequences.** Der Baum ist ausdrucksstärker als die Liste, weil rekursiv anwendbar, und
**kohäsiver wegen der einzelnen Wurzel**. Lern- und Verarbeitungsaufwand hängen stark von Tiefe,
Breite und Plattform ab: In mobilen SDKs ist JSON-Array-Verarbeitung nativ. Negativ: Baumnavigation
kommt als Herausforderung hinzu; aufgeblähte Strukturen erhöhen den Verarbeitungsaufwand und
verschwenden Netzkapazität — insbesondere Daten, die der Client gar nicht braucht. Sicherheitsseitig
ist gut, dass nur *eine* Struktur abzusichern ist; andererseits kann ihr Inhalt Daten mit
**unterschiedlichem Schutzbedarf** mischen (Personenname neben Kreditkartennummer), was feldgranulare
Absicherung erschwert — selektive attributbasierte Zugriffskontrolle kann nötig werden. Konkrete Schwelle: **Ab etwa fünf bis sieben Top-Level-Elementen** ist ein Wechsel zu
*Parameter Forest* zu erwägen. Hinweise: Postels Robustheitsprinzip anwenden (ausgehende Daten streng
validieren, eingehende so lax wie möglich) und zurückhaltend mit dynamisch interpretierten
Key-Value-Strukturen sein — erst domänenspezifische Namen ermöglichen Code-Completion und
Testautomatisierung.

### Parameter Forest

**Forces.** Identischer Kräftesatz wie beim *Parameter Tree*; Zusatzkraft ist, dass selbst mehrere
Nachrichten mit je einem komplexen Parameter den Informationsbedarf nicht decken.

**Consequences.** Ähnliches Profil wie *Parameter Tree*, eine Stufe weiter. Interessant ist die
historische Einordnung: In manchen Plattformen (JAX-WS) sind Baum und Wald kaum unterscheidbar, und
die Autoren stellen den Bezug zur alten SOAP-Debatte her — ein *wrapped document/literal*-Envelope ist
ein *Parameter Tree* (eine Wurzel), *rpc/encoded* kann echte *Parameter Forests* erzeugen. Positiv und
auf der Website nicht mehr vorhanden: Weil einzelne Felder für Signaturen oder Public-Key-Informationen
reserviert werden können, lässt sich das Pattern als **Security Enabler** einsetzen; ein „Zahn des
Kamms" trägt dann ausschließlich Sicherheitsinformation. Hinweise: Anzahl der Bäume auf drei bis fünf
begrenzen, maschinenlesbare Schemas mitliefern, und besonders vorsichtig mit optionalen Lücken *in der
Mitte* der Kammstruktur sein — sie erschweren Unmarshalling und Security-Policy-Verarbeitung und
können zu Mehrdeutigkeiten, falschen Ergebnissen und Audit-Fehlern führen.

### Pagination

**Forces.** Entwurfskriterien laut Paper: Größe des Datenbestands und Zugriffsprofil; Variabilität
der Daten (gleich strukturiert? wie oft ändern sich die Definitionen?); verfügbarer Speicher pro
Request auf *beiden* Seiten sowie Aktualität gegen Änderungsdynamik; Netzeigenschaften; Sicherheit
und Robustheit. Ausformuliert kommen Argumente hinzu, die die Website nicht führt:

- Textformate (XML, aber auch JSON) verursachen hohe Parsing-Kosten und Volumen; Binärformate wie Avro
  oder Protocol Buffers helfen, brauchen aber Marshalling-Bibliotheken, die nicht überall verfügbar
  sind — Webbrowser sind das genannte Beispiel.
- Transferzeiten sind **nichtlinear in der Datengröße**: 1500 Byte passen in ein Ethernet-Paket, ein
  Byte mehr erzwingt zwei zu koordinierende Pakete.
- Erzeugen und Codieren großer Datenmengen ist teuer und öffnet einen **Angriffsvektor für Denial of
  Service**; Übertragungen über unzuverlässige Netze (Mobilfunk) brechen ab.
- Unbegrenzte Query-Kontrakte führen zu Out-of-Memory-Fehlern. Entwickler unterschätzen den
  Speicherbedarf regelmäßig; das bleibt **unbemerkt, bis nebenläufige Last auftritt**.

**Consequences.** Positiv: verdauliche Datenmengen, direkte Navigation im Datenbestand, weniger
Speicher- und Netzbedarf pro Antwort — zum Preis des Verwaltungsaufwands für die Paginierung selbst.
Das Paper listet dann die **Folgeentscheidungen** auf, die durch das Pattern erst entstehen: wo, wann
und wie die Seitengröße festgelegt wird (beeinflusst die Chattiness); wie Ergebnisse geordnet und
Seiten zugewiesen werden; wo und wie lange Zwischenergebnisse gespeichert werden (Löschpolitik,
Timeouts); ob Requests idempotent sein müssen; wie Teilantworten korreliert werden. Dazu Caching,
Filterung, Query-Nachverarbeitung sowie Isolationslevel und Locking. Konsistenzanforderungen
unterscheiden sich nach Client-Typ; die Autoren trennen **virtuelle von technischer Paginierung**.
Negativ: Paginierung verlagert Programmieraufwand auf den Consumer und — der wichtigste Punkt —
**alle Funktionen über den vollständigen Datensatz, etwa Suche, funktionieren nicht oder nur mit
Zusatzaufwand**.

## Pattern-Sprache und Zusammenhänge

**Kategorien (Abb. 2, S. 7).** Sieben Kategorien: *Foundations*, *Identification*, *Structure*,
*Responsibility*, *Quality*, *Evolution*, *API Management* (letztere noch unbearbeitet). Die drei
mittleren bilden den Kern („RSQ") und beantworten **warum**, **was** und **wie**. Qualität ist
Kernkategorie, weil Repräsentationen qualitätsbezogene Informationen *enthalten* und ihre Größe und
Struktur die Dienstgüte beeinflussen.

**Strukturbeziehungen (Abb. 3, S. 8).** Die vier Basis-Patterns entstehen aus zwei Fragen — einfache
vs. strukturierte Datentypen, ein vs. mehrere Parameter:

- *Atomic Parameters* sind Bestandteile von *Atomic Parameter Lists*; eine solche Liste wird in einen
  *Parameter Tree* refaktoriert, wenn sie zu groß wird oder semantische Abhängigkeiten entstehen.
- Ein *Parameter Tree* hat genau eine Wurzel; *Atomic Parameters* sind seine Blätter.
- Ein *Parameter Forest* hat **keine** einzelne Wurzel, sondern enthält mehrere *Parameter Trees*
  und/oder *Atomic Parameter Lists*. Durch Einführen einer Wurzel wird er wieder zum *Parameter Tree*;
  zur Laufzeit leistet das ein *Splitter* (Hohpe/Woolf).
- *Pagination* baut auf allen vieren auf: Der Request ist typischerweise eine *Atomic Parameter List*,
  die Antwort ein *Parameter Tree* oder *Parameter Forest*. Varianten: *Offset-*, *Cursor-/Token-*
  und *Time-Based*.

**Bezüge nach außen.** Die vier Struktur-Patterns verfeinern *Command/Document/Event Message* aus den
EIP, indem sie in die Nutzlast hineinsehen; *Content Enricher* und *Content Filter* operieren auf
*Parameter Trees*. *Incremental State Build-up* (Conversation Patterns) hat die inverse Intention zu
*Pagination*; Daigneaus *Single Message Argument* ist der nächste Verwandte des *Parameter Tree*.

**Namenswandel.** Tabelle I (S. 31 f.) zeigt den Stand von 2017; am ehesten entsprechen sich
*Vertical/Horizontal Integration* → [Frontend](../foundation/FrontendIntegration.md)/[Backend
Integration](../foundation/BackendIntegration.md), *Service Contract* →
[API Description](../foundation/APIDescription.md), *Master Data Resource* →
[Master Data Holder](../responsibility/MasterDataHolder.md), *Transactional Data Resource* →
[Operational Data Holder](../responsibility/OperationalDataHolder.md), *Query Service* →
[Retrieval Operation](../responsibility/RetrievalOperation.md), *Command Service* →
[State Transition Operation](../responsibility/StateTransitionOperation.md), *Validation Service* →
[Computation Function](../responsibility/ComputationFunction.md).

## Bemerkenswerte Aussagen

- Die Autoren geben unumwunden zu, dass *Pagination* das Problem nicht allein löst: „More than a
  single pattern is required to solve the pagination problem" (S. 30).
- Aussagen wie „immer feingranular vor grobgranular" nennt das Paper unzureichend oder sogar
  unverantwortlich, weil Projektkontexte differieren (S. 33).
- Generische Flexibilität wird als Kostenfaktor bewertet, nicht als Vorteil.
- Kopplung entsteht auch über *implizit* geteiltes Wissen: Selbst Provenance-Angaben in der
  Dokumentation erhöhen sie, sobald Empfänger sich darauf verlassen.

## Bezug zu Kubernetes / KRM

Die vier Struktur-Patterns sind für den KRM-Nachrichtenrumpf keine offene Wahl mehr: Jedes Objekt ist
ein *Parameter Tree* mit genau einer Wurzel und den Top-Level-Feldern `apiVersion`, `kind`,
`metadata`, `spec`, `status`; ein *Parameter Forest* kommt nicht vor. Damit liegt KRM exakt auf der
Fünf-bis-sieben-Elemente-Schwelle — die Komplexität ist in die Tiefe verlagert, nicht in die Breite.
Den Rat, tiefe Verschachtelung zu vermeiden, bricht KRM bewusst (`PodSpec`); das deckt sich mit der
Ausnahme, die das Paper zulässt, und der Konsument ist ohnehin ein Controller, der alles liest.

Die Warnung vor generischen Key-Value-Strukturen trifft KRM an zwei Stellen: `labels` und
`annotations` in `ObjectMeta` sind genau solche `map[string]string` ohne Schema, und
`runtime.RawExtension` bzw. `unstructured.Unstructured` erlauben schemalose Teilbäume. Auch DTO vs.
DTR hat eine Entsprechung: Die versionierten externen Typen sind die DTRs, die internen Hub-Typen
gehen nie über die Leitung. Und *Atomic Parameter Tree* beschreibt, was `metav1.ListOptions` tut:
flache Query-Skalare, serverseitig in einen versionierten Struct überführt.

Bei *Pagination* liest sich das Paper wie eine Begründung der KRM-Entscheidung: Es hält fest, dass
Offset-basierte Paginierung bei veränderlichen Daten Einträge doppelt zeigen oder auslassen kann und
dass dann ein zeit- oder tokenbasiertes Verfahren zu wählen ist. Genau das tut KRM — `limit` plus
opakes `continue`-Token, kein Offset. Die geforderte Korrelation der Teilantworten entsteht dadurch,
dass das Token die `resourceVersion` des ersten Requests mitführt und alle Chunks aus demselben
Snapshot bedient werden; auf die Frage nach Speicherdauer und Löschpolitik antwortet KRM mit der
etcd-Kompaktierung, nach der Weiterblättern mit `410 Gone` und `Status.Reason: "Expired"` scheitert.
Und die negative Konsequenz — Funktionen über den Gesamtbestand funktionieren nicht gut — erklärt,
warum es keine serverseitige Volltextsuche gibt und warum `remainingItemCount` bei gesetzten
Selektoren fehlt.

## Verwandte Wiki-Seiten

[Atomic Parameter](../structure/AtomicParameter.md) · [Atomic Parameter List](../structure/AtomicParameterList.md) · [Parameter Tree](../structure/ParameterTree.md) · [Parameter Forest](../structure/ParameterForest.md) · [Pagination](../quality/Pagination.md) ·
[Data Element](../structure/DataElement.md) · [Metadata Element](../structure/MetadataElement.md) · [Id Element](../structure/IdElement.md) · [Link Element](../structure/LinkElement.md) · [API Key](../structure/APIKey.md) ·
[Error Report](../structure/ErrorReport.md) · [Context Representation](../structure/ContextRepresentation.md) · [Wish List](../quality/WishList.md) · [Embedded Entity](../quality/EmbeddedEntity.md) · [Linked Information Holder](../quality/LinkedInformationHolder.md) · [Kategorie Structure](../meta/category-structure.md)

---
[← Index](../README.md) · [Papers](../papers/) · [Quelle](http://eprints.cs.univie.ac.at/5161/1/WADE-EuroPlop2017Paper-FinalSubmissionOct19.pdf)
