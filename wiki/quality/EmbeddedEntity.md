---
title: Embedded Entity
kategorie: Quality
unterkategorie: Reference Management
quelle: https://microservice-api-patterns.org/patterns/quality/referenceManagement/EmbeddedEntity
---

# Embedded Entity

*a.k.a.* Inlined Entity Data, Embedded Document (Nesting)

**Kurzform:** Statt auf verwandte Informationselemente zu verweisen, werden deren Daten direkt in
die Nachricht eingebettet — der Client bekommt alles in einem Aufruf.

## Kontext

Die benötigte Information ist strukturiert und besteht aus mehreren Elementen, die zueinander in
Beziehung stehen. Ein Stammdatensatz wie ein Kundenprofil *enthält* Kontaktdaten, Adressen und
Telefonnummern; ein periodischer Ergebnisbericht *aggregiert* Einzeltransaktionen. Clients
verarbeiten beim Erstellen von Requests oder beim Auswerten von Responses mehrere dieser
zusammenhängenden Elemente.

## Problem

Wie vermeidet man mehrere Nachrichten, wenn der Empfänger Einblick in mehrere verwandte
Informationselemente braucht?

## Forces

- **Performance und Skalierbarkeit** — weniger Roundtrips vs. größere Nachrichten.
- **Modifizierbarkeit und Flexibilität** — eingebettete Struktur ist schwerer zu ändern.
- **Datenqualität** — wird ein Snapshot oder der aktuelle Stand geliefert?
- **Datenschutz** — Einbetten liefert Daten mit, die der Client vielleicht nicht sehen dürfte.
- **Datenaktualität und Konsistenz** — Teile mit unterschiedlicher Änderungsrate werden gemeinsam
  ausgeliefert.

Alle Beziehungen zu traversieren, um jede potenziell interessante Information mitzuliefern, führt zu
komplexen Repräsentationen und großen Nachrichten. Dass alle Empfänger denselben Inhalt brauchen,
ist unwahrscheinlich und schwer sicherzustellen.

## Lösung

Für jede Beziehung, der der Client folgen möchte, wird ein [Data Element](../structure/DataElement.md)
in die Nachricht eingebettet, das die Daten des Beziehungsziels enthält. Dieses *Embedded Entity*
sitzt innerhalb der Repräsentation der Beziehungsquelle.

## Beispiel

Antwort des `Customer Core`-Service von Lakeside Mutual; `customerProfile` und
`customerInteractionLog` sind vollständig enthalten, es gibt keine URIs auf andere Ressourcen.
`customerProfile` bettet seinerseits `currentAddress` und `moveHistory` ein:

```http
GET http://localhost:8080/customers/a51a-433f-979b-24e8f0

{
  "customer": { "id": "a51a-433f-979b-24e8f0" },
  "customerProfile": {
    "firstname": "Robbie",
    "lastname": "Davenhall",
    "birthday": "1961-08-11T23:00:00.000+0000",
    "currentAddress": {
      "streetAddress": "1 Dunning Trail",
      "postalCode": "9511",
      "city": "Banga"
    },
    "email": "rdavenhall0@example.com",
    "phoneNumber": "491 103 8336",
    "moveHistory": [{
      "streetAddress": "15 Briar Crest Center",
      "postalCode": "",
      "city": "Aeteke"
    }]
  },
  "customerInteractionLog": { "contactHistory": [], "classification": "??" }
}
```

## Konsequenzen

**Vorteile:**

- Weniger Aufrufe: Die Information ist da, der Client muss nicht nachfragen.
- Weniger Endpoints, weil für eingebettete Teile kein eigener Zugriffspunkt nötig ist.
- Bei tatsächlich benötigten Daten weniger Gesamt-Bandbreite als viele kleine Nachrichten mit je
  eigenem Header- und Metadaten-Overhead.
- Konsistenter Snapshot: Alle Teile stammen aus demselben Lesevorgang.

**Nachteile / Kosten:**

- Größere Nachrichten, längere Übertragungszeiten.
- Schwer vorhersehbar, was verschiedene Clients brauchen — die Tendenz geht zu "lieber mehr
  einbetten", besonders bei [Public APIs](../foundation/PublicAPI.md) mit unbekannten Clients.
- Ungenutzte Daten verbrauchen Bandbreite.
- Teile mit unterschiedlicher Änderungsgeschwindigkeit (schnelldrehende Transaktionsdaten neben
  unveränderlichen Stammdaten) machen Caching wirkungslos, weil sich immer *irgendetwas* ändert.
  In dem Fall ist [Linked Information Holder](LinkedInformationHolder.md) — ggf. plus
  [Conditional Request](ConditionalRequest.md) für den verlinkten Teil — die bessere Wahl.
- Einmal in der [API Description](../foundation/APIDescription.md) veröffentlicht, ist ein Embedded
  Entity kaum rückwärtskompatibel entfernbar.

## Bekannte Verwendungen

- **GitHub API v3**: Das Abrufen eines Issues liefert die vollständige Information zum zugeordneten
  Milestone mit.
- **Twitter REST API**: Ein Tweet enthält die komplette User-Information inklusive Follower-Zahl.
- **Microsoft Graph API**: `user`-Repräsentationen enthalten strukturierte Unter-Entitäten;
  `List events` liefert ein `attendees`-Array mit `type` und `status` pro Eintrag.
- Verbreitet in Enterprise-Informationssystemen und Master-Data-Management-Produkten.

