# Traefik Plugin Init

This repository provides a simple container that:

- Can be used as an `initContainer` in Kubernetes
- Avoids having to build custom `traefik` docker images to use [traefik plugins](https://doc.traefik.io/traefik-pilot/plugins/overview/)
    - (Specifically, [local plugins](https://github.com/traefik/plugindemo#local-mode), which do not require traefik pilot)
- Can be used with the [traefik helm chart](https://github.com/traefik/traefik-helm-chart)
- Allows multple plugins for a single traefik installation
- Allows specifying GitHub references per plugin (i.e. branch, tag, ref, etc.)

Please see the [./examples](./examples) directory for illustrations of how to use

## What does it do

On startup, the container loops through GitHub repositories and clones them
into the container. By mounting a shared volume between this container and the
traefik container, this will make the plugin sources available to traefik.

## Configuration

Recommended configuration is with environment variables:
- `TRAEFIK_PLUGIN_PATH`: the path repositories will be cloned within
- `TRAEFIK_PLUGIN_PREFIX`: the environment variable prefix to use to define GitHub repositories
- `TRAEFIK_PLUGIN_REPO_1`, `TRAEFIK_PLUGIN_REPO_2`, etc.: if
  `TRAEFIK_PLUGIN_PREFIX` = `TRAEFIK_PLUGIN_REPO_`, then these wouldl house
  your GitHub references

Git URL examples:
- `https://github.com/XciD/traefik-plugin-rewrite-headers` (will reference the
  default branch at container startup)
- `https://github.com/XciD/traefik-plugin-rewrite-headers@main` (chooses a
  branch or tag ref)
- `https://github.com/XciD/traefik-plugin-rewrite-headers@e66f2dc` (chooses a
  sha)

**NOTE:** in production, it is recommended to pin plugins to an immutable
reference (usually a `sha` or `tag`). Otherwise, restarting traefik pods can
change plugin versions

## Shortcomings

- Tested only against https://github.com . Semantics may be different for other
  git repository providers
