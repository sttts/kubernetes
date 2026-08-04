---
title: Data Transfer Resource
kategorie: Responsibility
unterkategorie: Information Holder Endpoint Types
quelle: https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/DataTransferResource
---

# Data Transfer Resource

*a.k.a.* Connector Resource, Integration Resource, Share, Temporary Data Store, Transient Information Holder

**Kurzform:** Ein geteilter Speicherendpunkt unter global eindeutiger Adresse, über den mehrere
Clients Daten austauschen, ohne einander zu kennen und ohne gleichzeitig aktiv zu sein — ein
Blackboard als Web-Ressource.

## Kontext

Zwei oder mehr Kommunikationsteilnehmer wollen Daten austauschen. Ihre Anzahl ändert sich über die
Zeit, und sie kennen einander nur teilweise. Sie sind nicht notwendigerweise gleichzeitig aktiv —
weitere Teilnehmer können auf dieselben Daten zugreifen wollen, *nachdem* die Quelle sie bereits
geteilt hat. Teilnehmer sind oft nur an der jeweils *neuesten* Version interessiert und müssen nicht
jede Änderung beobachten. Und sie können möglicherweise keine Messaging-Middleware jenseits einer
einfachen HTTP-Client-Bibliothek lokal installieren.

## Problem

Wie können zwei oder mehr Kommunikationsteilnehmer Daten austauschen, ohne einander zu kennen, ohne
gleichzeitig verfügbar zu sein, und selbst dann, wenn die Daten bereits gesendet wurden, bevor die
Empfänger bekannt waren?

## Forces

- **Kopplung in der Zeitdimension** — Sender und Empfänger müssen entkoppelt sein.
- **Kopplung in der Ortsdimension** — sie sollen einander nicht adressieren müssen.
- **Kommunikationsbeschränkungen** — nur HTTP verfügbar, keine Middleware installierbar.
- **Zuverlässigkeit** und **Skalierbarkeit.**
- **Speicherplatzeffizienz** — die Daten liegen dauerhaft, bis jemand aufräumt.
- **Latenz.**
- **Eigentümerschaft** — wem gehören die abgelegten Daten, und wer darf sie löschen?

## Lösung

Führe eine *Data Transfer Resource* als geteilten Speicherendpunkt ein, der von zwei oder mehr
API-Clients erreichbar ist. Gib dieser spezialisierten
[Information Holder Resource](InformationHolderResource.md) eine global eindeutige Netzwerkadresse,
damit die Clients sie als gemeinsamen Datenaustauschraum nutzen können. Statte sie mit mindestens
einer [State Creation Operation](StateCreationOperation.md) und einer
[Retrieval Operation](RetrievalOperation.md) aus.

Entscheide bewusst über Datenbesitz und dessen Übergang — und bevorzuge in diesem Fall
Client-Eigentümerschaft gegenüber Provider-Eigentümerschaft.

## Varianten

Die Quelle beschreibt eine direkte **HTTP-Realisierung**: Client A legt die Information per `PUT` auf
einer per URI eindeutig identifizierten geteilten Ressource ab, Client B holt sie per `GET`. Die
Information verschwindet nicht, solange kein Client explizit `DELETE` aufruft. Weil `PUT` idempotent
ist, kann A zuverlässig veröffentlichen; ein fehlgeschlagenes `GET` kann B einfach wiederholen. Eine
zweite, warteschlangenbasierte Variante wird unter *Bekannte Verwendungen* beschrieben.

## Beispiel

In einem Versicherungsszenario ist `ClaimReceptionSystemOfEngagement` die Datenquelle; eine
`ClaimTransferResource` entkoppelt davon die beiden Senken `ClaimProcessingSystemOfRecords` und
`FraudDetectionArchive`. Der Zugriff wird über einen [API Key](../structure/APIKey.md) kontrolliert.

## Konsequenzen

**Vorteile:**

