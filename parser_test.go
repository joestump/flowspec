package flowspec

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	input := []byte(`
name: test-workflow
trigger:
  cron: "0 8 * * *"
steps:
  - name: step-one
    agent: worker
    config:
      key: value
  - name: step-two
    agent: processor
    input: "$prev.output"
`)

	wf, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if wf.Name != "test-workflow" {
		t.Errorf("Name = %q, want %q", wf.Name, "test-workflow")
	}
	if wf.Trigger.Cron != "0 8 * * *" {
		t.Errorf("Trigger.Cron = %q, want %q", wf.Trigger.Cron, "0 8 * * *")
	}
	if len(wf.Steps) != 2 {
		t.Fatalf("len(Steps) = %d, want 2", len(wf.Steps))
	}
	if wf.Steps[0].Agent != "worker" {
		t.Errorf("Steps[0].Agent = %q, want %q", wf.Steps[0].Agent, "worker")
	}
	if wf.Steps[1].Input != "$prev.output" {
		t.Errorf("Steps[1].Input = %q, want %q", wf.Steps[1].Input, "$prev.output")
	}
}

func TestParseFile(t *testing.T) {
	wf, err := ParseFile("examples/morning-briefing.yaml")
	if err != nil {
		t.Fatalf("ParseFile() error: %v", err)
	}

	if wf.Name != "morning-briefing" {
		t.Errorf("Name = %q, want %q", wf.Name, "morning-briefing")
	}
	if len(wf.Steps) != 4 {
		t.Fatalf("len(Steps) = %d, want 4", len(wf.Steps))
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("nonexistent.yaml")
	if err == nil {
		t.Fatal("ParseFile() expected error for missing file")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		wf      *Workflow
		wantErr bool
	}{
		{
			name: "valid workflow",
			wf: &Workflow{
				Name:  "test",
				Steps: []Step{{Name: "s1", Agent: "a1"}},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			wf: &Workflow{
				Steps: []Step{{Name: "s1", Agent: "a1"}},
			},
			wantErr: true,
		},
		{
			name: "no steps",
			wf: &Workflow{
				Name:  "test",
				Steps: []Step{},
			},
			wantErr: true,
		},
		{
			name: "step missing name",
			wf: &Workflow{
				Name:  "test",
				Steps: []Step{{Agent: "a1"}},
			},
			wantErr: true,
		},
		{
			name: "step missing agent and parallel",
			wf: &Workflow{
				Name:  "test",
				Steps: []Step{{Name: "s1"}},
			},
			wantErr: true,
		},
		{
			name: "duplicate step names",
			wf: &Workflow{
				Name: "test",
				Steps: []Step{
					{Name: "s1", Agent: "a1"},
					{Name: "s1", Agent: "a2"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid on_error",
			wf: &Workflow{
				Name:  "test",
				Steps: []Step{{Name: "s1", Agent: "a1", OnError: "explode"}},
			},
			wantErr: true,
		},
		{
			name: "valid on_error",
			wf: &Workflow{
				Name:  "test",
				Steps: []Step{{Name: "s1", Agent: "a1", OnError: "retry"}},
			},
			wantErr: false,
		},
		{
			name: "parallel steps",
			wf: &Workflow{
				Name: "test",
				Steps: []Step{
					{
						Name: "fan-out",
						Parallel: []Step{
							{Name: "a", Agent: "worker-a"},
							{Name: "b", Agent: "worker-b"},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.wf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseAndValidateExamples(t *testing.T) {
	examples, err := os.ReadDir("examples")
	if err != nil {
		t.Fatalf("ReadDir(examples) error: %v", err)
	}

	for _, e := range examples {
		t.Run(e.Name(), func(t *testing.T) {
			wf, err := ParseFile("examples/" + e.Name())
			if err != nil {
				t.Fatalf("ParseFile() error: %v", err)
			}
			if err := Validate(wf); err != nil {
				t.Fatalf("Validate() error: %v", err)
			}
		})
	}
}
