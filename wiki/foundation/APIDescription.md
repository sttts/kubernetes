---
title: API Description
kategorie: Foundation
unterkategorie: API-Dokumentation
quelle: https://microservice-api-patterns.org/patterns/foundation/APIDescription
---

# API Description

*a.k.a.* API Documentation, Explicit Service Contract

**Kurzform:** Provider und Clients einigen sich auf ein explizites, dokumentiertes Wissen über die API —
Nachrichtenstrukturen, Fehlermeldung, Verhalten, Qualitäten und organisatorische Rahmenbedingungen.

## Kontext

Ein Service-Provider hat entschieden, eine oder mehrere API-Operationen in einem Endpunkt zu exponieren.
Client-Entwickler — Web- und Mobile-Entwickler für [Frontend Integration](FrontendIntegration.md),
Systemintegratoren für [Backend Integration](BackendIntegration.md) — können die Aufrufe noch nicht
codieren und wissen nicht, was in den Antworten zu erwarten ist. Auch ergänzende Erklärungen fehlen:
Bedeutung der Operationen und Parameter, Wirkung auf den Anwendungszustand, Qualitäten wie Idempotenz
und Transaktionalität.

## Problem

Welches Wissen muss zwischen Provider und Clients geteilt werden, und wie soll es dokumentiert werden?
Präziser:

- Wie machen Client und Provider ihre Einigung über die *funktionalen* Aspekte explizit —
  Transferrepräsentationen, Aufrufvoraussetzungen?
- Wie wird diese funktionale Information um weitere technische Spezifikationselemente (Protokollheader,
  Security-Policies, Fehlerdatensätze) und geschäftliche Dokumentation (Aufrufsemantik, API-Eigentümer,
  Abrechnung, Supportprozesse, Versionierung) ergänzt?

## Forces

- **Interoperabilität** — Clients und Provider laufen auf unterschiedlichen Middleware-Plattformen.
- **Consumability** — Verständlichkeit, Erlernbarkeit, Einfachheit.
- **Information Hiding** — Implementierungsdetails dürfen nicht durchschlagen.
- **Erweiterbarkeit und Evolvierbarkeit** als Facetten allgemeiner Änderbarkeit.

Diese Kräfte ziehen gegeneinander: eine vollständige Beschreibung erhöht die Consumability, legt aber
leicht Interna offen und friert sie ein.

## Lösung

Eine *API Description* wird erstellt, die Request- und Response-Strukturen, Fehlermeldung und weitere
technische Vereinbarungen festlegt. Neben statischer und struktureller Information deckt sie auch
dynamische bzw. verhaltensbezogene Aspekte ab: Aufrufreihenfolgen, Vor- und Nachbedingungen,
Invarianten. Die syntaktische Beschreibung wird um Qualitätsmanagement-Policies, semantische
Spezifikationen und organisatorische Angaben ergänzt.

Die Quelle unterscheidet dabei zwischen *minimalen* und *elaborierten* Beschreibungen und stellt eine
Template-Vorlage vor, die geschäftliche und funktional-technische Belange auf einer Seite zusammenführt.

## Varianten

- **Minimal Description** — nur die syntaktische Schnittstelle, meist maschinenlesbar.
- **Elaborate Description** — zusätzlich Semantik, Qualitätszusagen, Eigentümer, Support, Abrechnung.
  Der SOA-Reifegrad hilft laut Quelle bei der Wahl zwischen beiden.

## Konsequenzen

**Vorteile:**

- Client-Entwicklung wird ohne Rückfragen beim Provider möglich; Codegenerierung und Mocking auch.
- Die Beschreibung wird zum Prüfstein für Änderungen: Kompatibilität lässt sich maschinell testen.

**Nachteile / Kosten:**

- Die Beschreibung muss mit der Implementierung synchron gehalten werden, sonst wird sie schädlich.
- Detaillierte Beschreibungen können Interna zementieren und die Evolvierbarkeit senken.
- Verhaltensaspekte (Reihenfolgen, Invarianten) entziehen sich den gängigen IDLs und landen in Prosa.

## Bekannte Verwendungen

