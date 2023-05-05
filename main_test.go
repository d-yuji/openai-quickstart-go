package main

import "testing"

func TestOpenAIAPI(t *testing.T) {
	testCases := []struct {
		desc   string
		prompt string
	}{
		{
			desc:   "test",
			prompt: generatePrompt("small dog"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data, err := postOpenAIAPI(tC.prompt)
			t.Log(data)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
