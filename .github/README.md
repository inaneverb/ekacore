<p align="center">
  <img alt="Privet" height="125" src="https://raw.githubusercontent.com/qioalice/ekago/master/.github/logo.svg">
  <br>
</p>
<p align="center">
<sub>
Ekago is a useful Golang packages toolkit that I wrote for myself and published for others.
<br>
And this is the first toolkit that contains both of improved logging and error generating mechanisms which are strongly linked to each other.
Honestly, this project has been started as "another one logging library", but now it's not only that. 
<br>
The root of project does not provide any code. It's toolkit, remember?
All subpackages that this package has been named starting with "eka" avoiding your package's names import conflicting and import renaming. A trifle, but nice.
</sub>
</p>

---

Each **Ekago** subpackage has its own readme file where you can meet with them face by face, read more and get an answers to "What? Why? How?".

I strongly recommend to start your tour with the
- [`ekaerr`: _a package of error management that you desire_](ekaerr/),  and then go to
- [`ekalog`: _an intelligence  logging package for everything and especially for_ `ekaerr`](ekalog/)

It doesn't take a lot of time. See you!

-----

And if you need a quick description about all other packages, see:

- [`ekadeath`: _When_ `os.Exit(1)` _is not enough and you love destructors (and their calls)_](/ekadeath)
- [`ekadanger`: _Wanna see what_ `interface{}` _is? Maybe compare functions? Divide by zero?_](/ekadanger)
- [`ekamath`: _Bored about_ `math.Min(upperBound, math.Max(lowerBound, x))`_?_](/ekamath)
- [`ekatime`: _Adore unixtime? Want to extract daystart, dayend? What's about_ `time.Now()` _with rough precision?_](/ekatime)
- [`ekatype`: _Oh, you need UUID? Predefined interfaces? SQL type that if_ `NULL` _then JSON is_ `null`_?_](/ekatype)