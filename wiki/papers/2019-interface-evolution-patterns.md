---
title: Interface Evolution Patterns
kategorie: Paper
autoren: Daniel Lübke, Olaf Zimmermann, Cesare Pautasso, Uwe Zdun, Mirko Stocker
venue: EuroPLoP 2019, 24th European Conference on Pattern Languages of Programs, Irsee
quelle: http://eprints.cs.univie.ac.at/6082/1/WADE-EuroPlop2019Paper.pdf
---

# Interface Evolution Patterns — Balancing Compatibility and Extensibility across Service Life Cycles

## Bibliografische Angaben

Daniel Lübke (iQuest, Hannover), Olaf Zimmermann (OST Rapperswil), Cesare Pautasso (USI Lugano), Uwe Zdun
(Universität Wien), Mirko Stocker (OST Rapperswil): *Interface Evolution Patterns — Balancing Compatibility and
Extensibility across Service Life Cycles.* EuroPLoP '19, 3.–7. Juli 2019, Irsee. ACM, 24 Seiten. DOI
10.1145/3361149.3361164. Dritte Arbeit der MAP-Reihe nach den *Interface Representation Patterns* (2017) und den
*Interface Quality Patterns* (2018), gewonnen aus über 30 öffentlichen Web-APIs und aus Industrieprojekten der
Autoren.

## Worum es geht

Das Paper begründet die MAP-Kategorie *Evolution*. Leitfrage: welche Governance-Regeln Stabilität und Kompatibilität
gegen Wartbarkeit und Erweiterbarkeit ausbalancieren, wenn Provider und Client unterschiedlichen Lebenszyklen
folgen. Sechs Qualitätsziele werden vorangestellt: Kompatibilität und Developer Experience, entkoppelte
Lebenszyklen, minimale erzwungene Client-Änderungen, Freiheit des Providers zur Erweiterung, Vermeidung semantischer
Missverständnisse, minimaler Wartungsaufwand für alte Clients. Die Muster unterscheiden sich vor allem darin,
*welches* dieser Ziele sie priorisieren.

## Behandelte Patterns

| Pattern | Paper-Abschnitt | Wiki-Link |
|---|---|---|
| *API Description* | 4.1, S. 5–8 | [API Description](../foundation/APIDescription.md) |
| *Version Identifier* | 4.2, S. 8–11 | [Version Identifier](../evolution/VersionIdentifier.md) |
| *Semantic Versioning* | 4.3, S. 11–13 | [Semantic Versioning](../evolution/SemanticVersioning.md) |
| *Two in Production* | 4.4, S. 13–16 | [Two in Production](../evolution/TwoInProduction.md) |
| *Limited Lifetime Guarantee* | 4.5, S. 16–17 | [Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) |
| *Eternal Lifetime Guarantee* | 4.6, S. 17–19 | [Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) |
| *Aggressive Obsolescence* | 4.7, S. 19–21 | [Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) |
| *Experimental Preview* | 4.8, S. 21–23 | [Experimental Preview](../evolution/ExperimentalPreview.md) |

## Was das Paper gegenüber der Website ergänzt

Die Website kürzt Forces auf Stichworte und verweist bei den Consequences aufs Buch; im Paper stehen beide
vollständig, dazu bei mehreren Mustern eine *Non-solution*. Verbindlich werden die Zusagen laut Paper in
[API Description](../foundation/APIDescription.md) und
[Service Level Agreement](../quality/ServiceLevelAgreement.md).

### API Description

**Forces:** Information Hiding gegenüber Implementierungsdetails; Interoperabilität zwischen Clients und Providern
auf verschiedenen Middleware-Plattformen; Consumability (Erlernbarkeit, Einfachheit); Erweiterbarkeit und
Evolvierbarkeit als Facetten allgemeiner Modifizierbarkeit.

**Consequences:** Positiv — eine *Minimal Description* ist kompakt und leicht zu pflegen, eine *Elaborate
Description* ausdrucksstark und interoperabilitätsfördernd. Negativ — die minimale Variante verleitet zum Raten oder
Reverse Engineering, wodurch implizite Annahmen das Information Hiding verletzen und langfristig ungültig werden;
Mehrdeutigkeiten schaden der Interoperabilität; Test- und Wartungsaufwand steigen, wenn nicht rückwärtskompatible
Versionen nicht ausgewiesen sind. Die ausführliche Variante erzeugt Inkonsistenzen durch intrinsische Redundanz,
verletzt Information Hiding, sobald sie Downstream-Abhängigkeiten offenlegt, und muss systematisch nachgezogen
werden.

