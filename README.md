## [Prometheus](https://prometheus.io) High Availability proxy

> even though application is pretty stable, keep in mind it is still in active development, we use it in production, but it doesn't meen it is ready for you. 

### The Problem
Normal Prometheus upgrade procedure looks like this:
1. Switch traffic to fallback server (automatic with HA proxy, or manually)
2. Upgrade primary server
3. Wait for it to become available
4. Do 1-3 steps for secondary server

If you did everything right, it doesn't mean you will not have gaps on Grafana. 

### The Solution
In order to prevent gaps on Grafana, we query the same request from several servers. The resulting responses are merged into one unified response using, pretty naive, json merging techniques. 

```
docker run -d -p 9090:9090 anchorfree/prometheus-ha-proxy:master http://server1:9090 http://server2:9090
```

### what this project is
This project is simple and pretty dumb Prometheus output merger. With this project you can:
1. Use one endpoint in Grafana instead of several servers
2. Have truly highly available Prometheus setup for Grafana
3. Startup this proxy, and do upgrades/restarts without any problems. 

### what this project is NOT
1. It doesn't do any way of values aggregations over datasets, e.g. you can not calculate total bandwidth for all of your datacenters, you need [Federation Server](https://prometheus.io/docs/operating/federation/) for this. 
2. It doesn't solve any HA problems for AlertManager or any querying issues. 

### current state and feature plans
We would like to make a proof of concept and share the answer to Grafana drops problem we had with Prometheus at [Anchorfree](https://www.anchorfree.com). 

Evolution of this tool would be:
- Intergare Prometheus WEB GUI, in order to make look and feel like real Prometheus (currently only Grafana and curl queries are possible)
- Use remote reades + local TSDB storage to provide single point of contact for our Ops team
