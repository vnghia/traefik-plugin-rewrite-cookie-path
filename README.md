# Cookie Path Rewrite

Cookie Path Rewrite is a middleware plugin for [Traefik](https://traefik.io) which rewrites the path of a cookie in the response. Inspired by [SchmitzDan/traefik-plugin-cookie-path-prefix](https://github.com/SchmitzDan/traefik-plugin-cookie-path-prefix) and [XciD/traefik-plugin-rewrite-headers](https://github.com/XciD/traefik-plugin-rewrite-headers).

[![Build Status](https://github.com/vnghia/traefik-plugin-rewrite-cookie-path/workflows/Main/badge.svg?branch=main)](https://github.com/vnghia/traefik-plugin-rewrite-cookie-path/actions)

## Configuration

### Static

```yaml
experimental:
  plugins:
    cookiePathRewrite:
      modulename: "github.com/vnghia/traefik-plugin-rewrite-cookie-path"
      version: "v0.0.1"
```

### Dynamic

To configure the  plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/).
The following example creates and uses the cookie path prefix middleware plugin to replace the cookies path whose key is `someName` from `/foo` to `/bar`:

```yaml
http:
  middlewares:
    cookiePathRewrite:
      plugin:
        cookiePathRewrite:
          rewrites:
            - name: someName
              regex: "/foo"
              replacement: "/bar"
```

Configuration can also be set via toml or docker labels.
