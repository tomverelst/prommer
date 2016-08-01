FROM scratch

ADD prommer /bin/prommer

VOLUME /etc/prometheus

ENTRYPOINT ["/bin/prommer"]
CMD -target-file=/etc/prometheus/target-groups.json \
 -monitoring-label=prometheus-target
