---
title: "Data-Oriented Interface Responsibility Patterns: Types of Information Holder Resources"
kategorie: Paper
autoren: Olaf Zimmermann, Cesare Pautasso, Daniel Lübke, Uwe Zdun, Mirko Stocker
venue: EuroPLoP 2020
quelle: http://eprints.cs.univie.ac.at/6521/1/MAP-EuroPlop2020bPaper.pdf
---

# Data-Oriented Interface Responsibility Patterns: Types of Information Holder Resources

## Bibliografische Angaben

Olaf Zimmermann (HSR Rapperswil), Cesare Pautasso (USI Lugano), Daniel Lübke (iQuest, Hannover),
Uwe Zdun (Universität Wien), Mirko Stocker (HSR Rapperswil): *Data-Oriented Interface Responsibility
Patterns: Types of Information Holder Resources.* EuroPLoP '20, 1.–4. Juli 2020, virtuell,
Deutschland. ACM, 25 Seiten. DOI [10.1145/3424771.3424821](https://doi.org/10.1145/3424771.3424821).
Das Schwesterpapier *Interface Responsibility Patterns: Processing Resources and Operation
Responsibilities* behandelt die aktionsorientierten Endpunkte und die Operationsverantwortlichkeiten.

## Worum es geht

MAP unterscheidet zwei architektonische Rollen für API-Endpunkte: *Processing Resource* (nimmt
Kommandos entgegen) und *Information Holder Resource* (exponiert Speicherung und Verwaltung von
Daten). Das Paper führt die generische *Information Holder Resource* ein und verfeinert sie in fünf
Spezialisierungen: drei unterscheiden sich nur in **Lebensdauer, Änderbarkeit und Referenzstruktur**
der Daten, zwei haben einen Sonderzweck — lose gekoppelten Datenaustausch und Adressverwaltung.
Begriffliche Grundlage ist Responsibility-Driven Design, wo „information holder“ ein
Rollen-Stereotyp ist: wer Information kennt und liefert.

## Behandelte Patterns

| Pattern | Paper-Abschnitt | Wiki-Link |
|---|---|---|
| *Information Holder Resource* | 4.1 (S. 3–9) | [InformationHolderResource](../responsibility/InformationHolderResource.md) |
| *Operational Data Holder* | 4.2 (S. 9–11) | [OperationalDataHolder](../responsibility/OperationalDataHolder.md) |
| *Master Data Holder* | 4.3 (S. 11–14) | [MasterDataHolder](../responsibility/MasterDataHolder.md) |
| *Reference Data Holder* | 4.4 (S. 14–16) | [ReferenceDataHolder](../responsibility/ReferenceDataHolder.md) |
| *Data Transfer Resource* | 4.5 (S. 17–20) | [DataTransferResource](../responsibility/DataTransferResource.md) |
| *Link Lookup Resource* | 4.6 (S. 21–23) | [LinkLookupResource](../responsibility/LinkLookupResource.md) |

## Was das Paper gegenüber der Website ergänzt

Die Website kürzt die *Forces* auf Stichworte und ersetzt die *Consequences* durch einen
Buchverweis. Das Paper führt beides aus und hat je eine Rubrik *Non-solution*.

### Information Holder Resource

**Forces.** *Modellierungsansatz und Kopplungswirkung*: Ein datenzentrischer Schnitt erzeugt viele
CRUD-APIs; das schadet dem Abhängigkeitsmanagement (Information Hiding verletzt) **und** der
Datenqualität, weil jeder autorisierte Client beliebig manipulieren kann — CRUD erzeugt operative
*und* semantische Kopplung. *Qualitätskonflikte*: Aktualität gegen Konsistenzaufwand, Entwurfszeit-
gegen Laufzeit- gegen Evolutionsqualitäten, Sicherheit quer dazu; publizieren heißt zu entscheiden,
wer lesen und schreiben darf und was mit den Konsumenten passiert, wenn die Daten verschwinden.
*Architekturprinzipien*: lose Kopplung, Datenunabhängigkeit, unabhängige Deploybarkeit.

**Non-solutions.** Alles hinter Operationen und DTOs zu verstecken fördert Information Hiding,
begrenzt aber unabhängiges Deployen und Skalieren — es folgen geschwätzige Aufrufe oder redundante
Datenhaltung. Direkten Datenbankzugriff durchzureichen macht das Schema unantastbar.

**Consequences.**

- **+** Der Ansatz passt, wenn das Szenario datenzentrisch *ist*; Aktivitätsorientierung ist häufig
  vorzuziehen, aber nicht überall natürlich.
- **+** Die Verarbeitung wandert zum Konsumenten; der Endpunkt wird Quelle verlinkter Daten,
  Beziehungssenke (*Operational Data Holder*) oder beides (*Data Transfer Resource*).
- **−** Sicherheit, Datenschutz, Konsistenz, Verfügbarkeit und Kopplung sind einzeln abzuwägen;
  Konsumenten federn Ausfälle selbst ab, und *jede* Änderung an Inhalt, Metadaten und Format ist
  kontrollpflichtig. Dazu der Ruf, Kopplung zu erhöhen und Information Hiding zu verletzen.
- Gegen Nygards „entity service anti-pattern“ positionieren sich die Autoren ausdrücklich: immer
  wegzuentwickeln gehe zu weit — jede Verwendung müsse aber eine bewusste Entscheidung sein.

### Operational Data Holder

**Forces.** *Verarbeitungsgeschwindigkeit* — extrem niedrige Antwortzeiten für Lesen *und*
Schreiben. *Fachliche Agilität* — die Änderbarkeit muss bis auf die Schemaebene reichen; das Paper
nennt A/B-Tests mit einem Teil der Live-Nutzer. *Konzeptionelle Integrität der ausgehenden
Beziehungen* — hohe Genauigkeitsstandards bei prüfungsrelevanten Daten, referenzierte Entitäten oft
beim Integrationspartner, und Konsumenten erwarten sie danach korrekt erreichbar. Die
*Non-solution*, alles gleich zu behandeln, ergibt überkonstruierte Konsistenzverwaltung.

**Consequences.** Das Paper stuft das Pattern offen als Marker für die API-Dokumentation ein.

- **+** Je weniger eingehende Abhängigkeiten, desto leichter die Änderung; die begrenzte Lebensdauer
  hilft, bei der Evolution rückwärtskompatibel zu bleiben. Gelockerte Konsistenz erhöht die
  Verfügbarkeit, ein zustandsloser Endpunkt die horizontale Skalierbarkeit.
- **−** Konsistenz und Verfügbarkeit werden hier oft *anders* priorisiert als bei Stammdaten;
  eventual consistency kann die bessere Wahl sein. Mit ehrlichem Nachsatz: auch bei richtig
  gewähltem Pattern kann die Implementierung Performance und Verfügbarkeit ruinieren.

### Master Data Holder

**Forces.** *Stammdatenqualität* — die Kraft hat zwei Enden: liegen die Daten nicht an *einem* Ort,
führen unkoordinierte Updates zu schwer auffindbaren Inkonsistenzen; liegen sie zentral, wird der
Zugriff durch Contention langsam. *Stammdatenschutz* — ein attraktives Angriffsziel mit schweren
Folgen. *Daten unter externer Kontrolle* — Stammdaten gehören häufig einem MDM-System einer anderen
Organisationseinheit, haben Bilanzwert, unterliegen fremden Auditregeln und entwickeln sich **in
anderem Tempo** als die operativen Daten, die auf sie verweisen. Alles gleich zu behandeln (die
*Non-solution*) scheitert an Auditoren, Data Owners und den realen Personen hinter den Daten.

**Consequences.**

- **+** Die Kennzeichnung erzeugt den nötigen Fokus auf Datenqualität und Datenschutz.
- **−** Stammdaten haben viele eingehende und oft auch ausgehende Abhängigkeiten; für deren
  Konsistenz und Aktualität braucht es *weitere* Patterns.
- Die schärfste Selbstkritik der Sammlung: die Kennzeichnung allein löst **keine einzige** Kraft
  auf; erst die Implementierungshinweise tun das.
- Löschen ist wegen der vielen eingehenden Referenzen riskant, teils rechtlich verboten, teils
  geboten; üblich ist ein unveränderlicher Archivzustand statt physischer Löschung.

### Reference Data Holder

**Forces.** *Performance vs. Konsistenz beim Lesen* — Caching lohnt sich, muss aber so entworfen
sein, dass Caches nicht zu groß werden, Replikation Netzpartitionen übersteht und seltene Änderungen
dennoch konsistent ankommen. *DRY* — Daten hart zu verdrahten oder einmal zu holen und ewig lokal zu
halten funktioniert prächtig, bis sie sich ändern; sind die Clients dann außer Reichweite, ist die
Änderung nicht mehr durchsetzbar (zweistellige Jahreszahlen, Euro, neue Postleitzahlen). Die
*Non-solution* — statische wie dynamische Daten behandeln — verschenkt die Leseoptimierung.

**Consequences.**

- **+** DRY ist aufgelöst: ein zentraler Bezugspunkt verbreitet die Daten und behält zugleich die
  Kontrolle über sie.
- **+** Der Lesezugriff ist optimierbar, und unveränderliche Daten replizieren ohne
  Inkonsistenzrisiko — der Provider rüstet Proxies, Caches und Replikate **unsichtbar** nach.
- **−** Ein eigener Endpunkt muss entwickelt, dokumentiert, betrieben und gepflegt werden. Bringt er
  mehr Aufwand als Nutzen, ist das Refactoring benannt: die statischen Daten in einen bestehenden
  *Master Data Holder* hineinziehen. Und dies ist einer der wenigen Fälle für eine
  [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md).

### Data Transfer Resource

**Forces** — jede einzeln begründet: *zeitliche Kopplung* (je mehr Teilnehmer, desto
unwahrscheinlicher, dass alle gleichzeitig bereit sind); *örtliche Kopplung* (hinter NAT oder
Firewall nicht adressierbar); *Kommunikationsbeschränkungen* (Clients nehmen keine eingehenden
Verbindungen an, Middleware darf oft nicht installiert werden); *Zuverlässigkeit*; *Skalierbarkeit*
(Empfängerzahl beim Senden unbekannt, Datenmenge kann die Kapazität einzelner Nachrichten sprengen);
*Speicherplatzeffizienz*; *Latenz*; *Eigentümerschaft* — beim Aufräumen wollen drei Parteien
Unterschiedliches: der Sender maximale Reichweite, der Empfänger mehrfaches Lesen, der Betreiber
niedrige Speicherkosten. *Non-solution* wäre Messaging-Middleware, die lokal beim Client liefe.

**Consequences.**

- **+** Zeitliche und örtliche Entkopplung; Clients ohne direkte Verbindung nutzen die Ressource als
  Blackboard; Zuverlässigkeit durch idempotente Übertragung; Skalierbarkeit in Datenmenge *und*
  Client-Zahl; die Flexibilität der Eigentümerschaft hängt von der Variante ab.
- **−** Clients können keine Benachrichtigungen empfangen und müssen pollen; der Provider muss
  Speicher vorhalten; zwei Hops statt einem. Den Tausch verteidigt das Paper explizit: über große
  Zeiträume und viele Teilnehmer übertragen zu können, hat Vorrang vor der Einzelperformance.
- Danach folgen weitere Entwurfsfragen: Zugriffskontrolle, fehlende Koordination zwischen Lesern und
  Schreibern, optimistisches Sperren, Polling und Garbage Collection.

**Drei benannte Varianten**, die die Website nicht führt: *Relay Resource* (ein Schreiber, ein
Leser; die Eigentümerschaft wandert mit), *Published Resource* (ein Schreiber, unvorhersehbar viele
Leser, teils Jahre später) und *Conversation Resource* (alle lesen, schreiben und löschen).

### Link Lookup Resource

**Forces.** *Kohäsion und Kopplung* — Lookups in einen ohnehin reichen Endpunkt zu legen, macht ihn
schwer zu dokumentieren, zu warten und zu testen; das Single-Responsibility-Prinzip ist verletzt.
*Dynamische Endpunktreferenzen* — Bindung zur Entwurfs- oder Deploymentzeit reicht nicht, wenn
Endpunkte für Wartung verschwinden oder Vermittler nach einem Versionswechsel umleiten. *Anzahl der
Endpunkte* — im Extremfall verdoppelt ein Lookup je Ressource die Endpunktzahl. *Zentralisierung* —
jede zentrale Lösung zieht mehr Verkehr auf sich. *Nachrichtengrößen und Aufrufzahl* — Einbetten
vermeidet Referenzen, vergrößert aber Nachrichten. *Tote Links* — Konsumenten scheitern oder lesen
veraltete Daten. Die *Non-solution*, Lookups anzuhängen, zerstört die Kohäsion der Endpunkte.

**Consequences.**

- **+** Entkopplung in der Ortsautonomie; hohe Kohäsion, weil die Auflösungsverantwortung von
  Verarbeitung und Datenabruf getrennt ist.
- **−** Zusätzliche Aufrufe und mehr Endpunkte; laufende Betriebskosten, denn die Lookup-Ressource
  muss aktuell gehalten werden.
- Die Rechnung geht nur auf, wenn der Zusatzaufruf billiger ist als die Ersparnis an
  Nachrichtenlast. Sonst ist der Rückbau vorgezeichnet: erst den Lookup durch einen direkten Link
  ersetzen, und wenn der Dialog immer noch zu geschwätzig ist, die Daten einbetten.
- Als **Variante** nennt das Paper HATEOAS: zeigen die Links auf *Processing Resources*, werden
  Kontrollfluss und Anwendungszustand dynamisch und dezentral.

## Pattern-Sprache und Zusammenhänge

Abbildung 1 des Papers ist ein Entscheidungsbaum aus drei Fragen: **Zustand lesen oder schreiben?**
(nein → *Computation Function*; ja → aktivitäts- oder datenorientierte Endpunktrolle), **vom Client
änderbar?** (nein → *Reference Data Holder*) und **Lebensdauer?** (lang → *Master Data Holder*, kurz
→ *Operational Data Holder*). Quer dazu liegen *Data Transfer Resource* („Infrastruktur für Daten“)
und *Link Lookup Resource* („Realisierungsstrategie“) — sie klassifizieren den Zweck, nicht die
Daten.

| Typ | Lebensdauer | Änderbarkeit durch Clients | Referenzrichtung | Eigentum |
|---|---|---|---|---|
| *Operational Data Holder* | kurz bis mittel | häufig | viele ausgehende | Provider |
| *Master Data Holder* | lang | selten, aber möglich | viele eingehende | oft externes MDM-System |
| *Reference Data Holder* | sehr lang | keine (nur administrativ) | nur eingehende | Provider/Standardgeber |
| *Data Transfer Resource* | temporär | vollständig | keine | Clients |
| *Link Lookup Resource* | — | selten | hält nur Adressen | Provider |

Ob analytische oder Monitoring-Daten eigene Patterns bräuchten, verneint das Paper: analytische Daten
seien ein Sonderfall der Referenzdaten, Monitoring-Daten operativ. Wer sich mit der Einteilung
schwertut, darf von der generischen *Information Holder Resource* sprechen — für eine
Pattern-Sprache eine bemerkenswert entspannte Auskunft.

## Bemerkenswerte Aussagen

- Der *Operational Data Holder* diene vor allem als „marker pattern in API documentation“ (S. 10);
  beim *Master Data Holder* schärfer: die Kennzeichnung allein löse keine Kraft auf.
- Die Trennung von Stamm- und operativen Daten sei subjektiv und kontextabhängig — dieselbe
  Bestellung ist für den Käufer flüchtig, für den Shop-Betreiber ein dauerhafter Vermögenswert.

## Bezug zu Kubernetes / KRM

Die MAP-Dreiteilung Operational / Master / Reference Data existiert in KRM **nicht als
Modellkonzept**. KRM zieht dieselbe Achse *quer* dazu — nicht zwischen Ressourcentypen, sondern
durch **jedes** Objekt hindurch: `spec` gehört dem Nutzer und ist stammdatenartig, `status` gehört
dem Controller und ist operativ, durchgesetzt über die `/status`-Subresource mit eigenen
RBAC-Verben. Wo MAP zwei Endpunkttypen trennt, hat KRM zwei Hälften pro Objekt. Faktisch findet man
die drei Typen trotzdem wieder:

- **Operativ**: `Lease` (Node-Heartbeats und Leader Election, im Sekundentakt fortgeschrieben) und
  `Event`, das der API-Server nach `--event-ttl` wegräumt (Standardwert eine Stunde, gesetzt in
  `pkg/controlplane/apiserver/options/options.go`) — beide kurzlebig und nur *auf* langlebigere
  Objekte verweisend, die Referenzrichtung des Patterns.
- **Stammdaten**: `Node` und `Deployment` — administrativ gepflegt, selten geändert, von vielen
  Objekten per Name referenziert. Dass die hochfrequenten Heartbeats aus dem `Node`-Objekt in ein
  separates `Lease` wanderten, ist die Trennung von Stamm- und operativen Daten in Reinform.
- **Referenzdaten**: die Klassenobjekte `StorageClass`, `IngressClass` und `PriorityClass` sowie die
  ConfigMap `kube-root-ca.crt`, die ein Controller (`RootCACertConfigMapName` in
  `pkg/controller/certificates/rootcacertpublisher`) in jeden Namespace verteilt — der vom Pattern
  geforderte „Vollabruf und lokal kopieren“-Pfad, nur push- statt pull-basiert.

Die reinste Umsetzung eines Patterns aus diesem Paper ist die **Discovery-API** als *Link Lookup
Resource*: `/api` und `/apis` liefern den Katalog der Gruppen und Versionen,
`apidiscovery.k8s.io/v2` liefert ihn aggregiert in einem Roundtrip, und der `RESTMapper` in
apimachinery ist die Client-Seite davon — ohne ihn löst `kubectl` keinen Ressourcennamen auf. Der
charakteristische Unterschied: KRM löst **per Name statt per Link** auf. Zurückgegeben werden Namen
und `ObjectReference`s, nie URLs; die Adresse konstruiert der Client aus dem uniformen Pfadschema
selbst — das Force „Anzahl der Aufrufe“ ist einmal pro Sitzung statt pro Link bezahlt.

## Verwandte Wiki-Seiten

- [Information Holder Resource](../responsibility/InformationHolderResource.md) — das Oberpattern,
  verfeinert durch [Operational Data Holder](../responsibility/OperationalDataHolder.md),
  [Master Data Holder](../responsibility/MasterDataHolder.md),
  [Reference Data Holder](../responsibility/ReferenceDataHolder.md),
  [Data Transfer Resource](../responsibility/DataTransferResource.md) und
  [Link Lookup Resource](../responsibility/LinkLookupResource.md).
- [Embedded Entity](../quality/EmbeddedEntity.md) vs.
  [Linked Information Holder](../quality/LinkedInformationHolder.md), mit dem
  [Link Element](../structure/LinkElement.md) als Baustein — die Referenzentscheidung.
- [Kategorie Responsibility](../meta/category-responsibility.md) — Einordnung der Kategorie.

---
[← Index](../README.md) · [Papers](../papers/) · [Quelle](http://eprints.cs.univie.ac.at/6521/1/MAP-EuroPlop2020bPaper.pdf)
