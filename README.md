# shml
shml is a shell script based markup language and templating engine.


## Usage

```
m := map[string]string{"key":"value"}
template := shml.New()
template.Parse(`foo ${key}`)
out, _ := template.Execute(m)
fmt.Printf("%s", out)
```

This will produce

```
foo value
```

## Documentation

### Variables

```
${var}
```

### Transforms
Functions are used to mutate variable values.  To use functions you simply append a
`|` followed by the name of the transform.

```
${var|func}
```

#### json
This transforms the data to JSON format

Example:

```
${data|json}
```
