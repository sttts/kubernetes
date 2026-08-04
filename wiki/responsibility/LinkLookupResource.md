---
title: Link Lookup Resource
kategorie: Responsibility
unterkategorie: Information Holder Endpoint Types
quelle: https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/LinkLookupResource
---

# Link Lookup Resource

*a.k.a.* Address Data Holder, API Directory, Endpoint Repository, Inventory/Discovery Resource, Service Registry

**Kurzform:** Ein Endpunkt, der ausschließlich Adressen *anderer* Endpunkte hält und als
[Link Elements](../structure/LinkElement.md) herausgibt — damit Clients nicht an konkrete
Endpunktadressen gebunden sind.

## Kontext

Die Repräsentationen in Requests und Responses müssen den Informationsbedarf der Empfänger vollständig
decken und enthalten dafür manchmal Verweise auf andere API-Endpunkte. Solche Verweise direkt allen
Clients zu zeigen, erhöht die Kopplung und beschädigt Orts- und Referenzautonomie. Zwei Gründe
sprechen dagegen:

- Als Provider will ich Linkziele frei ändern können, wenn die API wächst und Anforderungen sich ändern.
- Als Client will ich Code und Konfiguration (z.B. Startprozeduren) nicht anpassen müssen, wenn sich
  Namens- und Strukturkonventionen für Links auf Providerseite ändern.

## Problem

Wie können Nachrichtenrepräsentationen auf andere — möglicherweise viele und häufig wechselnde —
API-Endpunkte und Operationen verweisen, ohne den Empfänger an deren tatsächliche Adressen zu binden?

## Forces

- **Kopplung zwischen Clients und Endpunkten.**
- **Dynamische Endpunktreferenzen** — Ziele sollen zur Laufzeit auflösbar sein.
- **Zentralisierung vs. Dezentralisierung** — eine Registry ist ein Single Point of Failure.
- **Nachrichtengrößen, Anzahl der Aufrufe, Ressourcenverbrauch** — jede Indirektion kostet einen Roundtrip.
- **Umgang mit toten Links.**
- **Anzahl der Endpunkte und API-Komplexität** — die Registry ist ein zusätzlicher Endpunkt.

## Lösung

Führe einen dedizierten *Link Lookup Resource*-Endpunkt als speziellen Typ von
[Information Holder Resource](InformationHolderResource.md) ein. Er exponiert besondere
[Retrieval Operations](RetrievalOperation.md), die einzelne oder ganze Sammlungen von
[Link Elements](../structure/LinkElement.md) zurückgeben — die aktuellen Adressen der referenzierten
API-Endpunkte.

## Beispiel

Zwei Operationen in MDSL-Notation, die Kunden-*Information Holder Resources* auffinden — einmal über
einen logischen Namen, einmal über Filterkriterien:

```
endpoint type LinkLookupResourceInterface
 exposes
  operation lookupInformationHolderByLogicalName
    expecting payload
      > "name": ID
    delivering payload
      > "endpointAddress": URI

  operation lookupInformationHolderByCriteria
    expecting payload  { "filter": P }
    delivering payload { > "uri": URI* }   // 0..m

API provider CustomerLookupResource
  offers LinkLookupResourceInterface
```

Liefert die Operation mehrere Treffer gleichen Typs, wird die *Link Lookup Resource* zu einer
*Collection Resource* im Sinne des RESTful Web Services Cookbook.

## Konsequenzen

**Vorteile:**

- Der Provider kann Adressen und Namenskonventionen ändern, ohne Clients anzufassen — Ortsautonomie.
- Ein logischer Name genügt dem Client; die Auflösung ist ein Implementierungsdetail des Providers.

**Nachteile / Kosten:**

- Eine zusätzliche Indirektion: mehr Aufrufe, mehr Latenz, ein Endpunkt mehr im Kontrakt.
- Die Registry wird zum kritischen Pfad; ihre Verfügbarkeit begrenzt die des Gesamtsystems.
- Tote Links werden zwar zentral, aber nicht automatisch korrekt behandelt — Caching-Strategie und
  Gültigkeitsdauer der Auflösung müssen vereinbart werden.

## Bekannte Verwendungen

- Das `Cargo Repository` im Cargo-Aggregate der DDD-Beispielanwendung mit zwei Find-Operationen.
- Slack: die „object types“ im veröffentlichten OpenAPI-Kontrakt, sofern als Endpunkte exponiert.
- Terravis betreibt intern einen Lookup-Dienst, der zu einem *Business Partner Identifier* (BPID)
  alle bekannten Endpunkte des identifizierten Partners liefert. Zusätzlich muss jede Partei einen
  Endpunkt bereitstellen, der alle unterstützten APIs samt
  [Version Identifier](../evolution/VersionIdentifier.md) und tatsächlicher Adresse zurückgibt;
  Terravis fragt ihn täglich ab und cacht das Ergebnis für den Tag.
- WSIL und UDDI aus der SOAP-Ära gelten als überholte Umsetzungen desselben Konzepts.

## Verwandte Patterns

- [Information Holder Resource](InformationHolderResource.md) — das Oberpattern; die Lookup-Ergebnisse
  zeigen häufig auf solche Endpunkte, prinzipiell aber auf jeden Endpunkttyp.
- [Link Element](../structure/LinkElement.md) — das zurückgegebene Strukturelement.
- [Id Element](../structure/IdElement.md) — eine [Retrieval Operation](RetrievalOperation.md) kann
  *Id Elements* liefern, die indirekt auf Endpunkte zeigen; die *Link Lookup Resource* verwandelt sie
  in *Link Elements*.
