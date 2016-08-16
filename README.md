## Kexpand

kexpand is a tool for expanding Kubernetes placeholder variables into their actual values.

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
    spec:
      containers:
      - name: $((name))
        image: tutum/hello-world
        ports:
        - containerPort: 80
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

The `expand` command will output the result of replacing all variables in the template file with the values from the
provided values file:

`kexpand expand deployment.yaml -f values.yaml`

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
    spec:
      containers:
      - name: $((name))
        image: tutum/hello-world
        ports:
        - containerPort: 80
```

values.yaml:

```
name: hello-world
replicas: 3
```

The result of expanding the template will be output to stdout, where it can be redirected to a file or piped to `kubectl`:

`kexpand expand deployment.yaml -f values.yaml | kubectl create -f -`

kexpand supports two different forms of variables: `$((key))` will output a non-quoted value (`$((replicas))` => `3`),
while `$(key)` will output a quoted value (`$(replicas)` => `"3"`). A legacy format, `{{key}}`, is also supported, but
not recommended for use.
