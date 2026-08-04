---
title: Frontend Integration
kategorie: Foundation
unterkategorie: Integrationsrichtung
quelle: https://microservice-api-patterns.org/patterns/foundation/FrontendIntegration
---

# Frontend Integration

*a.k.a.* Vertical Integration, Frontend-to-Backend Integration, Frontend API, North-South Connectivity

**Kurzform:** Ein Backend stellt seine Dienste über eine nachrichtenbasierte Remote-API bereit, damit
physisch getrennte Endbenutzer-Oberflächen Daten laden, darstellen und Aktionen im Backend auslösen können.

## Kontext

Es liegt eine Architekturentscheidung vor, das System in physische Tiers zu zerlegen. Das Frontend, das
die Benutzeroberfläche darstellt und steuert, läuft getrennt vom Backend, das Verarbeitung und
Datenhaltung verantwortet. Wegen dieser physischen Trennung kann das Frontend das Backend nicht über
lokale, sprach- oder betriebssysteminterne Aufrufe erreichen; Remoting ist zwingend.

## Problem

Wie werden clientseitige Oberflächen mit Rechenergebnissen, Suchergebnismengen und Detaildaten zu
Entitäten befüllt und aktualisiert? Wie stoßen sie Aktivitäten im Backend an oder laden Daten hinauf?
Und wie lassen sich mehrere Oberflächen für mehrere Plattformen unabhängig entwickeln, ohne sie an die
serverseitige Implementierung zu koppeln?

## Forces

- **Business needs** — welche fachlichen Abläufe die Oberfläche überhaupt unterstützen muss.
- **Client-Developer-Komfort vs. Risiko und Kosten** — bequeme, breit geschnittene Frontend-APIs
  erhöhen die Kopplung zwischen Oberfläche und Backend-Interna.
- **Sicherheit und Datenschutz** — das Frontend läuft auf fremdkontrollierten Geräten; alles, was es
  sehen darf, ist potentiell öffentlich.
- Konflikte zwischen diesen Kräften treten in der Frontend-Richtung anders auf als bei
  Backend-zu-Backend-Integration, etwa weil Clients nicht mitversioniert deployt werden können.

## Lösung

Das Backend der verteilten Anwendung exponiert seine Dienste gegenüber einem oder mehreren Application
Frontends über eine nachrichtenbasierte Remote-API. Die Quelle nennt zwei Integrationsrichtungen im
selben Bild: Frontend-to-Backend (dieses Pattern) und Backend-to-Backend
([*Backend Integration*](BackendIntegration.md)).

## Beispiel

Im Fallbeispiel *Lakeside Mutual* ist die Schnittstelle zwischen dem Customer-Self-Service-Frontend und
seinem Backend eine *Frontend Integration*. Ebenso die Schnittstelle zwischen der Customer-Management-
Rich-Internet-Application und der Interface-Schicht des Policy-Management-Backends: JavaScript im
Browser ist der API-Client, holt Vertrags- und Angebotsdaten, stellt sie dar und schickt
Änderungswünsche zurück.

## Konsequenzen

**Vorteile:**

- Oberflächen für verschiedene Kanäle (Web, Mobile, Desktop) können unabhängig vom Backend entwickelt
  und deployt werden.
- Die Backend-Logik bleibt an einer Stelle und ist für alle Kanäle wiederverwendbar.

**Nachteile / Kosten:**

- Latenz, Teilausfälle und Serialisierung werden zum Designthema, wo vorher lokale Aufrufe standen.
- Die API wird zur Angriffsfläche: Autorisierung und Datenminimierung müssen serverseitig erzwungen
  werden, weil der Client nicht vertrauenswürdig ist.
- Chattiness und Over-/Under-Fetching entstehen, wenn die API nicht auf die Darstellungsbedürfnisse
  zugeschnitten ist — daher die Ergänzungen *Backends for Frontends* und *API Gateway*.

## Bekannte Verwendungen

Praktisch jedes System mit Benutzeroberfläche. Die Quelle nennt das Terravis-Portal, die Trennung von
Frontend- und Backend-APIs bei einer großen Schweizer Bank (Murer/Hagen 2014, „Managed Evolution“) sowie
das *Dynamic Interface* aus Brandner et al. (2004), das über 1000 Backend-Funktionen an Online-Banking-
und Filialanwendungen exponiert.

## Verwandte Patterns

- [Backend Integration](BackendIntegration.md) — das Geschwister-Pattern, die andere Integrationsrichtung.
- [Public API](PublicAPI.md), [Community API](CommunityAPI.md),
  [Solution-Internal API](SolutionInternalAPI.md) — orthogonale Sichtbarkeitsstufen; eine
  *Frontend Integration* API kann jede davon haben.
- [API Description](APIDescription.md) — was Frontend-Entwickler brauchen, um gegen die API zu codieren.
- [Pagination](../quality/Pagination.md), [Wish List](../quality/WishList.md),
  [Wish Template](../quality/WishTemplate.md) — Antworten auf Over-/Under-Fetching in Oberflächen.
- *Backends for Frontends* (S. Newman) und *API Gateway* — kanalspezifische Edge-Services bzw.
  Entkopplung von Client und Provider; keine MAP-Patterns, aber in der Quelle referenziert.

## Bezug zu Kubernetes / KRM

