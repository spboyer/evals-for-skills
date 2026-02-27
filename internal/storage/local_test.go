package storage

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spboyer/waza/internal/models"
	"github.com/spboyer/waza/internal/projectconfig"
)

func makeOutcome(runID, skill, model string, passed, total int) *models.EvaluationOutcome {
	return &models.EvaluationOutcome{
		RunID:       runID,
		SkillTested: skill,
		BenchName:   "test-bench",
		Timestamp:   time.Date(2026, 2, 27, 12, 0, 0, 0, time.UTC),
		Setup:       models.OutcomeSetup{ModelID: model},
		Digest: models.OutcomeDigest{
			TotalTests:     total,
			Succeeded:      passed,
			Failed:         total - passed,
			AggregateScore: 0.85,
		},
		Measures: map[string]models.MeasureResult{
			"accuracy": {Value: 0.9},
		},
	}
}

func TestLocalStore_UploadAndDownload(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir)
	ctx := context.Background()

	outcome := makeOutcome("run-1", "code-gen", "gpt-4o", 8, 10)

	if err := store.Upload(ctx, outcome); err != nil {
		t.Fatalf("Upload() error: %v", err)
	}

	// File should exist on disk.
	path := filepath.Join(dir, "run-1.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file at %s: %v", path, err)
	}

	got, err := store.Download(ctx, "run-1")
	if err != nil {
		t.Fatalf("Download() error: %v", err)
	}

	if got.RunID != "run-1" {
		t.Errorf("RunID = %q, want %q", got.RunID, "run-1")
	}
	if got.SkillTested != "code-gen" {
		t.Errorf("SkillTested = %q, want %q", got.SkillTested, "code-gen")
	}
}

func TestLocalStore_UploadEmptyRunID(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir)

	outcome := makeOutcome("", "s", "m", 1, 1)
	err := store.Upload(context.Background(), outcome)
	if err == nil {
		t.Fatal("Upload() should error on empty RunID")
	}
}

func TestLocalStore_DownloadNotFound(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir)

	_, err := store.Download(context.Background(), "nonexistent")
	if err != ErrNotFound {
		t.Fatalf("Download() error = %v, want ErrNotFound", err)
	}
}

func TestLocalStore_List(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir)
	ctx := context.Background()

	o1 := makeOutcome("run-a", "skill-a", "gpt-4o", 5, 10)
	o1.Timestamp = time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	o2 := makeOutcome("run-b", "skill-b", "claude-sonnet", 9, 10)
	o2.Timestamp = time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)
	o3 := makeOutcome("run-c", "skill-a", "gpt-4o", 10, 10)
	o3.Timestamp = time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)

	for _, o := range []*models.EvaluationOutcome{o1, o2, o3} {
		if err := store.Upload(ctx, o); err != nil {
			t.Fatalf("Upload() error: %v", err)
		}
	}

	// List all.
	all, err := store.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("List() returned %d, want 3", len(all))
	}
	// Newest first.
	if all[0].RunID != "run-c" {
		t.Errorf("first result = %q, want run-c", all[0].RunID)
	}

	// Filter by skill.
	bySkill, err := store.List(ctx, ListOptions{Skill: "skill-a"})
	if err != nil {
		t.Fatalf("List(skill) error: %v", err)
	}
	if len(bySkill) != 2 {
		t.Fatalf("List(skill=skill-a) returned %d, want 2", len(bySkill))
	}

	// Filter by model.
	byModel, err := store.List(ctx, ListOptions{Model: "claude-sonnet"})
	if err != nil {
		t.Fatalf("List(model) error: %v", err)
	}
	if len(byModel) != 1 {
		t.Fatalf("List(model=claude-sonnet) returned %d, want 1", len(byModel))
	}

	// Limit.
	limited, err := store.List(ctx, ListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("List(limit) error: %v", err)
	}
	if len(limited) != 2 {
		t.Fatalf("List(limit=2) returned %d, want 2", len(limited))
	}

	// Since filter.
	since, err := store.List(ctx, ListOptions{Since: time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatalf("List(since) error: %v", err)
	}
	if len(since) != 2 {
		t.Fatalf("List(since=Feb10) returned %d, want 2", len(since))
	}
}

