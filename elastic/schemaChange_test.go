package elastic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUrlWithNonNumericRetry(t *testing.T) {
	_, _, err := parseURL("idm_employee_v5/_update_by_query?retry=a")
	assert.Error(t, err)
}

func TestParseUrlWithRetry(t *testing.T) {
	url, retry, err := parseURL("idm_employee_v5/_update_by_query?retry=3")
	assert.NoError(t, err)
	assert.Equal(t, "idm_employee_v5/_update_by_query", url)
	assert.Equal(t, 3, retry)
}

func TestParseUrlWithUppercaseRetry(t *testing.T) {
	url, retry, err := parseURL("idm_employee_v5/_update_by_query?RETRY=3")
	assert.NoError(t, err)
	assert.Equal(t, "idm_employee_v5/_update_by_query?RETRY=3", url)
	assert.Equal(t, 0, retry)
}

func TestParseUrlEndingWithRetry(t *testing.T) {
	url, retry, err := parseURL("idm_employee_v5/_update_by_query?foo=bar&retry=2")
	assert.NoError(t, err)
	assert.Equal(t, "idm_employee_v5/_update_by_query?foo=bar", url)
	assert.Equal(t, 2, retry)
}

func TestParseUrlNoRetry(t *testing.T) {
	url, retry, err := parseURL("idm_employee_v5/_update_by_query?foo=bar")
	assert.NoError(t, err)
	assert.Equal(t, "idm_employee_v5/_update_by_query?foo=bar", url)
	assert.Equal(t, 0, retry)
}

func TestShardAndReplicaTokenReplacementWithNoTokens(t *testing.T) {
	sc := NewSchemaChange("../tests/index_template.js", 2, 2)
	assert.Contains(t, sc.Action.JSON, `"index.number_of_shards": 5`)
	assert.Contains(t, sc.Action.JSON, `"index.number_of_replicas": 0`)
	assert.Equal(t, 2, sc.Shards)
	assert.Equal(t, 2, sc.Replicas)
}

func TestShardTokenReplacementWithTokens(t *testing.T) {
	sc := NewSchemaChange("../tests/index_template_with_shards.js", 2, 2)
	assert.Contains(t, sc.Action.JSON, `"index.number_of_shards": 2`)
	assert.Contains(t, sc.Action.JSON, `"index.number_of_replicas": 1`)
	assert.Equal(t, 2, sc.Shards)
	assert.Equal(t, 2, sc.Replicas)
}

func TestReplicaTokenReplacementWithNoTokens(t *testing.T) {
	sc := NewSchemaChange("../tests/index_template_with_replicas.js", 2, 2)
	assert.Contains(t, sc.Action.JSON, `"index.number_of_shards": 3`)
	assert.Contains(t, sc.Action.JSON, `"index.number_of_replicas": 2`)
	assert.Equal(t, 2, sc.Shards)
	assert.Equal(t, 2, sc.Replicas)
}

func TestShardAndReplicaTokenReplacementWithTokens(t *testing.T) {
	sc := NewSchemaChange("../tests/index_template_with_shards_replicas.js", 2, 2)
	assert.Contains(t, sc.Action.JSON, `"index.number_of_shards": 2`)
	assert.Contains(t, sc.Action.JSON, `"index.number_of_replicas": 2`)
	assert.Equal(t, 2, sc.Shards)
	assert.Equal(t, 2, sc.Replicas)
}
