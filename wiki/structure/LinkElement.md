---
title: Link Element
kategorie: Structure
unterkategorie: Element Stereotypes
quelle: https://microservice-api-patterns.org/patterns/structure/elementStereotypes/LinkElement
---

# Link Element

*a.k.a.* *Network-Accessible Identifier*, *Address Representation*, *Hypermedia Control*

**Kurzform:** Ein [Id Element](IdElement.md), das nicht nur eindeutig, sondern zugleich
netzwerk-adressierbar ist: ein maschinen- und menschenlesbarer Zeiger auf einen anderen
Endpoint oder eine andere Operation, der direkt aufgerufen werden kann.

## Kontext

Der Kontext setzt den von [Id Element](IdElement.md) fort: ein Domänenmodell mit mehreren
verwandten Elementen unterschiedlichen Lebenszyklus wurde in mehrere fernaufrufbare
API-Abstraktionen zerlegt. Konsumenten wollen Beziehungen folgen und Folgeoperationen aufrufen
— um Detailinformation nachzuladen oder den nächsten Verarbeitungsschritt anzustoßen. Die
Adresse dieses nächsten Schritts muss irgendwo stehen. Genau das verlangen das REST-Prinzip
HATEOAS und die Cursor-Variante von [Pagination](../quality/Pagination.md).

## Problem

Wie können Endpoints und Operationen in Request- und Response-Nachrichten so referenziert
werden, dass sie fern aufgerufen werden können?

## Forces

Dieselben wie bei [Id Element](IdElement.md), im Remoting-Kontext aber verschärft:

- **Aufwand vs. Stabilität** — eingebettete Adressen sind bequem, aber jede Umstrukturierung
  von Hosts, Pfaden oder Deployment-Topologie bricht gespeicherte Links.
- **Lesbarkeit für Menschen und Maschinen** — sprechende URLs sind explorierbar; Maschinen
  brauchen dagegen typisierte Relationen, nicht geratene Pfadstrukturen.
- **Sicherheit (Vertraulichkeit)** — URLs leaken Struktur und Fachdaten, landen in Logs,
  Referrern und Caches.

## Lösung

*Link Elements* in Request- oder Response-Nachrichten aufnehmen und sie als menschen- und
maschinenlesbare, netzwerk-adressierbare Zeiger auf andere Endpoints und Operationen wirken
lassen. Optional zusätzliche [Metadata Elements](MetadataElement.md) beifügen, die die Art der
Beziehung erklären (Relationstyp, erlaubte Methode, Medientyp).

## Beispiel

Die „remote issue links“ der Atlassian-JIRA-REST-API (Quellseite): Request

```json
{
    "object": {
        "url": "http://www.mycompany.com/support?id=1",
        "title": "Crazy customer support issue"
    }
}
```

Response mit `self`-Link, über den weitere Information abgerufen werden kann:

```json
{
    "self": "http://localhost:8090/jira/rest/api/latest/issue/TST-1/remotelink/100",
    "id": 100
}
```

## Varianten

Die typisierte Variante trägt neben der Adresse Metadaten zur Relation — die Quellseite nennt
Objekte mit `href`, `rel`, `method` und `type`, als Ausprägungen HAL, Hydra/JSON-LD,
Collection+JSON und Siren.

## Konsequenzen

**Vorteile:**

- Der Client muss URL-Strukturen nicht kennen oder selbst bauen; der Provider kann sie ändern.
- Erlaubte Folgeschritte werden zustandsabhängig mitgeliefert — die Grundidee von HATEOAS.
- Cursor-Pagination und lang laufende Operationen lassen sich ohne Out-of-Band-Wissen führen.

**Nachteile / Kosten:**

- Absolute URLs binden an Host, Schema, Port, Proxy- und Deployment-Topologie.
- Größere Nachrichten; Links müssen serverseitig erzeugt und gepflegt werden.
- Gespeicherte Links veralten; sie sind schlecht als persistente Fremdschlüssel geeignet.
- Der Nutzen tritt nur ein, wenn Clients Relationen wirklich auswerten — die Quellseite hält
  fest, dass das in der Praxis selten geschieht.

