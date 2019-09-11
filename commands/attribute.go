package commands
import (
	"fmt"
	"strings"

	"github.com/MichaelMure/git-bug/cache"
	"github.com/MichaelMure/git-bug/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

func runAttribute(cmd *cobra.Command, args []string) error {
	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	b, args, err := _select.ResolveBug(backend, args)
	if err != nil {
		return err
	}

	snap := b.Snapshot()

	for _, a := range snap.Attributes {
			fmt.Printf("%s:\t%s\n", a.Name(), a.Value())
	}

	return nil
}

var attributeCmd = &cobra.Command{
	Use:     "attr [<id>]",
	Short:   "Display, add or remove labels to/from a bug.",
	PreRunE: loadRepo,
	RunE:    runAttribute,
}

func runEditAttribute(cmd* cobra.Command, args []string, set bool) error {
	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)
	b, args, err := _select.ResolveBug(backend, args)
	if err != nil {
		return err
	}

	if set {
		value:=args[1];
		if strings.HasPrefix(value,"link:") {
			target, terr:=backend.ResolveBugPrefix(value)
			if terr != nil {
				return err
			}
			value=target.Id().String()
		}
		_, err = b.EditAttribute(args[0],args[1],set)
	} else {
		_, err = b.EditAttribute(args[0],"",set)
	}
	if err != nil {
		return err
	}
	return b.Commit()
}

func runAttributeSet(cmd *cobra.Command, args []string) error {
	return runEditAttribute(cmd,args,true)
}
func runAttributeUnset(cmd *cobra.Command, args []string) error {
	return runEditAttribute(cmd,args,false)
}

var attributeSetCmd = &cobra.Command {
	Use:	"set [<id>] <name> <value>",
	Short:	"Set an attribute to a bug.",
	PreRunE: loadRepo,
	RunE:	runAttributeSet,
}
var attributeUnsetCmd = &cobra.Command {
	Use:	"unset [<id>] <name> [<value>]",
	Short:	"Unset an attribute to a bug.",
	PreRunE: loadRepo,
	RunE:	runAttributeUnset,
}

func init() {
	RootCmd.AddCommand(attributeCmd)
	labelCmd.Flags().SortFlags = false
	attributeCmd.AddCommand(attributeSetCmd)
	attributeCmd.AddCommand(attributeUnsetCmd)
}

