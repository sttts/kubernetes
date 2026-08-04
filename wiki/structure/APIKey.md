---
title: API Key
kategorie: Structure
unterkategorie: Special Purpose Representations
quelle: https://microservice-api-patterns.org/patterns/structure/specialPurposeRepresentations/APIKey
---

# API Key

*a.k.a.* Access Token, Provider-Allocated Client Identifier

**Kurzform:** Der Provider vergibt jedem Client ein eindeutiges, providerseitig erzeugtes Token, das
der Client bei jedem Aufruf mitschickt. Damit lässt sich der Aufrufer identifizieren und
authentifizieren, ohne dass echte Benutzerkonto-Credentials über die Leitung gehen.

## Kontext

Ein API-Provider bedient nur registrierte Teilnehmer. Clients haben sich angemeldet — typischerweise,
weil ein *Rate Limit*, ein *Pricing Plan* oder eine Abrechnung daran hängt. Um solche Regeln
durchzusetzen, muss der Provider wissen, wer gerade ruft.

## Problem

Wie kann ein API-Provider seine Clients und deren einzelne Requests identifizieren und
authentifizieren?

## Forces

- **Basissicherheit** — irgendeine Form von Authentifizierung muss es geben, auch wenn keine
  vollständige Identitätsinfrastruktur vorhanden ist.
- **Zugriffskontrolle** — die Identität ist Voraussetzung für Autorisierungsentscheidungen.
- **Keine Benutzer-Credentials speichern oder übertragen** — der Provider will Passwörter der
  Endnutzer weder kennen noch halten.
- **Entkopplung von Client und Organisation** — ein Token soll unabhängig von Personal-Fluktuation
  in der Client-Organisation gültig bleiben (und einzeln widerrufbar sein).
- **Sicherheit vs. Bedienbarkeit** — ein Token ist trivial zu benutzen, aber auch trivial zu
  kopieren, wenn es einmal geleakt ist.
- **Performance** — die Prüfung passiert bei *jedem* Request und darf nicht teuer sein.

## Lösung

Jedem Client wird ein eindeutiges Token — der *API Key* — zugewiesen, das er dem Endpoint zur
Identifikation vorlegt. Üblich ist der Transport im `Authorization`-Header nach RFC 7235, meist mit
dem `Bearer`-Schema aus RFC 6750. Der Key ist providerseitig erzeugt, langlebiger als eine Session
und einem Client-Account (nicht einem Endnutzer) zugeordnet.

## Beispiel

Aufruf der Cloud-Convert-API; der Key identifiziert den Client für die Abrechnung:

```http
POST https://api.cloudconvert.com/process
Authorization: Bearer gqmbwwB74tToo4YOPEsev5
Content-Type: application/json

{
    "inputformat": "docx",
    "outputformat": "pdf"
}
```

## Konsequenzen

**Vorteile:**

- Identifikation ohne Übertragung von Benutzerkonto-Credentials.
- Voraussetzung für [Rate Limit](../quality/RateLimit.md), [Pricing Plan](../quality/PricingPlan.md)
  und nutzungsbasierte Abrechnung.
- Sehr geringer Implementierungsaufwand auf Client-Seite; ein Header genügt.
- Einzelne Keys können widerrufen oder eingeschränkt werden, ohne andere Clients zu stören.

**Nachteile / Kosten:**

- Ein Bearer-Token ist ein reines Inhaberpapier: Wer es hat, ist der Client. Ohne TLS und ohne
  sorgfältiges Secret-Handling ist es wertlos.
- Der Provider braucht Key-Lifecycle: Ausgabe, Rotation, Widerruf, Ablauf.
- Keys landen leicht in Repos, Logs, URLs oder Mobile-App-Binaries. Gegenmaßnahmen wie
  Einschränkung auf IP-Bereiche, Referrer oder App-Signatur sind Zusatzaufwand.
- Ein API Key allein sagt nichts über den *handelnden Endnutzer* — für dessen Identität braucht es
  OAuth 2.0 / OpenID Connect.

## Bekannte Verwendungen

- **YouTube Data API**: OAuth 2.0 und API Keys parallel; getrennte Server-, Browser-, iOS- und
  Android-Keys, jeweils auf App, IP-Bereich oder Domain eingeschränkt, damit ein aus einer App
  extrahierter Key außerhalb nutzlos ist.
- **GitHub API**: primär OAuth, aber Basic Auth mit Benutzername und Token wird unterstützt — der
  Key steht dann im Passwort-Feld statt in einem eigenen Header.
- **Stripe**: *publishable key* (nur Account-Identifikator) und *secret key* (der eigentliche API
  Key, im `Authorization`-Header). Anders als bei Amazons Secret Key, der nie übertragen, sondern
  nur zum Signieren benutzt wird.

## Verwandte Patterns

- [Metadata Element](MetadataElement.md) — ein API Key ist ein Metadatum, kein fachliches Datum.
- [Id Element](IdElement.md) — der Key identifiziert, ähnlich einem Identifier, aber mit
  Authentifizierungssemantik.
- [Context Representation](ContextRepresentation.md) — bündelt Credentials und andere
  Kontextinformationen protokollunabhängig in der Payload.
- [Rate Limit](../quality/RateLimit.md) und [Pricing Plan](../quality/PricingPlan.md) — brauchen
  einen identifizierten Client als Bezugspunkt.
