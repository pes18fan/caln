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

## Note on accuracy

While caln is still based on a lot of approximations as there is not really any
method I know of to do an AD-to-BS conversion exactly without any such
approximations. Regardless of that; from my testing, it is accurate to within
minutes for finding the corresponding Gregorian date for specific
[sankranti](https://en.wikipedia.org/wiki/Sankranti), which are the base for
calculating the first day of each month. The actual problem though is that the
assignment of the start of a month in the official BS calendar in Nepal is not
as simple as assigning the day when the sankranti occurs; there is seemingly no
specific known way to do so. Hence, it is quite possible that caln fails for
some dates, it is not fully reliable. But for the most part it works!

## Info about VSOP87 tables

`VSOP87D.ear.clean` is adapted from `VSOP87D.ear` from
[vsop87](https://github.com/ctdk/vsop87).

The .clean file is obtained via this command:

```bash
cat VSOP87D.ear | tr -s ' ' | sed '/[A-Z]/d; s/^ //g; s/ $//g' > VSOP87D.ear.clean
```