Technische Notationen: Swagger bzw. die *OpenAPI Specification*, WADL, RAML für RESTful HTTP; API
Blueprint; JSON:API und APIs.json; WSDL 1.1/2.0 für SOAP; `.proto`-Dateien in Protocol Buffers; das
GraphQL-Schema; Apache Thrift IDL; Apache Avro IDL; AsyncAPI für nachrichtengetriebene APIs; MDSL als
plattformunabhängige IDL, die das MAP-Domänenmodell integriert. Auf breiterer Ebene nennt die Quelle
USDL, SSDL/SOYA sowie SWIFT, das XML/WSDL als Beschreibungssprache mit abgestuften SLAs kombiniert.
Als Templates: die *Microservices Canvas* von C. Richardson und das informelle Vertragsmodell aus Agile
Modeling.

## Verwandte Patterns

- [Service Level Agreement](../quality/ServiceLevelAgreement.md) — ergänzt die Beschreibung um
  Qualitätsziele und die Folgen ihrer Verletzung.
- [Version Identifier](../evolution/VersionIdentifier.md),
  [Semantic Versioning](../evolution/SemanticVersioning.md),
  [Two in Production](../evolution/TwoInProduction.md) — Versions- und Evolutionsinformation gehört in
  die Beschreibung.
- [Error Report](../structure/ErrorReport.md) — das Fehlermodell ist Vertragsbestandteil.
- [Public API](PublicAPI.md), [Community API](CommunityAPI.md),
  [Solution-Internal API](SolutionInternalAPI.md) — die Sichtbarkeitsstufe bestimmt, mit wem die
  Beschreibung geteilt wird und wie elaboriert sie sein muss.
- [Frontend Integration](FrontendIntegration.md), [Backend Integration](BackendIntegration.md) — beide
  Zielgruppen haben unterschiedliche Detailbedürfnisse.
- [Pricing Plan](../quality/PricingPlan.md), [Rate Limit](../quality/RateLimit.md) — nichtfunktionale
  Vertragsbestandteile.

## Bezug zu Kubernetes / KRM

Kubernetes hat eine der vollständigsten maschinenlesbaren *API Descriptions* überhaupt, und sie ist zur
Laufzeit vom Server selbst abrufbar statt aus einem Portal.

**Discovery** beantwortet, *was* es gibt: `/api` und `/apis` liefern Gruppen und Versionen,
`/apis/<group>/<version>` die `APIResourceList` mit Name, `kind`, `namespaced`, `verbs`, `shortNames`
und `categories` je Ressource.

```console
$ kubectl get --raw /apis/apps/v1 | jq
```

```json
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "apps/v1",
  "resources": [
    { "name": "deployments", "singularName": "deployment", "namespaced": true, "kind": "Deployment",
      "verbs": ["create","delete","deletecollection","get","list","patch","update","watch"],
      "shortNames": ["deploy"],   // <- die kubectl-Kurzform steckt im Vertrag, nicht im Client
      "categories": ["all"] },
    { "name": "deployments/scale", "singularName": "", "namespaced": true,
      "group": "autoscaling", "version": "v1", "kind": "Scale",   // <- Subresource mit fremder GVK
      "verbs": ["get","patch","update"] }
    // ... deployments/status, statefulsets, daemonsets, replicasets, controllerrevisions
  ]
}
```

Seit der aggregierten Discovery kommt alles in einem Roundtrip, wenn der Client
`Accept: application/json;g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList` schickt
(`staging/src/k8s.io/apiserver/pkg/endpoints/discovery/aggregated/`).

**OpenAPI** beantwortet, *wie* die Objekte aussehen: `/openapi/v3` liefert einen Index von
Gruppenversionen auf Dokument-URLs wie `/openapi/v3/apis/apps/v1?hash=...`, wobei der `hash`-Parameter
clientseitiges Caching erlaubt (`staging/src/k8s.io/client-go/openapi/client.go`). `/openapi/v2`
existiert als monolithisches Legacy-Dokument. Die Feldbeschreibungen stammen aus den Go-Doc-Kommentaren
und werden nach `types_swagger_doc_generated.go` generiert; `kubectl explain` rendert daraus (Default
`-o plaintext` aus OpenAPI v3, `-o plaintext-openapiv2` als Rückfallpfad).

```console
$ kubectl explain deployment.spec.replicas
GROUP:      apps
KIND:       Deployment
VERSION:    v1

FIELD: replicas <integer>

DESCRIPTION:
    Number of desired pods. This is a pointer to distinguish between explicit
    zero and not specified. Defaults to 1.
```

