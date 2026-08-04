---
title: Eternal Lifetime Guarantee
kategorie: Evolution
unterkategorie: Lifecycle Management
quelle: https://microservice-api-patterns.org/patterns/evolution/EternalLifetimeGuarantee
---

# Eternal Lifetime Guarantee

*a.k.a.* *Here to Stay*, *Unlimited Support Period*

**Kurzform:** Der Provider sagt zu, eine veröffentlichte API-Version niemals zu brechen oder
abzuschalten — die stärkste Zusage der Evolution-Kategorie.

## Kontext

Eine API ist veröffentlicht, die Integration mit ihren Clients ist erfolgreich und diese
Clients laufen produktiv. Mindestens einer davon kann oder will nicht auf neuere Versionen
migrieren — etwa weil er gar nicht mehr weiterentwickelt wird oder weil der Aufwand nicht
darstellbar ist.

## Problem

Wie unterstützt ein Provider Clients, die nicht auf neuere API-Versionen migrieren können oder
wollen?

## Forces

- **Keine Client-Änderungen** aufgrund von API-Änderungen.
- Der Provider muss die API dennoch **für neue Anforderungen** anderer Clients weiterentwickeln
  können.
- **Minimaler Wartungsaufwand** für die Unterstützung alter Clients.
- Fähigkeit, die **API-Infrastruktur technologisch zu erneuern**.
- Fähigkeit, **Sicherheitslücken zu schließen**.
- **Machtverhältnisse** zwischen Provider und Client anerkennen — inwieweit kann der Client
  Design und Evolution der API überhaupt steuern?

## Lösung

Als Provider garantieren, den Zugriff auf eine veröffentlichte API-Version niemals zu brechen
und niemals einzustellen.

Die Zusage betrifft den **Vertrag**, nicht notwendigerweise Implementierung und Daten: Beides
darf sich ändern, solange die Schnittstelle stabil bleibt.

## Beispiel

Die Quellseite beschreibt eine Nationalbank, die einen Dienst zum Abruf der in einem Land zu
einem gegebenen Datum gültigen ISO-Währungscodes anbietet. Weil der Dienst sehr einfach ist und
viele Nutzer erwartet werden, wird er mit einer *Eternal Lifetime Guarantee* versehen: Clients
können sich darauf verlassen, solange die Nationalbank existiert. Stabil ist dabei nur der
Vertrag — die gelieferten Daten und die Implementierung dürfen sich ändern.

## Konsequenzen

**Vorteile:**

- Maximales Vertrauen; die API wird zur verlässlichen Infrastruktur, auf der andere aufbauen.
- Null Migrationskosten auf Client-Seite, auch für nicht mehr gepflegte Clients.
- Beseitigt die Notwendigkeit, Clients überhaupt zu kennen oder zu erreichen.

**Nachteile / Kosten:**

- Der Provider verliert dauerhaft Designfreiheit; jeder Fehler im Vertrag wird permanent.
- Technologiewechsel (Protokolle, Serialisierungsformate, TLS-Versionen) müssen mit dem alten
  Vertrag vereinbar bleiben oder gesondert verhandelt werden.
- Sicherheitslücken, die im Vertrag selbst liegen, sind kaum reparabel.
- Der Wartungsaufwand akkumuliert monoton — jede Version bleibt für immer.
- In der Praxis kaum eine echte Ewigkeitszusage, sondern eine ohne festes Enddatum.

## Bekannte Verwendungen

SAP bot seinen zahlenden Kunden für offizielle und zertifizierte Remote Function Calls (RFCs)
eine Zusage, die dem Muster nahekam. Die lokalen Java-APIs für Eclipse-Erweiterungen folgen der
API Prime Directive; die Quellseite zitiert sie als: „evolving the Component API from release
to release, do not break existing clients“ (microservice-api-patterns.org, Eternal Lifetime
Guarantee) — allerdings nur für öffentliche APIs, während viele Clients faktisch interne APIs
mitbenutzen. HTTP 1.1 und XML 1.0 gelten als W3C-Standards, die seit über einem Jahrzehnt
stabil und unterstützt sind.

## Verwandte Patterns

- [Limited Lifetime Guarantee](LimitedLifetimeGuarantee.md) — dieselbe Zusage, aber befristet.
- [Two in Production](TwoInProduction.md),
  [Aggressive Obsolescence](AggressiveObsolescence.md),
  [Experimental Preview](ExperimentalPreview.md) — alle drei begrenzen die Zusage; *Eternal
  Lifetime Guarantee* ist das Extrem am anderen Ende.
- [Version Identifier](VersionIdentifier.md) — kann verwendet werden, ist hier aber nicht
  zwingend: wenn nie etwas bricht, braucht es keine Unterscheidung.
- [Semantic Versioning](SemanticVersioning.md) — würde in diesem Modell nie eine
  Major-Erhöhung sehen.
- [API Description](../foundation/APIDescription.md),
  [Service Level Agreement](../quality/ServiceLevelAgreement.md) — dort wird die Zusage
  festgeschrieben („no phase-out planned“).
- [Public API](../foundation/PublicAPI.md) — der Kontext, in dem die Zusage am meisten wert
  ist und am teuersten kommt.

## Bezug zu Kubernetes / KRM

