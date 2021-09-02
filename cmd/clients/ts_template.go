package main

const tsIndexTemplate = `{{ range $service := .services }}import * as {{ $service.Name }} from './{{ $service.Name }}';
{{ end }}

export class Client {
	constructor(token: string) {
		{{ range $service := .services }}
		this.{{ $service.Name}}Service = new {{ $service.Name }}.{{ title $service.Name}}Service(token){{end}}
	}

{{ range $service := .services }}
	{{ $service.Name}}Service: {{ $service.Name }}.{{ title $service.Name}}Service;{{end}}
}
`

const tsServiceTemplate = `import * as m3o from '@m3o/m3o-node';

{{ $service := .service }}
export class {{ title $service.Name }}Service{
	private client: m3o.Client;

	constructor(token: string) {
		this.client = new m3o.Client({token: token})
	}
	{{ range $key, $req := $service.Spec.Components.RequestBodies }}
	{{ $endpointName := requestTypeToEndpointName $key}}{{ untitle $endpointName}}(request: {{ requestType $key }}): Promise<{{ requestTypeToResponseType $key }}> {
		return this.client.call("{{ $service.Name }}", "{{ requestTypeToEndpointPath $key}}", request) as Promise<{{ requestTypeToResponseType $key }}>;
	};
	{{ end }}
}

{{ range $typeName, $schema := $service.Spec.Components.Schemas }}
export interface {{ title $typeName }}{{ "{" }}
{{ recursiveTypeDefinition "typescript" $service.Name $typeName $service.Spec.Components.Schemas }}{{ "}" }}
{{end}}
`