### Version Identifier

**Forces:** Genauigkeit und exakte Identifikation der Version; Minimierung der Auswirkungen von API-Änderungen auf
Client-Seite; Garantie, dass Änderungen die Kompatibilität nicht *versehentlich* auf semantischer Ebene brechen;
Nachverfolgbarkeit der genutzten Versionen für die Governance. *Non-solution:* Versionierung „später“ nachrüsten —
fehlende Governance gilt dem Paper als einer der dominanten Gründe für das Scheitern früherer SOA-Initiativen.

**Consequences:** Positiv — klare Kommunikation über API, Operationen und Nachrichten; geringere Wahrscheinlichkeit
unentdeckter semantischer Änderungen; Nachvollziehbarkeit, welche Payload-Version Clients wirklich benutzen. Negativ
— ein geänderter Identifier kann Clients zum Upgrade zwingen, obwohl sich die von ihnen genutzte Funktionalität
nicht geändert hat. Wichtig aus der *Further discussion*: Das Muster entkoppelt die Lebenszyklen für sich genommen
**nicht**, es ist nur Voraussetzung dafür; und je feiner versioniert wird, desto geringer die Kopplung, aber desto
höher der Governance-Aufwand.

### Semantic Versioning

**Forces:** minimaler Aufwand, Inkompatibilität zu erkennen (besonders auf Client-Seite); Handhabbarkeit der
Versionen und zugehöriger Governance-Aufwand (Freigabeprozesse, Quality Gates, Zahl paralleler Versionen, Zahl der
Versionszweige); Klarheit über die Auswirkung einer Änderung; klare Trennung von Änderungen unterschiedlicher
Tragweite; Klarheit über die Evolutions-Zeitachse. *Non-solution:* laufende Nummern, die nur die Chronologie
kodieren und den Kompatibilitätsgraphen über mehrere Zweige unsichtbar lassen; ebenso Commit-IDs, da nicht jeder
Commit deployt wird.

**Consequences:** Positiv — hohe Klarheit über die Kompatibilitätswirkung zwischen zwei Versionen. Negativ —
erhöhter Aufwand bei der Vergabe, weil die Einordnung einer Änderung mitunter schwerfällt; die Handhabbarkeits-Force
wird nur *teilweise* aufgelöst; und wird das Muster nicht konsequent angewandt, schleichen sich Breaking Changes in
Minor-Updates ein. Der Implementierungshinweis mit der größten Sprengkraft: Patch- und ggf. Minor-Version vor
Clients verbergen — steht die Patch-Version in einem Namespace oder Attributnamen, bricht allein ihre Erhöhung die
Kompatibilität.

### Two in Production

Forces und Lösung sind auf der [Wiki-Seite](../evolution/TwoInProduction.md) bereits vollständig; das Paper ergänzt
vor allem die **Consequences**. Positiv — Clients können Änderungen weit im Voraus planen und müssen nicht genau
dann migrieren, wenn der Provider released; parallele Versionen bieten hohe Kompatibilität *und* Rollback-Fähigkeit;
das Halten alter Clients auf alten Versionen senkt die Wahrscheinlichkeit unentdeckter Kompatibilitätsänderungen;
die Wartungskosten sinken, weil technische Schulden zwischen den Versionen abgebaut werden dürfen, ohne
Rückwärtskompatibilität berücksichtigen zu müssen. Negativ — Clients müssen sich über die Zeit an inkompatible
Änderungen anpassen; die Reaktionsfähigkeit auf dringende Änderungswünsche ist eingeschränkt; der Parallelbetrieb
kostet zusätzlich; und je schneller der Provider ändert, desto kürzer werden die Migrationsintervalle.

### Limited Lifetime Guarantee

**Forces — es sind nur zwei:** clientseitige Änderungen infolge von API-Änderungen planbar machen; den
Wartungsaufwand für alte Clients begrenzen. Die Kürze ist selbst eine Aussage: ein reiner Kompromiss auf der
Zeitachse.

