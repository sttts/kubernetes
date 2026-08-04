---
title: Experimental Preview
kategorie: Evolution
unterkategorie: Lifecycle Management
quelle: https://microservice-api-patterns.org/patterns/evolution/ExperimentalPreview
---

# Experimental Preview

*a.k.a.* *Beta Program*, *Testing Sandbox*, *API Preview*

**Kurzform:** Der Provider gibt Zugriff auf eine noch in Entwicklung befindliche API auf
Best-Effort-Basis, ohne jede Zusage zu Funktionsumfang, Stabilität oder Lebensdauer — und sagt
das ausdrücklich.

## Kontext

Ein Provider entwickelt eine neue API oder eine API-Version, die sich stark von den
publizierten unterscheidet und noch intensiv verändert wird. Er will einerseits frei
umbauen können, andererseits frühen Zugriff geben, damit Clients bereits integrieren und
Rückmeldung zu Funktionsumfang und Struktur geben können.

## Problem

Wie machen Provider die Einführung einer neuen API oder API-Version für ihre Clients weniger
riskant und gewinnen Early-Adopter-Feedback, ohne das Design vorzeitig einfrieren zu müssen?

## Forces

- **Raum für Innovation und neue Features**, typischerweise iterativ und inkrementell
  entwickelt.
- **Frühes Feedback** für den Provider.
- **Fokussierung des Aufwands** — in der Frühphase z.B. API-Governance-Aufwand vermeiden.
- **Frühe Lernmöglichkeiten** für Konsumenten.
- **Stabilitätsbedürfnis** aus Client-Sicht — die Gegenkraft zu allem oben Genannten.

## Lösung

Zugriff auf die API auf Best-Effort-Basis gewähren, ohne Zusagen zu Funktionalität, Stabilität
und Lebensdauer zu machen. Diesen fehlenden Reifegrad klar und explizit kommunizieren, um
Client-Erwartungen zu steuern.

Entscheidend ist der zweite Satz: das Muster ist nicht „unfertige API veröffentlichen“, sondern
„unfertige API veröffentlichen **und als solche kennzeichnen**“.

## Beispiel

Die Quellseite skizziert ein Unternehmen, das bisher nur eine Web-Oberfläche für seinen
Build-and-Deploy-Cloud-Dienst anbietet. Großkunden fordern eine API zum Auslösen und Verwalten
von Builds sowie für Benachrichtigungen über Build-Zustände. Weil das Unternehmen bisher keine
APIs angeboten hat und entsprechend wenig Erfahrung besitzt, wählt es einen *Experimental
Preview* und verbessert die API kontinuierlich anhand des Feedbacks der Early Adopter.

## Konsequenzen

**Vorteile:**

- Designfehler werden entdeckt, solange sie noch billig zu beheben sind.
- Der Provider behält volle Änderungsfreiheit, ohne Vertrauen zu verspielen — weil er nichts
  versprochen hat.
- Early Adopter bekommen Vorlauf für ihre eigene Planung.

**Nachteile / Kosten:**

- Early Adopter tragen das Risiko und zahlen bei jeder Änderung Migrationskosten.
- Der Übergang in den produktiven Betrieb muss aktiv gemanagt werden — sonst bleibt die API
  dauerhaft „beta“ und die faktische Nutzung schafft eine Zusage, die nie gegeben wurde.
- Paralleler Betrieb von Preview und Produktion kostet Infrastruktur.
- Feedback aus einer kleinen, selbstselektierten Nutzergruppe ist verzerrt.

## Bekannte Verwendungen

GitHub bietet API Previews (etwa für die GraphQL-Unterstützung), Facebook betreibt
Beta-Programme wie das Audience Network SDK Beta, eBay kennt „Experimental APIs“. Google hatte
den Ruf, in langen Beta-Phasen steckenzubleiben. Sandboxes und Beta-Programme sind bei
Cloud-Providern verbreitet. James Higginbotham empfiehlt in einem Tyk-Blogpost die
Stabilitätsstufen *Experimental*, *Pre-Release*, *Supported*, *Deprecated*, *Retired* und rät,
unterstützte und nicht unterstützte Operationen strikt zu trennen.

## Verwandte Patterns

- [Aggressive Obsolescence](AggressiveObsolescence.md) — die nächststärkere Zusage; *Experimental
  Preview* ist das schwächste Commitment der Kategorie.
- [Two in Production](TwoInProduction.md) — beim Übergang in Produktion muss ein anderes
  Lifecycle-Muster gewählt werden; in der *N in Production*-Variante kann ein Preview eine der
  parallelen Versionen sein.
- [Limited Lifetime Guarantee](LimitedLifetimeGuarantee.md),
  [Eternal Lifetime Guarantee](EternalLifetimeGuarantee.md) — mögliche Zielzustände nach dem
  Preview.
- [Version Identifier](VersionIdentifier.md) — kann verwendet werden, ist hier nicht zwingend.
- [API Description](../foundation/APIDescription.md) — muss klar ausweisen, welche Version
  Preview und welche produktiv ist.
- [API Key](../structure/APIKey.md) — gezielte Vergabe von Preview-Zugängen an ausgewählte
  Clients.

## Bezug zu Kubernetes / KRM

Kubernetes erfüllt das Muster institutionalisiert und in zwei getrennten Ausprägungen, die man
nicht verwechseln sollte.

