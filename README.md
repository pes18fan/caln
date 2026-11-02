# caln

Gregorian to Bikram Sambat date converter that doesn't use a lookup table.

By lookup table I refer to a pre-filled table consisting of known mappings from
certain Bikram Sambat dates to their corresponding Gregorian dates. This works
but is fragile and needs frequent updating and maintenance to continue being
reliable. 

caln takes a different approach; it instead computes these tables from scratch
based on astronomical formulae. These calculations do involve using one kind of
lookup table which is a table of precalculated coefficients based on the VSOP87
mathematical model; however this is a different kind of lookup table because it
is just math results.

caln is still based on a lot of approximations, though, as there is not really
any method I know of to do an AD-to-BS conversion exactly without any such
approximations. Hence, I can't guarantee that it will ALWAYS be correct; but
for the most part it should be!

## Info about VSOP87 tables

`VSOP87D.ear.clean` is adapted from `VSOP87D.ear` from
[vsop87](https://github.com/ctdk/vsop87).

The .clean file is obtained via this command:

```bash
cat VSOP87D.ear | tr -s ' ' | sed '/[A-Z]/d; s/^ //g; s/ $//g' > VSOP87D.ear.clean
```