## Bekannte Verwendungen

Die Quellseite ist hier ungewöhnlich nüchtern: Obwohl HATEOAS bei Fielding eine verbindliche
REST-Randbedingung ist, sei „a minority of today's Web APIs are true Hypermedia APIs“
(microservice-api-patterns.org, *Link Element*). Genannt werden die JIRA-REST-API für Remote
Issue Links, cursor-basierte Pagination bei Facebook, GitHub, Google Calendar und YouTube,
Spring HATEOAS, HAL, Hydra sowie RESTBucks.

## Verwandte Patterns

- [Id Element](IdElement.md) — Geschwistermuster; ohne Netzwerkadresse und ohne Typinformation.
- [Data Element](DataElement.md) — gemeinsamer Oberbegriff.
- [Metadata Element](MetadataElement.md) — beschreibt die Art der Beziehung.
- [Atomic Parameter](AtomicParameter.md) — übliche syntaktische Form eines Links.
- [Pagination](../quality/Pagination.md) — `next`/`prev`-Links.
- [Linked Information Holder](../quality/LinkedInformationHolder.md) — verlinkt statt einbettet,
  Gegenstück zu [Embedded Entity](../quality/EmbeddedEntity.md).
- [Link Lookup Resource](../responsibility/LinkLookupResource.md) — löst
  [Id Elements](IdElement.md) zu *Link Elements* auf.
- [State Transition Operation](../responsibility/StateTransitionOperation.md) — hypermedia-
  getriebene Zustandsübergänge; [State Creation Operation](../responsibility/StateCreationOperation.md)
  liefert lokale IDs oder vollständige Links zurück.
- [Processing Resource](../responsibility/ProcessingResource.md) — orchestrierte Folge
  verlinkter Operationen.

Außerhalb von MAP: *Linked Service* (Daigneau); *Client-side Navigation following Hyperlinks*,
*Long Running Request* und *Resource Collection Traversal* (Pautasso/Ivanchikj/Schreier).

## Bezug zu Kubernetes / KRM

Hier weicht KRM am deutlichsten von MAP ab — und zwar bewusst. Das Muster wird in der Substanz
erfüllt, aber ohne Hyperlinks: **KRM verzichtet auf HATEOAS**. Referenzen sind keine URLs,
sondern schema-typisierte Tupel.

Die Bausteine im Kern-API:

| Typ | Felder | typische Verwendung |
|---|---|---|
| `ObjectReference` | `apiVersion`, `kind`, `namespace`, `name`, `uid`, `resourceVersion`, `fieldPath` | `Event.involvedObject`, `Event.related`, `EndpointSlice.endpoints[].targetRef` |
| `LocalObjectReference` | `name` | `Pod.spec.imagePullSecrets[]`, viele `secretRef`-Felder |
| `TypedLocalObjectReference` | `apiGroup`, `kind`, `name` | `PersistentVolumeClaim.spec.dataSource` |
| `TypedObjectReference` | `apiGroup`, `kind`, `name`, `namespace` | `PersistentVolumeClaim.spec.dataSourceRef` |
| `metav1.OwnerReference` | `apiVersion`, `kind`, `name`, `uid`, `controller`, `blockOwnerDeletion` | `metadata.ownerReferences[]` |

Die beiden Extreme stehen unmittelbar nebeneinander — voll qualifiziert gegen maximal kontextabhängig:

```yaml
apiVersion: v1
kind: Event
involvedObject:                  # <- ObjectReference: trägt alles selbst
  apiVersion: v1
  kind: Pod
  namespace: default
  name: web-6d4f8b7c9-x2klm
  uid: 8f2a1c0e-3b7d-4c11-9a55-2f0e6d7b1c34
  fieldPath: spec.containers{app}   # <- zeigt sogar in das Objekt hinein
# ... reason, message, type, eventTime
---
apiVersion: v1
kind: Pod
spec:
  imagePullSecrets:
    - name: registry-creds       # <- LocalObjectReference: nur ein Name.
                                 #    Typ und Namespace ergeben sich aus der Feldposition.
  # ... containers
```

