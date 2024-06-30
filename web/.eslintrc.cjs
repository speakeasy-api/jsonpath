module.exports = {
    extends: ['eslint:recommended', 'plugin:react/recommended', 'plugin:react-hooks/recommended', 'plugin:jsx-a11y/recommended'],
    plugins: ['@typescript-eslint', 'react', 'react-hooks', 'jsx-a11y', 'prettier'],
    parser: '@typescript-eslint/parser',
    parserOptions: {
        ecmaVersion: 2021,
        sourceType: 'module',
        ecmaFeatures: {
            jsx: true,
        },
    },
    rules: {
        'prettier/prettier': 'error',
        'react/jsx-uses-react': 'off',
        'react/react-in-jsx-scope': 'off',
        "no-unused-vars": "off",
        "no-undef": "off",
        "@typescript-eslint/no-unused-vars": ["error", {
            'argsIgnorePattern': '^_',
            'varsIgnorePattern': '^_',
            'caughtErrorsIgnorePattern': '^_'
        }]
    },
};
