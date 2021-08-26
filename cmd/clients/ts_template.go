package main

const tsTemplate = `import * as m3o from '@m3o/m3o-node';

class Client {
	constructor(token: string) {
		{{ range $service := .services }}
		this.{{ $service.Name}}Service = new {{ title $service.Name}}Service(token){{end}}
	}

{{ range $service := .services }}
	{{ $service.Name}}Service: {{ title $service.Name}}Service;{{end}}
}
{{ range $service := .services }}
export class {{ title $service.Name }}Service{
	private client: m3o.Client;

	constructor(token: string) {
		this.client = new m3o.Client({token: token})
	}
	{{ range $key, $req := $service.Spec.Components.RequestBodies }}
	{{ $endpointName := requestTypeToEndpointName $key}}{{ untitle $endpointName}}(request: {{ $key }}): Promise<{{ requestTypeToResponseType $key }}> {
		return this.client.call("{{ $service.Name }}", "{{ requestTypeToEndpointPath $key}}", request) as Promise<{{ requestTypeToResponseType $key }}>;
	};
	{{ end }}
}
{{end}}

{{ range $service := .services }}{{ range $typeName, $schema := .Spec.Components.Schemas }}
export interface {{ title $service.Name }}{{ title $typeName }}{{ "{" }}
{{ recursiveTypeDefinition "typescript" $service.Name $typeName $service.Spec.Components.Schemas }}{{ "}" }}
{{end}}{{end}}
`
