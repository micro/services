package main

const goTemplate = `package m3o

import(
	"github.com/m3o/m3o-go/client"
)

func NewClient(token string) *Client {
	return &Client{
		token: token,
		{{ range $service := .services }}
		{{ title $service.Name}}Service: New{{ title $service.Name}}Service(token),{{end}}
	}
}

type Client struct {
	token string
{{ range $service := .services }}
	{{ title $service.Name}}Service *{{ title $service.Name}}Service{{end}}
}
{{ range $service := .services }}
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
{{ $endpointName := requestTypeToEndpointName $key}}func (t *{{ title $service.Name }}Service) {{ $endpointName}}(request {{ $key }}) (*{{ requestTypeToResponseType $key }}, error) {
	rsp := &{{ requestTypeToResponseType $key }}{}
	return rsp, t.client.Call("{{ $service.Name }}", "{{ requestTypeToEndpointPath $key}}", request, rsp)
}
{{ end }}
{{end}}

{{ range $service := .services }}{{ range $typeName, $schema := .Spec.Components.Schemas }}
type {{ title $service.Name }}{{ title $typeName }} struct {{ "{" }}
{{ recursiveTypeDefinition "go" $service.Name $typeName $service.Spec.Components.Schemas }}{{ "}" }}
{{end}}{{end}}
`
