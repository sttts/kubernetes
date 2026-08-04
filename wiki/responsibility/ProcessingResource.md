---
title: Processing Resource
kategorie: Responsibility
unterkategorie: Endpoint Roles
quelle: https://microservice-api-patterns.org/patterns/responsibility/endpointRoles/ProcessingResource
---

# Processing Resource

*a.k.a.* Command Service, Controller Resource, Executor, Processing Endpoint

**Kurzform:** Ein API-Endpunkt, dessen Identität eine *Aktivität* ist, nicht ein Datensatz. Clients
stoßen damit eine Verarbeitung beim Provider an; Daten erscheinen nur als Request- und
Response-Payload, nicht als adressierbarer Zustand.

## Kontext

Die fachlichen Anforderungen einer Anwendung liegen vor (z.B. als User Stories oder
Geschäftsprozessmodelle). Aus ihnen folgt, dass entfernt aufrufbare Fähigkeiten benötigt werden —
[Frontend Integration](../foundation/FrontendIntegration.md) und/oder
[Backend Integration](../foundation/BackendIntegration.md). Eine (Micro-)Service-Architektur und die
Integrationsinfrastruktur sind grob festgelegt.

## Problem

Wie kann ein API-Provider seinen Clients erlauben, bei ihm eine Aktion auszulösen?

Solche Aktionen können eigenständige Kommandos sein (fachlich oder technisch) oder Aktivitäten
innerhalb eines Geschäftsprozesses. Sie dürfen providerseitigen Zustand lesen und schreiben, müssen
es aber nicht. Gewünscht ist ein Abstraktionsniveau, das die Aktion exponiert und die Daten so weit
wie möglich verbirgt.

## Forces

- **Ausdrucksstärke des Kontrakts vs. Service-Granularität** und deren Kopplungswirkung: je
  ausdrucksstärker, desto mehr muss gelernt, verwaltet und getestet werden.
- **Erlernbarkeit und Betreibbarkeit** — feingranulare Services sind leichter zu schützen und zu
  evolvieren, es gibt aber viele davon, die integriert werden müssen.
- **Semantische Interoperabilität** — Aktionsnamen tragen mehr implizite Bedeutung als Datenschemata.
- **Antwortzeit** — Verarbeitung kostet Zeit; synchron oder asynchron?
- **Sicherheit und Datenschutz.**
- **Kompatibilität und Evolvierbarkeit.**

Die Leitentscheidung: soll der Endpunkt aktivitätsorientierte oder datenorientierte Semantik haben?
Dieses Pattern beschreibt den aktivitätsorientierten Pol, das Geschwister-Pattern
[Information Holder Resource](InformationHolderResource.md) den datenorientierten.

## Lösung

Füge der API einen *Processing Resource*-Endpunkt hinzu, dessen Operationen fachliche Aktivitäten
oder Kommandos bündeln und kapseln.

Der Endpunkt kann zustandslos sein (reine Berechnung) oder providerseitigen Zustand über mehrere
Aufrufe fortschreiben. Auf Operationsebene wird das durch die vier Operation-Responsibility-Patterns
differenziert: [Computation Function](ComputationFunction.md) (zustandslos),
[Retrieval Operation](RetrievalOperation.md) (nur lesend),
[State Creation Operation](StateCreationOperation.md) (nur anlegend) und
[State Transition Operation](StateTransitionOperation.md) (lesend und schreibend).

## Beispiel

Aus der Lakeside-Mutual-Beispielanwendung: `RiskComputationService` ist eine zustandslose
*Processing Resource* mit genau einer *Computation Function*.

```java
@RestController
@RequestMapping("/riskfactor")
public class RiskComputationService {
    @PostMapping(value = "/compute")
    public ResponseEntity computeRiskFactor(
        @Valid @RequestBody RiskFactorRequestDto riskFactorRequest) {
        int riskFactor = computeRiskFactor(age, postalCode);
        return ResponseEntity.ok(new RiskFactorResponseDto(riskFactor));
    }
}
```

