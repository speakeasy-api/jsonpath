import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import posthog from "posthog-js";

if (process.env.PUBLIC_POSTHOG_API_KEY) {
  posthog.init(process.env.PUBLIC_POSTHOG_API_KEY, {
    api_host: "https://us.i.posthog.com",
    person_profiles: "identified_only",
  });
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
