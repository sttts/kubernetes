---
title: Version Identifier
kategorie: Evolution
unterkategorie: Versionierung
quelle: https://microservice-api-patterns.org/patterns/evolution/VersionIdentifier
---

# Version Identifier

*a.k.a.* *Explicit Versioning*, *Message-Level Version Information*

**Kurzform:** Ein expliziter Versionsindikator wird in die
[API Description](../foundation/APIDescription.md) **und** in die ausgetauschten Nachrichten
aufgenommen, damit Client und Provider erkennen, nach welchem Vertrag eine Nachricht zu
interpretieren ist.

## Kontext

Eine API läuft in Produktion und entwickelt sich weiter. Irgendwann ist eine Änderung nicht
mehr rückwärtskompatibel. Ohne erkennbare Versionsangabe merken bestehende Clients das nicht —
sie parsen die Nachricht weiterhin erfolgreich, interpretieren sie aber falsch.

## Problem

Wie kann ein Provider seine aktuellen Fähigkeiten und die Existenz möglicherweise
inkompatibler Änderungen so anzeigen, dass Clients nicht durch unbemerkte
Interpretationsfehler fehlfunktionieren?

## Forces

- **Genauigkeit und eindeutige Identifikation** der jeweils genutzten API-Version.
- **Kein versehentliches Brechen der Kompatibilität** auf semantischer Ebene — der gefährliche
  Fall ist nicht der Syntaxfehler, sondern die stillschweigend geänderte Bedeutung eines
  bestehenden Feldes.
- **Minimale Auswirkung auf die Client-Seite** bei API-Änderungen.
- **Nachvollziehbarkeit** der tatsächlich verwendeten Versionen für Governance.

## Lösung

Einen expliziten Versionsindikator einführen und ihn sowohl in der API-Beschreibung als auch in
den Nachrichten führen. Technisch wird er als [Metadata Element](../structure/MetadataElement.md)
an einer von drei Stellen untergebracht: in der Endpunktadresse, im Protokoll-Header oder im
Payload.

## Varianten

Die Quellseite unterscheidet vor allem nach Ort des Indikators:

| Ort | Beispiel | Granularität |
|---|---|---|
| URI-Pfad | `GET /v2/customers/1234` | Operation / API |
| Hostname | `Host: v2.api.service.com` | ganze API |
| Content-Type / `Accept` | `Accept: text/json+customer; version=1.0` | Repräsentationsformat |
| XML-Namespace | `http://www.ech.ch/xmlns/eCH-0134/1` | Nachrichtentyp |
| Payload-Feld | `{"version": "2.0", ...}` | Nachricht |

## Beispiel

Der von der Quellseite genannte Fall: In Version 1.0 sind alle Preise in Euro, in Version 2.0
kommt ein `currency`-Feld hinzu.

```json
{ "version": "1.0", "products": [ { "productId": "ABC123", "price": 5.00 } ] }
```

```json
{ "version": "2.0", "products": [ { "productId": "ABC123", "price": 5.00, "currency": "USD" } ] }
```

Ohne Versionsindikator würde ein alter Client die 5.00 aus der zweiten Nachricht weiterhin als
Euro lesen. Ein neu hinzugekommenes Attribut hat die Semantik eines bestehenden verändert —
genau der Fehlerfall, den das Muster verhindert.

## Konsequenzen

**Vorteile:**

- Inkompatibilitäten werden erkennbar, statt sich als Fehlinterpretation zu manifestieren.
- Voraussetzung für alle Lifecycle-Muster, insbesondere [Two in Production](TwoInProduction.md).
- Monitoring und Governance können messen, welche Versionen noch benutzt werden.

**Nachteile / Kosten:**

- Der Indikator selbst wird Vertragsbestandteil und muss gepflegt werden.
- Version im URI bricht die Identitätsstabilität von Ressourcen: dieselbe Entität hat pro
  Version eine andere URL.
- Version im `Accept`-Header ist REST-konformer, aber schlechter cachebar, schwerer zu testen
  und für Clients unbequemer.
- Versionierung verleitet dazu, statt sauberer Kompatibilität einfach eine neue Version zu
  schneiden — die Zahl parallel gepflegter Stände wächst.

## Bekannte Verwendungen

Die meisten öffentlichen Web-APIs nutzen einen einfachen, unstrukturierten *Version Identifier*
in der Request-URI; viele spiegeln ihn in Response-Headern. Facebook Graph API und die
Twitter-APIs verwenden eine `n.m`-Konvention, GitHub den HTTP-`Accept`-Header. Die
SOAP-basierten eCH-APIs realisieren das Muster über XML-Namespaces. Im WSDL-Umfeld werden
Versionen teils auf Operationsebene geführt und in Request- und Response-Nachrichten
mitgeschickt.

## Verwandte Patterns

- [Semantic Versioning](SemanticVersioning.md) — strukturiert den Indikator in `x.y.z`.
- [Two in Production](TwoInProduction.md) — setzt einen expliziten Indikator zwingend voraus.
- [Aggressive Obsolescence](AggressiveObsolescence.md),
  [Experimental Preview](ExperimentalPreview.md),
  [Limited Lifetime Guarantee](LimitedLifetimeGuarantee.md),
  [Eternal Lifetime Guarantee](EternalLifetimeGuarantee.md) — können den Indikator nutzen,
  müssen es aber nicht.
