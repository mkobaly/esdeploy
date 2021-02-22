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