- Vollständige zeitliche und örtliche Entkopplung ohne Messaging-Middleware — nur HTTP nötig.
- Nachzügler bekommen die Daten auch dann noch, wenn sie beim Senden noch nicht existierten.
- Idempotentes Schreiben und wiederholbares Lesen ergeben ein einfaches Zuverlässigkeitsmodell.

**Nachteile / Kosten:**

- Ohne explizites `DELETE` wächst der Speicher unbegrenzt; Lebenszyklus und Aufräumen müssen
  vereinbart werden.
- Nur der aktuelle Stand ist sichtbar — wer jede Änderung braucht, ist mit echtem Messaging besser
  bedient.
- Empfänger müssen pollen, was Latenz gegen Last eintauscht.
- Der Endpunkt wird zum geteilten Zustand zwischen sonst unabhängigen Parteien: eine klassische
  Blackboard-Kopplung.

## Bekannte Verwendungen

- Dropbox und ownCloud haben aus diesem Pattern ein SaaS-Geschäftsmodell gemacht und bieten
  entsprechende Integrations-APIs; Doodle ist keine API, nutzt das Muster aber ebenso.
- Ein deutscher Automobilhersteller betreibt eine warteschlangenbasierte Variante (Amazon SQS): eine
  API-Operation füllt die Queue, eine andere lässt Clients per Web-Ressource nach Neuzugängen sehen;
  gibt es keine, blockiert der Aufruf nicht, sondern antwortet nach 30 Sekunden mit „komm später
  wieder“. Clients wählen ihre Interessen, der Server filtert zusätzlich nach Rollen.

## Verwandte Patterns

- [Information Holder Resource](InformationHolderResource.md) — das Oberpattern. Der Unterschied:
  eine *Data Transfer Resource* besitzt und kontrolliert ihren Datenspeicher exklusiv, und der
  einzige Zugang ist ihre publizierte API. Andere Typen arbeiten oft mit Daten, auf die auch andere
  Parteien zugreifen. Zudem ist sie zugleich Quelle *und* Senke.
- [Link Lookup Resource](LinkLookupResource.md) — exponiert nur Metadaten, während eine
  *Data Transfer Resource* beliebige Daten hält.
- [Operational Data Holder](OperationalDataHolder.md), [Master Data Holder](MasterDataHolder.md),
  [Reference Data Holder](ReferenceDataHolder.md) — die Typen mit provider-definierter Datensemantik.
- [State Creation Operation](StateCreationOperation.md) und [Retrieval Operation](RetrievalOperation.md)
  — das Minimum an Operationen.
- [API Key](../structure/APIKey.md) — Zugriffskontrolle auf dem geteilten Raum.
- [Conditional Request](../quality/ConditionalRequest.md) — macht das Pollen der Empfänger billiger.

Konzeptionell ist das Pattern eine Web-Realisierung des *Message Channel* aus den Enterprise
Integration Patterns, verwandt mit *Blackboard* (POSA 1) und dem Remoting-Stil *shared repository*.

## Bezug zu Kubernetes / KRM

Kubernetes benutzt dieses Pattern **überall — aber ohne es je so zu nennen** (Fall b). Der gesamte
Controller-Verbund kommuniziert nicht per Nachrichten, sondern über den API-Server als geteiltes
Blackboard: ein Controller schreibt ein Objekt, ein anderer liest es, und beide kennen einander
nicht. Genau die Forces des Patterns — Entkopplung in Zeit und Ort, keine Middleware beim Client,
Zuverlässigkeit durch Wiederholbarkeit — sind die Begründung für dieses Design.

Konkrete Instanzen:

