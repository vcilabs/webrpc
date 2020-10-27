// Code generated by statik. DO NOT EDIT.

package embed

import (
	"github.com/rakyll/statik/fs"
)


func init() {
	data := "PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00client.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{define \"client\"}}\x0d\n{{- if .Services}}\x0d\n//\x0d\n// Client\x0d\n//\x0d\n\x0d\n{{- range .Services}}\x0d\nexport class {{.Name}} implements {{.Name | serviceInterfaceName}} {\x0d\n  private hostname: string\x0d\n  private fetch: Fetch\x0d\n  private path = '/rpc/{{.Name}}/'\x0d\n\x0d\n  constructor(hostname: string, fetch: Fetch) {\x0d\n    this.hostname = hostname\x0d\n    this.fetch = fetch\x0d\n  }\x0d\n\x0d\n  private url(name: string): string {\x0d\n    return this.hostname + this.path + name\x0d\n  }\x0d\n  {{range .Methods}}\x0d\n  {{.Name | methodName}} = ({{. | methodInputs}}): {{. | methodOutputs}} => {\x0d\n    return this.fetch(\x0d\n      this.url('{{.Name}}'),\x0d\n      {{- if .Inputs | len}}\x0d\n      createHTTPRequest(args, headers)\x0d\n      {{- else}}\x0d\n      createHTTPRequest({}, headers)\x0d\n      {{end -}}\x0d\n    ).then((res) => {\x0d\n      return buildResponse(res).then(_data => {\x0d\n        return {\x0d\n        {{- $outputsCount := .Outputs|len -}}\x0d\n        {{- range $i, $output := .Outputs}}\x0d\n          {{$output | newOutputArgResponse}}{{listComma $i $outputsCount}}\x0d\n        {{- end}}\x0d\n        }\x0d\n      })\x0d\n    })\x0d\n  }\x0d\n  {{end}}\x0d\n}\x0d\n{{end -}}\x0d\n{{end -}}\x0d\n{{end}}\x0d\nPK\x07\x08\xc7\xc9W\xa7J\x04\x00\x00J\x04\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x16\x00	\x00client_helpers.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{define \"client_helpers\"}}\x0d\nexport interface WebRPCError extends Error {\x0d\n  code: string\x0d\n  msg: string\x0d\n	status: number\x0d\n}\x0d\n\x0d\nconst createHTTPRequest = (body: object = {}, headers: object = {}): object => {\x0d\n  return {\x0d\n    method: 'POST',\x0d\n    headers: { ...headers, 'Content-Type': 'application/json' },\x0d\n    body: JSON.stringify(body || {})\x0d\n  }\x0d\n}\x0d\n\x0d\nconst buildResponse = (res: Response): Promise<any> => {\x0d\n  return res.text().then(text => {\x0d\n    let data\x0d\n    try {\x0d\n      data = JSON.parse(text)\x0d\n    } catch(err) {\x0d\n      throw { code: 'unknown', msg: `expecting JSON, got: ${text}`, status: res.status } as WebRPCError\x0d\n    }\x0d\n    if (!res.ok) {\x0d\n      throw data // webrpc error response\x0d\n    }\x0d\n    return data\x0d\n  })\x0d\n}\x0d\n\x0d\nexport type Fetch = (input: RequestInfo, init?: RequestInit) => Promise<Response>\x0d\n{{end}}\x0d\nPK\x07\x08\xb2\x85\x15R=\x03\x00\x00=\x03\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x11\x00	\x00proto.gen.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{- define \"proto\" -}}\x0d\n/* tslint:disable */\x0d\n// {{.Name}} {{.SchemaVersion}} {{.SchemaHash}}\x0d\n// --\x0d\n// This file has been generated by https://github.com/webrpc/webrpc using gen/typescript\x0d\n// Do not edit by hand. Update your webrpc schema and re-generate.\x0d\n\x0d\n// WebRPC description and code-gen version\x0d\nexport const WebRPCVersion = \"{{.WebRPCVersion}}\"\x0d\n\x0d\n// Schema version of your RIDL schema\x0d\nexport const WebRPCSchemaVersion = \"{{.SchemaVersion}}\"\x0d\n\x0d\n// Schema hash generated from your RIDL schema\x0d\nexport const WebRPCSchemaHash = \"{{.SchemaHash}}\"\x0d\n\x0d\n{{template \"types\" .}}\x0d\n\x0d\n{{- if .TargetOpts.Client}}\x0d\n  {{template \"client\" .}}\x0d\n  {{template \"client_helpers\" .}}\x0d\n{{- end}}\x0d\n\x0d\n{{- if .TargetOpts.Server}}\x0d\n  {{template \"server\" .}}\x0d\n  {{template \"server_helpers\" .}}\x0d\n{{- end}}\x0d\n\x0d\n{{- end}}\x0d\nPK\x07\x08\x1f\x0ev\xe9#\x03\x00\x00#\x03\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00server.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{define \"server\"}}\x0d\n\x0d\n{{- if .Services}}\x0d\n//\x0d\n// Server\x0d\n//\x0d\nexport class WebRPCError extends Error {\x0d\n    statusCode?: number\x0d\n\x0d\n    constructor(msg: string = \"error\", statusCode?: number) {\x0d\n        super(\"webrpc error: \" + msg);\x0d\n\x0d\n        Object.setPrototypeOf(this, WebRPCError.prototype);\x0d\n\x0d\n        this.statusCode = statusCode;\x0d\n    }\x0d\n}\x0d\n\x0d\nimport express from 'express'\x0d\n\x0d\n    {{- range .Services}}\x0d\n        {{$name := .Name}}\x0d\n        {{$serviceName := .Name | serviceInterfaceName}}\x0d\n\x0d\n        export type {{$serviceName}}Service = {\x0d\n            {{range .Methods}}\x0d\n                {{.Name}}: (args: {{.Name}}Args) => {{.Name}}Return | Promise<{{.Name}}Return>\x0d\n            {{end}}\x0d\n        }\x0d\n\x0d\n        export const create{{$serviceName}}App = (serviceImplementation: {{$serviceName}}Service) => {\x0d\n            const app = express();\x0d\n\x0d\n            app.use(express.json())\x0d\n\x0d\n            app.post('/*', async (req, res) => {\x0d\n                const requestPath = req.baseUrl + req.path\x0d\n\x0d\n                if (!req.body) {\x0d\n                    res.status(400).send(\"webrpc error: missing body\");\x0d\n\x0d\n                    return\x0d\n                }\x0d\n\x0d\n                switch(requestPath) {\x0d\n                    {{range .Methods}}\x0d\n\x0d\n                    case \"/rpc/{{$name}}/{{.Name}}\": {                        \x0d\n                        try {\x0d\n                            {{ range .Inputs }}\x0d\n                                {{- if not .Optional}}\x0d\n                                    if (!(\"{{ .Name }}\" in req.body)) {\x0d\n                                        throw new WebRPCError(\"Missing Argument `{{ .Name }}`\")\x0d\n                                    }\x0d\n                                {{end -}}\x0d\n\x0d\n                                if (\"{{ .Name }}\" in req.body && !validateType(req.body[\"{{ .Name }}\"], \"{{ .Type | jsFieldType }}\")) {\x0d\n                                    throw new WebRPCError(\"Invalid Argument: {{ .Name }}\")\x0d\n                                }\x0d\n                            {{end}}\x0d\n\x0d\n                            const response = await serviceImplementation[\"{{.Name}}\"](req.body);\x0d\n\x0d\n                            {{ range .Outputs}}\x0d\n                                if (!(\"{{ .Name }}\" in response)) {\x0d\n                                    throw new WebRPCError(\"internal\", 500);\x0d\n                                }\x0d\n                            {{end}}\x0d\n\x0d\n                            res.status(200).json(response);\x0d\n                        } catch (err) {\x0d\n                            if (err instanceof WebRPCError) {\x0d\n                                const statusCode = err.statusCode || 400\x0d\n                                const message = err.message\x0d\n\x0d\n                                res.status(statusCode).json({\x0d\n                                    msg: message,\x0d\n                                    status: statusCode,\x0d\n                                    code: \"\"\x0d\n                                });\x0d\n\x0d\n                                return\x0d\n                            }\x0d\n\x0d\n                            if (err.message) {\x0d\n                                res.status(400).send(err.message);\x0d\n\x0d\n                                return;\x0d\n                            }\x0d\n\x0d\n                            res.status(400).end();\x0d\n                        }\x0d\n                    }\x0d\n\x0d\n                    return;\x0d\n                    {{end}}\x0d\n\x0d\n                    default: {\x0d\n                        res.status(404).end()\x0d\n                    }\x0d\n                }\x0d\n            });\x0d\n\x0d\n            return app;\x0d\n        };\x0d\n    {{- end}}\x0d\n{{end -}}\x0d\n{{end}}\x0d\nPK\x07\x08L\xdc1\x96\xfe\x0d\x00\x00\xfe\x0d\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x16\x00	\x00server_helpers.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{ define \"server_helpers\" }}\x0d\n\x0d\nconst JS_TYPES = [\x0d\n    \"bigint\",\x0d\n    \"boolean\",\x0d\n    \"function\",\x0d\n    \"number\",\x0d\n    \"object\",\x0d\n    \"string\",\x0d\n    \"symbol\",\x0d\n    \"undefined\"\x0d\n]\x0d\n\x0d\n{{ range .Messages }}\x0d\n    const validate{{ .Name }} = (value: any) => {\x0d\n        {{ range .Fields }}\x0d\n            {{ if .Optional }}\x0d\n                if (\"{{ . | exportedJSONField }}\" in value && !validateType(value[\"{{ . | exportedJSONField }}\"], \"{{ .Type | jsFieldType }}\")) {\x0d\n                    return false\x0d\n                }\x0d\n            {{ else }}\x0d\n                if (!(\"{{ . | exportedJSONField }}\" in value) || !validateType(value[\"{{ . | exportedJSONField }}\"], \"{{ .Type | jsFieldType }}\")) {\x0d\n                    return false\x0d\n                }\x0d\n            {{ end }}\x0d\n        {{ end }}\x0d\n\x0d\n        return true\x0d\n    }\x0d\n{{ end }}\x0d\n\x0d\nconst TYPE_VALIDATORS: { [type: string]: (value: any) => boolean } = {\x0d\n    {{ range .Messages }}\x0d\n        {{ .Name }}: validate{{ .Name }},\x0d\n    {{ end }}\x0d\n}\x0d\n\x0d\nconst validateType = (value: any, type: string) => {\x0d\n    if (JS_TYPES.indexOf(type) > -1) {\x0d\n        return typeof value === type;\x0d\n    }\x0d\n\x0d\n    const validator = TYPE_VALIDATORS[type];\x0d\n\x0d\n    if (!validator) {\x0d\n        return false;\x0d\n    }\x0d\n\x0d\n    return validator(value);\x0d\n}\x0d\n\x0d\n{{ end }}PK\x07\x08_o\x9d\x06\x01\x05\x00\x00\x01\x05\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0d\x00	\x00types.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{define \"types\"}}\x0d\n//\x0d\n// Types\x0d\n//\x0d\n\x0d\n{{- if .Messages -}}\x0d\n{{range .Messages -}}\x0d\n\x0d\n{{if .Type | isEnum -}}\x0d\n{{$enumName := .Name}}\x0d\nexport enum {{$enumName}} {\x0d\n{{- range $i, $field := .Fields}}\x0d\n  {{- if $i}},{{end}}\x0d\n  {{$field.Name}} = '{{$field.Name}}'\x0d\n{{- end}}\x0d\n}\x0d\n{{end -}}\x0d\n\x0d\n{{- if .Type | isStruct  }}\x0d\nexport interface {{.Name | interfaceName}} {\x0d\n  {{- range .Fields}}\x0d\n  {{if . | exportableField -}}{{. | exportedJSONField}}{{if .Optional}}?{{end}}: {{.Type | fieldType}}{{- end -}}\x0d\n  {{- end}}\x0d\n}\x0d\n{{end -}}\x0d\n{{end -}}\x0d\n{{end -}}\x0d\n\x0d\n{{if .Services}}\x0d\n{{- range .Services}}\x0d\nexport interface {{.Name | serviceInterfaceName}} {\x0d\n{{- range .Methods}}\x0d\n  {{.Name | methodName}}({{. | methodInputs}}): {{. | methodOutputs}}\x0d\n{{- end}}\x0d\n}\x0d\n\x0d\n{{range .Methods -}}\x0d\nexport interface {{. | methodArgumentInputInterfaceName}} {\x0d\n{{- range .Inputs}}\x0d\n  {{.Name}}{{if .Optional}}?{{end}}: {{.Type | fieldType}}\x0d\n{{- end}}\x0d\n}\x0d\n\x0d\nexport interface {{. | methodArgumentOutputInterfaceName}} {\x0d\n{{- range .Outputs}}\x0d\n  {{.Name}}{{if .Optional}}?{{end}}: {{.Type | fieldType}}\x0d\n{{- end}}  \x0d\n}\x0d\n{{end}}\x0d\n\x0d\n{{- end}}\x0d\n{{end -}}\x0d\n{{end}}\x0d\nPK\x07\x08=&\xa0\x18r\x04\x00\x00r\x04\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQ\xc7\xc9W\xa7J\x04\x00\x00J\x04\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81\x00\x00\x00\x00client.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQ\xb2\x85\x15R=\x03\x00\x00=\x03\x00\x00\x16\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81\x8f\x04\x00\x00client_helpers.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQ\x1f\x0ev\xe9#\x03\x00\x00#\x03\x00\x00\x11\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81\x19\x08\x00\x00proto.gen.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQL\xdc1\x96\xfe\x0d\x00\x00\xfe\x0d\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81\x84\x0b\x00\x00server.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQ_o\x9d\x06\x01\x05\x00\x00\x01\x05\x00\x00\x16\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81\xc7\x19\x00\x00server_helpers.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQ=&\xa0\x18r\x04\x00\x00r\x04\x00\x00\x0d\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81\x15\x1f\x00\x00types.ts.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x05\x06\x00\x00\x00\x00\x06\x00\x06\x00\xb0\x01\x00\x00\xcb#\x00\x00\x00\x00"
		fs.Register(data)
	}
	