<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import type { DocumentRequest } from "#lib/schemas/documents.ts";
  import { requestTypes } from "#lib/schemas/documents.ts";
  import { timestamp } from "#lib/domain.ts";
  import Status from "#lib/components/Status.svelte";
  let { records, staff = false }: { records: DocumentRequest[]; staff?: boolean } = $props();
  const status = $derived(page.url.searchParams.get("status") ?? (staff ? "pending" : ""));
  const filtered = $derived(
    [...records]
      .filter((r) => !status || r.status === status)
      .sort((a, b) =>
        staff && status === "pending"
          ? a.submitted_at.localeCompare(b.submitted_at)
          : b.submitted_at.localeCompare(a.submitted_at),
      ),
  );
  function change(value: string) {
    const url = new URL(page.url.href);
    url.searchParams.set("status", value);
    void goto(url, { reset: false });
  }
</script>

<div class="toolbar">
  <label>
    Status
    <select value={status} onchange={(e) => change(e.currentTarget.value)}>
      <option value="">All requests</option>
      <option value="pending">Pending</option>
      <option value="approved">Approved</option>
      <option value="rejected">Rejected</option>
    </select>
  </label>
  <p class="muted small" aria-live="polite">
    {filtered.length} requests{staff && status === "pending" ? " · Oldest first" : ""}
  </p>
</div>
<div class="table-wrap">
  <table>
    <thead>
      <tr>
        <th>Request</th>
        {#if staff}
          <th>Student</th>
        {/if}
        <th>Submitted</th>
        <th>Status</th>
        <th><span class="sr-only">Action</span></th>
      </tr>
    </thead>
    <tbody>
      {#each filtered as record}
        <tr>
          <td data-label="Request">
            <a href="/{staff ? 'staff' : 'student'}/documents/{record.id}">
              {requestTypes[record.request_type]}
            </a>
            <span class="muted small">Reference {record.id}</span>
          </td>
          {#if staff}
            <td data-label="Student">
              <strong>{record.student.full_name}</strong>
              <span class="muted small">
                {record.student.institutional_id ?? "ID not provided"}
              </span>
            </td>
          {/if}
          <td data-label="Submitted">{timestamp(record.submitted_at)}</td>
          <td data-label="Status"><Status value={record.status} /></td>
          <td>
            <a href="/{staff ? 'staff' : 'student'}/documents/{record.id}">
              {staff && record.status === "pending" ? "Review" : "View"} →
            </a>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
{#if !filtered.length}
  <div class="empty">
    {records.length
      ? "No requests match this status."
      : staff
        ? "No document requests have been submitted."
        : "You haven’t requested a document yet."}
    {#if records.length}
      <button onclick={() => change("")}>Show all requests</button>
    {/if}
  </div>
{/if}
