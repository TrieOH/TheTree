const Module = require('node:module')
const path = require('node:path')

// TypeScript 7's native package does not expose the Compiler API yet.
// typescript-eslint still requires that API, so isolate its supported compiler
// version without changing the TypeScript version used by the application.
const eslintTypeScript = require('typescript-eslint-compiler')
const originalLoad = Module._load

Module._load = function load(request, parent, isMain) {
  if (request === 'typescript') {
    return eslintTypeScript
  }

  return originalLoad.call(this, request, parent, isMain)
}

const eslintBin = path.join(
  path.dirname(require.resolve('eslint/package.json')),
  'bin/eslint.js',
)
process.argv = [process.argv[0], eslintBin, ...process.argv.slice(2)]
require(eslintBin)
