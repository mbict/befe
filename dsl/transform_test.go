package dsl

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateWildcardTree_without_wildcard(t *testing.T) {
	in := patternTree{
		"data": patternTree{
			"1": patternTree{"bar": nil},
			"2": patternTree{"baz": nil},
		},
	}

	out := createWildcardTree(in)

	assert.Len(t, out, 1)
	assert.Len(t, out["data"], 2)
	assert.Len(t, out["data"]["1"], 1)
	assertPatternHasKey(t, "bar", out["data"]["1"])
	assert.Len(t, out["data"]["2"], 1)
	assertPatternHasKey(t, "baz", out["data"]["2"])
}

func TestCreateWildcardTree_with_wildcard(t *testing.T) {
	in := patternTree{
		"data": patternTree{
			"*": patternTree{"foo": nil},
			"1": patternTree{"bar": nil},
			"2": patternTree{"baz": nil},
		},
	}

	out := createWildcardTree(in)

	assert.Len(t, out, 1)
	assert.Len(t, out["data"], 3)
	assert.Len(t, out["data"]["*"], 1)
	assertPatternHasKey(t, "foo", out["data"]["*"])
	assert.Len(t, out["data"]["1"], 2)
	assertPatternHasKey(t, "foo", out["data"]["1"])
	assertPatternHasKey(t, "bar", out["data"]["1"])
	assert.Len(t, out["data"]["2"], 2)
	assertPatternHasKey(t, "foo", out["data"]["2"])
	assertPatternHasKey(t, "baz", out["data"]["2"])
}

func TestJSON_Filter_specific_slice_only_last_element(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
		},
	}

	err := Filter("data.1.*")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
		},
	}, data)
}

func TestJSON_Filter_nested_slice(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: []interface{}{

				0: map[string]interface{}{
					"foo": "foo0-0",
					"bar": "bar0-0",
					"baz": "baz0-0",
				},

				1: map[string]interface{}{
					"foo": "foo0-1",
					"bar": "bar0-1",
					"baz": "baz0-1",
				},
			},
			1: []interface{}{
				0: map[string]interface{}{
					"foo": "foo1-0",
					"bar": "bar1-0",
					"baz": "baz1-0",
				},
			},
		},
	}

	err := Filter("data.1.0.*")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: []interface{}{
				0: map[string]interface{}{
					"foo": "foo1-0",
					"bar": "bar1-0",
					"baz": "baz1-0",
				},
			},
		},
	}, data)
}

func TestJSON_Filter_nested_slice_without_wildcard(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: []interface{}{

				0: map[string]interface{}{
					"foo": "foo0-0",
					"bar": "bar0-0",
					"baz": "baz0-0",
				},

				1: map[string]interface{}{
					"foo": "foo0-1",
					"bar": "bar0-1",
					"baz": "baz0-1",
				},
			},
			1: []interface{}{
				0: map[string]interface{}{
					"foo": "foo1-0",
					"bar": "bar1-0",
					"baz": "baz1-0",
				},
			},
		},
	}

	err := Filter("data.0.1")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: []interface{}{
				0: map[string]interface{}{
					"foo": "foo0-1",
					"bar": "bar0-1",
					"baz": "baz0-1",
				},
			},
		},
	}, data)
}

func TestJSON_Filter_nested_slice_multiple_indexed_paths(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: []interface{}{

				0: map[string]interface{}{
					"foo": "foo0-0",
					"bar": "bar0-0",
					"baz": "baz0-0",
				},

				1: map[string]interface{}{
					"foo": "foo0-1",
					"bar": "bar0-1",
					"baz": "baz0-1",
				},
			},
			1: []interface{}{
				0: map[string]interface{}{
					"foo": "foo1-0",
					"bar": "bar1-0",
					"baz": "baz1-0",
				},
			},
		},
	}

	err := Filter("data.1.0", "data.0.1")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: []interface{}{
				0: map[string]interface{}{
					"foo": "foo0-1",
					"bar": "bar0-1",
					"baz": "baz0-1",
				},
			},
			1: []interface{}{
				0: map[string]interface{}{
					"foo": "foo1-0",
					"bar": "bar1-0",
					"baz": "baz1-0",
				},
			},
		},
	}, data)
}

func TestJSON_Filter_specific_slice_element_in_the_middle(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
			2: map[string]interface{}{
				"foo": "bar3",
				"bar": "foobar3",
				"baz": "test3",
			},
		},
	}

	err := Filter("data.1.*")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
		},
	}, data)
}

func TestJSON_Filter_specific_slice_element_first_and_last(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
			2: map[string]interface{}{
				"foo": "bar3",
				"bar": "foobar3",
				"baz": "test3",
			},
		},
	}

	err := Filter("data.0.*", "data.2.*")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar3",
				"bar": "foobar3",
				"baz": "test3",
			},
		},
	}, data)
}

func TestJSON_Filter_specific_slice_only_first_element(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
		},
	}

	err := Filter("data.0.*")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
		},
	}, data)
}

func TestJSON_Filter_data_array_multiple_fields_with_array_index(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
		},
	}

	err := Filter("data.*.foo", "data.1.bar")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
			},
		},
	}, data)
}

func TestJSON_Filter_data_with_wildcard_on_slice_and_fieldname(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
		},
	}

	err := Filter("data.*.*")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
		},
	}, data)
}

func TestJSON_Filter_data_array_multiple_fields(t *testing.T) {

	data := map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
				"baz": "test2",
			},
		},
	}

	err := Filter("data.*.foo", "data.*.bar")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"bar": "bazbar",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"bar": "foobar",
			},
		},
	}, data)
}

func TestJSON_Filter_data_array_single_field(t *testing.T) {
	data := map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
				"baz": "test1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
				"baz": "test2",
			},
		},
	}

	err := Filter("data.*.foo")(data)

	assert.NoError(t, err)
	AssertJson(t, map[string]interface{}{
		"data": []interface{}{
			0: map[string]interface{}{
				"foo": "bar1",
			},
			1: map[string]interface{}{
				"foo": "bar2",
			},
		},
	}, data)
}

func AssertJson(t *testing.T, expected interface{}, data interface{}) {
	expectedJson, _ := json.MarshalIndent(expected, "", " ")
	actualJson, _ := json.MarshalIndent(data, "", " ")
	assert.Equal(t, string(expectedJson), string(actualJson))
}

func assertPatternHasKey(t *testing.T, expectedKey string, data patternTree) {
	if _, ok := data[expectedKey]; !ok {
		assert.Fail(t, "key "+expectedKey+" not found")
	}
}
