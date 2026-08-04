---
title: Id Element
kategorie: Structure
unterkategorie: Element Stereotypes
quelle: https://microservice-api-patterns.org/patterns/structure/elementStereotypes/IdElement
---

# Id Element

*a.k.a.* *Identifier Element*, *Unique Identifier*, *Id Representation*

**Kurzform:** Ein spezialisiertes [Data Element](DataElement.md), das API-Endpoints,
Operationen und Repräsentationselemente eindeutig unterscheidbar macht — zur Design- wie zur
Laufzeit, mit bewusst gewähltem Gültigkeitsbereich.

## Kontext

Ein Domänenmodell ist definiert und implementiert; der Fernzugriff darauf wird gerade gebaut
(HTTP-Ressourcen, Web-Service-Operationen, gRPC-Methoden). Das Modell besteht aus mehreren
verwandten Elementen mit unterschiedlichen Lebenszyklen und Semantiken. Die gewählte Zerlegung
in eigenständig deploybare Endpoints legt nahe, diese Entitäten auf mehrere API-Elemente
aufzuteilen. Konsumenten wollen den Beziehungen zwischen diesen Elementen folgen können.

## Problem

Wie können API-Elemente zur Design- und zur Laufzeit voneinander unterschieden werden? Und wie
werden bei Domain-Driven Design die Elemente der *Published Language* identifiziert?

## Forces

- **Aufwand vs. Stabilität** — global eindeutige, unveränderliche IDs sind teurer zu erzeugen
  und zu verwalten als lokale Zähler, halten aber Umzügen und Reorganisationen stand.
- **Lesbarkeit für Menschen und Maschinen** — sprechende Namen helfen bei Debugging und
  Support, UUIDs sind kollisionsfrei, aber opak.
- **Sicherheit (Vertraulichkeit)** — IDs leaken Information: fortlaufende Zähler verraten
  Bestandsgrößen, fachlich abgeleitete IDs verraten Fachdaten, ratbare IDs erlauben Enumeration.

## Lösung

Einen speziellen Typ von *Data Element* einführen — ein eindeutiges *Id Element* — und es
konsistent in [API Description](../foundation/APIDescription.md) und Implementierung verwenden.
Explizit entscheiden, ob die ID **global** eindeutig ist oder nur **im Kontext einer bestimmten
API** gilt.

Die Quellseite weist darauf hin, dass der Gültigkeitsbereich nicht nur räumlich, sondern auch
**zeitlich** begrenzt sein kann — etwa ein Cursor, der nur bis zum Ende einer Session gilt.

## Beispiel

Pagination-Cursor der Twitter REST API (Quellseite): der Server garantiert Eindeutigkeit von
`next_cursor` mindestens bis zum Ablauf der Session und hält die Zuordnung zur Cursor-Position
vor.

```json
{
    "data": [],
    "data_type": "campaign",
    "next_cursor": "c-3yvu1pzhd3i7",
    "total_count": 200
}
```

## Konsequenzen

**Vorteile:**

- Referenzierbarkeit ohne Übertragung des vollständigen Zustands.
- Stabile Identität über Umbenennungen und Attributänderungen hinweg (bei opaken IDs).
- Voraussetzung für Idempotenz, Deduplizierung und Korrelation.

**Nachteile / Kosten:**

- Erzeugung und Verwaltung eindeutiger IDs kostet Aufwand und kann zum Engpass werden.
- Opake IDs sind schwer zu debuggen, sprechende IDs sind schwer stabil zu halten.
- IDs im API-Vertrag sind faktisch für immer festgeschrieben.
- Anders als [Link Elements](LinkElement.md) sind *Id Elements* nicht selbst auflösbar; der
  Client muss wissen, wo er sie einlösen kann.

## Bekannte Verwendungen

Die Quellseite nennt: numerische IDs der Facebook Graph API
(`"id": "164202036981850_1277103665691676"`); menschenlesbare, generierte Namen bei Heroku
(`peaceful-reaches-47689`) und Docker (`practical_carson`); URNs in der LinkedIn-REST-API;
DDD-Entities mit Identität und Lebenszyklus (Cargo-Tracking-Beispiel); global eindeutige IDs
wie DOI, ORCID, Steuer- und Sozialversicherungsnummern sowie der UID-Dienst der Schweizer
Verwaltung; die EGRID/EREID/EGBPID-Kennungen von Terravis.

## Verwandte Patterns

- [Data Element](DataElement.md) — Oberbegriff.
- [Link Element](LinkElement.md) — engster Verwandter: macht Referenzen zusätzlich
  netzwerk-adressierbar und trägt Typinformation; *Id Elements* tun beides typischerweise nicht.
- [Atomic Parameter](AtomicParameter.md) — übliche syntaktische Form.
- [Metadata Element](MetadataElement.md) — IDs können selbst Metadaten sein oder von ihnen
  begleitet werden.
- [API Key](APIKey.md), [Version Identifier](../evolution/VersionIdentifier.md) — Sonderformen
  von Identifikatoren.
- [Master Data Holder](../responsibility/MasterDataHolder.md) — braucht wegen der Langlebigkeit
  besonders starke Identifikationsschemata;
  [Operational Data Holder](../responsibility/OperationalDataHolder.md) ebenfalls eindeutig
  identifiziert; die Werte eines
  [Reference Data Holder](../responsibility/ReferenceDataHolder.md) können selbst als IDs
  dienen (z. B. Postleitzahlen).
