---
title: "One Schema to Rule Them All: The Kubernetes API, Measured Against a Pattern Language"
language: en
status: draft
---

# One Schema to Rule Them All

## The Kubernetes API, measured against a pattern language

There is a large body of academic and practitioner work on how to design message-based APIs. The
most systematic piece of it is *Microservice API Patterns* (MAP) by Olaf Zimmermann, Mirko Stocker,
Daniel Lübke, Cesare Pautasso and Uwe Zdun — 45 patterns mined over six years of EuroPLoP and ICSOC
papers, later consolidated into *Patterns for API Design* (Addison-Wesley, 2023).

Kubernetes has never been part of that conversation. Its API model grew out of Borg, out of `etcd`,
out of a specific idea about how distributed systems should be controlled — not out of the REST/SOA
design literature. And yet the Kubernetes Resource Model (KRM) is, by any reasonable measure, one of
the most thoroughly specified message-based APIs ever shipped. Millions of manifests depend on it.
Its deprecation policy is stricter than most commercial products'.

So: what happens when you hold KRM up against a pattern language it never consulted?

I went through all 45 patterns and classified each as (a) implemented, (b) solved differently, or
(c) absent. The result is not a report card. It is a map of where Kubernetes deliberately spent its
design budget, and what it bought.

---

## The thesis

**MAP is a decision language.** Every pattern answers a design question you can ask afresh for each
endpoint. What role does this endpoint play — does it hold data, or does it do something? What
operations does it offer? How big are its messages? Does the client get to shape the response?

The ICSOC 2018 paper puts numbers on this for the quality dimension alone: six architectural
decisions, forty options, forty-seven drivers. The authors are honest that even *with* their
decision model, the option space stays vast.

**KRM is the answer to almost all of those questions, taken once and frozen for every resource that
will ever exist.** One verb set. One message frame. One versioning discipline. One error type.

```yaml
# Every resource in the cluster — core types, CRDs written yesterday — has this shape.
apiVersion: apps/v1          # <- group/version, in the payload, not just the path
kind: Deployment
metadata:                    # <- uniform ObjectMeta, identical across all kinds
  name: web
  namespace: shop
  # ... uid, resourceVersion, generation, labels, annotations, ownerReferences, managedFields
spec:                        # <- desired state, owned by the client
  # ...
status:                      # <- observed state, owned by the controller
  # ...
```

```console
$ kubectl get --raw /apis/apps/v1 | jq '.resources[] | select(.name=="deployments") | .verbs'
["create","delete","deletecollection","get","list","patch","update","watch"]
```

The cost is expressiveness, precisely where MAP is most differentiated. The payoff is that generic
machinery — `kubectl`, informers, RBAC, admission, Server-Side Apply, version conversion, the
garbage collector — works on resources that did not exist when the machinery was written. A CRD
author gets all of it for free, and in exchange gets no say in any of it.

Here is the scorecard.

| | Count | |
|---|---|---|
| **(a) Same, often stricter** | 27 | KRM implements the pattern, frequently beyond what MAP asks |
| **(b) Solved differently** | 9 | Same force, different mechanism |
| **(c) Absent** | 9 | Deliberately missing — some real gaps, some non-goals |

---

## Where Kubernetes goes beyond the patterns

### 1. `spec`/`status` is a pattern MAP does not have

In MAP, a message is either input or output. A request representation and a response representation
are different things, designed separately.

In KRM, **the same object is both**, and write authority is split field-wise.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  generation: 7              # <- bumped by the apiserver on every spec change
spec:
  replicas: 5                # <- the client wrote this
  # ... selector, template
status:
  observedGeneration: 7      # <- the controller has seen generation 7
  replicas: 5
  readyReplicas: 4           # <- ... but only 4 are ready yet
  conditions:
    - type: Progressing
      status: "True"
      reason: NewReplicaSetAvailable
      lastTransitionTime: "2026-08-04T09:12:31Z"
