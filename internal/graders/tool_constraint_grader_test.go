package graders

import (
	"context"
	"testing"

	"github.com/spboyer/waza/internal/models"
)

func TestToolConstraintGrader_RequiresAtLeastOneConstraint(t *testing.T) {
	_, err := NewToolConstraintGrader("empty", ToolConstraintGraderParams{})
	if err == nil {
		t.Fatal("expected error for empty params")
	}
}

func TestToolConstraintGrader_ExpectTools_Pass(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		ExpectTools: []string{"bash", "edit"},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			ToolsUsed: []string{"bash", "edit", "view"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Feedback)
	}
	if result.Score != 1.0 {
		t.Errorf("expected score 1.0, got %f", result.Score)
	}
}

func TestToolConstraintGrader_ExpectTools_Fail(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		ExpectTools: []string{"bash", "edit"},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			ToolsUsed: []string{"bash", "view"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Error("expected fail, got pass")
	}
	if result.Score != 0.5 {
		t.Errorf("expected score 0.5, got %f", result.Score)
	}
}

func TestToolConstraintGrader_RejectTools_Pass(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		RejectTools: []string{"create_file"},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			ToolsUsed: []string{"bash", "edit"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Feedback)
	}
}

func TestToolConstraintGrader_RejectTools_Fail(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		RejectTools: []string{"create_file", "delete"},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			ToolsUsed: []string{"bash", "create_file"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Error("expected fail, got pass")
	}
	// 1 of 2 reject tools was used → 1 pass, 1 fail → 0.5
	if result.Score != 0.5 {
		t.Errorf("expected score 0.5, got %f", result.Score)
	}
}

func TestToolConstraintGrader_MaxTurns_Pass(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		MaxTurns: 15,
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			TotalTurns: 10,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Feedback)
	}
}

func TestToolConstraintGrader_MaxTurns_Fail(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		MaxTurns: 5,
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			TotalTurns: 10,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Error("expected fail, got pass")
	}
}

func TestToolConstraintGrader_MaxTurns_Boundary(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		MaxTurns: 10,
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			TotalTurns: 10,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass at boundary, got fail: %s", result.Feedback)
	}
}

func TestToolConstraintGrader_MaxTokens_Pass(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		MaxTokens: 5000,
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			TokensTotal: 3000,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Feedback)
	}
}

func TestToolConstraintGrader_MaxTokens_Fail(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		MaxTokens: 5000,
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			TokensTotal: 8000,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Error("expected fail, got pass")
	}
}

func TestToolConstraintGrader_AllConstraints_Pass(t *testing.T) {
	g, err := NewToolConstraintGrader("full", ToolConstraintGraderParams{
		ExpectTools: []string{"bash", "edit"},
		RejectTools: []string{"create_file"},
		MaxTurns:    15,
		MaxTokens:   5000,
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			ToolsUsed:   []string{"bash", "edit", "view"},
			TotalTurns:  10,
			TokensTotal: 4000,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Feedback)
	}
	if result.Score != 1.0 {
		t.Errorf("expected score 1.0, got %f", result.Score)
	}
}

func TestToolConstraintGrader_AllConstraints_PartialFail(t *testing.T) {
	g, err := NewToolConstraintGrader("partial", ToolConstraintGraderParams{
		ExpectTools: []string{"bash", "edit"},
		RejectTools: []string{"create_file"},
		MaxTurns:    15,
		MaxTokens:   5000,
	})
	if err != nil {
		t.Fatal(err)
	}

	// bash used, edit missing, create_file used, turns ok, tokens over
	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			ToolsUsed:   []string{"bash", "create_file"},
			TotalTurns:  10,
			TokensTotal: 8000,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Error("expected fail, got pass")
	}
	// 6 total checks: 2 expect + 1 reject + 1 turns + 1 tokens = 5? No:
	// expect_tools: bash(pass) + edit(fail) = 2 checks
	// reject_tools: create_file(fail) = 1 check
	// max_turns: pass = 1 check
	// max_tokens: fail = 1 check
	// total = 5 checks, 2 passed, score = 2/5 = 0.4
	if result.Score != 0.4 {
		t.Errorf("expected score 0.4, got %f", result.Score)
	}
}

func TestToolConstraintGrader_NilSession(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		MaxTurns: 10,
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Error("expected fail with nil session")
	}
	if result.Score != 0.0 {
		t.Errorf("expected score 0.0, got %f", result.Score)
	}
}

func TestToolConstraintGrader_Kind(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		MaxTurns: 10,
	})
	if err != nil {
		t.Fatal(err)
	}

	if g.Kind() != models.GraderKindToolConstraint {
		t.Errorf("expected kind %s, got %s", models.GraderKindToolConstraint, g.Kind())
	}
	if g.Name() != "test" {
		t.Errorf("expected name 'test', got '%s'", g.Name())
	}
}

func TestToolConstraintGrader_EmptyToolsUsed(t *testing.T) {
	g, err := NewToolConstraintGrader("test", ToolConstraintGraderParams{
		ExpectTools: []string{"bash"},
		RejectTools: []string{"delete"},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := g.Grade(context.Background(), &Context{
		Session: &models.SessionDigest{
			ToolsUsed: []string{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Error("expected fail — expected tool not found in empty list")
	}
	// expect: bash missing (fail), reject: delete not found (pass) → 1/2 = 0.5
	if result.Score != 0.5 {
		t.Errorf("expected score 0.5, got %f", result.Score)
	}
}
