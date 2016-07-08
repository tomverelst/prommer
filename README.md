# Prommer

**A silly contraction between Prometheus and Docker**

Prommer is a simple target discovery service for Prometheus using Docker.
It listens to Docker events to detect changes.
Only containers that are labeled with a certain configurable label will be used as a Prometheus target.
Prommer automatically writes these labeled containers as target groups to a JSON file.
This JSON file can be used for Prometheus' scrape configuration for dynamic target discovery.

**Prommer (currently) is only for a single host and not really tested.
Definitely do not use this in production.**

Why did I create this then, you ask?
I wanted a very simple solution for service discovery for Prometheus demo purposes on a single host.
I did not want to distract users with other complex service discovery setups like etcd and service registration agents.
I also wanted to use this learn Go,
so this is my first Go application.
Should you dare to look inside the source code,
be warned!

# Configuration

Starting Prommer is easy.

```
$ prommer [-target-file <target-file>] [-monitoring-label <label-name>]
```

Prometheus must be configured to use the target groups JSON file.
You must set the `file_sd_configs` field for the Prometheus scrape job to use the target groups

```yaml
scrape_configs:
  - job_name: 'prommer-job'
    scrape_interval: 5s
    file_sd_configs:
    - names: ['tgroups/*.json']
```

Prommer must be able to listen to the Docker events stream.
Therefor it is required to access the Docker daemon.
If you use Prommer inside a container (which you should),
this container must run with the `--privileged` flag.

# TODO

- [ ] Create and publish minimal Docker image that includes Prometheus and Prommer.
- [ ] Make it work for Docker Swarm so it might actually be useful.
- [ ] Optimize target update mechanism by doing incremental changes using the Docker events stream, instead of requesting all the labeled containers on each change.
- [ ] Become a Go ninja (unlikely) and probably rewrite this whole thing
