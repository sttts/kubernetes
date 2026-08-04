---
title: Parameter Forest
kategorie: Structure
unterkategorie: Representation Elements
quelle: https://microservice-api-patterns.org/patterns/structure/representationElements/ParameterForest
---

# Parameter Forest

*a.k.a.* Parameter Comb, Hybrid Parameter List

**Kurzform:** Zwei oder mehr gleichrangige [Parameter Trees](ParameterTree.md) bilden gemeinsam die
Nutzlast einer Operation — ohne dass ihnen eine gemeinsame, künstliche Wurzel übergestülpt wird.

## Kontext

Ein API-Provider bietet Operations an einem Endpoint an. Die auszutauschenden Daten zerfallen in
mehrere komplexe, aber inhaltlich getrennte Blöcke — etwa Nutzdaten und Steuer- bzw. Metadaten.

## Problem

Wie können mehrere Parameter Trees als Request- oder Response-Payload einer API-Operation exponiert
werden?

## Forces

Wie bei den drei übrigen Repräsentationspatterns:

- Struktur des Domänenmodells und des Systemverhaltens und deren Wirkung auf Verständlichkeit
- Zusätzlich zu übertragende Daten (Security-Informationen, Metadaten)
- Performance (Latenz, Verarbeitung) und Ressourcenverbrauch (Bandbreite, Speicher, CPU)
- Lose Kopplung und Interoperabilität
- Developer Convenience und Developer Experience
- Sicherheit und Datenschutz

## Lösung

Einen *Parameter Forest* aus zwei oder mehr Parameter Trees bilden. Die Mitglieder des Forests werden
über Position oder Name adressiert. Die Metapher der Quelle ist der Kamm: jeder Baum ist ein Zahn.

## Beispiel

Die Quelle illustriert alle vier Repräsentationspatterns in einer Java-Interface-Signatur mit ihren
Aliasnamen (`Dot`, `Dotted Line`, `Bar`, `Comb`):

```java
@WebService
public interface IRPService {
  boolean dotInDotOut(int singleScalarParameter);
  int dottedLineInDotOut(String scalarParameter1, String scalarParameter2);
  ResponseDTO barInBarOut(RequestDTO singleComplexParameter);
  ResponseDTO combInBarOut(RequestDTO complexParameter1,
                           AnotherRequestDTO complexParameter2);
}
```

`combInBarOut` hat zwei Parameter, die je einen „Zahn des Kamms" bilden und jeweils als DTO
strukturiert sind.

## Konsequenzen

**Vorteile:**

- Die Trennung inhaltlich unabhängiger Blöcke bleibt im Kontrakt sichtbar; keine erfundene Wurzel.
- Nutzdaten und Steuerinformationen (Cursor, Fehler, Metadaten) können unabhängig evolvieren.
- Erlaubt es, Sicherheits- oder Kontextinformationen anzuhängen, ohne die Domänenstruktur zu verbiegen.

**Nachteile / Kosten:**

- Komplexeste der vier Optionen: mehr zu spezifizieren, zu validieren und zu dokumentieren.
- Positionsbasierte Adressierung der Bäume ist brüchig, wenn sich die Signatur ändert.
- Viele Transportformate haben nur *einen* Body; der Forest muss dann doch in einen
  [Parameter Tree](ParameterTree.md) eingebettet werden, wodurch die Abgrenzung verschwimmt.
- Clients müssen mehrere Strukturen gleichzeitig verstehen — höhere kognitive Last.

## Bekannte Verwendungen

- Core-Banking-Integrationslösung bei Brandner et al. (2004) mit domänenspezifischen XSD-Typen.
- Flickr App Garden, z.B. `flickr.collections.getInfo`.
- Twitter REST API: ein Baum mit dem Objekt-Array, ein zweiter mit Kontrollinformationen und
  Metadaten wie Cursors und Page Tokens (siehe [Pagination](../quality/Pagination.md)).
- Protocol Buffers: eine Nachricht, die eine mit `repeated` markierte Nachricht enthält.
- JSON:API *Message Responses* mit `data`, `errors`, `meta` (verpflichtend) sowie `jsonapi`, `links`,
  `included` (optional) — jedes davon ein Parameter Tree.

## Verwandte Patterns

- [Parameter Tree](ParameterTree.md) — Baustein und zugleich Alternative, wenn eine gemeinsame Wurzel
  fachlich sinnvoll ist.
- [Atomic Parameter](AtomicParameter.md), [Atomic Parameter List](AtomicParameterList.md) — einfachere
  Alternativen bzw. Bausteine der einzelnen Bäume.
- [Pagination](../quality/Pagination.md) — Trennung von Seiteninhalt und Cursor-/Metadatenblock ist der
  Standardfall für einen Forest.
