# Dynatrace Extension


This image downloads and deploys an OneAgent extension

It is not very helpful on its own, but it is used by the [Dynatrace Extensions Operator](https://github.com/dlopes7/dynatrace-extensions-operator) to deploy extensions on kubernetes nodes where the user does not have access to the host.

To run with docker:

```shell
docker run --rm \
-e "DT_EXTENSION_NAME=rabbitmq" \
-e "DT_EXTENSION_LINK=http://my.server/custom.python.rabbitmq.zip" \
-v /opt/dynatrace/oneagent/plugin_deployment:/plugin_deployment \
quay.io/dlopes7/dt-extension
```

To use in kubernetes please check the [operator docs](https://github.com/dlopes7/dynatrace-extensions-operator).
