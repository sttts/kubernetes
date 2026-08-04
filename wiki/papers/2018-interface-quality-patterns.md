---
title: "Interface Quality Patterns — Communicating and Improving the Quality of Microservices APIs"
kategorie: Paper
autoren: Mirko Stocker, Olaf Zimmermann, Uwe Zdun, Daniel Lübke, Cesare Pautasso
venue: EuroPLoP '18, 4.–8. Juli 2018, Irsee (ACM, 16 Seiten, DOI 10.1145/3282308.3282319)
quelle: http://eprints.cs.univie.ac.at/5661/1/Interface%20Quality%20Patterns%20-%20Communicating%20and%20Improving%20the%20Quality%20of%20Microsevices%20APIs.pdf
---

# Interface Quality Patterns — Communicating and Improving the Quality of Microservices APIs

## Bibliografische Angaben

- **Autoren:** Mirko Stocker, Olaf Zimmermann (beide HSR Rapperswil), Uwe Zdun (Universität Wien),
  Daniel Lübke (iQuest, Hannover), Cesare Pautasso (USI Lugano)
- **Konferenz:** 23rd EuroPLoP, Irsee; 16 Seiten, fünf ausformulierte Patterns
- **Vorgängerarbeit:** *Interface Representation Patterns* (EuroPLoP '17), durchgehend referenziert

## Worum es geht

Das Paper ist die zweite Tranche der Pattern-Sprache, aus der später
[Microservice API Patterns](https://microservice-api-patterns.org) und das gleichnamige Buch wurden.
Während der Vorgänger fragte, *wie* Nachrichtenparameter strukturiert werden, fragt dieser Teil, wie
sich **beobachtbare Qualitätseigenschaften** einer API erreichen und — genauso wichtig —
**kommunizieren** lassen. Alle fünf Patterns beantworten dieselbe Designfrage: Wie erreicht man ein
bestimmtes Qualitätsniveau, ohne die eigenen Ressourcen unwirtschaftlich zu verbrauchen? Gewonnen
wurden sie aus 31 untersuchten Web-APIs. Die ersten vier Patterns adressieren Architekten und
Entwickler, die letzten beiden ausdrücklich API-Product-Owner — Qualität ist hier auch eine
Geschäftsfrage.

## Behandelte Patterns

| Pattern | Paper-Abschnitt | Wiki-Link |
| --- | --- | --- |
| *API Key* | 4.1 (S. 4–6) | [API Key](../structure/APIKey.md) |
| *Wish List* | 4.2 (S. 6–8) | [Wish List](../quality/WishList.md) |
| *Rate Limit* | 4.3 (S. 8–10) | [Rate Limit](../quality/RateLimit.md) |
| *Rate Plan* (heute *Pricing Plan*) | 4.4 (S. 11–12) | [Pricing Plan](../quality/PricingPlan.md) |
| *Service Level Agreement* | 4.5 (S. 12–15) | [Service Level Agreement](../quality/ServiceLevelAgreement.md) |

Referenziert, aber nicht eingeführt: *Atomic Parameter*, *Atomic Parameter List*, *Parameter Tree*,
*Parameter Forest*, [Pagination](../quality/Pagination.md) — Tabelle 1 fasst sie zusammen.

## Was das Paper gegenüber der Website ergänzt

Seit dem Buch kürzt microservice-api-patterns.org *Forces* und *Consequences* auf Stichworte oder
einen Buchverweis; das Paper enthält sie ausformuliert. Auffällig ist die Form: Forces durchgehend
als **Fragen**, Consequences als **Resolution of forces** mit `+`/`−`-Zeilen — ein Kraftfeld, das
ausdrücklich nicht restlos aufgelöst wird.

### API Key (S. 4–5)

**Forces** — zwei Gruppen. Die funktionalen: Wie identifizieren sich Client-*Programme* an einem
Endpoint, ohne Benutzerkonto-Credentials speichern und übertragen zu müssen? Wie lässt sich der
aufrufende Client von der Organisation des Kunden **entkoppeln**? Wie realisiert man *abgestufte*
Authentifizierungsniveaus je nach Sicherheitskritikalität? Dann die Zielkonflikte: identifizieren
und die API trotzdem einfach benutzbar halten; absichern und den Performance-Einfluss minimieren.
**Non-solution:** Das volle CIA-Portfolio ist für eine freie öffentliche API unwirtschaftlich; VPN
oder Two-Way-SSL verhindern Anwendungsfälle wie das Durchsetzen von *Rate Limits*.

**Consequences.** Ein *API Key* ist die leichtgewichtige Alternative zu einem vollen
Authentifizierungsprotokoll und tariert Basissicherheit gegen Overhead aus.

- `+` Als Shared Secret erlaubt er Identifikation und darauf aufbauend Authentifizierung und
  Autorisierung; wegen seiner geringen Größe kostet er pro Request kaum Performance.
- `+` Die Trennung vom Kundenkonto **entkoppelt Kundenrollen** (Administration, Business, Nutzung).
  Kunden können mehrere Keys mit unterschiedlichen Rechten pro Implementierung oder Standort
  anlegen; bei einem Leak wird ein einzelner Key unabhängig vom Konto widerrufen. Analytics und
  *Rate Limits* werden pro Key möglich.
- `−` Er wird bei jedem Request mitgeschickt, gehört also zwingend auf eine sichere Verbindung; die
  Alternativen (VPN, Public-Key-Kryptografie) kosten Konfigurationsaufwand und Performance.
- `−` Er ist **nur ein Bezeichner** und trägt keine Nutzlast — kein Ablaufzeitpunkt, keine
  Berechtigungen. Genau hier setzt die diskutierte JWT-Alternative an.

### Wish List (S. 6–7)

**Forces.** Wie befriedigt ein Provider die widersprüchlichen Informationsbedarfe einzelner Clients,
**ohne pro Client eigene Endpoints** zu bauen, und vermeidet trotzdem Under- und Over-Fetching? Wie
kann ein Client provider-seitige Selektionsfilter *spezifizieren und über sie lernen*? Und — die oft
übersehene Force — wie kommt der Provider mit der Komplexität zurecht, die maßgeschneiderte
Antworten für Sicherheit, Test und Wartung bedeuten? Dazu Antwortzeit, Durchsatz, Verarbeitungszeit.
**Non-solution:** Gateways und Caches senken die Last auch, verkomplizieren aber die Topologie.

**Consequences.**

- `+` Der Client drückt seine Wünsche durch An- und Abwählen von Attributen aus; damit ist
  Datensparsamkeit erfüllt.
- `+` Der Provider muss keine spezialisierten Operationsvarianten pflegen und **nicht raten**,
  welche Daten ein Use Case braucht; weniger Datenbank- und Netzlast.
- `−` Mehr Logik in der Service-Schicht, die bis in den Datenzugriff durchschlagen kann.
- `−` Der Provider **exponiert sein Datenmodell** und erhöht damit die Kopplung.
- `−` Die Liste selbst ist Aufwand: erzeugen, transportieren, verarbeiten.

Fallstrick aus der *Further discussion*: Eine kommaseparierte Attributliste ist nicht typsicher. Ein
Tippfehler führt bestenfalls zu einem Fehler, schlimmstenfalls wird der Wunsch **stillschweigend
ignoriert** — der Client glaubt dann, das Attribut existiere nicht. Umbenennungen wirken genauso.

### Rate Limit (S. 8–10)

**Forces.** Wie hält der Provider hohe Performance für *alle* Clients und wirtschaftet trotzdem mit
seinen Ressourcen? Wie verhindert er Missbrauch bzw. minimiert dessen Auswirkung — explizit unter
der Überschrift **Fairness**? Wie bleibt der Dienst zuverlässig und kosteneffizient, ohne einzelne
Clients übermäßig einzuschränken? Und, häufig vergessen: Wie kontrolliert der **Client** seinen
eigenen Verbrauch, wenn er Kapazität sparen muss? Was „exzessiv“ heißt, definiert laut Fußnote 29
der Provider. **Non-solution:** mehr Rechenleistung und Bandbreite kaufen; selten wirtschaftlich.

**Consequences.**

- `+` Schutz vor bösartigen Clients und unerwünschten Bots, Erhalt der Servicequalität.
- `+` Bessere Kapazitätsplanung durch **gedeckelte Maximalnutzung** — besser für alle.
- `−` Die Kalibrierung ist das eigentliche Problem: zu hoch wirkungslos, zu niedrig ärgerlich.
  Beispiel des Papers: 30.000 Requests/Monat, die ein Client ohne Zusatzschranke in einem einzigen
  Burst verfeuert — Gegenmittel ist eine zweite, kurzfristige Grenze (ein Request/Sekunde).
- `−` **Clients** müssen ihren Verbrauch verfolgen und den Limitfall behandeln (Tracing, Queuing,
  Caching, Priorisierung von Aufrufen).

Zwei Warnungen aus den Implementation Hints: Der Limit-Zustand gehört in eine Client-Datenbank und
**nicht in eine Session**, sonst umgeht der Client ihn durch eine neue Session; und die
Metering-Infrastruktur darf nicht mehr kosten als die abgewiesenen Aufrufe. Dazu die Wechselwirkung
mit *Wish List*: Wer Feldselektion oder GraphQL anbietet, kann nicht mehr fair nach Requests zählen.

### Rate Plan / Pricing Plan (S. 11–12)

**Forces.** Wie wählt ein Provider ein Preismodell, das seine wirtschaftlichen Interessen mit denen
der Kunden **und dem Wettbewerb** ausbalanciert? Wie feingranular muss gemessen werden, um den
Informationsbedarf der Kunden zu decken, ohne Performance- oder Verfügbarkeitseinbußen? Und wie
garantiert man die **Sicherheit der Messdaten** bei gleichzeitiger Abrechnungsgenauigkeit und
Auditierbarkeit? **Non-solution:** Eine pauschale Anmeldegebühr behandelt Hobbyisten und Großkunden
gleich — für die einen zu teuer, für die anderen zu billig.

**Consequences.** Das Paper ist hier ungewöhnlich zurückhaltend: Ein *Rate Plan* löst die meisten
Forces, **Sicherheitsgarantien muss aber die darunterliegende Implementierung liefern**.

- `+` Kunde und Provider haben eine klare Vereinbarung über anfallende Kosten und Pflichten.
- `−` Sinnvolle Preispläne zu schreiben ist schwer und verlangt viel Wissen über Interessen und
  Geschäftsmodelle **beider** Seiten.
- `−` Clients müssen über einen *API Key* oder ein anderes Verfahren identifizierbar sein.
- `−` Nutzungsbasierte Preise verlangen detailliertes Monitoring und — zur Streitvermeidung —
  Reporting an den Kunden; beides Aufwand beim Provider.

Selten explizit gestellt, hier schon: Was tun beim **Ausfall der Messfunktion**? Ohne Metering ist
keine spätere Abrechnung möglich — also API abschalten oder den Dienst so lange verschenken.

### Service Level Agreement (S. 13–14)

**Forces** — fünf Stück, die weit über Technik hinausgehen. Wie entscheidet ein Client, ob das
Angebot zu seinen Geschäftsanforderungen passt (Agilität, Überlebensfähigkeit des Anbieters)? Wie
erfährt er etwas über Regulierung, Sicherheits- und Datenschutzmaßnahmen? Wie kommuniziert ein
Provider Attraktivität, Verfügbarkeit und Performance **ohne unrealistische Versprechen**, die zu
Unzufriedenheit oder Verlusten führen? Wie balanciert er Ressourcenökonomie gegen Gewinn? Und was
ist der richtige Detailgrad — Unterspezifikation erzeugt Spannungen, Überspezifikation Aufwand?
**Non-solution:** dem Provider vertrauen — bei kritischer Nutzung untragbar; Freitext ist mehrdeutig.

**Consequences.**

- `+` Gemeinsames Verständnis über erwartbare Service- und Qualitätsniveaus.
- `+` Ein SLA kann alle Services oder gezielt einzelne Operationen adressieren: Datenschutz-SLOs
  gehören ins Gesamt-SLA, Data-Management-Ziele wie Backup-Frequenz können pro Endpoint variieren.
- `+` Gut gemachte SLAs mit messbaren SLOs sind ein **Reifeindikator**. Bemerkenswerterweise hatten
  damals viele öffentliche APIs und Cloud-Angebote gar keine oder nur schwache SLAs — zugeschrieben
  der Marktdynamik und fehlender Regulierung.
- `−` Der Provider **wird haftbar**; manche Organisationen wollen für ihre Fehler nicht einstehen,
  ein SLA kann daher auf internen Widerstand stoßen.

Zwei Nuancen: Ein dauerhaft **übererfülltes** SLA erzieht Kunden zu höheren Erwartungen — daraus der
von Google übernommene Gedanke, Verfügbarkeit nicht wesentlich besser als das SLO zu betreiben,
notfalls durch bewusste Ausfälle. Und es gibt das **interne SLA**: Der Provider misst sich selbst.

## Pattern-Sprache und Zusammenhänge

Abbildung 1 des Papers zeigt, warum die fünf Patterns zusammen ein Paket bilden:

- *Rate Limit* und *Rate Plan* **identifizieren den Client über** einen *API Key* — ohne
  Zurechenbarkeit gibt es weder Kontingent noch Rechnung.
- *Rate Plan* **nutzt** *Rate Limits* zur Durchsetzung von Abrechnungsstufen und **verweist auf**
  das *Service Level Agreement*, in dem beider Details stehen.
- *Wish List* **beeinflusst** *Rate Limit* — weniger Daten pro Aufruf, aber eine schwierigere
  Bemessungsgrundlage. Angewandt wird sie auf *Parameter Tree* und *Parameter Forest*; der
  *API Key* wird als *Atomic Parameter* repräsentiert.

Die Abgrenzung zu benachbarten Sprachen ist explizit: *Release It!* und Hanmers *Patterns for Fault
Tolerant Software* adressieren dieselben Qualitäten in der **internen Architektur**; diese Sprache
beschreibt Eigenschaften der **API-Beschreibung**.

## Bemerkenswerte Aussagen

Zur *Wish List*: „the desire for ‚Datensparsamkeit‘ (i.e., data parsimony) is met“ (S. 7) —
Datensparsamkeit erscheint als Qualitätsziel, nicht nur als Rechtsbegriff, und der deutsche Terminus
steht bewusst im englischen Text.

Sinngemäß zwei weitere Punkte: Clients umgehen ein *Rate Limit* durch Mini-Batch-Requests direkt
nach Fensteröffnung und erzeugen so gerade die Lastspitzen, die es verhindern sollte (beobachtet am
kostenlosen GitHub-Dienst). Und SLOs sollten **maschinenlesbar** vorliegen, damit Autoscaling
reagieren kann, bevor das SLA verletzt ist.

## Bezug zu Kubernetes / KRM

Kubernetes ist ein guter Gegentest: Es kennt fast alle Forces des Papers, löst mehrere davon aber
**anders auf** — und hat zwei Patterns schlicht nicht.

- ***Rate Limit*: Fairness statt Kontingent.** Die API Priority and Fairness (Gruppe
  `flowcontrol.apiserver.k8s.io`) limitiert mit `FlowSchema` und `PriorityLevelConfiguration`
  **Nebenläufigkeit (Seats)**, nicht Requests pro Zeitraum; ein „X Aufrufe pro Stunde“ existiert
  nicht. Überzählige Requests werden per Shuffle Sharding eingereiht statt abgewiesen, `429` fällt
  nur bei `limitResponse.type: Reject` oder Queue-Overflow an. Damit ist die Fairness-Force
  adressiert und die Kontingent-Force ausgelassen — passend dazu, dass die „Clients“ meist eigene
  Controller sind. Ein Teil der Begrenzung liegt bewusst **beim Client**: `rest.Config` trägt `QPS`
  und `Burst`.
- ***Conditional Request* gibt es nicht — `resourceVersion` und `watch` treten an seine Stelle.**
  Die API kennt kein ETag/`If-None-Match`; die „Ist es neu?“-Frage wird nicht pro GET gestellt,
  sondern einmal beim Aufbau eines `watch` ab einer `resourceVersion` — Polling wird überflüssig
  statt optimiert. Dieselbe `resourceVersion` trägt beim Schreiben die optimistische
  Nebenläufigkeit und übernimmt so die zweite ETag-Rolle.
- **Kein *Request Bundle*.** Es gibt keine Batch-Operation über mehrere Objekte: `kubectl apply` auf
  ein Manifest mit zwanzig Ressourcen erzeugt zwanzig Requests. Die einzigen Mengenoperationen,
  `list` und `deletecollection`, sind auf **eine** Ressourcenart begrenzt.
- **Kein Sparse Fieldset, also keine echte *Wish List*.** Statt eines `fields`-Parameters gibt es
  nur feste Projektionen per Content Negotiation (`as=PartialObjectMetadata`, `as=Table`);
  `fieldSelector` und `labelSelector` wählen **Zeilen, nicht Spalten**, und `kubectl -o jsonpath=...`
  ist client-seitig.
- ***ResourceQuota* ist kein *Pricing Plan*.** Sie begrenzt Objektzahlen und Compute-Ressourcen pro
  Namespace, kennt aber weder Metering zur Abrechnung noch Preise, Perioden oder Kunden — es fehlt
  genau die Dimension, die das Pattern ausmacht. Das Audit-Log dient der Nachvollziehbarkeit, nicht
  der Abrechnung.
- ***API Key* → gebundene Tokens.** Kubernetes wählt die Alternative gegen den Nachteil „ein Key
  trägt keine Nutzlast“: ServiceAccount-Tokens sind JWTs mit Audience, Ablaufzeitpunkt und
  Objektbindung aus der TokenRequest-API. Der statische Token-File-Authenticator — der reine
  *API Key* — ist der abgeratene Altfall.
- ***Service Level Agreement*: nur die interne Variante.** Upstream gibt keine
  Verfügbarkeitszusagen, SLAs liefern die Managed-Angebote. Was existiert, ist das interne SLA des
  Papers: SIG Scalability definiert SLIs/SLOs für API-Call-Latenzen und Pod-Startzeiten und misst sie
  in Tests — ohne Penalties, ohne Vertragspartner.

## Verwandte Wiki-Seiten

- [API Key](../structure/APIKey.md), [Wish List](../quality/WishList.md),
  [Rate Limit](../quality/RateLimit.md), [Pricing Plan](../quality/PricingPlan.md),
  [Service Level Agreement](../quality/ServiceLevelAgreement.md) — die fünf Patterns des Papers
- [Wish Template](../quality/WishTemplate.md) — im Buch aus der *Wish List* ausdifferenziert; hier
  erst Diskussionspunkt
- [Conditional Request](../quality/ConditionalRequest.md),
  [Request Bundle](../quality/RequestBundle.md), [Pagination](../quality/Pagination.md) — weitere
  Data-Transfer-Parsimony-Patterns, teils erst nach diesem Paper hinzugekommen
- [Embedded Entity](../quality/EmbeddedEntity.md),
  [Linked Information Holder](../quality/LinkedInformationHolder.md) — Designzeit-Antwort auf die
  Laufzeitfrage der *Wish List*
- [Error Report](../structure/ErrorReport.md) — Transportform für Limit- und Quota-Verletzungen
- [Kategorie Quality](../meta/category-quality.md) — Einordnung der Kategorie

---
[← Index](../README.md) · [Papers](../papers/) · [Quelle](http://eprints.cs.univie.ac.at/5661/1/Interface%20Quality%20Patterns%20-%20Communicating%20and%20Improving%20the%20Quality%20of%20Microsevices%20APIs.pdf)
