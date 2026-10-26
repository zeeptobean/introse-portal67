<script lang="ts">
  import { untrack, onMount, tick } from "svelte";
  import { page } from "$app/state";
  import { zod4Client } from "sveltekit-superforms/adapters";
  import type { SuperValidated } from "sveltekit-superforms";
  import { portalForm, actionUrl } from "#lib/form.ts";
  import { gradeBatchSchema, type GradeForm } from "#lib/schemas/academics.ts";
  import type { RosterEntry, Section } from "#lib/schemas/academics.ts";
  import { protectDraft } from "#lib/draft.svelte.ts";
  import FormFeedback from "#lib/components/FormFeedback.svelte";
  import Notice from "#lib/components/Notice.svelte";
  let {
    section,
    initial,
    initialForm,
  }: { section: Section; initial: RosterEntry[]; initialForm: SuperValidated<GradeForm> } =
    $props();
  let query = $state("");
  let ready = $state(false);
  onMount(() => {
    ready = true;
  });
  const controller = portalForm(
    untrack(() => initialForm),
    {
      validators: zod4Client(gradeBatchSchema),
      dataType: "json",
      invalidateAll: "pessimistic",
    },
  );
  const { form, errors, message, allErrors, submitting, enhance } = controller;
  const changed = $derived(
    initial.filter((row) => $form.grades.some((g) => g.enrollment_id === row.enrollment_id)),
  );
  const visible = $derived(
    initial.filter((row) =>
      `${row.full_name} ${row.institutional_id ?? ""}`.toLowerCase().includes(query.toLowerCase()),
    ),
  );
  const hidden = $derived(
    changed.filter((row) => !visible.some((v) => v.enrollment_id === row.enrollment_id)).length,
  );
  const closed = $derived(Date.now() > Date.parse(section.term.grading_closes_at));
  protectDraft(() => changed.length > 0);
  function edit(row: RosterEntry, value: string) {
    form.update((data) => ({
      grades: [
        ...data.grades.filter((g) => g.enrollment_id !== row.enrollment_id),
        ...(value === (row.score?.toString() ?? "")
          ? []
          : [{ enrollment_id: row.enrollment_id, score: value }]),
      ],
    }));
  }
  function rowError(id: string) {
    const index = $form.grades.findIndex((g) => g.enrollment_id === id);
    return index >= 0 ? $errors.grades?.[index]?.score : undefined;
  }
  async function focusError(path: string) {
    query = "";
    await tick();
    const index = Number(path.match(/^grades\[(\d+)\]/)?.[1]);
    const id = Number.isInteger(index) ? `grade-${$form.grades[index]?.enrollment_id}` : "grades";
    document.getElementById(id)?.focus();
  }
</script>

<p class="small muted">
  Save applies to every changed row, including filtered rows. Blank untouched scores remain
  ungraded. Clearing an existing score is unsupported.
</p>
<Notice>
  Changes become visible when saved. Concurrent edit detection is not available; coordinate with
  other editors.
</Notice>
{#if closed}
  <Notice tone="caution">Grading has closed for this semester. These records are read-only.</Notice>
{/if}
<FormFeedback message={$message} errors={$allErrors} onSelectError={focusError} />
<form method="POST" action={actionUrl(page.url, "save")} use:enhance novalidate>
  <div class="toolbar">
    <label class="grow">
      Find student
      <input type="search" bind:value={query} placeholder="Name or institutional ID" />
    </label>
    <button type="button" onclick={() => (query = "")}>Clear search</button>
  </div>
  <div class="table-wrap">
    <table class="grade-table" id="grades">
      <thead>
        <tr>
          <th>Student</th>
          <th>Institutional ID</th>
          <th>Saved score</th>
          <th>Score / 10</th>
          <th>State</th>
        </tr>
      </thead>
      <tbody>
        {#each visible as row}
          <tr>
            <td data-label="Student">
              <strong>{row.full_name}</strong>
              <span class="muted small">{row.email}</span>
            </td>
            <td data-label="Institutional ID">{row.institutional_id ?? "Not provided"}</td>
            <td data-label="Saved score">{row.score?.toFixed(2) ?? "Not graded"}</td>
            <td data-label="Score / 10">
              <input
                class="score-input"
                id="grade-{row.enrollment_id}"
                aria-label="Score for {row.full_name}"
                type="number"
                min="0"
                max="10"
                step="0.01"
                value={$form.grades.find((g) => g.enrollment_id === row.enrollment_id)?.score ??
                  row.score ??
                  ""}
                oninput={(event) => edit(row, event.currentTarget.value)}
                aria-invalid={!!rowError(row.enrollment_id)}
                aria-describedby={rowError(row.enrollment_id)
                  ? `error-${row.enrollment_id}`
                  : undefined}
                disabled={!ready || $submitting || closed || $message?.uncertain}
              />
              {#if rowError(row.enrollment_id)}
                <p class="field-error" id="error-{row.enrollment_id}">
                  {rowError(row.enrollment_id)?.join(" ")}
                </p>
              {/if}
            </td>
            <td data-label="State">
              {rowError(row.enrollment_id)
                ? "Invalid"
                : changed.some((c) => c.enrollment_id === row.enrollment_id)
                  ? "Edited"
                  : row.score === null
                    ? "Not graded"
                    : "Saved"}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
  {#if !visible.length}
    <div class="empty">
      {initial.length
        ? "No students match your search."
        : "No students are enrolled in this section."}
    </div>
  {/if}
  <div class="save-bar">
    <p aria-live="polite">
      {changed.length} changed rows{hidden ? ` · ${hidden} hidden by search` : ""}
    </p>
    <div class="actions">
      <button
        type="button"
        disabled={!changed.length || $submitting}
        onclick={() => controller.reset({ data: { grades: [] } })}
      >
        Discard edits
      </button>
      {#if $message?.uncertain}
        <button formaction={actionUrl(page.url, "verify")} disabled={$submitting}>
          Verify current grades
        </button>
      {/if}
      <button
        class="primary"
        disabled={!changed.length || $submitting || closed || $message?.uncertain}
      >
        {$submitting ? "Saving…" : `Save ${changed.length} grades`}
      </button>
    </div>
  </div>
  <noscript><p>This batch editor needs JavaScript to track changed rows.</p></noscript>
</form>

<style>
  .score-input {
    width: 100px;
    max-width: 100%;
    font-variant-numeric: tabular-nums;
  }
  .save-bar {
    position: sticky;
    bottom: 0;
    background: white;
    border: 1px solid var(--line);
    padding: 16px 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
    z-index: 5;
  }
  .save-bar p,
  .save-bar .actions {
    margin: 0;
  }
  @media (max-width: 620px) {
    .grade-table td[data-label="Student"] {
      margin-bottom: 8px;
    }
    .grade-table .field-error {
      width: 100%;
      text-align: left;
    }
    .score-input {
      width: 110px;
    }
    .save-bar {
      position: static;
      display: block;
      padding: 16px;
    }
    .save-bar .actions {
      margin-top: 12px;
    }
  }
</style>
