# Traefik Plugin Init

This repository provides a simple container that:

- Can be used as an `initContainer` in Kubernetes
- Avoids having to build custom `traefik` docker images to use [traefik plugins](https://doc.traefik.io/traefik-pilot/plugins/overview/)
    - (Specifically, [local plugins](https://github.com/traefik/plugindemo#local-mode), which do not require traefik pilot)
- Can be used with the [traefik helm chart](https://github.com/traefik/traefik-helm-chart)
- Allows multple plugins for a single traefik installation
- Allows specifying GitHub references per plugin (i.e. branch, tag, ref, etc.)

Please see the [./examples](./examples) directory for illustrations of how to use

## Shortcomings

- Tested only against https://github.com . Semantics may be different for other git repository providers
