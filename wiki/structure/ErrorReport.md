---
title: Error Report
kategorie: Structure
unterkategorie: Special Purpose Representations
quelle: https://microservice-api-patterns.org/patterns/structure/specialPurposeRepresentations/ErrorReport
---

# Error Report

*a.k.a.* API Fault Management, Error Code and Message

**Kurzform:** Fehler werden als eigenständige, maschinenlesbare Repräsentation in der
Antwort-Payload transportiert — Code plus Klassifikation plus menschenlesbarer Text —, statt sich
allein auf Protokoll-Statuszeilen zu verlassen.

## Kontext

Verteilte Kommunikationsteilnehmer müssen unerwartete Situationen zur Laufzeit zuverlässig
behandeln. Ein Client ruft eine Operation auf, der Provider kann sie nicht erfolgreich ausführen —
wegen falscher Requestdaten, ungültigen Serverzustands, fehlender Rechte oder Problemen in der
Infrastruktur. Die Ursache kann beim Client, beim Provider, in dessen Backend oder im Transport
liegen.

## Problem

Wie informiert ein API-Provider seine Clients über Kommunikations- und Verarbeitungsfehler? Und wie
lässt sich diese Information unabhängig von der zugrundeliegenden Kommunikationstechnologie und
Plattform machen, also etwa unabhängig von protokollspezifischen Status-Headern?

## Forces

- **Ausdrucksstärke und Zielgruppe** — Entwickler brauchen andere Information als Endnutzer oder
  Administratoren; ein Text bedient selten beide.
- **Robustheit und Verlässlichkeit** — der Client muss programmatisch entscheiden können, ob ein
  Retry sinnvoll ist. Das erfordert stabile, dokumentierte Codes.
- **Sicherheit und Performance** — detaillierte Fehler helfen beim Debugging, verraten aber auch
  interne Struktur. Stacktraces sind teuer und gefährlich.
- **Interoperabilität und Portabilität** — der Fehlerbericht soll auch überleben, wenn HTTP durch
  ein anderes Protokoll ersetzt oder ein Gateway dazwischengeschaltet wird.
- **Internationalisierung** — Codes sind sprachneutral, Texte nicht.

## Lösung

In der Antwort werden Fehlercodes zurückgegeben, die den Fehler einfach und maschinenlesbar
klassifizieren; ergänzt um textuelle Beschreibungen für die menschlichen Stakeholder des Clients
(Entwickler, Administratoren). Der Fehlerbericht ist eine eigene Struktur in der Payload und nicht
nur ein HTTP-Statuscode.

## Beispiel

Ein fehlgeschlagener Login bei Lakeside Mutual (Spring-Defaults):

```http
HTTP/1.1 401
Content-Type: application/json;charset=UTF-8

{
  "timestamp" : "2018-06-20T08:25:10.212+0000",
  "status" : 401,
  "error" : "Unauthorized",
  "message" : "Access Denied",
  "path" : "/auth"
}
```

RFC 7807 (`application/problem+json`) standardisiert dieselbe Idee mit `type` (URI),
`title` (Fehlerkategorie), `status` und `detail`:

```http
HTTP/1.1 400 Bad Request
Content-Type: application/problem+json

{
  "type": "https://myshop.org/out-of-stock",
  "title": "Out of Stock",
  "status": 400,
  "detail": "Item C330 is out of stock"
}
```

## Konsequenzen

**Vorteile:**

- Der Client kann Fehler programmatisch unterscheiden und differenziert reagieren (Retry,
  Nutzerdialog, Abbruch).
- Der Fehlerbericht überlebt Protokollwechsel und Gateways, weil er in der Payload steht.
- Ein einheitliches Fehlerformat über alle Endpoints senkt die Lernkurve deutlich.

**Nachteile / Kosten:**

- Der Fehlerraum wird Teil des API-Vertrags und damit evolutionspflichtig; ein Code, auf den Clients
  reagieren, lässt sich kaum noch entfernen.
- Detailtiefe steht im direkten Konflikt mit Informationspreisgabe.
- Redundanz zum Protokollstatus: Statuscode und Payload-Code können auseinanderlaufen.
- Übersetzte Texte und maschinenlesbare Codes müssen konsistent gepflegt werden.

## Bekannte Verwendungen

- **RFC 7807 / RFC 9457** als generische Standardisierung für HTTP-APIs, ergänzt durch die
  JSON-API-Spezifikation mit ihrem `errors`-Wurzelelement.
