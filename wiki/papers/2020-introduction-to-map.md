---
title: "Introduction to Microservice API Patterns (MAP)"
kategorie: Paper
autoren: Olaf Zimmermann, Mirko Stocker, Daniel Lübke, Cesare Pautasso, Uwe Zdun
venue: "Joint Post-proceedings of the First and Second International Conference on Microservices (Microservices 2017/2019), OASIcs Vol. 78, S. 4:1–4:17, Schloss Dagstuhl 2020 (DOI 10.4230/OASIcs.Microservices.2017-2019.4)"
quelle: https://drops.dagstuhl.de/opus/volltexte/2020/11826/
---

# Introduction to Microservice API Patterns (MAP)

## Bibliografische Angaben

- **Autoren:** Olaf Zimmermann, Mirko Stocker (beide University of Applied Sciences of Eastern
  Switzerland, Rapperswil), Daniel Lübke (iQuest GmbH, Hannover), Cesare Pautasso (Software
  Institute, USI Lugano), Uwe Zdun (Universität Wien)
- **Erschienen in:** Joint Post-proceedings of the First and Second International Conference on
  Microservices (Microservices 2017/2019), OASIcs Band 78, S. 4:1–4:17, Schloss Dagstuhl –
  Leibniz-Zentrum für Informatik, 6.2.2020
- **DOI:** 10.4230/OASIcs.Microservices.2017-2019.4 · **URN:** urn:nbn:de:0030-drops-118268
- **Lizenz:** Creative Commons Attribution 3.0 Unported (Open Access)
- **Umfang:** 17 Seiten, vier Abbildungen, zwei ausformulierte Patterns
- **Grundlage dieser Zusammenfassung:** das vollständige Original-PDF von DROPS; punktuell
  ergänzt um Formulierungen auf microservice-api-patterns.org, wo das kenntlich gemacht ist.

## Worum es geht

Dies ist das **Überblickspaper** der MAP-Reihe. Anders als die EuroPLoP-Beiträge, die jeweils eine
Handvoll Patterns ausformulieren, beantwortet es, *warum* es diese Pattern-Sprache gibt und *wie*
sie geschnitten ist. Der Anlass steht im Abstract: MAP ist unter diesem Namen samt Website auf der
Microservices-2019-Konferenz erstmals öffentlich aufgetreten.

Die Motivation (Abschnitt 1) hängt sich an fünf Designfragen auf, die aus den Erfahrungen früher
Microservice-Adopter destilliert sind: Wie viele Operationen gehören in eine API? Welche
Service-Schnitte liefern gemeinsam Nutzen, koppeln Client und Provider aber lose? Wie oft und wie
viel wird ausgetauscht? Welche Nachrichtenstrukturen und Verschachtelungstiefen passen — und wie
ändern sie sich über den Lebenszyklus? Und wie einigt man sich auf die *Bedeutung* von
Repräsentationen und hält sich langfristig daran?

Der Kern der Diagnose: Für vieles gibt es längst gute Ratgeber — RESTful HTTP, asynchrones
Messaging, Strategic Domain-Driven Design zur Service-Identifikation, Infrastrukturpatterns,
Datenhaltung. Was fehlt, ist eine Sprache für das **Strukturieren des Datenaustauschs selbst**. Die
Autoren greifen dafür Pat Hellands Unterscheidung von *data on the inside* und *data on the
outside* auf: Außendaten unterscheiden sich in Veränderlichkeit, Lebensdauer, Genauigkeit,
Konsistenz- und Schutzbedarf grundlegend von Innendaten, und Information Hiding über Servicegrenzen
hinweg bleibt hart — eine einzelne Lösung gibt es nicht (S. 4:3).

Dazu kommen vier strukturelle Spannungsfelder (S. 4:2 f.): *Requirements diversity* (Clients wollen
Unterschiedliches, und das ändert sich), *Design mismatches* (Backends können und strukturieren
anders, als Clients erwarten), *open vs. closed systems* (Publizieren heißt Kontrolle abgeben) und
*stability vs. flexibility* (DevOps liefert immer schneller, APIs sollen stabil bleiben). Daraus
leiten sich drei Trade-offs ab, die die ganze Sprache durchziehen: wenige große Aufrufe gegen viele
feingranulare; stabile, breite Schnittstellen gegen schnell wechselnde, fokussierte;
Datenkonsistenz gegen Verfügbarkeit und schnelle Antwortzeiten.

## Das MAP-Domänenmodell

