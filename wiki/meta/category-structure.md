---
title: Kategorie Structure
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/structure
---

# Kategorie: Structure Patterns

## Worum es geht

Ein API-Vertrag enthält die Adresse eines Endpunkts, den Bezeichner einer Operation — und die
Struktur der Request- und Response-Nachrichten. Um Letzteres geht es hier. Die Leitfrage der
Kategorie lautet sinngemäß: Wie viele Parameter braucht eine API-Nachricht, wie werden sie
strukturiert, und wie werden sie gruppiert und mit Nutzungshinweisen annotiert?

Die Kategorie beantwortet das auf drei Ebenen, die sich kombinieren lassen:

1. **Repräsentationselemente** — die *Form* der Nutzlast (flach, Liste, Baum, Wald).
2. **Element-Stereotypen** — die *Bedeutungsrolle* eines Elements innerhalb dieser Form
   (Nutzdatum, Metadatum, Identifikator, Verweis).
3. **Spezialrepräsentationen** — wiederkehrende, fachlich vordefinierte Nutzlastteile.

Form und Stereotyp sind orthogonal: Ein [Metadata Element](../structure/MetadataElement.md) kann
ein [Atomic Parameter](../structure/AtomicParameter.md) oder ein ganzer
[Parameter Tree](../structure/ParameterTree.md) sein.

## Inhalt

### Repräsentationselemente — die Form der Nutzlast

| Pattern | Einzeiler |
|---|---|
| [Atomic Parameter](../structure/AtomicParameter.md) | Ein einzelnes, unstrukturiertes Datum (Zahl, String, Boolean, Binärblock) mit Name, Typ, Kardinalität und Optionalität. |
| [Atomic Parameter List](../structure/AtomicParameterList.md) | Mehrere zusammengehörige *Atomic Parameters* in einem kohäsiven Element, adressiert über Position oder Schlüssel. |
| [Parameter Tree](../structure/ParameterTree.md) | Hierarchische Struktur mit einem Wurzelknoten und beliebig geschachtelten Kindknoten samt Kardinalitäten. |
| [Parameter Forest](../structure/ParameterForest.md) | Zwei oder mehr *Parameter Trees* nebeneinander als Nutzlast einer Operation. |

Die vier bilden eine Komplexitätsleiter. Die Empfehlung des [Cheat Sheets](cheatsheet.md) ist
schlicht: einfache Daten → *Atomic Parameter*/*Atomic Parameter List*, komplexe Daten →
*Parameter Tree*/*Parameter Forest*.

### Element-Stereotypen — die Bedeutungsrolle

| Pattern | Einzeiler |
|---|---|
| [Data Element](../structure/DataElement.md) | Ein eigenes Vokabular für Fachdaten in Nachrichten, das die providerinternen Datendefinitionen kapselt statt sie durchzureichen. |
| [Metadata Element](../structure/MetadataElement.md) | Zusatzinformation, die dem Empfänger die korrekte Interpretation der übrigen Elemente erlaubt, ohne Annahmen fest zu verdrahten. |
| [Id Element](../structure/IdElement.md) | Ein Spezialfall des *Data Element*, der Endpunkte, Operationen und Repräsentationselemente unterscheidbar macht — global oder API-lokal eindeutig. |
| [Link Element](../structure/LinkElement.md) | Ein Spezialfall des *Id Element*: ein netzwerkauflösbarer Zeiger auf einen anderen Endpunkt bzw. eine andere Operation (HATEOAS). |

Die Stereotypen sind eine Spezialisierungskette: *Link Element* ist ein *Id Element* ist ein
*Data Element*. Der Sprung vom *Id Element* zum *Link Element* ist im Cheat Sheet die empfohlene
Antwort auf „unflexibler, statischer Operationsablauf".

### Spezialrepräsentationen

| Pattern | Einzeiler |
|---|---|
| [API Key](../structure/APIKey.md) | Ein je Client eindeutiges Token zur Identifikation (und begrenzt Authentifizierung) der Aufrufe. |
| [Error Report](../structure/ErrorReport.md) | Maschinenlesbare Fehlercodes plus menschenlesbare Beschreibungen in der Antwort, unabhängig von Protokoll-Statuscodes. |
| [Context Representation](../structure/ContextRepresentation.md) | Alle Kontext-Metadaten (Identität, QoS, Konversationszustand) gebündelt in *einem* Nutzlastelement statt in Protokoll-Headern. |

