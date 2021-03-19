<!--
Copyright The Shipwright Contributors

SPDX-License-Identifier: Apache-2.0
-->

---
title: deprecate-runtime-image
authors:
  - "@ImJasonH"
reviewers:
  - "@otaviof"
  - "@zhangtbj"
approvers:
  - "@sbose78"
  - "@qu1queee"
creation-date: 2021-03-19
last-updated: 2021-03-19
status: provisional
see-also:
  - "/docs/proposals/runtime-image.md"
---

# Deprecate and remove `runtime` image support

## Release Signoff Checklist

- [ ] Enhancement is `implementable`
- [ ] Design details are appropriately documented from clear requirements
- [ ] Test plan is defined
- [ ] Graduation criteria for dev preview, tech preview, GA
- [ ] User-facing documentation is created in [docs](/docs/)

## Open Questions [optional]

- Are users relying on `runtime` in Shipwright today, that would be negatively impacted by its deprecation and removal?

- Will the users relying on `runtime` in OpenShift Builds v1 today have a good experience if they want to migrate to Shipwright?

## Summary

Support for `runtime` image was proposed and added in https://github.com/shipwright-io/build/pull/263, in June/July 2020.
Its main goal was to allow users to produce _lean images_, by pulling files from a previously built image and pulling them into an image produced by Shipwright.

This feature was inspired by [similar support in OpenShift Builds (v1)](https://www.openshift.com/blog/chaining-builds), which has supported this kind of feature since before Docker supported [multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/).
Since OpenShift Builds v1 introduced this support, there have been many improvements in the field of producing minimal images, including multi-stage Docker builds, buildpacks, and many others.

Multi-stage builds and buildpacks in particular seem to be the dominant mechanisms currently for producing minimal images. At the same time, `runtime` support in OpenShift Builds is OpenShift-specific and less widely used.

In Shipwright, `runtime` enables users to make relatively invasive changes to the input image, not only limited to changing the base image or copying out files, but also running new steps, which can have security considerations.

## Motivation

As we aim to produce a small, powerful, flexible and usable Shipwright API, we should actively examine areas of the API that are not "pulling their weight" and can be removed. Especially while the project is not in broad usage, we should take advantage of that fact and deprecate anything that isn't contributing to an ideal minimal API surface.

Potential security concerns around the `run` feature lead me to bias toward removing the feature before more users come to rely on it, which could lead to a more burdensome deprecation for contributors, and more users to impact.

If we find that users are not relying on this feature today, we should bias toward removing it until such a time as we believe users would benefit from its addition, at which point we can assess re-adding the feature based on user need at that time.

### Goals

- Responsibly deprecate and remove the `runtime` feature, giving users ample time to understand the coming change, and ample opportunity to push back on it.

### Non-Goals

- Replacing the functionality with a different Shipwright API surface directly

What is out of scope for this proposal? Listing non-goals helps to focus discussion and make
progress.

## Proposal

- Document that `runtime` is deprecated, and discourage new usage in docs.
- Notify shipwright-users@ of the deprecation and plan for removal.
- Include deprecation in the next release's release notes.
- After some time, remove support by deleting the field in build_types.go, any code that depends on it, any tests that exercise the code, and any documentation that describes it (except in docs/proposals/, for posterity).

### User Stories [optional]

#### New Shipwright User

As a user interested in trying Shipwright, I'd like to be able to build lean images. I should be able to find documentation about how to achieve this using existing technologies outside of the Shipwright project.

#### Existing Shipwright User Using `runtime`

As an existing Shipwright user whose workflow depends on the `runtime` image, I should be able to find documentation from Shipwright about how to achieve this, and should find help in the Shipwright mailing list and Slack channel to help me succeed.

### Risks and Mitigations

If users are heavily relying on `runtime`, we don't want to remove it and cause them pain. We can delay the removal and help them migrate to another solution, or we can decide that this signal means the feature should not be removed, and un-deprecate it.

If after removing the feature users tell us they'd like the feature back (or something like it), we can engage with them to re-add support based on their needs and their feedback. The feature we build might look exactly like `runtime` today, or it might not.

## Design Details

### Test Plan

As this is a strict code removal, tests of all kinds should remain in place during deprecation period, and can be deleted when the code they exercise is deleted.

### Graduation Criteria

<TODO>

Many sections of the proposal template have been removed as they cover gradation of a feature from Dev Preview -> Tech Preview -> GA, which is not a concern for this deprecation and removal.

### Upgrade / Downgrade Strategy

If a user relies on `runtime`, they can continue to run the latest version of Shipwright supports it for as long as they want. Given Shipwright's status, we already make no supportability or security patching promises, but they should consider migrating to a later version without `runtime` as soon as they can, potentially with our help if needed. Getting good signal that the feature is unused will be helpful in mitigating this concern (see [Risks and Mitigations](#risks-and-mitigations)).

### Version Skew Strategy

A Build CR might be stored on a cluster in etcd with the `runtime` field populated, and we need to make sure the controller doesn't panic when it sees this. Kubernetes' 

## Implementation History

_Major milestones in the life cycle of a proposal should be tracked in `Implementation History`._

## Drawbacks

Users might depend on this and we might not find out until after we've removed support, inflicting pain on users and churn on contributors while we figure out how/whether to re-add support.

## Alternatives

Continue to support `runtime` indefinitely, taking on future support burden including potential security concerns.

Remove support for more problematic features of `runtime` (e.g., `run`) and not others (e.g., `copy`). Deciding which features to cut and which to keep might be more complex than just removing the whole feature and re-adding it (conceptually) at a later date with more user input.
