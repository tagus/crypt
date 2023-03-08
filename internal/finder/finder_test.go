package finder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tagus/crypt/internal/crypt"
)

var (
	ebay = &crypt.Credential{
		Id:          "acct1",
		Service:     "eBay",
		Description: "electronic auction bay",
		Username:    "beanie_babies123",
		Password:    "mars321",
	}
	amazon = &crypt.Credential{
		Id:          "acct2",
		Service:     "Amazon Web Services",
		Description: "america's ali baba",
		Email:       "jeff.bezos@amazon.com",
		Password:    "123jupiter",
	}
	google = &crypt.Credential{
		Id:          "acct3",
		Service:     "Google",
		Description: "big brother",
		Email:       "obrien@google.com",
		Password:    "doublethink",
	}
	google2 = &crypt.Credential{
		Id:          "acct4",
		Service:     "Google Again",
		Description: "big brother 2",
		Email:       "winston@google.com",
		Password:    "newspeak",
	}
	hnr = &crypt.Credential{
		Id:          "acct5",
		Service:     "H&R",
		Description: "tax  service",
		Email:       "winston@google.com",
		Password:    "newspeak",
	}

	cr = &crypt.Crypt{
		Credentials: map[string]*crypt.Credential{
			"acct1": ebay,
			"acct2": amazon,
			"acct3": google,
			"acct4": google2,
			"acct5": hnr,
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

func TestFinder_FilterByServiceName(t *testing.T) {
	finder, err := New(cr)
	assert.NoError(t, err)

	svcs, err := finder.Filter("google")
	assert.NoError(t, err)
	assert.Len(t, svcs, 2)
}

func TestFinder_FilterById(t *testing.T) {
	finder, err := New(cr)
	assert.NoError(t, err)

	svcs, err := finder.Filter("acct1")
	assert.NoError(t, err)
	assert.Len(t, svcs, 1)

	assert.Equal(t, "eBay", svcs[0].Service)
}

func TestFinder_FilterWithSpecialCharacters(t *testing.T) {
	finder, err := New(cr)
	assert.NoError(t, err)

	svcs, err := finder.Filter("H&R")
	assert.NoError(t, err)
	assert.Len(t, svcs, 1)

	assert.Equal(t, "H&R", svcs[0].Service)
}

func TestFinder_Find(t *testing.T) {
	finder, err := New(cr)
	assert.NoError(t, err)

	svc, err := finder.Find("ebay")
	assert.NoError(t, err)
	assert.Equal(t, "acct1", svc.Id)
	assert.Equal(t, "eBay", svc.Service)
}
