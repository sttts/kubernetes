---
title: Atomic Parameter
kategorie: Structure
unterkategorie: Representation Elements
quelle: https://microservice-api-patterns.org/patterns/structure/representationElements/AtomicParameter
---

# Atomic Parameter

*a.k.a.* Single Scalar Representation, Dot

**Kurzform:** Ein einzelner, unstrukturierter Wert — Zahl, String, Boolean oder Binärblock — wird als
eigenständiges Nachrichten- oder Parameterelement definiert und mit Name, Typ, Kardinalität und
Optionalität in der API-Beschreibung dokumentiert.

## Kontext

Ein API-Provider bietet Operations an einem Endpoint an. Client und Provider müssen sich auf die
Struktur jeder Request- und Response-Nachricht einigen. *Atomic Parameter* ist die einfachste der vier
Optionen der Kategorie *Structural Representation*; die anderen drei sind
[Atomic Parameter List](AtomicParameterList.md), [Parameter Tree](ParameterTree.md) und
[Parameter Forest](ParameterForest.md).

## Problem

Wie lassen sich einfache, unstrukturierte Daten zwischen Client und Provider austauschen, ohne die
Nachricht unnötig aufzublähen?

## Forces

Die vier Repräsentationspatterns teilen sich denselben Kräftesatz — die Datenstruktur der Nachrichten
ist wesentlicher Teil des API-Kontrakts:

- Struktur des Domänenmodells und des Systemverhaltens und deren Wirkung auf Verständlichkeit
  (Einfachheit, Komplexität, Nachvollziehbarkeit)
- Zusätzlich zu übertragende Daten (Security-Informationen, Metadaten)
- Performance (Latenz, Nachrichtenverarbeitung) und Ressourcenverbrauch (Bandbreite, Speicher, CPU)
- Lose Kopplung und Interoperabilität
- Developer Convenience und Developer Experience
- Sicherheit und Datenschutz

## Lösung

Genau ein Parameter bzw. ein Body-Element definieren und ihm einen Basistyp aus dem Typsystem des
gewählten Austauschformats geben. Einen Namen bekommt der Parameter nur dann, wenn die
Empfängerseite ihn zur Verarbeitung braucht — bei einem alleinstehenden Body-Wert ist der Name
verzichtbar, bei einem Query-Parameter nicht. Name, Typ, Kardinalität und Optionalität gehören in die
API Description.

Transportvarianten in RESTful HTTP: als Pfadsegment (typischerweise eine ID), als Query-Parameter
(die Empfehlung der Quelle für alles, was in der URL steht und *keine* ID ist), oder als Body-Element.

## Beispiel

Pfadparameter und Body-Variante desselben Werts:

```
curl http://localhost:8080/claims/a1e00494-e982-45f3-aab1-78a10ae3e3bd
```

```json
{ "id": "a1e00494-e982-45f3-aab1-78a10ae3e3bd" }
```

## Konsequenzen

**Vorteile:**

- Minimale Nachrichtengröße, triviales Parsing, keine Schema-Evolution nötig.
- Maximal lose Kopplung: Client und Provider teilen nur einen Skalar, kein gemeinsames Objektschema.
- Gut cachebar und gut in URLs abbildbar (Pfad, Query-String).

**Nachteile / Kosten:**

- Keine Ausdrucksstärke: Zusammengehörigkeit mehrerer Werte ist nicht darstellbar.
- Semantik hängt vollständig an Name und Dokumentation; ein nackter String ist selbstbeschreibungsfrei.
- Evolution ist teuer — sobald ein zweiter Wert nötig wird, ändert sich die Signatur bzw. es entsteht
  eine [Atomic Parameter List](AtomicParameterList.md).
- Validierung (Wertebereiche, Formate) muss explizit ergänzt werden, sonst ist der Skalar ein
  Einfallstor für Injection und Fehlbedienung.

## Bekannte Verwendungen

- Praktisch jede nachrichtenbasierte Remote-API, z.B. Statusabfrage eines Langläufers per Prozess-ID.
- AWS S3 `createBucket` mit einem einzelnen `bucketName`-String.
- URI-Query-Strings in RESTful HTTP, etwa `event-id` in der Facebook Graph API.
- Skalare Wertetypen in Protocol Buffers, Primitive Types in Apache Avro.
- Das *parameter object* in Swagger/OpenAPI beschreibt laut Quelle "a single operation parameter".

## Verwandte Patterns

- [Atomic Parameter List](AtomicParameterList.md) — mehrere Atomic Parameters als zusammengehörige Gruppe.
- [Parameter Tree](ParameterTree.md) — *Atomic Parameter* sind die Blätter des Baums.
- [Parameter Forest](ParameterForest.md) — kann Atomic Parameters als Bausteine enthalten.
- [Id Element](IdElement.md) — der häufigste Anwendungsfall eines alleinstehenden Skalars.
- [Metadata Element](MetadataElement.md) und [Data Element](DataElement.md) — Element-Stereotypen, die
  als Atomic Parameter realisiert sein können.
