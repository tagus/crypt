package repos

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCredential_Clone(t *testing.T) {
	ts := time.Now()
	cr := &Credential{
		ID:          "id",
		Service:     "service",
		Domains:     []string{"domain"},
		Email:       "email",
		Username:    "username",
		Password:    "password",
		Description: "description",
		Details: &Details{
			SecurityQuestions: []SecurityQuestion{
				{
					Question: "question",
					Answer:   "answer",
				},
			},
		},
		Tags:          []string{"tag"},
		CreatedAt:     ts,
		UpdatedAt:     ts,
		AccessedAt:    &ts,
		AccessedCount: 1,
		Version:       1,
	}

	cloned := cr.Clone()
	require.Equal(t, cr, cloned)
}

func TestCredential_CloneWithNilValues(t *testing.T) {
	ts := time.Now()
	cr := &Credential{
		ID:            "id",
		Service:       "service",
		Domains:       []string{"domain"},
		Email:         "email",
		Username:      "username",
		Password:      "password",
		Description:   "description",
		Details:       nil,
		Tags:          []string{"tag"},
		CreatedAt:     ts,
		UpdatedAt:     ts,
		AccessedAt:    &ts,
		AccessedCount: 1,
		Version:       1,
	}

	cloned := cr.Clone()
	require.Equal(t, cr, cloned)
}
