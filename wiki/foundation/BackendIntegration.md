---
title: Backend Integration
kategorie: Foundation
unterkategorie: Integrationsrichtung
quelle: https://microservice-api-patterns.org/patterns/foundation/BackendIntegration
---

# Backend Integration

*a.k.a.* Horizontal Integration, Backend-to-Backend Integration, Remote Component API, East-West Connectivity

**Kurzform:** Unabhängig gebaute und getrennt deployte Systemteile tauschen über nachrichtenbasierte
Remote-APIs Daten aus und stoßen gegenseitig Aktivität an, ohne ihre begriffliche Integrität aufzugeben
oder unerwünschte Kopplung einzuführen.

## Kontext

Auf Unternehmensebene wurde entschieden, eine Systemlandschaft in mehrere verteilte Systeme bzw. Services
zu zerlegen. Diese müssen kommunizieren, um Ende-zu-Ende-Anwendungsfälle zu realisieren. Wegen der
physischen Trennung scheiden lokale Aufrufe aus. Die Quelle stellt fest, dass nachrichtenbasiertes
Remoting heute gegenüber Remote-Objekt-Technologie bevorzugt wird (Hohpe/Woolf 2003).

## Problem

Wie tauschen unabhängig gebaute, getrennt deployte verteilte Anwendungen und ihre Teile Daten aus und
lösen gegenseitig Aktivität aus — ohne die systeminterne begriffliche Integrität zu verlieren und ohne
unerwünschte Kopplung? Kurz: wie kommunizieren Microservices miteinander?

## Forces

- **Laufzeitqualitäten** wie Performance und Skalierbarkeit — bei East-West-Verkehr oft dominierend,
  weil ein Aufruf viele nachgelagerte Aufrufe auslöst.
- **Sicherheit** — auch interner Verkehr braucht Authentifizierung und Autorisierung.
- **Interoperabilität** — heterogene Middleware, Sprachen und Serialisierungsformate.
- **Entwicklungsbudgets** — Integration ist Aufwand ohne unmittelbar sichtbaren Fachwert.
- **Entwicklungskultur und Unternehmenspolitik** — wer definiert das Schema, wer trägt Breaking Changes.

## Lösung

Das Backend einer verteilten Anwendung wird mit einem oder mehreren anderen Backends integriert, indem
es seine Dienste über eine nachrichtenbasierte Remote-*Backend Integration* API exponiert. Die Quelle
zeigt beide Richtungen im selben Bild: Frontend-to-Backend
([*Frontend Integration*](FrontendIntegration.md)) und Backend-to-Backend.

## Beispiel

Im Fallbeispiel *Lakeside Mutual* ist die Schnittstelle zwischen dem Customer-Self-Service-Backend und
dem Customer-Core-Microservice, der Kundenstammdaten anwendungsübergreifend bündelt, eine
*Backend Integration*. Die Backend-Services kommunizieren teils über HTTP-Resource-APIs, teils über
Message Queues; ein queue-basiertes Beispiel ist das Policy-Update-Messaging zwischen dem
`RiskManagementMessageProducer` im Policy Management und dem Risk-Management-Server über ActiveMQ.

## Konsequenzen

**Vorteile:**

- Systemteile können unabhängig entwickelt, deployt und skaliert werden.
- Fachliche Zuständigkeiten (Bounded Contexts) bleiben klar getrennt und wiederverwendbar.

**Nachteile / Kosten:**

- Kopplung verschwindet nicht, sie wandert in Schema und Protokoll — Schemaänderungen werden zum
  organisationsübergreifenden Vorgang.
- Teilausfälle, Retries, Idempotenz und Konsistenz über Servicegrenzen hinweg müssen explizit behandelt
  werden.
- Ein eigener Produkt- und Beratungsmarkt (*Enterprise Application Integration*) existiert genau wegen
  dieses Aufwands.

## Bekannte Verwendungen

Allgegenwärtig: betriebliche Informationssysteme werden typischerweise so integriert. Die Quelle nennt
die Trennung von Frontend- und Backend-APIs bei einer großen Schweizer Bank (Murer/Hagen 2014,
Murer/Bonati/Furrer 2010) sowie die *Open Service Broker API*, über die Cloud-Anbieter ihre Angebote in
Cloud-Plattformen einbinden.

## Verwandte Patterns

- [Frontend Integration](FrontendIntegration.md) — das Geschwister-Pattern, die andere Richtung.
- [Public API](PublicAPI.md), [Community API](CommunityAPI.md),
  [Solution-Internal API](SolutionInternalAPI.md) — die orthogonalen Sichtbarkeitsstufen.
- [API Description](APIDescription.md) — Grundlage dafür, dass Integrationsentwickler Adapter schreiben können.
- [Data Transfer Resource](../responsibility/DataTransferResource.md) — entkoppelter Datenaustausch
  zwischen Backends über eine vermittelnde Ressource.
- [Master Data Holder](../responsibility/MasterDataHolder.md) — der typische Zweck eines
  backend-integrierten Kernservices.
- [Two in Production](../evolution/TwoInProduction.md) — weil Backends nicht synchron deployt werden.

## Bezug zu Kubernetes / KRM

Die gesamte Steuerebene von Kubernetes ist Backend-zu-Backend-Integration, aber in einer Topologie, die
MAP so nicht beschreibt: **alle Komponenten integrieren ausschließlich über den kube-apiserver**, nicht
direkt miteinander. kube-controller-manager, kube-scheduler, kubelet, kube-proxy und beliebige
Custom-Controller lesen und schreiben Ressourcen; keiner ruft einen anderen auf. Der apiserver mit etcd
ist damit ein geteilter, schemabehafteter Zustandsspeicher, über den integriert wird — konzeptionell
näher an einem Blackboard als an Punkt-zu-Punkt-Messaging.