```

This is not a convention — it is enforced. `/status` is a separate RBAC resource:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: deployment-controller
rules:
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "list", "watch"]          # <- may read spec, may not change it
  - apiGroups: ["apps"]
    resources: ["deployments/status"]        # <- separate resource, separate grant
    verbs: ["update", "patch"]
```

This is not a variant of an existing pattern. It is a missing element of the pattern language, and
almost everything else in KRM follows from it. Because an object is simultaneously an instruction
and a report, Kubernetes barely needs *Processing Resources* and has essentially no *State
Transition Operations*.

The 2020 EuroPLoP responsibility paper has a line about how provider-side statelessness is, in the
authors' words, "an illusion" (p. 5). KRM's answer is to stop pretending: state is explicit,
addressable, watchable, and split by ownership.

### 2. `watch` is an interaction style MAP does not model

MAP stays inside request/response throughout. *Conditional Request* — ETag, `If-None-Match`, `304
Not Modified` — optimises polling. It does not abolish it.

Kubernetes abolishes it.

```http
GET /api/v1/namespaces/shop/pods?watch=true
      &resourceVersion=48210
      &allowWatchBookmarks=true
      &sendInitialEvents=true
      &resourceVersionMatch=NotOlderThan
```

```json
{"type":"ADDED",   "object":{"kind":"Pod","metadata":{"name":"web-1","resourceVersion":"48211"}}}
{"type":"ADDED",   "object":{"kind":"Pod","metadata":{"name":"web-2","resourceVersion":"48212"}}}
{"type":"BOOKMARK","object":{"kind":"Pod","metadata":{"resourceVersion":"48212",
   "annotations":{"k8s.io/initial-events-end":"true"}}}}
{"type":"MODIFIED","object":{"kind":"Pod","metadata":{"name":"web-1","resourceVersion":"48260"}}}
{"type":"BOOKMARK","object":{"kind":"Pod","metadata":{"resourceVersion":"48294"}}}
```

Three details make this more than "long polling with extra steps":

- The `BOOKMARK` events carry **no object data**. They advance the client's `resourceVersion` so a
  reconnect resumes correctly. That is a "nothing new" signal costing **zero requests** — strictly
  cheaper than a `304`, which still needs a round trip.
- `sendInitialEvents=true` merges the initial full sync into the stream and terminates it with the
  annotated bookmark above. After that, a controller never issues a full list again.
- Missing an event is harmless, because reconciliation is level-triggered: the controller re-derives
  from observed state rather than replaying a log. No guaranteed delivery, no dead-letter queues, no
  idempotency keys.

### 3. *Two in Production*, taken far more seriously than the pattern requires

MAP's *Two in Production* says: run two versions side by side so clients can migrate. Crucially, the
pattern explicitly does **not** require the versions to be compatible. They are two contracts
offering similar functionality.

KRM requires the exact opposite, and enforces it.

Not two versions — arbitrarily many. `resource.k8s.io` today registers four in
`pkg/registry/resource/rest/storage_resource.go`: `v1`, `v1beta2`, `v1beta1` and `v1alpha3`.

And they are not two contracts — they are two views of one object. Here is a real field rename
across the two `flowcontrol` versions that were served in parallel in 1.26–1.28:

```yaml
apiVersion: flowcontrol.apiserver.k8s.io/v1beta2   # <- older view
kind: PriorityLevelConfiguration
metadata:
  name: workload-low
  uid: 9d5c1e3a-...-b21f                           # <- same uid
  resourceVersion: "77301"                         # <- same resourceVersion
spec:
  type: Limited
  limited:
    assuredConcurrencyShares: 30                   # <- the old field name
---
apiVersion: flowcontrol.apiserver.k8s.io/v1beta3   # <- newer view, SAME OBJECT
kind: PriorityLevelConfiguration
metadata:
  name: workload-low
  uid: 9d5c1e3a-...-b21f
  resourceVersion: "77301"
spec:
  type: Limited
  limited:
    nominalConcurrencyShares: 30                   # <- renamed, losslessly convertible
```

