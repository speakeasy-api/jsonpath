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
    "process.env.VSCODE_TEXTMATE_DEBUG": "false",
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },

});
