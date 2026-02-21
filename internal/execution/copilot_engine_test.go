package execution

import (
	"context"
	"errors"
	"os"
	"testing"

	copilot "github.com/github/copilot-sdk/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCopilotClient struct {
	startErr         error
	stopErr          error
	createSessionErr error
	session          copilotSession

	startCalls  int
	stopCalls   int
	createCalls int
	lastConfig  *copilot.SessionConfig
}

func (c *fakeCopilotClient) Start(ctx context.Context) error {
	c.startCalls++
	return c.startErr
}

func (c *fakeCopilotClient) Stop() error {
	c.stopCalls++
	return c.stopErr
}

func (c *fakeCopilotClient) CreateSession(ctx context.Context, config *copilot.SessionConfig) (copilotSession, error) {
	c.createCalls++
	c.lastConfig = config
	if c.createSessionErr != nil {
		return nil, c.createSessionErr
	}
	return c.session, nil
}

type fakeSession struct {
	id       string
	handlers []copilot.SessionEventHandler
	sendFn   func(context.Context, copilot.MessageOptions) (string, error)
}

func (s *fakeSession) On(handler copilot.SessionEventHandler) func() {
	s.handlers = append(s.handlers, handler)
	return func() {}
}

func (s *fakeSession) Send(ctx context.Context, opts copilot.MessageOptions) (string, error) {
	if s.sendFn != nil {
		return s.sendFn(ctx, opts)
	}
	return "", nil
}

func (s *fakeSession) SessionID() string {
	return s.id
}

func (s *fakeSession) emit(event copilot.SessionEvent) {
	for _, handler := range s.handlers {
		handler(event)
	}
}

func useFakeCopilotClient(t *testing.T, client copilotClient) {
	t.Helper()
	oldFactory := newCopilotClient
	newCopilotClient = func(opts *copilot.ClientOptions) copilotClient {
		return client
	}
	t.Cleanup(func() {
		newCopilotClient = oldFactory
	})
}

func TestCopilotEngine_Initialize(t *testing.T) {
	engine := NewCopilotEngineBuilder("test-model").Build()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := engine.Initialize(ctx)
	require.NoError(t, err)
}

func TestCopilotEngine_SetupResources(t *testing.T) {
	engine := NewCopilotEngineBuilder("test-model").Build()
	engine.workspace = t.TempDir()

	err := engine.setupResources([]ResourceFile{{Path: "data.txt", Content: "value"}})
	require.NoError(t, err)

	content, err := os.ReadFile(engine.workspace + "/data.txt")
	require.NoError(t, err)
	assert.Equal(t, "value", string(content))
}

func TestJoinStrings(t *testing.T) {
	assert.Equal(t, "", joinStrings(nil))
	assert.Equal(t, "abc", joinStrings([]string{"a", "b", "c"}))
}

func TestCopilotEngine_Execute_SetupResourcesErrorRotatesWorkspace(t *testing.T) {
	engine := NewCopilotEngineBuilder("test-model").Build()
	previousWorkspace := t.TempDir()
	engine.workspace = previousWorkspace

	resp, err := engine.Execute(context.Background(), &ExecutionRequest{
		Message: "hello",
		Resources: []ResourceFile{{
			Path:    "/absolute/path.txt",
			Content: "x",
		}},
		TimeoutSec: 1,
	})
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to setup resources")
	assert.Contains(t, engine.oldWorkspaces, previousWorkspace)
	assert.NotEmpty(t, engine.workspace)
	assert.NotEqual(t, previousWorkspace, engine.workspace)

	require.NoError(t, engine.Shutdown(context.Background()))
}

func TestCopilotEngine_Execute_StartError(t *testing.T) {
	oldClient := &fakeCopilotClient{stopErr: errors.New("old stop failed")}
	newClient := &fakeCopilotClient{startErr: errors.New("start failed")}
	useFakeCopilotClient(t, newClient)

	engine := NewCopilotEngineBuilder("test-model").Build()
	engine.client = oldClient

	resp, err := engine.Execute(context.Background(), &ExecutionRequest{Message: "hello", TimeoutSec: 1})
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to start copilot client")
	assert.Equal(t, 1, oldClient.stopCalls)
	assert.Equal(t, 1, newClient.startCalls)
}

func TestCopilotEngine_Execute_CreateSessionError(t *testing.T) {
	client := &fakeCopilotClient{createSessionErr: errors.New("session create failed")}
	useFakeCopilotClient(t, client)

	engine := NewCopilotEngineBuilder("test-model").Build()
	resp, err := engine.Execute(context.Background(), &ExecutionRequest{Message: "hello", TimeoutSec: 1})
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to create session")
	assert.Equal(t, 1, client.startCalls)
	assert.Equal(t, 1, client.createCalls)
}

func TestCopilotEngine_Execute_SendError(t *testing.T) {
	session := &fakeSession{
		id: "session-1",
		sendFn: func(ctx context.Context, opts copilot.MessageOptions) (string, error) {
			return "", errors.New("send failed")
		},
	}
	client := &fakeCopilotClient{session: session}
	useFakeCopilotClient(t, client)

	engine := NewCopilotEngineBuilder("test-model").Build()
	resp, err := engine.Execute(context.Background(), &ExecutionRequest{Message: "hello", TimeoutSec: 1})
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to send prompt")
}

