---
title: Context Representation
kategorie: Structure
unterkategorie: Special Purpose Representations
quelle: https://microservice-api-patterns.org/patterns/structure/specialPurposeRepresentations/ContextRepresentation
---

# Context Representation

*a.k.a.* Shared and Explicit Context, QoS Property Collection, Embedded Custom Header

**Kurzform:** Alle Kontext-Metadaten eines Aufrufs — Identität, QoS-Eigenschaften, Korrelations-IDs,
Präferenzen — werden zu *einem* klar markierten Element in der Nachrichten-Payload gebündelt, statt
sie auf Protokoll-Header zu verteilen.

## Kontext

Endpoint und Operationen sind definiert; zwischen Client und Provider muss Kontextinformation
fließen: Standort und Profildaten des Nutzers, Präferenzen im Sinne einer *Wish List*, oder
QoS-Kontrollen wie Credentials zur Authentifizierung, Autorisierung und Abrechnung — etwa
[API Keys](APIKey.md) oder JWT-Claims. Aufrufe sind oft Teil längerer *Conversations* aus mehreren
zusammenhängenden Operationen, wobei Provider selbst als Clients weiterer APIs auftreten. Ein Teil
des Kontexts ist lokal zum einzelnen Aufruf, ein anderer global und über die ganze Kette
weiterzureichen.

## Problem

Wie tauschen Client und Provider Kontextinformationen aus, ohne sich auf ein bestimmtes
Remoting-Protokoll festzulegen? Und wie werden Identitäts- und Qualitätsinformationen eines Requests
für nachfolgende Requests derselben Conversation sichtbar?

## Forces

- **Interoperabilität und Modifizierbarkeit auf technischer Ebene**, inklusive der Frage, ob Kontext
  zentral gehalten oder dezentral mitgeführt wird.
- **Abhängigkeit von sich entwickelnden Protokollen** — Header-Erweiterungsmechanismen sind
  protokollspezifisch und ändern sich.
- **Entwicklerproduktivität** — Kontrolle über die Struktur vs. Bequemlichkeit fertiger Header.
- **Diversität der Clients** — Browser, Rich Clients und Batch-Jobs brauchen unterschiedlichen
  Kontext.
- **Ende-zu-Ende-Sicherheit** über Service- und Protokollgrenzen hinweg; Gateways verschlucken oder
  überschreiben Header gern.
- **Logging und Auditing auf Fachebene** über mehrere Aufrufe hinweg.

## Lösung

Alle [Metadata Elements](MetadataElement.md), die den gewünschten Kontext tragen, werden zu einem
eigenen Repräsentationselement in Request und/oder Response gruppiert. Diese *Context
Representation* wird **nicht** in Protokoll-Headern transportiert, sondern in der Payload. Globaler
und lokaler Kontext einer Conversation werden strukturell getrennt. Das Element wird so positioniert
und markiert, dass es leicht auffindbar und klar von den fachlichen
[Data Elements](DataElement.md) unterscheidbar ist.

## Beispiel

Der MDSL-Vertragsauszug der Quellseite definiert eine `RequestContext`-Struktur, die als
gemeinsamer Kontext aller Operationen dient:

```
data type RequestContext {
    "apiKey":ID,
    "sessionId":D?,
    "qosPropertiesThatShouldNotGoToProtocolHeader":KeyValuePair*}

operation getCustomerAttributes
  expecting payload {
    > { "requestContextSharedByAllOperations": RequestContext,
        >"desiredCustomerAttributes":ID+ },
    > "searchParameters":D*
  }
```

Die `>`-Stereotypen markieren den Kontextblock, sodass er in der Payload sofort von den fachlichen
Daten zu unterscheiden ist. Der aus MDSL generierte OpenAPI-Auszug (Seite *Context Representation
Example*) bildet das auf ein wiederverwendbares Schema ab:

```yaml
components:
  schemas:
    RequestContext:
      type: object
      properties:
        apiKey: { type: string }
        sessionId: { type: integer, format: int32, nullable: true }
        qosPropertiesThatShouldNotGoToProtocolHeader:
          type: array
          items:
            $ref: '#/components/schemas/KeyValuePair'
```

Der Key-Value-Teil ist bewusst offen: Er ersetzt "zusätzliche freie Header", die man nicht mehr im
Protokoll unterbringen will. Die `sessionId` erzeugt der Provider nach erfolgreicher
Authentifizierung; der Client führt sie danach mit — der Conversation-Aspekt des Patterns.

Die Response desselben Beispiels kombiniert drei weitere Patterns:
`billingInfo`/`moreAnalytics` als Antwortkontext, `errorCode`/`errorMessage` als
[Error Report](ErrorReport.md) und `thisPageContent`/`previousPage`/`nextPage` als
[Pagination](../quality/Pagination.md); die `desiredCustomerAttributes` im Request sind eine
[Wish List](../quality/WishList.md).

## Konsequenzen