- [Service Level Agreement](../quality/ServiceLevelAgreement.md) — Zusicherungen gelten pro Client.
- [Public API](../foundation/PublicAPI.md) und [Community API](../foundation/CommunityAPI.md) —
  der typische Sichtbarkeitskontext für Key-basierte Registrierung.

Außerhalb von MAP: *Session Identifier* (Fowler 2002) löst ein ähnliches Problem, ist aber auf eine
einzelne Session begrenzt; RBAC/ABAC ergänzen den Key um Autorisierung.

## Bezug zu Kubernetes / KRM

Kubernetes erfüllt das Pattern **auf Transportebene genauso, im Payload-Design aber bewusst gar
nicht**. Credentials stehen nie im Request-Body einer Ressource: Der Client schickt
`Authorization: Bearer <token>` oder authentifiziert sich per Client-Zertifikat im TLS-Handshake.
Ein ServiceAccount-Token ist damit exakt ein API Key im Sinne des Patterns — nur dass es ein
signiertes JWT mit Audience und Ablaufzeit ist.

```http
POST /api/v1/namespaces/default/pods HTTP/1.1
Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ik...   # <- der Key, ausschließlich hier
Content-Type: application/yaml

apiVersion: v1
kind: Pod
metadata:
  name: web
spec:
  serviceAccountName: web        # <- Identität für *ausgehende* Aufrufe des Pods;
                                 #    im Payload steht kein Credential
  # ... containers
```

Ungewöhnlich und interessant ist, dass Kubernetes die *Ausstellung* der Credentials selbst
deklarativ modelliert. Die `TokenRequest`-Subresource (`serviceaccounts/token`,
`authentication.k8s.io/v1`) ist ein gewöhnliches KRM-Objekt mit `spec` und `status`:

```yaml
# POST /api/v1/namespaces/default/serviceaccounts/web/token
apiVersion: authentication.k8s.io/v1
kind: TokenRequest
spec:
  audiences: ["https://kubernetes.default.svc"]
  expirationSeconds: 3600
  boundObjectRef:                # <- Token gilt nur, solange dieses Objekt existiert
    kind: Pod
    apiVersion: v1
    name: web
    uid: 8f2a1c0e-3b7d-4c11-9a55-2f0e6d7b1c34
status:                          # <- vom Server befüllt, in der Antwort
  token: eyJhbGciOiJSUzI1NiIsImtpZCI6Ik...
  expirationTimestamp: "2026-08-04T11:00:00Z"
```

Ein Pod bekommt dasselbe Token über eine `projected`-Volume-Quelle; das Kubelet rotiert es
automatisch:

```yaml
spec:
  volumes:
    - name: token
      projected:
        sources:
          - serviceAccountToken:
              audience: vault.example.com   # <- Key ist an genau diesen Empfänger gebunden
              expirationSeconds: 3600       # <- kurzlebig statt statisch
              path: token
  # ... containers mit volumeMounts
```

Die Gegenrichtung — Zertifikate — läuft über `CertificateSigningRequest` (`certificates.k8s.io/v1`) mit `signerName`,
etwa `kubernetes.io/kube-apiserver-client`. Validiert wird ein vorgelegtes Token über die
`TokenReview`-API. Die Ausgabe des Keys ist also selbst eine KRM-Ressource, während der Key im
Betrieb reines Protokoll-Artefakt bleibt.

Die Trennung von Authentifizierung und Autorisierung ist bei Kubernetes schärfer als im Pattern:
Der Key liefert nur `username`, `uid`, `groups` und `extra`; die Entscheidung fällt danach in RBAC
(`Role`/`ClusterRole` plus Bindings) und ist über `SubjectAccessReview` selbst wieder als API
abfragbar. Zusätzlich gibt es Impersonation über die Header `Impersonate-User`,
`Impersonate-Group`, `Impersonate-Uid` und `Impersonate-Extra-<key>` — ein zweites Identitäts-Set
neben dem eigentlichen Key, ebenfalls strikt im Header, nicht in der Payload.

```http
GET /api/v1/namespaces/default/pods HTTP/1.1
Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ik...   # <- wer wirklich ruft
Impersonate-User: alice                                   # <- als wen er ruft
Impersonate-Group: developers
Impersonate-Uid: 3f7d1a20-9c44-4b2e-8f01-6de5a4c73b19
Impersonate-Extra-scopes: read-only
```

Clientseitig abstrahiert das `client.authentication.k8s.io`-`ExecCredential`-Protokoll die Beschaffung des
Tokens weg, sodass kubeconfig-Dateien nicht zwingend langlebige Keys enthalten.

Nicht abgebildet ist der Abrechnungsaspekt: Kubernetes kennt keinen Pricing Plan und identifiziert
Clients nicht, um sie zu belasten, sondern nur, um RBAC und (grob) Priority-and-Fairness-Flows zu
bestimmen. Statische, langlebige API Keys gelten als Anti-Pattern — der klassische
`kubernetes.io/service-account-token`-Secret ist der Legacy-Weg, der zugunsten kurzlebiger,
gebundener TokenRequest-Tokens verdrängt wurde. `bootstrap.kubernetes.io/token` bleibt als
bewusst kurzlebiger Sonderfall für Node-Joins.

---
[← Index](../README.md) · [Kategorie Structure](../meta/category-structure.md) · [Quelle](https://microservice-api-patterns.org/patterns/structure/specialPurposeRepresentations/APIKey)
