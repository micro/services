package main

const goIndexTemplate = `package micro

import(
	{{ range $service := .services }}"github.com/micro/micro-go/{{ $service.Name}}"
{{ end }}
)

func NewClient(token string) *Client {
	return &Client{
		token: token,
		{{ range $service := .services }}
		{{ title $service.Name }}Service: {{ $service.Name }}.New{{ title $service.Name}}Service(token),{{end}}
	}
}

type Client struct {
	token string
{{ range $service := .services }}
	{{ title $service.Name }}Service *{{ $service.Name }}.{{ title $service.Name }}Service{{end}}
}
`

const goServiceTemplate = `{{ $service := .service }}package {{ $service.Name }}

import(
	"github.com/m3o/m3o-go/client"
)

func New{{ title $service.Name }}Service(token string) *{{ title $service.Name }}Service {
	return &{{ title $service.Name }}Service{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type {{ title $service.Name }}Service struct {
	client *client.Client
}

{{ range $key, $req := $service.Spec.Components.RequestBodies }}
{{ $endpointName := requestTypeToEndpointName $key}}func (t *{{ title $service.Name }}Service) {{ $endpointName }}(request {{ requestType $key }}) (*{{ requestTypeToResponseType $key }}, error) {
	rsp := &{{ requestTypeToResponseType $key }}{}
	return rsp, t.client.Call("{{ $service.Name }}", "{{ requestTypeToEndpointPath $key}}", request, rsp)
}
{{ end }}


{{ range $typeName, $schema := $service.Spec.Components.Schemas }}
type {{ title $typeName }} struct {{ "{" }}
{{ recursiveTypeDefinition "go" $service.Name $typeName $service.Spec.Components.Schemas }}{{ "}" }}
{{end}}
`

const goExampleTemplate = `{{ $service := .service }}package example

import(
	"fmt"
	"github.com/micro/micro-go/{{ $service.Name }}"
)

{{if .example.Description}}
// {{ .example.Description }}{{end}}
func {{ .funcName }}() {
	{{ $service.Name }}Service := {{ $service.Name }}.New{{ title $service.Name }}Service("YOUR_MICRO_TOKEN_HERE")
	rsp, _ := {{ $service.Name }}Service.{{ title .endpoint }}({{ $service.Name }}.{{ title .endpoint }}Request{
		{{ goExampleRequest $service.Name .endpoint $service.Spec.Components.Schemas .example.Request }}
	})
	fmt.Println(rsp)
}
`
