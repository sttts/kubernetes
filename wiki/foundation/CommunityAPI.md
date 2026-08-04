---
title: Community API
kategorie: Foundation
unterkategorie: API-Sichtbarkeit
quelle: https://microservice-api-patterns.org/patterns/foundation/CommunityAPI
---

# Community API

**Kurzform:** Die API wird zugangsbeschränkt deployt — etwa in einem Extranet — und ihre
[API Description](APIDescription.md) nur an eine geschlossene Benutzergruppe verteilt, die sich über
mehrere Rechtssubjekte erstreckt.

## Kontext

Es wurde entschieden, ein (Sub-)System mit einer Remote-API auszustatten, aber die Sichtbarkeitsstufe
ist noch offen. Die vorgesehenen Clients liegen in verschiedenen Teilen der Organisation, die das System
besitzt und betreibt — bzw. bei mehreren beteiligten Organisationen.

## Problem

Wie lassen sich Sichtbarkeit und Zugriff auf eine API auf eine geschlossene Benutzergruppe beschränken,
die nicht zu einer einzigen Organisationseinheit gehört, sondern zu mehreren Rechtssubjekten —
Unternehmen, Non-Profits, Behörden? Die Sichtbarkeit wird dabei durch den Deployment-Ort und dessen
Netzanbindung bestimmt.

## Forces

Die Kräfte von [Public API](PublicAPI.md) gelten auch hier. Zusätzlich erschweren *Community APIs*
Konzeption, Implementierung und Deployment durch:

- **Größe, Ort und technische Präferenzen der Community** — heterogene Technologiestacks bei rechtlich
  getrennten Teilnehmern.
- **Stakeholder-Interessen** — niemand hat allein die Entscheidungshoheit über das Schema.
- **Lifecycle-Management** — Aufnahme, Ausschluss und Migration von Teilnehmern.
- **Budgets** für Entwicklung, Betrieb, Wartung und Evolution.
- **Community-spezifische Sicherheitserwägungen** — Föderation von Identitäten über Organisationsgrenzen.

## Lösung

API und Implementierungsressourcen werden sicher an einem zugangsbeschränkten Ort deployt, sodass nur
die gewünschte Benutzergruppe sie erreicht — beispielsweise in einem Extranet. Die
[API Description](APIDescription.md) wird ausschließlich mit dieser eingeschränkten Zielgruppe geteilt.

## Beispiel

In *Lakeside Mutual* qualifizieren sich die Schnittstellen zwischen den Bounded Contexts im Backend als
Beispiele, etwa die vom Customer-Core-Service exponierte und von Policy Management, Customer Management
und Customer Self Service konsumierte API.

## Konsequenzen

**Vorteile:**

- Die Clientmenge ist bekannt und benennbar; koordinierte Migrationen werden möglich.
- Weniger Angriffsfläche und geringerer Missbrauchsdruck als bei einer [Public API](PublicAPI.md).

**Nachteile / Kosten:**

- Governance über Organisationsgrenzen hinweg ist aufwendig: wer entscheidet über Schemaänderungen?
- Zugangsverwaltung (Onboarding, Zertifikate, Föderation) wird zum eigenen Betriebsthema.
- Die Netzwerkgrenze allein ist keine Sicherheitsgrenze; Authentifizierung bleibt nötig.

## Bekannte Verwendungen

Die Quelle nennt SWITCH edu-ID als Identitätslösung für Schweizer Hochschulen, die SOA-/API-Evolution
bei der Credit Suisse (Murer/Hagen 2014), das *Dynamic Interface* aus Brandner et al. (2004) sowie
Terravis (Lübke/van Lessen 2016), wo die Community rechtlich beschränkt ist — nur bestimmte
Organisationen dürfen auf Grundbuchdaten zugreifen. Ergänzend ein Wholesale-Retail-Order-Management-SOA
(Zimmermann et al. 2005). Viele solcher APIs sind nicht öffentlich dokumentiert.

## Verwandte Patterns

- [Public API](PublicAPI.md), [Solution-Internal API](SolutionInternalAPI.md) — die Geschwister mit
  weiterer bzw. engerer Sichtbarkeit.