Der Text stammt wörtlich aus dem Go-Kommentar über `DeploymentSpec.Replicas` — Beschreibung und
Implementierung können also nicht auseinanderlaufen, weil sie dieselbe Quelle haben.

Bemerkenswert ist, dass Kubernetes über die `x-kubernetes-*`-Erweiterungen genau das in die Beschreibung
bekommt, was MAP als *dynamische bzw. verhaltensbezogene* Aspekte fordert und was OpenAPI von Haus aus
nicht kann: `x-kubernetes-group-version-kind` bindet Schema an GVK, `x-kubernetes-list-type`
(`atomic`/`set`/`map`) mit `x-kubernetes-list-map-keys` und `x-kubernetes-map-type` definieren
Merge-Semantik für Server-Side Apply, `x-kubernetes-patch-strategy`/`patch-merge-key` dasselbe für
Strategic Merge Patch, `x-kubernetes-int-or-string` und `x-kubernetes-preserve-unknown-fields`
beschreiben Sonderfälle der Deserialisierung, und `x-kubernetes-validations` transportiert
CEL-Ausdrücke — also echte Invarianten im Schema statt in Prosa
(`staging/src/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/types_jsonschema.go`).

Für Erweiterungen ist die Beschreibung sogar **Teil der Ressource**: eine
`CustomResourceDefinition` trägt ihr Schema unter `spec.versions[].schema.openAPIV3Schema` (Structural
Schema), plus `additionalPrinterColumns` für die Tabellendarstellung und `subresources` für
`/status` und `/scale`. Die Beschreibung ist damit nicht nur Dokumentation, sondern die ausführbare
Quelle für Validierung, Defaulting, Pruning und Merge-Verhalten — ein Grad an Bindung, den MAP nicht
verlangt.

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: databases.example.com
spec:
  # ... group, names, scope
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}                   # <- allein diese Zeile erzeugt den Endpunkt /status
        scale:
          specReplicasPath: .spec.replicas
          statusReplicasPath: .status.replicas
          labelSelectorPath: .status.selector
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              required: ["engine"]
              x-kubernetes-validations:
                - rule: "self.minReplicas <= self.maxReplicas"   # <- CEL: Invariante über zwei Felder
                  message: "minReplicas must not exceed maxReplicas"
                - rule: "self.engine == oldSelf.engine"          # <- Transitionsregel: Feld ist immutabel
                  message: "engine is immutable"
              properties:
                engine:
                  type: string
                  enum: ["postgres", "mysql"]
                minReplicas: { type: integer, minimum: 1 }
                maxReplicas: { type: integer }
                backupTargets:
                  type: array
                  x-kubernetes-list-type: map        # <- Merge-Semantik für Server-Side Apply
                  x-kubernetes-list-map-keys: ["name"]
                  items:
                    type: object
                    required: ["name"]
                    properties:
                      name: { type: string }
                      bucket: { type: string }
            status:
              type: object
              # ... replicas, selector, conditions
```

Die beiden CEL-Regeln sind genau der Fall, den MAP als *dynamischen* Aspekt beschreibt und für den es
in klassischen IDLs nur Prosa gibt: eine feldübergreifende Invariante und eine Vorbedingung an den
Übergang vom alten zum neuen Zustand. Hier stehen sie im Schema und werden vom Server erzwungen.

Was die Beschreibung **nicht** abdeckt, sind die geschäftlichen Elemente aus dem MAP-Template:
API-Eigentümer, Supportprozess, Abrechnung und
[Service Level Agreement](../quality/ServiceLevelAgreement.md) haben kein Feld. Sie liegen out-of-band
in `OWNERS`-Dateien, der Deprecation Policy und den API-Konventionen — Prosa, die den maschinenlesbaren
Teil ergänzt. Bewertung: KRM erfüllt den technisch-funktionalen Teil des Patterns **umfassender als
üblich**, den organisatorisch-geschäftlichen Teil **außerhalb der API**.

---
[← Index](../README.md) · [Kategorie Foundation](../meta/category-foundation.md) · [Quelle](https://microservice-api-patterns.org/patterns/foundation/APIDescription)