**Vorteile:**

- Protokollunabhängigkeit: Der Kontext überlebt Übergänge zwischen HTTP, SOAP, Messaging und
  proprietären Backend-Protokollen unverändert.
- Der Kontext ist typisiert, versionierbar und Teil der
  [API Description](../foundation/APIDescription.md) — Header sind das oft nicht.
- Zwischenknoten reichen ihn verlässlicher weiter, weil Gateways Payloads seltener anfassen.
- Eine Stelle für Auditing, Korrelation und Compliance-Informationen.

**Nachteile / Kosten:**

- Deutlich mehr Design- und Implementierungsaufwand als "einfach einen Custom-Header setzen".
- Infrastruktur (Caches, Proxies, Load Balancer, Tracing) kann nicht mehr mitlesen.
- Die Payload wächst bei jedem Aufruf um denselben Block; Kontext wird gerne zur Müllhalde.
- Duplizierung, wenn dieselbe Information auch protokollseitig vorhanden ist.
- Der Kontext wird Teil des Vertrags und damit evolutionspflichtig.

## Bekannte Verwendungen

Im Public-API-Bereich selten, in *Community APIs* und *Solution-Internal APIs* verbreitet:

- Die Dynamic Interface zu Core-Banking-Diensten (Brandner et al. 2004) bleibt so von heterogenen
  Frontends bis zum Mainframe-Backend protokollunabhängig — Kontext in Request *und* Response.
- Terravis (Lübke/van Lessen 2016): Jede SOAP-Nachricht trägt Message-ID, globale
  Geschäftsprozess-ID sowie Authentifizierungs- und Autorisierungsdaten.
- Eine große Schweizer Bank propagiert Ursprungs-Request-ID, Nutzer und Rolle, Zeitstempel,
  Ursprungssystem und Länderkennung — letztere zur Durchsetzung von Compartmentalization.
- **JWT** (RFC 7519): Claims als Metadata Elements über den eingeloggten Nutzer, kodiert in einen
  einzigen [Atomic Parameter](AtomicParameter.md).
- Ein Schweizer Versicherungssoftware-Hersteller definiert das Pattern in internen REST-Guidelines
  als Bibliothek wiederverwendbarer Repräsentationsfragmente (Filter, Sortierung, Wish-List-Einträge)
  — eine Nutzung jenseits von QoS.
- JMS-Nachrichten mit Header- und Payload-Teil gelten als Instanz, sofern sie nicht nativ vom Broker
  verarbeitet werden.

Als **Gegenbeispiel** nennt die Quelle das Rate-Limit-Beispiel von Lakeside Mutual: Es liefert die
Limit-Information in eigenen HTTP-Headern — weniger Aufwand, dafür Protokollbindung.

## Verwandte Patterns

- [Metadata Element](MetadataElement.md) — die Bausteine, die eine Context Representation bündelt.
- [Error Report](ErrorReport.md) — kann Teil der Context Representation sein, etwa als
  Zusammenfassung für ein [Request Bundle](../quality/RequestBundle.md).
- [API Key](APIKey.md) — typischer Inhalt des Kontexts.
- [Id Element](IdElement.md) — Session- und Korrelations-IDs im Kontext.
- [Wish List](../quality/WishList.md) und [Wish Template](../quality/WishTemplate.md) — Präferenzen,
  die ebenfalls Kontextcharakter haben.
- [Pagination](../quality/Pagination.md) — im Beispiel gemeinsam mit dem Kontext modelliert.
- [Parameter Tree](ParameterTree.md) — die übliche Strukturform des Kontextblocks.

Außerhalb von MAP: *Context Object* (Alur et al.), *Invocation Context* (Voelter et al. 2004),
*Envelope Wrapper* und *Wire Tap* (Hohpe/Woolf).

## Bezug zu Kubernetes / KRM

Hier weicht Kubernetes **deutlich vom Pattern ab — überwiegend im Sinne von (c) gar nicht**, und
das ist eine bewusste Entscheidung. KRM-Objekte enthalten keinen standardisierten Kontextblock in
der Payload. Es gibt kein `spec.requestContext`, keine Session-ID, keine QoS-Properties im Body.
Der Grund ist das uniforme Schema: Jede Ressource hat `apiVersion`, `kind`, `metadata`, meist `spec`
und `status` — und `metadata` ist für die *Identität des Objekts* reserviert, nicht für den Kontext
des Aufrufs. Ein Aufrufkontext in der Payload würde außerdem den deklarativen Charakter brechen:
Der Body beschreibt einen *gewünschten Zustand*, der idempotent immer wieder angewendet wird; eine
Session-ID oder ein Zeitstempel darin wäre semantisch falsch, weil sie beim nächsten Apply anders
lauten müsste.