## Verwandte Patterns

- [Linked Information Holder](LinkedInformationHolder.md) — die komplementäre Lösung desselben
  Reference-Management-Problems und zugleich die Alternative.
- [Wish List](WishList.md) und [Wish Template](WishTemplate.md) — erlauben dem Client, den Umfang
  des Eingebetteten zu steuern, und entschärfen so das Overfetching.
- [Pagination](Pagination.md) — begrenzt eingebettete Kollektionen.
- [Conditional Request](ConditionalRequest.md) — Caching-Gegenmittel bei gemischten Änderungsraten.
- [Operational Data Holder](../responsibility/OperationalDataHolder.md) referenziert
  [Master Data Holder](../responsibility/MasterDataHolder.md) per Definition; diese Referenzen
  können eingebettet oder verlinkt werden.
- [Retrieval Operation](../responsibility/RetrievalOperation.md) — entscheidet konkret zwischen
  Einbetten und Verlinken.
- [Parameter Tree](../structure/ParameterTree.md) / [Parameter Forest](../structure/ParameterForest.md)
  — die Strukturformen, in denen Eingebettetes auftritt.

## Bezug zu Kubernetes / KRM

KRM nutzt beide Patterns, und die Grenze verläuft erstaunlich scharf: **Eingebettet wird, was
keinen eigenen Lebenszyklus hat; referenziert wird, was einen hat.** Ein Container ist kein
eigenständiges API-Objekt und wäre außerhalb des Pods sinnlos, also steht er vollständig im
`PodSpec`. Genauso wird `spec.template` (ein `PodTemplateSpec`) in Deployment, StatefulSet,
DaemonSet und Job eingebettet, inklusive vollständiger `metadata` und `spec` des zukünftigen Pods.
KRM erfüllt das Pattern hier also **genauso**, nur ohne dass es je eine Alternative gegeben hätte:
Es sind Value Objects.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  # ... namespace, labels, uid, resourceVersion
spec:
  replicas: 3
  template:                     # <- PodTemplateSpec: der ganze zukünftige Pod, eingebettet
    metadata:
      labels:
        app: web
    spec:
      containers:
        - name: app             # <- listMapKey: macht die eingebettete Liste feldgranular mergebar
          image: nginx:1.27
          ports:
            - containerPort: 8080
          resources:
            requests:
              cpu: 100m
          # ... env, volumeMounts, livenessProbe, securityContext
      volumes:
        - name: cache
          emptyDir: {}          # <- Value Object ohne eigenen Lebenszyklus
      # ... tolerations, affinity, nodeSelector
```

Der Preis dieser Einbettung ist in Kubernetes gut sichtbar. Ein `PodSpec` ist eines der größten
Schemata der API, und weil ein Deployment ein PodTemplate einbettet, das wiederum Container
einbettet, entsteht eine tiefe, schwer evolvierbare Struktur — jedes neue Container-Feld muss durch
alle Workload-Controller propagiert werden. Genau der von MAP genannte Nachteil "einmal eingebettet,
kaum entfernbar" gilt in voller Härte: Felder wie `spec.containers[].securityContext` können
deprecated, aber praktisch nie gelöscht werden.

Kubernetes hat zusätzlich eine Eigenheit, die MAP nicht kennt: **Server-Side Apply mit
`metadata.managedFields`** macht eingebettete Strukturen feldgranular besitzbar. Listen bekommen
über `+listType=map` und `+listMapKey` (etwa `name` bei `containers`) Merge-Semantik, sodass
verschiedene Controller unterschiedliche Teile derselben eingebetteten Entität verwalten können,
ohne sich zu überschreiben. Das mildert den klassischen Nachteil eingebetteter Daten — dass sie nur
als Ganzes geschrieben werden können — erheblich.

Ein zweiter KRM-spezifischer Fall ist die Einbettung *aggregierter* Information in `status`:
`Deployment.status.replicas`, `readyReplicas`, `availableReplicas` sowie
`Pod.status.containerStatuses[]` sind eingebettete, vom Controller berechnete Zusammenfassungen
über referenzierte Objekte. Das ist Embedded Entity als Antwort auf das N+1-Problem, das die
Referenz-Variante sonst erzeugen würde — der Client muss nicht alle Pods lesen, um zu wissen, wie
viele bereit sind.

```yaml
# Dasselbe Deployment, status-Seite: Zähler über N referenzierte Pods, eingebettet.
status:
  observedGeneration: 7
  replicas: 3
  readyReplicas: 3            # <- sonst: 3 zusätzliche GETs auf die Pods
  availableReplicas: 3
  updatedReplicas: 3
  unavailableReplicas: 0
  # ... conditions
```

Nicht abgebildet ist das Client-gesteuerte Einbetten. Es gibt kein `?embed=`- oder
`?expand=`-Query-Parameter, keine Wish List für Beziehungen. Was es gibt, sind
`?fieldSelector=` und `?labelSelector=` zum Filtern *von Objekten* sowie seit einiger Zeit
`--subresource` und Table-Konvertierung für `kubectl get` — aber nie ein dynamisches Auffalten von
Referenzen. Diese Zurückhaltung ist konsistent mit dem uniformen Schema: Ein Objekt sieht für alle
Clients gleich aus, und das Zusammenführen ist Aufgabe des Controllers oder von `kubectl describe`.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/referenceManagement/EmbeddedEntity)
