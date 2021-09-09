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
	{{ range $key, $req := $service.Spec.Components.RequestBodies }}{{ $endpointName := requestTypeToEndpointName $key}}{{ if endpointComment $endpointName $service.Spec.Components.Schemas }}{{ endpointComment $endpointName $service.Spec.Components.Schemas }}{{ end }}{{ untitle $endpointName}}(request: {{ requestType $key }}): Promise<{{ requestTypeToResponseType $key }}> {
		return this.client.call("{{ $service.Name }}", "{{ requestTypeToEndpointPath $key}}", request) as Promise<{{ requestTypeToResponseType $key }}>;
	};
	{{ end }}
}

{{ range $typeName, $schema := $service.Spec.Components.Schemas }}
export interface {{ title $typeName }}{{ "{" }}
{{ recursiveTypeDefinition "typescript" $service.Name $typeName $service.Spec.Components.Schemas }}{{ "}" }}
{{end}}
`

const tsExampleTemplate = `{{ $service := .service }}const {{ $service.Name }} = require('@m3o/services/{{ $service.Name }}');

{{ if endpointComment .endpoint $service.Spec.Components.Schemas }}{{ endpointComment .endpoint $service.Spec.Components.Schemas }}{{ end }}async function {{ .funcName }}() {
	let {{ $service.Name }}Service = {{ $service.Name }}.New{{ title $service.Name }}Service(process.env.MICRO_API_TOKEN)
	let rsp = await {{ $service.Name }}Service.{{ .endpoint }}({{ tsExampleRequest $service.Name .endpoint $service.Spec.Components.Schemas .example.Request }})
	console.log(rsp)
}

await {{ .funcName }}()`
