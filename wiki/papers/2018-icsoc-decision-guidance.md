---
title: Guiding Architectural Decision Making on Quality Aspects in Microservice APIs
kategorie: Paper
autoren: Uwe Zdun, Mirko Stocker, Olaf Zimmermann, Cesare Pautasso, Daniel Lübke
venue: ICSOC 2018 (16th International Conference on Service-Oriented Computing)
quelle: http://eprints.cs.univie.ac.at/5956/1/Guiding%20Architectural%20Decision%20Making%20on%20Quality%20Aspects%20in%20Microservice%20APIs.pdf
---

# Guiding Architectural Decision Making on Quality Aspects in Microservice APIs

## Bibliografische Angaben

| Feld | Wert |
|---|---|
| Autoren | Uwe Zdun (Univ. Wien), Mirko Stocker und Olaf Zimmermann (HSR Rapperswil), Cesare Pautasso (USI Lugano), Daniel Lübke (innoQ Schweiz) |
| Venue / Umfang | ICSOC 2018, 15 Seiten; Schwester-Publikation ist Stocker et al., *Interface Quality Patterns*, EuroPLoP 2018 — dort stehen die Patterns selbst |

## Worum es geht

Das Paper ist **kein Pattern-Paper**. Es nimmt die Interface-Quality-Patterns als gegeben und fragt
eine Ebene darüber: *Welche Architekturentscheidungen stehen dahinter, welche Optionen gibt es je
Entscheidung, welche Kriterien (decision drivers) entscheiden, und wie hängen die Entscheidungen
zusammen?* Ergebnis ist ein formales, wiederverwendbares Architectural-Design-Decision-Modell (ADD):

- **2 Decision Contexts** — `API & API Client` (fünf Entscheidungen) und `API Operation` (eine) —
  mit **6 Decisions**, **40 Decision Options** und **47 Decision Drivers**.
- Optionen sind typisiert als `Pattern`, `Practice` (Praxis ohne Pattern-Beschreibung) oder
  `Do Nothing`; Beziehungen als `«Variant»`, `«Includes»`, `«Realizes»`, `«Can Be Combined With»`,
  `«Can Use»`, `«Influences»`, `«Consider If Not Decided Yet»`.

Der zweite Beitrag ist eine **Uncertainty-Reduction-Schätzung**: Wie viel weniger muss ein Architekt
betrachten, wenn er das Modell benutzt statt dasselbe Wissen unstrukturiert? Scope sind allein die
**Message Representations im Interface-Kontrakt**.

## Methodik

| Aspekt | Vorgehen |
|---|---|
| Grundmethode | Pattern Mining als qualitative Forschung, an Grounded Theory angelehnt — konkret an Charmaz' *constructivist GT*, weil die Studie mit expliziten Forschungsfragen (RQ1/RQ2) startet |
| Kodierung | Open Coding und Axial Coding, `constant comparison`; textuelle Codes nur initial, dann Überführung in UML-Modelle und Pattern-Templates |
| Abbruchkriterium | Theoretical Saturation, sehr konservativ operationalisiert: Stopp erst, wenn 5–7 zusätzliche Quellen nichts Neues mehr beitragen |
| Quellen | 55 insgesamt: **31 real genutzte APIs** (AWS EC2/S3/Lambda, GitHub v3 und v4, Google Calendar, Microsoft Graph, Stripe, PayPal, Twitter, Salesforce, SWIFT, TMForum, Schweizer Banken- und Versicherungs-APIs …) plus **24 Spezifikationen, Standards, Technologien** (JSON API, HTTP/1.1 Conditional Requests, OAuth, OpenID Connect, SAML, Kerberos, LDAP, RFC 7519, OWASP REST Security …) |
| Auswahl / Review | breite reale Nutzung, moderne Service-Technologie, mindestens teilweise Microservice-Tenets; jeder Fund in mindestens fünf Iterationen von verschiedenen Autoren geprüft |

**Threats to Validity** laut Autoren: überwiegend öffentliche Internet-APIs, kaum In-House-APIs —
interne Praktiken könnten fehlen; **kein Vollständigkeitsanspruch**; die Uncertainty-Zahlen sind
Schätzungen, keine formale Evaluation.

## Die Entscheidungen im Detail

### Überblick: das Entscheidungsmodell als Baum

```
Quality Category
├── Endpoint-Specific Qualities      (Kontext: API + API Client)
│   ├── D1  Client Identification and Authentication
│   ├── D2  Communicate Errors
│   ├── D3  Prevent API Clients From Excessive API Usage  ──┐ AND-Gruppe; jede zieht
│   ├── D4  Metering and Charging for API Consumption     ──┤ D1 per „Consider If
│   └── D5  Explicit Spec. of Quality Objectives/Penalties──┘ Not Decided Yet"
└── Operation-Specific Qualities     (Kontext: einzelne API Operation)
    └── D6  Avoid Unnecessary Data Transfer
```