**Consequences:** Positiv — gut planbar durch feste, weit im Voraus bekannte Zeitfenster (übliche Werte sind
Vielfache von sechs Monaten). Negativ — eingeschränkte Reaktionsfähigkeit auf dringende Änderungswünsche; Clients
werden zu einem definierten Zeitpunkt zum Upgrade gezwungen, der mit ihrer Roadmap kollidieren kann; und das Muster
hat keine Antwort auf *abandoned clients*, also Clients, die noch produktiv laufen, aber nicht mehr weiterentwickelt
werden. Lange Garantien verleiten zudem dazu, die ferne Deadline zu ignorieren — worauf faktisch eine *Eternal
Lifetime Guarantee* entsteht.

### Eternal Lifetime Guarantee

**Forces — die längste Liste im Paper:** keine Client-Änderungen aufgrund von API-Änderungen; Möglichkeit des
Providers, die API zu verbessern und an neue Anforderungen anzupassen; minimaler Wartungsaufwand für alte Clients;
Fähigkeit, die API-Infrastrukturtechnologien zu aktualisieren; Fähigkeit, Sicherheitslücken zu schließen;
Anerkennung der Machtverhältnisse zwischen Provider und Client, insbesondere der Fähigkeit von Clients, API-Design
und -Evolution zu steuern.

**Consequences:** Positiv — Clients müssen sich nicht ändern; der Provider wird attraktiver, weil Clients mit
langfristiger Verfügbarkeit rechnen können. Negativ — Innovationschancen werden verpasst; technische Schulden
akkumulieren beim Provider und treiben Wartungs- und Betriebskosten. Die *Further discussion* ist schärfer als die
Website: Das Muster löst alle anderen Kräfte *kontraproduktiv* auf, und gebrochene kryptografische Verfahren lassen
sich womöglich nicht ersetzen, ohne die Rückwärtskompatibilität zu verletzen — daher die Empfehlung, sich das
Brechen der Garantie vorzubehalten.

### Aggressive Obsolescence

**Forces:** Minimierung des Wartungsaufwands; Reduktion erzwungener Client-Änderungen in einem gegebenen Zeitraum;
Anerkennung der Machtbalance zwischen Provider und Client; Respektierung kommerzieller Ziele und Randbedingungen,
etwa der Auswirkungen auf einen Rate Plan. *Non-solution:* gar keine Zusagen oder eine sehr kurze *Limited Lifetime
Guarantee* — beides minimiert die Auswirkung von Änderungen nicht wirklich; die API zum *Experimental Preview* zu
erklären, ist noch schwächer.

**Consequences:** Positiv — im Idealfall müssen Clients sich gar nicht ändern, sofern sie die abgekündigte
Funktionalität nicht nutzen; die Codebasis des Providers bleibt klein und wartbar. Negativ — der Provider muss
ankündigen, was wann abgekündigt und entfernt wird; Clients, die selten genutzte Features verwenden, werden nach
einem Zeitplan zu Änderungen gezwungen, der bei ihrem eigenen Release nicht bekannt war und sich später noch ändern
kann; und Clients müssen die obsoleten Features erst kennenlernen.

Für die Abgrenzung zentral: *Aggressive Obsolescence* arbeitet immer mit **relativen** Zeitfenstern ab dem
Deprecation-Zeitpunkt, *Two in Production* und *Limited Lifetime Guarantee* dagegen mit **absoluten**, am initialen
Release ausgerichteten Fristen — und es entfernt einzelne Repräsentationselemente statt ganzer Versionen.

### Experimental Preview

**Forces:** Raum für Innovation und neue Features; frühes Feedback für den Provider; Fokussierung des Aufwands in
der Frühphase, etwa durch Vermeidung von Governance-Aufwand; frühe Lernmöglichkeiten für Consumer; auf Client-Seite
der Wunsch, sich auf API-Stabilität verlassen zu können. *Non-solution:* erst nach Fertigstellung releasen oder sehr
häufige Releases, die den Governance-Aufwand hochtreiben.

**Consequences:** Positiv — Clients erhalten frühen Zugang zu Innovation und können das Design beeinflussen;
Provider können frei und schnell ändern, bevor sie „stabil“ erklären. Negativ — Provider gewinnen womöglich kaum
Clients, weil die API als unreif wahrgenommen wird; Clients müssen ihre Implementierung bis zum Stable-Release
fortlaufend anpassen; und Clients riskieren den Totalverlust ihrer Investition, falls nie eine stabile Version
erscheint oder der Preview plötzlich verschwindet. Nur im Paper: Ein „Beta for life“-Ansatz hilft Clients nicht,
weil sie irgendwann eine per SLA abgesicherte API brauchen.

## Pattern-Sprache und Zusammenhänge

