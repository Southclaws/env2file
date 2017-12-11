# env2file

A tiny utility for turning environment variable values into files.

```sh
go get github.com/Southclaws/env2file
```

```sh
docker pull southclaws/env2file
```

## Why

Sometimes you need to deploy a containerised application and it requires a file for configuration. That means you have
to think about where that file goes, how to automate its creation and maybe even version control so you can record
configuration change history.

This simplifies that by allowing you to store small-ish configuration files as environment variables in your container
management system.

It's kind of like ConfigMaps in Kubernetes.

## How

This is mainly designed for deployments that use Docker Compose. The intended usage is:

1. Add env2file to your compose config
2. Mount a volume/directory to both env2file and the app that wants a file
3. Set some environment variables on the env2file container
4. Everything boots up, env2file creates the files, they are visible to the app, app is happy!

env2file will search for environment variables that match the following format:

```env
EF_(name|data)_(\w+)
```

Where the first group is either `name` or `data` and the second group is some unique name.

Each "target" requires two variables: one for the filename and one for the contents. So if you wanted to create a file
named `config.json` with some JSON in it, you'd declare two variables:

```env
EF_name_cfg=config.json
EF_data_cfg={"some":"json"}
```

The unique key (`cfg` in the above example) permits for as many targets as you want:

```env
EF_name_cfg=config.json
EF_data_cfg={"some":"json"}

EF_name_auth=auth.yaml
EF_data_auth=some: yaml

EF_name_other=other_stuff
EF_data_other=my secret pizza recipe
```

```yaml
version: "3.5"
services:
  someapp:
    image: some/app
    volumes:
      - /shared/config:/etc/someapp/config
  env2file:
    image: southclaws/env2file
    environment:
      E2F_name_config: /config/someapp-config.json
      E2F_data_config: |
        {
          "host": "127.0.0.1",
          "port": 4444
        }
      E2F_name_clientid: /config/client_identifier
      E2F_data_clientid: 37d060be-fb2e-11e9-99d0-645aede9143b
    volumes:
      - /shared/config:/config
```

In the above example, `some/app` will see a file named `someapp-config.json` inside the `/etc/someapp/config` directory
with the contents:

```json
{
  "host": "127.0.0.1",
  "port": 4444
}
```

And a file named `client_identifier` inside the same directory that contains simply
`37d060be-fb2e-11e9-99d0-645aede9143b`.

---

You can demo/play locally with the following docker run line:

```sh
docker run -v$(pwd)/files:/files -e EF_name_target=/files/target.json -e EF_data_target='{"a":"b"}' southclaws/env2file
```
