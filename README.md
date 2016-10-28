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

## Multiline Expansion

kexpand is primarily intended for simple variable replacement, but since it will replace any arbitrary string, multiline templates can also be achieved. This can be useful when you need to set to set several keys in one environment, but in another they are not present at all.

deployment.yaml:

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: $((name))
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: $((name))
    spec:
      containers:
      - name: $((name))
        image: some-image
        ports:
        - containerPort: 80
        env:
        $((aws_access_keys))
```

values.yaml:

```
name: my-app
aws_access_keys: |-
  - name: AWS_ACCESS_KEY_ID
        value: fake
        - name: AWS_SECRET_ACCESS_KEY
        value: fake
```

`kexpand expand deployment.yaml -f values.yaml`:

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-app
        image: some-image
        ports:
        - containerPort: 80
        env:
        - name: AWS_ACCESS_KEY_ID
          value: fake
        - name: AWS_SECRET_ACCESS_KEY
          value: fake
```

Note that for this to work you must use `$((var))` style replacement, since `$(var)` will result in quotes that break the YAML. Additionally, you must handle proper indention levels yourself, since kexpand will just take the entire string as is.