Das Paper formuliert den Gegenstandsbereich in einem einzigen dichten Satz (Abschnitt 2): MAP
betrachtet API-Design und -Evolution aus der Perspektive der *data on the outside*, also der
Nachrichtenrepräsentationen und Payloads, die beim API-Aufruf ausgetauscht werden. Diese Nachrichten
seien aus **Representation Elements** aufgebaut, die sich in Bedeutung und Struktur unterscheiden,
*weil* API-Endpunkte und ihre Operationen unterschiedliche architektonische Rollen und
Verantwortlichkeiten haben.

Daraus ergibt sich die Begriffskette, auf der die gesamte Sprache steht:

| Begriff | Rolle im Modell |
|---|---|
| **API Provider** | Stellt die API bereit, betreibt die Endpunkte, gibt Kontrolle ab, sobald er publiziert. |
| **API Client** | Ruft auf; seine Informationsbedürfnisse sind divers und veränderlich. |
| **Endpoint** | Adressierbarer Anlaufpunkt der API; trägt eine architektonische *Rolle* (Verantwortlichkeitskategorie). |
| **Operation** | Am Endpunkt angebotener Aufruf; trägt eine eigene *Verantwortlichkeit* (lesend, schreibend, rechnend …). |
| **Message** | Request- bzw. Response-Nachricht, die beim Aufruf fließt — die eigentliche Betrachtungseinheit von MAP. |
| **Representation Element** | Baustein der Nachricht; Strukturform und Stereotyp sind getrennt modelliert (Strukturkategorie). |

Wichtig für die Einordnung: Das Paper zeichnet dieses Modell **nicht** als formales Diagramm. Es
etabliert die Begriffe im Fließtext und illustriert die Kommunikation in Figure 1 (S. 4:2) als
Hexagon-Bild — Microservices tauschen Repräsentationen über plattformunabhängige *Ports* und
technologiespezifische *Adapter* aus. Ein ausgezeichnetes „API domain model", das von
Protokolldetails abstrahiert, gibt es erst auf der Website und im späteren Buch. Dazu passt die
**Protokollneutralität**: Nachrichtenbasierte APIs — RESTful HTTP ebenso wie Event Sourcing und
Streaming — hätten die RPC-Ansätze verdrängt, JSON sei das dominante Austauschformat (S. 4:2).

## Die Kategorien der Pattern-Sprache

Abschnitt 2.2 leitet die Kategorien aus den Fragen ab, die sie beantworten. Vier werden inhaltlich
ausgeführt, zwei nur genannt: *Foundation* und *Identification* bleiben aus Platzgründen außen vor,
wobei *Identification* damals noch gar nicht existierte. Die Conclusion spricht von **sechs
Kategorien**.

| Kategorie | Leitfrage laut Paper (Abschnitt 2.2) | Wiki |
|---|---|---|
| Foundation | Wo und für wen ist die API sichtbar, was wird integriert, wie wird sie beschrieben? (im Paper nur erwähnt) | [Kategorie Foundation](../meta/category-foundation.md) |
| Responsibility | Welche architektonische Rolle spielt jeder Endpunkt und seine Operationen? Wie wirken Rollen und Verantwortlichkeiten auf Servicegröße und Granularität? | [Kategorie Responsibility](../meta/category-responsibility.md) |
| Structure | Wie viele Representation Elements sind angemessen? Wie werden sie strukturiert, gruppiert und mit Metadaten annotiert? | [Kategorie Structure](../meta/category-structure.md) |
| Quality | Wie erreicht ein Provider ein Qualitätsniveau bei wirtschaftlichem Ressourceneinsatz — und wie kommuniziert und verrechnet er die Trade-offs? | [Kategorie Quality](../meta/category-quality.md) |
| Evolution | Wie handhabt man Support-Zeiträume und Versionierung, fördert Rückwärtskompatibilität und kommuniziert brechende Änderungen? | [Kategorie Evolution](../meta/category-evolution.md) |
| Identification | Endpunkt- und Service-Identifikationsstrategien — im Paper als Lücke benannt, nicht als Kategorieinhalt. | — (existiert nicht) |

Die Aufzählungsreihenfolge im Paper ist *Structure, Quality, Responsibility, Evolution* — nicht die
heute übliche Lesereihenfolge, sondern die der Entstehung: Struktur- und Qualitätsfragen zuerst.

## Methodik und Entstehung der Pattern-Sprache

Abschnitt 2.1 begründet, warum die Autoren das Pattern-Format gewählt haben:

