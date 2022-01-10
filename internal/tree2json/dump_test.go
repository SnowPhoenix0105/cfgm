package tree2json

import (
	"github.com/SnowPhoenix0105/cfgm/internal/obj2tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAll(t *testing.T) {
	type Person struct {
		Name   string  `desc:"the name of the person"`
		Age    int     `desc:"the age of the person"`
		Height float64 `desc:"the height of the person (in meter)"`
	}
	type Date struct {
		Year  int
		Month string
	}
	type Family struct {
		Father   Person             `desc:"the man of the family"`
		Mother   Person             `desc:"the hostess of this family"`
		Children []*Person          `desc:"the children of the family"`
		Travel   map[string]*Date   `desc:"places that this family has traveled to"`
		Money    map[string]float64 `desc:"income and expenditure"`
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
			{
				Name:   "Tom",
				Age:    21,
				Height: 1.81,
			},
			{
				Name:   "Mary",
				Age:    19,
				Height: 1.73,
			},
			{
				Name:   "Child",
				Age:    0,
				Height: 0.50,
			},
			nil,
		}[:3],
		Travel: map[string]*Date{
			"__prototype__": {
				Year:  2000,
				Month: "January",
			},
			"Beijing": {
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
