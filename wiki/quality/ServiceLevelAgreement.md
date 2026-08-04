---
title: Service Level Agreement
kategorie: Quality
unterkategorie: Quality Management and Governance
quelle: https://microservice-api-patterns.org/patterns/quality/qualityManagementAndGovernance/ServiceLevelAgreement
---

# Service Level Agreement

*a.k.a.* Quality-of-Service Policies, Explicit and Structured Quality Goals

**Kurzform:** Ein strukturiertes, messbares Dokument neben der API Description, das Quality-of-Service-Zusagen
als testbare Service Level Objectives (SLOs) festhält — inklusive der Folgen ihrer Verletzung.

## Kontext

Für die API existiert ein Vertrag bzw. eine [API Description](../foundation/APIDescription.md): Operationen,
Request- und Response-Nachrichten mit Parametern sind spezifiziert. Das *dynamische* Verhalten ist damit
aber nicht beschrieben — weder qualitativ noch quantitativ — und der Support über den Lebenszyklus
(garantierte Lebensdauer, Mean Time to Repair) ebenfalls nicht.

## Problem

Wie erfährt ein Client die konkreten QoS-Eigenschaften einer API und ihrer Operationen? Und wie lassen
sich diese Eigenschaften sowie die Konsequenzen ihrer Nichteinhaltung messbar definieren und
kommunizieren?

## Forces

- **Business Agility:** Zu starre Zusagen fesseln den Provider an Implementierungsentscheidungen.
- **Attraktivität aus Konsumentensicht:** Ohne Zusagen kann der Client sein eigenes SLA nicht rechnen.
- **Verfügbarkeit:** Uptime ist die am häufigsten zugesagte und am schwersten sauber definierte Größe.
- **Performance und Skalierbarkeit:** Latenzzusagen brauchen eine Messstelle und ein Perzentil, sonst
  sind sie bedeutungslos.
- **Sicherheit und Datenschutz:** schwer quantifizierbar, landet daher meist in informell formulierten SLOs.
- **Regulierung und Rechtspflichten:** teils nicht verhandelbar.
- **Kosteneffizienz und Geschäftsrisiko des Providers:** Jede Zusage ist eine Wette mit Pönale.

## Lösung

Als API-Product-Owner ein strukturiertes, qualitätsorientiertes *Service Level Agreement* etablieren,
das **testbare** Service Level Objectives definiert. Zu jedem SLO gehören: die gemessene Größe (Service
Level Indicator), die Messstelle, das Zeitfenster, der Zielwert bzw. das Perzentil und das Remedy bei
Verfehlung.

## Varianten

- **SLA mit formal spezifizierten SLOs:** rechnerisch prüfbare Formeln (z. B. Monthly Uptime Percentage).
- **SLA mit informell spezifizierten SLOs:** Absichtserklärungen ohne Messvorschrift — üblich bei
  Sicherheitsaspekten.
- **Service Credits als einziges Remedy:** die verbreitetste Haftungsbegrenzung.

## Beispiel

Ein fiktiver SaaS-Anbieter eines Payroll-API sagt zu: „The payroll service has a response time of
maximally 0.93 seconds." Damit das prüfbar wird, muss die Messstelle mit definiert sein — gemessen vom
Eintreffen des Requests am API-Endpoint bis zur vollständigen Verarbeitung der Response, also *ohne*
Netzlaufzeit zwischen Provider- und Client-Endpoint. Dazu kommt der statistische Rahmen mit Remedy:
Das SLO gilt für 99 % der Requests; andernfalls erhält der Kunde 10 % Gutschrift auf die laufende
Abrechnungsperiode, muss den Vorfall aber mit Datum und Uhrzeit beim Support reklamieren.

## Konsequenzen

**Vorteile:**

- Der Client kann Verfügbarkeit und Latenz seiner eigenen Dienste kalkulieren, statt zu raten.
- Zwingt den Provider zu Messung und damit zu Erkenntnis über das eigene System.
- Schafft eine Verhandlungsbasis und begrenzt die Haftung explizit statt implizit.

