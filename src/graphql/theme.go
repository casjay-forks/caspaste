// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

// GraphiQL theming per AI.md PART 14
// Provides light and dark theme CSS for GraphiQL
package graphql

// GraphiQLLightCSS is the light theme CSS for GraphiQL
const GraphiQLLightCSS = `
body {
  margin: 0;
  padding: 0;
  height: 100vh;
  overflow: hidden;
}

#graphiql {
  height: 100vh;
}

.graphiql-container {
  font-family: system-ui, -apple-system, sans-serif;
}

.graphiql-session-header {
  background-color: #f5f5f5;
}

.graphiql-logo {
  color: #1a1a2e;
}

.graphiql-button {
  background-color: #4990e2;
  color: white;
  border: none;
  border-radius: 4px;
}

.graphiql-button:hover {
  background-color: #357abd;
}

.graphiql-toolbar {
  background-color: #fafafa;
}

.graphiql-editor {
  background-color: #ffffff;
}

.graphiql-response {
  background-color: #ffffff;
}
`

// GraphiQLDarkCSS is the dark theme CSS for GraphiQL
const GraphiQLDarkCSS = `
body {
  margin: 0;
  padding: 0;
  height: 100vh;
  overflow: hidden;
  background: #1a1a2e;
}

#graphiql {
  height: 100vh;
}

.graphiql-container {
  font-family: system-ui, -apple-system, sans-serif;
  --color-base: #f8f8f2;
  --color-primary: #8be9fd;
  --color-secondary: #bd93f9;
  --color-tertiary: #50fa7b;
  --color-info: #8be9fd;
  --color-success: #50fa7b;
  --color-warning: #ffb86c;
  --color-error: #ff5555;
  background-color: #282a36;
}

.graphiql-container,
.graphiql-container * {
  color: #f8f8f2;
}

.graphiql-session-header {
  background-color: #1e1f29 !important;
  border-bottom: 1px solid #44475a;
}

.graphiql-logo {
  color: #f8f8f2;
}

.graphiql-button {
  background-color: #6272a4 !important;
  color: #f8f8f2 !important;
  border: none !important;
  border-radius: 4px;
}

.graphiql-button:hover {
  background-color: #7082b4 !important;
}

.graphiql-toolbar {
  background-color: #282a36 !important;
  border-bottom: 1px solid #44475a;
}

.graphiql-sidebar {
  background-color: #1e1f29 !important;
  border-right: 1px solid #44475a;
}

.graphiql-sidebar-section {
  border-bottom: 1px solid #44475a;
}

.graphiql-doc-explorer-title {
  color: #f8f8f2;
}

.graphiql-doc-explorer-content {
  background-color: #282a36;
}

.graphiql-editor {
  background-color: #282a36 !important;
}

.graphiql-editor .CodeMirror {
  background-color: #282a36 !important;
  color: #f8f8f2 !important;
}

.graphiql-editor .CodeMirror-gutters {
  background-color: #1e1f29 !important;
  border-right: 1px solid #44475a;
}

.graphiql-editor .CodeMirror-linenumber {
  color: #6272a4 !important;
}

.graphiql-editor .CodeMirror-cursor {
  border-left: 1px solid #f8f8f2 !important;
}

.graphiql-editor .CodeMirror-selected {
  background: #44475a !important;
}

.graphiql-editor .cm-keyword {
  color: #ff79c6 !important;
}

.graphiql-editor .cm-def {
  color: #50fa7b !important;
}

.graphiql-editor .cm-property {
  color: #8be9fd !important;
}

.graphiql-editor .cm-string {
  color: #f1fa8c !important;
}

.graphiql-editor .cm-number {
  color: #bd93f9 !important;
}

.graphiql-editor .cm-atom {
  color: #bd93f9 !important;
}

.graphiql-editor .cm-punctuation {
  color: #f8f8f2 !important;
}

.graphiql-editor .cm-variable {
  color: #f8f8f2 !important;
}

.graphiql-response {
  background-color: #282a36 !important;
}

.graphiql-response .CodeMirror {
  background-color: #282a36 !important;
  color: #f8f8f2 !important;
}

.graphiql-tabs {
  background-color: #1e1f29 !important;
}

.graphiql-tab {
  background-color: #282a36 !important;
  border: 1px solid #44475a !important;
  color: #f8f8f2 !important;
}

.graphiql-tab-active {
  background-color: #44475a !important;
}

.graphiql-tab:hover {
  background-color: #44475a !important;
}

.graphiql-history-header {
  background-color: #1e1f29 !important;
}

.graphiql-history-item {
  border-bottom: 1px solid #44475a;
}

.graphiql-history-item:hover {
  background-color: #44475a !important;
}

.graphiql-doc-explorer {
  background-color: #282a36 !important;
}

.graphiql-doc-explorer-header {
  background-color: #1e1f29 !important;
  border-bottom: 1px solid #44475a;
}

.graphiql-doc-explorer-back {
  color: #8be9fd !important;
}

.graphiql-doc-explorer-type-name {
  color: #50fa7b !important;
}

.graphiql-doc-explorer-field-name {
  color: #8be9fd !important;
}

.graphiql-doc-explorer-argument-name {
  color: #ffb86c !important;
}

.graphiql-markdown-description {
  color: #c9c9c9 !important;
}

.graphiql-dropdown-menu {
  background-color: #282a36 !important;
  border: 1px solid #44475a !important;
}

.graphiql-dropdown-menu-item {
  color: #f8f8f2 !important;
}

.graphiql-dropdown-menu-item:hover {
  background-color: #44475a !important;
}

input, textarea, select {
  background-color: #44475a !important;
  color: #f8f8f2 !important;
  border: 1px solid #6272a4 !important;
}

input:focus, textarea:focus, select:focus {
  border-color: #8be9fd !important;
  outline: none !important;
}

::placeholder {
  color: #6272a4 !important;
}

/* Scrollbar */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #1e1f29;
}

::-webkit-scrollbar-thumb {
  background: #44475a;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #6272a4;
}
`
