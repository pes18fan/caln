# caln

Gregorian to BS date converter that doesn't use a lookup table

## Info about VSOP87 tables

`VSOP87D.ear.clean` is adapted from `VSOP87D.ear` from
[vsop87](https://github.com/ctdk/vsop87).

The .clean file is obtained via this command:

```bash
cat VSOP87D.ear | tr -s ' ' | sed '/[A-Z]/d; s/^ //g; s/ $//g' > VSOP87D.ear.clean
```
