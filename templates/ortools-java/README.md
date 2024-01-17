# Java OR-Tools Maven Example

Sample project explaining how to use Google OR-Tools in a maven project that can
be executed on the Nextmv platform.

## Usage

Java templates in Nextmv require a `main.jar` as an entry point. Running the
following command will generate a `main.jar` in the root direcotry of the
project.

```sh
mvn package
```

After that you can run the `main.jar` file with the following command:

```sh
java -jar main.jar --input input.json
```

You can also push the `main.jar` file to Nextmv platform and run it there. Take
a look at the `app.yml` file to see how the `main.jar` file is referenced.

```sh
nextmv app push -a <app-name>
```

## References

- [Google OR-Tools](https://github.com/or-tools/or-tools)
- [Google OR-Tools Java Examples](https://github.com/or-tools/java_or-tools)
