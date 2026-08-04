---
title: Rate Limit
kategorie: Quality
unterkategorie: Quality Management and Governance
quelle: https://microservice-api-patterns.org/patterns/quality/qualityManagementAndGovernance/RateLimit
---

# Rate Limit

*a.k.a.* Quota, Usage Limitation

**Kurzform:** Der Provider begrenzt die Zahl der Aufrufe, die ein identifizierter Client in einem
Zeitfenster absetzen darf, und weist Anfragen darüber hinaus ab — zum Schutz der eigenen Ressourcen
und der übrigen Clients.

## Kontext

API-Endpoint und Vertrag stehen; Operationen, Nachrichten und Datenrepräsentationen sind über eine
[API Description](../foundation/APIDescription.md) beschrieben. Clients haben sich beim Provider
registriert und ggf. Nutzungsbedingungen akzeptiert. Ein *Rate Limit* ist aber auch ohne
Vertragsverhältnis sinnvoll, etwa bei Open-Data-Diensten oder in Trial-Phasen.

## Problem

Wie verhindert ein API-Provider exzessive Nutzung durch einzelne Clients?

## Forces

- **Ökonomie:** Jeder Aufruf kostet Rechenzeit, Speicher und Bandbreite; unbegrenzte Nutzung ist
  nicht refinanzierbar.
- **Performance:** Ein einzelner heißlaufender Client darf die Latenz für alle anderen nicht ruinieren.
- **Zuverlässigkeit:** Überlast darf nicht in Kaskadenausfälle münden.
- **Risiko von API-Missbrauch:** Scraping, Credential Stuffing, DoS — Schwere und Eintrittswahr­scheinlichkeit
  bestimmen, wie hart limitiert wird.
- **Client-Awareness:** Ein Limit, das der Client nicht sehen kann, ist für ihn nur zufälliges Scheitern.
  Restkontingent und Reset-Zeitpunkt müssen kommunizierbar sein.

## Lösung

Ein *Rate Limit* pro identifiziertem Client einführen und durchsetzen: Aufrufe werden gezählt, bei
Überschreitung wird abgewiesen statt bedient. Voraussetzung ist eine Client-Identität — [API Key](../structure/APIKey.md),
IP-Adresse oder ein Authentifizierungsverfahren —, denn ohne Zurechenbarkeit gibt es keine Buchführung.

Typische Umsetzungsvarianten: fixes Zeitfenster, Sliding Window, Token Bucket bzw. *Leaky Bucket Counter*
(Hanmer 2007). Durchgesetzt wird häufig nicht im Endpoint selbst, sondern vorgelagert im *API Gateway*.

## Beispiel

Die GitHub-API antwortet bei Überschreitung mit `429 Too Many Requests` und transportiert den Zustand
des Limits in eigenen Headern:

```http
GET https://api.github.com/users/misto
HTTP/1.1 200 OK
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
X-RateLimit-Reset: 1498811560
```

`X-RateLimit-Reset` ist ein Unix-Timestamp, zu dem das Kontingent zurückgesetzt wird.

## Konsequenzen

**Vorteile:**

- Schützt Provider-Infrastruktur und schirmt Clients gegeneinander ab (Noisy-Neighbor-Dämpfung).
- Macht Kapazitätsplanung möglich, weil der Worst Case pro Client bekannt ist.
- Bildet die technische Durchsetzungsschicht für Abrechnungsstufen einer [Pricing Plan](PricingPlan.md).

**Nachteile / Kosten:**

- Erfordert Identifikation und persistentes Zählen pro Client — im horizontal skalierten Cluster ein
  verteiltes Zählproblem mit eigener Latenz und eigenen Konsistenzfragen.
- Falsch dimensionierte Limits blockieren legitime Nutzung; gut dimensionierte Limits verlangen
  Verständnis der Nutzungsprofile.
- Clients brauchen Retry- und Backoff-Logik; ohne sie erzeugt ein Limit erst recht Lastspitzen.
- Schützt nicht gegen verteilte Angriffe von vielen Identitäten.

