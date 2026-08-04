---
title: Reference Data Holder
kategorie: Responsibility
unterkategorie: Information Holder Endpoint Types
quelle: https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/ReferenceDataHolder
---

# Reference Data Holder

*a.k.a.* Immutable Endpoint / Immutable Data Holder, Static Data Resource, Reference Data Lookup Table

**Kurzform:** Eine [Information Holder Resource](InformationHolderResource.md), die statische, für
Clients unveränderliche Daten als einzigen Bezugspunkt anbietet — nur Leseoperationen, kein Create,
Update oder Delete.

## Kontext

Die Anforderungsanalyse zeigt, dass bestimmte Daten in nahezu allen Systemteilen referenziert werden,
sich aber kaum je ändern; wenn doch, dann administrativ und nicht durch Clients im Tagesgeschäft.
Solche Daten heißen *Referenzdaten*: Maßeinheiten, Postleitzahlen, Ländercodes, Währungscodes,
Geolokationen und Ähnliches. Die Repräsentationen in Requests und Responses können solche Daten
entweder *enthalten* oder auf sie *verweisen*.

## Problem

Wie sollen Daten in API-Endpunkten behandelt werden, die an vielen Stellen referenziert werden, lange
leben und für Clients unveränderlich sind? Und wie können sie in Requests an und Responses von
[Processing Resources](ProcessingResource.md) oder
[Information Holder Resources](InformationHolderResource.md) benutzt werden?

## Forces

- **Don't Repeat Yourself (DRY)** — dieselbe Codeliste in jedem Service zu duplizieren erzeugt
  Divergenz.
- **Performance vs. Konsistenz beim Lesen** — lokale Kopien sind schnell, veralten aber; jeder
  Zugriff über das Netz ist konsistent, aber teuer.

## Lösung

Stelle einen speziellen [Information Holder Resource](InformationHolderResource.md)-Endpunkt bereit —
einen *Reference Data Holder* — als einzigen Bezugspunkt für die statischen, unveränderlichen Daten.
Biete Leseoperationen an, aber keine Create-, Update- oder Delete-Operationen.

Der Endpunkt sollte drei Zugriffsmuster unterstützen:

1. **Vollabruf** des gesamten Datensatzes, damit Clients ihn lokal kopieren können.
2. **Gefilterten Abruf** eines Teils, etwa für Autovervollständigung in einem Formular.
3. **Einzel-Lookup** eines Eintrags, etwa zur Validierung.

## Konsequenzen

**Vorteile:**

- Eine einzige Quelle der Wahrheit; die Daten sind cache- und replizierbar, weil unveränderlich.
- Der Endpunkt ist trivial abzusichern (nur Lesezugriff) und leicht hoch verfügbar zu halten.

**Nachteile / Kosten:**

- „Unveränderlich“ ist eine Annahme mit Verfallsdatum: die Einführung des Euro und die Umstellung
  der deutschen Postleitzahlen in den 1990ern sind die Standardgegenbeispiele. APIs, die mit solchen
  Abstraktionen arbeiten, müssen **Katalogversion und Datentypen (z.B. Maßeinheiten) explizit
  machen** — sonst ist eine Änderung nicht von einem Fehler unterscheidbar.
- Ein zusätzlicher Endpunkt und ein zusätzlicher Netzwerkaufruf, wo Einbetten oft genügt hätte.

## Bekannte Verwendungen

- Die Länder-Endpunkte der spendenfinanzierten API RESTCountries; analoge Dienste für Währungscodes.
- Taxonomien digitaler Bibliotheken, etwa das ACM Computing Classification System und die
  Schlagwörter von IEEE Xplore, sofern über APIs exponiert.
- Geodaten: die Overpass API für OpenStreetMap; die Standortdatenbank hinter der Open Weather Map API
  mit eigenen `city ID`s für über 200.000 Städte sowie Länder- und Postleitzahlen.
- Die vom Schweizer Bund vergebenen Unternehmens-Identifikationsnummern (UID).
- Regions- und Produkt-/Marktkategoriecodes in der Kernbanken-Integrations-SOA (Brandner et al. 2004).
- Digitale Archive: Terravis bietet auf archivierten Dokumenten ausschließlich
  [Retrieval Operations](RetrievalOperation.md) an, ebenso auf Prozess-Metadaten abgeschlossener
  Prozessinstanzen.

## Verwandte Patterns

- [Information Holder Resource](InformationHolderResource.md) — das Oberpattern.
- [Master Data Holder](MasterDataHolder.md) — Alternative: ebenfalls langlebig, aber veränderlich.
- [Operational Data Holder](OperationalDataHolder.md) — Alternative für kurzlebige Daten.
- [Retrieval Operation](RetrievalOperation.md) — die einzige hier zugelassene Operationsverantwortlichkeit.
- [Embedded Entity](../quality/EmbeddedEntity.md) — einfache statische Daten werden oft eingebettet,
  was einen eigenen *Reference Data Holder* überflüssig macht.
