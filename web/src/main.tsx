import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";
import posthog from "posthog-js";

if (process.env.PUBLIC_POSTHOG_API_KEY) {
  posthog.init(process.env.PUBLIC_POSTHOG_API_KEY, {
    api_host: "https://metrics.speakeasy.com",
    person_profiles: "identified_only",
    disable_session_recording: true,
    autocapture: false,
  });

  posthog.capture("$pageview");
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
