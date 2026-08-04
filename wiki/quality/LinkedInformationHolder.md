---
title: Linked Information Holder
kategorie: Quality
unterkategorie: Reference Management
quelle: https://microservice-api-patterns.org/patterns/quality/referenceManagement/LinkedInformationHolder
---

# Linked Information Holder

*a.k.a.* Linked Entity, Data Reference, Compound Document (Sideloading)

**Kurzform:** Verwandte Informationselemente werden nicht mitgeliefert, sondern über ein
[Link Element](../structure/LinkElement.md) referenziert; der Client holt sie bei Bedarf mit einem
weiteren Aufruf.

## Kontext

Eine API stellt strukturierte Daten bereit, deren Elemente aufeinander verweisen — Produktstammdaten
*enthalten* Detailinformationen, ein Performance-Report *aggregiert* Einzelmessungen. Clients
arbeiten mit mehreren dieser Elemente, brauchen aber nicht immer alle in voller Tiefe.

## Problem

Wie bleiben Nachrichten klein, obwohl die API mit mehreren, sich gegenseitig referenzierenden
Informationselementen umgeht?

## Forces

Grundregel verteilter Systeme: Nachrichten dürfen nicht zu groß werden, sonst überlasten sie Netz
und Endpoint. Andererseits möchten Empfänger vielen oder allen Beziehungen folgen können. Sind die
verwandten Elemente nicht enthalten, braucht es Information über ihren *Ort*, ihren *Inhalt* und den
*Zugriff* darauf. Diese Informationsmenge muss entworfen, implementiert und evolviert werden, und
die entstehende Abhängigkeit muss verwaltet werden.

Abzuwägen sind:

- **Performance und Skalierbarkeit** — kleine Nachrichten, aber mehr Roundtrips.
- **Modifizierbarkeit und Flexibilität** — verlinkte Teile können unabhängig evolvieren.
- **Datenqualität** — der Link kann ins Leere zeigen.
- **Datenschutz** — der verlinkte Endpoint kann eigene Zugriffsregeln durchsetzen.
- **Datenaktualität und Konsistenz** — Quelle und Ziel werden zu verschiedenen Zeitpunkten gelesen.

## Lösung

Nachrichten, die mehrere verwandte Informationselemente betreffen, bekommen ein
[Link Element](../structure/LinkElement.md), das auf einen anderen API-Endpoint verweist, der das
verlinkte Element repräsentiert.

## Beispiel

Derselbe Kunde wie beim [Embedded Entity](EmbeddedEntity.md), nun mit ausgelagertem
`customerProfile` und `moveHistory`:

```http
GET http://localhost:8080/customers/a51a-433f-979b-24e8f0

{
  "customer": { "id": "a51a-433f-979b-24e8f0" },
  "links": [{
    "rel": "customerProfile",
    "href": "http://localhost:8080/customers/a51a-433f-979b-24e8f0/profile"
  }, {
    "rel": "moveHistory",
    "href": "http://localhost:8080/customers/a51a-433f-979b-24e8f0/moveHistory"
  }],
  "email": "rdavenhall0@example.com",
  "phoneNumber": "491 103 8336",
  "customerInteractionLog": { "contactHistory": [], "classification": "??" }
}
```

Das Pattern ist hier zweimal angewendet; der Rest bleibt eingebettet. Beide Patterns koexistieren
also in einer Repräsentation.

## Konsequenzen

**Vorteile:**

- Kleine, schnell übertragbare Nachrichten; der Client zahlt nur für das, was er wirklich liest.
- Verlinkte Elemente können unabhängig gecacht werden — besonders wertvoll, wenn sie sich langsamer
  ändern als die Quelle.
- Der verlinkte Endpoint kann eigene Autorisierung, eigene Evolution und eigenes
  [Rate Limit](RateLimit.md) haben.
- Der Vertrag der Quelle wird schlanker und dadurch länger stabil.

**Nachteile / Kosten:**

