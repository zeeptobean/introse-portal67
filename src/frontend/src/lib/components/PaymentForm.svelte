<script lang="ts">
  import { untrack } from "svelte";
  import { zod4Client } from "sveltekit-superforms/adapters";
  import type { SuperValidated } from "sveltekit-superforms";
  import { portalForm } from "#lib/form.ts";
  import { paymentFormSchema, type PaymentForm as PaymentData } from "#lib/schemas/tuition.ts";
  import type { Payment, TuitionSummary } from "#lib/schemas/tuition.ts";
  import type { Term } from "#lib/schemas/shared.ts";
  import { protectDraft } from "#lib/draft.svelte.ts";
  import { timestamp, money } from "#lib/domain.ts";
  import FormFeedback from "#lib/components/FormFeedback.svelte";
  import TuitionReview from "#lib/components/TuitionReview.svelte";
  import Notice from "#lib/components/Notice.svelte";
  let {
    summary,
    term,
    payment,
    initial,
  }: {
    summary: TuitionSummary;
    term: Term;
    payment: Payment | null;
    initial: SuperValidated<PaymentData>;
  } = $props();
  let review = $state(false);
  const controller = portalForm<PaymentData>(
    untrack(() => initial),
    {
      validators: zod4Client(paymentFormSchema),
    },
  );
  const { form, errors, message, allErrors, submitting, enhance } = controller;
  const locked = $derived(payment?.status === "voided");
  const expected = $derived(
    $form.balance_vnd +
      (payment?.amount_vnd ?? 0) -
      ($form.voiding ? 0 : Number($form.amount || 0)),
  );
  protectDraft(() => controller.isTainted());
  async function prepare() {
    const valid = await controller.validateForm();
    if (valid.valid) review = true;
  }
</script>

<FormFeedback message={$message} errors={$allErrors} />
{#if $message}
  <a href="/staff/tuition/students/{summary.student.id}?term={term.id}">
    Inspect the latest ledger
  </a>
{/if}
{#if locked}
  <Notice>
    This payment is voided and cannot be changed. Its history remains available below.
  </Notice>
{:else}
  <form method="POST" use:enhance novalidate>
    <input type="hidden" name="balance_vnd" value={$form.balance_vnd} />
    <input type="hidden" name="revision" value={$form.revision} />
    {#if !review}
      <label for="amount">Amount (VND)</label>
      <input
        id="amount"
        name="amount"
        inputmode="numeric"
        bind:value={$form.amount}
        aria-invalid={!!$errors.amount}
        aria-describedby="amount-error"
      />
      <p class="field-error" id="amount-error">{$errors.amount?.join(" ") ?? ""}</p>
      <label for="paid_at">Payment date and time (campus time UTC+7)</label>
      <input
        id="paid_at"
        name="paid_at"
        type="datetime-local"
        bind:value={$form.paid_at}
        aria-invalid={!!$errors.paid_at}
        aria-describedby="paid_at-error"
      />
      <p class="field-error" id="paid_at-error">{$errors.paid_at?.join(" ") ?? ""}</p>
      <label for="note">Note (optional, 500 characters maximum)</label>
      <textarea
        id="note"
        name="note"
        rows="3"
        bind:value={$form.note}
        aria-invalid={!!$errors.note}
        aria-describedby="note-error"></textarea>
      <p class="field-error" id="note-error">{$errors.note?.join(" ") ?? ""}</p>
      {#if payment}
        <label class="radio-label">
          <input type="checkbox" name="voiding" bind:checked={$form.voiding} />
          Void this payment
        </label>
        <p class="hint">
          Voiding retains the record and removes it from payment totals. It cannot be undone here.
        </p>
      {/if}
      <button class="primary" type="button" onclick={prepare} disabled={$message?.uncertain}>
        Review {$form.voiding ? "void" : "payment"}
      </button>
    {:else}
      <input type="hidden" name="amount" value={$form.amount} />
      <input type="hidden" name="paid_at" value={$form.paid_at} />
      <input type="hidden" name="note" value={$form.note} />
      <input type="hidden" name="voiding" value={$form.voiding ? "true" : ""} />
      <TuitionReview
        student={summary.student}
        {term}
        title={$form.voiding ? "Void payment" : payment ? "Correct payment" : "Record payment"}
        amount={$form.amount}
        before={$form.balance_vnd}
        after={expected}
        description={`${timestamp(new Date(`${$form.paid_at}+07:00`).toISOString())}\n${$form.note}`}
      />
      {#if $form.voiding}
        <p>This removes the payment from totals and retains its read-only history.</p>
      {/if}
      <div class="actions">
        <button
          class={$form.voiding ? "danger" : "primary"}
          disabled={$submitting || $message?.uncertain}
        >
          {$submitting ? "Saving…" : "Confirm and save"}
        </button>
        <button type="button" onclick={() => (review = false)} disabled={$submitting}>
          Edit values
        </button>
      </div>
    {/if}
    <noscript><button class="primary">Confirm and save payment</button></noscript>
  </form>
{/if}
{#if payment}
  <details class="panel" open>
    <summary>Payment history · Reference {payment.id}</summary>
    {#each payment.history as event}
      <article class="history-event">
        <strong>{event.action} · {money(event.amount_vnd)}</strong>
        <p>{event.actor.full_name} · {timestamp(event.created_at)}</p>
        <p>{event.note ?? "No note"}</p>
      </article>
    {/each}
  </details>
{/if}

<style>
  .history-event {
    padding: 16px 0;
    border-bottom: 1px solid var(--line);
    font-size: 13px;
  }
  .history-event p {
    margin: 8px 0;
  }
</style>
