package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// GenerateCompletion outputs shell completion scripts for the specified shell.
func GenerateCompletion(shell string) error {
	return generateCompletion(os.Stdout, shell)
}

func generateCompletion(w io.Writer, shell string) error {
	binaryName := filepath.Base(os.Args[0])

	switch shell {
	case "bash":
		return generateBashCompletion(w, binaryName)
	case "zsh":
		return generateZshCompletion(w, binaryName)
	case "fish":
		return generateFishCompletion(w, binaryName)
	case "powershell":
		return generatePowerShellCompletion(w, binaryName)
	default:
		return fmt.Errorf("unsupported shell: %s (supported: bash, zsh, fish, powershell)", shell)
	}
}

func generateBashCompletion(w io.Writer, name string) error {
	script := fmt.Sprintf(`# Bash completion for %[1]s
_%[1]s_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local opts="-login -codex-login -claude-login -qwen-login -iflow-login -iflow-cookie -no-browser -antigravity-login -project_id -config -vertex-import -completion -help"
    
    COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
    return 0
}

complete -F _%[1]s_completions %[1]s
`, name)
	_, err := fmt.Fprint(w, script)
	return err
}

func generateZshCompletion(w io.Writer, name string) error {
	script := fmt.Sprintf(`#compdef %[1]s

__%[1]s_complete() {
    local -a opts
    opts=(
        '-login[Login Google Account]'
        '-codex-login[Login to Codex using OAuth]'
        '-claude-login[Login to Claude using OAuth]'
        '-qwen-login[Login to Qwen using OAuth]'
        '-iflow-login[Login to iFlow using OAuth]'
        '-iflow-cookie[Login to iFlow using Cookie]'
        '-no-browser[Do not open browser automatically for OAuth]'
        '-antigravity-login[Login to Antigravity using OAuth]'
        '-project_id[Project ID (Gemini only, not required)]:project id:'
        '-config[Configure File Path]:config file:_files'
        '-vertex-import[Import Vertex service account key JSON file]:json file:_files -g "*.json"'
        '-completion[Generate shell completion script]:shell:(bash zsh fish powershell)'
        '-help[Show help]'
    )
    _describe 'options' opts
}

compdef __%[1]s_complete %[1]s
`, name)
	_, err := fmt.Fprint(w, script)
	return err
}

func generateFishCompletion(w io.Writer, name string) error {
	script := fmt.Sprintf(`# Fish completion for %[1]s

complete -c %[1]s -f
complete -c %[1]s -s login -d 'Login Google Account'
complete -c %[1]s -l codex-login -d 'Login to Codex using OAuth'
complete -c %[1]s -l claude-login -d 'Login to Claude using OAuth'
complete -c %[1]s -l qwen-login -d 'Login to Qwen using OAuth'
complete -c %[1]s -l iflow-login -d 'Login to iFlow using OAuth'
complete -c %[1]s -l iflow-cookie -d 'Login to iFlow using Cookie'
complete -c %[1]s -l no-browser -d 'Do not open browser automatically for OAuth'
complete -c %[1]s -l antigravity-login -d 'Login to Antigravity using OAuth'
complete -c %[1]s -l project_id -d 'Project ID (Gemini only)' -r
complete -c %[1]s -l config -d 'Configure File Path' -r -F
complete -c %[1]s -l vertex-import -d 'Import Vertex service account key JSON file' -r -F
complete -c %[1]s -l completion -d 'Generate shell completion script' -r -a 'bash zsh fish powershell'
complete -c %[1]s -s h -l help -d 'Show help'
`, name)
	_, err := fmt.Fprint(w, script)
	return err
}

func generatePowerShellCompletion(w io.Writer, name string) error {
	script := fmt.Sprintf(`# PowerShell completion for %[1]s

Register-ArgumentCompleter -Native -CommandName %[1]s -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    $completions = @(
        [CompletionResult]::new('-login', '-login', [CompletionResultType]::ParameterName, 'Login Google Account')
        [CompletionResult]::new('-codex-login', '-codex-login', [CompletionResultType]::ParameterName, 'Login to Codex using OAuth')
        [CompletionResult]::new('-claude-login', '-claude-login', [CompletionResultType]::ParameterName, 'Login to Claude using OAuth')
        [CompletionResult]::new('-qwen-login', '-qwen-login', [CompletionResultType]::ParameterName, 'Login to Qwen using OAuth')
        [CompletionResult]::new('-iflow-login', '-iflow-login', [CompletionResultType]::ParameterName, 'Login to iFlow using OAuth')
        [CompletionResult]::new('-iflow-cookie', '-iflow-cookie', [CompletionResultType]::ParameterName, 'Login to iFlow using Cookie')
        [CompletionResult]::new('-no-browser', '-no-browser', [CompletionResultType]::ParameterName, 'Do not open browser automatically for OAuth')
        [CompletionResult]::new('-antigravity-login', '-antigravity-login', [CompletionResultType]::ParameterName, 'Login to Antigravity using OAuth')
        [CompletionResult]::new('-project_id', '-project_id', [CompletionResultType]::ParameterName, 'Project ID (Gemini only)')
        [CompletionResult]::new('-config', '-config', [CompletionResultType]::ParameterName, 'Configure File Path')
        [CompletionResult]::new('-vertex-import', '-vertex-import', [CompletionResultType]::ParameterName, 'Import Vertex service account key JSON file')
        [CompletionResult]::new('-completion', '-completion', [CompletionResultType]::ParameterName, 'Generate shell completion script')
        [CompletionResult]::new('-help', '-help', [CompletionResultType]::ParameterName, 'Show help')
    )

    $completions | Where-Object { $_.CompletionText -like "$wordToComplete*" }
}
`, name)
	_, err := fmt.Fprint(w, script)
	return err
}
