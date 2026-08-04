---
title: KRM vs. MAP — das Kubernetes Resource Model im Licht einer API-Pattern-Sprache
kategorie: Analyse
---

# KRM vs. MAP

Eine Einordnung des Kubernetes Resource Model (KRM) gegen die
[Microservice API Patterns](meta/overview.md). Grundlage sind die 45 Pattern-Seiten dieses Wikis,
die zugehörigen [Papers](README.md#papers) und der Quellcode in diesem Repository.

## Die These in einem Absatz

MAP ist eine **Entscheidungssprache**. Jedes Pattern beantwortet eine Entwurfsfrage, die man pro
Endpunkt neu stellen kann: Welche Rolle hat dieser Endpunkt? Welche Operationen bietet er? Wie
groß sind seine Nachrichten? Das ICSOC-Paper von 2018 beziffert allein für die Quality-Kategorie
6 Entscheidungen mit 40 Optionen und 47 Einflussfaktoren
([Details](papers/2018-icsoc-decision-guidance.md)).

KRM ist die **einmal getroffene, für alle Ressourcen eingefrorene Antwort** auf fast alle diese
Fragen. Der Entwurfsraum, den MAP aufspannt, wird auf einen Punkt zusammengezogen: ein festes
Verbset, ein festes Nachrichtengerüst, eine feste Versionierungsdisziplin. Der Preis ist
Ausdrucksverlust genau dort, wo MAP am differenziertesten ist. Der Gewinn ist, dass die gesamte
generische Maschinerie — `kubectl`, Informer, RBAC, Admission, Server-Side Apply,
Versionskonvertierung — für Ressourcen funktioniert, die es zum Zeitpunkt ihres Baus noch gar
nicht gab.

Wer eine CRD entwirft, hat deshalb null offene Entscheidungen in den Kategorien
[Responsibility](meta/category-responsibility.md) und [Evolution](meta/category-evolution.md) —
und trägt die volle Last in [Structure](meta/category-structure.md).

## Die Bilanz über alle 45 Patterns

| | Anzahl | Bedeutung |
|---|---|---|
| **(a) genauso oder stärker** | 27 | KRM setzt das Pattern um, oft rigider als MAP verlangt |
| **(b) anders gelöst** | 9 | Dieselbe Kraft, anderer Mechanismus |
| **(c) gar nicht** | 9 | Bewusst nicht vorhanden — teils Lücke, teils Nicht-Ziel |

### (a) Umgesetzt, oft übererfüllt

| Pattern | KRM |
|---|---|
| [API Description](foundation/APIDescription.md) | Discovery + OpenAPI v3 + `x-kubernetes-*` + CRD-Schema als ausführbare Spezifikation |
| [Community API](foundation/CommunityAPI.md) | RBAC + Namespaces; aber Description für alle Authentifizierten |
| [Solution-Internal API](foundation/SolutionInternalAPI.md) | Nicht-serialisierte Hub-Versionen `__internal`, Subresources als RBAC-Grenze |
| [Information Holder Resource](responsibility/InformationHolderResource.md) | Praktisch jedes Objekt; kompromissloser als MAP fordert |
| [Operational / Master / Reference Data Holder](responsibility/OperationalDataHolder.md) | Faktisch vorhanden, aber quer geschnitten: `status` vs. `spec` statt Endpunkttypen |
| [Data Transfer Resource](responsibility/DataTransferResource.md) | Der API-Server *ist* das Blackboard (`cluster-info`, CSR, PVC→PV) |
| [Link Lookup Resource](responsibility/LinkLookupResource.md) | `/api`, `/apis`, Aggregated Discovery, `RESTMapper` |
| [State Creation Operation](responsibility/StateCreationOperation.md) | `create`, aber mit vollem Objekt als Antwort |
| [Retrieval Operation](responsibility/RetrievalOperation.md) | `get`/`list` plus `watch` als Streaming-Erweiterung |
| [Atomic Parameter](structure/AtomicParameter.md) | Query-Parameter, typisiert als `metav1.ListOptions` |
| [Parameter Tree](structure/ParameterTree.md) | Erzwungen für jede Ressource |
| [Data Element](structure/DataElement.md) | `spec`/`status` — plus die KRM-Erfindung der Zuständigkeitstrennung |
| [Metadata Element](structure/MetadataElement.md) | `ObjectMeta` uniform über *alle* Ressourcen, inkl. `managedFields` |
| [Id Element](structure/IdElement.md) | `name`, `uid`, `resourceVersion`, `generation` — vier Ids mit klar getrennten Rollen |
| [API Key](structure/APIKey.md) | Bearer-Token/Client-Zertifikat; Ausstellung selbst deklarativ (`TokenRequest`, CSR) |
| [Error Report](structure/ErrorReport.md) | `metav1.Status` mit `reason`/`code`/`details`, clientseitig via `apierrors.Is*()` |
| [Embedded Entity](quality/EmbeddedEntity.md) | Value Objects (Container, PodTemplate); SSA macht sie feldgranular besitzbar |
| [Linked Information Holder](quality/LinkedInformationHolder.md) | Alles mit eigenem Lebenszyklus — per Name, nicht per URI |
| [Pagination](quality/Pagination.md) | `limit`/`continue`, ausschließlich Cursor |
| [Version Identifier](evolution/VersionIdentifier.md) | Im Pfad **und** im Objekt (`apiVersion`) |
| [Two in Production](evolution/TwoInProduction.md) | *n* Versionen, verlustfrei konvertierend — siehe unten |
| [Limited Lifetime Guarantee](evolution/LimitedLifetimeGuarantee.md) | Deprecation Policy, `Warning`-Header, Metrik |
| [Eternal Lifetime Guarantee](evolution/EternalLifetimeGuarantee.md) | Faktisch für `v1` |
| [Experimental Preview](evolution/ExperimentalPreview.md) | Alpha-Versionen + Feature Gates, default-off |
| [Aggressive Obsolescence](evolution/AggressiveObsolescence.md) | Unterhalb GA, mit maschinenlesbarer Ankündigung |

### (b) Anders gelöst

| Pattern | MAP-Mechanismus | KRM-Mechanismus |
|---|---|---|
| [Frontend Integration](foundation/FrontendIntegration.md) | Eigene, kanalspezifische API | Content Negotiation (`as=Table`) auf derselben API |
| [Backend Integration](foundation/BackendIntegration.md) | Punkt-zu-Punkt-Messaging | Sterntopologie über den API-Server, `watch`-Streams |
| [Computation Function](responsibility/ComputationFunction.md) | Eigene Operation | `create` auf ein nie persistiertes Objekt (`*Review`), `dryRun=All` |
| [State Transition Operation](responsibility/StateTransitionOperation.md) | Benannte Übergangsoperation | `spec` ändern, Controller konvergiert level-triggered |
| [Atomic Parameter List](structure/AtomicParameterList.md) | Lose Parametersignatur | Versionierter Options-Typ; Top-Level-Arrays verboten |
| [Link Element](structure/LinkElement.md) | URL im Payload | `ObjectReference` / `ownerReferences`; `selfLink` wurde sogar wieder entfernt |
| [Conditional Request](quality/ConditionalRequest.md) | ETag / `If-None-Match` / 304 | `resourceVersion` überall — und `watch` statt Polling |
| [Rate Limit](quality/RateLimit.md) | Kontingent pro Client und Zeitfenster | APF: Nebenläufigkeits-Seats + Shuffle Sharding |
| [Service Level Agreement](quality/ServiceLevelAgreement.md) | QoS-Zusagen mit Remedy | Stabilitäts- und Lebensdauerzusagen statt Latenzzusagen |

### (c) Nicht vorhanden

| Pattern | Bewertung |
|---|---|
| [Public API](foundation/PublicAPI.md) | Nicht-Ziel: der Cluster-Endpunkt ist nie public, nur die *Spezifikation* ist es |
| [Processing Resource](responsibility/ProcessingResource.md) | Bewusst vermieden; Ausnahmen sind Subresources (`pods/exec`, `pods/eviction`, `pods/binding`) |
| [Parameter Forest](structure/ParameterForest.md) | Bewusst verboten — alles wird unter eine Wurzel gezwungen |
| [Context Representation](structure/ContextRepresentation.md) | Kontext reist in Headern und `context.Context`, nicht in der Payload |
| [Wish List](quality/WishList.md) | **Echte Lücke.** Nur `PartialObjectMetadata` und `Table` als feste Projektionen |
| [Wish Template](quality/WishTemplate.md) | **Echte Lücke.** `fieldsV1` ist formal ein Wish Template — aber auf der Schreibseite |
| [Request Bundle](quality/RequestBundle.md) | **Echte Lücke.** Kein Batching; Ausnahmen: Aggregated Discovery, `SelfSubjectRulesReview` |
| [Pricing Plan](quality/PricingPlan.md) | Nicht-Ziel; `ResourceQuota` ist ein Ressourcen-, kein Aufrufkontingent |
| [Semantic Versioning](evolution/SemanticVersioning.md) | Bewusst ersetzt durch die Reifegradskala `v1alpha1`/`v1beta1`/`v1` |

## Wo KRM besser ist als das Pattern

### 1. `spec`/`status` — die Kategorie, die MAP fehlt

MAP kennt Nachrichten als Eingabe *oder* Ausgabe. KRM macht **dasselbe Objekt zu beidem** und
teilt die Schreibzuständigkeit feldweise: der Client besitzt `spec`, der Controller besitzt
`status`, durchgesetzt über die `/status`-Subresource als eigene RBAC-Ressource. Dazu kommt
`metadata.generation` gegen `status.observedGeneration` als Fortschrittsanzeige.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  generation: 7              # <- vom apiserver bei jeder spec-Änderung hochgezählt
spec:
  replicas: 5                # <- gehört dem Client
  # ... selector, template
status:
  observedGeneration: 7      # <- Controller hat generation 7 gesehen
  readyReplicas: 4           # <- ... aber erst 4 sind bereit
  conditions:
    - type: Progressing
      status: "True"
      reason: NewReplicaSetAvailable
      lastTransitionTime: "2026-08-04T09:12:31Z"
```

Durchgesetzt wird die Trennung über RBAC, nicht über Konvention:

```yaml
rules:
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "list", "watch"]      # <- darf spec lesen, nicht ändern
  - apiGroups: ["apps"]
    resources: ["deployments/status"]    # <- eigene Ressource, eigene Erlaubnis
    verbs: ["update", "patch"]
