package finder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tagus/crypt/internal/crypt"
)

var (
	ebay = &crypt.Credential{
		Id:          "acct-1",
		Service:     "eBay",
		Description: "electronic auction bay",
		Username:    "beanie_babies123",
		Password:    "mars321",
	}
	amazon = &crypt.Credential{
		Id:          "acct-2",
		Service:     "Amazon Web Services",
		Description: "america's ali baba",
		Email:       "jeff.bezos@amazon.com",
		Password:    "123jupiter",
	}
	google = &crypt.Credential{
		Id:          "acct-3",
		Service:     "Google",
		Description: "big brother",
		Email:       "obrien@google.com",
		Password:    "doublethink",
	}
	google2 = &crypt.Credential{
		Id:          "acct-4",
		Service:     "Google Again",
		Description: "big brother 2",
		Email:       "winston@google.com",
		Password:    "newspeak",
	}

	cr = &crypt.Crypt{
		Credentials: map[string]*crypt.Credential{
			"acct-2": amazon,
			"acct-1": ebay,
			"acct-3": google,
			"acct-4": google2,
		},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
)

func TestFinder_FilterForNonExistentService(t *testing.T) {
	finder, err := New(cr)
	assert.NoError(t, err)
	svcs, err := finder.Filter("non-existent")
	assert.NoError(t, err)
	assert.Empty(t, svcs)
}

func TestFinder_Filter(t *testing.T) {
	finder, err := New(cr)
	assert.NoError(t, err)

	svcs, err := finder.Filter("google")
	assert.NoError(t, err)
	assert.Len(t, svcs, 2)
}

func TestFinder_Find(t *testing.T) {
	finder, err := New(cr)
	assert.NoError(t, err)

	svc, err := finder.Find("ebay")
	assert.NoError(t, err)
	assert.Equal(t, "acct-1", svc.Id)
	assert.Equal(t, "eBay", svc.Service)
}
