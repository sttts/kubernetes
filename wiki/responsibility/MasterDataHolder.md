---
title: Master Data Holder
kategorie: Responsibility
unterkategorie: Information Holder Endpoint Types
quelle: https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/MasterDataHolder
---

# Master Data Holder

*a.k.a.* Master Data Resource, Primary Data Access and Modification

**Kurzform:** Eine [Information Holder Resource](InformationHolderResource.md) für langlebige, selten
geänderte Daten, auf die viele Clients und viele andere Entitäten verweisen — mit entsprechend hohen
Qualitäts- und Schutzanforderungen.

## Kontext

Ein Domänenmodell, ein ER-Diagramm oder ein Glossar liegt vor, und es ist entschieden, einige
Entitäten als *Information Holder Resources* zu exponieren. Die Spezifikation zeigt stark
unterschiedliche Lebensdauern und Änderungszyklen. Langlebige Daten haben typischerweise viele
*eingehende* Beziehungen, kurzlebige Daten verweisen *auf* sie. Vielfach referenzierte, langlebige
Daten haben in vielen Szenarien hohe Anforderungen an Datenqualität und Datenschutz; die
Zugriffsprofile beider Datenarten unterscheiden sich substanziell.

## Problem

Wie entwerfe ich eine API, die Zugriff auf Stammdaten gibt — Daten, die lange leben, sich selten
ändern und von vielen Clients referenziert werden?

## Forces

Über die allgemeinen Kräfte der [Information Holder Resource](InformationHolderResource.md) hinaus:

- **Stammdatenqualität** — Fehler propagieren über alle eingehenden Referenzen und sind teuer zu
  korrigieren.
- **Stammdatenschutz** — Stammdaten sind oft personenbezogen und regulatorisch relevant (DSGVO).
- **Daten unter externer Kontrolle** — Stammdaten liegen häufig in einem MDM-System, das nicht der
  API gehört; die API ist dann Fassade, nicht Eigentümerin.

## Lösung

Kennzeichne eine [Information Holder Resource](InformationHolderResource.md) als dedizierten *Master
Data Holder*, der Stammdatenzugriff und -änderung so bündelt, dass Konsistenz gewahrt und Referenzen
angemessen verwaltet werden. Behandle Löschoperationen als Spezialfall von Updates — Stammdaten
werden markiert und archiviert, nicht physisch entfernt, weil eingehende Referenzen sonst brechen.

Optional bietet der Endpunkt weitere Lebenszyklusereignisse oder Zustandsübergänge sowie
domänenspezifische Zusatzverantwortung an; ein Archiv etwa zeitorientierte Abfragen, Bulk-Anlage und
Purge-Operationen.

## Konsequenzen

**Vorteile:**

- Ein einziger, verantwortlicher Ort für viel referenzierte Daten — Voraussetzung dafür, dass
  Qualitätsregeln und Zugriffsschutz zentral durchsetzbar sind.
- Änderungen sind selten, also sind aggressives Caching und
  [Conditional Requests](../quality/ConditionalRequest.md) besonders wirksam.

**Nachteile / Kosten:**

- Der Endpunkt wird zum Kopplungszentrum: viele eingehende Referenzen machen Schemaevolution und
  Verfügbarkeitsanforderungen zum Problem aller Clients.
- Löschen als Update erzeugt Aufräum-, Aufbewahrungs- und Compliance-Fragen (DSGVO-Löschpflichten
  vs. referenzielle Integrität).

## Bekannte Verwendungen

- Häfen und mögliche Routen in der DDD-Beispielanwendung „Cargo Tracking“.
- Inventory Service und Account Service im E-Commerce-Beispiel von Chris Richardsons
  Microservices-Patterns.
- Die „users“-Methodenfamilie der Slack Web API.
- Kunden mit Tarifplänen und das verwaltete Telefonienetz in der Order-Management-SOA aus
  Zimmermann et al. (2005).
- Der UID-Dienst der Schweizer Bundesverwaltung (Firmeninformationen, Firmensuche und
  Steuernummernvalidierung im öffentlichen Teil).
- Ein deutscher Automobilhersteller betreibt einen REST-Level-2-Dienst für Nutzerprofile: strenge
  Validierung, Normalisierung von Adressen und Telefonnummern auf Länderstandards, Create/Read/Update
  ohne Suche — die Suchlücke ist bewusst gesetzt, um DSGVO-konform zu bleiben.
- Terravis: föderierte Abfrage von Grundstücken, Rechten und Personen im Schweizer Grundbuch sowie
  ein Dienst für Stammdaten aller teilnehmenden Banken, Notariate und Grundbuchämter.

## Verwandte Patterns

- [Information Holder Resource](InformationHolderResource.md) — das Oberpattern.
- [Reference Data Holder](ReferenceDataHolder.md) — Alternative für unveränderliche Daten.
- [Operational Data Holder](OperationalDataHolder.md) — Alternative für kurzlebige Daten mit weniger
  eingehenden Referenzen; verweist typischerweise auf *Master Data Holders*.
- [Linked Information Holder](../quality/LinkedInformationHolder.md) — der übliche Weg, Stammdaten aus
  anderen Repräsentationen zu referenzieren statt sie einzubetten.
- [Atomic Parameter List](../structure/AtomicParameterList.md),
  [Parameter Tree](../structure/ParameterTree.md) — Nachrichtenstrukturen der Operationen.
