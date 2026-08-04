---
title: Wish Template
kategorie: Quality
unterkategorie: Data Transfer Parsimony
quelle: https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/WishTemplate
---

# Wish Template

*a.k.a.* Wish Skeleton, Response Mock, Attribute/Entity/Representation Stencil

**Kurzform:** Der Client legt dem Request ein Skelett-Objekt bei, das die Struktur der gewünschten
Antwort spiegelt; Anwesenheit bzw. Belegung eines Knotens signalisiert Interesse, Abwesenheit
„don't care".

## Kontext

Wie bei [Wish List](WishList.md): Ein Endpoint bedient Clients mit sehr unterschiedlichem
Informationsbedarf. Der Unterschied liegt in der Datenform — die Antwort ist tief verschachtelt,
und der Client braucht bestimmte Teilbäume, nicht bestimmte Top-Level-Attribute.

## Problem

Wie teilt ein Client dem Provider mit, an welchen **verschachtelten** Daten er interessiert ist,
und wie lassen sich solche Präferenzen flexibel und dynamisch ausdrücken?

## Forces

Identisch zu [Wish List](WishList.md):

- **Performance, Skalierbarkeit, Ressourcenverbrauch**
- **Individuelle Informationsbedarfe** der Clients
- **Lose Kopplung und Interoperabilität** — eine strukturspiegelnde Vorlage koppelt noch stärker
  an das Antwortschema als eine flache Namensliste.
- **Developer Convenience** — die Vorlage ist ausdrucksstark, aber aufwendig zu konstruieren.
- **Sicherheit und Datenschutz**
- **Test- und Wartungsaufwand** — der Raum möglicher Antwortformen ist nun ein Baum, kein Powerset
  einer flachen Menge.

## Lösung

Der Request bekommt einen oder mehrere zusätzliche Parameter, die die **hierarchische Struktur der
Antwortparameter nachbilden**. Diese Parameter werden optional gemacht oder auf `Boolean` getypt,
sodass ihr Vorhandensein bzw. ihr Wert angibt, ob der zugehörige Antwortteil geliefert werden soll.
Die Vorlage wird damit Teil des [Parameter Tree](../structure/ParameterTree.md) der Anfrage.

## Beispiel

Vertragsskizze aus der Quelle: Ein `mockObject` mit derselben Struktur wie das Ergebnis reist im
Request mit; optionale Teilstrukturen (`PersonalData?`, `Address?`) drücken die Wünsche aus.

```
data type CustomerEntity { PersonalData?, Address? }

endpoint type CustomerInformationHolderService
  exposes
    operation getCustomerAttributes
      expecting payload {
        "requestPayload": D,
        "mockObject": CustomerEntity   // gleiche Struktur wie das gewünschte Ergebnis
      }
      delivering payload CustomerEntity*
```

## Konsequenzen

**Vorteile:**

- Auch tiefe, heterogene Strukturen lassen sich präzise selektieren — inklusive „nur dieses eine
  Feld aus diesem einen Teilbaum".
- Die Vorlage ist selbstbeschreibend: Wer das Antwortschema kennt, kann die Anfrage bilden.
- Reduziert sowohl Volumen als auch Anzahl der Aufrufe und wirkt damit doppelt auf ein
  [Rate Limit](RateLimit.md).

**Nachteile / Kosten:**

- Höchste Kopplung aller Parsimony-Patterns: Der Client repliziert die Providerstruktur.
- Provider-seitig braucht es einen Interpreter für die Vorlage plus Validierung; Kosten- und
  Autorisierungsabschätzung pro Anfrage wird schwer.
- Die Vorlage selbst kostet Bandbreite und muss bei Schemaevolution mitwandern.

## Bekannte Verwendungen

Das „Dynamic Interface" bei Brandner et al. (2004) nutzt XML-Template-Objekte, deren Struktur die
Antwort spiegelt; leere Werte bedeuten Desinteresse. Die Google Calendar API mischt Wish Template
und [Wish List](WishList.md): Die `fields`-Syntax ist „loosely based" auf XPath und erlaubt damit
strukturierte Ausdrücke. Die Facebook Graph API nutzt das Pattern ebenfalls. **GraphQL** — ebenfalls
ursprünglich von Facebook — lässt sich als fortgeschrittene Realisierung lesen: Die Query *ist* eine
deklarative Beschreibung des Informationsbedarfs. Der Preis ist ein eigener GraphQL-Server als
zusätzliche Endpoint-Schicht über den eigentlichen Ressourcen.

## Verwandte Patterns

- [Wish List](WishList.md) — gleiches Problem, flache Aufzählung statt Mock-Objekt.
- [Pagination](Pagination.md) — reduziert die Länge statt der Tiefe.
- [Conditional Request](ConditionalRequest.md) — kombinierbar: Vorlage greift nur, wenn überhaupt gesendet wird.
- [Request Bundle](RequestBundle.md) — orthogonal; laut Quelle nur mit Bedacht zu kombinieren.
- [Parameter Tree](../structure/ParameterTree.md) — die Vorlage wird Teil des Request-Baums.
- [Parameter Forest](../structure/ParameterForest.md) — Zielstruktur bei mehreren Wurzeln.
- [Embedded Entity](EmbeddedEntity.md) / [Linked Information Holder](LinkedInformationHolder.md) — statische Designzeit-Antwort auf dieselbe Frage.
- [Rate Limit](RateLimit.md) — profitiert doppelt.

