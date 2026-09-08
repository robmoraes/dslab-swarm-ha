# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
Versions identify published revisions of the supporting material for the
Distributed Systems Lab series.

## [Unreleased]

## [1.0.0] - 2026-09-08

### Added

- Added the complete supporting material for season 1 of the Docker Swarm HA
  cluster on AWS series, from the initial AWS and EC2 setup through public HTTPS
  exposure.
- Added an Ubuntu bootstrap script for installing Docker Engine, the CLI,
  Buildx and the Compose plugin.
- Added standalone container examples and a Go diagnostic API for inspecting
  requests, runtime, host, container, network and EC2 instance metadata.
- Added Docker Swarm stack examples for the diagnostic API and an OWASP
  ModSecurity WAF service.
- Added Traefik edge-routing configuration with shared overlay networking,
  DNS-based routing and automatic TLS certificates through Let's Encrypt.
- Added architecture diagrams and operational notes used throughout the video
  series.

[Unreleased]: https://github.com/robmoraes/dslab-swarm-ha/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/robmoraes/dslab-swarm-ha/releases/tag/v1.0.0