- [Conditional Request](../quality/ConditionalRequest.md) — nutzt die Seltenheit der Änderungen aus.
- [Two in Production](../evolution/TwoInProduction.md) — weil viele Clients gleichzeitig migriert
  werden müssen.

DDD unterscheidet in seinen taktischen Patterns *nicht* zwischen Stamm- und operativen Daten; beide
können als *Entities* in *Aggregates* und *Bounded Contexts* auftreten und Teil der *Published
Language* sein. Der Begriff stammt aus der Informationsintegration und der Wirtschaftsinformatik und
spielt in OLAP, Data Warehousing und BI eine tragende Rolle.

## Bezug zu Kubernetes / KRM

KRM kennt keinen Modellierungsbegriff für Stammdaten (Fall c auf Schema-Ebene), realisiert das
Muster aber faktisch — mit einer eigenwilligen Wendung (Fall b): der Master-Data-Charakter sitzt
nicht in einer Ressource, sondern in deren **`spec`-Hälfte plus `metadata`**.

Beispiele im Cluster:

- **`Node`** — der Klassiker. Was in `metadata` und `spec` steht, ist langlebig und wird von
  Scheduler, Autoscaler und Workloads massiv referenziert; `status` und der Heartbeat im separaten
  `Lease` sind demgegenüber operativ. Genau diese Aufspaltung entspricht der Trennung, die MAP über
  zwei Endpunkttypen erreicht.

  ```yaml
  apiVersion: v1
  kind: Node
  metadata:
    name: node-1
    labels:
      kubernetes.io/hostname: node-1
      topology.kubernetes.io/region: eu-central-1
      topology.kubernetes.io/zone: eu-central-1a      # <- Stammdatum: Scheduling-Constraints hängen daran
      node.kubernetes.io/instance-type: m5.large
    # ... uid, resourceVersion, annotations
  spec:                                               # langlebig, administrativ gepflegt
    providerID: aws:///eu-central-1a/i-0a1b2c3d4e5f
    unschedulable: true
    taints:
      - key: node.kubernetes.io/unschedulable
        effect: NoSchedule
    # ... podCIDRs
  status:                                             # <- operative Hälfte, sekündlich geschrieben
    conditions:
      - type: Ready
        status: "True"
        lastHeartbeatTime: "2026-08-04T10:00:07Z"
        lastTransitionTime: "2026-07-29T08:12:44Z"
    # ... capacity, allocatable, addresses, nodeInfo, images
  ```

- **`Namespace`, `ServiceAccount`, `Role`/`ClusterRole`, `RoleBinding`/`ClusterRoleBinding`,
  `CustomResourceDefinition`, `APIService`** — administrativ gepflegt, selten geändert, von sehr
  vielen Objekten per Name referenziert.
- **`PersistentVolume`** und die von einem CSI-Treiber registrierten `CSIDriver`/`CSINode`-Objekte.

Zwei Aspekte löst KRM systematisch anders:

- **„Löschen als Update“** hat in KRM eine allgemeine Entsprechung: `metadata.finalizers` plus
  `metadata.deletionTimestamp`. Ein `DELETE` markiert das Objekt nur; erst wenn alle Finalizer
  entfernt sind, verschwindet es. Damit können abhängige Objekte referenzielle Aufräumarbeit
  leisten, bevor die Stammdatenentität wirklich weg ist — und `metadata.ownerReferences` mit
  Garbage Collection kaskadiert Löschungen in die Gegenrichtung.

  ```yaml
  apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: pvc-a3f7c2e1-9b0d-4f2a-8c31-6d5e0b7a1f42
    deletionTimestamp: "2026-08-04T10:00:00Z"   # <- das DELETE ist längst quittiert …
    finalizers:
      - kubernetes.io/pv-protection             # <- … das Objekt bleibt trotzdem sichtbar und lesbar
    # ... uid, resourceVersion, labels, annotations
  spec:
    # ... capacity, accessModes, claimRef, storageClassName
  status:
    phase: Bound
  ```

  Für Clients heißt das: ein erfolgreiches `DELETE` bedeutet nicht „weg“, sondern „zum Löschen
  vorgemerkt“. Wer referenzielle Integrität braucht, liest `metadata.deletionTimestamp` und
  behandelt das Objekt ab dann als scheidend.

- **Stabile Identität**: `metadata.uid` unterscheidet ein wiederverwendetes Objekt gleichen Namens
  von seinem Vorgänger. Für einen Stammdatenendpunkt, dessen Identität lange getragen werden muss,
  ist das genau die Eigenschaft, die MAP mit dem [Id Element](../structure/IdElement.md) fordert.

Die Qualitäts- und Schutzkräfte des Patterns adressiert KRM nicht im Endpunkt, sondern in Querschichten:
Validierung und Defaulting im API-Server (OpenAPI-Schema, CEL-Validierungsregeln, Validating/Mutating
Admission Webhooks) und Zugriffsschutz über RBAC pro Gruppe/Ressource/Namespace/Verb. Weil jede
Ressource denselben uniformen Kontrakt hat, sind diese Mechanismen einmal gebaut und gelten für alle
Stammdatenobjekte — inklusive selbstdefinierter CRDs.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/MasterDataHolder)