Stattdessen lebt Kontext bei Kubernetes an drei anderen Orten. **Erstens im HTTP-Protokoll**, also
exakt dort, wo das Pattern ihn nicht haben will: `Authorization`, `User-Agent` (den der apiserver
für Audit und Metriken auswertet), `Accept` für Content-Negotiation, und die Impersonation-Header
`Impersonate-User`, `Impersonate-Group`, `Impersonate-Uid`, `Impersonate-Extra-<key>`. Dazu
Query-Parameter wie `dryRun=All`, `fieldManager`, `fieldValidation` und `resourceVersion` —
funktional Kontext, technisch Protokoll.

```http
# Der gesamte Aufrufkontext steht oberhalb der Leerzeile; der Body ist reiner Zielzustand.
PATCH /apis/apps/v1/namespaces/default/deployments/web?fieldManager=kubectl&fieldValidation=Strict&dryRun=All HTTP/1.1
Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ik...
Impersonate-User: alice                         # <- zweites Identitäts-Set
Impersonate-Group: developers
User-Agent: kubectl/v1.36.0 (darwin/arm64) kubernetes/abc1234
Accept: application/json
Content-Type: application/strategic-merge-patch+json
traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01

{"spec":{"replicas":3}}
```

**Zweitens im serverseitigen Request-Kontext.** Der apiserver reichert die Anfrage über eine
Filterkette an: `WithRequestInfo` legt aufgelöste Gruppe/Version/Ressource/Verb in den
Go-`context.Context`, Authentifizierungsfilter hängen den `user.Info` an, und
`audit.AddAuditAnnotation` erlaubt Komponenten (etwa Admission-Webhooks), strukturierte
Kontextannotationen an das Audit-Event zu heften. Das ist eine Context Representation im Geiste des
Patterns — nur ist sie prozessintern und wird nicht über die Leitung geschickt. Sichtbar wird sie
erst im Audit-Log, wo sie als eigenes KRM-Objekt materialisiert:

```json
{
  "kind": "Event",
  "apiVersion": "audit.k8s.io/v1",
  "level": "Metadata",
  "auditID": "3f7d1a20-9c44-4b2e-8f01-6de5a4c73b19",
  "stage": "ResponseComplete",
  "requestURI": "/apis/apps/v1/namespaces/default/deployments/web?fieldManager=kubectl",
  "verb": "patch",
  "user": { "username": "system:serviceaccount:ci:deployer",
            "groups": ["system:serviceaccounts"] },
  "impersonatedUser": { "username": "alice", "groups": ["developers"] },
  "userAgent": "kubectl/v1.36.0 (darwin/arm64) kubernetes/abc1234",
  "sourceIPs": ["10.0.2.15"],
  "objectRef": { "apiGroup": "apps", "apiVersion": "v1", "resource": "deployments",
                 "namespace": "default", "name": "web" },
  "responseStatus": { "code": 200 },
  "annotations": {
    "authorization.k8s.io/decision": "allow",
    "authorization.k8s.io/reason": "RBAC: allowed by RoleBinding \"deployer/default\""
  }
}
```

**Drittens in `metadata.annotations` und `metadata.labels`.** Annotationen sind der einzige
payload-seitige, generische Key-Value-Bereich und damit das KRM-Pendant zu den
`qosPropertiesThatShouldNotGoToProtocolHeader` des Beispiels; genutzt etwa für
`kubectl.kubernetes.io/last-applied-configuration` oder
`deprecated.daemonset.template.generation`.

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-agent
  annotations:                   # <- das KRM-Pendant zu den "qosProperties…" des MDSL-Beispiels
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"DaemonSet", ...}
    deprecated.daemonset.template.generation: "3"
  labels:
    environment: production      # <- anders als Payload-Kontext zusätzlich *abfragbar*:
    team: platform               #    kubectl get ds -l environment=production
```

Der entscheidende Unterschied: Sie sind Kontext des
*Objekts*, nicht des *Aufrufs*, überleben den Request und sind für alle Leser sichtbar. Labels plus
Selektoren machen Kontext (Umgebung, Team, Tier) zusätzlich *abfragbar* — was eine
Payload-Context-Representation nicht leistet.

Ehrlich bewertet: Kubernetes standardisiert Fachkontext bewusst kaum. Es gibt keinen definierten
Ort für "Mandant", "Region", "Auftragsnummer" oder "Korrelations-ID" — diese Felder existieren
schlicht nicht im API-Vertrag; Konventionen dafür sind Sache der jeweiligen Erweiterung, und
Distributed Tracing nutzt W3C-`traceparent`-Header, also wieder Protokoll statt Payload. Für den
eigentlichen Anwendungsfall des Patterns — protokollunabhängige Kontextweitergabe über
Conversations — besteht in KRM auch kein Bedarf: Es gibt keine Conversations. Level-getriggerte
Reconciliation ersetzt Aufrufketten durch wiederholtes Abgleichen des Zielzustands, und dieser
Zustand steht ohnehin schon im Cluster.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/specialPurposeRepresentations/ContextRepresentation)