- **Twitter**: `{"errors":[{"message":"Sorry, that page does not exist","code":34}]}`.
- **Facebook Graph API**: `error`-Objekt mit `message`, `type`, `code` und `fbtrace_id`; ein
  Debug-Parameter im Query-String erhöht den Detailgrad.
- **JIRA Cloud API** (per JSON Schema definierte Fehler-Bodies), **Microsoft Graph API** mit
  API-spezifischen Codes wie `maxQueryLengthExceeded`, **Kloudless** mit einem vereinheitlichten
  Fehlerformat über heterogene Backends.
- Eine Core-Banking-Integration (Brandner et al. 2004) transportiert Fehlerobjekte als Teil einer
  standardisierten [Context Representation](ContextRepresentation.md), per XML Schema
  API-weit vereinheitlicht.

## Verwandte Patterns

- [Context Representation](ContextRepresentation.md) — ein Error Report kann darin eingebettet
  sein, etwa als Sammelbericht für ein [Request Bundle](../quality/RequestBundle.md).
- [Metadata Element](MetadataElement.md) — der Bericht besteht aus Metadaten, nicht aus Fachdaten;
  er kann Hinweise auf nächste Schritte enthalten.
- [Atomic Parameter List](AtomicParameterList.md) / [Parameter Tree](ParameterTree.md) — typische
  Strukturformen des Fehlerberichts (flach bzw. mit `causes`-Unterstruktur).
- [Request Bundle](../quality/RequestBundle.md) — braucht Teilfehler pro Bundle-Element.
- [Rate Limit](../quality/RateLimit.md) — meldet Überschreitungen über einen Error Report.
- [API Description](../foundation/APIDescription.md) — der Fehlerraum gehört in den Vertrag.

Außerhalb von MAP: *Remoting Error* (Voelter et al. 2004) als Middleware-nahes Pendant,
*Circuit Breaker* (Nygard) und *Dead Letter Channel* (Hohpe/Woolf) als ergänzende Bausteine.

## Bezug zu Kubernetes / KRM

Kubernetes erfüllt dieses Pattern **sehr sauber und ungewöhnlich konsequent** — es hat einen
einzigen, uniformen Error Report für *alle* Ressourcen, Gruppen und Versionen:
`metav1.Status` aus `k8s.io/apimachinery/pkg/apis/meta/v1`. Jede fehlgeschlagene Anfrage an den
kube-apiserver liefert ein Objekt mit `kind: Status` und den Feldern `status` (`"Success"` oder
`"Failure"`), `message` (menschenlesbar), `reason` (maschinenlesbar), `details` und `code` (der
vorgeschlagene HTTP-Status). Bemerkenswert: Der Fehlerbericht ist selbst ein reguläres
KRM-Objekt mit `apiVersion`/`kind` und wird durch dieselbe Serialisierungs-Pipeline (JSON, YAML,
Protobuf) geschickt wie Fachobjekte — Protokollunabhängigkeit ergibt sich hier nicht durch
Sonderbehandlung, sondern weil Fehler einfach ein Typ unter vielen sind.

`reason` ist eine geschlossene Konstantenliste vom Typ `StatusReason`: `NotFound`, `AlreadyExists`,
`Conflict`, `Invalid`, `Forbidden`, `Unauthorized`, `Gone`, `Expired`, `Timeout`, `ServerTimeout`,
`TooManyRequests`, `RequestEntityTooLarge`, `UnsupportedMediaType`, `MethodNotAllowed`,
`ServiceUnavailable`, `InternalError` und einige weitere. Genau das, was das Pattern unter
"maschinenlesbare Klassifikation" fordert — und was HTTP-Statuscodes allein nicht leisten, weil
etwa `409 Conflict` sowohl `AlreadyExists` als auch einen `resourceVersion`-Konflikt bedeuten kann.

```http
HTTP/1.1 404 Not Found
Content-Type: application/json

{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure",
 "message":"deployments.apps \"web\" not found",
 "reason":"NotFound",
 "details":{"group":"apps","kind":"deployments","name":"web"},
 "code":404}
```

```http
# Derselbe Statuscode kann zwei verschiedene Dinge heißen — erst `reason` entscheidet:
HTTP/1.1 409 Conflict
Content-Type: application/json

{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure",
 "message":"Operation cannot be fulfilled on deployments.apps \"web\": the object has been modified; please apply your changes to the latest version and try again",
 "reason":"Conflict",
 "details":{"group":"apps","kind":"deployments","name":"web"},
 "code":409}
```