- **Bootstrap-Discovery**: die `cluster-info`-ConfigMap im `kube-public`-Namespace, die kubeadm
  anlegt, damit beitretende Knoten Endpunkt und CA finden, bevor sie authentifiziert sind. Ein
  Schreiber, viele unbekannte, später auftauchende Leser — der Lehrbuchfall.

  ```yaml
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: cluster-info
    namespace: kube-public       # <- kubeadm bindet system:anonymous per Role auf genau dieses Objekt
  data:
    kubeconfig: |                # <- Adresse und CA, damit ein Knoten überhaupt erst vertrauen kann
      apiVersion: v1
      kind: Config
      clusters:
        - cluster:
            server: https://10.0.0.1:6443
            certificate-authority-data: LS0tLS1CRUdJTi...
      # ... contexts, users — leer, hier ist noch niemand authentifiziert
    jws-kubeconfig-07401b: eyJhbGciOiJIUzI1NiIsImtpZCI6IjA3NDAxYiJ9..DFRxs   # <- Signatur je Bootstrap-Token
  ```

- **Zertifikatsaustausch**: ein Client legt eine `CertificateSigningRequest` ab, ein Signer-Controller
  füllt `status.certificate`, der Client holt sie ab. Weder kennen sich die beiden, noch müssen sie
  gleichzeitig laufen.

  ```yaml
  # Schritt 1 — der Kandidat legt ab, was nur er hat, und geht:
  apiVersion: certificates.k8s.io/v1
  kind: CertificateSigningRequest
  metadata:
    name: node-csr-bWFpbg
  spec:
    request: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVFVRVNULS0t...   # <- base64(PEM des CSR)
    signerName: kubernetes.io/kube-apiserver-client-kubelet    # <- adressiert den Empfänger, nicht per URL
    usages: ["digital signature", "client auth"]
    # ... username, groups, uid — vom API-Server aus dem Aufrufer gefüllt
  ```

  ```yaml
  # Schritt 2 — irgendwann später legt der Signer seine Antwort an dieselbe Stelle,
  # über PUT .../certificatesigningrequests/node-csr-bWFpbg/status:
  status:
    conditions:
      - type: Approved       # zuvor vom Approver über das eigene /approval-Subresource gesetzt
        status: "True"
        reason: AutoApproved
        lastTransitionTime: "2026-08-04T10:00:01Z"
    certificate: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t...       # <- die Senke wird zur Quelle
  ```

- **Volume-Übergabe**: `PersistentVolumeClaim` → `PersistentVolume` → `VolumeAttachment`. Jedes
  Objekt ist ein Ablageort, an dem eine Partei etwas hinterlegt und eine andere es aufnimmt.
- **`Lease` bei der Leader Election** — ein geteilter Ort, an dem Kandidaten ihre Identität
  hinterlegen und alle anderen sie lesen.

Zwei Abweichungen von der MAP-Lösung sind aufschlussreich. Erstens ersetzt **`watch` das Polling**:
statt „komm später wieder“ hält der Empfänger einen Änderungsstrom offen, was Latenz *und* Last
senkt. Zweitens ist **Eigentümerschaft anders gelöst**, als das Pattern empfiehlt: MAP rät zu
Client-Eigentümerschaft, KRM legt sie in die Ressource selbst — `metadata.ownerReferences` plus
Garbage Collection räumen abgeleitete Objekte automatisch auf, wenn ihr Besitzer verschwindet, und
`metadata.finalizers` verzögert das Löschen, bis alle Beteiligten fertig sind. Damit löst KRM das
Speicherplatz-Force strukturell, statt es Clients per `DELETE`-Disziplin aufzubürden. Der bewusst
nicht gelöste Teil bleibt: der API-Server ist *kein* Message Broker, jede Änderung geht durch etcd,
und ein Empfänger sieht bei Verbindungsverlust nur den aktuellen Stand — Level-triggered statt
edge-triggered, exakt die Einschränkung, die das Pattern als „nur die neueste Version“ beschreibt.

---
[← Index](../README.md) · [Kategorie Responsibility](../meta/category-responsibility.md) · [Quelle](https://microservice-api-patterns.org/patterns/responsibility/informationHolderEndpointTypes/DataTransferResource)
