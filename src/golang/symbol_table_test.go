package golang

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSymbolTable(t *testing.T) {
	test := SymTab

	test = PushTable(test, "forBlock")
	require.Equal(t, test.PrevTable.Scope, "main")
	require.Equal(t, test.Scope, "main/forBlock")

	test = PopTable(test)
	require.Equal(t, test.Scope, "main")

	test = PushTable(test, "forBlock")
	test = PushTable(test, "ifBlock")
	require.Equal(t, test.Scope, "main/forBlock/ifBlock")
	require.Equal(t, test.PrevTable.PrevTable.Scope, "main")

	copy := SnapshotTable(SymTab)
	require.Equal(t, copy.Scope, "main/forBlock/ifBlock")
	require.Equal(t, copy.PrevTable.Scope, "main/forBlock")
	require.Equal(t, copy.PrevTable.PrevTable.Scope, "main")
}
