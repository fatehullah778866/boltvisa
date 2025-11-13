export default function GuideTroubleshooting() {
  return (
    <div className="prose max-w-none">
      <h1>Troubleshooting</h1>
      <h3>“Invalid user data” or redirects to login</h3>
      <ul>
        <li>Sign up or log in again; token may be missing or expired.</li>
        <li>Ensure the API is running at <code>NEXT_PUBLIC_API_URL</code>.</li>
      </ul>
      <h3>API 404s</h3>
      <ul>
        <li>Restart the Go API after code changes.</li>
        <li>Verify routes (e.g., <code>POST /api/v1/applications</code>).</li>
      </ul>
      <h3>Build issues</h3>
      <ul>
        <li>Run <code>pnpm -w -r build</code> and check TypeScript errors.</li>
      </ul>
    </div>
  );
}