- N+1-Problem: Für eine Liste mit N Einträgen können N zusätzliche Aufrufe entstehen.
- Keine referenzielle Integrität — Links können brechen; das Gegenmittel ist eine
  [Link Lookup Resource](../responsibility/LinkLookupResource.md).
- Keine Konsistenzgarantie zwischen Quelle und Ziel, weil sie zu verschiedenen Zeitpunkten gelesen
  werden.
- Der Client muss die Beziehungssemantik (`rel`-Werte) kennen und implementieren.
- Link-Format und URI-Struktur werden selbst zum Vertragsbestandteil.

## Bekannte Verwendungen

- **JIRA Cloud REST API** beim Melden von Issue-Verknüpfungen ("Get issue link").
- **Microsoft Graph API**: `user`-Repräsentationen enthalten skalare und komplexe Attribute als
  "Properties", verlinken aber auf Ressourcen wie Calendar unter "Relationships".
- RESTful-HTTP-APIs auf **Richardson-Maturity-Level 3** wenden das Pattern an, wenn die
  Zustandsübergangs-Links sowohl Stamm- als auch Transaktionsdaten betreffen; Beispiel: Spring
  Restbucks.

## Verwandte Patterns

- [Embedded Entity](EmbeddedEntity.md) — das Geschwister-Pattern und die direkte Alternative.
- [Link Element](../structure/LinkElement.md) — der strukturelle Baustein.
- [Information Holder Resource](../responsibility/InformationHolderResource.md) — das typische
  Linkziel.
- [Link Lookup Resource](../responsibility/LinkLookupResource.md) — begegnet brechenden Links durch
  Indirektion.
- [Operational Data Holder](../responsibility/OperationalDataHolder.md) referenziert
  [Master Data Holder](../responsibility/MasterDataHolder.md) per Definition — entweder verlinkt
  oder eingebettet.
- [Conditional Request](ConditionalRequest.md) — macht das Nachladen billig, wenn sich nichts
  geändert hat.
- [Pagination](Pagination.md) — arbeitet selbst mit Links auf Folgeseiten.
- [Id Element](../structure/IdElement.md) — ein reiner Identifier als leichtgewichtige Referenzform.

## Bezug zu Kubernetes / KRM

Für alles mit eigenem Lebenszyklus erfüllt KRM dieses Pattern **konsequent — aber in einer eigenen
Variante: Referenziert wird per Name, nicht per URI.** Derselbe `PodSpec`, der Container einbettet,
ist auf der anderen Seite ein Geflecht aus Namensverweisen — jeder davon zeigt auf ein Objekt, das
unabhängig existiert, versioniert und autorisiert wird:

```yaml
# PodSpec, Referenzseite — jeder Name unten ist ein eigenes Objekt im selben Namespace.
spec:
  serviceAccountName: web              # <- Identität, eigenes Objekt
  containers:
    - name: app
      image: nginx:1.27
      env:
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: web-config         # <- Verweis auf eine ConfigMap ...
              key: logLevel            #    ... und auf einen Schlüssel darin
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: web-db
              key: password
      # ... volumeMounts, resources
  volumes:
    - name: config
      configMap:
        name: web-config
    - name: tls
      secret:
        secretName: web-tls            # <- heißt hier secretName, nicht name
    - name: data
      persistentVolumeClaim:
        claimName: web-data
```

Ein `PersistentVolumeClaim` verweist über `spec.storageClassName` auf eine StorageClass. Für
generische Verweise gibt es `ObjectReference` (mit `apiVersion`, `kind`, `namespace`, `name`,
`uid`, `resourceVersion`, `fieldPath`), `LocalObjectReference` und `TypedLocalObjectReference`.

Daneben steht eine Referenzform, die MAP so nicht kennt: die *mengenwertige, label-basierte*.
Ein Service nennt kein Ziel, sondern ein Prädikat, und die Zielmenge ergibt sich zur Laufzeit:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web
spec:
  selector:
    app: web                 # <- kein Name, sondern ein Prädikat über Pod-Labels
  ports:
    - port: 80
      targetPort: 8080
  # ... clusterIP, type
