# New hydration process

```
src/...    html -> hydration -> translation -> dist/
             ↓                      ↑
          gettext                   |
         extraction                 |
             ↓                      |
          combine                   |
             ↓                      |
           build -------------------|


       svg,png,jpeg -> image processing     -> assets/
       styl         -> stylus compiler      -> assets/
       ls           -> livescript compiler  -> assets/

```

## Hydration & translation

### Regular pages

1. Execute template with regular data
1. For both languages:
    1. Translate
    1. Write to source path with `src/` → `dist/`

### Dynamic pages (`_{work,tag}.html,using/_technology.html`)

1. For every resource, for both languages:
    1. Execute template with regular data _and the language's resource's_,
    1. Translate
    1. Write to appropriate path in `dist/`

## Watchers

### HTML

#### Content changes (new file or content modified)

- `{components,layouts}/FILE`:
    1. update the modified file
    1. update files containing `{{ template "FILE" }}`
- others:
    1. update the modified file

#### Suppressions

1. Warn when some linked to that page

#### Moves

1. Change links referencing the file (triggers _Content changes_)

### Scripting & styling

Just plain-old watch-and-compile for: `*.{ls,styl}`

### Metadata database changes

(database that keeps valid tags & technologies)

#### Addition

1. Build a new page from `_tag.html` or `using/_technology.html`
1. (technologies only) process the tech's logo

### Database changes

#### Project addition

1. Crawl & build a database containing only that project
1. Merge it with `database.json`
1. Build a new page from `_work.html`

### Messages (translations) changes

1. Re-do everything
