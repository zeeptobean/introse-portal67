<script lang="ts">
  import PageHeader from "#lib/components/PageHeader.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<PageHeader title="Account & help" subtitle="Your identity and the scope of this prototype." />
<div class="form-width">
  <dl class="account-details">
    <div>
      <dt>Name</dt>
      <dd>{data.user.full_name}</dd>
    </div>
    <div>
      <dt>Institutional ID</dt>
      <dd>{data.user.institutional_id ?? "Not provided"}</dd>
    </div>
    <div>
      <dt>Email</dt>
      <dd>{data.user.email}</dd>
    </div>
    <div>
      <dt>Role</dt>
      <dd class="role-label">{data.user.role}</dd>
    </div>
  </dl>
  <h2>Working with Portal67</h2>
  <details>
    <summary>Registration and schedule conflicts</summary>
    <p>
      A preview does not reserve a seat. The API checks seats, duplicate registrations, registration
      windows, and overlapping meetings. Cancelling releases your seat; registering again depends on
      availability.
    </p>
  </details>
  <details>
    <summary>Results and estimated GPA</summary>
    <p>
      Only graded credits contribute to GPA. A zero score is a grade; an untouched blank score is
      not. The Results page explains the current sample scale.
    </p>
  </details>
  <details>
    <summary>Simulated tuition</summary>
    <p>
      Charges and recorded payments are sample records. Portal67 does not transfer money, process
      bank payments, or connect to official university systems. Voided payments remain visible but
      do not reduce the balance.
    </p>
  </details>
  <details>
    <summary>Document decisions</summary>
    <p>
      Requests are Pending, Approved, or Rejected. A staff comment explains the decision. Approval
      does not issue a legally valid document.
    </p>
  </details>
  <details>
    <summary>Accounts and access</summary>
    <p>
      Accounts are provided by the project team. Password recovery and account administration are
      not available in this prototype.
    </p>
  </details>
</div>

<style>
  .account-details {
    margin: 0 0 40px;
  }
  .account-details > div {
    display: grid;
    grid-template-columns: 160px minmax(0, 1fr);
    gap: 16px;
    padding: 16px 0;
    border-bottom: 1px solid var(--line);
  }
  @media (max-width: 620px) {
    .account-details > div {
      grid-template-columns: 1fr;
      gap: 8px;
    }
  }
</style>
