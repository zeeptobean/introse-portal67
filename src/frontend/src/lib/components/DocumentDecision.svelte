<script lang="ts">
  import { untrack } from "svelte";
  import { zod4Client } from "sveltekit-superforms/adapters";
  import type { SuperValidated } from "sveltekit-superforms";
  import { portalForm } from "#lib/form.ts";
  import { decisionInputSchema, type DecisionInput } from "#lib/schemas/documents.ts";
  import type { DocumentRequest } from "#lib/schemas/documents.ts";
  import { protectDraft } from "#lib/draft.svelte.ts";
  import FormFeedback from "#lib/components/FormFeedback.svelte";
  let { record, initial }: { record: DocumentRequest; initial: SuperValidated<DecisionInput> } =
    $props();
  let review = $state(false);
  const controller = portalForm(
    untrack(() => initial),
    {
      validators: zod4Client(decisionInputSchema),
      invalidateAll: "pessimistic",
      onUpdated: ({ form }) => {
        if (form.valid) review = false;
      },
    },
  );
  const { form, errors, message, allErrors, submitting, enhance } = controller;
  protectDraft(() => controller.isTainted());
  async function prepare() {
    const valid = await controller.validateForm();
    if (valid.valid) review = true;
  }
</script>

<form class="panel" method="POST" use:enhance novalidate>
  <h2>Review decision</h2>
  <FormFeedback message={$message} errors={$allErrors} />
  {#if record.status === "pending"}
    {#if !review}
      <fieldset>
        <legend>Decision</legend>
        <label class="radio-label">
          <input type="radio" name="status" bind:group={$form.status} value="approved" />
          Approve
        </label>
        <label class="radio-label">
          <input type="radio" name="status" bind:group={$form.status} value="rejected" />
          Reject
        </label>
      </fieldset>
      <label for="comment">Comment visible to the student (required)</label>
      <textarea
        id="comment"
        name="comment"
        rows="4"
        bind:value={$form.comment}
        aria-invalid={!!$errors.comment}
        aria-describedby="comment-error"></textarea>
      <p class="hint" id="comment-error">
        {$errors.comment?.join(" ") ?? `${[...$form.comment].length} / 1,000 characters`}
      </p>
      <button type="button" class="primary" onclick={prepare} disabled={$message?.uncertain}>
        Review decision
      </button>
    {:else}
      <input type="hidden" name="status" value={$form.status} />
      <input type="hidden" name="comment" value={$form.comment} />
      <div class="review-box">
        <h3>{$form.status === "approved" ? "Approve" : "Reject"} request {record.id}</h3>
        <p>{record.student.full_name} · {record.student.institutional_id ?? "ID not provided"}</p>
        <p class="pre-wrap">{$form.comment}</p>
        <p>This decision is final. The student will see your comment.</p>
      </div>
      <div class="actions">
        <button
          class={$form.status === "rejected" ? "danger" : "primary"}
          disabled={$submitting || $message?.uncertain}
        >
          {$submitting ? "Saving…" : "Confirm decision"}
        </button>
        <button type="button" onclick={() => (review = false)} disabled={$submitting}>
          Edit decision
        </button>
      </div>
    {/if}
  {:else}
    <p>This request has already been reviewed. Its confirmed decision is shown above.</p>
  {/if}
  <noscript><button class="primary">Confirm decision</button></noscript>
</form>
