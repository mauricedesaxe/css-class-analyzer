# How fast can you analyze all the classes in a project?

A challenge, if you will, to see how fast I can analyze all CSS classes
in a HTML project.

## How to use

1. Ensure Go is installed on your system.
2. Clone the repository and navigate to the root directory.
3. Do not run the tests yet. Run these commands to set up the pre-commit hook:
```sh
cp hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```
4. To commit your changes, ensure no tests are failing and you only have no diff to the `performance_results.json` file. The pre-commit hook in will run the tests again, abort the commit if any test fails and update the `performance_results.json` with the latest performance metrics.
5. Check the `performance_results.json` file to see the performance metrics of your runs. This file is automatically updated and should not be manually modified.

Note: I'm aware this is a hacky way of testing performance. It's okay for now.