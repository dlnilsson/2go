# 2go
Instantly converts JSON|YAML into a Go type definition in the terminal.

## Installation
```bash
go install github.com/dlnilsson/2go@latest
```

## Example

### JSON from stdin
```bash
curl -s https://httpbin.org/json | 2go
```

#### output

```go
type Autogenerated struct {
	Slideshow Slideshow `json:"slideshow,omitempty"`
}

type Slides struct {
	Title string `json:"title,omitempty"`
	Type  string `json:"type,omitempty"`
}

type Slideshow struct {
	Author string   `json:"author,omitempty"`
	Date   string   `json:"date,omitempty"`
	Slides []Slides `json:"slides,omitempty"`
	Title  string   `json:"title,omitempty"`
}
```
### YAML from stdin

```bash
cat togo/testdata/simple.yaml | 2go
```

```go
type Autogenerated struct {
	Kind     string   `yaml:"kind,omitempty"`
	Metadata Metadata `yaml:"metadata,omitempty"`
}

type Metadata struct {
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}
```

### with jq
```bash
jq -n --arg favorite "$(echo 'hello world' | base64)" '{fruits: ["apple", "banana", "cherry"], favorite: $favorite}' | 2go

type Autogenerated struct {
	Favorite []byte   `json:"favorite,omitempty"`
	Fruits   []string `json:"fruits,omitempty"`
}

```
