---
title: State Transition Operation
kategorie: Responsibility
unterkategorie: Operation Responsibilities
quelle: https://microservice-api-patterns.org/patterns/responsibility/operationResponsibilities/StateTransitionOperation
---

# State Transition Operation

*a.k.a.* Read-Write Operation, Data Change Operation, Business Activity Processor

**Kurzform:** Eine Operation, die Client-Eingabe mit vorhandenem Provider-Zustand kombiniert, einen
gültigen Zustandsübergang ausführt und ein nicht-triviales Ergebnis zurückgibt — der API-seitige
Ausdruck einer Geschäftsaktivität.

## Kontext

Geschäftsfunktionalität soll in einer API exponiert und in mehrere Aktivitäten zerlegt werden, deren
Ausführungszustand in der API sichtbar ist, damit Clients ihn vorantreiben können. Typisch für
langlaufende Geschäftsprozesse mit inkrementellen Updates und koordiniertem Zustandsmanagement, oft
vorab modelliert in BPMN, UML-Aktivitätsdiagrammen oder Use-Case-Modellen.

## Problem

Wie kann ein Client eine Verarbeitung anstoßen, die den Provider-seitigen Anwendungszustand ändert? Und
wie teilen Client und Provider die Verantwortung für Ausführung und Steuerung von Geschäftsprozessen?
Die Quelle unterscheidet zwei Pole: *Frontend BPM* (der Client besitzt den Prozesszustand und ruft
feingranulare Aktivitäten auf) und *BPM services* (der Provider besitzt Prozess und Zustand, der Client
initiiert und verfolgt nur). Kanonisches Beispiel ist die Schadensbearbeitung einer Versicherung:
Validierung, Betrugsprüfung, Korrespondenz, Entscheidung, Auszahlung, Archivierung — teils parallel,
teils streng sequenziell, über Tage bis Jahre.

## Forces

- **Service-Granularität:** atomare Aktivität, Teilprozess oder ganzer Prozess pro Operation?
- **Konsistenz und Auditierbarkeit:** Zustandsübergänge müssen nachvollziehbar und prüfbar sein.
- **Abhängigkeiten von vorherigen Zustandsänderungen**, die mit anderen kollidieren können
  (Transaktionsproblematik).
- **Lastmanagement.**
- **Netzwerkeffizienz gegen Datensparsamkeit** (Nachrichtengrößen).
- Zeitverhalten und Zuverlässigkeit gelten ebenfalls, werden aber unter
  [State Creation Operation](StateCreationOperation.md) diskutiert.

## Lösung

Eine Operation `sto: (in, S) -> (out, S')` wird eingeführt. Die gültigen Zustandsübergänge werden im
Endpunkt modelliert, und die Gültigkeit eingehender Änderungswünsche wird zur Laufzeit geprüft. Auf
Nachrichtenebene wird ein *Command Message* mit einem *Document Message* gepaart: Das Kommando beschreibt
die gewünschte Aktion, die Antwort liefert Quittung oder Ergebnis.

## Beispiel

Aus dem Online-Shop: „proceed to checkout and pay“ ist eine Aktivität im Bestellprozess, „add item to
shopping basket“ eine Aktivität im Teilprozess Katalog-Browsing. Beide ändern Provider-Zustand, tragen
Geschäftssemantik und haben nicht-triviale Vor- und Nachbedingungen sowie Invarianten — etwa: nicht
ausliefern und fakturieren, bevor der Kunde bestellt und bestätigt hat.

## Konsequenzen

**Vorteile:**

- Der Endpunkt kapselt den Zustandsautomaten; ungültige Übergänge werden zentral abgewiesen, statt in
  jedem Client nachgebildet zu werden.
- Geschäftssemantik ist im Kontrakt sichtbar und damit auditierbar.
- Sowohl Push- als auch Pull-Datenflüsse möglich.

**Nachteile / Kosten:**

- Stärkste Kopplung aller vier Operationsverantwortlichkeiten: Der Client muss den Zustandsautomaten
  kennen, um ihn korrekt zu bedienen.
- Zustandsbehaftete, oft langlaufende *conversations* mit entsprechenden Transaktions-, Timeout- und
  Kompensationsproblemen.
- Die Quelle warnt, nicht die gesamte Domäne als *Published Language* zu exponieren — das koppelt
  Clients an die Provider-Implementierung.

## Bekannte Verwendungen

- `BookingService.java` und `CargoInspectionService.java` in der DDD Sample Application.
- PayPals *controller resources*; M. Nygard empfiehlt verhaltensorientierte Activity Sets / Process
  Services.
- Der Großteil der 1000+ Services einer Core-Banking-SOA (z. B. Geldtransfers), Telco-Order-Management
  (z. B. Technikertermin bei Umzug), Terravis (Hypothekenübertrag, Einreichen signierter Dokumente).

## Verwandte Patterns

- [State Creation Operation](StateCreationOperation.md) — schreibt nur (append), referenziert
  üblicherweise kein bestehendes Zustandselement.
- [Retrieval Operation](RetrievalOperation.md) — liest, schreibt nicht.
- [Computation Function](ComputationFunction.md) — berührt Zustand überhaupt nicht.
- [Processing Resource](ProcessingResource.md) — die typische Endpunktrolle für Geschäftsaktivitäten.
- [Information Holder Resource](InformationHolderResource.md) — kann ebenfalls Zustandsübergänge
  anbieten.
- [Community API](../foundation/CommunityAPI.md) — häufiger Sichtbarkeitsbereich solcher Operationen,
  abgesichert per [API Key](../structure/APIKey.md) und geregelt per
  [Service Level Agreement](../quality/ServiceLevelAgreement.md).

