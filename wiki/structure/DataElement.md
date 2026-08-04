---
title: Data Element
kategorie: Structure
unterkategorie: Element Stereotypes
quelle: https://microservice-api-patterns.org/patterns/structure/elementStereotypes/DataElement
---

# Data Element

*a.k.a.* *Text Element*, *Payload Content*, *Data Representation*

**Kurzform:** Fachliche Information wird über ein eigenes, für die API definiertes Vokabular
ausgetauscht, statt die internen Datenstrukturen der Implementierung nach außen durchzureichen.
Das entkoppelt Client und Provider auf der Ebene des Datenmanagements.

## Kontext

Endpoints und ihre Operationen sind identifiziert — sei es vorwärts aus einem Domänenmodell
heraus, sei es rückwärts, weil ein bestehendes System oder eine bestehende Datenbank über eine
API geöffnet werden soll. Die Basis-Strukturmuster (*Atomic Parameter*, *Parameter Tree* usw.)
sind grundsätzlich gewählt, aber die konkreten Request- und Response-Nachrichten sind noch
nicht fixiert. Die auszutauschenden Daten sind häufig Teil des persistenten Anwendungszustands
oder beeinflussen ihn.

## Problem

Wie lässt sich fachliche Information zwischen Client und Provider austauschen, ohne
provider-interne Datendefinitionen in der API offenzulegen — und ohne dass beide Seiten
datenseitig aneinandergekettet werden?

## Forces

- **Loose Coupling** — jedes exponierte Feld wird zum Vertragsbestandteil und schränkt spätere
  Änderungen an der Implementierung ein.
- **Reichhaltige Funktionalität vs. Verarbeitbarkeit und Performance** — viele, tief
  geschachtelte Elemente erhöhen Aussagekraft, aber auch Nachrichtengröße und Parsing-Aufwand.
- **Sicherheit und Datenschutz vs. Konfigurationsaufwand** — je feiner man exponiert, desto
  granularer muss autorisiert und gefiltert werden.
- **Wartbarkeit vs. Flexibilität** — generische, schwach typisierte Repräsentationen sind
  flexibel, aber schwerer zu validieren und zu verstehen.

## Lösung

Ein eigenes Vokabular von *Data Elements* für Request- und Response-Nachrichten definieren, das
die relevanten Teile der Daten aus der Business-Logik der API-Implementierung kapselt bzw. auf
sie abbildet. Die *Data Elements* gehören zum API-Vertrag, nicht zum Domänenmodell — auch wenn
sie aus ihm abgeleitet werden.

*Data Element* ist im MAP-Modell der Oberbegriff der vier Element-Stereotypen: jedes
[Metadata Element](MetadataElement.md), jedes [Id Element](IdElement.md) und jedes
[Link Element](LinkElement.md) ist zugleich ein *Data Element*, aber nicht umgekehrt.

## Beispiel

Aus der Quellseite, in der MAP-eigenen Notation (`D` = Data, `ID` = Identifier):

```
data type Customer {
    "name": ("first":D, "last":D),
    "phoneNumber":D
}

endpoint type CustomerRelationshipManagementService
  exposes
      operation getCustomer
        expecting payload "customerId": ID
        delivering payload Customer
```

`Customer` ist ein [Parameter Tree](ParameterTree.md), der ein strukturiertes (`name`) und ein
flaches (`phoneNumber`) *Data Element* kombiniert; `customerId` ist ein
[Atomic Parameter](AtomicParameter.md) in der Rolle eines *Id Elements*.

## Konsequenzen

**Vorteile:**

- Die interne Datenhaltung bleibt austauschbar; das Domänenmodell kann sich unabhängig vom
  Vertrag weiterentwickeln (*Published Language* statt durchgereichtes Schema).
- Explizite Elemente sind dokumentier-, validier- und versionierbar
  ([API Description](../foundation/APIDescription.md)).
- Sensible Attribute können gezielt weggelassen oder maskiert werden.

**Nachteile / Kosten:**

- Zusätzliche Mapping-Schicht zwischen Implementierungsmodell und Vertragsmodell — Aufwand und
  eine weitere Fehlerquelle.
- Bei redundanten Definitionen driften Vertrag und Implementierung auseinander.
- Über- oder Untergranularität schlägt direkt auf Nachrichtengröße bzw. Chattiness durch
  ([Wish List](../quality/WishList.md), [Embedded Entity](../quality/EmbeddedEntity.md)).

## Bekannte Verwendungen

Praktisch jede fachliche API. Die Quellseite nennt Tweets in der Twitter-API; Users, Projects
und Commits in der GitHub-API; das `Balance`-Objekt und die übrigen „core resources“ der
Stripe-API. Siren (`application/vnd.siren+json`) realisiert die Entity-Variante des Musters mit
*entities* und *sub-entities*. In der Versicherungs- und Bankenintegration werden Verträge,
Konten und Finanzprodukte so übertragen.

