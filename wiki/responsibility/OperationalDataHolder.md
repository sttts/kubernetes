---
title: Operational Data Holder
kategorie: Responsibility
unterkategorie: Information Holder Endpoint Types
quelle: https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/OperationalDataHolder
---

# Operational Data Holder

*a.k.a.* Transaction(al) Data Holder, Secondary Data Access and Modification

**Kurzform:** Eine [Information Holder Resource](InformationHolderResource.md) für kurzlebige,
häufig geänderte Daten des Tagesgeschäfts, die überwiegend *auf* langlebigere Daten verweisen.

## Kontext

Ein Domänenmodell oder Glossar liegt vor, und es ist entschieden, einige Entitäten als
*Information Holder Resources* zu exponieren. Die Datenspezifikation zeigt, dass Lebensdauern und
Änderungszyklen stark differieren — von Sekunden und Minuten bis zu Jahren und Jahrzehnten — und
dass die schnell veränderlichen Entitäten in Beziehungen zu den langsameren stehen. Schnelle Daten
sind meist Link-*Quelle*, langsame meist Link-*Ziel*.

## Problem

Wie unterstützt eine API Clients, die Instanzen von Domänenentitäten anlegen, lesen, ändern und
löschen wollen, die *operative Daten* darstellen: eher kurzlebig, im Tagesgeschäft häufig geändert
und mit vielen ausgehenden Beziehungen?

## Forces

- **Verarbeitungsgeschwindigkeit** für Lese- und Änderungsoperationen — das Änderungsvolumen ist hoch.
- **Fachliche Agilität und Schema-Flexibilität** — operative Daten folgen dem Prozess und ändern
  ihre Form häufiger als Stammdaten.
- **Konzeptionelle Integrität und Konsistenz der Beziehungen** — die vielen ausgehenden Referenzen
  müssen auch dann gültig bleiben, wenn die Ziele sich unabhängig entwickeln.

## Lösung

Markiere eine [Information Holder Resource](InformationHolderResource.md) als *Operational Data
Holder* und statte sie mit Operationen aus, die häufiges und schnelles Anlegen, Lesen, Ändern und
Löschen erlauben.

Optional erhält sie fachliche Zusatzverantwortung. Ein Warenkorb kann etwa Gebühren- und
Steuerberechnung, Preisänderungsbenachrichtigungen und Rabattierung anbieten — also
[State Transition Operations](StateTransitionOperation.md) über das reine CRUD hinaus.

## Beispiel

In Lakeside Mutual sind Schadensmeldungen (`claims`) und Risikoeinschätzungen operative Daten:
Angebote referenzieren Verträge und Kunden, Verträge referenzieren Kunden. Für die Nachrichten
gelten alle Struktur-Patterns — das Einlegen eines Artikels in den Warenkorb erwartet einen
[Parameter Tree](../structure/ParameterTree.md) und liefert ein
[Atomic Parameter](../structure/AtomicParameter.md) als Erfolgskennzeichen zurück; der Checkout
erwartet einen [Parameter Forest](../structure/ParameterForest.md) und liefert Bestellnummer und
Liefertermin als [Atomic Parameter List](../structure/AtomicParameterList.md); das Löschen wird durch
ein [Id Element](../structure/IdElement.md) ausgelöst und beantwortet ein Erfolgskennzeichen bzw.
einen [Error Report](../structure/ErrorReport.md).

## Konsequenzen

**Vorteile:**

- Die explizite Klassifikation trennt Zugriffs- und Änderungsprofile: operative Endpunkte dürfen
  schneller evolvieren, weil weniger auf sie verweist.
- Schreiblast und Caching-Strategie lassen sich getrennt von den Stammdaten optimieren.

**Nachteile / Kosten:**

- Die vielen ausgehenden Referenzen erzeugen Laufzeitkopplung an
  [Master Data Holder](MasterDataHolder.md)-Endpunkte — ausgehende Aufrufe brauchen Absicherung.
- Klassifikationsfragen sind nicht immer eindeutig (siehe Konten vs. Kontobewegungen unten).

Für die Absicherung empfiehlt die Quelle Microservices-Infrastruktur-Patterns nach Nygard: ein
*Circuit Breaker* auf den ausgehenden Aufrufen vom *Operational Data Holder* zum
[Master Data Holder](MasterDataHolder.md), ein gemeinsames *Bulkhead* für eng gekoppelte
*Operational Data Holders*, die zusammen ein DDD-*Aggregate* bilden.

## Bekannte Verwendungen

- Die DDD-Beispielanwendung „Cargo Tracking“ mit `Cargo`, `RouteSpecification` und `Itinerary`
  (Ankünfte als Handling Events).
- Die „conversations“-Methodenfamilie der Slack Web API.
- Tweets und Posts in sozialen Netzwerken, sofern über [Public APIs](../foundation/PublicAPI.md)
  exponiert.
- Terravis: Zahlungsversprechen und Gläubigerentlassung bei der Übertragung eines Darlehens zwischen
  zwei Banken — beide existieren nur bis zum Abschluss der Transaktion.
- Kernbanken-APIs: Kontobewegungen sind operative Daten, während die Konten selbst eher Stammdaten
  sind — ein Beispiel dafür, dass die Grenze fachlich gezogen werden muss.

## Verwandte Patterns

- [Information Holder Resource](InformationHolderResource.md) — das Oberpattern.
- [Master Data Holder](MasterDataHolder.md) und [Reference Data Holder](ReferenceDataHolder.md) —
  die Alternativen für längerlebige Daten mit mehr eingehenden Referenzen.
