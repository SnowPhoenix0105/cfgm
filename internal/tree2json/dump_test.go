package tree2json

import (
	"github.com/SnowPhoenix0105/cfgm/internal/obj2tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAll(t *testing.T) {
	type Person struct {
		Name   string
		Age    int
		Height float64
	}
	type Date struct {
		Year  int
		Month string
	}
	type Family struct {
		Father   Person
		Mother   Person
		Children []*Person
		Travel   map[string]*Date
		Money    map[string]float64
	}

	obj := Family{
		Father: Person{
			Name:   "Bob",
			Age:    46,
			Height: 1.76,
		},
		Mother: Person{
			Name:   "Emily",
			Age:    46,
			Height: 1.70,
		},
		Children: []*Person{
			&Person{
				Name:   "Tom",
				Age:    21,
				Height: 1.81,
			},
			&Person{
				Name:   "Mary",
				Age:    19,
				Height: 1.73,
			},
			&Person{
				Name:   "Child",
				Age:    0,
				Height: 0.50,
			},
			nil,
		}[:3],
		Travel: map[string]*Date{
			"__prototype__": &Date{
				Year:  2000,
				Month: "January",
			},
			"Beijing": &Date{
				Year:  2021,
				Month: "February",
			},
		},
		Money: map[string]float64{
			"Tuition":            -2000,
			"Salary":             +10000,
			"GovernmentBenefits": +2000,
		},
	}
	root, err := obj2tree.BuildFrom(&obj, 1)
	assert.Nil(t, err)
	t.Log(DumpToString(root))
	assert.Equal(t, 1, len(obj.Travel))
	assert.Equal(t, 2, len(obj.Children))
}
