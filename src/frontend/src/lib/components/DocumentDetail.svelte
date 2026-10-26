<script lang="ts">
  import type { SuperValidated } from "sveltekit-superforms";
  import { requestTypes, type DecisionInput } from "#lib/schemas/documents.ts";
  import type { DocumentRequest } from "#lib/schemas/documents.ts";
  import { timestamp } from "#lib/domain.ts";
  import Status from "#lib/components/Status.svelte";
  import DocumentDecision from "#lib/components/DocumentDecision.svelte";
  let {
    record,
    staff = false,
    decisionForm,
  }: {
    record: DocumentRequest;
    staff?: boolean;
    decisionForm?: SuperValidated<DecisionInput>;
  } = $props();
</script>

<div class="detail-layout">
  <section>
    <div class="panel">
      <div class="section-heading">
        <h2>{requestTypes[record.request_type]}</h2>
        <Status value={record.status} />
      </div>
      <p class="muted small">Reference {record.id} · Submitted {timestamp(record.submitted_at)}</p>
      <section class="identity-strip">
        <strong>{record.student.full_name}</strong>
        <span>{record.student.institutional_id ?? "ID not provided"} · {record.student.email}</span>
      </section>
      <h3>Purpose</h3>
      <p class="pre-wrap">{record.purpose}</p>
      {#if record.decision_comment}
        <div class="decision">
          <h3>{record.status === "approved" ? "Approved" : "Rejected"} · Decision comment</h3>
          <p class="pre-wrap">{record.decision_comment}</p>
          {#if record.decided_at}
            <p class="small muted">{timestamp(record.decided_at)}</p>
          {/if}
        </div>
      {/if}
      <p class="small muted">
        A prototype request record. Approval does not issue a document or create a legally valid
        certificate.
      </p>
    </div>
    {#if staff && decisionForm}
      <DocumentDecision {record} initial={decisionForm} />
    {/if}
  </section>
  <aside class="panel">
    <h2>Request history</h2>
    <ol class="timeline">
      {#each record.history as event}
        <li>
          <Status value={event.to_status} />
          <p>{event.actor.full_name} · {event.actor.role}</p>
          <p class="muted small">{timestamp(event.created_at)}</p>
          {#if event.comment}
            <p class="pre-wrap">{event.comment}</p>
          {/if}
        </li>
      {/each}
    </ol>
  </aside>
</div>

<style>
  .detail-layout {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 320px;
    gap: 32px;
    align-items: start;
  }
  @media (max-width: 1200px) {
    .detail-layout {
      grid-template-columns: minmax(0, 1fr) 280px;
      gap: 24px;
    }
  }
  @media (max-width: 900px) {
    .detail-layout {
      grid-template-columns: 1fr;
    }
  }
  .decision {
    border-top: 1px solid var(--line);
    margin-top: 24px;
    padding-top: 24px;
  }
  .timeline {
    padding: 0;
    list-style: none;
  }
  .timeline li {
    border-left: 2px solid var(--line);
    padding: 0 0 24px 16px;
    margin-top: 20px;
  }
  .timeline p {
    font-size: 13px;
    margin: 8px 0;
  }
</style>