`details` (`StatusDetails`) trägt `group`, `kind`, `name`, `uid`, `retryAfterSeconds` und vor allem
`causes` — eine Liste von `StatusCause` mit `reason` (`CauseType`), `message` und `field`. Das
`field` verwendet JSON-Pfadnotation inklusive Array-Index, etwa `spec.containers[0].image`.
Validierungsfehler sind damit feldgenau und clientseitig auf UI-Felder abbildbar; das ist deutlich
mehr, als RFC 7807 im Kern anbietet. `apierrors.NewInvalid` erzeugt genau diese Struktur mit
`code: 422` und `reason: Invalid`.

```http
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json

{
  "kind": "Status", "apiVersion": "v1", "metadata": {},
  "status": "Failure",
  "message": "Deployment.apps \"web\" is invalid: [spec.replicas: Invalid value: -1: must be greater than or equal to 0, spec.template.spec.containers[0].image: Required value]",
  "reason": "Invalid",
  "details": {
    "group": "apps",
    "kind": "Deployment",
    "name": "web",
    "causes": [
      { "reason": "FieldValueInvalid",
        "message": "Invalid value: -1: must be greater than or equal to 0",
        "field": "spec.replicas" },
      { "reason": "FieldValueRequired",
        "message": "Required value",
        "field": "spec.template.spec.containers[0].image" }
    ]
  },
  "code": 422
}
```

Die `field`-Pfade sind direkt auf Eingabefelder einer UI abbildbar; `message` in der Wurzel ist die
für Menschen aggregierte Fassung derselben Information.

Auf Client-Seite existiert dazu eine vollständige Prädikatenbibliothek in
`k8s.io/apimachinery/pkg/api/errors`: `IsNotFound()`, `IsAlreadyExists()`, `IsConflict()`,
`IsInvalid()`, `IsForbidden()`, `IsTooManyRequests()`, `IsResourceExpired()` und so weiter. Das
`if apierrors.IsNotFound(err) { ... }` in Controllern ist das direkte Konsumieren des Error Reports
und ein Kernbaustein level-getriggerter Reconciliation: `NotFound` heißt "Objekt weg, Zielzustand
anpassen", `Conflict` heißt "Optimistic-Concurrency-Konflikt, requeue", `TooManyRequests` heißt
"Backoff nach `retryAfterSeconds`". Der Fehlerbericht ist hier also nicht nur Diagnose, sondern
Steuerinformation für die Regelschleife.

```go
deployment, err := c.lister.Deployments(ns).Get(name)
if apierrors.IsNotFound(err) {
    return nil            // <- Objekt weg: nichts zu tun, kein Requeue
}
// ...
if _, err := c.client.AppsV1().Deployments(ns).UpdateStatus(ctx, deployment, metav1.UpdateOptions{}); err != nil {
    return err            // <- IsConflict: Workqueue requeued mit Backoff, dann neu lesen
}
```

Zwei Ergänzungen jenseits des Patterns: Erstens gibt es mit `Warning`-Headern (HTTP 299) einen
Kanal für nicht-fatale Hinweise, etwa auf deprecated API-Versionen — Erfolg und Beanstandung
zugleich. Zweitens ist bei asynchronen, deklarativen Operationen der eigentliche "Fehlerbericht"
gar nicht die HTTP-Antwort: Ein `kubectl apply` kann `201 Created` liefern, während das
Fachproblem erst später in `status.conditions` oder als `Event` erscheint.

```yaml
# 201 Created war die Antwort — der eigentliche Fehler steht Sekunden später hier:
status:
  observedGeneration: 7
  conditions:
    - type: Progressing
      status: "False"
      reason: ProgressDeadlineExceeded        # <- maschinenlesbar, wie `reason` in Status
      message: ReplicaSet "web-6d4f8b7c9" has timed out progressing.
      lastTransitionTime: "2026-08-04T10:12:00Z"
  # ... replicas, readyReplicas
```

KRM hat damit *zwei* Error-Report-Ebenen — eine synchrone für den API-Aufruf und
eine level-getriggerte für die Konvergenz —, was im ursprünglichen Request/Response-Bild des
Patterns nicht vorgesehen ist.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/specialPurposeRepresentations/ErrorReport)
