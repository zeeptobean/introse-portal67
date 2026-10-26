<script lang="ts">
  import { untrack } from "svelte";
  import { tick } from "svelte";
  import { page } from "$app/state";
  import { goto, invalidate } from "$app/navigation";
  import { portalForm, actionUrl } from "#lib/form.ts";
  import type { Enrollment } from "#lib/schemas/academics.ts";
  import { meetingLabel, registrationOpen } from "#lib/domain.ts";
  import PageHeader from "#lib/components/PageHeader.svelte";
  import TermSelect from "#lib/components/TermSelect.svelte";
  import FormFeedback from "#lib/components/FormFeedback.svelte";
  import ConfirmDialog from "#lib/components/ConfirmDialog.svelte";
  import SectionPreview from "#lib/components/SectionPreview.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
  let cancelling = $state<Enrollment | null>(null);
  let cancelElement: HTMLFormElement;
  const controller = portalForm(
    untrack(() => data.cancelForm),
    { invalidateAll: "pessimistic" },
  );
  const { form, message, allErrors, submitting, enhance } = controller;
  const q = $derived(page.url.searchParams.get("q") ?? "");
  const kind = $derived(page.url.searchParams.get("kind") ?? "");
  const sections = $derived(
    data.sections.filter(
      (s) =>
        (!kind || s.course.kind === kind) &&
        `${s.course.code} ${s.course.name}`.toLowerCase().includes(q.toLowerCase()),
    ),
  );
  const courses = $derived([...new Map(sections.map((s) => [s.course.id, s.course])).values()]);
  const active = $derived(data.enrollments.filter((e) => e.status === "enrolled"));
  const preview = $derived(
    data.sections.find((s) => s.id === page.url.searchParams.get("section")),
  );
  function filter(key: string, value: string) {
    const url = new URL(page.url.href);
    if (value) url.searchParams.set(key, value);
    else url.searchParams.delete(key);
    void goto(url, { reset: false, replace: true });
  }
  async function cancel() {
    if (!cancelling) return;
    form.set({ enrollment_id: cancelling.id });
    await tick();
    cancelElement.requestSubmit();
    cancelling = null;
  }
</script>

<PageHeader title="Registration" subtitle="Compare sections. Build a week that works.">
  <TermSelect terms={data.terms} value={data.termId} />
</PageHeader>
<FormFeedback message={$message} errors={$allErrors} />
<form method="POST" action={actionUrl(page.url, "cancel")} use:enhance bind:this={cancelElement}>
  <input type="hidden" name="enrollment_id" value={$form.enrollment_id} />