## Verwandte Wiki-Seiten

- [Kategorie Responsibility](category-responsibility.md) — legt fest, *welche* Nachrichten es gibt
- [Kategorie Quality](category-quality.md) — optimiert die hier festgelegten Strukturen
  ([Embedded Entity](../quality/EmbeddedEntity.md) vs. [Linked Information Holder](../quality/LinkedInformationHolder.md))
- [Terminologie](terms.md) — Message, Representation, Published Language
- [Cheat Sheet](cheatsheet.md)

## Bezug zu Kubernetes / KRM

KRM ist ein extremer Sonderfall dieser Kategorie: Es gibt genau *eine* Nachrichtenform für alle
Ressourcen. Jedes Objekt ist ein [Parameter Tree](../structure/ParameterTree.md) mit vier
vorgegebenen Wurzelknoten — `apiVersion`, `kind`, `metadata`, und (fast immer) `spec` und `status`.
Ein [Parameter Forest](../structure/ParameterForest.md) tritt praktisch nicht auf; das
Nächstliegende sind `List`-Typen mit `metadata` (`ListMeta`) und `items`, was aber wieder ein Baum
ist. Diese Uniformität ist die Voraussetzung für generische Werkzeuge (`kubectl`, Server-Side
Apply, `unstructured.Unstructured`) und wird in den API-Konventionen ausdrücklich verlangt.

Die Stereotypen sind in KRM nicht bloß konventionell, sondern strukturell durchgesetzt:
`metadata` (Typ `metav1.ObjectMeta`) ist der reservierte Ort für
[Metadata Elements](../structure/MetadataElement.md) — `labels`, `annotations`,
`resourceVersion`, `generation`, `creationTimestamp`, `finalizers`, `managedFields`. Die
[Id Elements](../structure/IdElement.md) sind `metadata.name` (namensraum-eindeutig),
`metadata.namespace` und `metadata.uid` (global eindeutig, von der API-Machinery vergeben, stabil
über Neuerzeugung hinweg). `spec` und `status` enthalten die
[Data Elements](../structure/DataElement.md); die Trennung Wunschzustand/Ist-Zustand ist eine
KRM-Erfindung, die MAP nicht kennt.

Beim [Link Element](../structure/LinkElement.md) geht KRM einen eigenen Weg: Statt URLs enthalten
Objekte strukturierte Verweise — `ObjectReference`, `metadata.ownerReferences`,
`LocalObjectReference`, oder feldspezifische Paare wie `spec.nodeName`. Die Auflösung erfolgt
nicht durch Dereferenzieren einer mitgelieferten Adresse, sondern durch Rekonstruktion des
Pfades aus Group-Version-Resource, Namespace und Name. HATEOAS im Wortsinn gibt es in KRM nicht;
die „Hypermedia" steckt in Discovery plus Konvention.

Der [Error Report](../structure/ErrorReport.md) ist in KRM als eigener Kind ausgeführt:
`metav1.Status` mit `status`, `reason` (z. B. `NotFound`, `Conflict`, `AlreadyExists`), `code`,
`message` und `details` — maschinenlesbar auswertbar über `k8s.io/apimachinery/pkg/api/errors`.
Das entspricht dem Pattern sehr genau, inklusive der Forderung, sich nicht allein auf
Protokoll-Statuscodes zu verlassen. Die [Context Representation](../structure/ContextRepresentation.md)
dagegen fehlt: Identitäts- und Kontextinformation reist in KRM in HTTP-Headern (Bearer-Token,
Impersonation-Header `Impersonate-User`) und in `context.Context` auf Serverseite, gerade *nicht*
in der Nutzlast. Ebenso ist der [API Key](../structure/APIKey.md) kein Nutzlastelement, sondern ein
ServiceAccount-Token oder Client-Zertifikat auf Transportebene.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure)
