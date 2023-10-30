import { useRouteError } from "react-router-dom";

function ErrorPage() {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const error = useRouteError() as any;

  return (
    <main>
      <div className="container">
        <h1>An error occurred</h1>
        <p>{error.statusText || error.message}</p>
      </div>
    </main>
  )
}

export default ErrorPage
