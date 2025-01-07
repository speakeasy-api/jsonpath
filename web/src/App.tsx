import { MoonshineConfigProvider } from "@speakeasy-api/moonshine";
import Playground from "./Playground.tsx";

export default function App() {
  return (
    <MoonshineConfigProvider themeElement={document.body}>
      <Playground />
    </MoonshineConfigProvider>
  );
}
