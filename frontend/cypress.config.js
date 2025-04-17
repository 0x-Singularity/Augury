const { defineConfig } = require("Cypress");

module.exports = defineConfig({
  e2e: {
    baseUrl: "http://localhost:3000",
    supportFile: false,
    specPattern: "test/**/*.spec.js",
  },
});
