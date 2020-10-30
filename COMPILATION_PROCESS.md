# Compilation process

## Compile


_not shown: compilation of static assets (TS→JS, Stylus→CSS)_
```
            src/*.pug
                ↓  go template rendering
         trans/src/*.pug
                ↓  gettext-extract
             ↙      ↘
      fr_FR.po   trans/src/*.pug
        ↓ msgfmt         ↓ pug-cli
      fr_FR.mo   trans/src/*.html
        ↘                ↙
                ↓  go apply .mo file
        dist/{fr,en}/*.html
```

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