```

Das ist keine Variante eines MAP-Patterns, sondern ein fehlendes Element der Pattern-Sprache. Ihm
folgt fast alles andere: Weil ein Objekt zugleich Auftrag und Bericht ist, braucht KRM kaum
[Processing Resources](responsibility/ProcessingResource.md) und praktisch keine
[State Transition Operations](responsibility/StateTransitionOperation.md).

### 2. `watch` — der Interaktionsstil, den MAP nicht abbildet

MAP bleibt durchgehend im Request-Response-Modell.
[Conditional Request](quality/ConditionalRequest.md) optimiert das Polling; es schafft es nicht ab.
KRM schafft es ab: Ein Client listet einmal ab einer `resourceVersion` und konsumiert danach nur
noch Deltas.

```http
GET /api/v1/namespaces/shop/pods?watch=true&resourceVersion=48210
      &allowWatchBookmarks=true&sendInitialEvents=true&resourceVersionMatch=NotOlderThan
```
```json
{"type":"ADDED",   "object":{"kind":"Pod","metadata":{"name":"web-1","resourceVersion":"48211"}}}
{"type":"BOOKMARK","object":{"kind":"Pod","metadata":{"resourceVersion":"48211",
   "annotations":{"k8s.io/initial-events-end":"true"}}}}
{"type":"MODIFIED","object":{"kind":"Pod","metadata":{"name":"web-1","resourceVersion":"48260"}}}
{"type":"BOOKMARK","object":{"kind":"Pod","metadata":{"resourceVersion":"48294"}}}
```

Die `BOOKMARK`-Events tragen **keine Objektdaten**. Sie schreiben nur die `resourceVersion` fort, an
der ein Reconnect ansetzt — ein „nichts Neues"-Signal, das *keinen Request kostet* und damit
sparsamer ist als jedes `304`. `sendInitialEvents` verschmilzt den Initialabgleich mit dem Stream
und beendet ihn mit dem oben annotierten Bookmark; danach listet ein Controller nie wieder
vollständig.

Das ist der Grund, warum die drei „echten Lücken" der Quality-Kategorie —
[Wish List](quality/WishList.md), [Wish Template](quality/WishTemplate.md),
[Request Bundle](quality/RequestBundle.md) — weniger wehtun, als man erwarten würde: Wer ein
Objekt ohnehin nur einmal vollständig lädt und danach inkrementell fortschreibt, gewinnt durch
Response Shaping wenig.

### 3. *Two in Production*, aber ernsthaft

MAP verlangt für [Two in Production](evolution/TwoInProduction.md) ausdrücklich **nicht**, dass
die koexistierenden Versionen kompatibel sind. KRM verlangt das Gegenteil und erzwingt es:

- Beliebig viele Versionen gleichzeitig, nicht zwei: `resource.k8s.io` registriert heute vier
  (`v1`, `v1beta2`, `v1beta1`, `v1alpha3`, siehe `pkg/registry/resource/rest/storage_resource.go`).
- Sie sind nicht zwei Verträge, sondern **zwei Sichten auf dasselbe Objekt** — gleiche
  `metadata.uid`, gleiche `resourceVersion`, gleiche Position im Watch-Stream.
- Konvertierung sternförmig über eine interne Hub-Version (`runtime.APIVersionInternal` =
  `"__internal"`), O(n) statt O(n²) Konverter, abgesichert durch fuzzing-basierte
  Round-Trip-Tests.
- Die Speicherversion ist davon entkoppelt und migrierbar (`storageVersionHash`,
  `StorageVersionMigration`).

Dasselbe Objekt über zwei Pfade — `uid` und `resourceVersion` sind identisch. Hier eine echte
Feldumbenennung zwischen den beiden `flowcontrol`-Versionen, die in 1.26–1.28 parallel serviert
wurden:

```yaml
apiVersion: flowcontrol.apiserver.k8s.io/v1beta2   # <- ältere Sicht
kind: PriorityLevelConfiguration
metadata:
  name: workload-low
  uid: 9d5c1e3a-...-b21f
  resourceVersion: "77301"