For CRDs the same guarantee is configured declaratively:

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
spec:
  versions:
    - name: v1alpha1
      served: true            # <- still answered
      storage: false
      deprecated: true
      deprecationWarning: "example.com/v1alpha1 Widget is deprecated; use v1"
      # ... schema
    - name: v1
      served: true
      storage: true           # <- exactly one version may be the storage version
      # ... schema
  conversion:
    strategy: Webhook         # <- or None, which only rewrites apiVersion
    webhook:
      conversionReviewVersions: ["v1"]
      clientConfig:
        service: { namespace: example-system, name: conversion-webhook, path: /convert }
```

For built-in types, conversion is hub-and-spoke through an internal version that is never
serialised (`runtime.APIVersionInternal`, the string `"__internal"`), turning O(n²) converters into
O(n). Losslessness is enforced by fuzzing-based round-trip tests; where a field has no home in an
older version, it gets parked in an annotation so the round trip still holds.

The bill comes due at design time: a field that cannot be losslessly mapped back onto the older
version is effectively unbuildable while that version is still served. This single constraint
explains why Kubernetes almost never bumps a major API version. It does not manage incompatibility —
it forbids it.

### 4. The description is executable

For MAP, an *API Description* is a document. For a CRD, the description **is part of the resource**,
and it is simultaneously the source of truth for validation, defaulting, pruning and merge
behaviour:

```yaml
schema:
  openAPIV3Schema:
    type: object
    properties:
      spec:
        type: object
        properties:
          replicas:
            type: integer
            minimum: 0
          containers:
            type: array
            x-kubernetes-list-type: map          # <- merge semantics for Server-Side Apply
            x-kubernetes-list-map-keys: ["name"]
            items:
              type: object
              # ... properties
        x-kubernetes-validations:
          - rule: "self.replicas <= self.maxReplicas"    # <- CEL invariant, enforced by the server
            message: "replicas must not exceed maxReplicas"
```

The `x-kubernetes-*` extensions carry exactly what MAP asks for under "dynamic and behavioural
aspects" and what OpenAPI cannot express on its own. What the description does *not* cover is the
commercial half of MAP's template: owner, support process, billing, SLA. Those live out of band, in
`OWNERS` files and policy documents. A recurring theme: Kubernetes nails the technical half of every
governance pattern and ignores the business half entirely.

### 5. Errors are objects, and there are two channels

`metav1.Status` is a `Kind` like any other:

```json
{
  "kind": "Status", "apiVersion": "v1",
  "status": "Failure",
  "message": "Deployment.apps \"web\" is invalid: spec.replicas: Invalid value: -1: must be >= 0",
  "reason": "Invalid",
  "details": {
    "group": "apps", "kind": "Deployment", "name": "web",
    "causes": [
      { "reason": "FieldValueInvalid",
        "message": "Invalid value: -1: must be >= 0",
        "field": "spec.replicas" }
    ]
  },
  "code": 422
}
```

Client code does not parse strings; it calls `apierrors.IsNotFound(err)`, `apierrors.IsConflict(err)`
and friends. One uniform error report across every resource, including CRDs nobody had written yet.

Beyond the pattern, KRM has a **second, asynchronous error channel**. A failure that only surfaces
during reconciliation has a defined place to live:

```yaml
status:
  containerStatuses:
    - name: app
      state:
        waiting:
          reason: CreateContainerConfigError      # <- the ConfigMap referenced in spec is missing
          message: 'configmap "app-config" not found'
```

Under MAP, such a failure has nowhere to go: the HTTP response was sent long ago.

---

## Where Kubernetes is genuinely worse

### 1. There is no read-side parsimony

This is the hard one. Three of MAP's five data-transfer-parsimony patterns are simply missing.

**No *Wish List*.** There is no sparse fieldset protocol:

```http
# Does not exist. There is no ?fields= parameter anywhere in the API.
GET /api/v1/namespaces/shop/pods?fields=metadata.name,status.phase
```

What exists instead are two hard-wired projections, negotiated via `Accept` parameters:

```http
GET /api/v1/namespaces/shop/pods
Accept: application/json;as=PartialObjectMetadataList;g=meta.k8s.io;v=v1
```
```yaml
apiVersion: meta.k8s.io/v1
kind: PartialObjectMetadataList        # <- TypeMeta + ObjectMeta only, no spec, no status
items:
  - metadata:
      name: web-1
      resourceVersion: "48211"
      # ... labels, annotations
