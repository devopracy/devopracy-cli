# Devopracy CLI

[![Build Status][travis-badge]][travis]

[travis-badge]: https://travis-ci.org/devopracy/devopracy-cli.svg?branch=master
[travis]: https://travis-ci.org/devopracy/devopracy-cli

The devopracy command line interface is a meta-application for building and deploying prefabricated infrastructure on public cloud platforms. It is meant to enable fast, secure, consistent end-to-end environment builds with a minimum of configuration.This is the third version of the CLI which was originally built at Prime 8 Consulting for exploratory testing and small cloud development. The CLI comes with a handful of buildfiles and modules that can be built and deployed with a single configuration file. The base cloud contains five services: buildserver, vpn, secretserver, monitor, dns. From this basic setup, a team can have an insta-environment for secure development. The environment is customizable by adding new builds and modules to the infrastructure model.

Built on Hashicorp's Packer and Terraform, the devopracy-cli was built as in-house tooling for an innovation consultancy that ran many small environments for various prototyping projects. The advantage of many small pipelines is that a project can have a small dedicated buildserver allowing engineers to self serve. Furthermore the entire generated code for the infrastructure can be handed off to a client and reproduced exactly. This is useful for both disaster recovery and cloud defense. The model for generating small clouds in a touchless paradigm was invented for online deliberation and cloud-based election systems (online voting) to support technologically transparent systems for citizens to interact with government. 