```text
MAP-Bild — jede Integration ist eine Kante zwischen zwei Komponenten:

    scheduler ───► kubelet ───► controller
        ▲                           │
        └───────────────────────────┘

    bis zu N*(N-1) Kanten, je eigenes Schema und Protokoll


KRM-Bild — keine Kante zwischen zwei Komponenten, nur Kanten zum apiserver:

    scheduler ─┐                ┌─ controller-manager
               ▼                ▼
            ┌────────────────────┐
            │   kube-apiserver   │ ──► etcd
            └────────────────────┘
               ▲                ▲
    kubelet ───┘                └─ custom-controller

    N Kanten, ein Schema und ein Protokoll für alle
```

Die Nachrichtenschicht ist **watch**: statt Request/Response oder einer Message Queue etablieren Clients
über `?watch=true` einen Stream von `ADDED`/`MODIFIED`/`DELETED`/`BOOKMARK`-Events, gepuffert in
Informern und Sharded Caches (`k8s.io/client-go/tools/cache`). Entscheidend ist die Semantik:
Level-Triggered Reconciliation. Ein verpasstes Event ist unkritisch, weil der Controller aus dem
Ist-Zustand neu ableitet — genau deshalb braucht KRM keine garantierte Zustellung, keine
Dead-Letter-Queues und keine Idempotenzschlüssel, die bei nachrichtenbasierter Integration üblich sind.
`resourceVersion` und optimistische Nebenläufigkeit (409 Conflict) ersetzen Transaktionskoordination.

```http
GET /apis/apps/v1/namespaces/default/deployments?watch=true&resourceVersion=482931&allowWatchBookmarks=true HTTP/1.1

HTTP/1.1 200 OK
Content-Type: application/json
Transfer-Encoding: chunked      # <- ein JSON-Objekt pro Frame, der Stream endet nie von selbst
```

```json
{"type":"MODIFIED","object":{"kind":"Deployment","apiVersion":"apps/v1",
  "metadata":{"name":"web","namespace":"default","generation":8,"resourceVersion":"482940"},
  "spec":{"replicas":5},
  "status":{"observedGeneration":7,"readyReplicas":3}}}
{"type":"BOOKMARK","object":{"kind":"Deployment","apiVersion":"apps/v1",
  "metadata":{"resourceVersion":"482944"}}}
// <- trägt nur die resourceVersion: Wiederaufsetzpunkt, kein Domain-Event
{"type":"DELETED","object":{"kind":"Deployment","apiVersion":"apps/v1",
  "metadata":{"name":"web","namespace":"default","resourceVersion":"482951"}}}
```

Ein verpasstes `MODIFIED` ändert nichts: der Controller liest `spec.replicas` und
`status.observedGeneration` im nächsten Durchlauf ohnehin neu. Genau deshalb ist der `BOOKMARK` kein
Zustellungsbeleg, sondern nur eine `resourceVersion` zum Wiederaufsetzen — bei `sendInitialEvents=true`
markiert der erste Bookmark zusätzlich mit der Annotation `k8s.io/initial-events-end: "true"` das Ende
des synthetischen Anfangszustands.

Die `Event`-Ressourcen (`core/v1`, `events.k8s.io/v1`) sind trotz des Namens **keine** Domain-Events zur
Integration, sondern kurzlebige Diagnoseinformation; sich für die Reconciliation darauf zu verlassen,
widerspricht den API-Konventionen.

Es gibt allerdings auch klassische Backend-Integration, bei der der apiserver selbst Client ist:
Admission- und Conversion-Webhooks (`admissionregistration.k8s.io`,
`apiextensions.k8s.io` `spec.conversion.webhook`) sind synchrone HTTP-Aufrufe nach außen, ebenso
`TokenReview`/`SubjectAccessReview` gegen externe Authentifizierer und die Weiterleitung an
aggregierte API-Server (`apiregistration.k8s.io/v1 APIService`). Außerhalb der KRM integriert das
kubelet über gRPC-Schnittstellen (CRI, CSI, Device Plugins) — dort gilt MAPs Bild von
Punkt-zu-Punkt-Backend-APIs unverändert.

Die Webhook-Konfiguration ist selbst wieder eine Ressource — die Integrationsbeziehung wird also
deklariert, nicht einprogrammiert, und trägt genau die Merkmale, die MAP von einer synchronen
Backend-API verlangt: Endpunkt, Version, Timeout, Fehlerverhalten und Seiteneffektzusage.

```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: policy.example.com
webhooks:
  - name: validate.policy.example.com
    clientConfig:
      service:
        namespace: policy-system
        name: policy-webhook
        path: /validate           # <- ab hier ist der apiserver Client, nicht Server
        port: 443
      # ... caBundle
    rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: ["apps"]
        apiVersions: ["v1"]
        resources: ["deployments"]
    admissionReviewVersions: ["v1"]
    sideEffects: None             # <- Vertragszusage, damit --dry-run erlaubt bleibt
    failurePolicy: Fail           # <- ist der Webhook nicht erreichbar, scheitern die passenden Writes
    timeoutSeconds: 10
```

Bewertung: KRM erfüllt das Pattern **anders**. Die Absicht (Entkopplung unabhängig deployter Teile) ist
erfüllt, das Mittel ist aber ein zentraler deklarativer Ressourcenspeicher mit Watch-Streams statt
direkter Service-zu-Service-Nachrichten.

---
[← Index](../README.md) · [Kategorie Foundation](../meta/category-foundation.md) · [Quelle](https://microservice-api-patterns.org/patterns/foundation/BackendIntegration)
