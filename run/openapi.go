package run

import (
	"bytes"
	_ "embed"
	"net/http"
	"text/template"

	"github.com/swaggest/openapi-go/openapi3"
)

//go:embed openapi.html
var html string

//go:embed openapi.json
var schemaBytes []byte

// Schema returns OpenAPI schema.
func Schema[Input, Option, Solution any]() (string, error) {
	var s openapi3.Spec

	if err := s.UnmarshalJSON(schemaBytes); err != nil {
		return "", err
	}

	reflector := openapi3.Reflector{}
	reflector.Spec = &s

	summary := "Post a new run"
	t := true

	description := "The input of the algorithm"
	body := openapi3.RequestBodyOrRef{
		RequestBody: &openapi3.RequestBody{
			Description: &description,
			Required:    &t,
		},
	}

	postRun := openapi3.Operation{
		Summary:     &summary,
		RequestBody: &body,
	}

	err := reflector.SetupRequest(openapi3.OperationContext{
		Operation:  &postRun,
		Input:      new(Input),
		HTTPMethod: http.MethodPost,
		HTTPStatus: int(http.StatusOK),
	})
	if err != nil {
		return "", err
	}

	err = reflector.SetJSONResponse(
		&postRun, new(Solution), http.StatusOK,
	)
	if err != nil {
		return "", err
	}

	err = reflector.SetJSONResponse(&postRun, nil, http.StatusNotFound)
	if err != nil {
		return "", err
	}

	err = reflector.Spec.AddOperation(http.MethodPost, "/", postRun)
	if err != nil {
		return "", err
	}

	schema, err := reflector.Spec.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(schema), nil
}

func openAPI[Input, Option, Solution any](
	config HTTPRunnerConfig,
) (http.Handler, error) {
	rawSchema, err := Schema[Input, Option, Solution]()
	if err != nil {
		return nil, err
	}
	schemaTmpl, err := template.New("schema").Parse(rawSchema)
	if err != nil {
		return nil, err
	}
	addressStruct := struct {
		Address string
	}{
		Address: config.Runner.HTTP.Address,
	}

	var schema bytes.Buffer
	err = schemaTmpl.Execute(&schema, addressStruct)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("html").Parse(html)
	if err != nil {
		return nil, err
	}
	specStruct := struct {
		Spec string
	}{
		Spec: schema.String(),
	}
	var spec bytes.Buffer
	err = tmpl.Execute(&spec, specStruct)
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := w.Write(spec.Bytes())
		if err != nil {
			panic(err)
		}
	}), nil
}
