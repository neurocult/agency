package openai

import (
	"context"

	"github.com/eqtlab/lib/core"
)

var promptTemplate = `Extract the desired information from the following passage.

Only extract the properties mentioned in the 'information_extraction' function.

Passage:
{%s}
`

type TagParams struct {
	Props map[string]string
}

func (f Factory) Tag(params TagParams) *core.Pipe {
	return core.NewPipe(func(ctx context.Context, msg core.Message, opts ...core.PipeOption) (core.Message, error) {
		// promptTemplate
		
		return nil, nil
	})
}

// 	   function = _get_tagging_function(schema)
//     prompt = prompt or ChatPromptTemplate.from_template(_TAGGING_TEMPLATE)
//     output_parser = JsonOutputFunctionsParser()
//     llm_kwargs = get_llm_kwargs(function)
//     chain = LLMChain(
//         llm=llm,
//         prompt=prompt,
//         llm_kwargs=llm_kwargs,
//         output_parser=output_parser,
//         **kwargs,
//     )
//     return chain

// def _get_tagging_function(schema: dict) -> dict:
//     return {
//         "name": "information_extraction",
//         "description": "Extracts the relevant information from the passage.",
//         "parameters": _convert_schema(schema),
//     }

// def _convert_schema(schema: dict) -> dict:
//     props = {k: {"title": k, **v} for k, v in schema["properties"].items()}
//     return {
//         "type": "object",
//         "properties": props,
//         "required": schema.get("required", []),
//     }