**Das Kompatibilitätsmodell (S. 5).** Kompatibilität ist ausdrücklich eine *Eigenschaft der Beziehung* zwischen
Provider und Client, nicht der API: Beide sind kompatibel, wenn sie ihren Nachrichtenaustausch durchführen und alle
Nachrichten gemäß der Semantik der jeweiligen Version korrekt interpretieren können. Provider und Client derselben
Version *n* sind per Definition kompatibel; ist ein Client für *n* kompatibel mit dem Provider *n − 1*, ist dieser
**vorwärtskompatibel**, ist er kompatibel mit dem Provider *n + 1*, ist dieser **rückwärtskompatibel**. Daraus folgt
die zentrale Beobachtung: Sobald sich die Lebenszyklen nicht mehr synchronisieren lassen — und beim
Zero-Downtime-Deployment mit rollierendem Instanztausch ist das nie der Fall —, muss beim Entwurf mit *mehreren
Client-Versionen gegen mehrere API-Versionen* gerechnet werden.

**Die Lebenszyklus-Modelle in den Abbildungen.** Abbildung 7 zeigt *Two in Production* als **gleitendes Fenster**
aktiver Versionen: Beim Release wird die älteste noch laufende zurückgezogen, verbleibende Clients werden informiert
und auf Protokollebene umgeleitet. Abbildung 8 zerlegt *Aggressive Obsolescence* in einen **vierstufigen Prozess**:
(0) Version läuft produktiv, (1) Deprecation mit Entfernungstermin, (2) Clients migrieren oder wechseln den
Provider, (3) Removal — Anfragen scheitern oder werden umgeleitet, (4) nicht migrierte Clients funktionieren nicht
mehr. Abbildung 9 kombiniert *Experimental Preview* in einer ungoverneten Sandbox mit *Two in Production* für den
produktiven Teil.

**Die Abhängigkeitsstruktur (Abb. 2).** *Version Identifier* ist die Wurzel: *Semantic Versioning* setzt ihn voraus,
*Two in Production* benötigt ihn zwingend, die übrigen Lebenszyklus-Muster *können* ihn nutzen. Die
Lebenszyklus-Muster bilden eine Skala des Commitments — von *Experimental Preview* (schwächste Zusage, gefolgt von
*Aggressive Obsolescence*) über *Two in Production* und *Limited Lifetime Guarantee* bis zur *Eternal Lifetime
Guarantee*. Dazu beschreibt das Paper Degenerationspfade: Eine nicht durchgesetzte Löschpolitik in *Two in
Production* wird implizit zur *Eternal Lifetime Guarantee*, eine zu lange *Limited Lifetime Guarantee* ebenso;
umgekehrt kippt eine untragbar gewordene *Eternal Lifetime Guarantee* in eines der begrenzten Muster.

## Bemerkenswerte Aussagen

Zur *Eternal Lifetime Guarantee*: „designers essentially freeze the API at the cost of also freezing innovation“ (S.
18) — die Garantie friert nicht nur den Vertrag ein, sondern den technischen Fortschritt gleich mit. Dass die
Strategiewahl nicht allein beim Provider liegt, zeigt das IBAN-Beispiel: Die Frist setzte das Europäische Parlament.

## Bezug zu Kubernetes / KRM

KRM geht in der Versionierung am eigenständigsten vor — in vier Punkten.

**1. Der Version Identifier steckt in der Payload, nicht nur im Pfad.** Das Paper empfiehlt, den Identifier an
*genau einer* Stelle zu führen; Kubernetes führt ihn im Pfad (`/apis/apps/v1/deployments`) *und* im Objekt
(`apiVersion: apps/v1` aus `TypeMeta`), hält beide aber konsistent. Entscheidend ist die Payload-Variante: Weil
`apiVersion` Teil des serialisierten Objekts ist, trägt ein Manifest auf der Platte, ein Eintrag in etcd und ein
Ereignis im `watch`-Stream seine Version mit sich — die Voraussetzung dafür, dass `kubectl apply -f` ohne Pfadwissen
funktioniert.

