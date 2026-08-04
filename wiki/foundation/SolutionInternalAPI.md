---
title: Solution-Internal API
kategorie: Foundation
unterkategorie: API-Sichtbarkeit
quelle: https://microservice-api-patterns.org/patterns/foundation/SolutionInternalAPI
---

# Solution-Internal API

*a.k.a.* Application-Internal Programming Interface

**Kurzform:** Die Anwendung wird logisch in Komponenten zerlegt, die lokale oder remote APIs exponieren;
diese werden ausschließlich systeminternen Kommunikationspartnern angeboten.

## Kontext

Es wurde entschieden, ein Subsystem einer Anwendung mit einer Remote-API auszustatten, aber die
Sichtbarkeitsstufe ist noch offen. Die vorgesehenen Clients gehören zur selben Anwendung bzw. Lösung.

## Problem

Wie werden Zugriff auf und Nutzung einer API auf eine Anwendung begrenzt — etwa auf Komponenten
derselben oder einer anderen logischen Schicht bzw. eines anderen physischen Tiers?

## Forces

Die Kräfte von [Public API](PublicAPI.md) und [Community API](CommunityAPI.md) gelten auch hier und
werden von der Quelle nicht wiederholt:

- Geschäftsmodell und Stakeholder-Interessen
- Größe, Ort und Heterogenität der Zielgruppe
- Komplexität und Reife der benötigten Backend-Systeme und Datenspeicher
- Sicherheitserwägungen
- Technische Präferenzen
- Lifecycle-Management
- Budgets

## Lösung

Die Anwendung wird logisch in Komponenten zerlegt. Diese exponieren lokale oder remote APIs, die nur
systeminternen Partnern angeboten werden — etwa anderen Services im Anwendungs-Backend.

## Beispiel

In *Lakeside Mutual* zählen die Schnittstellen zwischen Customer Core und seinen Hilfsdiensten
(Ortssuche per Postleitzahl, Telefonnummernvalidierung) dazu, ebenso der Risk-Calculation-Service. Die
Quelle zeigt dessen OpenAPI-Auszug:

```json
"/riskfactor/compute": {
  "post": {
    "tags": ["risk-computation-service"],
    "summary": "Computes the risk factor for a given customer.",
    "operationId": "computeRiskFactorUsingPOST",
    "consumes": ["application/json"],
    "parameters": [{
      "in": "body",
      "name": "riskFactorRequest",
      "required": true,
      "schema": { "$ref": "#/definitions/RiskFactorRequestDto" }
    }],
    "responses": {
      "200": { "schema": { "$ref": "#/definitions/RiskFactorResponseDto" } },
      "401": { "description": "Unauthorized" },
      "403": { "description": "Forbidden" }
    }
  }
}
```

## Konsequenzen

**Vorteile:**

- Die Clientmenge ist vollständig bekannt und wird vom selben Team bzw. derselben Organisation
  kontrolliert — Breaking Changes sind koordinierbar.
- Weniger Aufwand für Dokumentation, Rate Limiting, Abrechnung und Support als bei weiterer Sichtbarkeit.

**Nachteile / Kosten:**

- Die Grenze ist organisatorisch, nicht technisch erzwungen; interne APIs erodieren leicht zu De-facto-
  Community-APIs, sobald Nachbarteams sie entdecken.
- Fehlende Beschreibungsdisziplin rächt sich später bei jeder Beförderung der Sichtbarkeitsstufe.

## Bekannte Verwendungen

Viele interne APIs in Unternehmen und Behörden fallen hierunter, sind aber selten öffentlich
dokumentiert. Die Quelle nennt die produktinternen Management-APIs von Apigee Edge (Analytics,
Billing), das Eclipse Communication Framework mit Distribution Providern für ActiveMQ, gRPC, MQTT und
JAX-RS sowie Terravis, wo eine Prozess-Engine systeminterne Operationen zur Dokumentenerzeugung und für
nicht-öffentliche Stammdaten aufruft.

## Verwandte Patterns

- [Public API](PublicAPI.md), [Community API](CommunityAPI.md) — die Geschwister mit weiterer Sichtbarkeit.
- [Frontend Integration](FrontendIntegration.md), [Backend Integration](BackendIntegration.md) — beide
  Richtungen sind möglich.
- [API Description](APIDescription.md) — auch intern nötig, oft nur in minimaler Ausprägung.
- [Computation Function](../responsibility/ComputationFunction.md) — der Risk-Calculation-Service im
  Beispiel ist genau das.
- [Solution-Internal API](SolutionInternalAPI.md) wird durch das API-Refactoring *Extract Endpoint*
  weiter zerlegt (kein MAP-Pattern, sondern ein Refactoring aus derselben Quelle).

## Bezug zu Kubernetes / KRM

