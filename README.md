# nagios-heartbeat-plugin

OpsGenie heartbeat plugin that can be used in Nagios; with Linux and Windows distributions.

1- Please [create an OpsGenie account](https://www.opsgenie.com/#signup) if you haven't done already.

2- Add a heartbeat in OpsGenie [Heartbeats page](https://opsgenie.com/heartbeat).

3- Create a Heartbeat from [Heartbeats](https://app.opsgenie.com/heartbeat) page. 

4- Copy api key of Default API integration.

5- Download the binary file, _heartbeat_ (provided below), into your Nagios libexec directory.

6- In Nagios, define a command like so: 

```
define command{
    command_name    opsgenie_heartbeat
    command_line    /usr/local/nagios/libexec/heartbeat -apiKey $ARG1$ -name $ARG2$
}
```

7- Define a service that will run the command like so:

```
define service {
    service_description     OpsGenie Heartbeat
    host_name               localhost
    check_interval          10
    check_period            24x7
    max_check_attempts      60
    retry_interval          1
    notification_interval   60
    check_command           opsgenie_heartbeat!API_KEY!HEARTBEAT_NAME
}
```

where API_KEY is the api key you acquired from OpsGenie integration, and HEARTBEAT_NAME is the name of the heartbeat you added.

##### Troubleshooting

Make sure the binary file _heartbeat_ is executable.