Demgegenüber ist `InsuranceQuoteRequestProcessingResource` eine zustandsbehaftete *Processing
Resource*, deren *State Transition Operations* einen Angebotsantrag durch mehrere Stadien bewegen.

## Konsequenzen

**Vorteile:**

- Die Aktion ist die Abstraktionseinheit: interne Datenmodelle bleiben verborgen, die Kopplung an
  providerseitige Schemata sinkt.
- Fachliche Invarianten lassen sich in einem Aufruf durchsetzen; keine Sequenz von CRUD-Calls, die
  der Client selbst korrekt orchestrieren müsste.
- Ideal für Prozessschritte, Validierungen und Berechnungen ohne natürliches Datenpendant.

**Nachteile / Kosten:**

- Jede Operation hat eigene Semantik — mehr zu lernen, zu dokumentieren, zu versionieren.
- Uniforme HTTP-Semantik (Caching, Idempotenz, bedingte Anfragen) greift kaum; Kommandonamen in URIs
  gelten Teilen der Web-API-Community als REST-Antipattern.
- Generische Werkzeuge (Clients, Caches, Generatoren) können weniger beitragen.

## Bekannte Verwendungen

- Die Slack Web API ist verarbeitungsorientiert (`https://slack.com/api/METHOD_FAMILY.method`) und
  spricht selbst von „HTTP RPC methods“.
- `BookingService` und `CargoInspectionService` in der Domain-Driven-Design-Beispielanwendung.
- SOA-Landschaften in Unternehmen: das „Dynamic Interface“ eines Kernbankensystems (produktiv seit
  2003), das Terravis-Prozess-Hub für Schweizer Hypothekarprozesse (*Command Services*, z.B.
  Vertragsunterzeichnung), sowie die Order-Management-SOA aus Zimmermann et al. (2005).

## Verwandte Patterns

- [Information Holder Resource](InformationHolderResource.md) — die gegenteilige Semantik und die
  direkte Alternative zu diesem Pattern.
- [Computation Function](ComputationFunction.md) — zustandslose Operation innerhalb einer *Processing Resource*.
- [State Creation Operation](StateCreationOperation.md), [State Transition Operation](StateTransitionOperation.md),
  [Retrieval Operation](RetrievalOperation.md) — die übrigen Operationsverantwortlichkeiten.
- [Public API](../foundation/PublicAPI.md), [Community API](../foundation/CommunityAPI.md) — die Sichtbarkeiten,
  in denen *Processing Resources* typischerweise angeboten werden.
- [API Key](../structure/APIKey.md), [Rate Limit](../quality/RateLimit.md),
  [Service Level Agreement](../quality/ServiceLevelAgreement.md) — Schutz und Governance des Endpunkts.
- [Context Representation](../structure/ContextRepresentation.md) — hält technische Parameter aus dem
  fachlichen Payload heraus.
- [Error Report](../structure/ErrorReport.md) — Fehlermeldung bei fehlgeschlagener Verarbeitung.

## Bezug zu Kubernetes / KRM

Das Kubernetes Resource Model erfüllt dieses Pattern bewusst **fast gar nicht** (Fall c). Die
API-Konventionen fordern, dass praktisch alles eine Ressource mit dem uniformen Verb-Satz
`get`/`list`/`watch`/`create`/`update`/`patch`/`delete`/`deletecollection` ist; RBAC, Auditing,
Admission, Watch-Cache und `kubectl` setzen genau diesen uniformen Kontrakt voraus. Statt eine Aktion
aufzurufen, schreibt ein Client den *gewünschten Zustand* und ein Controller gleicht ihn
level-triggered ab — die Aktivität liegt hinter der API, nicht in ihr.

Echte *Processing Resources* existieren in Kubernetes nur als eng umgrenzte Ausnahmen, und fast alle
sind Subresources oder synthetische „Review“-Objekte:

- **Streaming-Subresources**: `pods/exec`, `pods/attach`, `pods/portforward`, `pods/proxy`,
  `nodes/proxy`, `services/proxy`. Sie verlassen die CRUD-Semantik vollständig und upgraden die
  Verbindung auf einen Stream (SPDY bzw. WebSocket), obwohl nichts persistiert wird.

  ```http
  # Aktivitätsorientiert: die Identität des Endpunkts ist "führe dieses Kommando aus".
  POST /api/v1/namespaces/default/pods/web-7d9f/exec?container=app&command=sh&command=-c&command=id&stdout=true&stderr=true&tty=true
  Connection: Upgrade
  Upgrade: SPDY/3.1
  X-Stream-Protocol-Version: v4.channel.k8s.io    # <- ab hier ein Stream, kein Request/Response mehr
  ```

  ```http
  # Datenorientiert: dieselbe Ressource als Information Holder — Pfad ohne Verb, Body ist ein Objekt.
  GET /api/v1/namespaces/default/pods/web-7d9f
  Accept: application/json
  ```

  Die Query-Parameter sind kein freies Format, sondern die serialisierten Felder von
  `PodExecOptions` (`stdin`, `stdout`, `stderr`, `tty`, `container`, `command`) — ein
  Kommandoparametersatz, der sich als Options-Objekt tarnt. `kubectl exec` spricht wahlweise SPDY
  über `POST` oder WebSocket (`v5.channel.k8s.io`) über `GET`; der Endpunkt akzeptiert beide
  Methoden, weshalb je nach Transport das RBAC-Verb `create` oder `get` greift.

- **Write-only-Kommandos als POST auf ein Subresource-Objekt**: `pods/eviction` (nimmt ein
  `policy/v1` `Eviction`-Objekt entgegen und respektiert PodDisruptionBudgets), `pods/binding`
  (Scheduler-Zuweisung), `namespaces/finalize`,
  `serviceaccounts/token` (TokenRequest, liefert ein frisch signiertes Token zurück).

  ```yaml
  # POST /api/v1/namespaces/default/pods/web-7d9f/eviction
  apiVersion: policy/v1
  kind: Eviction
  metadata:
    name: web-7d9f                # <- muss den Pod aus dem Pfad benennen
    namespace: default
  deleteOptions:                  # <- kein spec/status: der ganze Body ist Kommandoparameter
    gracePeriodSeconds: 30
  # Antwort: metav1.Status mit status: Success — bzw. 429 TooManyRequests,
  # wenn ein PodDisruptionBudget die Verdrängung gerade verbietet.
  ```

- **Request/Response-Objekte, die nie gespeichert werden**: `SubjectAccessReview`,
  `SelfSubjectAccessReview`, `SelfSubjectRulesReview`, `SelfSubjectReview`, `TokenReview`. Sie werden
  per `create` gesendet, der API-Server füllt `status` und gibt das Objekt zurück, ohne es in etcd zu
  schreiben. Das ist ein RPC im Ressourcen-Kostüm — MAP würde das als
  [Computation Function](ComputationFunction.md) auf einer *Processing Resource* einordnen.
- **Grenzfall `scale`**: `deployments/scale` & Co. sehen datenorientiert aus (ein `Scale`-Objekt mit
  `spec.replicas`/`status.replicas`), werden aber als Aktion benutzt (`kubectl scale`). Formal bleibt
  es ein Information Holder mit projiziertem Schema.
- Außerhalb des Ressourcenmodells liegen die *nonResourceURLs* (`/healthz`, `/readyz`, `/livez`,
  `/metrics`, `/version`); RBAC hat für sie eine eigene Regelart, weil Gruppe/Ressource/Namespace
  dort nicht greifen.

Der Preis dieser Strenge ist sichtbar: Wo Kubernetes doch ein Kommando braucht, muss es das
Ressourcenmodell verbiegen (Objekte, die nie gespeichert werden; `create`-Verben ohne Persistenz).
Der Gewinn ist der uniforme Kontrakt — Watch, Server-Side Apply, Feldselektoren, generische Clients
und Autorisierung funktionieren für jede Ressource gleich, ohne pro Aktion neu spezifiziert zu werden.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/endpointRoles/ProcessingResource)