- **Pattern-Namen sollen ein Vokabular bilden**, eine *ubiquitous language* im Sinne von Evans.
  Vorbild sind Hohpe/Woolfs *Enterprise Integration Patterns*, die zur Lingua franca des
  queue-basierten Messaging geworden sind; genau eine solche Sprache fehle dem API-Design.
- **Forces und Consequences tragen die Entscheidungsunterstützung.** Nicht die Lösungsskizze ist
  der eigentliche Wert, sondern die explizit gemachten Kräfte und die Nachteile einer Wahl.
- **Patterns sind „weich an den Rändern"** — sie skizzieren Lösungen und liefern keine Blaupausen
  zum blinden Nachbauen.
- **Patterns werden nicht erfunden, sondern gefunden.** Quellen sind öffentliche Web-APIs sowie
  Anwendungs- und Integrationsprojekte der Autoren und ihrer Industriepartner; anschließend
  Kuratierung und Härtung im Peer-Feedback der Writers' Workshops.

Das **Pattern-Template** wird in Abschnitt 3 offengelegt und ist für die Lektüre aller MAP-Texte
nützlich: *Context* (Vorbedingungen der Anwendbarkeit) → *Problem* (als Frage) → *Forces* (warum
es schwer ist; oft mit Non-Solution) → *Solution* (mit Funktionsweise, Beispiel und
Implementierungshinweisen) → *Consequences* (Auflösung der Forces plus Vor- und Nachteile, ggf.
Alternativen) → *Known Uses* → *More Information* (Pattern-Relationen und weitere Quellen).

Zum Stand der Sprache: Vor diesem Paper waren 18 Patterns auf Pattern-Konferenzen publiziert; dieses
Paper bringt zwei weitere. Der Ausblick nennt als offene Arbeit weitere Strukturrepräsentationen,
Endpunkt- und Operationsrollen sowie — als größte Lücke — Identifikationsstrategien, für die Context
Maps, Bounded Contexts und Aggregates aus DDD als Ansatzpunkte gelten. Angekündigt werden außerdem
eine technologieunabhängige Service-Contract-Sprache mit den Patterns als First-Class-Elementen und
Werkzeuge, die Pattern-Instanzen in bestehenden APIs aufspüren.

## Die Pattern-Sprache im Überblick

Figure 2 (S. 4:5) ist die Landkarte: fünf Kategoriekästen, drei davon mit Unterkategorien. Fett
gesetzte Namen waren zum Redaktionsschluss online verfügbar, ausgegraute wurden noch „gemined" —
sichtbar wird damit der Reifegrad: *Quality* und *Evolution* vollständig ausgearbeitet, *Foundation*
und *Responsibility* fast vollständig offen. Die Einträge in Figure-2-Reihenfolge, Namen mit †
hießen 2020 anders (siehe unten):

| Kategorie | Untergruppe | Patterns |
|---|---|---|
| Foundation | — | [Frontend Integration](../foundation/FrontendIntegration.md), [Backend Integration](../foundation/BackendIntegration.md), [Public API](../foundation/PublicAPI.md), [Community API](../foundation/CommunityAPI.md), [Solution Internal API](../foundation/SolutionInternalAPI.md), [API Description](../foundation/APIDescription.md) |
| Responsibility | Endpoint Roles | [Processing Resource](../responsibility/ProcessingResource.md), [Information Holder Resource](../responsibility/InformationHolderResource.md) |
| Responsibility | Processing Responsibilities | [Computation Function](../responsibility/ComputationFunction.md), *Event Processor*, [Retrieval Operation](../responsibility/RetrievalOperation.md), [State Transition Operation](../responsibility/StateTransitionOperation.md) † |
| Responsibility | Information Holders | [Operational Data Holder](../responsibility/OperationalDataHolder.md) †, [Master Data Holder](../responsibility/MasterDataHolder.md), [Reference Data Holder](../responsibility/ReferenceDataHolder.md) †, [Data Transfer Resource](../responsibility/DataTransferResource.md) †, [Link Lookup Resource](../responsibility/LinkLookupResource.md) † |
| Structure | Representation Elements | [Atomic Parameter](../structure/AtomicParameter.md), [Atomic Parameter List](../structure/AtomicParameterList.md), [Parameter Tree](../structure/ParameterTree.md), [Parameter Forest](../structure/ParameterForest.md) |
| Structure | Element Stereotypes | [Data Element](../structure/DataElement.md), [Id Element](../structure/IdElement.md), [Link Element](../structure/LinkElement.md), [Metadata Element](../structure/MetadataElement.md) |
| Structure | Composite Representations | *Annotated Parameter Collection*, [Context Representation](../structure/ContextRepresentation.md), [Pagination](../quality/Pagination.md) |
| Quality | Quality Management and Governance | [API Key](../structure/APIKey.md), [Rate Limit](../quality/RateLimit.md), [Pricing Plan](../quality/PricingPlan.md) †, [Service Level Agreement](../quality/ServiceLevelAgreement.md), [Error Report](../structure/ErrorReport.md) |
| Quality | Data Transfer Parsimony | [Conditional Request](../quality/ConditionalRequest.md), [Request Bundle](../quality/RequestBundle.md), [Wish List](../quality/WishList.md), [Wish Template](../quality/WishTemplate.md) |
| Quality | Reference Management | [Embedded Entity](../quality/EmbeddedEntity.md), [Linked Information Holder](../quality/LinkedInformationHolder.md) |
| Evolution | — | [Version Identifier](../evolution/VersionIdentifier.md), [Semantic Versioning](../evolution/SemanticVersioning.md), [Two In Production](../evolution/TwoInProduction.md), [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md), [Experimental Preview](../evolution/ExperimentalPreview.md), [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md), [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) |

