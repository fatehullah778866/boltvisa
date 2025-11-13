export default function GuideDashboard() {
  return (
    <div className="prose max-w-none">
      <h1>Dashboard</h1>
      <p>The dashboard shows your profile, quick actions, live stats, recent applications, and a category strip.</p>
      <ul>
        <li><strong>Quick Actions:</strong> New Application, Browse Categories, View Applications, Notifications.</li>
        <li><strong>Stats:</strong> Totals by status (draft, submitted, in review, approved, rejected).</li>
        <li><strong>Recent:</strong> Last few applications you touched.</li>
      </ul>
    </div>
  );
}
