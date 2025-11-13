export default function GuideApplications() {
  return (
    <div className="prose max-w-none">
      <h1>Applications</h1>
      <h3>Create</h3>
      <p>Use <strong>Applications → New</strong> or “Apply” from Visa Categories. Fill details and decide whether to submit now or save as draft.</p>
      <h3>Track</h3>
      <p>Visit <strong>/applications</strong> to see status and timestamps. Click an item (coming soon: detail view) for more info.</p>
      <h3>Statuses</h3>
      <ul>
        <li>Draft → Submitted → In Review → Approved/Rejected (Cancelled if withdrawn)</li>
      </ul>
    </div>
  );
}
