# Tree-sitter-go v0.25.0 Investigation

## Problem Summary

After updating tree-sitter-go from commit 7cb21a6 to v0.25.0 (commit 2346a3a), the parser produces empty root nodes with 0 children, despite successful grammar loading and parsing operations.

## Debug Output

```
Loading grammar for language: go (normalized: go) ✓
Grammar loaded successfully for go ✓
Grammar set on parser for go ✓
Parser obtained for go, starting parse... ✓
Parse successful for go ✓
Root node kind: [EMPTY STRING] ✗  <- Should be "source_file"
Root node has 0 named children ✗    <- Should have package_clause, imports, etc.
```

## Root Cause Analysis

Between commit 7cb21a6 (original) and 2346a3a (v0.25.0), there were major changes:

1. **Commit 1496eb7: "feat: use the new reserved rules API"**
   - Rewrote `src/parser.c` with 1432 changed lines
   - Updated `src/tree_sitter/parser.h` with new API structures
   - Added `TSLanguageMetadata` structure
   - Changed `TSFieldMapSlice` to `TSMapSlice`
   - Added 57 new error test cases

2. **Other significant changes**:
   - `feat: support generic type aliases` (e1076e5)
   - `fix: give index expressions a dynamic precedence of 1` (edea6bf)
   - `feat: expose statement list` (179ca03)

## Current State

I've reverted the WASM binary to the original version (pre-v0.25.0):

- **Current**: `lib/ts.wasm.br` - Original WASM binary with commit 7cb21a6
- **Backup**: `lib/ts.wasm.br.v0.25.0` - Broken v0.25.0 WASM binary
- **grammars.json**: Still references v0.25.0 (metadata only, doesn't affect WASM)

## Testing Required

The original version (7cb21a6) was reported to have stack overflow crashes. We need to test:

1. **Does the original WASM binary still crash?**
   - If NO → Use original version, problem solved
   - If YES → We need a different solution

2. **Why does v0.25.0 produce empty root nodes?**
   - Possible API incompatibility between tree-sitter core and v0.25.0 grammar
   - Possible WASM compilation issue
   - Possible bug in v0.25.0 itself

## Options Moving Forward

### Option 1: Use Original Version (7cb21a6)
- **Pros**: May work without crashes
- **Cons**: Older grammar, misses bug fixes and features

### Option 2: Debug v0.25.0 Issue
- **Pros**: Latest grammar with all fixes
- **Cons**: Requires understanding WASM/tree-sitter internals

### Option 3: Try Intermediate Version
- Test versions between 7cb21a6 and 1496eb7 (before "new reserved rules API")
- Find the last working version before the breaking change

### Option 4: File Upstream Bug Report
- Report empty root node issue to tree-sitter-go
- Wait for fix in newer version

## Commit Timeline (7cb21a6 → 2346a3a)

```
2346a3a - ci: bump tree-sitter/parser-test-action from 2 to 3
1547678 - 0.25.0 (RELEASE TAG)
3f912e9 - chore: generate
179ca03 - feat: expose statement list
e25214e - fix: allow the terminator to be omitted for the last element
edea6bf - fix: give index expressions a dynamic precedence of 1
e1076e5 - feat: support generic type aliases
00a299e - ci: update test failures, use macos-15
93c2bb6 - build: update bindings
1496eb7 - feat: use the new reserved rules API ⚠️ BREAKING CHANGE
c350fa5 - ci: bump actions/checkout from 4 to 5
5e73f47 - Include LICENSE file
7cb21a6 - chore: add FUNDING.yml (ORIGINAL VERSION)
```

## Recommendation

Test the original WASM binary first to see if it works without crashes. If it does, we should use it and investigate v0.25.0 separately. If it still crashes, we need to try Option 3 (intermediate versions) or consider alternative approaches like trying a different tree-sitter wrapper.
