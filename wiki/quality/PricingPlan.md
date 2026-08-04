---
title: Pricing Plan
kategorie: Quality
unterkategorie: Quality Management and Governance
quelle: https://microservice-api-patterns.org/patterns/quality/qualityManagementAndGovernance/PricingPlan
---

# Pricing Plan

*a.k.a.* Metering and Billing, Accounting

**Kurzform:** Der Provider hängt an die [API Description](../foundation/APIDescription.md) ein
Preismodell, definiert messbare Nutzungsmetriken und rechnet Clients (oder Werbetreibende) danach ab.

## Kontext

Eine API ist ein Asset mit monetärem und immateriellem Wert; Entwicklung und Betrieb müssen finanziert
werden. Die Refinanzierung kann über Nutzungsgebühren der Clients laufen, aber auch über Werbung oder
andere Quellen.

## Problem

Wie misst der Provider die Inanspruchnahme des API-Dienstes und rechnet sie ab?

## Forces

- **Ökonomie:** Das Modell muss die Kosten decken und zugleich für Kunden attraktiv und
  vergleichbar bleiben; zu komplexe Tarife schrecken ab.
- **Genauigkeit:** Zählwerte müssen belastbar sein — verlorene oder doppelt gezählte Aufrufe sind
  direkt Geld.
- **Granularität des Meters:** Pro Aufruf? Pro Operation? Pro übertragenem Byte, pro Rechenzeit, pro
  Datensatz? Feinere Metriken sind fairer, aber teurer zu erheben und schwerer zu erklären.
- **Sicherheit:** Der Meter ist ein Angriffsziel — sowohl für Kunden (Umgehen der Zählung) als auch
  gegen Kunden (Zurechnung fremder Last).

## Lösung

Der API Description einen *Pricing Plan* zuordnen, der die Abrechnung von Kunden, Werbetreibenden oder
sonstigen Stakeholdern trägt. Dazu Nutzungsmetriken definieren und überwachen, z. B. Aufrufstatistiken
pro Operation.

## Varianten

- **Usage-based / verbrauchsabhängig:** Preis skaliert mit gemessener Nutzung, oft in Stufen
  (Freemium-Stufe, Staffelpreise, Volumenrabatte).
- **Flat-Rate-Subscription:** fester Betrag pro Periode, minimaler Messaufwand, kein Kostenrisiko für
  den Kunden — dafür Kostenrisiko beim Provider.
- **Market-based:** Preis ergibt sich aus Angebot und Nachfrage (z. B. Spot-Märkte).

## Beispiel

Ein fiktiver Anbieter einer E-Mail-Versand-API mit verbrauchsabhängigem Plan:

| E-Mails pro Monat (bis) | Preis pro Monat |
|---|---|
| 100 | kostenlos |
| 10 000 | 20 $ |
| 100 000 | 150 $ |
| 1 000 000 | 1000 $ |

Ein Wettbewerber, der Monitoring minimieren will, bietet stattdessen 50 $/Monat flat für unbegrenzten
Versand.

## Konsequenzen

**Vorteile:**

- Macht die API zu einem finanzierbaren Produkt statt zu einem Kostenblock.
- Die Freemium-Stufe senkt die Einstiegshürde und liefert zugleich Nutzungsdaten.
- Nutzungsmetriken sind auch ohne Abrechnung wertvoll: Kapazitätsplanung, Deprecation-Entscheidungen.

**Nachteile / Kosten:**

- Metering-, Rating- und Billing-Infrastruktur ist ein eigenes Subsystem mit eigener Verfügbarkeits-
  und Korrektheitsanforderung.
- Preisänderungen sind fast so heikel wie inkompatible API-Änderungen — sie brechen die
  Geschäftsgrundlage der Clients.
- Feingranulare Tarife verleiten Clients zu Optimierungen, die die API-Nutzung verzerren.

## Bekannte Verwendungen

- *Business Support Systems (BSS)* der Telekommunikation; das TM Forum Applications Framework 3.0
  beschreibt diese Fähigkeiten plattformneutral.
- Ein dynamisches Interface zu einem Kernbankensystem (Brandner et al. 2004) mit Preisen pro
  Operationstyp — ein Portfolio-Lookup kostet anderes als eine Saldoabfrage.
