import { Outlet } from "react-router-dom";
import { useEffect } from "react";
import { useTheme } from "@/hooks/use-theme";


const telegram = window.Telegram.WebApp;

function App() {
  const { setTheme } = useTheme()

  useEffect(() => {
    telegram.ready()
    setTheme(telegram.colorScheme ?? 'system')
    return () => {
      telegram.close();
    }
  }, [setTheme])

  return (
    <main>
      <div className="container">
        <Outlet />
      </div>
    </main>
  )

}

export default App

