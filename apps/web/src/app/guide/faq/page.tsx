export default function GuideFAQ() {
  return (
    <div className="prose max-w-none">
      <h1>FAQ</h1>
      <h3>Can I edit an application after submitting?</h3>
      <p>Edits will be supported soon via an Application Detail page with update actions.</p>
      <h3>How are notifications delivered?</h3>
      <p>In-app first; email/webhook integrations are planned.</p>
      <h3>Where are documents stored?</h3>
      <p>In dev, no storage; in prod, Google Cloud Storage.</p>
    </div>
  );
}
