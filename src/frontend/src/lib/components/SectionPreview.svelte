<script lang="ts">
  import { untrack } from "svelte";
  import { page } from "$app/state";
  import type { SuperValidated } from "sveltekit-superforms";
  import { portalForm, actionUrl } from "#lib/form.ts";
  import type { Section, TimetableEntry } from "#lib/schemas/academics.ts";
  import { conflicts, meetingLabel, registrationOpen, timestamp } from "#lib/domain.ts";
  import Notice from "#lib/components/Notice.svelte";
  import Status from "#lib/components/Status.svelte";
  import Agenda from "#lib/components/Agenda.svelte";
  import FormFeedback from "#lib/components/FormFeedback.svelte";
  let {
    section,
    timetable,
    initial,
  }: {
    section: Section;
    timetable: TimetableEntry[];
    initial: SuperValidated<{ section_id: string }>;
  } = $props();
  const { message, allErrors, submitting, enhance } = portalForm(
    untrack(() => ({ ...initial, data: { section_id: section.id } })),
    { invalidateAll: "pessimistic" },
  );
  const clashes = $derived(
    conflicts(
      section.meetings,
      timetable.filter((t) => t.section_id !== section.id),
    ),
  );
</script>

<div class="panel preview-panel" id="section-preview">
  <p class="eyebrow">Section detail</p>
  <h2>
    <span class="course-code">{section.course.code}</span>
    · {section.section_code}
  </h2>
  <p>{section.course.name}</p>
  <p class="muted">{section.course.credits} credits · {section.lecturer.full_name}</p>
  <ul class="meeting-list">
    {#each section.meetings as meeting}
      <li>{meetingLabel(meeting)}</li>
    {/each}
  </ul>
  <p>{section.seats_available} of {section.capacity} seats available</p>
  <p class="small muted">
    Registration: {timestamp(section.term.registration_opens_at)} – {timestamp(
      section.term.registration_closes_at,
    )}
  </p>
  <FormFeedback message={$message} errors={$allErrors} />
  {#if section.enrollment_id}
    <p>
      <Status value="enrolled" />
    </p>
  {:else}
    {#if clashes.length}
      <Notice tone="caution">
        <strong>This section overlaps your timetable</strong>
        {#each clashes as clash}
          <p>
            {clash.course_code} · {clash.section_code}: {meetingLabel(clash)}
          </p>
        {/each}
      </Notice>
    {/if}
    <form method="POST" action={actionUrl(page.url, "register")} use:enhance>
      <input type="hidden" name="section_id" value={section.id} />
      <button
        class="primary"
        disabled={$submitting ||
          $message?.uncertain ||
          !registrationOpen(section.term) ||
          section.seats_available === 0 ||
          clashes.length > 0}
      >
        {$submitting
          ? "Registering…"
          : section.seats_available === 0
            ? "Full"
            : !registrationOpen(section.term)
              ? "Registration is closed"
              : `Register for section ${section.section_code}`}
      </button>
    </form>
  {/if}
  <details class="comparison" open>
    <summary>Schedule comparison · Preview only</summary>
    <p class="small muted">This preview does not reserve a seat or change your registrations.</p>
    <div class="candidate">
      <strong>Preview only · {section.course.code}</strong>
      {#each section.meetings as meeting}
        <p>{meetingLabel(meeting)}</p>
      {/each}
    </div>
    <h3>Currently registered</h3>
    <Agenda entries={timetable} />
  </details>
</div>

<style>
  .preview-panel h2 {
    font-size: 18px;
  }
  /* The comparison agenda is a compact list inside this narrow panel. */
  .preview-panel :global(h3) {
    font-size: 14px;
  }
  .preview-panel :global(.agenda-row) {
    display: block;
    padding: 16px 0;
  }
  .preview-panel :global(.agenda-time) {
    margin-bottom: 8px;
  }
  .meeting-list {
    padding-left: 20px;
    font-size: 14px;
  }
  .meeting-list li {
    margin: 8px 0;
  }
  .comparison {
    margin-top: 24px;
  }
  .candidate {
    border: 1px dashed var(--ink);
    padding: 16px;
    margin: 16px 0;
    background: var(--tint);
    font-size: 13px;
  }
  .candidate p {
    margin: 8px 0;
  }
</style>
