## [Prometheus](https://prometheus.io) High Availability proxy

### it is currently in development, please use with caution 

In order to prevent gaps on Grafana, we query the same query from different servers. The resulting queries are merged into one response. 

```
docker run -d -p 9090:9090 anchorfree/prometheus-ha-proxy:master http://server1:9090 server2:9090
```