## Bekannte Verwendungen

- GitHub API v3: 5000 Requests/Stunde für authentifizierte, 60/Stunde für anonyme Clients. Die
  GraphQL-basierte v4-API rechnet zusätzlich die Zahl der abgefragten Knoten an.
- Open Weather Map („access limitation"), Quandl (zusätzlich Limit auf parallele Requests), Twitter
  REST API (15-Minuten-Fenster) — jeweils abhängig vom Abo-Level.
- UID-Register der Schweizer Bundesverwaltung: 20 Requests/Minute, darüber Fehler `Request_limit_exceeded`.
- Let's Encrypt (ACME): wöchentliches Zertifikatslimit pro registrierter Domain plus Limit auf
  Account-Registrierungen pro IP und Stunde.
- API-Gateways wie der MuleSoft API Manager bieten Rate Limiting und Throttling als Konfiguration.

## Verwandte Patterns

- [Pricing Plan](PricingPlan.md) — Limits staffeln oft die Abrechnungsstufen; das Rate Limit ist deren Durchsetzung.
- [Service Level Agreement](ServiceLevelAgreement.md) — Limits sind typischerweise Teil des SLA-Texts.
- [API Key](../structure/APIKey.md) — liefert die Client-Identität, ohne die keine Zählung möglich ist.
- [Wish List](WishList.md), [Wish Template](WishTemplate.md) — reduzieren Datenvolumen und helfen, datenmengenbasierte Limits einzuhalten.
- [Context Representation](../structure/ContextRepresentation.md) — Transportvehikel für den aktuellen Limit-Zustand (Restkontingent, Reset).
- [Pagination](Pagination.md) — begrenzt Antwortgrößen und wirkt damit in dieselbe Richtung, allerdings pro Aufruf statt pro Zeitfenster.
- [API Description](../foundation/APIDescription.md) — Ort, an dem das Limit dokumentiert wird.

## Bezug zu Kubernetes / KRM

Kubernetes erfüllt das Pattern **anders (b)**: Es gibt serverseitige Überlastkontrolle, aber sie ist
*fairness*-orientiert, nicht kontingent-orientiert. Die API Priority and Fairness (APF) modelliert die
Kontrolle selbst als KRM-Ressourcen in der Gruppe `flowcontrol.apiserver.k8s.io` (heute `v1`):
`FlowSchema` klassifiziert eingehende Requests über `spec.rules` (Subjects, `resourceRules`,
`nonResourceRules`), ordnet sie per `spec.matchingPrecedence` und leitet sie an eine
`PriorityLevelConfiguration` weiter. Der Flow Distinguisher (`spec.distinguisherMethod.type`, Werte
`ByUser` oder `ByNamespace`) definiert, was als „ein Client" gilt. Das Priority Level vergibt über
`spec.limited.nominalConcurrencyShares` Anteile an der Gesamtnebenläufigkeit des Servers und kann per
`lendablePercent`/`borrowingLimitPercent` Kapazität an andere Level verleihen bzw. von ihnen borgen.
Beide sind gewöhnliche KRM-Objekte; die folgenden Werte sind die vom Apiserver ausgelieferten
Defaults:

```yaml
apiVersion: flowcontrol.apiserver.k8s.io/v1
kind: FlowSchema
metadata:
  name: service-accounts
spec:
  priorityLevelConfiguration:
    name: workload-low
  matchingPrecedence: 9000          # <- kleinere Zahl zuerst; das erste passende Schema gewinnt
  distinguisherMethod:
    type: ByUser                    # <- definiert, was hier „ein Client" ist (oder ByNamespace)
  rules:
    - subjects:
        - kind: Group
          group:
            name: system:serviceaccounts
      resourceRules:
        - verbs: ["*"]
          apiGroups: ["*"]
          resources: ["*"]
          namespaces: ["*"]
          clusterScope: true
      # ... nonResourceRules
---
apiVersion: flowcontrol.apiserver.k8s.io/v1
kind: PriorityLevelConfiguration
metadata:
  name: workload-low
spec:
  type: Limited
  limited:
    nominalConcurrencyShares: 100   # <- Anteil an gleichzeitigen Seats, keine Requests/Stunde
    lendablePercent: 90
    limitResponse:
      type: Queue                   # <- Alternative: Reject, dann sofort 429
      queuing:
        queues: 128
        handSize: 6                 # <- Shuffle Sharding: 6 von 128 Queues pro Client-Hash
        queueLengthLimit: 50
```

Entscheidend ist die Abweichung: APF limitiert **Nebenläufigkeit (Seats), nicht Requests pro Zeitraum**.
Es gibt kein „X Aufrufe pro Stunde pro Consumer". Statt abzuweisen, wird bei `limitResponse.type: Queue`
eingereiht; das `queuing`-Objekt oben implementiert Shuffle Sharding: ein Client wird per Hash auf
`handSize` von `queues` Warteschlangen abgebildet, sodass ein Vielschreiber nur einen Bruchteil der
anderen Clients trifft, statt sie alle zu blockieren. Nur `limitResponse.type: Reject` (und der
Queue-Overflow) führt zu `429`; `tooManyRequests()` in
`staging/src/k8s.io/apiserver/pkg/server/filters/priority-and-fairness.go` setzt dann `Retry-After` und
`http.StatusTooManyRequests`. Zur Nachvollziehbarkeit setzt der Server auf jede Antwort zwei Header
mit **UIDs statt Namen** — damit die Klassifikation nicht die intern gewählten
Prioritätsbezeichnungen preisgibt:

```http
HTTP/1.1 429 Too Many Requests
Retry-After: 1
X-Kubernetes-PF-FlowSchema-UID: 4f2c1e8a-...      # <- keine Namen, absichtlich
X-Kubernetes-PF-PriorityLevel-UID: b7d0a913-...
Content-Type: text/plain; charset=utf-8

Too many requests, please try again later.
```

Beachtenswert: Das ist kein `metav1.Status`, sondern die Klartext-Antwort von `http.Error` — der
Request wird im Filter verworfen, bevor die API-Maschinerie ihn überhaupt sieht.

Ein `exempt`-Level nimmt Systemverkehr ganz von der Begrenzung aus. Ältere, gröbere Stellschrauben
sind die Flags `--max-requests-inflight` und `--max-mutating-requests-inflight`.

Für Ablehnungen jenseits von APF existiert das Muster ebenfalls: `errors.NewTooManyRequests(message, retryAfterSeconds)`
in apimachinery füllt `Status.details.retryAfterSeconds`, und der Filter
`staging/src/k8s.io/apiserver/pkg/server/filters/with_retry_after.go` sendet während des Shutdowns
`Retry-After: 5`, damit Clients auf eine andere Apiserver-Instanz ausweichen. Ein Client-Kontingent im
MAP-Sinn existiert nur punktuell, etwa im Admission-Plugin `EventRateLimit`
(`plugin/pkg/admission/eventratelimit`), das Event-Schreibraten pro Server, Namespace, User oder
Quelle+Objekt begrenzt.

Bemerkenswert ist die **Verlagerung eines Teils des Rate Limits auf den Client**: `rest.Config` in
client-go trägt `QPS` und `Burst` sowie einen austauschbaren `RateLimiter` aus
`k8s.io/client-go/util/flowcontrol`:

```go
cfg.QPS = 5.0   // rest.DefaultQPS   — das Kontingent setzt sich der Client selbst
cfg.Burst = 10  // rest.DefaultBurst
```

Dazu kommen Backoff-Manager und die Rate Limiter der Workqueues für Controller. Kubernetes
verlässt sich also auf kooperative Clients plus serverseitige
Fairness statt auf ein hartes Consumer-Kontingent — passend zu einem Modell, in dem die „Clients"
überwiegend eigene Controller im selben Cluster und keine fremden Vertragspartner sind.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/qualityManagementAndGovernance/RateLimit)
