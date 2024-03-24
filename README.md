# How fast can you analyze all the classes in a project?

A challenge, if you will, to see how fast I can analyze all CSS classes
in a HTML project.

## How to use

1. Ensure Go is installed on your system.
2. Clone the repository and navigate to the root directory.
3. Use standard Go commands to run the project, tests, etc.

## Performance profile

At the time of writing, the performance profile (I ran 100 times) is as follows:

```
Total average Duration: 285.482459ms
Total median Duration: 322.099042ms
Average time per LoC: 569 ns
Average time per File: 917950 ns
```

This can certainly be improved upon, but it's a good start. A few ideas to explore going forward:
- Profile go routine management overhead and see if it can be optimized.
- Profile memory management outside of the html parser and see if it can be optimized.
- Profile the html parser and see if it can be optimized, altough I have a gut feeling this is not the main bottleneck right now.

## Feature parity

The feature set is incomplete as of now as the parser can only parse `.html` files. In the future,
you'd add support for other templating languages, such as `.jsx`, `.tsx`, and make use of dependency
injection (i.e.: `main.go` detects which language is being used and uses the appropriate parser, most of the
go routine and dir walking logic remaining similar).

## Contributing

If you would like to contribute, please open an issue or pull request. Make sure
to add the performance report (from github actions) to the pull request's description.