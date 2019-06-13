package subcommands

import (
	"fmt"
)

type CompletionCommand struct {
	Bash BashCompletionCommand `command:"bash" description:"bash completion script"`
}

type BashCompletionCommand struct {
}

func (c *BashCompletionCommand) Execute(args []string) error {
	fmt.Printf(`# This script allows Pcf to do autocompletion via Bash.
# Add the following to your .bashrc file, making sure that the path matches your system:
# eval "$(/path/to/pcf completion bash)"

_complete_pcf() {
  args=("${COMP_WORDS[@]:1:$COMP_CWORD}") # Skip first arg
  local IFS=$'\n' # Split into lines
  COMPREPLY=($(GO_FLAGS_COMPLETION=1 ${COMP_WORDS[0]} "${args[@]}"))
  return 0
}

complete -F _complete_pcf pcf
`)
	return nil
}
