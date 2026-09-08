package seed

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"gopkg.in/yaml.v3"
)

func TestCompleteSeed(t *testing.T) {
	var configuration spec.Config
	if err := yaml.Unmarshal([]byte(YAML), &configuration); err != nil {
		t.Fatal(err)
	}
	if err := validate.Config(configuration); err != nil {
		t.Fatal(err)
	}
	var document yaml.Node
	if err := yaml.Unmarshal([]byte(YAML), &document); err != nil {
		t.Fatal(err)
	}
	checkFields(t, document.Content[0], reflect.TypeOf(configuration), "seed")
}

func checkFields(t *testing.T, node *yaml.Node, typ reflect.Type, path string) {
	t.Helper()
	switch typ.Kind() {
	case reflect.Struct:
		fields := map[string]*yaml.Node{}
		for i := 0; i < len(node.Content); i += 2 {
			fields[node.Content[i].Value] = node.Content[i+1]
		}
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			name := strings.Split(field.Tag.Get("yaml"), ",")[0]
			child, ok := fields[name]
			if !ok {
				t.Errorf("missing %s.%s", path, name)
				continue
			}
			checkFields(t, child, field.Type, path+"."+name)
		}
	case reflect.Slice:
		for _, child := range node.Content {
			checkFields(t, child, typ.Elem(), path+"[]")
		}
	}
}
