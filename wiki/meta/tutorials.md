---
title: Tutorials 0–2
kategorie: Meta
quelle: https://microservice-api-patterns.org/patterns/tutorials
---

# Tutorials 0–2

## Worum es geht

MAP bietet drei aufeinander aufbauende, geführte Tutorials an, die Website-Navigation mit echten
Entwurfsaufgaben mischen. Jede Aufgabe hat eine ausklappbare Musterlösung. Verträge werden in
MDSL notiert (Microservice Domain Specific Language, Anhang C des Buchs); zusätzlich existiert zu
jedem Tutorial eine Selbststudium-Variante mit Wiederholungsfragen.

| Tutorial | Thema | Dauer |
|---|---|---|
| 0 | Orientierung und ein einzelnes Pattern ([Pagination](../quality/Pagination.md)) | 20–25 min |
| 1 | API-Qualität: Nachrichtengrößen, Zugriffskontrolle, SLA | 45 min |
| 2 | Alle Kategorien an einem realistischen Fall, in fünf Schritten | 60–90 min |

## Inhalt

### Tutorial 0 — Orientierung und ein einzelnes Pattern

Das „Hello World" des API-Designs mit MAP. Ausgangslage ist ein Kundensuch-Vertrag, dessen
`findCustomer`-Operation einen einfachen Suchfilter entgegennimmt und *alle* Treffer mit *allen*
Kundendaten zurückliefert — teuer auf beiden Seiten.

Der Ablauf: erst Startseite, [Cheat Sheet](cheatsheet.md), [Kategorien](overview.md) und die
Filter [nach Scope](navigation-byscope.md), [Phase](navigation-byphase.md) und
[Qualität](navigation-byquality.md) erkunden; dann über den Cheat-Sheet-Eintrag zu
Performance-Problemen [Pagination](../quality/Pagination.md) finden, Problem und Solution lesen,
und den Vertrag umbauen. Die Musterlösung führt eine `CustomerPage` mit `pageContent`, einem
`PageCursor` (`thisPage`/`previousPage`/`nextPage`) und optionalen `PageMetadata`
(`totalResultCount`, `pageSize`, `offsetNumber`) ein — und behält die alte Operation daneben.
Die genannten Varianten sind Offset-, Cursor-/Token- und Time-Based Pagination.

Nebenbei wird die MDSL-Notation eingeführt: `D` für ein [Data Element](../structure/DataElement.md),
`MD` für ein [Metadata Element](../structure/MetadataElement.md), `L` für ein
[Link Element](../structure/LinkElement.md), `ID` für ein [Id Element](../structure/IdElement.md).
Die Element-Stereotypen sind also im Typsystem der DSL verankert.

### Tutorial 1 — API-Qualität

Laufendes Beispiel ist der `policies`-Endpunkt der fiktiven Versicherung *Lakeside Mutual*, ein
[Master Data Holder](../responsibility/MasterDataHolder.md), dessen `PolicyDto` auf viele weitere
DTOs verweist. Drei Schritte:

1. **Data Transfer Parsimony.** [Wish List](../quality/WishList.md) lesen und anwenden: Für die
   nächste Hauptversion soll `customer` nicht mehr per Default geliefert werden, sondern über
   einen `expand`-Query-Parameter — eine
   [Atomic Parameter List](../structure/AtomicParameterList.md) — anforderbar sein. Anschließend
   Vergleich mit [Wish Template](../quality/WishTemplate.md). Nützlicher Nebenbefund der
   Musterlösung: Die *expansion*-Variante der Wish List optimiert nicht die Datenmenge, sondern
   die Anzahl der Requests.
2. **Reference Management.** Umbau von [Embedded Entity](../quality/EmbeddedEntity.md) auf
   [Linked Information Holder](../quality/LinkedInformationHolder.md).
3. **Quality Management and Governance.** Basiszugriffskontrolle über
   [API Key](../structure/APIKey.md) einführen und Bedingungen in einem
   [Service Level Agreement](../quality/ServiceLevelAgreement.md) dokumentieren. Die Musterlösung
   ist an dieser Stelle bemerkenswert unideologisch: Innerhalb einer Organisation braucht es oft
   kein SLA, wohl aber SLOs — als SRE-Praxis, bei der die SLO-Verletzungsrate gegen das
   Error Budget gerechnet wird und darüber Releases freigegeben werden.

### Tutorial 2 — Alle Kategorien

„A Complete Guide through MAP": Lakeside Mutual will ein Kunden-Self-Service-Frontend bauen, mit
dem Kunden ihre Kontaktdaten selbst ändern können. Gegeben sind eine Context Map aus strategischem
DDD (mit Customer-Supplier-, Open-Host-Service-, Published-Language- und Conformist-Beziehungen),
eine User Story, zwei nichtfunktionale Anforderungen (80 % der Vorgänge unter 2 Sekunden;
10.000 Kunden, davon 10 % gleichzeitig) und ein Architekturüberblicksdiagramm.

