package cli

import (
	"errors"
	"fmt"
)

var (
	errorInvalidHeaderFormat        = errors.New("header should be in the format \"key:value\"")
	errorInvalidMethod              = errors.New("method not recognized")
	errorHeaderContainsNewlineChars = errors.New("header should not contain newline characters")
	errorFileEmpty                  = errors.New("empty file")
	errorFileUnchanged              = errors.New("no changes")
	errorNoEditorFound              = errors.New("no editor found")

	missingFlagBase  = "expected flag missing: %s"
	missingFlagsBase = "expected flags missing: %s"
	missingArgBase   = "expected arg missing: %s"
	missingArgsBase  = "expected args missing: %s"
)

func errorMissingFlag(flag string) error {
	return errors.New(fmt.Sprintf(missingFlagBase, flag))
}
func errorMissingFlags(flags string) error {
	return errors.New(fmt.Sprintf(missingFlagsBase, flags))
}
func errorMissingArg(arg string) error {
	return errors.New(fmt.Sprintf(missingArgBase, arg))
}
func errorMissingArgs(args string) error {
	return errors.New(fmt.Sprintf(missingArgsBase, args))
}
