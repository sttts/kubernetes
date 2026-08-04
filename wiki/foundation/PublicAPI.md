---
title: Public API
kategorie: Foundation
unterkategorie: API-Sichtbarkeit
quelle: https://microservice-api-patterns.org/patterns/foundation/PublicAPI
---

# Public API

*a.k.a.* Open API — wobei eine wirklich *offene* API laut Quelle eine *Public API* ohne
[API Key](../structure/APIKey.md) oder sonstiges Authentifizierungsmittel ist.

**Kurzform:** Die API wird im öffentlichen Internet exponiert und mit einer detaillierten
[API Description](APIDescription.md) versehen, damit eine unbegrenzte und potentiell unbekannte Menge
organisationsfremder Clients sie nutzen kann.

## Kontext

Es wurde entschieden, ein System mit einer Remote-API und einem oder mehreren API-Endpunkten
auszustatten. Die vorgesehenen Clients liegen in anderen Organisationen, möglicherweise in anderen
Ländern, und sind dem Provider unter Umständen gar nicht bekannt.

## Problem

Wie wird eine API einer unbegrenzten und/oder unbekannten Menge organisationsfremder Clients verfügbar
gemacht, die global, national oder regional verteilt sind?

## Forces

Die Sichtbarkeit einer API wird durch den Hosting-Ort und dessen Netzanbindung bestimmt (Internet,
Extranet, Firmennetz, einzelnes Rechenzentrum). Entscheidungskriterien:

- **Geschäftsmodell** — wird die API selbst monetarisiert oder ist sie Mittel zum Zweck?
- **Größe, Ort und Heterogenität der Zielgruppe**
- **Komplexität und Reife der benötigten Backend-Systeme und Datenspeicher** — öffentlicher Traffic ist
  unvorhersehbar.
- **Sicherheitserwägungen** — jeder Client ist potentiell feindlich.
- **Budgets** für Entwicklung, Betrieb, Wartung und Evolution.

## Lösung

Die API wird im öffentlichen Internet exponiert, zusammen mit einer detaillierten
[API Description](APIDescription.md), die funktionale *und* nichtfunktionale Eigenschaften beschreibt.

## Beispiel

Im Fallbeispiel *Lakeside Mutual* ist die Schnittstelle zwischen dem als Browser-JavaScript
implementierten Customer-Self-Service-Frontend und dem zugehörigen Backend eine *Public API*.

## Konsequenzen

**Vorteile:**

- Maximale Reichweite; Ökosystem- und Plattformeffekte werden möglich.
- Clients können unabhängig vom Provider entstehen, ohne vorherige Absprache.

**Nachteile / Kosten:**

- Die Menge der Clients ist unbekannt und nicht koordinierbar — Breaking Changes sind praktisch
  ausgeschlossen, siehe [Two in Production](../evolution/TwoInProduction.md) und
  [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md).
- Missbrauchs-, DoS- und Kostenrisiken erzwingen [Rate Limit](../quality/RateLimit.md),
  [API Key](../structure/APIKey.md) und ggf. [Pricing Plan](../quality/PricingPlan.md).
- Betrieb, Support und Dokumentation werden zu Dauerkosten.

## Bekannte Verwendungen

Die Quelle verweist auf hunderte bis tausende Einträge im Programmable-Web-Verzeichnis und nennt unter
anderem Google Calendar/Maps/Knowledge Graph, Facebook Graph API, Atlassian JIRA REST APIs, GitHub API,
YouTube Data API v3, LinkedIn REST API, AWS-Dienste wie S3 und EC2, Stripe, PayPal, Microsoft Graph
(OData-basiert), Heroku Platform API, Flickr und eBay.

## Verwandte Patterns

- [Community API](CommunityAPI.md), [Solution-Internal API](SolutionInternalAPI.md) — die Geschwister
  mit engerer Sichtbarkeit.
- [Frontend Integration](FrontendIntegration.md), [Backend Integration](BackendIntegration.md) — eine
  *Public API* unterstützt stets eine der beiden Integrationsrichtungen.
- [API Description](APIDescription.md) — zwingender Bestandteil der Lösung.
- [API Key](../structure/APIKey.md), [Rate Limit](../quality/RateLimit.md),
  [Pricing Plan](../quality/PricingPlan.md),
  [Service Level Agreement](../quality/ServiceLevelAgreement.md) — die üblichen Begleiter.
