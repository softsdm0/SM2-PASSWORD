import { defineConfig } from "wxt";

// See https://wxt.dev/api/config.html
export default defineConfig({
  modules: ["@wxt-dev/module-vue"],
  srcDir: "src",
  manifest: {
    action: {},
    permissions: [
      "storage",
      "activeTab",
      "webNavigation",
      "tabs",
      "cookies",
      "<all_urls>"
    ],
  },
});
