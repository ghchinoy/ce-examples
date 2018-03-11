# eventson

Using the existing cectl toml file (`~/.config/ce/cectl.toml`), go through all profiles available and display a table of Element Instances and what types of events are enabled.

A sample output is below.

There are three flags:
* `--filter` - this is a profile filter, only the profiles that match this flag will be displayed, ex. `--filter prod`
* `--element` - only the elements specified in this flag will be displayed, ex `--element sfdc` will show only the Salesforce Element Instances
* `--disable` - this boolean flag determines whether the Instances will then be subject to having their events disabled after reporting, see an example below

With the output, one can the use `cectl` to enable or disable Element Instance events (`cectl instances events-enable <id> [true|false]`) or disable the Element Instance itself (`cectl instances disable <id>`) as needed.


```
$ eventson
2018/03/10 11:14:59 Querying 86 profiles
+-------------------+-----------------------+--------+--------------------------------+----------+-----------+----------+
|      PROFILE      |        ELEMENT        |   ID   |              NAME              | DISABLED | EVENTTYPE | INTERVAL |
+-------------------+-----------------------+--------+--------------------------------+----------+-----------+----------+
|              2080 | sqlserver             |  67375 | Tuesday                        | false    | polling   |        1 |
+                   +                       +--------+--------------------------------+          +           +          +
|                   |                       |  64578 | TwentyEighty MS SQL            |          |           |          |
+                   +                       +--------+--------------------------------+          +           +----------+
|                   |                       |  65109 | Jbond SQL Test                 |          |           |       15 |
+-------------------+-----------------------+--------+--------------------------------+          +-----------+----------+
| default           | connectwisehd         |  91549 | ConnectWise Dev                |          |           |          |
+                   +-----------------------+--------+--------------------------------+          +-----------+----------+
|                   | slack                 |  75552 | Slack Message Receiver         |          | webhooks  |          |
+-------------------+-----------------------+--------+--------------------------------+----------+-----------+----------+
| demo              | box                   |   8694 | Sales Demo Account with events | false    |           |          |
|                   |                       |        |                              3 |          |           |          |
+-------------------+-----------------------+--------+--------------------------------+----------+-----------+----------+
...
```

This example shows all flags in effect, disabling events after reporting:

```
$ eventson --filter prod --element sfdc --disable
2018/03/11 11:11:57 Querying 5/86 profiles: prod
2018/03/11 11:11:57 Filtering by Element key: sfdc
2018/03/11 11:12:09 prod-apitester instances Non-200 Status: 404
2018/03/11 11:12:10 prod-uk instances Non-200 Status: 404
+---------+---------+--------+---------------------------+----------+-----------+----------+
| PROFILE | ELEMENT |   ID   |           NAME            | DISABLED | EVENTTYPE | INTERVAL |
+---------+---------+--------+---------------------------+----------+-----------+----------+
| prod    | sfdc    | 452323 | STRIDE_SFDC_1509685135052 | true     | polling   |        5 |
+---------+---------+--------+---------------------------+----------+-----------+----------+
2018/03/11 11:12:10 Disabling instance id 452323 from profile prod
2018/03/11 11:12:15 Disabled events on sfdc/STRIDE_SFDC_1509685135052
```