**Nachteile / Kosten:**

- Messinfrastruktur, Reporting und Claim-Bearbeitung sind laufender Aufwand.
- Zusagen beschränken künftige Architekturentscheidungen.
- Schlecht definierte SLOs (fehlende Messstelle, fehlendes Perzentil) erzeugen Streit statt Klarheit.
- Die Beweislast liegt in der Praxis oft beim Kunden — was das SLO wirtschaftlich entwertet.

## Bekannte Verwendungen

Viele öffentliche Web-APIs haben *kein* explizites SLA, sondern nur Terms and Conditions ohne harte
Garantien. Public-Cloud-Anbieter dagegen schon:

- Amazon EC2 sagt eine „Monthly Uptime Percentage" zu, definiert über Minuten im Zustand
  „Region Unavailable".
- Microsoft Azure (u. a. für Functions) rechnet
  `Monthly Uptime % = (Maximum Available Minutes − Downtime) / Maximum Available Minutes × 100` und
  begrenzt die Haftung auf Service Credits.
- Google Compute Engine koppelt die Uptime-Garantie an Multi-Zone-Deployments und zählt nur
  zusammenhängende Downtime-Perioden ab einer Minute.
- SLAs im strategischen Outsourcing und Application Management (Miksovic/Zimmermann 2011), etwa für
  Helpdesk-Reaktionszeiten und Patch-Lieferung nach Schweregrad.
