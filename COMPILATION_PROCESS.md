# Compilation process

## Compile

```
   +------------------------------+
   | ~/projects/*/.portfoliodb/** |
   +---------------+--------------+
                   |
                   |   pnpm database:crawl
                   v
+------------------+----------------+
| ~/projects/.buildarea/database/** |
+------------------+----------------+
                   |
                   |   pnpm database:build
                   v
        +----------+----------+     +-----------+                             +----------+          +-----------+         +-----------+
        | database/works.json |     | src/*.pug |                             | src/*.ts |          | src/*.css |         | static/** |
        +----------+----------+     +-----+-----+                             +----+-----+          +-----+-----+         +-----+-----+
                   |                      |                                        |                      |                     |
                   +--------------------->+   pnpm pug:hydrate                     |                      |                     |
                                          v                                        |                      |                     |
                         +----------------+---------------+                        |                      |                     |
                         | artifacts/phase_1/{lang}/*.pug |                        |                      |                     |
                         +----------------+---------------+                        |                      |                     |
                                          |                                        |                      |                     |
                                          |   pnpm messages:extract                |                      |                     |
                                          |                                        |                      |                     |
                              +-----------+-----------+                            |                      |                     |
                              |                       |                            |                      |                     |
                              v                       v                            |                      |                     |
                   +----------+--------+  +-----------+--------------------+       |                      |                     |
                   | messages/fr_FR.po |  | artifacts/phase_2/{lang}/*.pug |       |                      |                     |
                   +----------+--------+  +-----------+--------------------+       |                      |                     |
                              |                       |                            |                      |                     |
        pnpm messages:build   |                       |   pnpm pug:build           |                      |                     |
                              v                       v                            |                      |                     |
                   +----------+--------+  +-----------+--------------------+       |                      |                     |
                   | messages/fr_FR.mo |  |artifacts/phase_3/{lang}/*.html |       |                      |                     |
                   +----------+--------+  +-----------+--------------------+       |                      |                     |
                              |                       |                            |                      |                     |
                              |                       |                            |                      |                     |
                              +-----------+-----------+                            |                      |                     |
                                          |                                        |                      |                     |
                                          |   pnpm html:translate                  |                      |                     |
                                          v                                        |   pnpm js:build      |   pnpm css:build    |   pnpm static:build
                               +----------+---------+                              |                      |                     |
                               | dist/{lang}/*.html +<-----------------------------+----------------------+---------------------+
                               +--------------------+
```

## 

<dl>
  <dt></dt>
  <dd></dd>
</dl>

## Cleanup

- `rm -r trans/`
- `rm messages/dictionary.mo`

## Workflow

1. Edit `{src,}/**` files
2. Run `pnpm prepare:i18n`
3. Add missing translations to `messages/fr_FR.po`
4. Run `pnpm makeclean`

Or:

1. Edit source files
2. Run `./imake`
