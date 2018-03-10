# eventson

Using the existing cectl toml file (`~/.config/ce/cectl.toml`), go through all profiles available and display a table of Element Instances and what types of events are enabled.

A sample output is below.

There are two flags:
* `--filter` - this is a profile filter, only the profiles that match this flag will be displayed, ex. `--filter prod`
* `--element` - only the elements specified in this flag will be displayed, ex `--element sfdc` will show only the Salesforce Element Instances

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