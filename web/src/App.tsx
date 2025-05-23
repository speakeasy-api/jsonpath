import { MoonshineConfigProvider } from "@speakeasy-api/moonshine";
import Playground from "./Playground";
import "@speakeasy-api/moonshine/moonshine.css";

export default function App() {
  return (
    <MoonshineConfigProvider themeElement={document.body}>
      <Playground />
    </MoonshineConfigProvider>
  );
}