```

```http
GET /api/v1/namespaces/shop/pods
Accept: application/json;as=Table;g=meta.k8s.io;v=v1
```
```yaml
kind: Table
columnDefinitions:                     # <- columns chosen by the PROVIDER, not the client
  - {name: Name, type: string}
  - {name: Ready, type: string}
  - {name: Status, type: string}
rows:
  - cells: ["web-1", "1/1", "Running"]
    # ... object, per TableOptions.includeObject: None | Metadata | Object
```

Worth stating plainly, because it trips people up: `kubectl get -o jsonpath=…` and
`-o custom-columns=…` are **client-side**. They save exactly zero bytes on the wire.

**No *Wish Template*.** Even less. The structurally closest thing is on the *write* side —
Server-Side Apply's `managedFields`:

```json
"managedFields": [{
  "manager": "kube-controller-manager", "operation": "Update", "apiVersion": "v1",
  "fieldsType": "FieldsV1",
  "fieldsV1": {
    "f:spec": {
      "f:containers": {
        "k:{\"name\":\"nginx\"}": { "f:image": {}, "f:resources": { "f:requests": { "f:cpu": {} } } }
      }
    }
  }
}]
```

That tree — object structure mirrored with `f:`-prefixed keys and empty leaves — is formally a wish
template. But it expresses *ownership of written fields*, not *interest in read fields*, and no
operation accepts it as a read selector.

**No *Request Bundle*.** No batch endpoint, no multi-object transaction, no atomicity across
objects:

```console
$ kubectl apply -f manifests/          # 50 objects
# -> 50 separate HTTP requests; ApplyOptions.Run iterates and calls applyOneObject per object.
# Failures are partial. There is no rollback.
```

`deletecollection` is the only genuine collection write, and it is one operation over a selected
set, not n bundled operations:

```http
DELETE /api/v1/namespaces/shop/pods?labelSelector=app%3Dweb
```

Two places do implement bundling, both in the meta layer — and, tellingly, the discovery one is also
one of the very few endpoints in the whole apiserver that supports ETags:

```http
GET /apis
Accept: application/json;g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList
If-None-Match: "a1b2c3d4"
# -> the entire catalogue in one request, or 304 Not Modified
```

**Why it hurts less than you would think.** The expensive path in Kubernetes is not writing many
objects once; it is watching continuously — and that path is already bundled into a single
long-running request. Clients that would benefit from response shaping mostly hold informer caches
of complete objects instead. But where informers do not fit — short-lived CLI tools, web UIs,
one-shot jobs — you pay the full price with no mitigation available.

There is also a structural reason the gap will not close easily. A per-request partial object is not
a valid object in this model: you could not distinguish "field omitted from the response" from
"field genuinely empty", which makes it useless for level-triggered reconciliation, and cache
entries would differ per requesting controller.

### 2. Operations get bent out of shape

Where a transition really is imperative and cannot be phrased as a desired state, KRM has to distort
the model. A `SubjectAccessReview` is an RPC wearing a resource costume — you `create` an object and
get it back with `status` filled in, and nothing is ever stored:

```yaml
# Request: POST /apis/authorization.k8s.io/v1/subjectaccessreviews
apiVersion: authorization.k8s.io/v1
kind: SubjectAccessReview
spec:
  user: alice
  resourceAttributes: { namespace: shop, verb: delete, group: apps, resource: deployments }
```
```yaml
# Response — same kind, status filled, nothing persisted. No get, no list, no watch exist.
apiVersion: authorization.k8s.io/v1
kind: SubjectAccessReview
spec: { ... }
status:
  allowed: false
  reason: 'RBAC: no rules allow user "alice" to delete deployments in namespace "shop"'