- [Linked Information Holder](../quality/LinkedInformationHolder.md) — die Alternative: der Link zeigt
  auf den *Reference Data Holder*.
- [Version Identifier](../evolution/VersionIdentifier.md) — für die explizite Katalogversion.
- [Conditional Request](../quality/ConditionalRequest.md) — nutzt die Unveränderlichkeit maximal aus.

## Bezug zu Kubernetes / KRM

KRM hat für Referenzdaten **kein eigenes Endpunktkonzept** (Fall c), aber ein sehr klar erkennbares
Idiom: *Klassenobjekte*, die per Name aus vielen Workloads referenziert werden.

- **`StorageClass`, `IngressClass`, `PriorityClass`, `RuntimeClass`, `DeviceClass`** sind
  clusterweit, administrativ gepflegt und werden von Workloads nur gelesen. `PriorityClass` liefert
  über `value` eine Zahl, die der Scheduler beim Zulassen in `pod.spec.priority` einsetzt — ein
  Lookup einer statischen Codeliste im Wortsinn.

  ```yaml
  apiVersion: scheduling.k8s.io/v1
  kind: PriorityClass
  metadata:
    name: system-cluster-critical
  value: 2000000000        # <- der Codelisteneintrag; flach, ohne spec/status
  globalDefault: false
  description: "Used for system critical pods that must run in the cluster, but can be moved to another node if necessary."
  ---
  # Pod-Ausschnitt: der Client referenziert nur den Namen …
  spec:
    priorityClassName: system-cluster-critical
    priority: 2000000000   # <- … die Admission löst auf und schreibt den Wert in den Pod
  ```

- **`ConfigMap` und `Secret`** sind der generische Referenzdatenbehälter. Beide haben seit
  Kubernetes 1.21 ein echtes Feld `immutable` (`*bool`); ist es gesetzt, weist der API-Server jede
  Änderung an `data`/`binaryData` zurück, und Kubelets müssen das Objekt nicht mehr überwachen — ein
  ausdrücklicher Performance-für-Unveränderlichkeit-Handel, wie ihn das Force „Performance vs.
  Konsistenz“ beschreibt.

  ```yaml
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: currency-codes
    namespace: default
  immutable: true                   # <- danach sind data/binaryData gesperrt, nur metadata bleibt änderbar
  data:
    codes.json: |
      {"CHF": 756, "EUR": 978}
      # ...
  ```

- **`kube-root-ca.crt`** — der API-Server verteilt die Cluster-CA automatisch als ConfigMap in jeden
  Namespace, damit Clients sie lokal haben statt sie einzeln abzufragen. Das ist der vom Pattern
  beschriebene „Vollabruf und lokal kopieren“-Zugriffspfad, nur push- statt pull-basiert.

  ```yaml
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: kube-root-ca.crt
    namespace: default              # <- identisch in *jedem* Namespace, angelegt vom rootcacertpublisher
    annotations:
      kubernetes.io/description: "Contains a CA bundle that can be used to verify the kube-apiserver when using internal endpoints ..."
  data:
    ca.crt: |                       # <- ein einziger Schlüssel; der Katalog wird gepusht, nicht gepollt
      -----BEGIN CERTIFICATE-----
      # ...
      -----END CERTIFICATE-----
  ```

- **`/openid/v1/jwks`** und `/.well-known/openid-configuration` liefern die Signaturschlüssel für
  ServiceAccount-Tokens — statische, nur lesbare Daten außerhalb des Ressourcenmodells.

Die drei vom Pattern geforderten Zugriffsmuster deckt KRM uniform ab: Vollabruf per `list`,
gefilterter Abruf per `labelSelector`/`fieldSelector` und Einzel-Lookup per `get <name>`. Statt
Client-seitigem Polling gibt es `watch` — der Client hält seine lokale Kopie über einen Informer
aktuell, was das Konsistenzproblem der „lokalen Kopie“ elegant auflöst.

Das Force „Katalogversion explizit machen“ löst KRM nicht auf Datenebene, sondern über den
API-Kontrakt: `apiVersion` versioniert das *Schema*, während der *Inhalt* einer ConfigMap über
`metadata.resourceVersion` bzw. `metadata.generation` identifiziert wird. Ein verbreitetes
Deployment-Idiom besteht deshalb darin, den Namen der ConfigMap mit einem Inhaltshash zu suffigieren
und so Unveränderlichkeit über die Objektidentität statt über ein Versionsfeld zu erzwingen.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/ReferenceDataHolder)