spec:
  limited:
    assuredConcurrencyShares: 30                   # <- alter Feldname
---
apiVersion: flowcontrol.apiserver.k8s.io/v1beta3   # <- neuere Sicht, DASSELBE Objekt
kind: PriorityLevelConfiguration
metadata:
  name: workload-low
  uid: 9d5c1e3a-...-b21f
  resourceVersion: "77301"
spec:
  limited:
    nominalConcurrencyShares: 30                   # <- umbenannt, verlustfrei konvertierbar
```

Bei CRDs wird dieselbe Zusage deklarativ konfiguriert:

```yaml
spec:
  versions:
    - name: v1alpha1
      served: true             # <- wird noch beantwortet
      storage: false
      deprecated: true
      deprecationWarning: "example.com/v1alpha1 Widget is deprecated; use v1"
    - name: v1
      served: true
      storage: true            # <- genau eine Version ist Speicherversion
  conversion:
    strategy: Webhook          # <- oder None, das nur apiVersion umschreibt
    webhook:
      conversionReviewVersions: ["v1"]
      clientConfig:
        service: { namespace: example-system, name: conversion-webhook, path: /convert }
```

Der Preis steht ehrlich in der Bilanz: Ein Feld, das sich nicht verlustfrei auf die alte Version
abbilden lässt, ist nicht baubar, solange die alte Version serviert wird.

### 4. Die Beschreibung ist ausführbar

Bei MAP ist die [API Description](foundation/APIDescription.md) ein Dokument. Bei einer CRD ist
sie **Teil der Ressource** (`spec.versions[].schema.openAPIV3Schema`) und zugleich die Quelle für
Validierung, Defaulting, Pruning und Merge-Verhalten. Die `x-kubernetes-*`-Erweiterungen
transportieren genau das, was MAP als „dynamische/verhaltensbezogene Aspekte" fordert und was
OpenAPI selbst nicht kann — bis hin zu CEL-Invarianten in `x-kubernetes-validations`.

### 5. Der Error Report ist selbst ein Objekt

`metav1.Status` ist ein `Kind` wie jedes andere, mit stabilen `reason`-Konstanten und
`details.causes` samt Feldpfaden. Darüber hinaus hat KRM einen **zweiten, asynchronen Fehlerkanal**,
den MAP nicht kennt: `status.conditions` und Events. Ein Fehler, der erst bei der Reconciliation
auftritt, hat einen definierten Ort — bei MAP gäbe es ihn nur als HTTP-Antwort, die längst
verschickt ist.

## Wo KRM schlechter ist

### 1. Keine Datensparsamkeit auf der Leseseite

Das ist die härteste Lücke. Es gibt kein Sparse Fieldset, keine Feldmaske, kein Batching. Ein
`kubectl apply -f` über fünfzig Manifeste sind fünfzig HTTP-Requests. Die verfügbaren
Projektionen — `PartialObjectMetadata` und `Table` — sind fest verdrahtet und
provider-definiert.

```http
# Gibt es nicht. Es existiert nirgends ein ?fields=-Parameter.
GET /api/v1/namespaces/shop/pods?fields=metadata.name,status.phase
```

```http
# Stattdessen: zwei fest verdrahtete Projektionen über Content Negotiation.
GET /api/v1/namespaces/shop/pods
Accept: application/json;as=PartialObjectMetadataList;g=meta.k8s.io;v=v1
```
```yaml
kind: PartialObjectMetadataList     # <- nur TypeMeta + ObjectMeta, kein spec, kein status
items:
  - metadata:
      name: web-1
      resourceVersion: "48211"
      # ... labels, annotations
