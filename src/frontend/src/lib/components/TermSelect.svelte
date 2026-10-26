<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import type { Term } from "#lib/schemas/shared.ts";
  let { terms, value }: { terms: Term[]; value: string } = $props();
  function change(event: Event) {
    const url = new URL(page.url.href);
    url.searchParams.set("term", (event.target as HTMLSelectElement).value);
    void goto(url, { reset: false });
  }
</script>

<label class="term-select">
  Semester
  <select aria-label="Semester" {value} onchange={change}>
    {#each terms as term}
      <option value={term.id}>{term.name} · {term.code}</option>
    {/each}
  </select>
</label>

<style>
  .term-select {
    font-size: 12px;
    color: var(--muted);
  }
  .term-select select {
    margin-top: 8px;
    max-width: 340px;
    font-size: 13px;
    color: #18232d;
  }
  @media (max-width: 900px) {
    .term-select select {
      max-width: 280px;
    }
  }
  @media (max-width: 620px) {
    .term-select {
      width: 100%;
    }
    .term-select select {
      width: 100%;
      max-width: none;
    }
  }
</style>
