# How fast can you analyze all the classes in a project?

A challenge, if you will, to see how fast I can analyze all CSS classes
in a HTML project.

## How to use

1. Ensure Go is installed on your system.
2. Clone the repository and navigate to the root directory.
3. Use standard Go commands to run the project, tests, etc.

## Performance profile of the analyzer (that's the core of the project)

At the time of writing, the performance profile (I ran 100 times) is as follows:

From my local machine (Apple M2 Max 64GB RAM):
```
Total average Duration: 285.482459ms
Total median Duration: 322.099042ms
Average time per LoC: 569 ns
Average time per File: 917950 ns
```

From Github Actions:
```
Total average Duration: 237.337131ms
Total median Duration: 223.923411ms
Average time per LoC: 473 ns
Average time per File: 763141 ns
```

This can certainly be improved upon, but it's a good start. A few ideas to explore going forward:
- Profile go routine management overhead and see if it can be optimized.
- Profile memory management outside of the html parser and see if it can be optimized.
- Profile the html parser and see if it can be optimized, altough I have a gut feeling this is not the main bottleneck right now.

## Feature parity of the analyzer

The feature set is incomplete as of now as the parser can only parse `.html` files. In the future,
you'd add support for other templating languages, such as `.jsx`, `.tsx`, and make use of dependency
injection (i.e.: `main.go` detects which language is being used and uses the appropriate parser, most of the
go routine and dir walking logic remaining similar).

## The web component

The API is a simple web server that provides access to the analyzer. It's currently very hacky
and under-optimized, but it works. It's meant to help create a fun little web tool. Ultimately,
the analyzer should be developed and used as a CLI tool.

## Contributing

If you would like to contribute, please open an issue or pull request. Make sure
to add the performance report (from github actions) to the pull request's description.