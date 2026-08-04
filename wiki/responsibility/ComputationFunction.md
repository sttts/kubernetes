---
title: Computation Function
kategorie: Responsibility
unterkategorie: Operation Responsibilities
quelle: https://microservice-api-patterns.org/patterns/responsibility/operationResponsibilities/ComputationFunction
---

# Computation Function

*a.k.a.* Stateless Computation Operation, Calculation Action, Side-Effect-Free/Stateless Operation

**Kurzform:** Eine Operation, die ausschließlich aus ihrer Eingabe ein Ergebnis berechnet — sie liest
keinen Provider-Zustand und schreibt keinen. Die reine Funktion unter den vier
Operationsverantwortlichkeiten.

## Kontext

Die Anforderungen verlangen eine Berechnung. Eingabe liegt lokal vor, das Ergebnis wird lokal
weiterverwendet — dennoch soll die Berechnung nicht lokal laufen, aus Kosten-, Effizienz-, Last-,
Sicherheits- oder anderen Gründen. Beispiele: prüfen lassen, ob Daten bestimmte Bedingungen erfüllen,
oder sie von einem Format in ein anderes konvertieren.

## Problem

Wie kann ein Client seiteneffektfreie Verarbeitung auf der Provider-Seite anstoßen, um aus seiner
Eingabe ein Ergebnis berechnen zu lassen?

## Forces

- **Reproduzierbarkeit und Vertrauen:** Gleiche Eingabe muss gleiches Ergebnis liefern; der Client muss
  der fremden Berechnung trauen können (Korrektheit, aber auch Vertraulichkeit der übergebenen Daten).
- **Performance:** Der Netzwerk-Roundtrip muss sich gegenüber lokaler Berechnung lohnen.
- **Lastmanagement:** Rechenintensive Funktionen sind ein Angriffs- und Überlastvektor, da beliebig oft
  wiederholbar.

## Lösung

Eine Operation `cf: in -> out` wird dem Endpunkt hinzugefügt — meist einer
[Processing Resource](ProcessingResource.md). Sie validiert die eingehende Nachricht, führt die Funktion
aus und gibt das Ergebnis zurück. Sie greift weder lesend noch schreibend auf den Anwendungszustand zu.

## Varianten

- **Validation Service:** prüft übergebene Daten gegen Regeln oder Schemata und meldet Verstöße;
  optional mit Korrekturvorschlägen.
- **Transformation Service:** konvertiert Daten von einem Format oder Standard in ein anderes.
- **Health Check / Ping:** die minimale Ausprägung — eine „I am alive“-Operation mit trivialen Vor- und
  Nachbedingungen, üblicher Bestandteil einer Systems-Management-Strategie für geschäftskritische APIs.

## Konsequenzen

**Vorteile:**

- Maximale Entkopplung: kein geteilter Zustand, keine Reihenfolgeabhängigkeit, keine Transaktionen und
  keine fachliche Kompensation nötig.
- Trivial idempotent, beliebig retrybar und horizontal skalierbar — jede Instanz kann jede Anfrage
  beantworten.
- Ergebnisse sind über die Eingabe cachebar.

**Nachteile / Kosten:**

- Alle nötigen Daten müssen mitgeschickt werden; das kann große Nachrichten und Datenschutzfragen
  erzeugen.
- Der Client trägt die Beweislast für Aktualität der Eingabedaten — der Provider zieht nichts nach.
- Rechenlast ohne natürliche Begrenzung; [Rate Limit](../quality/RateLimit.md) und
  [Pricing Plan](../quality/PricingPlan.md) werden schnell notwendig.

## Bekannte Verwendungen

- Serverless-Lambdas in AWS oder Azure — solange sie nicht mit Cloud-Storage kombiniert werden und
  dadurch zustandsbehaftet werden.
- Online-JSON- und JSON-Schema-Validatoren; Bild- und PDF-Verarbeitungsdienste.
- Ein deutscher Automobilhersteller betreibt einen REST-Level-2-Microservice zur Profilverwaltung, dessen
  *analysis*-Endpunkt per `POST` Daten strikt validiert, Korrekturen vorschlägt und Adressen sowie
  Telefonnummern auf Landesstandards normalisiert — ohne sie zu speichern.
- Terravis berechnet aus einem Wunschdatum unter Berücksichtigung von Feiertagen und Wochenenden das
  nächste gültige Zahlungsdatum; ebenso werden Vertragsdokumente zustandslos als PDF generiert.

## Verwandte Patterns

- [Retrieval Operation](RetrievalOperation.md) — ändert ebenfalls keinen Zustand, konsultiert ihn aber
  lesend; die Computation Function bekommt alles vom Client.
- [State Creation Operation](StateCreationOperation.md) — bekommt ebenfalls alles vom Client, schreibt
  aber Zustand; ihre Clients erwarten nur eine Quittung, nicht ein Ergebnis.
- [State Transition Operation](StateTransitionOperation.md) — liefert ebenfalls nicht-triviale Daten,
  liest und schreibt aber Zustand.
