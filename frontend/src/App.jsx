import './App.scss'
import Layout from "./layout"
import { ExplorerProvider } from "./context/explorer_context"
import { UiProvider } from "./context/ui_context"

function App() {
  return (
    <UiProvider>
      <ExplorerProvider>
        <Layout />
      </ExplorerProvider>
    </UiProvider>
  )
}

export default App