```

The other contortions: `create` verbs that create nothing (`pods/exec`, `pods/attach`,
`pods/portforward` abandon CRUD entirely and upgrade the connection to a stream), and write-only
commands POSTed to subresources (`pods/binding`, `pods/eviction`, `namespaces/finalize`,
`serviceaccounts/token`, `certificatesigningrequests/approval`).

Compare what the same intent looks like when it *can* be phrased declaratively:

```http
# Does not exist. There is no verb for business transitions.
POST /apis/apps/v1/namespaces/shop/deployments/web/restart
```
```yaml
# Instead: change the desired state; the controller converges. Idempotent by construction.
spec:
  template:
    metadata:
      annotations:
        kubectl.kubernetes.io/restartedAt: "2026-08-04T10:00:00Z"
```

Each workaround works. Collectively they are the sound of a uniform model being levered open.

### 3. There is no referential integrity

References are name strings, not URIs — deliberately, so a manifest is portable across clusters:

```yaml
spec:
  serviceAccountName: web                       # <- a name, resolved by convention
  volumes:
    - name: config
      configMap:
        name: app-config                        # <- may not exist. Admission will not object.
    - name: data
      persistentVolumeClaim:
        claimName: web-data
  containers:
    - name: app
      env:
        - name: API_TOKEN
          valueFrom:
            secretKeyRef: { name: api-creds, key: token }
```

Nothing validates any of these at write time; the failure surfaces later, level-triggered, as the
`CreateContainerConfigError` shown earlier. The consequences compound: renaming is effectively
delete-and-recreate, and "who references me?" cannot be answered without scanning every potential
referrer.

`metadata.ownerReferences` fixes integrity for *ownership* relations — and only those:

```yaml
metadata:
  ownerReferences:
    - apiVersion: apps/v1
      kind: ReplicaSet
      name: web-7d9f
      uid: 3f1c...-88ab          # <- uid guards against name reuse
      controller: true
      blockOwnerDeletion: true   # <- drives garbage collection ordering
```

This is a deliberate choice: validating references would turn every write into a distributed
two-phase problem. It is still a cost.

### 4. Mistakes are permanent

The flip side of the de facto *Eternal Lifetime Guarantee* for `v1`. In a GA version, deprecation is
advice, nothing more, because *Two in Production*'s round-trip requirement forbids removal:

```yaml
spec:
  serviceAccountName: web
  serviceAccount: web          # <- deprecated alias; Go type field is DeprecatedServiceAccount.
                               #    Still in the v1 PodSpec. Will never be removed.
  volumes:
    - name: legacy
      # ... gitRepo, glusterfs, rbd, flexVolume, cinder — in-tree types whose
      #     implementations were removed years ago. The fields remain.
```

What makes this sharper is the **ordering problem**. In MAP's own tutorial, switching from *Embedded
Entity* to *Linked Information Holder* is step 4 — a normal improvement you make after measuring. In
KRM it is a breaking schema change requiring a new group-version with conversion. The decision has
to be made before graduation, not after the first performance complaint. Alpha exists precisely for
this, which is good design — but the pressure to get structure right up front is far higher than the
pattern language assumes.

### 5. Discovery leaks the entire schema

A small finding, but a real violation. MAP's *Community API* requires sharing the API description
only with the intended audience. In Kubernetes:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:discovery
rules:
  - verbs: ["get"]
    nonResourceURLs:
      - /livez
      - /readyz
      - /healthz
      - /version
      - /version/
      - /openapi
      - /openapi/*        # <- the complete OpenAPI schema of every group
      - /api
      - /api/*
      - /apis
      - /apis/*           # <- including every CRD in the cluster
```
```yaml
kind: ClusterRoleBinding
metadata:
  name: system:discovery
roleRef: { kind: ClusterRole, name: system:discovery }
subjects:
  - kind: Group
    name: system:authenticated     # <- anyone who can log in, with no permissions on anything
```

In a multi-tenant cluster, tenant A learns exactly which operators tenant B is running. There is a
good reason for it — discovery and version negotiation have to work before authorization can be
meaningfully evaluated — but it is an information disclosure the pattern language explicitly rules
out, and it is worth knowing you have made that trade.

