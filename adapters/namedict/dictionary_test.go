package namedict

import (
	"github.com/gissleh/litxap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	res := New("nor", "no", "ta-*mu", "kelnì")
	assert.Equal(t, []string{"nor: ", "no: -r "}, res.(*nameDict).table["nor"])
	assert.Equal(t, []string{"ta.*mu: -ri "}, res.(*nameDict).table["tamuri"])

	entries, err := res.LookupEntries("nor")
	assert.NoError(t, err)
	assert.Equal(t, []litxap.Entry{
		*litxap.ParseEntry("nor: : Custom Name"),
		*litxap.ParseEntry("no: -r: Custom Name"),
	}, entries)

	entries, err = res.LookupEntries("tamul")
	assert.NoError(t, err)
	assert.Equal(t, []litxap.Entry{
		*litxap.ParseEntry("ta.*mu: -l: Custom Name"),
	}, entries)

	entries, err = res.LookupEntries("kelnur")
	assert.NoError(t, err)
	assert.Equal(t, []litxap.Entry{
		*litxap.ParseEntry("kelnì: -ur: Custom Name"),
	}, entries)

	entries, err = res.LookupEntries("neytiriti")
	assert.ErrorIs(t, err, litxap.ErrEntryNotFound)
	assert.Nil(t, entries)
}
