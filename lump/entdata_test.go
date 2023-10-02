package lump

import (
	"bytes"
	"github.com/galaco/bsp/lump/primitive/entities"
	"reflect"
	"testing"
)

func TestEntData_Contents(t *testing.T) {
	sut := EntData{
		Data: []entities.Entity{
			{
				entities.KeyValue{Key: "classname", Value: "worldspawn"},
				entities.KeyValue{Key: "skyname", Value: "sky_dust"},
			},
		},
	}

	if !reflect.DeepEqual(sut.Contents(), sut.Data) {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
}

func TestEntData_FromBytes(t *testing.T) {
	testCases := []struct {
		name     string
		entdata  []byte
		expected []entities.Entity
	}{
		{
			name: "simple",
			entdata: []byte(`{
"classname" "worldspawn"
"skyname" "sky_dust"
}
` + "\x00"),
			expected: []entities.Entity{
				{
					entities.KeyValue{Key: "classname", Value: "worldspawn"},
					entities.KeyValue{Key: "skyname", Value: "sky_dust"},
				},
			},
		},
		{
			name: "multiple entities",
			entdata: []byte(`{
"classname" "worldspawn"
"skyname" "sky_dust"
}
{
"classname" "prop_dynamic"
"name" "bottle_01"
}
` + "\x00"),
			expected: []entities.Entity{
				{
					entities.KeyValue{Key: "classname", Value: "worldspawn"},
					entities.KeyValue{Key: "skyname", Value: "sky_dust"},
				},
				{
					entities.KeyValue{Key: "classname", Value: "prop_dynamic"},
					entities.KeyValue{Key: "name", Value: "bottle_01"},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ent := EntData{}
			err := ent.FromBytes(tc.entdata)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(ent.Data, tc.expected) {
				t.Error("mismatched between expected and actual unmarshalled bytes")
			}
		})
	}
}

func TestEntData_ToBytes(t *testing.T) {
	testCases := []struct {
		name     string
		ents     []entities.Entity
		expected []byte
	}{
		{
			name: "simple",
			ents: []entities.Entity{
				{
					entities.KeyValue{Key: "classname", Value: "worldspawn"},
					entities.KeyValue{Key: "skyname", Value: "sky_dust"},
				},
			},
			expected: []byte(`{
"classname" "worldspawn"
"skyname" "sky_dust"
}
` + "\x00"),
		},
		{
			name: "multiple entities",
			ents: []entities.Entity{
				{
					entities.KeyValue{Key: "classname", Value: "worldspawn"},
					entities.KeyValue{Key: "skyname", Value: "sky_dust"},
				},
				{
					entities.KeyValue{Key: "classname", Value: "prop_dynamic"},
					entities.KeyValue{Key: "name", Value: "bottle_01"},
				},
			},
			expected: []byte(`{
"classname" "worldspawn"
"skyname" "sky_dust"
}
{
"classname" "prop_dynamic"
"name" "bottle_01"
}
` + "\x00"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := EntData{
				Data: tc.ents,
			}
			res, err := l.ToBytes()
			if err != nil {
				t.Error(err)
			}
			if !bytes.Equal(res, tc.expected) {
				t.Errorf("expected %s, got %s", tc.expected, res)
			}
		})
	}
}
