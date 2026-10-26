<script lang="ts">
  import type { TimetableEntry } from "#lib/schemas/academics.ts";
  import { days, meetingLanes } from "#lib/domain.ts";
  let { entries }: { entries: TimetableEntry[] } = $props();
  const minutes = (time: string) => Number(time.slice(0, 2)) * 60 + Number(time.slice(3, 5));
  const start = $derived(
    Math.floor(Math.min(420, ...entries.map((e) => minutes(e.starts_at))) / 60) * 60,
  );
  const end = $derived(
    Math.ceil(Math.max(1080, ...entries.map((e) => minutes(e.ends_at))) / 60) * 60,
  );
  const hours = $derived(Array.from({ length: (end - start) / 60 }, (_, i) => start + i * 60));
  const positions = $derived(meetingLanes(entries));
</script>

<!-- svelte-ignore a11y_no_noninteractive_tabindex (keyboard users need to scroll this region) -->
<div
  class="timetable-scroll"
  role="region"
  aria-label="Weekly timetable; agenda follows below"
  tabindex="0"
>
  <div class="timetable-grid" style:--height="{end - start}px">
    <div class="time-column">
      <h3>Time</h3>
      <div class="time-axis">
        {#each hours as hour}
          <span style:top="{hour - start}px">
            {String(hour / 60).padStart(2, "0")}:00
          </span>
        {/each}
      </div>
    </div>
    {#each days.slice(1) as day, index}
      <section class="grid-day">
        <h3>{day}</h3>
        <div class="day-slots">
          {#each hours as hour}
            <div class="hour-line" style:top="{hour - start}px"></div>
          {/each}
          {#each entries.filter((e) => e.weekday === index + 1) as entry}
            {@const position = positions.get(entry.id)!}
            <a
              class="timetable-block"
              style:top="{minutes(entry.starts_at) - start}px"
              style:height="{minutes(entry.ends_at) - minutes(entry.starts_at)}px"
              style:left="calc({(position.index / position.count) * 100}% + 2px)"
              style:width="calc({100 / position.count}% - 4px)"
              href="/student/registration/sections/{entry.section_id}?term={entry.term.id}"
              aria-label="{entry.course_name}, section {entry.section_code}, {day} {entry.starts_at} to {entry.ends_at}, {entry.room ||
                'Room not provided'}"
            >
              <strong>{entry.starts_at}–{entry.ends_at}</strong>
              <span>{entry.course_code} · {entry.section_code}</span>
              <span>{entry.room || "Room not provided"}</span>
              {#if entry.conflicts_with.length}
                <span class="conflict">Overlap</span>
              {/if}
            </a>
          {/each}
        </div>
      </section>
    {/each}
  </div>
</div>

<style>
  .timetable-scroll {
    overflow: auto;
    border: 1px solid var(--line);
    background: white;
    margin: 24px 0;
  }
  .timetable-grid {
    display: grid;
    grid-template-columns: 64px repeat(7, minmax(140px, 1fr));
    min-width: 1050px;
  }
  .timetable-grid h3 {
    font-size: 12px;
    padding: 12px 8px;
    margin: 0;
    height: 44px;
    border-bottom: 1px solid var(--line);
  }
  .time-axis,
  .day-slots {
    position: relative;
    height: var(--height);
  }
  .time-axis > span {
    position: absolute;
    font-size: 11px;
    color: var(--muted);
    padding: 4px 6px;
  }
  .grid-day {
    border-left: 1px solid var(--line);
  }
  .hour-line {
    position: absolute;
    left: 0;
    right: 0;
    border-top: 1px solid #edf0f2;
    pointer-events: none;
  }
  .timetable-block {
    position: absolute;
    border-left: 3px solid var(--ink);
    padding: 8px 6px;
    display: flex;
    flex-direction: column;
    gap: 6px;
    font-size: 11px;
    background: var(--tint);
    text-decoration: none;
    overflow: auto;
    line-height: 1.3;
  }
  .timetable-block:hover {
    background: #dce9f2;
  }
</style>