Kubernetes kennt eine harte, im Code verankerte Ausprägung dieses Patterns: die **internen
Versionen**. Jede API-Gruppe hat neben den serialisierten Außenversionen (`v1`, `v1beta1`, …) eine
Hub-Version mit der Kennung `runtime.APIVersionInternal`, konstant `"__internal"`
(`staging/src/k8s.io/apimachinery/pkg/runtime/interfaces.go`). Die zugehörigen Typen unter
`pkg/apis/<group>/types.go` werden nie über die Leitung geschickt und nie persistiert; sie existieren
nur, damit Konversion zwischen Außenversionen sternförmig statt paarweise erfolgt. Das ist eine echte
*Solution-Internal API* im MAP-Sinn — nur eben in-process statt remote.

```go
// staging/src/k8s.io/apimachinery/pkg/runtime/interfaces.go
// APIVersionInternal may be used if you are registering a type that should not
// be considered stable or serialized - it is a convention only ...
const APIVersionInternal = "__internal"

// pkg/apis/apps/register.go — dieselbe Gruppe, aber die nicht serialisierte Hub-Version
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}
```

```yaml
# So etwas existiert nicht und kann nicht existieren — "__internal" ist nie ein apiVersion-Wert
# auf der Leitung oder in etcd:
apiVersion: apps/__internal
kind: Deployment

---
# Sichtbar sind nur die Außenversionen. Die interne Version ist der Knoten dazwischen:
# v1 -> __internal -> v2, nie v1 -> v2 direkt.
apiVersion: apps/v1     # <- das ist alles, was ein Client je zu sehen bekommt
kind: Deployment
metadata:
  name: web
  # ... namespace, resourceVersion
```

Auf der Remote-Ebene gibt es mehrere Gruppen, deren Zielgruppe ausschließlich Systemkomponenten sind:
`authentication.k8s.io` (`TokenReview`), `authorization.k8s.io` (`SubjectAccessReview`),
`certificates.k8s.io`, `coordination.k8s.io` (`Lease` für Leader Election) und
`flowcontrol.apiserver.k8s.io`. Sie sind für Benutzer sichtbar (Discovery ist nicht RBAC-gefiltert,
siehe [Community API](CommunityAPI.md)), aber praktisch nur für die Steuerebene gedacht. Feiner
granular wirkt die Trennung über **Subresourcen**: `/status` und `/scale` sind eigene
RBAC-Ressourcen, sodass Controller `pods/status` schreiben dürfen, während Benutzer nur `pods` ändern
— die Spec/Status-Trennung ist damit zugleich eine Sichtbarkeitsgrenze zwischen Benutzer-API und
controller-interner API.

```yaml
# ClusterRole system:kube-scheduler (Auszug)
rules:
  - apiGroups: [""]
    resources: ["pods/status"]      # <- eigene RBAC-Ressource, von "pods" nicht mit abgedeckt
    verbs: ["patch", "update"]

---
# Aggregierte ClusterRole system:aggregate-to-edit (Auszug) — volle Schreibrechte auf die Spec …
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["create", "delete", "deletecollection", "patch", "update"]
    # <- "pods/status" fehlt hier bewusst: kein Benutzer der Rolle "edit" schreibt Status

---
# … und in system:aggregate-to-view taucht der Status auf, aber nur lesend:
rules:
  - apiGroups: [""]
    resources: ["pods/status"]
    verbs: ["get", "list", "watch"]
```

Nicht als KRM modelliert und deshalb umso deutlicher solution-internal sind die Komponenten-Configs
(`apiserver.config.k8s.io`, `kubelet.config.k8s.io`, `kubescheduler.config.k8s.io`): versionierte
API-Objekte, die aber nur aus Dateien gelesen und nie über den apiserver serviert werden. Ebenso die
gRPC-Schnittstellen des kubelets (CRI, CSI, Device Plugin API), die node-lokal über Unix-Sockets laufen.

```yaml
# /var/lib/kubelet/config.yaml — gelesen per --config, nie über GET /apis/... erreichbar
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration        # <- kein metadata: kein Name, kein Namespace, keine resourceVersion
clusterDomain: cluster.local
cgroupDriver: systemd
authentication:
  anonymous:
    enabled: false
  webhook:
    enabled: true
authorization:
  mode: Webhook
# ... syncFrequency, staticPodPath, clusterDNS
```

```console
$ kubectl get --raw /apis | jq -r '.groups[].name' | grep config.k8s.io
# (leer) — die Gruppe ist versioniert und konvertierbar, aber nicht Teil der Remote-API
```

Bewertung: **erfüllt, mit eigener Ausprägung**. KRM zieht die Grenze nicht am Deployment-Ort, wie MAP es
vorschlägt, sondern über drei andere Mechanismen: nicht-serialisierte interne Versionen, separate
RBAC-Ressourcen für Subresourcen und Konvention (`system:`-Präfixe, `*.k8s.io`-Gruppen). Grund ist das
uniforme Modell: es gibt genau einen Endpunkt für alles, also muss Sichtbarkeit im Schema und in der
Autorisierung ausgedrückt werden statt in der Netzwerktopologie.

---
[← Index](../README.md) · [Kategorie Foundation](../meta/category-foundation.md) · [Quelle](https://microservice-api-patterns.org/patterns/foundation/SolutionInternalAPI)
