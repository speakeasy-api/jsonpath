import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import svgr from "vite-plugin-svgr";
import path from "path";
import vercel from "vite-plugin-vercel";
import { createHtmlPlugin } from "vite-plugin-html";

const htmlPlugin = createHtmlPlugin({
  minify: true,
  entry: "src/main.tsx",
  inject:
    process.env.ANALYTICS_SCRIPT_HEAD && process.env.ANALYTICS_SCRIPT_BODY
      ? {
          data: {
            injectScriptHead: process.env.ANALYTICS_SCRIPT_HEAD,
            injectScriptBody: process.env.ANALYTICS_SCRIPT_BODY,
          },
        }
      : {
          data: {
            injectScriptHead: ``,
            injectScriptBody: ``,
          },
        },
});

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), svgr(), vercel(), htmlPlugin],
  define: {
    // add nodejs shims: moonshine requires them.
    global: {},
    // importing moonshine error'd without this.
    "process.env.VSCODE_TEXTMATE_DEBUG": "false",
    "process.env.PUBLIC_POSTHOG_API_KEY": JSON.stringify(
      process.env.PUBLIC_POSTHOG_API_KEY,
    ),
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