```

```console
$ kubectl apply -f manifests/       # 50 Objekte
# -> 50 einzelne HTTP-Requests. ApplyOptions.Run iteriert und ruft applyOneObject pro Objekt.
#    Fehlschläge sind partiell, es gibt kein Rollback.
```

Kompensiert wird das durch Informer-Caches; wo die nicht passen (kurzlebige Jobs, CLI-Tools,
Web-UIs), zahlt man voll. `kubectl get -o jsonpath` spart **kein einziges Byte** auf der Leitung.

### 2. Der Ausdrucksverlust bei Operationen

Wo ein Übergang wirklich imperativ und nicht als Zielzustand formulierbar ist, muss KRM das Modell
verbiegen. `SubjectAccessReview` ist ein RPC im Ressourcen-Kostüm — man `create`t ein Objekt,
bekommt es mit gefülltem `status` zurück, und gespeichert wird nichts:

```yaml
# Request: POST /apis/authorization.k8s.io/v1/subjectaccessreviews
apiVersion: authorization.k8s.io/v1
kind: SubjectAccessReview
spec:
  user: alice
  resourceAttributes: { namespace: shop, verb: delete, group: apps, resource: deployments }
---
# Response — dasselbe Kind, status gefüllt. Es gibt kein get, kein list, kein watch dafür.
status:
  allowed: false
  reason: 'RBAC: no rules allow user "alice" to delete deployments in namespace "shop"'