---

## The deeper trade: decisions made once

Step back from the individual patterns and the shape of the thing becomes clear.

MAP distributes design authority across three independent levels: which endpoints exist, which
operations each offers, and how messages are structured. KRM standardises the middle level out of
existence. What remains is: which resources exist, and what their schemas look like.

This has a measurable consequence. Walk MAP's cheat sheet top to bottom and most rows are
pre-answered:

| MAP design issue | KRM's answer |
|---|---|
| Visibility | Not an API property — RBAC and network |
| Endpoint role | Information Holder Resource with `spec`/`status`; anything else needs a strong justification |
| Operation responsibility | Not selectable — the verb set is fixed |
| Data contract | Always a Parameter Tree; OpenAPI v3 in the CRD |
| Identification | `metadata.name`, `metadata.namespace`, `metadata.uid` |
| Hypermedia | Not applicable — `ObjectReference` plus discovery |
| Error report | `metav1.Status` |
| Rate limiting | API Priority and Fairness |
| Versioning | Group-version in path *and* payload; maturity encoded in the name |
| Multiple versions | All served versions, losslessly converting |
| Deprecation | Policy, `Warning` header, removal at the announced release |

Two of those deserve a closer look, because they show how thoroughly Kubernetes has industrialised
what MAP treats as a per-API decision.

**Rate limiting** is not a per-consumer quota at all. API Priority and Fairness limits *concurrency*
— seats — and distributes damage by shuffle sharding rather than rationing calls:

```yaml
apiVersion: flowcontrol.apiserver.k8s.io/v1
kind: FlowSchema
metadata: { name: workload-low }
spec:
  matchingPrecedence: 900
  priorityLevelConfiguration: { name: workload-low }
  distinguisherMethod: { type: ByUser }     # <- what counts as "one client"
  rules:
    - subjects: [{ kind: ServiceAccount, serviceAccount: { namespace: shop, name: "*" } }]
      resourceRules:
        - verbs: ["*"]
          apiGroups: ["*"]
          resources: ["*"]
          namespaces: ["*"]
---
apiVersion: flowcontrol.apiserver.k8s.io/v1
kind: PriorityLevelConfiguration
metadata: { name: workload-low }
spec:
  type: Limited
  limited:
    nominalConcurrencyShares: 30            # <- share of total server concurrency, not calls/hour
    limitResponse:
      type: Queue                           # <- wait, don't reject
      queuing:
        queues: 64
        handSize: 6                         # <- shuffle sharding: one heavy client hits 6 of 64
        queueLengthLimit: 50
```

Only `type: Reject` and queue overflow produce a 429:

```http
HTTP/1.1 429 Too Many Requests
Retry-After: 1
X-Kubernetes-PF-FlowSchema-UID: 9c2f...
X-Kubernetes-PF-PriorityLevel-UID: 4a71...
```

**Deprecation** travels in-band, in every single response:

```http
HTTP/1.1 200 OK
Warning: 299 - "storage.k8s.io/v1beta1 VolumeAttributesClass is deprecated in v1.34+,
  unavailable in v1.37+; use storage.k8s.io/v1 VolumeAttributesClass"
```

Deprecated-since, removed-in and replacement — all three in every single response, assembled by
`deprecation.WarningMessage` from generated lifecycle methods. `kubectl` and every client-go client
print it to stderr, so it shows up in terminals and CI logs rather than only in release notes.

…and the provider can measure who is still calling:

```
apiserver_requested_deprecated_apis{group="storage.k8s.io",version="v1beta1",
  resource="volumeattributesclasses",subresource="",removed_release="1.37"} 1
```

So what is left for the API designer? Two rows: **what is a `Kind`**, and **what does its schema look
like**.