- [Link Element](LinkElement.md) — eine URI als einzelner String ist ein Atomic Parameter.
- [API Key](APIKey.md) — klassisches skalares Credential im Header oder Query-String.
- [Version Identifier](../evolution/VersionIdentifier.md) — skalarer Versionsstring.
- [Pagination](../quality/Pagination.md) — Cursor/Offset/Limit sind einzelne Skalare.

Aus Hohpe/Woolf (2003) verwendbar in *Command Message*, *Document Message* und *Event Message*, sofern
sich der Inhalt als Skalar darstellen lässt.

## Bezug zu Kubernetes / KRM

Im KRM sind *Atomic Parameters* fast ausschließlich außerhalb des Payloads zu finden: in Pfad und
Query-String. Der Pfadparameter `{name}` in `/apis/apps/v1/namespaces/{namespace}/deployments/{name}`
ist die kanonische Instanz — ein DNS-Subdomain-String, kein UUID wie im MAP-Beispiel. Die
Query-Parameter der List-/Watch-Operationen sind ebenfalls Skalare: `labelSelector`, `fieldSelector`,
`resourceVersion`, `resourceVersionMatch`, `limit`, `continue`, `watch`, `allowWatchBookmarks`,
`timeoutSeconds`, `sendInitialEvents`. Bemerkenswert ist, dass Kubernetes diese Skalare *nicht* als
lose Parameterliste modelliert, sondern als Go-Struct `metav1.ListOptions` mit `TypeMeta` — der
Query-String wird per generierter `Convert_url_Values_To_v1_ListOptions` in einen
[Parameter Tree](ParameterTree.md) konvertiert (`staging/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go`,
`zz_generated.conversion.go`). Das ist Punkt (b): dieselbe Wirkung, aber ein einheitlicher,
versionierter Objekttyp statt einer ad-hoc-Signatur.

```http
# Pfadparameter {name} — ein einzelner Skalar, wie im MAP-Beispiel:
GET /apis/apps/v1/namespaces/default/deployments/web

# Query-String einer Watch-Anfrage (zur Lesbarkeit umgebrochen).
# Jeder Parameter ist für sich ein Atomic Parameter:
GET /apis/apps/v1/namespaces/default/deployments
      ?labelSelector=app%3Dweb          # <- string
      &limit=500                        # <- int64
      &resourceVersion=0
      &resourceVersionMatch=NotOlderThan
      &watch=true                       # <- bool
      &allowWatchBookmarks=true
      &timeoutSeconds=300               # <- *int64, fehlt wenn nicht gesetzt
Accept: application/json
```

Weitere Skalare an Schreiboperationen: `dryRun` (technisch `[]string`, praktisch nur `All`),
`fieldManager`, `fieldValidation`, `force` in `metav1.PatchOptions`/`CreateOptions`/`UpdateOptions`,
sowie `gracePeriodSeconds` und `propagationPolicy` in `metav1.DeleteOptions`.

Innerhalb der Payloads treten Atomic Parameters nur als Blätter auf. Zwei apimachinery-Eigenheiten
sind hier relevant: `resource.Quantity` und `intstr.IntOrString` implementieren eigene
`MarshalJSON`/`UnmarshalJSON` und erscheinen dadurch nach außen als einfacher Skalar (`"100Mi"`,
`"25%"`, `8080`), obwohl sie intern strukturiert sind — ein bewusst domänenspezifischer Atomic
Parameter statt eines generischen Trees.

```yaml
apiVersion: v1
kind: Service
spec:
  ports:
    - port: 80
      targetPort: 8080          # <- intstr.IntOrString, hier die Int-Variante
```

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  strategy:
    rollingUpdate:
      maxUnavailable: "25%"     # <- derselbe Go-Typ, hier die String-Variante
      maxSurge: 1               # <- und daneben wieder Int
  template:
    spec:
      containers:
        - name: app
          # ... image, ports
          resources:
            limits:
              memory: 100Mi     # <- resource.Quantity: intern strukturiert, außen ein Skalar
              cpu: 500m
```

`omitempty` plus Pointer-Typen (`*int64`, `*bool`) regeln, ob ein nicht gesetzter Skalar aus der
Serialisierung verschwindet oder explizit als Wert auftaucht; ein nackter `int64` ohne Pointer kann
„nicht gesetzt" und „0" nicht unterscheiden.

```yaml
# Beide Felder sollen 0 sein — aber nur eines übersteht das Marshalling:
spec:
  replicas: 0                   # <- *int32 + omitempty: explizite 0 bleibt erhalten
  minReadySeconds: 0            # <- int32 + omitempty: fällt beim Serialisieren weg
  # Fehlt replicas dagegen ganz, heißt das nicht „0", sondern „nicht gesetzt" —
  # dann setzt das Defaulting 1 ein.
```

Diese Unterscheidung ist im MAP-Pattern als Optionalität nur benannt, in Kubernetes ist sie eine
Konventionsfrage mit direkten Auswirkungen auf Defaulting, Patch-Semantik und Server-Side Apply.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/representationElements/AtomicParameter)