## Bezug zu Kubernetes / KRM

Von allen fünf Data-Transfer-Parsimony-Patterns ist *Wish Template* dasjenige, das KRM **am
wenigsten** umsetzt — es gibt schlicht **keinen Mechanismus, mit dem ein Client eine strukturierte
Antwortvorlage mitschickt**. Es existiert kein GraphQL-artiger Endpoint, keine Feldmaske und keinen
`mockObject`-Parameter. Ein `GET` liefert das ganze Objekt in der angeforderten Group/Version; die
einzigen Verkleinerungen sind die in [Wish List](WishList.md) beschriebenen, fest verdrahteten
Projektionen `as=PartialObjectMetadata` und `as=Table`.

Die strukturell nächste Verwandtschaft findet sich auf der **Schreibseite**, nicht auf der Leseseite:
Server-Side Apply. Ein Apply-Request enthält ein partielles Objekt, das genau die Felder trägt, für
die der Aufrufer Verantwortung übernimmt, identifiziert durch `fieldManager`:

```http
PATCH /apis/apps/v1/namespaces/default/deployments/web?fieldManager=my-operator&force=true
Content-Type: application/apply-patch+yaml    # <- types.ApplyYAMLPatchType

apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
spec:
  template:
    spec:
      containers:
        - name: app                            # <- Schlüssel, nicht Wunsch
          image: nginx:1.27                    # <- das einzige Feld, das dieser Manager besitzt
# alles Nicht-Genannte bleibt unangetastet und in fremdem Besitz
```

Der Server merged das und führt in `metadata.managedFields` pro Manager ein `fieldsV1`-Set. Dieses
Set ist formal exakt ein Wish Template: ein Baum, der die Objektstruktur mit `f:`-präfigierten
Schlüsseln spiegelt, Listenelemente über `k:{…}` adressiert und dessen Blätter leere Objekte sind
(`metav1.FieldsV1`, Serialisierung aus `sigs.k8s.io/structured-merge-diff`):

```yaml
metadata:
  managedFields:
    - manager: my-operator
      operation: Apply                         # <- oder Update bei gewöhnlichen Schreibzugriffen
      apiVersion: apps/v1
      fieldsType: FieldsV1
      fieldsV1:
        f:spec:
          f:template:
            f:spec:
              f:containers:
                k:{"name":"app"}:              # <- Listenelement per listMapKey adressiert
                  .: {}                        # <- das Element selbst
                  f:image: {}                  # <- Blatt = leeres Objekt
    # ... weitere Manager, etwa kube-controller-manager für f:status
```

Nur: Das beschreibt Eigentum an geschriebenen Feldern, nicht Interesse an gelesenen. Es gibt keine
Operation, die ein solches Set als Leseselektor akzeptiert.

Als grobkörnige, **provider-definierte** Vorlagen kann man die Subresources lesen: `/scale` liefert
für Deployments, ReplicaSets und StatefulSets ein `autoscaling/v1.Scale` mit `spec.replicas`,
`status.replicas` und `status.selector` — eine drastisch reduzierte Sicht auf ein großes Objekt,
die genau den Bedarf des HorizontalPodAutoscalers deckt. Ebenso ist `/status` eine strukturelle
Teilsicht mit eigener Autorisierung.

```http
GET /apis/apps/v1/namespaces/default/deployments/web/scale

{
  "kind": "Scale",
  "apiVersion": "autoscaling/v1",     // <- andere Group/Version als das Deployment selbst
  "metadata": { "name": "web", "namespace": "default", "resourceVersion": "8675309" },
  "spec":   { "replicas": 3 },
  "status": { "replicas": 3, "selector": "app=web" }
  // der Rest des Deployments — template, strategy, conditions — kommt gar nicht erst mit
}
```

Das ist der KRM-Weg: Wo MAP dem Client erlaubt, die Vorlage zur Laufzeit zu formulieren, definiert
Kubernetes sie zur Designzeit als eigenen, typisierten und RBAC-fähigen API-Pfad.

Der Grund für diese Verweigerung ist konsistent mit dem Rest des Modells: Reconciliation-Schleifen
arbeiten gegen einen lokalen Cache vollständiger Objekte, den ein Informer per `watch` inkrementell
aktuell hält. Wenn ein Objekt ohnehin nur einmal vollständig geladen und danach nur noch per Delta
fortgeschrieben wird, ist der Nutzen einer Antwortvorlage gering — und der Preis, dass Cache-Einträge
je nach anfragendem Controller unterschiedlich beschnitten wären, sehr hoch. Die Effizienz holt sich
KRM über [Conditional Request](ConditionalRequest.md)-Mechanismen und Watch, nicht über
Response Shaping.

---
[← Index](../README.md) · [Kategorie Quality](../meta/category-quality.md) · [Quelle](https://microservice-api-patterns.org/patterns/quality/dataTransferParsimony/WishTemplate)