func TestLocalStore_Compare(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir)
	ctx := context.Background()

	o1 := makeOutcome("run-1", "skill-x", "gpt-4o", 6, 10)
	o1.Digest.AggregateScore = 0.6
	o1.Measures = map[string]models.MeasureResult{
		"accuracy": {Value: 0.6},
		"speed":    {Value: 100},
	}
	o2 := makeOutcome("run-2", "skill-x", "gpt-4o", 9, 10)
	o2.Digest.AggregateScore = 0.9
	o2.Measures = map[string]models.MeasureResult{
		"accuracy": {Value: 0.9},
		"speed":    {Value: 80},
	}

	_ = store.Upload(ctx, o1)
	_ = store.Upload(ctx, o2)

	report, err := store.Compare(ctx, "run-1", "run-2")
	if err != nil {
		t.Fatalf("Compare() error: %v", err)
	}

	if report.PassDelta != 30.0 {
		t.Errorf("PassDelta = %v, want 30.0", report.PassDelta)
	}
	wantScoreDelta := 0.3
	if diff := report.ScoreDelta - wantScoreDelta; diff < -0.001 || diff > 0.001 {
		t.Errorf("ScoreDelta = %v, want ~%v", report.ScoreDelta, wantScoreDelta)
	}
	if len(report.Metrics) != 2 {
		t.Fatalf("Metrics has %d entries, want 2", len(report.Metrics))
	}
	if report.Metrics["accuracy"].Delta != 0.3 {
		t.Errorf("accuracy delta = %v, want 0.3", report.Metrics["accuracy"].Delta)
	}
}

func TestLocalStore_LoadFromExistingFiles(t *testing.T) {
	dir := t.TempDir()

	// Write a result file directly (simulating pre-existing data).
	outcome := makeOutcome("preexisting", "my-skill", "model-x", 7, 10)
	data, _ := json.MarshalIndent(outcome, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, "preexisting.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	store := NewLocalStore(dir)
	got, err := store.Download(context.Background(), "preexisting")
	if err != nil {
		t.Fatalf("Download() error: %v", err)
	}
	if got.SkillTested != "my-skill" {
		t.Errorf("SkillTested = %q, want %q", got.SkillTested, "my-skill")
	}
}

func TestLocalStore_EmptyDir(t *testing.T) {
	store := NewLocalStore("")
	results, err := store.List(context.Background(), ListOptions{})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("List() returned %d, want 0", len(results))
	}
}

func TestLocalStore_NonexistentDir(t *testing.T) {
	store := NewLocalStore("/tmp/waza-nonexistent-test-dir-12345")
	results, err := store.List(context.Background(), ListOptions{})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("List() returned %d, want 0", len(results))
	}
}

func TestNewStore_LocalDefault(t *testing.T) {
	dir := t.TempDir()
	store, err := NewStore(nil, dir)
	if err != nil {
		t.Fatalf("NewStore() error: %v", err)
	}
	if _, ok := store.(*LocalStore); !ok {
		t.Error("NewStore(nil) should return *LocalStore")
	}
}

func TestNewStore_AzureNotImplemented(t *testing.T) {
	dir := t.TempDir()
	cfg := &projectconfig.StorageConfig{
		Provider: "azure-blob",
		Enabled:  true,
	}
	_, err := NewStore(cfg, dir)
	if err == nil {
		t.Fatal("NewStore(azure-blob) should return error")
	}
}

func TestNewStore_AzureDisabled(t *testing.T) {
	dir := t.TempDir()
	cfg := &projectconfig.StorageConfig{
		Provider: "azure-blob",
		Enabled:  false,
	}
	store, err := NewStore(cfg, dir)
	if err != nil {
		t.Fatalf("NewStore() error: %v", err)
	}
	if _, ok := store.(*LocalStore); !ok {
		t.Error("NewStore(azure-blob, disabled) should return *LocalStore")
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"path/to/run", "path_to_run"},
		{"run:123", "run_123"},
		{"has spaces", "has_spaces"},
	}
	for _, tt := range tests {
		got := sanitizeFilename(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
