<script lang="ts">
  import type { TimetableEntry } from "#lib/schemas/academics.ts";
  import { days } from "#lib/domain.ts";
  let {
    entries,
    empty = "No classes in this schedule.",
  }: { entries: TimetableEntry[]; empty?: string } = $props();
</script>

{#if !entries.length}
  <div class="empty">{empty}</div>
{:else}
  <div class="agenda">
    {#each [...entries].sort((a, b) => a.weekday - b.weekday || a.starts_at.localeCompare(b.starts_at)) as entry}
      <article class="agenda-row">
        <div class="agenda-time">
          <strong>{days[entry.weekday]}</strong>
          <span>{entry.starts_at}–{entry.ends_at}</span>
        </div>
        <div>
          <a
            class="course-code"
            href="/student/registration/sections/{entry.section_id}?term={entry.term.id}"
          >
            {entry.course_code} · {entry.section_code}
          </a>
          <h3>{entry.course_name}</h3>
          <p class="muted">{entry.room || "Room not provided"} · {entry.lecturer_name}</p>
          {#if entry.conflicts_with.length}
            <p class="conflict">
              Overlaps {entries
                .filter((e) => entry.conflicts_with.includes(e.id))
                .map((e) => `${e.course_code} ${e.starts_at}–${e.ends_at}`)
                .join(", ")}
            </p>
          {/if}
        </div>
      </article>
    {/each}
  </div>
{/if}

<style>
  .agenda {
    background: white;
    border-top: 1px solid var(--line);
  }
  .agenda-row {
    display: grid;
    grid-template-columns: 148px minmax(0, 1fr);
    padding: 24px;
    border-bottom: 1px solid var(--line);
    gap: 24px;
  }
  .agenda-row h3 {
    margin: 4px 0;
  }
  .agenda-row p {
    margin: 4px 0 0;
  }
  .agenda-time {
    display: flex;
    flex-direction: column;
    font-variant-numeric: tabular-nums;
    font-size: 14px;
    gap: 4px;
  }
  .agenda-time span {
    color: var(--muted);
  }
  @media (max-width: 620px) {
    .agenda-row {
      grid-template-columns: 1fr;
      padding: 20px 16px;
      gap: 12px;
    }
    .agenda-time {
      flex-direction: row;
      flex-wrap: wrap;
      gap: 12px;
    }
  }
</style>