</form>
<div class="registration-layout">
  <section>
    <div class="toolbar">
      <label class="grow">
        Search courses
        <input
          type="search"
          value={q}
          oninput={(e) => filter("q", e.currentTarget.value)}
          placeholder="Course code or name"
        />
      </label>
      <label>
        Course type
        <select value={kind} onchange={(e) => filter("kind", e.currentTarget.value)}>
          <option value="">All types</option>
          <option value="course">Courses</option>
          <option value="specialized_topic">Specialized topics</option>
        </select>
      </label>
    </div>
    <div class="section-heading">
      <p class="muted small" aria-live="polite">{sections.length} sections available to compare</p>
      <button class="text-button" onclick={() => invalidate("portal:academics")}>
        Refresh availability
      </button>
    </div>
    {#each courses as course}
      <section class="course-group">
        <header>
          <span class="course-code">{course.code}</span>
          <h2>{course.name}</h2>
          <p class="muted small">
            {course.credits} credits · {course.kind === "specialized_topic"
              ? "Specialized topic"
              : "Course"}
          </p>
        </header>
        {#each sections.filter((s) => s.course.id === course.id) as section}
          <article class="section-row">
            <div>
              <strong>Section {section.section_code}</strong>
              <p class="muted small">{section.lecturer.full_name}</p>
            </div>
            <div class="small">
              {#each section.meetings as meeting}
                <p>{meetingLabel(meeting)}</p>
              {/each}
            </div>
            <div>
              <span class:positive={!!section.enrollment_id} class="small">
                {section.enrollment_id
                  ? "✓ Registered"
                  : section.seats_available === 0
                    ? "Full"
                    : `${section.seats_available} / ${section.capacity} available`}
              </span>
              <a
                class="button compact"
                href="?{new URLSearchParams({
                  ...Object.fromEntries(page.url.searchParams),
                  section: section.id,
                })}#section-preview"
              >
                View section
              </a>
            </div>
          </article>
        {/each}
      </section>
    {/each}
    {#if !sections.length}
      <div class="empty">
        {q || kind
          ? "No sections match your filters."
          : "No sections are available for this semester."}
        {#if q || kind}
          <button onclick={() => goto(`?term=${data.termId}`)}>Clear filters</button>
        {/if}
      </div>
    {/if}
  </section>
  <aside class="registration-aside">
    <section class="panel">
      <div class="section-heading">
        <h2>My registrations</h2>
        <span class="count">{active.length}</span>
      </div>
      <p class="muted small">
        {active.reduce((n, e) => n + e.course.credits, 0)} credits · Confirmed registrations
      </p>
      {#each active as enrollment}
        <article class="registered-row">
          <a
            href="?term={data.termId}&section={enrollment.section_id}#section-preview"
            class="course-code"
          >
            {enrollment.course.code} · {enrollment.section_code}
          </a>
          <p>{enrollment.course.name}</p>
          <p class="muted small">{enrollment.course.credits} credits</p>
          {#if enrollment.graded}
            <span class="small muted">Graded · cannot cancel</span>
          {:else if !registrationOpen(enrollment.term)}
            <span class="small muted">Registration is closed</span>
          {:else}
            <button
              class="text-button"
              disabled={$submitting || $message?.uncertain}
              onclick={() => (cancelling = enrollment)}
            >
              Cancel
            </button>
          {/if}
        </article>
      {/each}
      {#if !active.length}
        <p class="empty">No registrations for this semester yet.</p>
      {/if}
      <a class="arrow-link" href="/student/schedule/classes?term={data.termId}">
        View timetable <span>→</span>
      </a>
    </section>
    {#if preview}
      {#key preview.id}
        <SectionPreview section={preview} initial={data.registerForm} timetable={data.timetable} />
      {/key}
    {/if}
    <details>
      <summary>Cancelled registrations</summary>
      {#each data.enrollments.filter((e) => e.status === "cancelled") as enrollment}
        <p class="small">
          {enrollment.course.code} · {enrollment.section_code} · Cancelled
        </p>
      {/each}
    </details>
  </aside>
</div>
{#if cancelling}
  <ConfirmDialog
    title="Cancel {cancelling.course.code} · section {cancelling.section_code}?"
    action="Cancel registration"
    oncancel={() => (cancelling = null)}
    onconfirm={cancel}
  >
    <p>{cancelling.course.name} · {cancelling.term.name}</p>
    <p>Your seat will be released. Registering again depends on availability.</p>
    {#each data.sections.find((s) => s.id === cancelling?.section_id)?.meetings ?? [] as meeting}
      <p>
        {meetingLabel(meeting)}
      </p>
    {/each}
  </ConfirmDialog>
{/if}

<style>
  .registration-layout {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 304px;
    gap: 28px;
    align-items: start;
  }
  .registration-aside {
    min-width: 0;
  }
  @media (min-width: 1500px) {
    .registration-layout {
      grid-template-columns: minmax(0, 1fr) 320px;
      gap: 32px;
    }
  }
  @media (max-width: 1200px) {
    .registration-layout {
      grid-template-columns: minmax(0, 1fr);
    }
    .registration-aside {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 24px;
      align-items: start;
    }
    .registration-aside > details {
      grid-column: 1/-1;
    }
  }
  @media (max-width: 620px) {
    .registration-aside {
      display: block;
    }
  }
  .course-group {
    margin-bottom: 24px;
    background: white;
    border: 1px solid var(--line);
    border-radius: 8px;
    overflow: hidden;
  }
  .course-group > header {
    padding: 20px 20px 16px;
    border-bottom: 1px solid var(--line);
  }
  .course-group h2 {
    font-size: 17px;
    margin: 4px 0;
  }
  .course-group header p {
    margin: 6px 0 0;
  }
  .section-row {
    padding: 18px 20px;
    display: grid;
    grid-template-columns: minmax(100px, 0.8fr) minmax(160px, 1.4fr) minmax(100px, 0.8fr);
    gap: 16px;
    border-bottom: 1px solid var(--line);
    font-size: 14px;
  }
  .section-row:last-child {
    border-bottom: 0;
  }
  .section-row p {
    margin: 4px 0;
  }
  .section-row .button {
    display: flex;
    margin-top: 12px;
  }
  @media (max-width: 620px) {
    .course-group > header,
    .section-row {
      padding: 18px 16px;
    }
    .section-row {
      grid-template-columns: 1fr;
      gap: 8px;
    }
    .section-row > div:last-child {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 16px;
    }
    .section-row .button {
      margin: 0;
    }
  }
  .registered-row {
    padding: 16px 0;
    border-bottom: 1px solid var(--line);
    font-size: 14px;
  }
  .registered-row p {
    margin: 4px 0;
  }
  .registered-row button {
    font-size: 12px;
  }
  .count {
    padding: 3px 9px;
    background: var(--tint);
    color: var(--ink);
    border-radius: 4px;
    font-size: 13px;
  }
</style>