- [Id Element](../structure/IdElement.md), [Metadata Element](../structure/MetadataElement.md) —
  ein *Version Identifier* ist ein Spezialfall beider Stereotypen.
- [API Description](../foundation/APIDescription.md) — dort wird die Version dokumentiert.
- [Public API](../foundation/PublicAPI.md), [Community API](../foundation/CommunityAPI.md),
  [Solution-Internal API](../foundation/SolutionInternalAPI.md) — die Sichtbarkeit bestimmt,
  wie streng versioniert werden muss.

Außerhalb von MAP: *Tolerant Reader* (Daigneau) als komplementäre Client-seitige Technik.

## Bezug zu Kubernetes / KRM

KRM erfüllt das Muster, aber deutlich radikaler als MAP es beschreibt: die Version steht
**gleichzeitig im Pfad und im Objekt selbst**. Der Pfad lautet
`/apis/<group>/<version>/namespaces/<ns>/<resource>/<name>` (die Core-Gruppe hat den
historischen Sonderpfad `/api/v1`), und dasselbe Objekt trägt in seiner `TypeMeta` die Felder
`apiVersion: <group>/<version>` und `kind`:

```http
GET /apis/apps/v1/namespaces/default/deployments/web    # <- Version im Pfad
Accept: application/yaml
```

```yaml
apiVersion: apps/v1        # <- dieselbe Version noch einmal, jetzt im Payload
kind: Deployment
metadata:
  name: web
  namespace: default
  # ... uid, resourceVersion, generation, creationTimestamp
spec:
  # ... replicas, selector, template
```

Jedes Objekt ist damit *self-describing*: eine YAML- oder JSON-Datei ist ohne Kenntnis des
Transports vollständig interpretierbar, und ein Stream gemischter Objekte ist auflösbar. Ein
Multi-Doc-Manifest braucht keinen Kontext von außen — jedes Dokument trägt seine eigene
Zieladresse:

```yaml
apiVersion: v1             # <- Core-Gruppe ohne Gruppennamen, Pfad /api/v1/...
kind: ConfigMap
metadata:
  name: web-config
data:
  # ... Schlüssel/Werte
---
apiVersion: apps/v1        # <- /apis/apps/v1/...
kind: Deployment
metadata:
  name: web
spec:
  # ... replicas, selector, template
---
apiVersion: batch/v1       # <- /apis/batch/v1/...
kind: CronJob
metadata:
  name: web-cleanup
spec:
  schedule: "0 3 * * *"
  # ... jobTemplate
```

MAP kennt Payload-Versionierung zwar als Variante, aber nicht als durchgängige, für *alle*
Ressourcen verbindliche Konvention.

Genau diese Doppelung ist die Voraussetzung für das übrige Maschinenwerk: Der Decoder in
`k8s.io/apimachinery/pkg/runtime` wählt anhand von `apiVersion`+`kind` den passenden Go-Typ,
`schema.GroupVersionKind` und `schema.GroupVersionResource` sind die zentralen Identitätstypen
der API-Machinery, und `RESTMapper` bildet Kind auf Resource ab. Weil die Version im Objekt
steht, kann der API-Server ein per `v1beta1` geschriebenes Objekt als `v1` zurückliefern und
das Ergebnis bleibt eindeutig lesbar.

Der Indikator ist bewusst *kein* Freitext. `k8s.io/apimachinery/pkg/version` kodiert das
erwartete Schema als Regex und macht daraus eine Ordnung:

```go
// staging/src/k8s.io/apimachinery/pkg/version/helpers.go
var kubeVersionRegex = regexp.MustCompile("^v([\\d]+)(?:(alpha|beta)([\\d]+))?$")

const (
	// Bigger the version type number, higher priority it is
	versionTypeAlpha versionType = iota   // <- v1alpha1
	versionTypeBeta                       // <- v1beta1
	versionTypeGA                         // <- v1
)
```

`CompareKubeAwareVersionStrings` sortiert damit GA vor Beta vor Alpha, innerhalb einer Stufe
nach Zahl — `v2, v1, v1beta2, v1beta1, v1alpha1`. Strings, die nicht auf den Regex passen,
werden nicht abgelehnt, sondern landen grundsätzlich *unter* allen passenden. Die Konvention
wird also nicht durch Validierung erzwungen, sondern durch die Rangfolge, mit der Discovery und
`kubectl` eine bevorzugte Version auswählen. Der Versionsindikator kodiert damit nicht nur
Identität, sondern auch Reifegrad (siehe [Semantic Versioning](SemanticVersioning.md)).

Ein Detail, das MAP so nicht vorsieht: Die Version ist bei KRM keine Eigenschaft der
*Nachricht*, sondern eine *Sicht auf dasselbe Objekt*. `metadata.uid`,
`metadata.resourceVersion` und der Speicherinhalt sind versionsübergreifend identisch — siehe
[Two in Production](TwoInProduction.md).

---
[← Index](../README.md) · [Kategorie Evolution](../meta/category-evolution.md) · [Quelle](https://microservice-api-patterns.org/patterns/evolution/VersionIdentifier)
