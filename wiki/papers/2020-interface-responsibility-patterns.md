---
title: "Interface Responsibility Patterns: Processing Resources and Operation Responsibilities"
kategorie: Paper
autoren: Olaf Zimmermann, Daniel Lübke, Uwe Zdun, Cesare Pautasso, Mirko Stocker
venue: EuroPLoP '20, 1.–4. Juli 2020, Virtual Event, Germany (ACM)
quelle: http://eprints.cs.univie.ac.at/6520/1/MAP-EuroPlop2020aPaper.pdf
---

# Interface Responsibility Patterns: Processing Resources and Operation Responsibilities

## Bibliografische Angaben

Olaf Zimmermann und Mirko Stocker (OST Rapperswil), Daniel Lübke (iQuest, Hannover), Uwe Zdun
(Universität Wien), Cesare Pautasso (USI Lugano): *Interface Responsibility Patterns: Processing
Resources and Operation Responsibilities.* EuroPLoP '20, 1.–4. Juli 2020, Virtual Event, Germany. ACM,
24 Seiten. DOI [10.1145/3424771.3424822](https://doi.org/10.1145/3424771.3424822). Primärquelle für fünf
Patterns der Kategorie [Responsibility](../meta/category-responsibility.md); das Begleitpaper derselben
Autoren (EuroPLoP '19, „Data-Oriented Interface Responsibility Patterns“) deckt die datenorientierte
Gegenseite ab, also [Information Holder Resource](../responsibility/InformationHolderResource.md).

## Worum es geht

Zwei gestaffelte Entwurfsfragen. Auf **Endpunktebene**: Soll *Verarbeitung* oder sollen *Daten* das
leitende Konzept sein? Auf **Operationsebene**: Welche Verantwortlichkeit trägt jede Operation beim
Lesen und Schreiben von Provider-Zustand? Antwort auf die zweite Frage ist eine 2×2-Matrix:

| | schreibt nicht | schreibt |
|---|---|---|
| **liest nicht** | *Computation Function* `f: in -> out` | *State Creation Operation* `f: in -> (out, S')` |
| **liest** | *Retrieval Operation* `f: (in, S) -> out` | *State Transition Operation* `f: (in, S) -> (out, S')` |

## Behandelte Patterns

| Pattern | Paper-Abschnitt | Wiki-Link |
|---|---|---|
| Processing Resource | 4.1 (S. 4–9) | [ProcessingResource.md](../responsibility/ProcessingResource.md) |
| Computation Function | 4.2 (S. 9–12) | [ComputationFunction.md](../responsibility/ComputationFunction.md) |
| State Creation Operation | 4.3 (S. 12–15) | [StateCreationOperation.md](../responsibility/StateCreationOperation.md) |
| Retrieval Operation | 4.4 (S. 15–18) | [RetrievalOperation.md](../responsibility/RetrievalOperation.md) |
| State Transition Operation | 4.5 (S. 18–22) | [StateTransitionOperation.md](../responsibility/StateTransitionOperation.md) |
| Information Holder Resource | nur als Kontrastfolie (Abschn. 1, Abb. 1) | [InformationHolderResource.md](../responsibility/InformationHolderResource.md) |

## Was das Paper gegenüber der Website ergänzt

Die MAP-Website kürzt *Forces* und *Consequences* auf Stichworte oder verweist aufs Buch; das Paper
führt beide aus und ergänzt je einen *Non-solution*-Abschnitt — den verworfenen naheliegenden Entwurf.

### Processing Resource

**Forces.** *Kontraktausdrucksstärke und Granularität:* Viele einfache Interaktionen geben dem Client
Kontrolle und Effizienz, erzeugen aber Koordinationsaufwand und Evolutionsprobleme; wenige reichhaltige
Fähigkeiten fördern Konsistenz, passen aber nicht zu jedem Client und verschwenden Ressourcen.
Mehrdeutige Aufrufsemantik schadet der Interoperabilität und führt zu ungültigen Ergebnissen. —
*Erlernbarkeit und Verwaltbarkeit:* Zu viele Aktionen erzeugen Orientierungsprobleme für
Client-Programmierer, Tester und Wartungspersonal, das nicht mehr das ursprüngliche Team sein muss. —
*Semantische Interoperabilität:* Die API Description muss festhalten, was
eine Operation tut **und was nicht** — bezogen auf Zustandsänderungen, Idempotenz, Transaktionalität,
Event-Emission und Verbrauch nachgelagerter Ressourcen. — *Antwortzeit:* Je länger der Client blockiert,
desto wahrscheinlicher bricht etwas; wartende Endanwender drücken Refresh und erzeugen Zusatzlast. —
*Sicherheit und Datenschutz:* Bei Audit-Pflicht ist Provider-Zustandslosigkeit nicht erreichbar; nicht
jeder Client darf jede Aktion auslösen (Beispiel des Papers: Angestellte dürfen ihr eigenes Gehalt nicht
erhöhen); das Sicherheitsdesign muss Policy Decision und Enforcement Points berücksichtigen, zwischen
RBAC und ABAC entscheiden und DoS, gefälschte Bestellungen und betrügerische Schadensmeldungen im
Bedrohungsmodell führen. — *Kompatibilität:* Ändert sich der Datenkontrakt (Maßeinheiten, neue optionale
Parameter), muss der Client das bemerken und reagieren können.

**Non-solution.** Eine *Shared Database* mit Stored Procedures — verbreitet, aber Single Point of
Failure, skaliert nicht und ist nicht unabhängig deploybar.

**Consequences.** (+) Aktivitäts- und Prozessorientierung kann Kopplung senken und Information Hiding
fördern. (−) In vielen Integrationsszenarien müsste sie dem Entwurf aufgezwungen werden, was ihn schwer
wartbar macht; dann ist *Information Holder Resource* die bessere Wahl. Die Auflösung **aller übrigen
Kräfte** delegiert das Paper auf die Operationsebene — die Endpunktrolle legt nur die Perspektive fest.

### Computation Function

**Forces.** *Netzwerkeffizienz gegen Datensparsamkeit:* Kleine Nachrichten bedeuten viele Nachrichten;
wenige große erzeugen weniger Verkehr, sind aber schwerer zu erzeugen und zu verarbeiten. —
*Reproduzierbarkeit:* Lokale Aufrufe lassen sich leicht protokollieren und wiederholen; Auslagerung
kostet Kontrolle und Latenz. — *Lastmanagement:* Rechenintensive oder lange laufende Funktionen
gefährden Skalierbarkeit und SLA; greift die Berechnung auf serverseitige Ressourcen zu (Logger,
Backend-Dienste), wird sie zustandsbehaftet und skaliert nicht mehr leicht horizontal.

**Non-solution.** Lokal rechnen — kann große Datenmengen erfordern, den Client verlangsamen und führt
letztlich zu einer monolithischen Architektur.

**Consequences.** (+) Lastmanagement wird einfacher, weil zustandslose Operationen frei verschoben
werden können. (−) Reproduzierbarkeit und Auditierbarkeit leiden: Es entsteht eine externe Abhängigkeit
außerhalb der Client-Kontrolle; er muss darauf vertrauen, dass wiederholte Aufrufe dasselbe liefern.
(−) Nachrichten werden größer, weil zustandslose Server keine Zwischenergebnisse aus eigenen Speichern
nachziehen können. Ist ein Remote-Aufruf zu teuer, bleibt die lokale Bibliothek die billigere
Alternative. Varianten: *Transformation Service*, *Validation Service* und *Long Running Computation*
(asynchrones Messaging, Call-mit-Callback oder Long-Running-Request mit Polling-Link). Bezieht ein
Validator Provider-Zustand ein, wird er zum „Business Rule Validator“.

### State Creation Operation

**Forces.** *Kopplungs-Trade-off:* Damit die Provider-Verarbeitung einfach bleibt, sollte die Meldung
selbsttragend und unabhängig von anderen Ereignissen sein; damit Client-Konstruktion schlank bleibt,
Transportkapazität gespart und Interna verborgen werden, sollte sie nur das Minimum enthalten — beides
zugleich geht nicht. — *Konsistenzeffekte:* Wenn Provider-Zustand nicht gelesen werden kann oder soll,
ist schwer zu prüfen, ob die Verarbeitung Invarianten verletzt. — *Zeitverhalten:* Eintritt, Meldung und
Eintreffen fallen auseinander; die Reihenfolge von Vorfällen verschiedener Clients ist womöglich nicht
bestimmbar — Zeitsynchronisation ist eine Schranke verteilter Systeme. — *Zuverlässigkeit:* Meldungen
kommen in anderer Reihenfolge an, gehen verloren oder treffen mehrfach ein.

**Non-solution.** Einfach eine weitere Operation ohne besondere Semantik ergänzen: Die
Integrationsannahmen müssen dann in Dokumentation und Beispielen explizit gemacht werden, sonst geraten
sie in Vergessenheit; die Kohäsion leidet, DevOps muss raten, wo zu deployen ist.

**Consequences.** (+) Lose Kopplung, weil Client und Provider keinen Anwendungszustand teilen: Beim
Eintreffen wird kein Provider-Zustand gelesen. (−) Genau deshalb sind keine Provider-seitigen Prüfungen
möglich; Konsistenz lässt sich nicht vollständig sicherstellen. (−) Zeitmanagement bleibt aus demselben
Grund schwierig. (−) Zuverlässigkeit leidet ohne Bestätigung oder Zustandsbezeichner; kommt einer
zurück, muss der Client ihn korrekt interpretieren. Ein [Rate Limit](../quality/RateLimit.md) ist
heikel, weil Events mit hohem Durchsatz auftreten. Varianten: *Event Notification Operation*
(Grundlage für Event Sourcing), *Bulk Report*.

### Retrieval Operation

**Forces.** *Veracity, Variety, Velocity, Volume* — die vier V der Big-Data-Diskussion; Daten kommen in
vielen Formen, das Client-Interesse variiert. Dazu *Lastmanagement* und *Netzwerkeffizienz gegen
Datensparsamkeit*, beide wie bei *Computation Function*.

**Non-solution.** Alle Daten periodisch „hinter den Kulissen“ replizieren — mit gravierenden Mängeln
bei Konsistenz, Verwaltbarkeit und Datenaktualität.

**Consequences.** (+) Lastmanagement: Wegen ihrer Nur-lesend-Natur skalieren sie durch Replikation der
Daten. (+) Netzwerkeffizienz: Sie können Identifikatoren voll ausnutzen und lokale Daten bedarfsweise
holen, cachen und optimieren — es muss nicht alles in der Anfrage stehen. (−) Sie werden zum Engpass,
wenn Informationsbedarf und angebotene Abfragefähigkeiten nicht zusammenpassen.
Ausdrücklich **nicht** aufgelöst ist das Vier-V-Force: [Pagination](../quality/Pagination.md) adressiert
„Volume“, „Velocity“ bräuchte Stream Processing. Varianten: *Status Check*, *Time-Bound Report* und
*Business Rule Validator* (validiert vorhandene Provider-Daten statt übergebener).

### State Transition Operation

**Forces.** *Service-Granularität:* Große Services tragen reichhaltigen Zustand, der nur in wenigen
Übergängen aktualisiert wird; kleine sind einfach, aber gesprächig. — *Konsistenz:* Prozessinstanzen
unterliegen oft der Auditpflicht; abhängig vom aktuellen Zustand dürfen bestimmte Aktivitäten *nicht*
ausgeführt werden; manche müssen in Zeitfenstern abgeschlossen werden, weil sie reservierte Ressourcen
belegen; geht etwas schief, muss rückgängig gemacht werden. — *Abhängigkeiten von zuvor erfolgten
Zustandsänderungen:* Kollisionen mit Transaktionen anderer Clients, externen Ereignissen oder
Provider-internen Batch-Jobs. — *Netzwerkeffizienz gegen Datensparsamkeit* (nur das Delta schicken) und
*Lastmanagement*. Zeitverhalten und Zuverlässigkeit gelten auch, stehen aber unter *State Creation
Operation*.

**Non-solution.** Provider-Zustand vollständig verbieten (realistisch nur bei Taschenrechnern), oder
zustandslose Operationen anbieten und den Zustand jedes Mal mitübertragen (*Client Session State*;
HATEOAS begünstigt das). Das skaliert gut, öffnet aber Sicherheitsrisiken und untergräbt die
Auditierbarkeit: Wie garantiert man dann, dass nur gültige Abläufe vorkommen
(`order->pay->deliver->return->refund`) und betrügerische ausgeschlossen sind?

**Consequences.** (+) Netzwerkeffizienz: Ein RESTful-Entwurf kann über Zustandstransfers und
Ressourcenschnitt eine Balance aus Ausdrucksstärke und Effizienz finden. (+) Granularität: Das Pattern
verträgt kleinere wie größere „service cuts“ und fördert Agilität. (+) Konsistenz: Solche Operationen
können und müssen Geschäfts- und Systemtransaktionsmanagement intern behandeln. (−) Abhängigkeiten von
vorher erfolgten Zustandsänderungen können kollidieren. (−) Lastmanagement: Zustandsbehaftete
Operationen skalieren nicht leicht und lassen sich nicht nahtlos verlagern — im Widerspruch zu den
IDEAL-Eigenschaften von Cloud-Anwendungen.

Die Lösung ist weiter ausgeführt als auf der Website: generische Prozesssteuerungs-Primitive (*prepare*,
*start*, *suspend*/*resume*, *cancel*, *undo*, *restart*, *cleanup*) plus die Ereignisse *completed*
(finished/failed/aborted) und *stateChanged*, als Zustandsautomat mit Ready, Running, Suspended,
Completed. Varianten: *Full Overwrite* (`PUT`) gegen *Partial Change* (`PATCH`), auf Nachrichtenebene
*Full Report* gegen *Delta Report*. Absolute Updates („setze x auf y“) sind inkrementellen vorzuziehen,
weil letztere bei Duplikaten Daten verfälschen.

## Pattern-Sprache und Zusammenhänge

Abbildung 1 ordnet die Patterns dreistufig: oben die Endpunktebene (*funktionale gegen datenorientierte
Perspektive*), darunter die Operationsebene mit den Achsen *State Read* und *State Write*, unten die
Realisierungsstrategie. Wichtig ist die Nichtsymmetrie: Eine *Processing Resource* enthält
typischerweise *Computation Functions*, *State Creation Operations* und *State Transition Operations*;
*Retrieval Operations* sollten dort auf Statusabfragen beschränkt bleiben und gehören eher in eine
*Information Holder Resource*.

Methodisch stützt sich das Paper auf Responsibility-Driven Design: Operationen übernehmen
*responsibilities*, Endpunkte bündeln sie zu *roles*, Aufrufe sind *collaborations*, die API Description
ist der *contract*; eine *Processing Resource* ist ein *interfacer*, der Zugriff auf *service providers*,
*controllers* und *coordinators* bereitstellt und schützt. Eine Fußnote erklärt die Namensgebung:
zustandserhaltende Verantwortlichkeiten heißen *functions*, zustandsändernde *operations*.

Nach außen verbindet sich die Sprache mit den übrigen MAP-Kategorien (Nachrichtenstrukturen wie
[Data Element](../structure/DataElement.md), Schutz über API Key und Rate Limit, Versionierung über die
Evolution Patterns) und mit fremden Sprachen: *Command Message*, *Document Message*, *Request-Reply* und
*Idempotent Receiver* aus den Enterprise Integration Patterns, *Business Transaction* und *Client Session
State* aus Fowlers PoEAA, *Command* aus GoF, Service, Aggregate, Entity und Domain Event aus DDD.

## Bemerkenswerte Aussagen

Zum Sicherheits-Force der *Processing Resource* halten die Autoren fest, dass bei Audit-Pflicht
„statelessness on the provider side is an illusion“ (S. 5) — selbst dann, wenn die Fachlichkeit gar
keinen Zustand verlangt. Als Härtetest für fachliche Ausrichtung schlagen sie vor (S. 21), die
API-Dokumentation Domänenexperten ohne Informatikausbildung zum Kommentieren vorzulegen.

## Bezug zu Kubernetes / KRM

Liest man das Paper mit dem Kubernetes Resource Model im Kopf, ist der Befund einseitig: **Kubernetes
kennt praktisch keine *Processing Resources* und praktisch keine *State Transition Operations*.** Das
Paper beschreibt einen Entwurfsraum mit zwei Polen — KRM besetzt konsequent nur den datenorientierten.
Der Grund ist das feste, ressourcenunabhängige Verbset
`get`/`list`/`watch`/`create`/`update`/`patch`/`delete`/`deletecollection` (Abbildung von HTTP-Methoden
auf Verbnamen: `staging/src/k8s.io/apiserver/pkg/endpoints/installer.go`). Damit entfällt die zentrale
Entwurfsfreiheit des Papers: Keine benannte Operation drückt eine fachliche Aktivität aus — kein
`approve`, kein `promote`, kein `checkout`. Ein uniformes Verbprofil, das für Core-Typen wie für CRDs
gilt, ersetzt operationsspezifische Verantwortlichkeiten.

Der Zustandsübergang verschwindet dadurch nicht, er wandert aus der Operation heraus: in die Trennung
von `spec` (gewünschter Zustand) und `status` (beobachteter Zustand) plus einen level-triggered
reconciliierenden Controller. Der Client schreibt einen Zielzustand, der Controller organisiert die
Reise dorthin. Das dreht die Force-Bilanz um: „Abhängigkeiten von zuvor erfolgten Zustandsänderungen“
löst KRM über optimistische Nebenläufigkeitskontrolle (`metadata.resourceVersion`) statt über
Transaktionsklammern; absolute statt inkrementeller Updates sind keine Empfehlung, sondern die
Grundform; und der Prozesszustand liegt in Feldern statt Operationen (`status.conditions`,
`metadata.generation` gegen `status.observedGeneration`, `metadata.deletionTimestamp` mit
`metadata.finalizers`).

Die Ausnahmen sind zählbar und fast durchweg Subresources oder synthetische Objekte:

- **Streaming-Subresources** wie `pods/exec`, `pods/attach`, `pods/portforward` — sie verlassen die
  CRUD-Semantik vollständig und upgraden die Verbindung auf einen Stream.
- **Imperative Kommandos als `create` auf ein Subresource-Objekt**: `pods/eviction` (nimmt ein
  `Eviction`-Objekt entgegen, respektiert PodDisruptionBudgets,
  `pkg/registry/core/pod/storage/eviction.go`), `pods/binding` (weist einen Pod einem Node zu,
  `pkg/registry/core/pod/storage/storage.go`).
- **Eine State Transition Operation in Reinform**: `certificatesigningrequests/approval` mit eigener
  Update-Strategie (`ApprovalREST` in `pkg/registry/certificates/certificates/storage/storage.go`) —
  ein fachlicher Zustandsübergang mit eigenem RBAC-Verb.
- **Nie persistierte `*Review`-Objekte** unter `pkg/registry/authorization/` (`SubjectAccessReview`,
  `SelfSubjectAccessReview`, `LocalSubjectAccessReview`, `SelfSubjectRulesReview`) und
  `pkg/registry/authentication/` (`TokenReview`, `SelfSubjectReview`): per `create` geschickt, der
  Server füllt `status` und gibt das Objekt zurück, ohne es zu speichern — in der Sprache des Papers
  *Computation Functions* auf einer *Processing Resource*, nur im Ressourcen-Kostüm.

Was das Paper als Trade-off je Endpunkt darstellt, hat Kubernetes einmal global entschieden. Der Preis:
Wo ein Kommando unvermeidlich ist, muss das Modell verbogen werden — `create` ohne Persistenz, Objekte,
die nie in etcd landen. Der Gewinn ist der uniforme Kontrakt: Watch, Server-Side Apply, Feldselektoren,
RBAC, Auditing, Admission und `kubectl` gelten für jede Ressource gleich.

## Verwandte Wiki-Seiten

- Endpunktrollen: [Processing Resource](../responsibility/ProcessingResource.md),
  [Information Holder Resource](../responsibility/InformationHolderResource.md)
- Operationen: [Computation Function](../responsibility/ComputationFunction.md),
  [State Creation Operation](../responsibility/StateCreationOperation.md),
  [Retrieval Operation](../responsibility/RetrievalOperation.md),
  [State Transition Operation](../responsibility/StateTransitionOperation.md)
- Umfeld: [Data Element](../structure/DataElement.md), [Kategorie Responsibility](../meta/category-responsibility.md)

---
[← Index](../README.md) · [Papers](../papers/) · [Quelle](http://eprints.cs.univie.ac.at/6520/1/MAP-EuroPlop2020aPaper.pdf)
