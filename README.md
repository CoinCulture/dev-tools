# breaking
Find all breaking changes to functions and methods in a Go project using git

# Run

```
breaking
```

will check the diff between the current branch and `master`.

To run it on a per-package level, pass a list of packages as arguments.
For instance, if using `glide` for dependencies, you can run

```
breaking $(glide novendor)
```
