# Logs

## TiUP playground

- [Install TiUP locally](https://docs.pingcap.com/tidb/stable/tiup-overview/#install-tiup)

```shell
# Start a local TiDB cluster with version 7
tiup playground ^7
```

- Confirm default username and password ref. [What is the default username and password in tiup playground? #60983](https://github.com/orgs/pingcap/discussions/60983)
  - TiDB Dashboard: http://127.0.0.1:2379/dashboard (default username: `root`, password: ``)
  - Grafana: http://127.0.0.1:3000 (default username: `admin`, password: `admin`)
- Connect to the TiDB cluster using MySQL client:
  - Install mysql client
    - on macOS: `brew install mysql-client` ref. https://formulae.brew.sh/formula/mysql-client

```shell
mysql --comments --host 127.0.0.1 --port 4000 -u root
```

## Kubernetes

- [Deploy TiDB Operator on Kubernetes](https://docs.pingcap.com/tidb-in-kubernetes/stable/deploy-tidb-operator/)
- [TiDB実践入門 サンプルコード](https://github.com/makocchi-git/tidb-practical-book-sample)
