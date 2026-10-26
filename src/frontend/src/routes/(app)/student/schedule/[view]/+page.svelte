<script lang="ts">
  import TimetableGrid from "#lib/components/TimetableGrid.svelte";
  import PageHeader from "#lib/components/PageHeader.svelte";
  import TermSelect from "#lib/components/TermSelect.svelte";
  import Agenda from "#lib/components/Agenda.svelte";
  import { timestamp } from "#lib/domain.ts";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
  let grid = $state(false);
</script>

<PageHeader title="Schedule" subtitle="Classes and exams · Campus time UTC+7">
  <TermSelect terms={data.terms} value={data.termId} />
</PageHeader>
<nav class="tabs" aria-label="Schedule views">
  <a class:chosen={data.view === "classes"} href="/student/schedule/classes?term={data.termId}">
    Classes
  </a>
  <a class:chosen={data.view === "exams"} href="/student/schedule/exams?term={data.termId}">
    Exams
  </a>
</nav>
{#if data.view === "classes"}
  <div class="section-heading">
    <h2>Weekly classes</h2>
    <button onclick={() => (grid = !grid)}>{grid ? "Show agenda" : "Show weekly grid"}</button>
  </div>
  <p class="muted small">
    From your weekly timetable. Holidays and make-up sessions are not included.
  </p>
  {#if grid}
    <TimetableGrid entries={data.entries} />
  {/if}
  <div class:agenda-under-grid={grid}><Agenda entries={data.entries} /></div>
{:else}
  <h2>Exam schedule</h2>
  {#each [...data.exams].sort((a, b) => a.starts_at.localeCompare(b.starts_at)) as exam}
    <article class="record" id="exam-{exam.id}">
      <div>
        <span class="course-code">{exam.course_code} · {exam.section_code}</span>
        <h3>{exam.course_name} · {exam.label}</h3>
        <p>{timestamp(exam.starts_at)} – {timestamp(exam.ends_at)}</p>
        <p class="muted">{exam.room || "Room not provided"}</p>
        {#each exam.conflicts_with as id}
          <p class="conflict">
            Overlaps <a href="#exam-{id}">
              {data.exams.find((e) => e.id === id)?.course_code ?? "another exam"} · {data.exams.find(
                (e) => e.id === id,
              )
                ? timestamp(data.exams.find((e) => e.id === id)!.starts_at)
                : ""}
            </a>
          </p>
        {/each}
      </div>
    </article>
  {/each}
  {#if !data.exams.length}
    <div class="empty">No exams are scheduled for your registered sections in this semester.</div>
  {/if}
{/if}

<style>
  .tabs {
    display: flex;
    border-bottom: 1px solid var(--line);
    margin: 0 0 28px;
    gap: 24px;
  }
  .tabs a {
    padding: 12px 4px;
    min-height: 44px;
    text-decoration: none;
    color: var(--muted);
    border-bottom: 3px solid transparent;
    font-size: 14px;
  }
  .tabs .chosen {
    color: var(--ink);
    border-bottom-color: var(--ink);
    font-weight: 600;
  }
  .agenda-under-grid {
    margin-top: 24px;
  }
</style>
