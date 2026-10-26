<script lang="ts">
  import type { TuitionSummary } from "#lib/schemas/tuition.ts";
  import { money, timestamp } from "#lib/domain.ts";
  import Status from "#lib/components/Status.svelte";
  let {
    summary,
    staff = false,
    termId,
  }: { summary: TuitionSummary; staff?: boolean; termId: string } = $props();
  const hasRecords = $derived(summary.terms.some((t) => t.charges.length + t.payments.length > 0));
</script>

<section class="identity-strip">
  <strong>{summary.student.full_name}</strong>
  <span>{summary.student.institutional_id ?? "ID not provided"} · {summary.student.email}</span>
</section>
<div class="ledger-summary">
  <section>
    <p class="eyebrow">Semester charges</p>
    <strong>{money(summary.charged_vnd)}</strong>
  </section>
  <span class="math-sign" aria-hidden="true">−</span>
  <section>
    <p class="eyebrow">Recorded payments</p>
    <strong>{money(summary.paid_vnd)}</strong>
  </section>
  <span class="math-sign" aria-hidden="true">=</span>
  <section>
    <p class="eyebrow">{summary.balance_vnd < 0 ? "Credit balance" : "Balance to review"}</p>
    <strong>{money(Math.abs(summary.balance_vnd))}</strong>
    <Status value={summary.status} />
  </section>
</div>
<p class="muted small">
  Simulated tuition records. No money is transferred through Portal67.{summary.balance_vnd < 0
    ? " Recorded payments exceed charges; the difference is a credit in your favor."
    : ""}
</p>
{#if staff}
  <div class="actions">
    <a
      class="button primary"
      href="/staff/tuition/students/{summary.student.id}/payments/new?term={termId}"
    >
      Record payment
    </a>
    <a class="button" href="/staff/tuition/students/{summary.student.id}/charges/new?term={termId}">
      Add charge
    </a>
  </div>
{/if}
{#if !hasRecords}
  <div class="empty">No tuition records for this semester.</div>
{/if}
{#each summary.terms as term}
  <section class="ledger-term">
    <h2>{term.term.name}</h2>
    <h3>Charges</h3>
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Description</th>
            <th>Recorded</th>
            <th class="numeric">Amount</th>
          </tr>
        </thead>
        <tbody>
          {#each term.charges as charge}
            <tr>
              <td data-label="Description">
                {charge.description}
                <span class="muted small">Reference {charge.id}</span>
              </td>
              <td data-label="Recorded">{timestamp(charge.created_at)}</td>
              <td class="numeric" data-label="Amount">{money(charge.amount_vnd)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
    {#if !term.charges.length}
      <p class="empty">No charges recorded.</p>
    {/if}
    <h3>Payments</h3>
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Reference / note</th>
            <th>Payment date</th>
            <th>Status</th>
            <th class="numeric">Amount</th>
            {#if staff}
              <th>Action</th>
            {/if}
          </tr>
        </thead>
        <tbody>
          {#each term.payments as payment}
            <tr class:voided={payment.status === "voided"}>
              <td data-label="Reference / note">
                {payment.id}
                <span class="small pre-wrap">{payment.note ?? "No note"}</span>
              </td>
              <td data-label="Payment date">{timestamp(payment.paid_at)}</td>
              <td data-label="Status"><Status value={payment.status} /></td>
              <td data-label="Amount" class="numeric">
                {money(payment.amount_vnd)}
                {#if payment.status === "voided"}
                  <span class="small muted">Excluded from totals</span>
                {/if}
              </td>
              {#if staff}
                <td>
                  <a
                    href="/staff/tuition/students/{summary.student
                      .id}/payments/{payment.id}?term={termId}"
                  >
                    {payment.status === "voided" ? "View history" : "Correct / void"}
                  </a>
                </td>
              {/if}
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
    {#if !term.payments.length}
      <p class="empty">No payments recorded.</p>
    {/if}
  </section>
{/each}

<style>
  .ledger-summary {
    padding: 28px 24px;
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr) auto minmax(0, 1fr);
    gap: 24px;
    background: white;
    border: 1px solid var(--line);
    border-radius: 8px;
    align-items: start;
    margin-bottom: 16px;
  }
  .ledger-summary strong {
    font-size: 22px;
    font-variant-numeric: tabular-nums;
    display: block;
    overflow-wrap: anywhere;
  }
  .ledger-summary section :global(.status) {
    margin-top: 10px;
  }
  .math-sign {
    font-size: 24px;
    color: var(--muted);
    padding-top: 20px;
  }
  @media (max-width: 1200px) {
    .ledger-summary {
      gap: 12px;
      padding: 24px 16px;
    }
    .ledger-summary strong {
      font-size: 19px;
    }
    .ledger-summary .eyebrow {
      font-size: 10px;
      letter-spacing: 0.6px;
    }
  }
  @media (max-width: 620px) {
    .ledger-summary {
      grid-template-columns: 1fr;
      gap: 16px;
      padding: 24px;
    }
    .ledger-summary section {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 12px;
      flex-wrap: wrap;
    }
    .ledger-summary .eyebrow {
      font-size: 11px;
      max-width: 100%;
      margin: 4px 0;
    }
    .ledger-summary strong {
      font-size: 20px;
      margin-left: auto;
    }
    .ledger-summary section :global(.status) {
      margin-top: 0;
      width: 100%;
      background: none;
      padding: 0;
      justify-content: flex-end;
    }
    .math-sign {
      display: none;
    }
  }
  .ledger-term {
    margin-top: 40px;
  }
  .ledger-term > h3 {
    margin-top: 24px;
  }
  .voided {
    color: var(--muted);
  }
</style>