**D1 ist die Voraussetzungsentscheidung**: Man kann weder limitieren noch abrechnen noch Garantien
zusichern, ohne den Client zu kennen. D6 ist die einzige Entscheidung **pro Operation**.

### D1 — Client Identification and Authentication

| Option | Typ | Beschreibung |
|---|---|---|
| no secure identification and authentication needed | Do Nothing | tragfähig nur bei begrenzter Client-Zahl und geringem Missbrauchsrisiko |
| identification and authentication via shared secret | Pattern (API Key) | eindeutiges Token pro Client, minimalistisch |
| … secured with secret key | Pattern (API Key + Secret Key) | zusätzlicher, **nicht übertragener** Schlüssel — erst damit echte Authentifizierung |
| via a dedicated protocol | Practice (Authentication / Authorization Protocol) | OAuth, SAML, Kerberos, LDAP |

| Kriterium | Wirkung der Optionen |
|---|---|
| Security-Level | Do Nothing ≪ API Key ≪ API Key + Secret Key ≈ Protokolle |
| Ease of Use, Credential Management, Performance | API Key kaum schlechter als Do Nothing und mit geringem Overhead; Protokolle brauchen komplexe Protokoll-APIs, Infrastruktur, beidseitige Account-Verwaltung und kosten mehr Laufzeit (weil mehr Features) |
| Kopplung | API Keys **entkoppeln** den aufrufenden Client von der Organisation des Kunden — Account-Credentials gäben Admins und Entwicklern unnötig Vollzugriff |

### D2 — Communicate Errors

| Option | Typ | Beschreibung |
|---|---|---|
| provide no specific solution | Do Nothing | für Produktiv-APIs in der Regel nicht ratsam |
| nur Protokoll-Fehlercodes | Practice (Protocol-Level Error Codes) | z. B. HTTP-Statuscodes; funktioniert nur bei einem einzigen Protokoll-Stack |
| API-level Error Reporting | Pattern (Error Reporting) | maschinenlesbarer Code **plus** Textbeschreibung, mit Parametern/Konstanten für i18n; `«Can Be Combined With»` den Fehlercodes — der Regelfall |

| Kriterium | Wirkung |
|---|---|
| Defect Fixing, Robustness, Reliability, Maintainability, Evolvability | Haupttreiber; je mehr die Meldung über die Ursache sagt, desto geringer der Aufwand der Fehlersuche — Error Reporting > Fehlercodes |
| Interoperability & Portability | Error Reporting besser, weil Protokoll-, Format- und Plattform-Autonomie möglich |
| Security | **gegenläufig**: detaillierte Meldungen legen Interna offen und öffnen Angriffsvektoren |
| Internationalisierung | Error Reporting erzeugt Übersetzungsaufwand |

### D3 — Prevent API Clients From Excessive API Usage

Binär `yes` (Pattern **Rate Limit**) / `no`; `«Can Use» Authentication`. Das Limit gilt als Requests pro Zeitraum; danach wird abgelehnt, verzögert oder niedriger priorisiert.

| Kriterium | Wirkung |
|---|---|
| Scalability, Performance, Resilience, Fault Tolerance | müssen gehalten werden und sind bei Missbrauch in Gefahr; Limits machen Abuse schwer |
| Client Awareness | Clients müssen erfahren können, wie viel ihres Limits verbraucht ist |
| Impact/Severity des Missbrauchsrisikos, Economic Aspects | bestimmen, ob sich der Aufwand lohnt; Limits kosten Geld, werden von Clients kritisch gesehen, und verhandelbare Limits erzeugen Zusatzkomplexität |

Die **schwächste** Entscheidung im Modell: kein Kriterium ist vorentscheidbar, die Unsicherheitsreduktion beträgt nur 50 % (und 0 % bei Decision Nodes und Outcomes).

### D4 — Metering and Charging for API Consumption

Binär `yes`/`no`; `yes` = Pattern **Rate Plan**, das `«Can Use»` Rate Limit (Fair Use) und
Authentication. Varianten: **Usage-based Pricing** (pro Call), **Market-based Allocation** (mit
Untervariante **Auction-style Allocation**), **Flat-rate Subscription** — jede davon und die
Basisvariante kombinierbar mit einem **Freemium Model**.