- AWS Lambda: verbrauchsabhängig mit Free Tier von 1 Mio. Requests/Monat.
- Amazon S3: Preis nach Speicherkapazität und Operationen, mit Stufen und Volumenrabatten.
- CloudConvert: Prepaid-Pakete an Konvertierungsminuten neben monatlichen Abos.
- Amazon EC2 Spot Instances als Beispiel für marktbasierte Preisbildung.

## Verwandte Patterns

- [Rate Limit](RateLimit.md) — setzt die Abrechnungsstufen technisch durch.
- [Service Level Agreement](ServiceLevelAgreement.md) — der Pricing Plan sollte darauf verweisen; Preis und zugesagte Qualität gehören zusammen.
- [API Key](../structure/APIKey.md) — identifiziert den abzurechnenden Client.
- [Wish List](WishList.md), [Wish Template](WishTemplate.md) — halten das übertragene Datenvolumen klein, wenn Volumen Teil des Preismodells ist.
- [API Description](../foundation/APIDescription.md) — Trägerdokument des Plans.
- [Public API](../foundation/PublicAPI.md) — der Sichtbarkeitsbereich, in dem Pricing Plans überhaupt relevant werden.

## Bezug zu Kubernetes / KRM

Hier ist die ehrliche Antwort **(c) gar nicht**. Das Kubernetes Resource Model kennt kein
Abrechnungskonzept: Es gibt keine API-Gruppe für Metering oder Billing, keine Ressource, die Preise
oder Tarifstufen beschreibt, und der Apiserver zählt Requests nicht pro Konsument zur Verrechnung.
Die eingebaute Zurechnung endet bei Authentifizierung (`user.Info` mit Name, Groups, Extra) und
Audit-Logging; das Audit-Log wäre die nächstliegende Datenquelle für ein externes Metering, ist aber
nicht dafür gebaut.

Die naheliegende Analogie führt in die Irre. `ResourceQuota` sieht auf den ersten Blick aus wie ein
Tarif mit Zählwerk — Soll gegen Ist, sauber pro Namespace geführt:

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: team-a
  namespace: team-a
spec:
  hard:
    requests.cpu: "20"
    limits.memory: 64Gi
    pods: "150"
    count/deployments.apps: "50"    # <- count/<resource>.<group> zählt Objekte, nicht Aufrufe
    services.loadbalancers: "2"
  # ... scopes, scopeSelector
status:
  used:                             # <- Ist-Stand, vom Quota-Controller fortgeschrieben
    requests.cpu: "12500m"
    limits.memory: 40Gi
    pods: "87"
    count/deployments.apps: "31"
    services.loadbalancers: "1"
# ACHTUNG, falsche Analogie: Das ist kein Meter. Es zählt, was im Cluster existiert,
# nicht, wer wie oft die API aufruft — und es trägt keinen Preis und keine Periode.
```

Dasselbe gilt für `LimitRange` (`spec.limits` mit `max`, `min`, `default`, `defaultRequest`,
`maxLimitRequestRatio`). Man kann beide als Durchsetzungsmechanismus einer extern definierten
Tarifstufe verwenden (viele Plattformteams tun genau das), aber der Tarif selbst lebt außerhalb
von KRM.

Das ist kein Versehen, sondern folgt aus dem Zuschnitt: Kubernetes ist eine
[Solution-Internal API](../foundation/SolutionInternalAPI.md) bzw. eine
[Community API](../foundation/CommunityAPI.md) für Cluster-interne Controller und Administratoren, keine
kommerzielle [Public API](../foundation/PublicAPI.md) mit externen Vertragspartnern. Managed-Kubernetes-Anbieter
(GKE, EKS, AKS und andere) haben sehr wohl Pricing Plans — pro Cluster-Stunde, pro Node, pro
Control-Plane —, aber diese leben in der Provider-eigenen Cloud-API und im Vertrag, nicht im
Kubernetes-API-Vertrag. Wer Abrechnung *auf* KRM aufsetzen will, tut das per Custom Resource
Definition und eigenem Controller, der Audit-Events oder Metriken aggregiert; das Kernmodell steuert
dazu nichts bei außer der Erweiterbarkeit.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/qualityManagementAndGovernance/PricingPlan)