**2. Reifegrad statt SemVer.** Für API-Gruppen gibt es kein Dreier-Tripel, sondern `v1alpha1`, `v1beta1`, `v1` — der
Identifier kodiert nicht die Änderungstiefe, sondern die *Stabilitätszusage*. Aus dem Vergleich zweier Nummern die
Kompatibilität abzulesen entfällt damit; stattdessen liest man ab, *welches Lebenszyklusmuster* gilt: Alpha
entspricht dem *Experimental Preview*, Beta einer *Limited Lifetime Guarantee*, GA kommt innerhalb von `v1` einer
*Eternal Lifetime Guarantee* nahe. Semantic Versioning gilt in Kubernetes für das *Release* (`v1.34.2`), nicht für
die API-Gruppenversion — genau die vom Paper angemahnte Trennung von Interface- und Implementierungsversion.

**3. Two in Production, aber auf demselben Objekt.** Der API-Server serviert beliebig viele Versionen einer
Ressource gleichzeitig — `flowcontrol.apiserver.k8s.io` hatte zeitweise `v1beta1`, `v1beta2`, `v1beta3` und `v1`
parallel. Anders als im Paper, wo die parallelen Versionen ausdrücklich *nicht* kompatibel sein müssen, sind es in
KRM Repräsentationen **desselben Objekts**: gleiche `metadata.uid`, gleicher Datensatz. Erzwungen wird das durch
verlustfreie Round-Trip-Konversion über eine interne Hub-Version — `runtime.APIVersionInternal` = `__internal` in
`staging/src/k8s.io/apimachinery/pkg/runtime/interfaces.go`. Konvertiert wird sternförmig (`v1beta1 → __internal →
v1`), was den Aufwand von O(n²) auf O(n) Konverter senkt. Für CRDs übernimmt `spec.conversion.strategy` dieselbe
Rolle mit den Werten `None` und `Webhook` (`NoneConverter`, `WebhookConverter` in
`.../apiextensions-apiserver/pkg/apis/apiextensions/v1/types.go`); pro Eintrag in `spec.versions[]` steuern `served`
und `storage`, ob eine Version ausgeliefert bzw. persistiert wird. Damit dreht KRM einen Punkt um, den das Paper nur
als Kostenposten führt: Ein Feld, das sich nicht verlustfrei auf die alte Version abbilden lässt, ist nicht baubar,
solange diese serviert wird.

**4. Deprecation ist maschinenlesbar und wird gemessen.** *Aggressive Obsolescence* verlangt, Abkündigungen „früh,
klar und öffentlich“ zu kommunizieren — bei anonymen Clients das schwierigste Teilproblem. Kubernetes löst es im
Protokoll: Antworten für abgekündigte Versionen tragen einen HTTP-`Warning`-Header mit Code 299
(`.../apiserver/pkg/endpoints/filters/warning.go`), dessen Text aus den generierten `APILifecycle*`-Funktionen der
`zz_generated.prerelease-lifecycle.go`-Dateien stammt (eingehängt über `deprecation.WarningMessage`); für CRDs gibt
es dieselbe Zusage deklarativ als `spec.versions[].deprecated` und `.deprecationWarning`. Beispiel: `FlowSchema` in
`flowcontrol/v1beta3` ist mit „eingeführt 1.26, deprecated 1.29, entfernt 1.32, Ersatz
`flowcontrol.apiserver.k8s.io/v1`“ annotiert — die relative Frist ab Deprecation, hier als Standardregel
„Deprecation plus drei Minor-Releases“. Die Governance-Force aus *Version Identifier* ist ihrerseits als Metrik
implementiert: `apiserver_requested_deprecated_apis` (`.../apiserver/pkg/endpoints/metrics/metrics.go`) mit den
Labels `group`, `version`, `resource`, `subresource` und `removed_release`. Genau das empfiehlt das Paper — hier ist
es kein Ratschlag, sondern eine stabile Metrik.

## Verwandte Wiki-Seiten

Evolution: [Version Identifier](../evolution/VersionIdentifier.md) ·
[Semantic Versioning](../evolution/SemanticVersioning.md) · [Two in Production](../evolution/TwoInProduction.md) ·
[Limited Lifetime Guarantee](../evolution/LimitedLifetimeGuarantee.md) ·
[Eternal Lifetime Guarantee](../evolution/EternalLifetimeGuarantee.md) ·
[Experimental Preview](../evolution/ExperimentalPreview.md) ·
[Aggressive Obsolescence](../evolution/AggressiveObsolescence.md) ·
[Kategorie Evolution](../meta/category-evolution.md).

---
[← Index](../README.md) · [Papers](../papers/) · [Quelle](http://eprints.cs.univie.ac.at/6082/1/WADE-EuroPlop2019Paper.pdf)
