package transform

import (
	"github.com/tinygo-org/tinygo/compileopts"
	"github.com/tinygo-org/tinygo/compiler/llvmutil"
	"tinygo.org/x/go-llvm"
)

// CreateStackSizeLoads replaces internal/task.getGoroutineStackSize calls with
// loads from internal/task.stackSizes that will be updated after linking. This
// way the stack sizes are loaded from a separate section and can easily be
// modified after linking.
func CreateStackSizeLoads(mod llvm.Module, config *compileopts.Config) []string {
	functionMap := map[llvm.Value][]llvm.Value{}
	var functions []llvm.Value // ptrtoint values of functions
	var functionNames []string
	var functionValues []llvm.Value // direct references to functions
	for _, use := range getUses(mod.NamedFunction("internal/task.getGoroutineStackSize")) {
		if use.FirstUse().IsNil() {
			// Apparently this stack size isn't used.
			use.EraseFromParentAsInstruction()
			continue
		}
		ptrtoint := use.Operand(0)
		if _, ok := functionMap[ptrtoint]; !ok {
			functions = append(functions, ptrtoint)
			functionNames = append(functionNames, ptrtoint.Operand(0).Name())
			functionValues = append(functionValues, ptrtoint.Operand(0))
		}
		functionMap[ptrtoint] = append(functionMap[ptrtoint], use)
	}

	if len(functions) == 0 {
		// Nothing to do.
		return nil
	}

	ctx := mod.Context()
	targetData := llvm.NewTargetData(mod.DataLayout())
	defer targetData.Dispose()
	uintptrType := ctx.IntType(targetData.PointerSize() * 8)

	// Create the new global with stack sizes, that will be put in a new section
	// just for itself.
	stackSizesGlobalType := llvm.ArrayType(functions[0].Type(), len(functions))
	stackSizesGlobal := llvm.AddGlobal(mod, stackSizesGlobalType, "internal/task.stackSizes")
	stackSizesGlobal.SetSection(".tinygo_stacksizes")
	defaultStackSizes := make([]llvm.Value, len(functions))
	defaultStackSize := llvm.ConstInt(functions[0].Type(), config.StackSize(), false)
	for i := range defaultStackSizes {
		defaultStackSizes[i] = defaultStackSize
	}
	stackSizesGlobal.SetInitializer(llvm.ConstArray(functions[0].Type(), defaultStackSizes))

	// Add all relevant values to llvm.used (for LTO).
	llvmutil.AppendToGlobal(mod, "llvm.used", append([]llvm.Value{stackSizesGlobal}, functionValues...)...)

	// Replace the calls with loads from the new global with stack sizes.
	irbuilder := ctx.NewBuilder()
	defer irbuilder.Dispose()
	for i, function := range functions {
		for _, use := range functionMap[function] {
			ptr := llvm.ConstGEP(stackSizesGlobalType, stackSizesGlobal, []llvm.Value{
				llvm.ConstInt(ctx.Int32Type(), 0, false),
				llvm.ConstInt(ctx.Int32Type(), uint64(i), false),
			})
			irbuilder.SetInsertPointBefore(use)
			stacksize := irbuilder.CreateLoad(uintptrType, ptr, "stacksize")
			use.ReplaceAllUsesWith(stacksize)
			use.EraseFromParentAsInstruction()
		}
	}

	return functionNames
}
