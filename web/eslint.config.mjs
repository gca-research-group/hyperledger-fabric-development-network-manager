import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { fixupConfigRules } from '@eslint/compat';
import js from '@eslint/js';
import { FlatCompat } from '@eslint/eslintrc';
import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';
import globals from 'globals';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all,
});

export default tseslint.config(
  {
    ignores: ['projects/**/*'],
  },
  eslint.configs.recommended,
  ...tseslint.configs.recommendedTypeChecked,
  eslintPluginPrettierRecommended,
  ...fixupConfigRules(
    compat.extends(
      //'eslint:recommended',
      //'plugin:@typescript-eslint/recommended',
      'plugin:@angular-eslint/recommended',
      'plugin:@angular-eslint/template/process-inline-templates',
      'plugin:import/typescript',
      'plugin:import/recommended',
      'plugin:import/typescript',
    ),
  ),
  {
    languageOptions: {
      globals: {
        ...globals.es2015,
      },
      ecmaVersion: 5,
      sourceType: 'module',
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
  },
  {
    files: ['**/*.ts'],
    rules: {
      '@angular-eslint/directive-selector': [
        'error',
        {
          type: 'attribute',
          prefix: 'app',
          style: 'camelCase',
        },
      ],

      '@angular-eslint/component-selector': [
        'error',
        {
          type: 'element',
          prefix: 'app',
          style: 'kebab-case',
        },
      ],
      '@typescript-eslint/no-useless-constructor': 'error',
      //'@typescript-eslint/semi': 'error',
      'prettier/prettier': 'error',
      eqeqeq: 'error',

      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
        },
      ],
      'no-unused-vars': 'off',
      'import/no-named-as-default-member': 'off',
      'no-trailing-spaces': 'error',
      indent: 'off',
      'linebreak-style': 'off',
      'eol-last': ['error', 'always'],
      'no-multiple-empty-lines': [
        'error',
        {
          max: 1,
        },
      ],
      'prefer-const': 'error',
      'no-var': 'error',
      'import/no-unresolved': 'off',
      'import/order': [
        'error',
        {
          groups: [
            ['builtin', 'external'],
            ['parent', 'sibling', 'index', 'unknown'],
            ['internal'],
          ],

          pathGroups: [
            {
              pattern: '@angular/**',
              group: 'external',
              position: 'after',
            },
            {
              pattern: 'src/**',
              group: 'internal',
            },
            {
              pattern: '@app/**',
              group: 'external',
              position: 'after',
            },
          ],

          pathGroupsExcludedImportTypes: ['internal'],
          'newlines-between': 'always',

          alphabetize: {
            order: 'asc',
            caseInsensitive: true,
          },
        },
      ],
    },
  },
  // ...compat
  //   .extends(
  //     'plugin:@angular-eslint/template/recommended',
  //     'plugin:@angular-eslint/template/accessibility',
  //   )
  //   .map(config => ({
  //     ...config,
  //     files: ['**/*.html'],
  //   }))
);