- [Error Report](ErrorReport.md) — Fehlerblock neben dem Datenblock.
- [Context Representation](ContextRepresentation.md) — Kontext-/Security-Baum neben den Nutzdaten.
- [Metadata Element](MetadataElement.md) — die Metadaten-„Zähne" des Kamms.
- [Embedded Entity](../quality/EmbeddedEntity.md) — beeinflusst, wie voll die einzelnen Bäume werden.

Verwendbar in *Command Message*, *Document Message* und *Event Message* (Hohpe/Woolf 2003), etwa wenn
zu einem Dokument noch Security- oder sonstige Metadaten gehören. *Content Filter* und
*Content Enricher* (ebd.) erzeugen bzw. reduzieren solche Strukturen.

## Bezug zu Kubernetes / KRM

KRM erfüllt dieses Pattern bewusst *nicht* in seiner reinen Form — und das ist eine der markantesten
Abweichungen. Die Kubernetes-API-Konventionen verlangen, dass jeder Request- und Response-Body ein
einzelnes Objekt mit `apiVersion` und `kind` ist. Ein Top-Level-Array oder ein Payload aus mehreren
gleichrangigen Wurzeln ist ausgeschlossen, weil sonst weder Typ-Discovery noch generische Clients noch
additive Erweiterbarkeit funktionieren. Wo MAP einen Forest vorsähe, verpackt KRM ihn in einen
[Parameter Tree](ParameterTree.md).

Die Listenobjekte sind das beste Beispiel. In MAP-Notation läge ein Forest nahe — zwei Zähne, die
inhaltlich nichts miteinander zu tun haben:

```
// hypothetisch, so gibt es das im KRM nicht:
operation listPods
  delivering payload
    > "listControl": {"resourceVersion":D, "continue":ID, "remainingItemCount":D},
    > "pods": Pod*
```

```yaml
# Tatsächlich: dieselben zwei Blöcke, aber unter eine gemeinsame Wurzel gezwungen.
apiVersion: v1
kind: PodList                  # <- die Wurzel, die den Forest verhindert
metadata:                      # <- Zahn 1: Kontroll-/Aggregat-Metadaten (metav1.ListMeta)
  resourceVersion: "12345"
  continue: eyJ2IjoibWV0YS5rOHMuaW8vdjEi...
  remainingItemCount: 1200
items:                         # <- Zahn 2: die Nutzdaten
  - metadata:
      name: web-1
    # ... spec, status
  - metadata:
      name: web-2
    # ... spec, status
```

Der Nutzdatenblock (`items`) und der Kontroll-/Metadatenblock (`metadata`) sind inhaltlich genau die
zwei „Zähne", die die Quelle am Twitter-Beispiel beschreibt — nur eben unter eine
gemeinsame Wurzel gezwungen. Der Forest ist da, aber als Teilbaumpaar. Dasselbe Muster zeigt
`metav1.Table` für die `kubectl`-Ausgabe: `columnDefinitions` (Schema) und `rows` (Daten) sind zwei
unabhängige Bäume in einem Objekt.

Auch die Watch-Semantik ist ein Fall dieser Umformung: Statt eines Payloads aus mehreren Bäumen liefert
`GET ...?watch=true` einen Stream von einzelnen `metav1.WatchEvent`-Objekten (`type` plus `object` als
`runtime.RawExtension`, siehe `apimachinery/pkg/apis/meta/v1/watch.go`).

```http
GET /api/v1/namespaces/default/pods?watch=true&resourceVersion=12345&allowWatchBookmarks=true

HTTP/1.1 200 OK
Content-Type: application/json
Transfer-Encoding: chunked

{"type":"ADDED","object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"web-3", ...}}}
{"type":"MODIFIED","object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"web-1", ...}}}
{"type":"BOOKMARK","object":{"kind":"Pod","apiVersion":"v1","metadata":{"resourceVersion":"12403"}}}
{"type":"DELETED","object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"web-2", ...}}}
```

Die Sequenz ist kein Array-Payload, sondern ein Chunked-Stream gleichförmiger Trees — eine Konsequenz der
Level-Triggered-Reconciliation, bei der ein Client den Zustand fortlaufend nachführt statt einmalig
abzuholen.

Und schließlich gilt: sobald eine Struktur unter einer Wurzel liegt, greift die
Server-Side-Apply-Merge-Semantik über `x-kubernetes-list-type` (`atomic`/`set`/`map`, mit
`x-kubernetes-list-map-keys`) und `x-kubernetes-map-type` (`atomic`/`granular`). Ein echter Forest aus
unabhängigen Top-Level-Bäumen hätte kein gemeinsames Feldpfad-Universum und damit keinen definierten
Besitzbegriff pro Feld. Die Aufgabe des Forest-Patterns ist im KRM also der Preis für
Multi-Writer-Deklarativität.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/representationElements/ParameterForest)
