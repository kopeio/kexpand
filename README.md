## Kexpand

kexpand is a tool for expanding Kubernetes placeholder variables into their actual values.
It implements the upcoming kubernetes templating specification client-side (with some extensions based
on real-world requirements).

You can use templates today, and when k8s implements it on the server, you should hopefully
not have to rewrite your manifests.

## Syntax

quoted form: `$(key) => "value"`

unquoted form: `$((key)) => value`

## Example

kexpand turns this:

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: $((name))
spec:
  replicas: $((replicas))
  template:
    metadata:
      labels:
        app: $((name))
        release: $((release))
    spec:
      containers:
      - name: $((name))
        image: tutum/hello-world
        ports:
        - containerPort: $((port))
```

Into this:

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: hello-world
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: hello-world
        release: stable
    spec:
      containers:
      - name: hello-world
        image: tutum/hello-world
        ports:
        - containerPort: 80
```

## Installation

Build the code (make sure you have set GOPATH):
```
go get -d github.com/kopeio/kexpand
cd ${GOPATH}/src/github.com/kopeio/kexpand
make
```

## Expanding variables

The `expand` command will output the result of replacing all variables in the template file with:

* the values provided as command line parameters (-k)
* values from the provided values file ('-f')
* values defined in the environment ('-e')

`release=stable kexpand expand deployment.yaml -f values.yaml -k port=80 -e`

deployment.yaml:

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: $((name))
spec:
  replicas: $((replicas))
  template:
    metadata:
      labels:
        app: $((name))
        release: $((release))
    spec:
      containers:
      - name: $((name))
        image: tutum/hello-world
        ports:
        - containerPort: $((port))
```

values.yaml:

```
name: hello-world
replicas: 3
```

The result of expanding the template will be output to stdout, where it can be redirected to a file or piped to `kubectl`:

`RELEASE=stable kexpand expand deployment.yaml -f values.yaml -k port=80 -e | kubectl create -f -`

kexpand supports two different forms of variables: `$((key))` will output a non-quoted value (`$((replicas))` => `3`),
while `$(key)` will output a quoted value (`$(replicas)` => `"3"`). A legacy format, `{{key}}`, is also supported, but
not recommended for use.

## Variable names
Variables names can include any alphanumeric character along with hypens(-), period(.), and underscores(_).

## Multiple files
kexpand supprts passing multiple files and wildcards for templates to specify multiple files at once. kexpand will add `---` 
between each filename as required by kubectl.

`kexpand *.tmpl.yaml -f values.yaml`

## Base64 support.
Base64 encoding of values is supported by adding `|base64` to the key as in `$(key|base64)`, `$((key|base64))`, and `{{keyi|base64}}`.  This
is useful for secret definitions.
