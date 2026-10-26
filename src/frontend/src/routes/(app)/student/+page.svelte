<script lang="ts">
  import PageHeader from "#lib/components/PageHeader.svelte";
  import TermSelect from "#lib/components/TermSelect.svelte";
  import Agenda from "#lib/components/Agenda.svelte";
  import Notice from "#lib/components/Notice.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
  const now = new Date();
  const campus = new Intl.DateTimeFormat("en-CA", {
    timeZone: "Asia/Ho_Chi_Minh",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(now);
  const weekdayName = new Intl.DateTimeFormat("en-US", {
    timeZone: "Asia/Ho_Chi_Minh",
    weekday: "short",
  }).format(now);
  const weekday = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"].indexOf(weekdayName) + 1;
  const today = $derived(
    data.entries.filter(
      (e) => e.weekday === weekday && e.term.starts_on <= campus && e.term.ends_on >= campus,
    ),
  );
  const active = $derived(data.enrollments.filter((e) => e.status === "enrolled"));
</script>

<PageHeader
  title="Today"
  subtitle={new Intl.DateTimeFormat("en-GB", {
    dateStyle: "full",
    timeZone: "Asia/Ho_Chi_Minh",
  }).format(now)}
>
  <TermSelect terms={data.terms} value={data.termId} />
</PageHeader>
<div class="today-layout">
  <section>
    <div class="section-heading">
      <h2>Your day, at a glance</h2>
      <span class="muted small">From your weekly timetable</span>
    </div>
    <Agenda entries={today} empty="No classes today in the selected semester." />
    <a class="arrow-link" href="/student/schedule/classes?term={data.termId}">
      View full weekly schedule <span>→</span>
    </a>
    <div class="registration-summary">
      <span class="eyebrow">Your semester</span>
      <h2>
        {active.length} registered sections · {active.reduce((sum, e) => sum + e.course.credits, 0)} credits
      </h2>
      <a href="/student/registration?term={data.termId}">Review your registrations →</a>
    </div>
  </section>
  <aside class="panel">
    <p class="eyebrow">Plan ahead</p>
    <h2>Exams & exceptions</h2>
    {#if data.exams.some((e) => e.conflicts_with.length)}
      <Notice tone="caution">
        Your exam schedule contains an overlap. <a
          href="/student/schedule/exams?term={data.termId}"
        >
          Review conflicting exams
        </a>
      </Notice>
    {/if}
    {#if data.entries.some((e) => e.conflicts_with.length)}
      <Notice tone="caution">
        <a href="/student/schedule/classes?term={data.termId}">Review class overlaps</a>
      </Notice>
    {/if}
    <p class="muted">Check dates, rooms, and any clashes before exam week.</p>
    <a class="arrow-link" href="/student/schedule/exams?term={data.termId}">
      View exam schedule <span>→</span>
    </a>
  </aside>
</div>

<style>
  .today-layout {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 320px;
    gap: 32px;
    align-items: start;
  }
  .today-layout > aside {
    margin-top: 0;
  }
  @media (max-width: 1200px) {
    .today-layout {
      grid-template-columns: minmax(0, 1fr) 280px;
      gap: 24px;
    }
  }
  @media (max-width: 900px) {
    .today-layout {
      grid-template-columns: 1fr;
    }
  }
  .registration-summary {
    border-top: 1px solid var(--line);
    margin-top: 40px;
    padding-top: 28px;
  }
</style>