- [Frontend Integration](FrontendIntegration.md), [Backend Integration](BackendIntegration.md) — eine
  *Community API* unterstützt eine der beiden Richtungen.
- [API Description](APIDescription.md) — hier bewusst nur eingeschränkt verteilt.
- [Service Level Agreement](../quality/ServiceLevelAgreement.md) — bei mehreren Rechtssubjekten
  praktisch unvermeidlich.
- [Two in Production](../evolution/TwoInProduction.md) — koordinierte, aber nicht synchrone Migration
  bekannter Teilnehmer.

## Bezug zu Kubernetes / KRM

Die API eines konkreten Clusters passt am ehesten hierher, sobald mehrere Teams oder Mandanten sie
nutzen: der Endpunkt ist netzseitig beschränkt (privater Endpunkt, VPN, autorisierte CIDRs), die
Teilnehmer sind bekannt, und die Zugehörigkeit wird über Identitäten und RBAC-Bindings verwaltet.
Namespaces plus `Role`/`RoleBinding` bilden dabei die organisatorischen Grenzen ab, `ClusterRole`/
`ClusterRoleBinding` die clusterweiten. Bei mehreren Rechtssubjekten kommt in der Praxis
OIDC-Föderation oder ein externer Webhook-Authentifizierer hinzu.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: deployer
  namespace: team-a          # <- die Mandantengrenze steckt im Objekt, nicht im Netz
rules:
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: deployer
  namespace: team-a
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: deployer
subjects:
  - kind: Group
    apiGroup: rbac.authorization.k8s.io
    name: oidc:team-a        # <- föderierte Identität aus dem OIDC-Provider
```

Ein Punkt weicht hier deutlich vom Pattern ab. MAP verlangt, die
[API Description](APIDescription.md) *nur* mit der Zielgruppe zu teilen — Kubernetes teilt sie mit
**jedem authentifizierten Benutzer**, unabhängig von dessen Rechten.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:discovery
rules:
  - nonResourceURLs:
      - /api
      - /api/*
      - /apis
      - /apis/*          # <- Discovery aller Gruppen, auch der fremden
      - /openapi
      - /openapi/*       # <- das vollständige Schema, auch fremder CRDs
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
  name: system:discovery
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:discovery
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated    # <- jeder, der sich anmelden kann — kein Rechtebezug
```

Wer sich anmelden kann, sieht also das vollständige Schema aller Gruppen und Ressourcen inklusive aller
CRDs — auch derjenigen, auf die er keinerlei Zugriff hat.

```console
$ kubectl auth can-i list certificates.cert-manager.io --namespace team-b
no

$ kubectl get --raw /apis/cert-manager.io/v1 | jq -r '.resources[].name'
certificates
certificates/status
certificaterequests
# ... Beschreibung sichtbar, Daten nicht — genau die Trennung, die MAP nicht vorsieht
```

Das ist eine bewusste Entscheidung: Discovery und Versionsaushandlung müssen funktionieren, bevor
Autorisierung überhaupt sinnvoll auswertbar ist, und das uniforme Schema soll für alle Clients
identisch sein.

Auf der Ebene der *Spezifikation* trifft „Community“ dagegen gut: die KRM-API wird von einer
Multi-Organisations-Community (SIGs, KEP-Prozess, API-Reviewer) verwaltet, mit
`api-conventions.md` als geteiltem Vertrag. Die Gruppendomäne kodiert das sogar: `*.k8s.io` ist der
Community, herstellereigene CRDs sind an eigene Domänen gebunden — eine strukturelle Zuordnung von
API-Oberfläche zu Eigentümer, für die MAP kein Gegenstück hat.

Bewertung: **im Kern erfüllt, in einem Detail bewusst anders** — beschränkter Zugriff ja, beschränkte
Beschreibung nein.

---
[← Index](../README.md) · [Kategorie Foundation](../meta/category-foundation.md) · [Quelle](https://microservice-api-patterns.org/patterns/foundation/CommunityAPI)
