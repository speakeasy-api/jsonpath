import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import svgr from "vite-plugin-svgr";
import path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), svgr()],
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