- [Link Lookup Resource](../responsibility/LinkLookupResource.md) — nimmt *Id Elements* entgegen
  und liefert [Link Elements](LinkElement.md).
- [Data Transfer Resource](../responsibility/DataTransferResource.md) — nutzt IDs zur Definition
  von Transfereinheiten und Speicherorten.
- [Pagination](../quality/Pagination.md) — Cursor als zeitlich begrenzte IDs.

Außerhalb von MAP: *Correlation Identifier*, *Return Address*, *Claim Check* und
*Format Identifier* (Hohpe/Woolf); RFC 4122 für UUIDs.

## Bezug zu Kubernetes / KRM

KRM erfüllt das Muster, aber nicht mit einem Identifikator, sondern mit vier klar getrennten,
in `metadata` uniform verorteten — ein ungewöhnlich sauber ausdifferenziertes Beispiel für
dieses Pattern:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web                                   # <- fachlich, menschenlesbar; eindeutig nur im
  namespace: default                          #    Tupel (group, resource, namespace, name)
  # alternativ serverseitig erzeugt:
  # generateName: web-      ->  name: web-x2klm
  uid: 8f2a1c0e-3b7d-4c11-9a55-2f0e6d7b1c34   # <- global eindeutig, unveränderlich, opak
  resourceVersion: "1250043"                  # <- identifiziert keine Entität, sondern einen
                                              #    Zustand; opak, nur zeitlich gültig
  generation: 7                               # <- Zähler der spec-Änderungen
  # ... labels, annotations, ownerReferences
spec:
  # ... replicas, selector, template
status:
  observedGeneration: 7                       # <- Controller meldet zurück, welche generation
  # ... readyReplicas, conditions             #    er verarbeitet hat
```

- **`metadata.name`** — der fachliche, menschenlesbare Identifikator. Sein Gültigkeitsbereich
  ist explizit gewählt: eindeutig **pro Namespace und pro Resource-Typ** (bzw. clusterweit bei
  cluster-scoped Ressourcen). Der volle Schlüssel ist damit das Tupel
  (group, resource, namespace, name). Die Syntax ist eingeschränkt (meist RFC-1123-Label bzw.
  -Subdomain), damit Namen in DNS-Namen und URL-Pfadsegmenten verwendbar sind. Ergänzt wird das
  durch **`metadata.generateName`**, mit dem der Server ein zufälliges Suffix anhängt — die
  serverseitige Erzeugung eindeutiger Namen, die MAP als Lösungsvariante beschreibt.
- **`metadata.uid`** — global eindeutig (UUID), vom Server bei der Erstellung vergeben,
  unveränderlich. Genau der von MAP geforderte „globally unique“-Fall. Er löst das Problem, das
  `name` nicht lösen kann: nach Löschen und Neuanlegen unter demselben Namen ist es ein anderes
  Objekt. Deshalb tragen `ownerReferences` und `ObjectReference` die `uid` mit — der Garbage
  Collector prüft sie, um verwaiste Referenzen auf wiederverwendete Namen zu erkennen.
  Derselbe Schutz lässt sich beim Löschen erzwingen:

  ```http
  DELETE /apis/apps/v1/namespaces/default/deployments/web HTTP/1.1
  Content-Type: application/json

  {"kind":"DeleteOptions","apiVersion":"v1",
   "preconditions":{"uid":"8f2a1c0e-3b7d-4c11-9a55-2f0e6d7b1c34"},
   "propagationPolicy":"Foreground"}
  ```

  Stimmt die `uid` nicht — weil inzwischen ein neues Objekt gleichen Namens existiert —, scheitert
  der Aufruf mit `409 Conflict`, statt das falsche Objekt zu löschen. `preconditions` akzeptiert
  daneben `resourceVersion` für dieselbe Prüfung auf Zustandsebene.
- **`metadata.resourceVersion`** — ein für Clients ausdrücklich **opaker** Wert, der die interne
  Version eines Objekts kennzeichnet. Er identifiziert keine Entität, sondern einen Zustand, und
  dient Optimistic Concurrency (ein `update` mit veraltetem Wert scheitert mit `409 Conflict`),
  Änderungserkennung und dem Wiederaufsetzen von `watch`. Das ist Ausdruck derselben
  MAP-Beobachtung, dass ein Identifikator auch zeitlich begrenzt gültig sein kann. Als
  Listen-Cursor tritt zusätzlich `metadata.continue` in `metav1.ListMeta` auf.
- **`metadata.generation`** — ein vom Server geführter Zähler der `spec`-Änderungen. Controller
  spiegeln ihn nach `status.observedGeneration` bzw. in `status.conditions[].observedGeneration`
  zurück und machen so maschinell entscheidbar, ob ein beobachteter Zustand zur aktuellen
  gewünschten Version gehört. Ein solches Konstrukt kennt MAP nicht; es ist eine direkte Folge
  des level-triggered Reconciliation-Modells.

Bemerkenswert ist, wie stark KRM die MAP-Force „Lesbarkeit für Menschen und Maschinen“ auflöst,
statt sie zu vertagen: `name` bedient Menschen und `kubectl`, `uid` bedient die Maschinerie,
und beide existieren nebeneinander, statt einen Kompromiss zu erzwingen.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/elementStereotypes/IdElement)