- [Version Identifier](../evolution/VersionIdentifier.md),
  [Semantic Versioning](../evolution/SemanticVersioning.md) — Evolution ohne Client-Koordination.

## Bezug zu Kubernetes / KRM

Die Kubernetes-API eines konkreten Clusters ist **keine Public API** im Sinne dieses Patterns, auch wenn
der Endpunkt bei verwalteten Angeboten (EKS, GKE, AKS) über eine öffentliche IP erreichbar sein kann.
Entscheidend ist das MAP-Kriterium „unbegrenzte und unbekannte Clientmenge“: jeder Zugriff verlangt
Authentifizierung, und die Autorisierung ist per RBAC an konkret benannte Subjekte gebunden. Der
zutreffendere Sichtbarkeitsgrad ist [Solution-Internal API](SolutionInternalAPI.md) bzw. für
Multi-Tenant-Plattformen [Community API](CommunityAPI.md).

Die einzige tatsächlich unauthentifiziert exponierte Fläche ist winzig und explizit definiert
(`plugin/pkg/auth/authorizer/rbac/bootstrappolicy/policy.go`): fünf nicht-Ressourcen-Pfade, gebunden an
die Gruppe `system:unauthenticated`.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:public-info-viewer
  # ... labels, annotations (kubernetes.io/bootstrapping: rbac-defaults)
rules:
  - nonResourceURLs:      # <- keine `resources`: das sind nackte HTTP-Pfade, keine KRM-Objekte
      - /healthz
      - /livez
      - /readyz
      - /version
      - /version/
    verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:public-info-viewer
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:public-info-viewer
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:unauthenticated    # <- die gesamte „öffentliche“ API-Oberfläche des Clusters
```

Anonyme Requests laufen unter dem Benutzer `system:anonymous`; jenseits dieser fünf Pfade endet die
Öffentlichkeit sofort — inklusive Discovery und OpenAPI, siehe [Community API](CommunityAPI.md).

```http
GET /livez HTTP/1.1
# kein Authorization-Header

HTTP/1.1 200 OK
ok
```

```http
GET /api/v1/namespaces/default/pods HTTP/1.1
# ebenfalls kein Authorization-Header

HTTP/1.1 403 Forbidden
Content-Type: application/json
```

```json
{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure",
 "message":"pods is forbidden: User \"system:anonymous\" cannot list resource \"pods\" in API group \"\" in the namespace \"default\"",
 "reason":"Forbidden","details":{"kind":"pods"},"code":403}
```

Public im MAP-Sinn ist dagegen die **Spezifikation**: die API-Konventionen, die Go-Typen unter
`staging/src/k8s.io/api/`, das generierte OpenAPI-Dokument und die Referenzdokumentation sind frei
zugänglich und werden über den KEP-Prozess öffentlich weiterentwickelt. Kubernetes trennt also, was MAP
zusammendenkt: die *Beschreibung* ist öffentlich, der *Endpunkt* ist es nicht.

Die typischen Public-API-Begleiter fehlen entsprechend oder sehen anders aus. Es gibt keinen
[API Key](../structure/APIKey.md) als Sichtbarkeitsmittel (stattdessen ServiceAccount-Tokens, Client-
Zertifikate, OIDC); Rate Limiting ist nicht kommerziell, sondern als API Priority and Fairness
(`flowcontrol.apiserver.k8s.io`) ein Überlastschutz; ein [Pricing Plan](../quality/PricingPlan.md)
existiert im Kern nicht. Die Evolutionsgarantien dagegen sind strenger als bei den meisten öffentlichen
APIs: eine GA-Gruppenversion wie `apps/v1` darf nicht inkompatibel geändert werden, was faktisch einer
[Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) nahekommt — die Deprecation
Policy verlangt für GA-APIs mindestens zwölf Monate bzw. drei Releases Vorlauf.

Bewertung: **teils gar nicht** (der Cluster-Endpunkt ist nie öffentlich im Pattern-Sinn), **teils
anders** (die API-Spezifikation ist öffentlich, ohne dass ein Endpunkt es wäre).

---
[← Index](../README.md) · [Kategorie Foundation](../meta/category-foundation.md) · [Quelle](https://microservice-api-patterns.org/patterns/foundation/PublicAPI)