Kubernetes betreibt **keine getrennte Frontend-API**. `kubectl`, das Dashboard, `k9s` und jedes
client-go-basierte Werkzeug sprechen exakt dieselbe HTTPS/JSON-REST-API des kube-apiservers wie
Controller und kubelet. Das ist eine bewusste Abweichung von der Trennung, die MAP zwischen
*Frontend Integration* und [Backend Integration](BackendIntegration.md) zieht: das uniforme
Ressourcenschema über alle Gruppen hinweg macht ein separates Frontend-API-Layer weitgehend überflüssig,
und dieselbe Autorisierung (RBAC) gilt für Mensch und Maschine.

```http
# So würde MAP das aufteilen — je Integrationsrichtung ein eigener Endpunkt:
GET /frontend/v1/deployments?view=list     # zugeschnitten auf die Oberfläche
GET /internal/v1/deployments               # für andere Backends
```

```http
# KRM: ein Pfad für beide. Der einzige Unterschied ist die Identität im Token.
GET /apis/apps/v1/namespaces/default/deployments
Authorization: Bearer <token>   # <- kubectl-Benutzer oder deployment-controller — derselbe Handler
Accept: application/json
```

Eine echte frontend-spezifische Repräsentation existiert trotzdem: das **Table-Format**. Ein Client
schickt `Accept: application/json;as=Table;v=v1;g=meta.k8s.io` und erhält statt der Ressource eine
`metav1.Table` mit Spalten und Zeilen — genau das, was `kubectl get` in Tabellenform ausgibt. Die
Spaltendefinition liefert der Server (`TableConvertor` im apiserver-Endpoint-Handler, bei CRDs
`spec.versions[].additionalPrinterColumns`). Das ist ein serverseitig gerendertes, rein
darstellungsorientiertes View-Modell — funktional das, was MAP als kanalspezifischen Zuschnitt
beschreibt, nur als Content-Negotiation statt als eigener Endpunkt.

```http
GET /apis/apps/v1/namespaces/default/deployments HTTP/1.1
Accept: application/json;as=Table;v=v1;g=meta.k8s.io, application/json   # <- Fallback, falls der Server kein Table kann
```

```json
{
  "kind": "Table",
  "apiVersion": "meta.k8s.io/v1",
  "metadata": { "resourceVersion": "482931" },
  "columnDefinitions": [
    { "name": "Name",  "type": "string", "format": "name", "description": "Name must be unique within a namespace. ...", "priority": 0 },
    { "name": "Ready", "type": "string", "format": "", "description": "Number of the pod with ready state", "priority": 0 },
    // ... Up-to-date, Available, Age (priority 0)
    { "name": "Containers", "type": "string", "format": "", "description": "Names of each container in the template.", "priority": 1 }
    // <- priority != 0: von kubectl nur bei -o wide gezeigt, der Server liefert die Spalte immer
    // ... Images, Selector — ebenfalls priority 1
  ],
  "rows": [
    {
      "cells": ["web", "3/3", "3", "3", "12d", "app", "nginx:1.27", "app=web"],
      // <- so viele Zellen wie Spaltendefinitionen, in derselben Reihenfolge
      "object": {
        "kind": "PartialObjectMetadata", "apiVersion": "meta.k8s.io/v1",
        "metadata": { "name": "web", "namespace": "default" }
        // <- includeObject=Metadata ist der Default; None und Object sind die Alternativen
      }
    }
  ]
}
```

Für CustomResources ist derselbe Mechanismus deklarativ zugänglich: die CRD trägt die Spalten der
Tabellenansicht als Teil ihres eigenen Schemas.

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: certificates.cert-manager.io
spec:
  # ... group, names, scope
  versions:
    - name: v1
      served: true
      storage: true
      # ... schema.openAPIV3Schema, subresources
      additionalPrinterColumns:
        - name: Ready                # <- wird zum Spaltenkopf in `kubectl get`
          type: string
          jsonPath: .status.conditions[?(@.type=="Ready")].status
        - name: Age
          type: date
          jsonPath: .metadata.creationTimestamp
        - name: Issuer
          type: string
          priority: 1                # <- erst bei `kubectl get -o wide`
          jsonPath: .spec.issuerRef.name
```

Weitere frontend-orientierte Mechanismen: `kubectl explain` rendert aus dem OpenAPI-v3-Dokument (siehe
[API Description](APIDescription.md)), Label- und Feld-Selektoren sowie `limit`/`continue`
([Pagination](../quality/Pagination.md)) begrenzen die Übertragungsmenge, und `kubectl proxy` stellt
einen lokalen, bereits authentifizierten Zugang für Browser-Frontends bereit.

Ein *Backends for Frontends* im Newman'schen Sinne findet sich in Kubernetes nur außerhalb des Kerns:
das Kubernetes Dashboard betreibt einen eigenen Backend-Prozess, der gegenüber dem Browser als BFF
auftritt und selbst client-go gegen den apiserver spricht. Der apiserver selbst kennt keine
kanalspezifischen Endpunkte. Bewertung: KRM erfüllt das Pattern **anders** — die Integrationsrichtung
ist kein API-Designkriterium, sondern nur eine Frage von Identität und RBAC-Rechten des Aufrufers.

---
[← Index](../README.md) · [Kategorie Foundation](../meta/category-foundation.md) · [Quelle](https://microservice-api-patterns.org/patterns/foundation/FrontendIntegration)
