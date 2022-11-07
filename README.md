![](https://avatars.githubusercontent.com/k8s-social)

# k8s.social - Infra

Infrastructure automation for [k8s.social](https://k8s.social)

> Infrastructure kindly provided by [Civo](https://civo.com). 💙 Thank you Civo!

## Includes

- Kubernetes cluster
  - CNI = cilium
  - Ingress = Nginx
  - Gitops = ArgoCD
  - Secrets = SealedSecrets
- Firewall
- S3 compatible object store bucket
- DNS

## Info

The aim of this repo is to bootstrap the initial infrastructure enough to then pass control over to [k8s-social/gitops](https://github.com/k8s-social/gitops) for handling of applications.

Infrastructure is handled using Pulumi (in Go).

All (approved) pull requests will perform a `pulumi preview` and display the resulting changes as a comment on the PR.

Merges to the `main` branch result in `pulumi up` being run to apply all infrastructure changes.