```

Zum Vergleich derselbe Wunsch, sobald er *deklarativ* formulierbar ist:

```http
# Gibt es nicht. Kein Verb für Geschäftsübergänge.
POST /apis/apps/v1/namespaces/shop/deployments/web/restart
```
```yaml
# Sondern: Zielzustand ändern, Controller konvergiert. Idempotent von Bauart.
spec:
  template:
    metadata:
      annotations:
        kubectl.kubernetes.io/restartedAt: "2026-08-04T10:00:00Z"
```

Dazu kommen `create`-Verben, die nichts erzeugen (`pods/exec`, `pods/attach`, `pods/portforward`
verlassen CRUD ganz und upgraden die Verbindung auf einen Stream) und Write-only-Kommandos als POST
auf Subresources (`pods/binding`, `pods/eviction`, `namespaces/finalize`, `serviceaccounts/token`,
`certificatesigningrequests/approval`). Jedes für sich funktioniert; zusammen sind sie das Geräusch
eines uniformen Modells, das aufgehebelt wird.

### 3. Keine referenzielle Integrität

Referenzen sind Namensstrings. Ein Pod darf eine ConfigMap nennen, die es nicht gibt; Admission
lehnt das nicht ab:

```yaml
spec:
  serviceAccountName: web
  volumes:
    - name: config
      configMap:
        name: app-config          # <- muss nicht existieren; kein Schreibfehler
---
status:
  containerStatuses:
    - name: app
      state:
        waiting:
          reason: CreateContainerConfigError    # <- der Fehler kommt erst level-getriggert
          message: 'configmap "app-config" not found'
```

Umbenennen ist faktisch Löschen und Neuanlegen, und „wer referenziert mich?" ist ohne vollständigen
Scan aller potenziellen Referenzierer nicht beantwortbar. Nur für *Besitz*beziehungen ist Integrität
gelöst:

```yaml
metadata:
  ownerReferences:
    - apiVersion: apps/v1
      kind: ReplicaSet
      name: web-7d9f
      uid: 3f1c...-88ab          # <- uid schützt gegen Namens-Wiederverwendung
      controller: true
      blockOwnerDeletion: true   # <- steuert die Reihenfolge der Garbage Collection