| Kriterium | Wirkung |
|---|---|
| Economic Aspects | Haupttreiber; die Variante muss zum Geschäftsmodell von Provider oder Consumer passen |
| Accuracy und Meter Granularity | Clients erwarten, nur Konsumiertes zu zahlen — das setzt eine angemessene Messgranularität voraus |
| Security | **vorentschieden negativ**: Abrechnungsdaten verraten, wie gut ein Kunde im Markt dasteht, brauchen also Extraschutz — `Do Nothing` ist hier der sicherere Weg |

### D5 — Explicit Specification of Quality Objectives and Penalties

Binär `yes`/`no`; `yes` = Pattern **Service Level Agreement** als Erweiterung der API Description mit
messbaren SLOs und Strafen; `«Can Use»` API Description, Authentication, Rate Limit, Rate Plan.
Varianten: SLA nur intern, SLA mit **formal** und SLA mit **informell** spezifizierten SLOs.

| Kriterium | Wirkung |
|---|---|
| Attractiveness (Consumer) vs. Cost-Efficiency und Business Risk (Provider) | die zentrale Abwägung |
| Regulation und Legal Obligations | z. B. DSGVO/GDPR erzwingen Zusagen |
| Die zugesicherten Qualitäten selbst, Business Agility & Vitality | typischerweise Availability, Performance/Scalability, Security/Privacy — sie werden selbst zu Decision Drivers, weil das Geschäftsmodell des Clients von ihnen abhängen kann |

### D6 — Avoid Unnecessary Data Transfer (pro Operation)

| Option | Typ | Auslösende Situation |
|---|---|---|
| no data transfer reduction possible or wanted | Do Nothing | — |
| use simple list to provide the information | Pattern (Wish List) | Client-Bedarfe unterschiedlich oder unvorhersehbar, flache Auswahl reicht |
| use a structured template | Pattern (Wish Template) | verschachtelte oder repetitive Parameterstrukturen |
| make data transfer dependent on a condition | Pattern (Conditional Request) | viele Clients fragen wiederholt dieselben, selten geänderten Daten ab |
| bundle multiple requests in a container message | Pattern (Request Bundle) | ein Client stellt viele logisch zusammengehörige Requests |

