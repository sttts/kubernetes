---
title: State Creation Operation
kategorie: Responsibility
unterkategorie: Operation Responsibilities
quelle: https://microservice-api-patterns.org/patterns/responsibility/operationResponsibilities/StateCreationOperation
---

# State Creation Operation

*a.k.a.* Write-Only Operation, Data Insertion Operation

**Kurzform:** Eine Operation, deren Verantwortung ausschließlich das Anlegen von Provider-seitigem
Zustand ist. Der Client meldet ein Ereignis oder legt einen neuen Datensatz an und erwartet nicht mehr
als eine Empfangsbestätigung — keine fachliche Antwort, keine Zustandsauswertung.

## Kontext

Ein API-Endpunkt existiert bereits, funktionale und nicht-funktionale Anforderungen sind erhoben. Der
Client will den Provider über clientseitig eingetretene Vorfälle informieren, ohne dass eine unmittelbare
inhaltliche Antwort nötig wäre. Typische Auslöser: das Anstoßen einer langlaufenden Geschäftstransaktion
(Bestell- und Fulfillment-Prozess) oder die Meldung, dass ein clientseitiger Batch-Lauf abgeschlossen ist.
Solche Ereignisse führen Provider-seitig zu Dateneinfügungen, die für den Client jedoch nicht sichtbar
werden.

## Problem

Wie kann ein API-Provider seinen Clients erlauben zu melden, dass etwas geschehen ist, das der Provider
wissen muss — etwa um sofortige oder spätere Verarbeitung anzustoßen?

## Forces

- **Coupling-Trade-off:** Genauigkeit und Ausdrucksstärke der Meldung gegen Informationssparsamkeit.
  Je mehr Kontext im Event steckt, desto stärker koppelt es an das Provider-Domänenmodell.
- **Timing:** Wann wird verarbeitet — synchron beim Eintreffen oder verzögert? Die Operation garantiert
  Annahme, nicht Verarbeitung.
- **Konsistenzeffekte:** Der Client sieht das Ergebnis der Einfügung nicht; Lesen-nach-Schreiben ist
  nicht Teil des Kontrakts.
- **Zuverlässigkeit:** Verlorene oder doppelte Meldungen. Duplikaterkennung erfordert zwar Lesezugriff,
  ändert aber nichts am schreibenden Charakter der Operation.

## Lösung

Dem Endpunkt — einer [Processing Resource](ProcessingResource.md) oder einer
[Information Holder Resource](InformationHolderResource.md) — wird eine Operation der Signatur
`sco: in -> (out, S')` hinzugefügt. „Write-only“ ist dabei als Verantwortlichkeit zu lesen, nicht als
technisches Verbot: Die Operation darf lesen (etwa um Schlüsselkollisionen zu prüfen), ihr Zweck
bleibt aber das Erzeugen von Zustand. Die Antwort `out` ist typischerweise eine Quittung.

## Varianten

- **Event Notification Operation:** Der Client meldet ein fachliches Ereignis („order placed“,
  „product XYZ created“). Diese Variante ist die Brücke zu *Domain Events* aus Domain-Driven Design
  und zu Event-Sourcing-/CQRS-Architekturen.

## Beispiel

MDSL-Dekorator aus der Quelle:

```
"STATE_CREATION_OPERATION" @PaperItemDTO
createPaperItem (String who, String what, String where);
```

## Konsequenzen

**Vorteile:**

- Minimale Kopplung an Provider-Interna: Der Client braucht kein Wissen über Zustandsmodelle oder
  Vorbedingungen; er meldet nur, was bei ihm geschehen ist.
- Gut asynchron implementierbar (Message Queue, Event-Driven Consumer) und dadurch entkoppelt in
  Verfügbarkeit und Last.
- Klare Trennung von Kommando und Abfrage (Command-Query Separation auf Architekturebene).

**Nachteile / Kosten:**

- Der Client erfährt nichts über den Verarbeitungserfolg; Fehlerbehandlung muss über andere Kanäle
  laufen (Callbacks, Statusabfragen, Reklamationsprozesse).
- Idempotenz muss explizit entworfen werden, sonst erzeugen Retries Duplikate.
- Ohne Rückgabe eines Identifikators ist der erzeugte Zustand später schwer referenzierbar.

## Bekannte Verwendungen

- `submitReport` des Handling Report Service in der DDD Sample Application (Cargo-Tracking).
- Die Slack Event API als umfassende Umsetzung API-basierter Ereignisverarbeitung.
- Terravis: Banken starten Prozessinstanzen über eine „start process“-Operation.
- Ein Schweizer Banking-Softwarehaus meldet Änderungen an Geschäftsobjekten per Event an interessierte
  Microservices.

## Verwandte Patterns