## Bezug zu Kubernetes / KRM

Hier weicht KRM am deutlichsten vom Pattern ab: **Das Kubernetes Resource Model hat praktisch keine
expliziten State Transition Operations.** Das feste Verbset `get`, `list`, `watch`, `create`, `update`,
`patch`, `delete`, `deletecollection` enthält kein Verb, das eine benannte Geschäftsaktivität ausdrückt.
Es gibt kein `approve`, `promote`, `restart` oder `checkout` als API-Verb — solche Wünsche werden
stattdessen deklarativ als Änderung an `spec` formuliert und per `update` oder `patch` geschrieben.

```http
# So nicht — es gibt kein Verb für Geschäftsübergänge; der Pfad existiert schlicht nicht:
POST /apis/apps/v1/namespaces/default/deployments/web/restart
```

```yaml
# Sondern: Zielzustand deklarieren, der Controller organisiert den Übergang.
# PATCH /apis/apps/v1/namespaces/default/deployments/web
spec:
  template:
    metadata:
      annotations:
        kubectl.kubernetes.io/restartedAt: "2026-08-04T10:00:00Z"   # <- genau das schreibt `kubectl rollout restart`
```

Die Annotation ist inhaltlich bedeutungslos; sie ändert nur den Pod-Template-Hash und löst damit
einen Rollout aus. Das „Kommando“ ist eine Datenänderung, die zufällig einen Effekt hat.

Der Zustandsübergang findet damit nicht in der API-Operation statt, sondern danach: Ein Controller
beobachtet die Ressource per `watch`, vergleicht `spec` (gewünschter Zustand) mit `status`
(beobachteter Zustand) und arbeitet auf Angleichung hin — level-triggered, nicht edge-triggered.
Der Aufruf ist damit idempotent und wiederholbar; ein zweimal geschriebenes `replicas: 3` erzeugt
keinen zweiten Übergang. Genau deshalb ist die im Pattern zentrale Laufzeitprüfung „ist dieser Übergang
vom aktuellen Zustand aus gültig?“ in KRM die Ausnahme: Admission-Plugins und Validierungs-Strategien
prüfen zwar Invarianten und teilweise Übergänge (etwa Immutability bestimmter Felder), aber der
Normalfall ist, dass ein beliebiger Zielzustand deklariert werden darf und der Controller die Reise
dorthin organisiert.

Prozesszustand wird nicht in Operationen, sondern in Feldern geführt: `status.conditions`,
`metadata.generation` gegen `status.observedGeneration`, sowie `metadata.deletionTimestamp` plus
`metadata.finalizers` für den Löschprozess.

```yaml
metadata:
  generation: 8              # <- der API-Server zählt hoch, sobald spec sich ändert
  # ... name, uid, resourceVersion
spec:
  replicas: 5
status:
  observedGeneration: 7      # <- Controller hat Generation 8 noch nicht verarbeitet: Übergang läuft
  readyReplicas: 3
  conditions:
    - type: Progressing
      status: "True"
      reason: ReplicaSetUpdated
      lastUpdateTime: "2026-08-04T10:00:05Z"
      lastTransitionTime: "2026-08-04T10:00:02Z"   # <- wann der Übergang stattfand — als Datum, nicht als Aufruf
    - type: Available
      status: "False"
      reason: MinimumReplicasUnavailable
      lastTransitionTime: "2026-08-04T10:00:02Z"
```

Ein Client, der wissen will, ob „sein“ Übergang durch ist, vergleicht also zwei Zahlen und liest
Bedingungen — er bekommt kein Ergebnis auf seinen Schreibaufruf zurück. `delete` ist dadurch selbst
kein sofortiger Übergang, sondern setzt bei vorhandenen Finalizern nur `deletionTimestamp` und startet
eine Reconciliation-Kaskade. Optimistische Nebenläufigkeitskontrolle über `metadata.resourceVersion`
ersetzt die Transaktionsklammer, die das Pattern diskutiert.

Ausnahmen existieren, sind aber wenige und bewusst gesetzte imperative Subresources:
`pods/binding` (`create`, weist einen Pod einem Node zu, siehe `pkg/registry/core/pod/storage/storage.go`),
`pods/eviction` (`create`, prüft PodDisruptionBudgets, siehe `.../storage/eviction.go`),
`certificatesigningrequests/approval` (`update` mit eigener Strategie, siehe
`pkg/registry/certificates/certificates/storage/storage.go`) und
`serviceaccounts/token` (`create`, TokenRequest).

```yaml
# PUT /apis/certificates.k8s.io/v1/certificatesigningrequests/node-csr-bWFpbg/approval
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: node-csr-bWFpbg
status:
  conditions:
    - type: Approved       # <- hier *ist* die Bedingung der Übergang, nicht dessen Beobachtung
      status: "True"
      reason: KubectlApprove
      message: "This CSR was approved by kubectl certificate approve."
# Die approval-Strategie lässt ausschließlich Änderungen an conditions zu — spec und
# status.certificate sind über diesen Pfad gesperrt.
```

Auch `scale` ist ein Subresource, bleibt aber
deklarativ: `get`/`update`/`patch` auf ein `Scale`-Objekt. Die Existenz dieser Handvoll Fälle bestätigt
die Regel — wo ein Übergang wirklich imperativ, nicht wiederholbar und nicht als Zielzustand
formulierbar ist, braucht KRM einen Sonderweg. Fazit: (b) anders und eigenwillig gelöst, und zwar
absichtlich, weil ein uniformes Schema über alle Ressourcen generische Clients, generische Controller
und generische Werkzeuge wie `kubectl` erst möglich macht.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/operationResponsibilities/StateTransitionOperation)