Kombinationen: Conditional Request `«Can Be Combined With»` Wish List **oder** Wish Template
(„welche Teilmenge der geänderten Daten will ich?"), Request Bundle mit allen vorherigen — jede
Kombination erhöht aber die API-Komplexität. Alle vier `«Influences»` **Rate Limit**.

| Kriterium | Wirkung |
|---|---|
| Individual Information Needs | der Haupttreiber; ohne Analyse keine sinnvolle Pattern-Wahl |
| Data Parsimony, Performance | Bandbreite, Response Time, Throughput, Processing Time, Kosten |
| Security | **zweischneidig**: Wish List/Template können ungewollt sensible Daten exponieren oder Angriffsvektoren öffnen — andererseits kann nicht Übertragenes nicht gestohlen werden |
| Complexity of API Design and Programming, Test- und Wartungsaufwand | steigen bei allen vier, weil der Client zur Laufzeit bestimmt, was er bekommt, und die Sonderfälle die Testmatrix vervielfachen; GraphQL gilt den Autoren als Extremform des Wish Template |

### Die Uncertainty-Reduction-Schätzung

Gemessen werden **Decision Nodes** (`ndec`), nötige **Criteria Assessments** (`ncri`) und mögliche
**Outcomes** (`ndo`), je mit (⊕) und ohne (⊖) Modell. Ohne Modell ist jede Nicht-`Do-Nothing`-Option
ein eigener Ja/Nein-Knoten, und weil zulässige Kombinationen unbekannt sind, muss die
**Potenzmenge** aller Knoten betrachtet werden.

| Decision | ndec ⊕/⊖ | ncri ⊕/⊖ | ndo ⊕/⊖ | Reduktion (ndec / ncri / ndo) |
|---|---|---|---|---|
| D1 Client Identification | 1 / 4 | 1 / 44 | 5 / 16 | 75 % / 97,73 % / 68,75 % |
| D2 Communicate Errors | 1 / 2 | 1 / 18 | 4 / 4 | 50 % / 94,44 % / 0 % |
| D3 Excessive API Usage | 1 / 1 | 4 / 8 | 2 / 2 | 0 % / 50 % / 0 % |
| D4 Metering and Charging | 1 / 6 | 4 / 24 | 12 / 64 | 83,33 % / 83,33 % / 81,25 % |
| D5 Quality Objectives | 1 / 4 | 2 / 28 | 5 / 16 | 75 % / 92,86 % / 68,75 % |
| **Kontext API/Client gesamt** | 5 / 17 | 12 / 122 | 2,68·10⁸ / 4,83·10¹¹ | 70,59 % / 90,16 % / 99,94 % |
| D6 = Kontext Operation gesamt | 1 / 4 | 2 / 32 | 12 / 16 | 75 % / 93,75 % / 25 % |

Der große Hebel bei `ncri` kommt von den **vorentschiedenen** Kriterien: Man bewertet einen Vektor
bereits geklärter Kriterien und prüft nur den Rest. Die Grenze benennen die Autoren selbst — auch
*mit* Modell bleiben 2,68·10⁸ Outcomes; **Kombinationen mehrerer Entscheidungen** brauchen Tooling.

## Bezug zu den MAP-Patterns

| Decision | Option(en) im Modell | Wiki-Seite |
|---|---|---|
| D1 Client Identification and Authentication | API Key, API Key + Secret Key | [API Key](../structure/APIKey.md) |
| D2 Communicate Errors | Error Reporting | [Error Report](../structure/ErrorReport.md) |
| D3 Prevent Excessive API Usage | Rate Limit | [Rate Limit](../quality/RateLimit.md) |
| D4 Metering and Charging | Rate Plan (+ Varianten) | [Pricing Plan](../quality/PricingPlan.md) |
| D5 Quality Objectives and Penalties | Service Level Agreement (+ Varianten) | [Service Level Agreement](../quality/ServiceLevelAgreement.md) |
| D6 Avoid Unnecessary Data Transfer | Wish List, Wish Template, Conditional Request, Request Bundle | [Wish List](../quality/WishList.md) · [Wish Template](../quality/WishTemplate.md) · [Conditional Request](../quality/ConditionalRequest.md) · [Request Bundle](../quality/RequestBundle.md) |

Nicht Teil des Modells, aber aus derselben Pattern-Familie: [Pagination](../quality/Pagination.md)
(Reduktion über die Länge statt die Breite), [Embedded Entity](../quality/EmbeddedEntity.md) und
[Linked Information Holder](../quality/LinkedInformationHolder.md) (dieselbe Abwägung zur Designzeit)
sowie deren Strukturgrundlagen [Parameter Tree](../structure/ParameterTree.md) und [Link Element](../structure/LinkElement.md).

## Bemerkenswerte Aussagen

Zum zweischneidigen Sicherheitsargument der Data-Transfer-Patterns:

> „data that is not transferred cannot be stolen and cannot be tampered with"
> — S. 10 des PDF

Er steht direkt neben der gegenteiligen Feststellung, dass client-gesteuerte Feldauswahl sensible
Daten *exponieren* kann — diese Doppelrichtung eines einzigen Kriteriums begründet die
Unterscheidung zwischen vorentschiedenen und offen bleibenden Kriterien.

## Bezug zu Kubernetes / KRM

Kubernetes hat **jede dieser sechs Entscheidungen genau einmal zentral** getroffen und für alle
Ressourcen festgeschrieben — statt sie pro Endpoint oder pro Operation offenzulassen.

| Decision | KRM-Festlegung |
|---|---|
| D1 Identification/Authentication | Im API-Server, nicht pro Endpoint: Authenticator-Kette aus X.509-Client-Zertifikaten, Bearer Tokens (insb. ServiceAccount-JWTs), OIDC und Webhook. Kein API-Key-Konzept; die Autorisierung folgt uniform über RBAC auf `apiGroup`/`resource`/`verb`. |
| D2 Communicate Errors | Beide Optionen des Modells kombiniert, aber einmal für alles: `metav1.Status` mit `code`, `reason`, `message`, `details` liegt **jedem** Fehler jeder Ressource bei, zusätzlich zum HTTP-Statuscode. `apierrors.IsNotFound`, `IsConflict` usw. funktionieren deshalb generisch. |
| D3 Excessive Usage | API Priority and Fairness (`flowcontrol.apiserver.k8s.io` mit `FlowSchema` + `PriorityLevelConfiguration`) plus die älteren `--max-requests-inflight`/`--max-mutating-requests-inflight`. Limitiert werden Request-*Flows* am zentralen Handler, nicht einzelne Endpoints; `client-go` ergänzt client-seitiges Rate Limiting. |
| D4 Metering and Charging | Bewusst `Do Nothing` — die Kubernetes-API ist kein kommerzielles Angebot; Abrechnung findet, wenn überhaupt, bei Cloud-Providern eine Ebene darüber statt. |
| D5 Quality Objectives | Kein SLA-Objekt in der API. Die verbindlichen Zusagen sind Prozesszusagen: die API-Deprecation-Policy (GA-Versionen mindestens 12 Monate bzw. 3 Releases nach Deprecation weiter unterstützt, Beta 9 Monate bzw. 3 Releases, Alpha ohne Garantie) und die von SIG Scalability gepflegten SLIs/SLOs für API-Call-Latenzen. |
| D6 Avoid Unnecessary Data Transfer | Keine Wish List, kein Wish Template — kein generisches Sparse-Fieldset-Protokoll. Stattdessen wenige vordefinierte Projektionen per Content Negotiation (`as=PartialObjectMetadata`, `as=Table`) und Zeilen- statt Feldselektion (`labelSelector`, `fieldSelector`). Statt Conditional Request via ETag/If-None-Match dient `resourceVersion` als Optimistic-Concurrency-Precondition (409 bei Konflikt) und als Watch-Startpunkt. Statt Request Bundle nur `deletecollection`, sonst ein Request pro Objekt. Und statt Polling der Default-Weg `watch` bzw. `WatchList` mit Informer-Cache. |

Dazu kommen zwei Festlegungen, die das Modell nicht führt, weil MAP sie unter Reference Management
ablegt: **Referenz statt Einbettung als Default** (`ownerReferences`, `secretKeyRef`,
`configMapRef`, Volume-Referenzen) und **Cursor-Pagination** (`limit` + opaker `continue`-Token in
`ListOptions`, `remainingItemCount` in `ListMeta`; kein Offset, keine Seitennummern).

**Was das kostet.** Das Paper sagt für D6 ausdrücklich, die Entscheidung müsse *pro Operation*
fallen, weil nur eine Analyse der individuellen Client-Bedarfe zeige, ob sich Transfer reduzieren
lässt. KRM verzichtet genau darauf. Wer von 5.000 Pods nur `status.phase` braucht, bekommt 5.000
vollständige Pod-Objekte, sofern er nicht mit Metadata-only oder Table auskommt; große `ConfigMap`s
und `Secret`s kommen immer ganz. Diese Kosten sind der Grund für Metadata-only-Informer,
`WatchList`/`sendInitialEvents` und die 1,5-MiB-Objektgrenze. Auch D5 fehlt maschinenlesbar:
Qualitätszusagen lassen sich nicht abfragen, nur in der Dokumentation nachlesen.

**Was es bringt.** KRM zäunt die Uncertainty-Rechnung des Papers von der anderen Seite auf: Statt
den Entscheidungsraum zu *strukturieren*, eliminiert es ihn. Für den Autor einer CRD ist `ndec` über
alle sechs Entscheidungen **null** — Authentifizierung, Fehlerformat, Rate Limiting, Pagination,
Watch-Semantik und Referenzstil liefert der API-Server. Das ist stärkere Guidance, als ein
Entscheidungsmodell sie geben kann, und zahlt sich aus: Generische Maschinerie funktioniert über
Ressourcen, die sie nie gesehen hat — `kubectl get` zeigt eine fremde CRD an, der Informer-Cache
indiziert sie, Server-Side Apply merged sie, der Garbage Collector räumt sie auf, RBAC autorisiert
sie. Jeder eigene Fehlercode, jedes eigene Pagination-Schema, jeder Sparse-Fieldset-Dialekt pro
Endpoint würde das zerstören. KRM tauscht also **Optimierbarkeit pro Endpoint** gegen
**Uniformität über alle Endpoints**.

## Verwandte Wiki-Seiten

- D6-Optionen: [Wish List](../quality/WishList.md) · [Wish Template](../quality/WishTemplate.md) · [Conditional Request](../quality/ConditionalRequest.md) · [Request Bundle](../quality/RequestBundle.md)
- Patterns hinter D1–D5: [API Key](../structure/APIKey.md) · [Error Report](../structure/ErrorReport.md) · [Rate Limit](../quality/RateLimit.md) · [Pricing Plan](../quality/PricingPlan.md) · [Service Level Agreement](../quality/ServiceLevelAgreement.md)
- Benachbart, nicht modelliert: [Pagination](../quality/Pagination.md) · [Embedded Entity](../quality/EmbeddedEntity.md) · [Linked Information Holder](../quality/LinkedInformationHolder.md) · [Parameter Tree](../structure/ParameterTree.md) · [Link Element](../structure/LinkElement.md)
- [Kategorie Quality](../meta/category-quality.md) — Überblick über alle Quality-Patterns

---
[← Index](../README.md) · [Papers](../papers/) · [Quelle](http://eprints.cs.univie.ac.at/5956/1/Guiding%20Architectural%20Decision%20Making%20on%20Quality%20Aspects%20in%20Microservice%20APIs.pdf)