Formal gibt Kubernetes keine Ewigkeitsgarantie: Die Deprecation Policy nennt für GA-APIs eine
*Mindest*frist von 12 Monaten oder 3 Releases (siehe
[Limited Lifetime Guarantee](LimitedLifetimeGuarantee.md)) — theoretisch also ein Enddatum.
Faktisch ist die Zusage für `v1`-Kernressourcen aber eine *Eternal Lifetime Guarantee*: `Pod`,
`Service`, `Node`, `ConfigMap`, `Secret`, `Namespace` und die übrigen Objekte der Core-Gruppe
werden seit 2015 unter `/api/v1` serviert und sind nie inkompatibel geändert worden. Ein
`Pod`-Manifest aus der Frühzeit ist gegen einen heutigen API-Server weiterhin gültig:

```yaml
apiVersion: v1              # <- seit 2015 unverändert
kind: Pod
metadata:
  name: web
  labels:
    app: web
spec:
  containers:
    - name: app
      image: nginx:1.27
      ports:
        - containerPort: 80
      volumeMounts:
        - name: config
          mountPath: /etc/nginx/conf.d
      # ... env, resources, probes, securityContext
  volumes:
    - name: config
      configMap:
        name: web-config
  restartPolicy: Always
```

Jedes Feld in diesem Manifest existiert unverändert seit `v1`; alles, was seither dazugekommen
ist — `initContainers`, `topologySpreadConstraints`, `resourceClaims`, `resources` auf
Pod-Ebene — ist optional und additiv. Auch bei den Gruppen-APIs gilt: entfernt wurden bislang
Alpha- und Beta-Versionen (`extensions/v1beta1`, `apps/v1beta1`,
`flowcontrol.apiserver.k8s.io/v1beta*`), aber keine GA-Version einer noch verwendeten
Ressource.

Warum das funktioniert, ist strukturell und nicht bloß Disziplin. Die API-Konventionen erlauben
nur additive, optionale Änderungen an einer GA-Version, und
[Two in Production](TwoInProduction.md) erzwingt verlustfreie Konversion zwischen allen
servierten Versionen über die interne Hub-Version `__internal`. Neue Funktionalität landet
deshalb typischerweise als *neues optionales Feld in derselben `v1`* oder als *neue Ressource*,
nicht als neue inkompatible Version. Genau deshalb ist Kubernetes bei
`v1` stehen geblieben, statt `v2` einzuführen: der Preis eines Bruchs wäre eine
Ökosystem-Migration von Millionen Manifesten, Helm-Charts, Operatoren und CI-Pipelines.

Die MAP-Force „Ability to upgrade API infrastructure technologies“ hat Kubernetes dabei
tatsächlich einlösen können, weil der Vertrag von der Transportform entkoppelt ist: dasselbe
`v1`-Objekt wird als JSON, YAML, Protobuf und (neu) CBOR serialisiert, ohne dass sich das
Schema ändert; Storage-Backend, Verschlüsselung und Speicherversion sind hinter
`storageVersionHash` und `StorageVersionMigration` austauschbar. Der Vertrag ist stabil, die
Implementierung darunter nicht — exakt die Unterscheidung, die das Nationalbank-Beispiel der
Quellseite macht.

Die Kehrseite sieht man ebenfalls: Fehlentscheidungen in `v1` sind permanent. `PodSpec` trägt
weiterhin das Feld `serviceAccount`, einen Alias für `serviceAccountName`. Der API-Server füllt
ihn bei *jedem* Lesezugriff aktiv mit auf, obwohl kein Client ihn mehr setzt:

```yaml
# Geschrieben:
spec:
  serviceAccountName: web
```

```yaml
# Zurückgelesen — für immer:
spec:
  serviceAccountName: web
  serviceAccount: web        # <- Alias, vom Server ergänzt
```

```go
// pkg/apis/core/v1/conversion.go — beim Schreiben gewinnt das neue Feld …
if in.ServiceAccountName == "" {
	out.ServiceAccountName = in.DeprecatedServiceAccount
}
// … beim Lesen wird das alte bedingungslos wiederhergestellt.
out.DeprecatedServiceAccount = in.ServiceAccountName
```

Der Go-Typ heißt bezeichnenderweise `DeprecatedServiceAccount` — der Name trägt die Warnung, das
Feld bleibt. Dasselbe gilt für `VolumeSource`: die In-Tree-Volume-Typen stehen weiterhin im
Schema, obwohl ihre Implementierung längst entfernt oder durch CSI ersetzt ist.

```yaml
spec:
  volumes:
    - name: legacy
      gitRepo:                     # <- „DEPRECATED“ im Godoc, trotzdem noch im Schema
        repository: https://example.com/config.git
        revision: HEAD
        # ... directory
    - name: modern
      csi:                         # <- der Ersatz, additiv daneben
        driver: example.csi.k8s.io
        # ... fsType, volumeAttributes, nodePublishSecretRef
    # ... glusterfs, rbd, flexVolume, cinder, scaleIO, storageos — alle noch gültige VolumeSource-Felder
```

Die Felder verschwinden nicht — Deprecation wirkt in einer GA-Version nur als Empfehlung. Genau
die Kostenseite, die die Quellseite dem Muster zuschreibt.

---
[← Index](../README.md) · [Kategorie Evolution](../meta/category-evolution.md) · [Quelle](https://microservice-api-patterns.org/patterns/evolution/EternalLifetimeGuarantee)
