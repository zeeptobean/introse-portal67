<script lang="ts">
  import { untrack } from "svelte";
  import { zod4Client } from "sveltekit-superforms/adapters";
  import type { SuperValidated } from "sveltekit-superforms";
  import { portalForm } from "#lib/form.ts";
  import { chargeFormSchema, type ChargeForm as ChargeData } from "#lib/schemas/tuition.ts";
  import type { Term } from "#lib/schemas/shared.ts";
  import type { TuitionSummary } from "#lib/schemas/tuition.ts";
  import { protectDraft } from "#lib/draft.svelte.ts";
  import FormFeedback from "#lib/components/FormFeedback.svelte";
  import TuitionReview from "#lib/components/TuitionReview.svelte";
  let {
    summary,
    term,
    initial,
  }: { summary: TuitionSummary; term: Term; initial: SuperValidated<ChargeData> } = $props();
  let review = $state(false);
  const controller = portalForm(
    untrack(() => initial),
    { validators: zod4Client(chargeFormSchema) },
  );
  const { form, errors, message, allErrors, submitting, enhance } = controller;
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
<form method="POST" use:enhance novalidate>
  <input type="hidden" name="balance_vnd" value={$form.balance_vnd} />
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
    <label for="description">Charge description (required)</label>
    <textarea
      id="description"
      name="description"
      rows="3"
      bind:value={$form.description}
      aria-invalid={!!$errors.description}
      aria-describedby="description-error"></textarea>
    <p class="field-error" id="description-error">{$errors.description?.join(" ") ?? ""}</p>
    <button type="button" class="primary" onclick={prepare} disabled={$message?.uncertain}>
      Review charge
    </button>
  {:else}
    <input type="hidden" name="amount" value={$form.amount} />
    <input type="hidden" name="description" value={$form.description} />
    <TuitionReview
      student={summary.student}
      {term}
      title="Add charge"
      amount={$form.amount}
      before={$form.balance_vnd}
      after={$form.balance_vnd + Number($form.amount)}
      description={$form.description}
    />
    <div class="actions">
      <button class="primary" disabled={$submitting || $message?.uncertain}>
        {$submitting ? "Saving…" : "Confirm and save"}
      </button>
      <button type="button" onclick={() => (review = false)}>Edit values</button>
    </div>
  {/if}
  <noscript><button class="primary">Confirm and save charge</button></noscript>
</form>
