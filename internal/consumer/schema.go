package consumer

import "encoding/json"

var Schema json.RawMessage = []byte(`{
		"type": "object",
		"properties": {
			"problems": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"problem": {
						"type": "string"
					},
					"description": {
						"type": "string"
					},
					"weight": {
						"type": "number"
					},
					"datapoints": {
						"type": "number"
					}
				},
				"required": ["problem", "description", "weight", "datapoints"],
				"additionalProperties": false
				}
			}
		},
		"required": ["problems"],
		"additionalProperties": false
	}`)