The first is genuinely hard and MAP does not help much — the best available framing comes from
Domain-Driven Design, which the pattern language references but does not own. A `Kind` corresponds
most closely to an Aggregate, and Kubernetes' per-object atomicity guarantee is precisely the
Aggregate consistency boundary. The `Published Language` concept describes a versioned API group
more accurately than anything in MAP itself: a published exchange vocabulary that binds neither side
internally. The internal Go types are explicitly not part of it and are never serialised.

The second — schema design — is where all remaining effort goes: embed or reference, list merge
semantics, optionality and defaulting, and above all whether every planned extension can be
losslessly converted back to today's version.

That is the whole trade. Kubernetes traded the ability to design each endpoint well for the ability
to build tooling once. Given that CRDs let anyone add resources at runtime, and that the tooling has
to work on resources it has never seen, it is hard to argue the trade was wrong.

---

## What each side could borrow

**What API designers could take from KRM:**

- **Split the representation by write authority, not by direction.** `spec`/`status` with separate
  authorization beats "request DTO and response DTO" for anything with an asynchronous back end.
  Server-Side Apply's per-field ownership generalises it to n writers.
- **Put the version in the payload, not just the path.** Self-describing objects make manifests,
  streams and mixed documents tractable — a file stays interpretable with no knowledge of how it was
  transported.
- **Make deprecation travel in-band**, with a server-side metric so the provider knows who is still
  calling. MAP puts the end date in a document; Kubernetes puts it in every response and in a time
  series.
- **Consider whether you need polling at all.** A watch stream with cost-free "nothing new" signals
  beats any conditional-request scheme, and level-triggered semantics make delivery guarantees
  unnecessary.

**What Kubernetes could take from MAP:**

- **A read-side projection mechanism.** Not GraphQL — something narrow and cacheable, closer to
  *Wish Template*, and honest about the "omitted vs. empty" problem. `fieldsV1` already proves the
  encoding exists.
- **The vocabulary for talking about data types.** MAP's Operational / Master / Reference Data
  Holder distinction exists in Kubernetes in practice — `Lease` and `Event` behave one way, `Node`
  and `Namespace` another, `StorageClass` a third — but there is no name for it, so the reasoning
  gets rediscovered per KEP. `Lease` was introduced specifically to pull high-frequency heartbeats
  out of the large, heavily-referenced `Node` object:

  ```yaml
  apiVersion: coordination.k8s.io/v1
  kind: Lease
  metadata: { name: node-1, namespace: kube-node-lease }
  spec:
    holderIdentity: node-1
    leaseDurationSeconds: 40
    renewTime: "2026-08-04T09:14:02.113Z"   # <- rewritten every ~10s, referenced by nobody
  ```

  That is a textbook operational/master data split, argued from first principles because the
  vocabulary was not there.
- **Forces and consequences written down.** MAP's format — forces as questions, consequences with
  explicit `+`/`−` lines — would improve KEPs considerably. (Ironically, the pattern website itself
  now truncates these, referring readers to the book. The original EuroPLoP papers still have them
  in full, and they are the best part.)

---

## Coda

The most striking thing about running this comparison is how rarely Kubernetes is simply *missing*
something. Nine patterns are absent, and of those, four are non-goals for an infrastructure control
plane (*Public API*, *Pricing Plan*, and arguably *Context Representation* and *Semantic
Versioning*). The genuine gaps cluster in exactly one place: letting the client shape what comes
back.

That is not an accident or an oversight. It is the direct consequence of the founding decision — one
schema, uniform across every resource, so that generic machinery works on resources nobody has seen
yet. Response shaping is the one thing that decision makes structurally impossible.

Kubernetes is an unusually opinionated API. Measured against a pattern language built for a world
where each endpoint gets designed on its own terms, it looks less like a violation of good practice
and more like a bet: that uniformity compounds, and expressiveness does not.

---

*Sources: the [MAP wiki](../README.md) in this repository summarises all 45 patterns with a
Kubernetes cross-reference each, plus the seven underlying papers. The full comparison is in
[KRM vs. MAP](../KRM-vs-MAP.md). Pattern language: <https://microservice-api-patterns.org/>,
Zimmermann, Stocker, Lübke, Pautasso and Zdun.*
