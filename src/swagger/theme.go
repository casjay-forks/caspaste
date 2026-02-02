// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

// Swagger UI theming per AI.md PART 14
// Provides light and dark theme CSS for Swagger UI
package swagger

// SwaggerLightCSS is the light theme CSS for Swagger UI
const SwaggerLightCSS = `
body {
  margin: 0;
  padding: 0;
  background: #fafafa;
}

.swagger-ui .topbar {
  background-color: #1a1a2e;
}

.swagger-ui .info .title {
  color: #3b4151;
}

.swagger-ui .info .description {
  color: #3b4151;
}

.swagger-ui .opblock-tag {
  color: #3b4151;
  border-bottom: 1px solid rgba(59,65,81,.3);
}

.swagger-ui .opblock.opblock-get {
  border-color: #61affe;
  background: rgba(97,175,254,.1);
}

.swagger-ui .opblock.opblock-post {
  border-color: #49cc90;
  background: rgba(73,204,144,.1);
}

.swagger-ui .opblock.opblock-put {
  border-color: #fca130;
  background: rgba(252,161,48,.1);
}

.swagger-ui .opblock.opblock-delete {
  border-color: #f93e3e;
  background: rgba(249,62,62,.1);
}

.swagger-ui .btn {
  border-radius: 4px;
}

.swagger-ui .btn.execute {
  background-color: #4990e2;
  border-color: #4990e2;
}

.swagger-ui .btn.execute:hover {
  background-color: #357abd;
}
`

// SwaggerDarkCSS is the dark theme CSS for Swagger UI
const SwaggerDarkCSS = `
body {
  margin: 0;
  padding: 0;
  background: #1a1a2e;
}

.swagger-ui {
  background: #1a1a2e;
}

.swagger-ui .topbar {
  background-color: #0f0f1a;
}

.swagger-ui .info .title {
  color: #f8f8f2;
}

.swagger-ui .info .description {
  color: #c9c9c9;
}

.swagger-ui .info .description p {
  color: #c9c9c9;
}

.swagger-ui .opblock-tag {
  color: #f8f8f2;
  border-bottom: 1px solid rgba(255,255,255,.2);
}

.swagger-ui .opblock-tag:hover {
  background: rgba(255,255,255,.05);
}

.swagger-ui .opblock {
  background: rgba(255,255,255,.05);
  border-radius: 4px;
}

.swagger-ui .opblock.opblock-get {
  border-color: #61affe;
  background: rgba(97,175,254,.15);
}

.swagger-ui .opblock.opblock-post {
  border-color: #49cc90;
  background: rgba(73,204,144,.15);
}

.swagger-ui .opblock.opblock-put {
  border-color: #fca130;
  background: rgba(252,161,48,.15);
}

.swagger-ui .opblock.opblock-delete {
  border-color: #f93e3e;
  background: rgba(249,62,62,.15);
}

.swagger-ui .opblock .opblock-summary-method {
  border-radius: 3px;
}

.swagger-ui .opblock .opblock-summary-description {
  color: #c9c9c9;
}

.swagger-ui .opblock .opblock-section-header {
  background: rgba(255,255,255,.05);
}

.swagger-ui .opblock .opblock-section-header h4 {
  color: #f8f8f2;
}

.swagger-ui table thead tr th {
  color: #f8f8f2;
  border-bottom: 1px solid rgba(255,255,255,.2);
}

.swagger-ui table tbody tr td {
  color: #c9c9c9;
}

.swagger-ui .parameter__name {
  color: #f8f8f2;
}

.swagger-ui .parameter__type {
  color: #8be9fd;
}

.swagger-ui .parameter__in {
  color: #bd93f9;
}

.swagger-ui .response-col_status {
  color: #50fa7b;
}

.swagger-ui .response-col_description {
  color: #c9c9c9;
}

.swagger-ui .model-title {
  color: #f8f8f2;
}

.swagger-ui .model {
  color: #c9c9c9;
}

.swagger-ui .model .property {
  color: #f8f8f2;
}

.swagger-ui .model .property.primitive {
  color: #8be9fd;
}

.swagger-ui section.models {
  border: 1px solid rgba(255,255,255,.2);
}

.swagger-ui section.models h4 {
  color: #f8f8f2;
}

.swagger-ui .model-box {
  background: rgba(255,255,255,.05);
}

.swagger-ui .btn {
  border-radius: 4px;
  color: #f8f8f2;
}

.swagger-ui .btn.execute {
  background-color: #6272a4;
  border-color: #6272a4;
  color: #f8f8f2;
}

.swagger-ui .btn.execute:hover {
  background-color: #7082b4;
}

.swagger-ui .btn.cancel {
  background-color: #44475a;
  border-color: #44475a;
  color: #f8f8f2;
}

.swagger-ui select {
  background: #44475a;
  color: #f8f8f2;
  border: 1px solid #6272a4;
}

.swagger-ui input[type="text"],
.swagger-ui textarea {
  background: #44475a;
  color: #f8f8f2;
  border: 1px solid #6272a4;
}

.swagger-ui .highlight-code {
  background: #282a36;
}

.swagger-ui .highlight-code pre {
  color: #f8f8f2;
}

.swagger-ui .response-control-media-type__accept-message {
  color: #50fa7b;
}

.swagger-ui .loading-container .loading {
  background-color: #44475a;
}

.swagger-ui .scheme-container {
  background: #282a36;
}

.swagger-ui .servers-title {
  color: #f8f8f2;
}

.swagger-ui .servers>label {
  color: #c9c9c9;
}

.swagger-ui .auth-wrapper {
  background: rgba(255,255,255,.05);
}

.swagger-ui .authorization__btn {
  background: transparent;
  border-color: #6272a4;
}

.swagger-ui .authorization__btn svg {
  fill: #f8f8f2;
}

.swagger-ui .dialog-ux .modal-ux {
  background: #282a36;
  border: 1px solid #44475a;
}

.swagger-ui .dialog-ux .modal-ux-header h3 {
  color: #f8f8f2;
}

.swagger-ui .dialog-ux .modal-ux-content {
  color: #c9c9c9;
}

.swagger-ui .markdown code {
  background: #44475a;
  color: #f8f8f2;
}

.swagger-ui .markdown pre {
  background: #282a36;
  color: #f8f8f2;
}

.swagger-ui .renderedMarkdown p {
  color: #c9c9c9;
}

.swagger-ui .errors-wrapper {
  background: rgba(249,62,62,.2);
}

.swagger-ui .errors h4 {
  color: #ff5555;
}
`