**Alpha-API-Versionen.** Eine Ressource in `v1alpha1` ist ein *Experimental Preview* im
Wortsinn: keinerlei Kompatibilitätszusage, Feldumbenennungen und Semantikänderungen zwischen
zwei Minor-Releases erlaubt, Entfernung jederzeit möglich, keine Migrationspflicht des
Providers. Der Reifegrad steckt dabei direkt im
[Version Identifier](VersionIdentifier.md) — der Client sieht ihn in jedem einzelnen Objekt und
kann sich nicht darüber täuschen:

```yaml
apiVersion: resource.k8s.io/v1alpha3       # <- der Reifegrad steht im Objekt
kind: DeviceTaintRule
metadata:
  name: gpu-maintenance
spec:
  deviceSelector:
    driver: gpu.example.com
    pool: node-01
    # ... device
  taint:
    key: maintenance.example.com/planned
    value: "true"
    effect: NoExecute                      # <- daneben: NoSchedule, None
    # ... timeAdded
status:
  conditions:
    - type: EvictionInProgress
      status: "False"
      # ... lastTransitionTime, reason, message
```

Wie ernst „keine Kompatibilitätszusage“ gemeint ist, zeigt derselbe Typ im Repo: zwei
Selektorfelder wurden ersatzlos entfernt und stehen nur noch als auskommentierter Grabstein im
Quelltext.

```go
// staging/src/k8s.io/api/resource/v1alpha3/types.go
// Tombstoned since 1.35 because it turned out that supporting this in all cases
// would depend on copying the device attributes into the ResourceClaim allocation
// result. [...]
//
// DeviceClassName *string `json:"deviceClassName,omitempty" protobuf:"bytes,1,opt,name=deviceClassName"`
```

In einer GA-Version wäre das undenkbar; übrig bleibt nur die auskommentierte Zeile, damit die
Protobuf-Feldnummer nicht versehentlich neu vergeben wird. Genau die „klare und explizite
Kommunikation des fehlenden Reifegrads“, die die Lösung des Musters fordert, ist hier
syntaktisch erzwungen statt dokumentiert.

**Default-off.** Alpha-Versionen werden nicht serviert, solange der Betreiber sie nicht
freischaltet. `DefaultAPIResourceConfigSource()` in `pkg/controlplane/instance.go` aktiviert
ausschließlich stabile Gruppenversionen und deaktiviert Alpha- und Beta-Versionen explizit; der
Kommentar im Code ist unmissverständlich: „Don't put alpha or beta versions in the list“.
Eingeschaltet wird über `--runtime-config`:

```console
# Eine einzelne Gruppenversion freischalten:
kube-apiserver --runtime-config=resource.k8s.io/v1alpha3=true

# Sammelschalter (staging/src/k8s.io/apiserver/pkg/server/resourceconfig/helpers.go):
#   api/all   api/ga   api/beta   api/alpha
# Konkreter Eintrag schlägt Sammelschalter:
kube-apiserver --runtime-config=api/alpha=false,resource.k8s.io/v1alpha3=true
```

Ein Preview ist damit eine bewusste Betreiberentscheidung pro Cluster — MAP kennt dieses Konzept
eines „eingeschalteten Sandkastens im Produktivsystem“ nicht.

**Feature Gates.** Orthogonal dazu steuert `--feature-gates` einzelne Verhaltensweisen. Die
Stufen sind in `k8s.io/component-base/featuregate` als `PreAlpha`, `Alpha`, `Beta`, `GA` und
`Deprecated` kodiert. Ein Gate ist keine Momentaufnahme, sondern eine versionierte Leiter —
derselbe Schalter trägt seine ganze Reifungsgeschichte:

```go
// pkg/features/kube_features.go
DRADeviceTaints: {
	{Version: version.MustParse("1.33"), Default: false, PreRelease: featuregate.Alpha},
	{Version: version.MustParse("1.36"), Default: true,  PreRelease: featuregate.Beta},  // <- Default kippt
},
```

```console
kube-apiserver \
  --runtime-config=resource.k8s.io/v1alpha3=true \   # <- API-Version servieren
  --feature-gates=DRADeviceTaints=true               # <- Verhalten einschalten
```

Beide Schalter sind nötig und unabhängig: der eine öffnet den Endpunkt, der andere das
Verhalten dahinter. Alpha-Gates sind per Default aus, und `featuregate` erzwingt beim
Registrieren eine monotone Reifung — Alpha darf nicht auf ein früheres Beta oder GA folgen,
`Deprecated` ist Endzustand. Ein neues Feld einer bestehenden GA-Ressource kann so als Alpha
ausgeliefert werden, ohne die Version der Gruppe anzufassen: hinter geschlossenem Gate wird das
Feld beim Schreiben ausgefiltert. Das ist eine Feinkörnigkeit, die MAP erst bei
[Aggressive Obsolescence](AggressiveObsolescence.md) für die Gegenrichtung vorsieht.

Beta ist bei Kubernetes historisch der wackligste Punkt: Beta-APIs wurden früher default-on
serviert und dadurch faktisch produktiv genutzt, obwohl das Commitment deutlich schwächer ist
als bei GA — das Muster wurde durch Betriebsrealität ausgehöhlt. Inzwischen ist das korrigiert:
in `pkg/controlplane/instance.go` stehen sämtliche Beta-Gruppenversionen in den Listen
`betaAPIGroupVersionsDisabledByDefault` bzw. `genericBetaAPIGroupVersionsDisabledByDefault` und
müssen wie Alpha-Versionen per `--runtime-config` eingeschaltet werden. Zusage und tatsächliche
Nutzung sind damit wieder in Deckung.

---
[← Index](../README.md) · [Kategorie Evolution](../meta/category-evolution.md) · [Quelle](https://microservice-api-patterns.org/patterns/evolution/ExperimentalPreview)