### Namen, die sich seither geändert haben

Die Landkarte von 2020 ist nicht deckungsgleich mit dem heutigen Bestand. Belegt durch die
*a.k.a.*-Zeilen der jeweiligen Pattern-Seiten: *Transactional Data Holder* → *Operational Data
Holder*, *Static Data Holder* → *Reference Data Holder*, *Transfer Resource* → *Data Transfer
Resource*, *Lookup Resource* → *Link Lookup Resource*, *Business Activity Processor* → *State
Transition Operation*, *Rate Plan* → *Pricing Plan*. *Event Processor* und *Annotated Parameter
Collection* kommen unter diesen Namen nicht mehr vor; die Rasterstelle des *Event Processor* —
schreibende Operation, die den Vorzustand nicht liest — besetzt inzwischen die
[State Creation Operation](../responsibility/StateCreationOperation.md). Auch die
Kategoriezuordnung hat sich verschoben: Figure 2 führt *API Key* und *Error Report* unter Quality
und *Pagination* unter Structure; das Wiki folgt der heutigen Einordnung der Website.

### Die beiden ausformulierten Patterns

Abschnitt 3 zeigt exemplarisch, wie ein MAP-Pattern aussieht, und wählt dafür ein Geschwisterpaar
aus dem *Reference Management*: [Embedded Entity](../quality/EmbeddedEntity.md) (verschachteln) und
[Linked Information Holder](../quality/LinkedInformationHolder.md) (verlinken). Beide teilen sich
fast denselben Kontext und beantworten dieselbe Frage gegensätzlich — eingebettete Daten sparen
Roundtrips, kosten Nachrichtengröße und Cachebarkeit; Links sparen Bandbreite, kosten zusätzliche
Requests und Endpunkte. Als Known Uses nennt das Paper unter anderem die GitHub-v3-API, die
Twitter-REST-API und die Microsoft Graph API, die beide Muster nebeneinander verwendet.

## Bemerkenswerte Aussagen

Der methodische Kern steht in Abschnitt 2.1 (S. 4:4): „Patterns are not invented, but mined from
practical experience". Das ist keine Koketterie, sondern die Selbstverpflichtung, die den Ton aller
MAP-Texte erklärt — der Pattern-Autor versteht sich laut Website ausdrücklich als Journalist und
Reporter, nicht als Erfinder.

Zwei weitere Punkte lohnen das Festhalten: dass Datenaustausch ohne Verletzung des Information
Hiding zu strukturieren schwer bleibt und keine einzelne Lösung existiert (S. 4:3) — MAP ist kein
Rezeptbuch, sondern ein geordneter Lösungsraum; und die Selbsteinschätzung im Ausblick, dass MAP ein
Freiwilligenprojekt ist und die Identifikationskategorie — die Frage, *welche* Endpunkte es
überhaupt geben soll — 2020 noch nicht bearbeitet war.

## Bezug zu Kubernetes / KRM

Das Kubernetes Resource Model lässt sich auf MAPs Begriffskette abbilden — die Abbildung ist
erstaunlich glatt, und genau darin liegt der interessante Befund:

| MAP | KRM |
|---|---|
| API Provider | kube-apiserver (inkl. aggregierter Apiserver und CRD-Handler) |
| API Client | kubectl, Controller, Operatoren, Client-Bibliotheken |
| Endpoint | Group-Version-Resource-Pfad, z. B. `/apis/apps/v1/namespaces/{ns}/deployments/{name}`, dazu Subresources wie `.../status` oder `.../scale` |
| Operation | Verb aus einem festen Set: `get`, `list`, `watch`, `create`, `update`, `patch`, `delete`, `deletecollection` |
| Message | das Ressourcenobjekt selbst (bzw. die `List`-Hülle) |
| Representation Element | Felder unter `apiVersion`, `kind`, `metadata`, `spec`, `status` |