func TestCopilotEngine_Execute_Success(t *testing.T) {
	session := &fakeSession{id: "session-success"}
	session.sendFn = func(ctx context.Context, opts copilot.MessageOptions) (string, error) {
		require.Equal(t, "hello", opts.Prompt)
		delta := "Hello "
		message := "world"
		skillName := "example"
		skillPath := "/tmp/example/SKILL.md"
		session.emit(copilot.SessionEvent{Type: copilot.AssistantMessageDelta, Data: copilot.Data{Content: &delta}})
		session.emit(copilot.SessionEvent{Type: copilot.AssistantMessage, Data: copilot.Data{Content: &message}})
		session.emit(copilot.SessionEvent{Type: copilot.SkillInvoked, Data: copilot.Data{Name: &skillName, Path: &skillPath}})
		session.emit(copilot.SessionEvent{Type: copilot.SessionIdle})
		return "message-1", nil
	}

	client := &fakeCopilotClient{session: session}
	useFakeCopilotClient(t, client)

	engine := NewCopilotEngineBuilder("test-model").Build()
	cwd, err := os.Getwd()
	require.NoError(t, err)

	resp, err := engine.Execute(context.Background(), &ExecutionRequest{
		Message:    "hello",
		SkillPaths: []string{"/custom/skills", "/custom/skills"},
		TimeoutSec: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.True(t, resp.Success)
	assert.Equal(t, "", resp.ErrorMsg)
	assert.Equal(t, "Hello world", resp.FinalOutput)
	assert.Equal(t, "session-success", resp.SessionID)
	require.Len(t, resp.SkillInvocations, 1)
	assert.Equal(t, "example", resp.SkillInvocations[0].Name)
	assert.Equal(t, "/tmp/example/SKILL.md", resp.SkillInvocations[0].Path)

	require.NotNil(t, client.lastConfig)
	assert.Equal(t, "test-model", client.lastConfig.Model)
	require.Len(t, client.lastConfig.SkillDirectories, 2)
	assert.Equal(t, cwd, client.lastConfig.SkillDirectories[0])
	assert.Equal(t, "/custom/skills", client.lastConfig.SkillDirectories[1])
}

func TestCopilotEngine_Execute_Timeout(t *testing.T) {
	session := &fakeSession{id: "session-timeout"}
	session.sendFn = func(ctx context.Context, opts copilot.MessageOptions) (string, error) {
		return "message-timeout", nil
	}

	client := &fakeCopilotClient{session: session}
	useFakeCopilotClient(t, client)

	engine := NewCopilotEngineBuilder("test-model").Build()
	resp, err := engine.Execute(context.Background(), &ExecutionRequest{Message: "hello", TimeoutSec: 0})
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.False(t, resp.Success)
	assert.Contains(t, resp.ErrorMsg, "execution timed out after 0s")
	assert.Equal(t, "session-timeout", resp.SessionID)
}

func TestCopilotEngine_Execute_ContextCanceled(t *testing.T) {
	session := &fakeSession{id: "session-cancel"}
	session.sendFn = func(ctx context.Context, opts copilot.MessageOptions) (string, error) {
		session.emit(copilot.SessionEvent{Type: copilot.SessionError, Data: copilot.Data{Message: nil}})
		return "message-cancel", nil
	}

	client := &fakeCopilotClient{session: session}
	useFakeCopilotClient(t, client)

	engine := NewCopilotEngineBuilder("test-model").Build()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	resp, err := engine.Execute(ctx, &ExecutionRequest{Message: "hello", TimeoutSec: 1})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Contains(t, []string{"execution timed out after 1s", sessionFailedUnknown}, resp.ErrorMsg)
}

func TestCopilotEngine_Execute_SessionError(t *testing.T) {
	session := &fakeSession{id: "session-error"}
	session.sendFn = func(ctx context.Context, opts copilot.MessageOptions) (string, error) {
		session.emit(copilot.SessionEvent{Type: copilot.SessionError, Data: copilot.Data{Message: nil}})
		return "message-error", nil
	}

	client := &fakeCopilotClient{session: session}
	useFakeCopilotClient(t, client)

	engine := NewCopilotEngineBuilder("test-model").Build()
	resp, err := engine.Execute(context.Background(), &ExecutionRequest{Message: "hello", TimeoutSec: 1})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Equal(t, sessionFailedUnknown, resp.ErrorMsg)
}

func TestCopilotEngine_Shutdown_StopsClientAndCleansWorkspaces(t *testing.T) {
	engine := NewCopilotEngineBuilder("test-model").Build()
	client := &fakeCopilotClient{stopErr: errors.New("stop failed")}
	currentWorkspace := t.TempDir()
	oldWorkspace := t.TempDir()

	engine.client = client
	engine.workspace = currentWorkspace
	engine.oldWorkspaces = []string{oldWorkspace}

	err := engine.Shutdown(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, client.stopCalls)
	assert.Nil(t, engine.client)
	assert.Empty(t, engine.workspace)
	assert.Nil(t, engine.oldWorkspaces)

	_, err = os.Stat(currentWorkspace)
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(oldWorkspace)
	assert.True(t, os.IsNotExist(err))
}
