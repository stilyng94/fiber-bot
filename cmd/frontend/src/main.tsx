import React from 'react'
import ReactDOM from 'react-dom/client'
import '@/index.css'
import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import CheckoutPage, { loader as checkoutPageLoader } from '@/pages/CheckoutPage'
import ErrorPage from '@/pages/ErrorPage'
import Homepage, { loader as homepageLoader } from '@/pages/HomePage'
import App from '@/App'
import { ThemeProvider } from '@/components/custom/theme-provider'

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <ErrorPage />,
    children: [
      {
        element: <Homepage />,
        index: true,
        loader: homepageLoader
      },
      {
        path: "checkout",
        element: <CheckoutPage />,
        loader: checkoutPageLoader
      }
    ]
  }
])



ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
      <RouterProvider router={router} />
    </ThemeProvider>
  </React.StrictMode>,
)