## Verwandte Patterns

- [Parameter Tree](ParameterTree.md), [Parameter Forest](ParameterForest.md) — die Strukturen,
  in denen *Data Elements* reisen.
- [Atomic Parameter](AtomicParameter.md), [Atomic Parameter List](AtomicParameterList.md) —
  einfachste syntaktische Formen eines *Data Elements*.
- [Metadata Element](MetadataElement.md), [Id Element](IdElement.md),
  [Link Element](LinkElement.md) — Spezialisierungen desselben Stereotyps.
- [Master Data Holder](../responsibility/MasterDataHolder.md),
  [Operational Data Holder](../responsibility/OperationalDataHolder.md),
  [Reference Data Holder](../responsibility/ReferenceDataHolder.md) — Endpoint-Typen, deren
  Lese- und Schreibzugriffe *Data Elements* benötigen.
- [Embedded Entity](../quality/EmbeddedEntity.md),
  [Linked Information Holder](../quality/LinkedInformationHolder.md) — Einbetten vs. Verlinken
  zusammengehöriger *Data Elements*.
- [API Description](../foundation/APIDescription.md) — dort werden alle Elemente erklärt.

Außerhalb von MAP verwandt: *Data Transfer Object* (Alur/Malks/Crupi, Fowler), *Entity* und
*Value Object* aus DDD (Evans). Die Quellseite warnt explizit davor, DDD-Konstrukte 1:1 in
API-Designs zu übersetzen.

## Bezug zu Kubernetes / KRM

KRM erfüllt das Muster, teilt es aber auf eine im MAP-Kontext ungewöhnliche Weise: der
fachliche Nutzlastteil jeder Ressource zerfällt in `spec` (gewünschter Zustand, vom Autor
geschrieben) und `status` (beobachteter Zustand, vom Controller geschrieben). Das ist keine
kosmetische Gruppierung, sondern die zentrale Erfindung des Modells: dieselbe Ressource ist
Eingabe- und Ausgabe-Nachricht zugleich, und die Trennung entscheidet, wer welches Feld
besitzen darf.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  generation: 7               # <- Server zählt bei jeder spec-Änderung hoch
  # ... namespace, uid, resourceVersion, labels
spec:                         # <- Vokabular des Autors: gewünschter Zustand
  replicas: 3
  # ... selector, template
status:                       # <- Vokabular des Controllers: beobachteter Zustand
  observedGeneration: 7       # <- genau diese spec wurde verarbeitet
  readyReplicas: 3
  conditions:
    - type: Available
      status: "True"
      reason: MinimumReplicasAvailable
      message: Deployment has minimum availability.
      lastTransitionTime: "2026-08-04T10:00:00Z"
      # ... lastUpdateTime
```

Technisch wird die Trennung über die Subresource `/status` erzwungen — ein Update auf
die Hauptressource lässt `status` unangetastet und umgekehrt, und RBAC kann auf
`<resource>/status` separat vergeben werden.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
# ... metadata
rules:
  - apiGroups: ["apps"]
    resources: ["deployments"]          # <- darf spec schreiben
    verbs: ["get", "list", "watch", "update", "patch"]
  - apiGroups: ["apps"]
    resources: ["deployments/status"]   # <- eigenes Recht für denselben Objekttyp
    verbs: ["get", "update", "patch"]
```

Die Kopplung von Schreiber und Leser wird dadurch level-triggered statt call/response: ein
Client schreibt `spec` und beobachtet später `status`, statt eine Antwort auf einen Aufruf zu
erhalten. Der übliche Rückkanal ist `status.conditions` (`metav1.Condition` mit `type`,
`status`, `reason`, `message`, `lastTransitionTime`, `observedGeneration`) — also ein
standardisiertes *Data Element* für „was der Provider gerade sieht“.

Anders als in MAP gibt es keine getrennten Request/Response-DTOs pro Operation: Es existiert
genau ein versioniertes Schema pro Group/Version/Kind, und alle Verben (`get`, `list`, `watch`,
`create`, `update`, `patch`, `delete`) arbeiten mit demselben Typ. Die MAP-Empfehlung, das
Interface-Datenmodell vom internen Modell zu entkoppeln, wird über die Trennung von
*internal type* und *external/versioned type* plus generierten Konversions- und
Defaulting-Funktionen umgesetzt (`pkg/apis/<group>/types.go` gegenüber
`staging/src/k8s.io/api/<group>/<version>/types.go`, Konversion via
`k8s.io/apimachinery/pkg/conversion`). Der Vertrag ist damit explizit ein anderes Artefakt als
die interne Repräsentation — genau die Absicht des Musters.

Für CRDs übernimmt das OpenAPI-v3-Schema in `spec.versions[].schema` diese Rolle; per Default
werden unbekannte Felder verworfen (Structural Schemas / Pruning), was das Vokabular
tatsächlich schließt statt es nur zu dokumentieren.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/elementStereotypes/DataElement)