Warum keine URLs? Erstens ist die Adressierung aus dem Schema ableitbar: ein Client kombiniert
Discovery (`/apis`, `/apis/<group>/<version>`) mit der Group/Version/Resource-Zuordnung
(RESTMapper in `k8s.io/apimachinery/pkg/api/meta`) und baut den Pfad
`/apis/<group>/<version>/namespaces/<ns>/<resource>/<name>` selbst. Der Link trägt also keine
Information, die der Client nicht ohnehin hat. Zweitens sind Referenzen in KRM primär
*Deklarationen*, keine Navigationsaufforderungen: ein Controller löst sie level-triggered auf,
oft über einen Informer-Cache statt über einen HTTP-Aufruf, und beobachtet den Referenten
dauerhaft per `watch` — eine URL, die man einmal abruft, passt zu diesem Modell nicht. Drittens
wäre eine URL nicht cluster-portabel; das Tupel ist es.

Der historische Beleg ist eindeutig: `metadata.selfLink` und `ListMeta.selfLink` — die einzigen
echten Hyperlinks im Modell — sind als „legacy read-only field that is no longer populated by
the system“ markiert. KRM hat seinen einzigen Hypermedia-Zeiger aktiv wieder entfernt.

```yaml
metadata:
  name: web-6d4f8b7c9-x2klm
  namespace: default
  # selfLink: /api/v1/namespaces/default/pods/web-6d4f8b7c9-x2klm
  #   ^ der einzige echte Hyperlink des Modells; deprecated und vom System
  #     nicht mehr befüllt — hier nur zur Illustration auskommentiert.
  ownerReferences:
    - apiVersion: apps/v1
      kind: ReplicaSet
      name: web-6d4f8b7c9
      uid: 1b0c9d55-7e21-4a3f-8de0-5c9a2b6f4e77
      controller: true           # <- genau ein verwaltender Besitzer
      blockOwnerDeletion: true   # <- Steuerinformation für den Garbage Collector
```

Referenzen müssen im KRM auch nicht einzeln sein: Ein Selektor ist eine mengenwertige Referenz, die
nicht auf ein bestimmtes Objekt zeigt, sondern auf alle, die gerade passen — und deren Auflösung
sich ändert, ohne dass die Referenz angefasst wird.

```yaml
apiVersion: v1
kind: Service
spec:
  selector:                      # <- map[string]string, kein LabelSelector-Typ
    app: web
  # ... ports
---
apiVersion: apps/v1
kind: Deployment
spec:
  selector:                      # <- metav1.LabelSelector, ausdrucksstärker
    matchLabels:
      app: web
    matchExpressions:
      - key: tier
        operator: In
        values: ["frontend", "edge"]
  # ... replicas, template
```

Die interessanteste Ausprägung ist `ownerReferences`: Sie ist gleichzeitig
[Link Element](LinkElement.md) und Steuerinformation. Über `uid` ist die Referenz gegen
Namenswiederverwendung abgesichert, `controller: true` markiert genau einen verwaltenden
Besitzer, und `blockOwnerDeletion` steuert das Verhalten des Garbage Collectors bei
Foreground-Deletion. Ein deklaratives Feld ersetzt hier also das, was in einer Hypermedia-API
eine Menge von Aktionslinks wäre.

Kritisch anzumerken ist, dass KRM diese Referenztypen selbst nicht mehr für neue APIs empfiehlt:
Die Kommentare an `LocalObjectReference` und `TypedLocalObjectReference` raten von neuer
Verwendung ab („New uses of this type are discouraged“) — wegen ungenauer Validierbarkeit und
weil `kind` keine präzise Abbildung auf eine Resource ist. Empfohlen wird stattdessen ein
eng gefasster, lokal definierter Referenztyp pro Anwendungsfall. Das ist im Kern die MAP-Aussage
zu [Data Element](DataElement.md): das Vokabular der eigenen API selbst definieren, statt
generische Strukturen durchzureichen.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/elementStereotypes/LinkElement)