- Optimizely als Beispiel für informelle Sicherheits-SLOs („commercially reasonable technical and
  organizational measures").

Beyer et al. (2016) widmen SLOs und ihren Indikatoren (SLIs) ein eigenes Kapitel.

## Verwandte Patterns

- [API Description](../foundation/APIDescription.md) — das SLA begleitet sie; funktionaler und nicht-funktionaler Vertrag gehören nebeneinander.
- [Rate Limit](RateLimit.md) — dessen Details werden typischerweise ins SLA aufgenommen.
- [Pricing Plan](PricingPlan.md) — Preis und zugesagte Qualität bedingen einander.
- [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md), [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) — Lebensdauerzusagen sind SLA-Bestandteile.
- [Experimental Preview](../evolution/ExperimentalPreview.md) — der explizite Verzicht auf Zusagen.
- [Two in Production](../evolution/TwoInProduction.md), [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) — regeln, wie lange ein Client sich auf eine Version verlassen darf.
- [Public API](../foundation/PublicAPI.md) — SLAs sind vor allem hier relevant.

## Bezug zu Kubernetes / KRM

Kubernetes erfüllt das Pattern **anders (b)**, mit einer klaren Verschiebung: Im API-Vertrag selbst
gibt es **kein** SLA. Es gibt keine Zusage über Latenz oder Verfügbarkeit des Apiservers, kein Remedy,
keine Service Credits — schon deshalb, weil das Projekt Software liefert und keinen Dienst betreibt.
Was Kubernetes stattdessen sehr präzise garantiert, sind **Stabilitäts- und Lebensdauerzusagen** —
funktional der SLA-Ersatz, aber auf der Evolutionsachse statt auf der QoS-Achse.

Konkret: Die *API Deprecation Policy* legt Mindest-Support-Zeiträume pro Reifegrad fest — GA-Versionen
(`v1`) laufen mindestens 12 Monate oder drei Releases nach ihrer Deprecation weiter, Beta mindestens
neun Monate oder drei Releases, Alpha kann pro Release wegfallen. Der Reifegrad steht dabei direkt im
API-Pfad (`apps/v1`, `storage.k8s.io/v1beta1`, `resource.k8s.io/v1alpha3`), sodass ein
Client seine Risikoklasse an der URL ablesen kann — eine bemerkenswert direkte Umsetzung von
[Version Identifier](../evolution/VersionIdentifier.md) als Trägermedium der Garantie. Deprecated
Ressourcen liefern zusätzlich einen `Warning`-Header (der Installer in
`staging/src/k8s.io/apiserver/pkg/endpoints/installer.go` hängt ihn über `AddWarningsHandler` an, wenn
`deprecation.IsDeprecated` für die GVK zutrifft), sodass die Ankündigung in-band beim Client ankommt;
parallel markiert der Apiserver solche Requests im Audit-Log mit den Annotationen `k8s.io/deprecated`
und `k8s.io/removed-release`, was dem Betreiber die Migrationsplanung erlaubt.

```http
GET /apis/storage.k8s.io/v1beta1/volumeattributesclasses

HTTP/1.1 200 OK
Warning: 299 - "storage.k8s.io/v1beta1 VolumeAttributesClass is deprecated in v1.34+, unavailable in v1.37+; use storage.k8s.io/v1 VolumeAttributesClass"
```

Drei Aussagen in einem Header: seit wann deprecated, ab wann weg, was stattdessen — genau die
testbaren Bestandteile eines SLO, nur eben auf der Evolutions- statt auf der QoS-Achse. Die Zahlen
stammen nicht aus einem Dokument, sondern aus dem Code: `APILifecycleDeprecated()` und
`APILifecycleRemoved()` in `zz_generated.prerelease-lifecycle.go` der jeweiligen API-Gruppe.

Wer eigene APIs anbietet, bekommt denselben Mechanismus deklarativ — eine CRD-Version kann ihre
eigene Zusage formulieren:

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
spec:
  versions:
    - name: v1beta1
      served: true                 # <- noch bedient ...
      storage: false
      deprecated: true             # <- ... aber angekündigt: löst denselben Warning-Header aus
      deprecationWarning: "example.com/v1beta1 Widget is deprecated; use example.com/v1 Widget"
                                   # <- überschreibt den Default-Text; nur zulässig bei deprecated: true
      # ... schema, additionalPrinterColumns
    - name: v1
      served: true
      storage: true                # <- genau eine Version darf das sein
```

Der Betreiber sieht dieselbe Information aggregiert als Metrik — die einzige Stelle im Projekt, an
der eine Lebensdauerzusage messbar gegen die tatsächliche Nutzung gehalten wird:

```text
apiserver_requested_deprecated_apis{group="storage.k8s.io",version="v1beta1",resource="volumeattributesclasses",subresource="",removed_release="1.37"} 1
```

Feature Gates ergänzen das auf Verhaltensebene: Ein Gate durchläuft Alpha (default off), Beta und GA;
`k8s.io/component-base/featuregate` trägt diese Stufen samt `PreRelease` und Removal-Versionen im Code,
und die Stufe ist die Aussage darüber, wie sehr man sich auf das Verhalten verlassen darf. Die
*Version Skew Policy* schließlich ist die Verfügbarkeitszusage in verkleideter Form: Sie definiert,
welche Kombinationen aus kube-apiserver, kubelet, kube-controller-manager, kube-scheduler und
kubectl unterstützt sind (kubelet höchstens drei Minor-Versionen älter als der Apiserver, kubectl eine
Minor-Version in beide Richtungen) — damit ein rollierendes Upgrade ohne Ausfall überhaupt planbar
ist.

Die Bewertung: Kubernetes hat die *testbare* Hälfte des Patterns (klare, überprüfbare Zusagen mit
Zeiträumen) sehr gut, die *QoS-* und *Remedy*-Hälfte gar nicht. Das ist konsistent mit dem
deklarativen, level-triggered Modell: Ein Controller, der bis zur Konvergenz reconciled, braucht keine
Latenzzusage pro Aufruf — er braucht die Zusage, dass das Schema, auf das er reconciled, morgen noch
existiert. Wo QoS-Zusagen doch gebraucht werden, entstehen sie außerhalb des Kern-APIs: bei
Managed-Kubernetes-Anbietern für die Control-Plane-Uptime, oder Cluster-intern über eigene SLO-Objekte
per CRD samt Monitoring-Stack.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/qualityManagementAndGovernance/ServiceLevelAgreement)
