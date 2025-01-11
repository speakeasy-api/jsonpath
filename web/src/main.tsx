import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";
import posthog from "posthog-js";

if (process.env.PUBLIC_POSTHOG_API_KEY) {
  posthog.init(process.env.PUBLIC_POSTHOG_API_KEY, {
    api_host: "https://metrics.speakeasy.com",
    disable_session_recording: true,
    autocapture: false,
    before_send: function (event) {
      if (event?.properties?.$current_url) {
        const url = new URL(event.properties.$current_url);
        url.hash = "";
        event.properties.$current_url = url.toString();
      }
      return event;
    },
  });
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
