import js from '@eslint/js';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  js.configs.recommended,
  ...tseslint.configs.recommended,
  { ignores: ['dist', 'node_modules'] },
  {
    files: ['src/**/*.{ts,tsx}'],
    ignores: ['src/**/*.test.{ts,tsx}'],
    rules: {
      'no-restricted-imports': ['error', { patterns: [{ group: ['../*', '../../*', '../../../*'], message: 'Use a stable layer alias instead of reaching across module folders.' }] }],
    },
  },
  {
    files: ['src/shared/**/*.{ts,tsx}'],
    ignores: ['src/**/*.test.{ts,tsx}'],
    rules: {
      'no-restricted-imports': ['error', { patterns: [{ group: ['@app/*', '@pages/*', '@widgets/*', '@features/*', '@entities/*'], message: 'Shared code cannot depend on a higher frontend layer.' }] }],
    },
  },
  {
    files: ['src/entities/**/*.{ts,tsx}'],
    ignores: ['src/**/*.test.{ts,tsx}'],
    rules: {
      'no-restricted-imports': ['error', { patterns: [{ group: ['@app/*', '@pages/*', '@widgets/*', '@features/*'], message: 'Entities may depend only on entities and shared infrastructure.' }] }],
    },
  },
  {
    files: ['src/features/**/*.{ts,tsx}'],
    ignores: ['src/**/*.test.{ts,tsx}'],
    rules: {
      'no-restricted-imports': ['error', { patterns: [{ group: ['@app/*', '@pages/*', '@widgets/*'], message: 'Features cannot depend on composition or route layers.' }] }],
    },
  },
  {
    files: ['src/widgets/**/*.{ts,tsx}'],
    ignores: ['src/**/*.test.{ts,tsx}'],
    rules: {
      'no-restricted-imports': ['error', { patterns: [{ group: ['@app/*', '@pages/*'], message: 'Widgets cannot depend on app composition or route pages.' }] }],
    },
  },
  {
    files: ['src/pages/**/*.{ts,tsx}'],
    ignores: ['src/**/*.test.{ts,tsx}'],
    rules: {
      'no-restricted-imports': ['error', { patterns: [{ group: ['@app/*', '@pages/*'], message: 'Pages compose lower layers and must not depend on app or sibling pages.' }] }],
    },
  },
);
