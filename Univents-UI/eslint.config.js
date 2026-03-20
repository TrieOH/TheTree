//  @ts-check

import { tanstackConfig } from '@tanstack/eslint-config'
import tsEslint from '@typescript-eslint/eslint-plugin'
import tsParser from '@typescript-eslint/parser'

export default [
  ...tanstackConfig,
  {
    files: ['**/*.ts', '**/*.tsx'],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
    plugins: {
      '@typescript-eslint': tsEslint,
    },
    rules: {
      ...tsEslint.configs.recommended.rules,
      ...tsEslint.configs['recommended-type-checked'].rules,
      ...tsEslint.configs['strict-type-checked'].rules,
      ...tsEslint.configs['stylistic-type-checked'].rules,

      // Additional type-aware rules or overrides
      '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
      '@typescript-eslint/restrict-template-expressions': ['error', { allowNumber: true, allowBoolean: true }],
      '@typescript-eslint/consistent-type-imports': 'error',
      '@typescript-eslint/no-import-type-side-effects': 'error',
      '@typescript-eslint/consistent-type-exports': 'error',
      '@typescript-eslint/no-floating-promises': 'error',
      '@typescript-eslint/no-misused-promises': 'error',
      '@typescript-eslint/require-await': 'off'
    },
  },
  {
    rules: {
      'import/no-cycle': 'error',
      'import/order': 'error',
      'sort-imports': 'off',
      '@typescript-eslint/array-type': 'off',
      '@typescript-eslint/require-await': 'error',
      'pnpm/json-enforce-catalog': 'off',
    },
  },
  {
    ignores: ['eslint.config.js', 'prettier.config.js', '.content-collections/'],
  },
  {
    files: ['src/routes/**/*.{ts,tsx}'],
    rules: {
      '@typescript-eslint/only-throw-error': 'off',
    },
  }
]