- [Processing Resource](ProcessingResource.md) — die weniger daten-, mehr aktionsorientierte Alternative.
- [State Creation Operation](StateCreationOperation.md), [State Transition Operation](StateTransitionOperation.md),
  [Retrieval Operation](RetrievalOperation.md) — alle Operationsverantwortlichkeiten sind hier zulässig.
- [Embedded Entity](../quality/EmbeddedEntity.md) — für Referenzen auf andere *Operational Data Holders*.
- [Linked Information Holder](../quality/LinkedInformationHolder.md) — für Referenzen auf
  [Master Data Holder](MasterDataHolder.md); diese werden typischerweise *nicht* eingebettet.
- [Pagination](../quality/Pagination.md) — für die oft großen Ergebnismengen.

## Bezug zu Kubernetes / KRM

KRM kennt die Unterscheidung operativ/Stamm/Referenz **nicht als Modellkonzept** (Fall c auf
Schema-Ebene), realisiert sie aber sehr wohl faktisch — an anderer Stelle: nicht pro Ressourcentyp,
sondern pro `spec`/`status`-Hälfte und über den Speicherpfad.

Klare *Operational Data Holder* im Cluster sind:

- **`coordination.k8s.io/v1` `Lease`** — der prototypische Fall. Die Node-Heartbeats
  (`kube-node-lease`-Namespace) und alle Leader Elections schreiben im Sekunden- bis
  Zehn-Sekunden-Takt `spec.renewTime` fort; dazu kommen `spec.holderIdentity`,
  `spec.leaseDurationSeconds`, `spec.acquireTime`, `spec.leaseTransitions` und für die
  koordinierte Leader Election `spec.strategy` und `spec.preferredHolder`. Ein `status` fehlt
  vollständig — das Objekt *ist* der beobachtete Zustand. Die Lease-Objekte wurden 2018/19 gerade deshalb
  eingeführt, um die hochfrequenten Heartbeats aus dem großen, viel referenzierten `Node`-Objekt
  herauszuziehen — eine lehrbuchhafte Trennung von operativen und Stammdaten.

  ```yaml
  apiVersion: coordination.k8s.io/v1
  kind: Lease
  metadata:
    name: node-1
    namespace: kube-node-lease
    # ... ownerReferences auf den Node; resourceVersion steigt im Sekundentakt
  spec:
    holderIdentity: node-1
    leaseDurationSeconds: 40
    renewTime: "2026-08-04T10:00:03.412000Z"   # <- das einzige Feld, das sich ständig ändert
    # ... acquireTime, leaseTransitions
  # kein status: an einem reinen Heartbeat gibt es nichts zu beobachten
  ```

- **`Event`** (`core/v1` bzw. `events.k8s.io/v1`) — per Definition kurzlebig; der API-Server löscht
  sie nach `--event-ttl` (Standard eine Stunde), und `series` fasst Wiederholungen zusammen, um die
  Schreiblast zu begrenzen. Events referenzieren über `involvedObject` bzw. `regarding` fast immer
  auf langlebigere Objekte — die vom Pattern beschriebene Richtung „viele ausgehende Referenzen“.

  ```yaml
  apiVersion: events.k8s.io/v1
  kind: Event
  metadata:
    name: web-7d9f.17f2a1c0d4e5b6a7
    namespace: default
  eventTime: "2026-08-04T10:00:00.000000Z"     # Felder liegen flach: kein spec, kein status
  type: Warning
  reason: BackOff
  action: Pulling
  note: 'Back-off pulling image "nginx:1.27"'
  reportingController: kubernetes.io/kubelet
  reportingInstance: kubelet-node-1
  regarding:                                   # <- ausgehende Referenz auf das langlebigere Objekt
    apiVersion: v1
    kind: Pod
    name: web-7d9f
    namespace: default
    uid: a3f7c2e1-9b0d-4f2a-8c31-6d5e0b7a1f42
  series:
    count: 12                                  # <- verdichtet statt zwölf Einzelobjekte zu schreiben
    lastObservedTime: "2026-08-04T10:04:31.000000Z"
  ```

- **`status`-Subresources allgemein** — `pods/status`, `nodes/status`, `deployments/status`. Der
  `status`-Teil eines Objekts ist die operative Hälfte: hochfrequent geschrieben, von niemandem
  referenziert, jederzeit rekonstruierbar. Dass er eine eigene RBAC-*Ressource* ist — `pods/status`
  steht als eigener Eintrag unter `resources`, mit den Verben `get`/`update`/`patch` — ist die
  KRM-Antwort auf das Force „unterschiedliche Zugriffsprofile“.
- **`EndpointSlice`** und `resource.k8s.io`-Objekte wie `ResourceClaim` sind ebenfalls operativ und
  verweisen auf langlebigere Pods, Services und Klassen.

KRM löst die Forces anders als MAP: statt ausgehende Aufrufe zwischen Endpunkten mit Circuit
Breakers abzusichern, gibt es *keine* synchronen Aufrufe zwischen Ressourcen. Referenzen sind
Namensstrings bzw. `ObjectReference`s, die ein Controller asynchron über seinen Informer-Cache
auflöst; ein nicht auflösbares Ziel führt zu einem Requeue und einer `status.conditions`-Meldung,
nicht zu einem Fehler beim Client. Level-Triggered Reconciliation ersetzt den Circuit Breaker.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/OperationalDataHolder)
