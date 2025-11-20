module.exports = {
  // TypeScript/JavaScript files
  'apps/web/**/*.{ts,tsx,js,jsx}': ['eslint --fix', 'prettier --write'],

  // Go files
  'apps/api/**/*.go': ['gofmt -w', 'goimports -w'],

  // JSON files
  '**/*.json': ['prettier --write'],

  // Markdown files
  '**/*.md': ['prettier --write'],

  // YAML files
  '**/*.{yml,yaml}': ['prettier --write'],
};
