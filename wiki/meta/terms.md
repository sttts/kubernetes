---
title: Terminologie / Glossar
kategorie: Meta
quelle: https://microservice-api-patterns.org/terms
---

# Terminologie / Glossar

## Worum es geht

MAP definiert seine Begriffe nicht auf einer eigenen Glossarseite, sondern in einem
*API-Domänenmodell*, das protokollunabhängig von HTTP, SOAP oder Messaging abstrahiert. Beschrieben
ist es in Kapitel 1 des [Buchs](book-and-resources.md) („A Domain Model for Remote APIs") und in
Auszügen auf der [Introduction](introduction.md)-Seite; die Pattern-Filter
[nach Scope](navigation-byscope.md) verwenden dieselben Ebenen als Gliederung.

> *Hinweis zur Quelle:* Die Seite `/terms` der MAP-Website enthält die Nutzungsbedingungen
> (Terms and Conditions of Use), kein Glossar. Die Definitionen unten sind aus
> [Introduction](introduction.md), [Primer](primer.md), den Kategorieseiten und dem
> Buch-Inhaltsverzeichnis zusammengetragen; die Kubernetes-Entsprechungen stammen aus den
> Kubernetes-API-Konventionen und der API-Machinery.

Aus den Nutzungsbedingungen ist praktisch relevant: Die Pattern-Icons stehen unter CC BY 4.0; für
das Anwenden oder Zitieren der Patterns wird keine formale Zitation verlangt, eine Nennung der
Website oder des Buchs ist aber erwünscht.

## Inhalt

### Kernbegriffe des API-Domänenmodells

| MAP-Begriff | Bedeutung | Kubernetes / KRM |
|---|---|---|
| **API** | Menge dokumentierter Netzwerk-Endpunkte, über die Komponenten einander Dienste anbieten. | Der `kube-apiserver` als Ganzes; feiner: eine API-Gruppe (`apps`, `batch`, `rbac.authorization.k8s.io`). |
| **API Provider** | Die Seite, die Endpunkte implementiert, betreibt und ihren Vertrag verantwortet. | `kube-apiserver` selbst, ein Aggregated API Server (`APIService`), oder — bei CRDs — der zuständige Controller/Operator. |
| **API Client** | Die Seite, die Endpunkte aufruft. | `kubectl`, Controller, Kubelet, Scheduler, Operatoren; technisch `client-go` mit `rest.Config`. |
| **API Endpoint** | Adressierbarer Zugangspunkt mit einer Rolle und einer Menge von Operationen. | Group-Version-Resource, z. B. `apps/v1, Resource=deployments` unter `/apis/apps/v1/namespaces/{ns}/deployments`. |
| **Operation** | Aufrufbare Funktion eines Endpunkts, definiert durch Signatur und Zustandsverhalten. | *Verb*: `get`, `list`, `watch`, `create`, `update`, `patch`, `delete`, `deletecollection`. Die Menge ist fix, nicht pro Ressource frei definierbar. |
| **Message** | Die tatsächlich übertragene Nachricht (Request oder Response). | HTTP-Request/-Response-Body; serialisiert als JSON, YAML oder Protobuf, ausgehandelt über `Content-Type`/`Accept`. |
| **Representation** | Die Struktur einer Nachricht: wie Daten für die Übertragung geformt sind. | Das Objekt in einer konkreten `apiVersion`; dieselbe interne Entität hat je Group-Version eine andere Repräsentation. |
| **Representation Element** | Ein einzelnes Element innerhalb einer Repräsentation. | Ein Feld im Objektschema, adressierbar über einen Feldpfad (`spec.template.spec.containers[0].image`). |

### Element-Stereotypen

| MAP-Begriff | Bedeutung | Kubernetes / KRM |
|---|---|---|
| **[Data Element](../structure/DataElement.md)** | Fachliches Nutzdatum, gekapselt gegenüber der Providerimplementierung. | Felder unter `spec` (Wunschzustand) und `status` (beobachteter Zustand). Die Trennung kennt MAP nicht. |
| **[Metadata Element](../structure/MetadataElement.md)** | Zusatzinformation zur Interpretation der übrigen Elemente. | Alles unter `metadata` (`metav1.ObjectMeta`): `labels`, `annotations`, `resourceVersion`, `generation`, `creationTimestamp`, `deletionTimestamp`, `finalizers`, `managedFields`. |
| **[Id Element](../structure/IdElement.md)** | Bezeichner, der Elemente unterscheidbar macht; global oder API-lokal eindeutig. | `metadata.name` (eindeutig je Namespace und Ressource), `metadata.namespace`, `metadata.uid` (global eindeutig, serverseitig vergeben, nicht wiederverwendet). |
| **[Link Element](../structure/LinkElement.md)** | Netzwerkauflösbarer Verweis auf einen anderen Endpunkt. | `ObjectReference`, `metadata.ownerReferences`, `LocalObjectReference` — strukturierte Verweise statt URLs; die Adresse wird aus GVR + Namespace + Name rekonstruiert. |

### Weitere Begriffe aus MAP und seinem Umfeld

| Begriff | Bedeutung | Kubernetes / KRM |
|---|---|---|
| **[API Description](../foundation/APIDescription.md)** (a.k.a. Service Contract) | Das Artefakt, das Struktur, Fehler, Ablauf, Vor-/Nachbedingungen und Policies festhält. | Discovery unter `/api` und `/apis` (aggregiert über `apidiscovery.k8s.io/v2`), OpenAPI v2 (`/openapi/v2`) und v3 (`/openapi/v3`), `kubectl explain`, plus `api-conventions.md`. |
| **Published Language** (DDD) | Ein gemeinsam vereinbartes Austauschvokabular zwischen Kontexten, das keine Seite intern bindet. | Das versionierte API-Schema selbst. `v1` ist die veröffentlichte Sprache; die internen Go-Typen (`pkg/apis/.../types.go`) sind es ausdrücklich *nicht* und werden nie serialisiert. |
| **Bounded Context** (DDD) | Abgegrenzter Modellbereich mit eigener Ubiquitous Language. | Am ehesten eine API-Gruppe samt zugehöriger SIG und Controller. |
| **Aggregate** (DDD) | Konsistenzgrenze über mehrere Entitäten mit einer Wurzel. | Ein Objekt (`Kind`) — Konsistenz wird atomar je Objekt garantiert, nie objektübergreifend. |
| **Kind** | *(kein MAP-Begriff)* | Der Typname im Schema (`Deployment`). Das Tripel Group/Version/Kind (GVK) benennt den Typ, Group/Version/Resource (GVR) den Endpunkt; die Abbildung leistet der RESTMapper. |
| **Subresource** | *(kein MAP-Begriff)* | Zweiter Endpunkt auf demselben Objekt mit eigener Autorisierung und eigenem Schreibpfad: `/status`, `/scale`, `/exec`, `/eviction`, `/token`. Erlaubt Rollentrennung ohne neue Ressource. |
| **Watch** | *(kein MAP-Begriff)* | Langlebiger Stream von `ADDED`/`MODIFIED`/`DELETED`/`BOOKMARK`-Ereignissen ab einer `resourceVersion`. Ersetzt in KRM das, was MAP über [Conditional Request](../quality/ConditionalRequest.md) und Polling löst. |
| **Conversation** | Mehrere zusammengehörige Nachrichtenaustausche zwischen Client und Provider. | Die Kombination aus initialem `list` (mit `resourceVersion`) und anschließendem `watch` — die kanonische KRM-Konversation, implementiert im Reflector/Informer. |
| **Data on the outside** | Ausgetauschte Daten, im Gegensatz zu „data on the inside" (gespeicherte Daten). | Genau die Unterscheidung zwischen externen versionierten Typen (`staging/src/k8s.io/api/...`) und internen Typen; die Konvertierung dazwischen ist generierter Code. |
| **Forces** | Die widerstreitenden Kräfte, die ein Problem schwer lösbar machen. | Im KEP-Prozess unter „Motivation", „Goals/Non-Goals" und „Risks and Mitigations" abgebildet. |
| **Known Uses** | Belegte reale Anwendungen eines Patterns. | Im KEP-Format am ehesten „Alternatives" und die Graduierungskriterien. |

## Verwandte Wiki-Seiten

- [Überblick MAP](overview.md) — Einstieg und die fünf Kategorien
- [Introduction](introduction.md) — dort wird das Domänenmodell motiviert
- [Kategorie Structure](category-structure.md) — die Stereotypen im Detail
- [Patterns nach Scope](navigation-byscope.md) — dieselben Ebenen als Navigationsraster
- [KRM vs. MAP](../KRM-vs-MAP.md) — die durchgehende Gegenüberstellung

## Bezug zu Kubernetes / KRM

Die begriffliche Hauptdifferenz ist der Ort der Gestaltungsfreiheit. In MAP entwirft man je
Endpunkt eine eigene Menge von Operationen mit eigenen Signaturen; Endpunkt, Operation und
Nachricht sind drei unabhängig gestaltbare Ebenen. In KRM ist die Operationsebene vollständig
standardisiert — die Verbmenge ist über alle Ressourcen identisch, ihre Semantik in
`api-conventions.md` festgelegt. Übrig bleiben zwei Ebenen: welche Ressourcen es gibt und wie ihr
Schema aussieht.

Daraus folgt eine zweite Differenz: MAP kennt keine Entsprechung für das Paar `spec`/`status`.
Bei MAP ist eine Nachricht entweder Eingabe oder Ausgabe; bei KRM ist *dasselbe Objekt* beides,
und die Zuständigkeit ist feldweise geteilt (Client schreibt `spec`, Controller schreibt `status`
über die `/status`-Subresource). Server-Side Apply macht diese Aufteilung explizit und
mehrschreiberfähig, indem `metadata.managedFields` je Feld den verantwortlichen `fieldManager`
festhält — ein Konzept, für das es in MAP kein Gegenstück gibt, weil dort nie zwei Parteien
dieselbe Repräsentation teilweise besitzen.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/terms)