| Schritt | Kategorie | Ergebnis |
|---|---|---|
| 1 | [Foundation](category-foundation.md) | Domänenanalyse, API-Scoping, *Candidate Endpoint List* (Endpunkt / Operation / Daten) |
| 2 | [Responsibility](category-responsibility.md) | Rollen und Verantwortungen ergänzen → *Refined Endpoint List*; erste MDSL-Fassung |
| 3 | [Structure](category-structure.md) | Nutzlastdesign und Entscheidungsbegründung |
| 4 | [Quality](category-quality.md) | Nachrichtengrößen balancieren |
| 5 | [Evolution](category-evolution.md) | Evolutionsstrategie wählen und im Vertrag festhalten |

Kernentscheidungen der Musterlösung: Die Kundenressource ist ein
[Master Data Holder](../responsibility/MasterDataHolder.md) (lange Lebensdauer, datenorientiert,
viele eingehende Referenzen). Die GET-Operationen sind
[Retrieval Operations](../responsibility/RetrievalOperation.md), das PUT ist eine
[State Creation Operation](../responsibility/StateCreationOperation.md) — begründet damit, dass
die Änderungsanfrage das verarbeitete Ereignis ist. In Schritt 4 wird die Adressbeziehung von
[Embedded Entity](../quality/EmbeddedEntity.md) auf
[Linked Information Holder](../quality/LinkedInformationHolder.md) umgestellt, weil die
vollständige Umzugshistorie sonst bei jedem Aufruf mitgeliefert würde und die NFRs gefährdet.
In Schritt 5 fällt die Wahl auf [Experimental Preview](../evolution/ExperimentalPreview.md) für
den Erstrelease, mit [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) oder
[Two in Production](../evolution/TwoInProduction.md) sowie
[Version Identifier](../evolution/VersionIdentifier.md) nach
[Semantic Versioning](../evolution/SemanticVersioning.md) für später.

Auch der Begriffsapparat wird geklärt: *Rolle* und *Verantwortung* stammen aus
Responsibility-Driven Design — eine Rolle ist ein Satz zusammengehöriger Verantwortungen (der
Endpunkt), eine Verantwortung eine Verpflichtung, eine Aufgabe auszuführen oder eine Information
zu kennen (die Operation).

## Verwandte Wiki-Seiten

- [Überblick MAP](overview.md) · [Cheat Sheet](cheatsheet.md) · [Terminologie](terms.md)
- Kategorien in Tutorial-2-Reihenfolge: [Foundation](category-foundation.md) → [Responsibility](category-responsibility.md) → [Structure](category-structure.md) → [Quality](category-quality.md) → [Evolution](category-evolution.md)
- [Buch und Ressourcen](book-and-resources.md) — Kapitel 3 des Buchs ist das großformatige Tutorial

## Bezug zu Kubernetes / KRM

Der Fünfschritt aus Tutorial 2 hat im Kubernetes-Umfeld eine klare Entsprechung, allerdings mit
stark verschobenen Gewichten. Schritt 1 und 2 fallen weitgehend weg: Endpunktrolle und
Operationsmenge sind vom Modell vorgegeben (eine Ressource mit `spec`/`status`, Standardverben).
Die *Candidate Endpoint List* wird zur Liste der geplanten Ressourcen und Subresources, die
*Refined Endpoint List* zu den RBAC-Verben je Ressource.

Schritt 3 und 4 sind die eigentliche Arbeit und laufen anders ab: Wo Tutorial 2 in Schritt 4 von
[Embedded Entity](../quality/EmbeddedEntity.md) auf
[Linked Information Holder](../quality/LinkedInformationHolder.md) umstellt, ist genau das in KRM
die schwierigste Änderung überhaupt — bei einer GA-Ressource ist sie brechend und erfordert eine
neue Group-Version mit Konvertierung. Die KRM-Konsequenz lautet daher: Diese Entscheidung muss vor
der Graduierung getroffen werden, nicht nach den ersten Performance-Beschwerden. Genau dafür
existiert die Alpha-Stufe, die dem in Schritt 5 gewählten
[Experimental Preview](../evolution/ExperimentalPreview.md) entspricht.

Für die konkrete Aufgabe „Kontaktdaten ändern" würde KRM zudem nicht auf eine
[State Creation Operation](../responsibility/StateCreationOperation.md) hinauslaufen, sondern auf
`patch` auf `spec` mit anschließender Konvergenz durch einen Controller — Ereignis und Übergang
werden entkoppelt.

---
[← Index](../README.md) · [Überblick](overview.md) · [Quelle](https://microservice-api-patterns.org/patterns/tutorials)