- [State Transition Operation](StateTransitionOperation.md) — Schwesterpattern; liest *und* schreibt
  Zustand und referenziert dabei üblicherweise ein bestehendes Zustandselement (etwa eine Order-Id).
- [Retrieval Operation](RetrievalOperation.md) — komplementär: zieht Daten, statt sie zu schieben.
- [Computation Function](ComputationFunction.md) — erhält ebenfalls alle Daten vom Client, berührt aber
  keinerlei Provider-Zustand.
- [Processing Resource](ProcessingResource.md) und
  [Information Holder Resource](InformationHolderResource.md) — die Endpunkt-Rollen, die Instanzen
  dieses Patterns enthalten können.
- [Operational Data Holder](OperationalDataHolder.md) — kurzlebige, häufig geschriebene Daten sind der
  typische Zielspeicher solcher Ereignismeldungen.

## Bezug zu Kubernetes / KRM

Das Kubernetes Resource Model kennt ein festes, ressourcenunabhängiges Verbset:
`get`, `list`, `watch`, `create`, `update`, `patch`, `delete`, `deletecollection` (siehe die Abbildung
von HTTP-Methoden auf Verbnamen in `staging/src/k8s.io/apiserver/pkg/endpoints/installer.go`). Die
*State Creation Operation* bildet darin exakt auf `create` ab — POST auf eine Collection. Weil das
Verbset uniform ist, gibt es keine benannten fachlichen Anlege-Operationen wie `createPaperItem`;
die Fachlichkeit steckt in `Kind` und Schema, nicht im Operationsnamen.

Die Event-Notification-Variante hat in KRM eine wörtliche Entsprechung: `Event`-Objekte
(`events.k8s.io`) werden von Kubelet, Scheduler und Controllern per `create` gemeldet, um zu berichten,
dass etwas geschehen ist — ohne dass der Sender eine fachliche Antwort erwartet. Auch
`SubjectAccessReview` & Co. nutzen `create`, sind aber gerade *keine* State Creation Operations, sondern
[Computation Functions](ComputationFunction.md) ohne Persistenz.

Zwei KRM-Spezifika weichen vom Pattern ab. Erstens ist `create` in Kubernetes nicht write-only im Sinne
der Antwort: Der API-Server gibt das vollständige, mit `metadata.uid`, `metadata.resourceVersion`,
`metadata.creationTimestamp` und Defaults angereicherte Objekt zurück. Das ist mehr als die „got it“-
Quittung des Patterns und macht `create` zugleich zur Quelle des künftigen Referenzschlüssels.

```yaml
# Request: POST /api/v1/namespaces/default/pods
apiVersion: v1
kind: Pod
metadata:
  generateName: web-        # <- der Client vergibt bewusst keinen Namen
  namespace: default
spec:
  # ... containers
```

```yaml
# Response: 201 Created — dasselbe Objekt, vom Server angereichert
apiVersion: v1
kind: Pod
metadata:
  name: web-7d9f2                               # <- generateName + Zufallssuffix
  generateName: web-
  namespace: default
  uid: a3f7c2e1-9b0d-4f2a-8c31-6d5e0b7a1f42     # <- ab hier der stabile Referenzschlüssel
  resourceVersion: "184203"
  creationTimestamp: "2026-08-04T10:00:00Z"
  # ... managedFields
spec:
  # ... containers, dazu Defaults: restartPolicy, schedulerName, tolerations, serviceAccountName
status:
  phase: Pending            # <- schon gesetzt, obwohl noch kein Controller gearbeitet hat
```

Zweitens löst KRM Idempotenz namensbasiert: Ein zweites `create` mit gleichem `metadata.name`
scheitert mit `409 Conflict` und einem `metav1.Status` als Body — was Retries unschädlich macht.
Wer das nicht will, nutzt `metadata.generateName` wie oben und akzeptiert bewusst serverseitig
vergebene Namen.

```json
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "pods \"web\" already exists",
  "reason": "AlreadyExists",
  "details": { "name": "web", "kind": "pods" },
  "code": 409
}
```

Bemerkenswert ist außerdem, dass `create` in KRM keine Verarbeitung *auslöst*, sondern nur gewünschten
Zustand deklariert. Ein `Job`-Objekt anzulegen ist syntaktisch dieselbe Operation wie ein `ConfigMap`
anzulegen; dass daraus Arbeit entsteht, ist Sache des level-triggered reconciliierenden Controllers, der
per `watch` auf die Erzeugung reagiert. Die im Pattern beschriebene Asynchronität ist damit nicht
optional, sondern architektonisch fest verdrahtet. Insgesamt erfüllt KRM das Pattern — aber mit
uniformem Verb statt fachlicher Operation und mit einer deutlich reichhaltigeren Antwort.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/operationResponsibilities/StateCreationOperation)
