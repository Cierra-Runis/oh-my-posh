package segments

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/cache"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"
	"github.com/jandedobbeleer/oh-my-posh/src/template"

	"github.com/stretchr/testify/assert"
)

const (
	WorkingDirRoot = "/home/user/dev/my-app"
)

type testCase struct {
	Case            string
	Template        string
	ExpectedString  string
	UI5YamlFilename string
	WorkingDir      string
	Version         string
	DisplayMode     string
}

func TestUI5Tooling(t *testing.T) {
	cases := []testCase{
		{
			Case:            "1) ui5tooling 2.12.1 - file ui5.yaml present in cwd; DisplayMode = files",
			ExpectedString:  "2.12.1",
			UI5YamlFilename: "ui5.yaml",
			Version:         `2.12.1 (from C:\somewhere\cli\bin\ui5.js)`,
			DisplayMode:     DisplayModeFiles,
		},
		{
			Case:            "2) ui5tooling 2.12.2 - file ui5.yaml present in cwd; default display mode (context)",
			ExpectedString:  "2.12.2",
			UI5YamlFilename: "ui5.yaml",
			Version:         `2.12.2 (from C:\somewhere\cli\bin\ui5.js)`,
		},
		{
			Case:            "3) ui5tooling 2.12.3 - file ui5.yaml present; cwd is sub dir, default display mode (context)",
			ExpectedString:  "2.12.3",
			WorkingDir:      WorkingDirRoot + "/subdir",
			UI5YamlFilename: "ui5.yaml",
			Version:         `2.12.3 (from C:\somewhere\cli\bin\ui5.js)`,
		},
		{
			Case:            "5) ui5tooling 2.12.4 - file ui5-dist.yml present in cwd",
			ExpectedString:  "2.12.4",
			UI5YamlFilename: "ui5-dist.yml",
			Version:         `2.12.4 (from C:\somewhere\cli\bin\ui5.js)`,
			DisplayMode:     DisplayModeFiles,
		},
		{
			Case:           "8) ui5tooling 11.0.0-rc1, no ui5.yaml file but display mode = always",
			Template:       "{{ .Major }}",
			ExpectedString: "11",
			Version:        `11.0.0-rc1 (from C:\somewhere\cli\bin\ui5.js)`,
			DisplayMode:    DisplayModeAlways,
		},
	}

	for _, tc := range cases {
		params := &mockedLanguageParams{
			cmd:           "ui5",
			versionParam:  "--version",
			versionOutput: tc.Version,
			extension:     UI5ToolingYamlPattern,
		}
		env, props := getMockedLanguageEnv(params)

		if tc.DisplayMode == "" {
			tc.DisplayMode = DisplayModeContext
		}

		props[DisplayMode] = tc.DisplayMode

		ui5tooling := &UI5Tooling{}
		ui5tooling.Init(props, env)

		err := mockFilePresence(&tc, env)

		if err != nil {
			t.Fail()
		}

		if tc.Template == "" {
			tc.Template = ui5tooling.Template()
		}

		// this is needed to build the version URL as before renderTemplate, the template is not initialized
		env.On("Shell").Return("foo")
		template.Cache = &cache.Template{}
		template.Init(env, nil, nil)

		failMsg := fmt.Sprintf("Failed in case: %s", tc.Case)
		assert.True(t, ui5tooling.Enabled(), failMsg)
		assert.Equal(t, tc.ExpectedString, renderTemplate(env, tc.Template, ui5tooling), failMsg)
	}
}

func mockFilePresence(tc *testCase, env *mock.Environment) error {
	for _, f := range []string{UI5ToolingYamlPattern} {
		match, err := filepath.Match(f, tc.UI5YamlFilename)

		if err != nil {
			return err
		}

		if match {
			if tc.DisplayMode == DisplayModeFiles && tc.WorkingDir == WorkingDirRoot { //nolint: gocritic
				env.On("HasFiles", f).Return(true)
				env.On("HasFileInParentDirs", f, uint(4)).Return(false)
				// mode context, working dir != working dir root
			} else if tc.DisplayMode == DisplayModeContext {
				env.On("HasFileInParentDirs", f, uint(4)).Return(false)
				env.On("HasFiles", f).Return(true)
			} else {
				env.On("HasFileInParentDirs", f, uint(4)).Return(false)
				env.On("HasFiles", f).Return(false)
			}
		} else {
			env.On("HasFileInParentDirs", f, uint(4)).Return(false)
			env.On("HasFiles", f).Return(false)
		}
	}

	return nil
}