- [Retrieval Operation](RetrievalOperation.md) — die Operationsverantwortlichkeit dieses Endpunkts.
- [Data Transfer Resource](DataTransferResource.md) — hält beliebige Daten, während eine
  *Link Lookup Resource* nur Metadaten exponiert.
- [Linked Information Holder](../quality/LinkedInformationHolder.md) — das Qualitätspattern, das die
  Links in Repräsentationen erzeugt, die hier aufgelöst werden.
- [API Description](../foundation/APIDescription.md) — die statische Entsprechung zur dynamischen Auflösung.

## Bezug zu Kubernetes / KRM

Kubernetes hat mit seiner **Discovery-API** eine der reinsten Umsetzungen dieses Patterns (Fall a):

- `GET /api` liefert eine `APIVersions`-Liste, `GET /apis` eine `APIGroupList` mit allen
  API-Gruppen und deren Versionen samt `preferredVersion`.

  ```console
  $ kubectl get --raw /apis | jq -c '.groups[] | {name, preferred: .preferredVersion.groupVersion}'
  {"name":"apps","preferred":"apps/v1"}
  {"name":"batch","preferred":"batch/v1"}
  {"name":"certificates.k8s.io","preferred":"certificates.k8s.io/v1"}
  ... rund zwei Dutzend weitere Gruppen
  ```

- `GET /apis/<group>/<version>` liefert eine `APIResourceList`. Jede `APIResource` darin ist
  Adresse *und* Fähigkeitsprofil des Endpunkts in einem Datensatz.

  ```console
  $ kubectl get --raw /api/v1 | jq '.resources[] | select(.name=="pods")'
  {
    "name": "pods",
    "singularName": "pod",
    "namespaced": true,
    "kind": "Pod",
    "verbs": ["create","delete","deletecollection","get","list","patch","update","watch"],
    "shortNames": ["po"],
    "categories": ["all"]
  }
  ```

  `shortNames` ist der eigentliche Lookup-Schlüssel des Patterns: `po` ist ein logischer Name ohne
  jede Adressähnlichkeit, und erst dieser Katalog macht ihn auflösbar. `categories` erlaubt denselben
  Lookup für eine ganze Gruppe von Endpunkten (`kubectl get all`).

- Seit der aggregierten Discovery (`apidiscovery.k8s.io/v2`, `APIGroupDiscoveryList`) kommt der
  gesamte Katalog in einem einzigen Roundtrip statt in einem Aufruf pro Gruppe — angefordert über
  einen Accept-Parameter (`Accept: application/json;g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList`),
  nicht über eine eigene URL. Das ist genau die Antwort auf das Force „Anzahl der Aufrufe“.

Das ist kein akademisches Konstrukt: `kubectl` *muss* diesen Lookup durchführen, bevor es irgendetwas
tun kann. Ein `kubectl get po` erfordert die Auflösung des Kurznamens `po` auf
`{Gruppe: "", Version: "v1", Ressource: "pods"}` — genau die „Lookup über logischen Namen“-Operation
aus dem Pattern. Die Implementierung heißt in apimachinery `RESTMapper` bzw.
`DeferredDiscoveryRESTMapper`, und `kubectl` cacht das Ergebnis lokal, wie es Terravis mit seinem
Tagescache tut. `kubectl api-resources` gibt den Katalog direkt aus.

Weitere Instanzen: `APIService` (registriert eine Gruppe/Version bei einem aggregierten API-Server,
mit `spec.service` als Adresse), `CustomResourceDefinition` (registriert eine Gruppe/Version im
API-Server selbst), sowie `EndpointSlice` als Adress-Holder auf Datenebene — `endpoints[].addresses`
und `ports` sind exakt die „aktuellen Adressen“, die kube-proxy und Gateway-Implementierungen
auflösen.

```yaml
apiVersion: apiregistration.k8s.io/v1
kind: APIService
metadata:
  name: v1beta1.metrics.k8s.io      # <- der Name ist der logische Schlüssel: <version>.<group>
spec:
  group: metrics.k8s.io
  version: v1beta1
  service:                          # <- die aufgelöste Adresse — kein Hostname, sondern eine Referenz
    name: metrics-server
    namespace: kube-system
    port: 443
  groupPriorityMinimum: 100
  versionPriority: 100
  # ... caBundle, insecureSkipTLSVerify
status:
  # ... conditions: Available=True, sobald der Backend-Service erreichbar ist
```

Auch hier zeigt der Eintrag nicht auf eine URL, sondern auf ein weiteres Objekt (`Service`), das
seinerseits erst zur Laufzeit zu Adressen aufgelöst wird — die Auflösungskette ist zweistufig.

Charakteristisch für KRM ist, dass die Auflösung **nicht per Link, sondern per Name** erfolgt: KRM
gibt in Repräsentationen praktisch nie URLs zurück, sondern Namen und `ObjectReference`s aus
`{apiVersion, kind, name, namespace, uid}`. Die URL wird clientseitig aus Discovery-Informationen
plus Namen konstruiert. Das ist eine bewusste Abweichung von HATEOAS und dem
[Link Element](../structure/LinkElement.md): weil das Adressschema `/apis/<group>/<version>/namespaces/<ns>/<resource>/<name>`
für *alle* Ressourcen uniform gilt, genügt ein einziger Katalog-Lookup, um jede Adresse im Cluster
zu berechnen — die Indirektion kostet dann einen Aufruf für die gesamte Sitzung statt einen pro Link.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/LinkLookupResource)