```

Bewusst *fehlt* dabei die URI. Ein KRM-Link ist ein Name innerhalb eines Namespace, kein `href`.
Der Grund ist Ortsunabhängigkeit: Dasselbe Manifest muss in jedem Cluster gelten, und ein
eingebetteter Hostname würde es an eine Installation binden. Der Client löst die Referenz über die
ihm bekannten Discovery-Regeln auf (`/apis/<group>/<version>/namespaces/<ns>/<resource>/<name>`).
Das ist Level-2-REST, nicht HATEOAS — Kubernetes ist ausdrücklich kein Maturity-Level-3-API.

Die Rückwärts-Referenz `metadata.ownerReferences[]` (mit `apiVersion`, `kind`, `name`, `uid`,
`controller`, `blockOwnerDeletion`) ist die interessanteste KRM-Ausprägung. Sie zeigt vom Kind zum
Eltern-Objekt statt umgekehrt, enthält die `uid` zur Absicherung gegen Namens-Wiederverwendung und
steuert Garbage Collection. Damit löst KRM genau das Integritätsproblem, das MAP als Nachteil des
Patterns nennt — allerdings nur für Besitzbeziehungen, nicht für lose Verweise.

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: web-7d9f8c
  ownerReferences:
    - apiVersion: apps/v1
      kind: Deployment
      name: web
      uid: 8f2e1c4a-...          # <- schützt gegen ein neues Deployment gleichen Namens
      controller: true           # <- genau ein Eintrag darf das sein: der zuständige Controller
      blockOwnerDeletion: true   # <- Eltern-Delete wartet, bis dieses Kind weg ist
  # ... labels, resourceVersion, generation
```

**Warum KRM Referenzen bevorzugt**, ist an vier Punkten festzumachen. Erstens unabhängige
Lebenszyklen: Eine ConfigMap überlebt den Pod, wird von mehreren Pods geteilt und separat versioniert.
Zweitens RBAC-Granularität: Secrets brauchen andere Rechte als Deployments; wären sie eingebettet,
würde jeder mit Deployment-Leserecht die Secrets mitlesen — der Datenschutz-Force des Patterns in
Reinform. Drittens eigene Watch-Streams: Jede Ressource ist einzeln beobachtbar, und Controller
reagieren gezielt auf Änderungen genau des referenzierten Objekts, statt große Objekte zu
repollen. Viertens Nachrichtengröße und etcd-Limits: etcd begrenzt Objekte auf rund 1,5 MB, ein
voll aufgefaltetes Deployment wäre schnell darüber.

**Was das kostet**, zeigt der Alltag jedes Controllers. N+1 ist real: Um zu wissen, ob ein Pod
starten kann, müssen ServiceAccount, ConfigMaps, Secrets und PVCs einzeln gelesen werden. KRM
mildert das nicht durch Auffalten, sondern durch client-seitiges Caching — `client-go`-Informer und
Lister halten pro Ressourcentyp einen watch-gespeisten lokalen Store, sodass der "Join" im Speicher
und nicht über das Netz passiert. Referenzielle Integrität gibt es nicht: Ein Pod darf eine
ConfigMap nennen, die es nicht gibt; die Admission lehnt das nicht ab. Der Fehler taucht erst
level-getriggert auf — als `CreateContainerConfigError` in `status.containerStatuses[]` oder als
Event. Das ist die bewusste Entscheidung, Validierung nicht zum verteilten Zwei-Phasen-Problem zu
machen. Und weil Referenzen namensbasiert sind, ist Umbenennen faktisch Löschen und Neuanlegen,
und "wer referenziert mich?" ist ohne vollständigen Scan aller potenziellen Referenzierer nicht
beantwortbar.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/referenceManagement/LinkedInformationHolder)
