# go-logging
A simple wrapper around `logrus` to use in libraries, enabling easy replacement with other logging libraries.

## Usage
Simply use `Logger` interface in your code and use dependency injection to provide logging functionality to deeper layers.
You may use `logrus` sub-package to easily use `logrus` in your project as a default implementation.

## License
Apache 2.0 License