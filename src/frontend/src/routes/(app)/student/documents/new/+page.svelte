<script lang="ts">
  import { untrack } from "svelte";
  import { zod4Client } from "sveltekit-superforms/adapters";
  import { portalForm } from "#lib/form.ts";
  import { documentInputSchema, requestTypes } from "#lib/schemas/documents.ts";
  import { protectDraft } from "#lib/draft.svelte.ts";
  import FormFeedback from "#lib/components/FormFeedback.svelte";
  import PageHeader from "#lib/components/PageHeader.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
  const controller = portalForm(
    untrack(() => data.form),
    { validators: zod4Client(documentInputSchema) },
  );
  const { form, errors, message, allErrors, submitting, enhance } = controller;
  protectDraft(() => controller.isTainted());
</script>

<a class="breadcrumb" href="/student/documents">← Documents</a>
<PageHeader title="New document request" subtitle="Tell us what you need and why." />
<div class="form-width">
  <FormFeedback message={$message} errors={$allErrors} />
  {#if $message?.uncertain}
    <a href="/student/documents">Inspect your latest requests</a>
  {/if}
  <form method="POST" use:enhance novalidate>
    <section class="identity-strip">
      <strong>{data.user.full_name}</strong>
      <span>{data.user.institutional_id ?? "ID not provided"} · {data.user.email}</span>
    </section>
    <label for="request_type">Document type</label>
    <select id="request_type" name="request_type" bind:value={$form.request_type}>
      {#each Object.entries(requestTypes) as [value, label]}
        <option {value}>{label}</option>
      {/each}
    </select>
    <p class="hint">
      {$form.request_type === "transcript"
        ? "A request for academic results documentation."
        : $form.request_type === "enrollment_confirmation"
          ? "A request for confirmation of your enrollment."
          : "Describe the document you need in the purpose."}
    </p>
    <label for="purpose">Purpose (required)</label>
    <textarea
      id="purpose"
      name="purpose"
      rows="5"
      bind:value={$form.purpose}
      aria-invalid={!!$errors.purpose}
      aria-describedby="purpose-hint"></textarea>
    <p class="hint" id="purpose-hint">
      {$errors.purpose?.join(" ") ??
        `${[...$form.purpose].length} / 1,000 characters · Include enough context for the reviewer.`}
    </p>
    <section class="review-box">
      <h2>Review your request</h2>
      <p>{requestTypes[$form.request_type]} for {data.user.full_name}</p>
      <p class="pre-wrap">{$form.purpose || "Your purpose will appear here."}</p>
      <p class="small muted">
        This submits a sample request. Approval does not issue a legally valid document.
      </p>
    </section>
    <div class="actions">
      <button class="primary" disabled={$submitting || $message?.uncertain}>
        {$submitting ? "Submitting…" : "Submit request"}
      </button>
      <a href="/student/documents">Cancel</a>
    </div>
  </form>
</div>
