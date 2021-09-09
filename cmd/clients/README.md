# Client and example generation

To run the code generation, from the repo root issue:


```sh
go install ./cmd/clients; clients .
```

The generated clients will end up in `./clients`.

Take inspiration from the `.github/workflows/publish.yml` to see how to publish the NPM package.


# Typescript gotchas

There is some funkiness going on with the package names in the generator - 