```

### 4. Fehlentscheidungen sind permanent

Die Kehrseite der faktischen [Eternal Lifetime Guarantee](evolution/EternalLifetimeGuarantee.md):
Deprecation wirkt in einer GA-Version nur als Empfehlung, weil die Round-Trip-Anforderung aus
*Two in Production* das Entfernen verbietet.

```yaml
spec:
  serviceAccountName: web
  serviceAccount: web          # <- deprecated Alias; Go-Feld heißt DeprecatedServiceAccount.
                               #    Steht bis heute im v1-PodSpec und verschwindet nie.
  volumes:
    - name: legacy
      # ... gitRepo, glusterfs, rbd, flexVolume, cinder — In-Tree-Typen, deren
      #     Implementierung längst entfernt ist. Die Felder bleiben.
```

Verschärft wird das durch die Reihenfolge der Entscheidungen: Der Wechsel von
[Embedded Entity](quality/EmbeddedEntity.md) zu
[Linked Information Holder](quality/LinkedInformationHolder.md) — bei MAP eine normale
Verbesserungsmaßnahme, im MAP-Tutorial sogar explizit als Schritt 4 vorgesehen — ist bei einer
GA-Ressource eine brechende Schemaänderung. Sie muss vor der Graduierung fallen, nicht nach den
ersten Performance-Beschwerden.

### 5. Eine überraschende Verletzung: Discovery ist nicht RBAC-gefiltert

MAP verlangt für eine [Community API](foundation/CommunityAPI.md), die Beschreibung nur mit der
Zielgruppe zu teilen. Die ClusterRole `system:discovery` ist an die Gruppe `system:authenticated`
gebunden:

```yaml
kind: ClusterRole
metadata: { name: system:discovery }
rules:
  - verbs: ["get"]
    nonResourceURLs:
      - /openapi
      - /openapi/*        # <- vollständiges OpenAPI-Schema aller Gruppen
      - /api
      - /api/*
      - /apis
      - /apis/*           # <- inklusive jeder CRD im Cluster
      # ... /livez, /readyz, /healthz, /version
---
kind: ClusterRoleBinding
metadata: { name: system:discovery }
roleRef: { kind: ClusterRole, name: system:discovery }
subjects:
  - kind: Group
    name: system:authenticated     # <- jeder, der sich anmelden kann — ohne jedes Recht
```

Wer sich anmelden kann, sieht das vollständige Schema aller Gruppen und Ressourcen — inklusive aller
CRDs, auf die er keinerlei Zugriff hat. In einem Multi-Tenant-Cluster erfährt Mandant A damit genau,
welche Operatoren Mandant B betreibt. Das ist bewusst so (Discovery und Versionsaushandlung müssen
vor der Autorisierung funktionieren), aber es ist eine Informationspreisgabe, die MAP explizit
ausschließt.

## Was das für den API-Entwurf bedeutet

Die Cheat-Sheet-Zeilen von MAP lassen sich für einen KRM-Entwurf fast vollständig vorab abhaken —
[die Tabelle steht im Cheat Sheet](meta/cheatsheet.md#bezug-zu-kubernetes--krm). Übrig bleibt echte
Designarbeit an genau zwei Stellen:

1. **Ressourcenschnitt.** Was ist ein `Kind`? Die Atomaritätsgarantie je Objekt ist exakt die
   Aggregate-Konsistenzgrenze aus DDD — der `Published Language`-Begriff beschreibt eine
   versionierte API-Gruppe präziser als alles in MAP.
2. **Schema.** Einbetten oder referenzieren, Listen-Merge-Semantik (`x-kubernetes-list-type`),
   Optionalität und Defaulting, und die Frage, ob sich jede geplante Erweiterung verlustfrei auf
   die heutige Version zurückkonvertieren lässt.

Alles andere hat das Modell schon entschieden.

## Weiterlesen

- [Überblick über MAP](meta/overview.md) · [Cheat Sheet](meta/cheatsheet.md) · [Terminologie](meta/terms.md)
- Die Kategorien: [Foundation](meta/category-foundation.md) ·
  [Responsibility](meta/category-responsibility.md) · [Structure](meta/category-structure.md) ·
  [Quality](meta/category-quality.md) · [Evolution](meta/category-evolution.md)
- Die [Papers](README.md#papers), insbesondere
  [das Entscheidungsmodell](papers/2018-icsoc-decision-guidance.md) und
  [Interface Evolution Patterns](papers/2019-interface-evolution-patterns.md)

---
[← Index](README.md)
