---
title: Kategorie Foundation
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/foundation
---

# Kategorie: Foundation Patterns

## Worum es geht

Die *Foundation Patterns* beantworten die Fragen, die ganz am Anfang eines API-Vorhabens stehen,
bevor über Endpunkte, Nachrichten oder Qualität entschieden wird:

- Welche Art von (Teil-)Systemen und Komponenten wird integriert?
- Von wo aus soll die API erreichbar sein — und für wen?
- Wie wird das geteilte Wissen zwischen Provider und Client dokumentiert?

Die Kategorie zerfällt damit in drei Themen: *Integrationstyp*, *Sichtbarkeit* und *Dokumentation*.
Im Buch bilden die ersten beiden zusammen den Abschnitt „Foundations: API Visibility and
Integration Types" (Kapitel 4); die Dokumentation bekommt ein eigenes Kapitel (Kapitel 9).

Die Entscheidungen dieser Kategorie sind nicht orthogonal, aber auch nicht redundant: Integrationstyp
und Sichtbarkeit werden unabhängig voneinander gewählt (eine *Frontend Integration* kann
*Solution-Internal* oder *Public* sein), und eine *API Description* braucht jede API.

## Inhalt

### Integrationstypen — wer redet mit wem?

| Pattern | Einzeiler |
|---|---|
| [Frontend Integration](../foundation/FrontendIntegration.md) | Ein Backend stellt seine Dienste einer oder mehreren Anwendungs-Frontends über eine nachrichtenbasierte Remote-API bereit. |
| [Backend Integration](../foundation/BackendIntegration.md) | Backends verschiedener, unabhängig gebauter und deployter Anwendungen tauschen Daten aus und stoßen einander Aktivitäten an, ohne konzeptionelle Integrität aufzugeben. |

Der Unterschied ist nicht technischer, sondern kopplungsbezogener Natur: Bei *Frontend Integration*
ist der Client eine Benutzeroberfläche mit menschlichem Nutzungsprofil (Interaktivität, Latenz,
Darstellungswünsche), bei *Backend Integration* ein Programm mit maschinellem Profil (Durchsatz,
Konsistenz, Idempotenz).

### Sichtbarkeit — wer darf die API sehen und rufen?

| Pattern | Einzeiler |
|---|---|
| [Public API](../foundation/PublicAPI.md) | Die API steht im öffentlichen Internet, für eine unbegrenzte und unbekannte Zahl von Clients, begleitet von einer detaillierten *API Description*. |
| [Community API](../foundation/CommunityAPI.md) | Zugriff und Sichtbarkeit sind auf eine geschlossene Nutzergruppe über mehrere Rechtsträger hinweg beschränkt (z. B. Extranet); auch die Beschreibung wird nur dort geteilt. |
| [Solution-Internal API](../foundation/SolutionInternalAPI.md) | Die API wird nur systemintern angeboten, etwa zwischen Komponenten derselben Anwendung oder Schicht. |

Die drei Sichtbarkeiten sind eine Skala, keine Menge disjunkter Optionen: Von links nach rechts
sinken Reichweite und Kompatibilitätsdruck, gleichzeitig sinken Betriebs-, Support- und
Governance-Kosten. Die Sichtbarkeit determiniert später viele Entscheidungen aus den Kategorien
[Quality](category-quality.md) und [Evolution](category-evolution.md) — etwa ob ein
[Rate Limit](../quality/RateLimit.md) oder ein
[Pricing Plan](../quality/PricingPlan.md) überhaupt nötig ist.

### Dokumentation

| Pattern | Einzeiler |
|---|---|
| [API Description](../foundation/APIDescription.md) | Ein explizites Artefakt beschreibt Request-/Response-Strukturen, Fehlerbehandlung, Aufrufreihenfolgen, Vor-/Nachbedingungen, Invarianten sowie Qualitäts- und Organisationsinformation. |

*API Description* ist das Pattern mit den meisten eingehenden Verweisen im gesamten Sprachgebäude:
Fast alle Evolution- und Governance-Patterns setzen ein Beschreibungsartefakt voraus, in das sie
ihre Aussagen (Version, Lebensdauer, SLO, Preis) eintragen können.

## Verwandte Wiki-Seiten

- [Überblick MAP](overview.md) — Einstieg
- [Cheat Sheet](cheatsheet.md) — „Ich will eine API anbieten" → welche Foundation-Patterns
- [Patterns nach Phase](navigation-byphase.md) — die meisten Foundation-Patterns gehören in Sprint 0
- [Kategorie Responsibility](category-responsibility.md) — der nächste Schritt: Endpunkt-Rollen

## Bezug zu Kubernetes / KRM

Der `kube-apiserver` ist ein einziger Endpunkt-Verbund, der *beide* Integrationstypen gleichzeitig
bedient: `kubectl` und Dashboards sind *Frontend-Integration*-Clients, Controller, Kubelets und
Scheduler sind *Backend-Integration*-Clients. MAP würde daraus zwei APIs mit unterschiedlichen
Nachrichtengrößen und -formen machen; KRM macht bewusst das Gegenteil und bietet allen dasselbe
uniforme Schema an. Die Anpassung an das Frontend geschieht nicht durch eine eigene API, sondern
durch Content Negotiation: `kubectl get` fordert per `Accept: application/json;as=Table;g=meta.k8s.io;v=v1`
eine tabellarische Repräsentation an, die serverseitig aus demselben Objekt erzeugt wird.

Sichtbarkeit ist in KRM keine Eigenschaft der API, sondern des Deployments und der Autorisierung:
Ein Cluster-Endpunkt ist praktisch nie eine *Public API* im Pattern-Sinn (unbekannte,
unbegrenzte Clientmenge), sondern eine *Community API* oder *Solution-Internal API*, abgesichert
über RBAC (`rbac.authorization.k8s.io`), Authentifizierung und Netzwerkgrenzen. Bemerkenswert ist
die Umkehrung: Die *Spezifikation* der Kubernetes-API ist vollständig öffentlich (OpenAPI-Schemata
im Repository, api-conventions.md, Deprecation Policy), obwohl kein einzelner Endpunkt es ist. MAP
kennt diese Trennung von öffentlicher Beschreibung und nicht-öffentlichem Endpunkt nicht explizit.

Die *API Description* existiert in KRM in ungewöhnlich starker Form, weil sie zur Laufzeit vom
Server selbst ausgeliefert und von Werkzeugen konsumiert wird: Discovery unter `/api` und `/apis`
(seit v1.30 als Aggregated Discovery über `apidiscovery.k8s.io/v2`), OpenAPI v2 unter `/openapi/v2`
und OpenAPI v3 unter `/openapi/v3`, plus `kubectl explain` als menschenlesbare Sicht darauf. Bei
CRDs ist das Schema (`spec.versions[].schema.openAPIV3Schema`) sogar Teil der Ressource selbst und
wird vom Server zur Validierung verwendet — Beschreibung und Durchsetzung fallen zusammen, was MAP
so nicht vorsieht.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/foundation)