- [Processing Resource](ProcessingResource.md) — die typische Endpunktrolle.
- [Rate Limit](../quality/RateLimit.md) / [Service Level Agreement](../quality/ServiceLevelAgreement.md) —
  Steuerung der Rechenlast.
- [Error Report](../structure/ErrorReport.md) — Transportformat für Validierungsergebnisse.

## Bezug zu Kubernetes / KRM

Das feste Verbset `get`, `list`, `watch`, `create`, `update`, `patch`, `delete`, `deletecollection` sieht
keine berechnende, zustandsneutrale Operation vor. KRM löst das mit einem charakteristischen Kniff:
**`create` auf eine Ressource, die nicht persistiert wird.** Die `*Review`-Ressourcen sind genau das —
das Request-Objekt trägt die Eingabe, das zurückgegebene Objekt trägt in seinem `status` das Ergebnis,
und in etcd landet nichts.

Konkrete Instanzen im Repository: `SubjectAccessReview`, `SelfSubjectAccessReview`,
`LocalSubjectAccessReview` und `SelfSubjectRulesReview` unter `pkg/registry/authorization/` sowie
`TokenReview` und `SelfSubjectReview` unter `pkg/registry/authentication/`. Alle implementieren nur
`Create` und keinerlei Store-Semantik — es gibt für sie kein `get`, `list` oder `watch`, weil es nichts
zu lesen gibt. Ein `SubjectAccessReview` beantwortet „darf Subjekt X das Verb Y auf Ressource Z?“; ein
`TokenReview` prüft ein Bearer-Token und liefert `status.authenticated` samt `status.user`. Das ist
strukturell exakt eine *Computation Function*, nur ausgedrückt im Ressourcen-Vokabular statt als RPC.

```yaml
# Request: POST /apis/authorization.k8s.io/v1/subjectaccessreviews
apiVersion: authorization.k8s.io/v1
kind: SubjectAccessReview
spec:                       # <- die komplette Eingabe steckt im Body, kein Objektname im Pfad
  user: system:serviceaccount:default:builder
  groups: ["system:serviceaccounts", "system:authenticated"]
  resourceAttributes:
    namespace: default
    verb: create
    group: apps
    resource: deployments
  # ... alternativ nonResourceAttributes für /healthz & Co.
```

```yaml
# Response: 201 Created — dasselbe Objekt, status gefüllt, nichts persistiert
apiVersion: authorization.k8s.io/v1
kind: SubjectAccessReview
metadata:
  creationTimestamp: null   # <- kein name, keine uid, keine resourceVersion: es gibt kein Objekt
spec:
  # ... unverändert wie oben zurückgespiegelt
status:
  allowed: true             # <- das Rechenergebnis
  reason: 'RBAC: allowed by RoleBinding "deployer/default" of Role "deployer" to ServiceAccount "builder/default"'
  # ... denied, evaluationError
```

Das leere `metadata` ist der Beleg für die Behauptung: ein `create`, dessen Antwort keine Identität
trägt, hat nichts angelegt.

Die Variante *Validation Service* hat ebenfalls eine wörtliche Entsprechung: Der Parameter
`dryRun=All` (Feld `DryRun` in `CreateOptions`, `UpdateOptions`, `PatchOptions`, `DeleteOptions`) lässt
eine Schreiboperation die komplette Kette aus Decoding, Defaulting, Validierung und Admission durchlaufen
und das Ergebnisobjekt zurückgeben, ohne zu persistieren. `kubectl apply --dry-run=server` und
`kubectl diff` bauen darauf auf: dieselbe Operation, einmal als Zustandsänderung, einmal als reine
Berechnung — unterschieden nur durch einen Query-Parameter.

```http
PATCH /apis/apps/v1/namespaces/default/deployments/web?fieldManager=kubectl&dryRun=All
Content-Type: application/apply-patch+yaml
Accept: application/json

HTTP/1.1 200 OK

# Body: das Deployment, wie es nach Defaulting, Validierung und Admission aussähe.
# Ohne dryRun=All wäre exakt derselbe Aufruf ein Schreibzugriff.
```

Zwei Einschränkungen der Analogie sind ehrlich zu benennen. Erstens sind die Reviews nicht wirklich
eingabe-vollständig im Sinne des Patterns: Ein `SubjectAccessReview` konsultiert die RBAC-Regeln des
Clusters, ein `TokenReview` die Signaturschlüssel — sie lesen also Provider-Zustand und liegen damit
formal näher an einer [Retrieval Operation](RetrievalOperation.md) mit komplexer Auswertung. Zweitens
ist `serviceaccounts/token` (TokenRequest, `pkg/registry/core/serviceaccount/storage/token.go`) zwar
ebenfalls ein nicht-persistierendes `create`, aber nicht seiteneffektfrei, weil es ein neues Credential
ausstellt. KRM erfüllt das Pattern damit (b) anders: nicht durch ein eigenes Verb, sondern durch die
Umdeutung von `create` zu „schicke ein Eingabeobjekt, bekomme ein ausgefülltes Objekt zurück“ — eine
Konsequenz daraus, dass das uniforme Schema über alle Ressourcen wichtiger war als ein passgenaues Verb.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/operationResponsibilities/ComputationFunction)