Die MAP-Strukturkategorie greift unmittelbar: `metadata.name` und `metadata.uid` sind *Id Elements*,
`metadata.resourceVersion` und `metadata.generation` *Metadata Elements*, `spec` und `status` sind
*Parameter Trees* aus *Data Elements*. Auch das Geschwisterpaar aus Abschnitt 3 findet sich wieder:
Ein `Pod` bettet Container und Volumes unter `spec` vollständig ein (*Embedded Entity*), verweist
auf `ConfigMap` und `Secret` aber nur über deren Namen (*Linked Information Holder*), und
`metadata.ownerReferences` ist eine Liste ebensolcher Verweise.

**Die entscheidende strukturelle Abweichung** liegt eine Ebene höher. MAP modelliert *pro Endpunkt
individuell*: Jeder Endpunkt bekommt seine Rolle, jede Operation ihre Verantwortlichkeit, jede
Nachricht ihre eigens entworfene Repräsentationsstruktur — die Leitfrage „Wie viele Representation
Elements sind angemessen?" setzt genau diese Entwurfsfreiheit voraus. KRM dreht das um und erzwingt
**ein einziges uniformes Schema über alle Ressourcen hinweg**: Die Aufteilung in
`apiVersion`/`kind`/`metadata`/`spec`/`status` gilt für `Pod` wie für jede fremde CRD, und das
Verbset ist geschlossen. Wo MAP fragt, welche Operationen dieser Endpunkt braucht, lautet die
KRM-Antwort: dieselben wie überall. Die Freiheitsgrade der Structure-Kategorie sind vorentschieden —
der Preis für generische Clients, uniformes RBAC, generisches Caching und ein einheitliches
Apply-Modell.

Die zweite Abweichung ist ein Zugewinn auf KRM-Seite: **`watch` ist ein Interaktionsstil, den MAP
nicht abdeckt.** Figure 2 kennt kein Pattern für eine langlaufende Subscription, die einen
Initialzustand und danach einen fortlaufenden Änderungsstrom liefert. Das Paper erwähnt Event
Sourcing und Streaming als verbreitete Protokollwahl (S. 4:2) und nennt die Frage, ob
Zustandsänderungen über API-Aufrufe oder reaktives Event-Streaming gemeldet werden sollen, als
offenen Trade-off (S. 4:3) — die Pattern-Sprache beantwortet ihn 2020 aber nicht. Damit fehlt ihr
genau der Mechanismus, auf dem Informer, Caches und die Controller-Schleife aufsetzen. Auch die
Datensparsamkeitspatterns lesen sich so anders: `watch` mit `resourceVersion` löst dieselbe Aufgabe
wie *Conditional Request*, nur abonnierend statt pollend.

Fazit: Das Domänenmodell **passt** — KRM ist eine legitime MAP-API. Die Pattern-*Sprache* passt nur
teilweise, weil sie Entwurfsfreiheit unterstellt, die KRM bewusst aufgegeben hat, und weil ihr der
Interaktionsstil fehlt, der KRM trägt.

## Verwandte Wiki-Seiten

- Die fünf Kategorien, die dieses Paper als Landkarte einführt:
  [Foundation](../meta/category-foundation.md), [Responsibility](../meta/category-responsibility.md),
  [Structure](../meta/category-structure.md), [Quality](../meta/category-quality.md),
  [Evolution](../meta/category-evolution.md)
- [Embedded Entity](../quality/EmbeddedEntity.md) und
  [Linked Information Holder](../quality/LinkedInformationHolder.md) — die beiden Patterns, die
  dieses Paper erstmals ausformuliert
- [Interface Quality Patterns (EuroPLoP '18)](2018-interface-quality-patterns.md) — eine der drei
  Vorgängerarbeiten, aus denen die damals 18 publizierten Patterns stammen
- [API Description](../foundation/APIDescription.md) — laut Figure 2 das einzige damals fertige
  Foundation-Pattern
- [Information Holder Resource](../responsibility/InformationHolderResource.md) — die Endpunktrolle,
  die KRM flächendeckend realisiert

---
[← Index](../README.md) · [Papers](../papers/) · [Quelle](https://drops.dagstuhl.de/opus/volltexte/2020/11826/)
