<script lang="ts">
  import type { Snippet } from "svelte";
  let {
    title,
    action,
    onconfirm,
    oncancel,
    children,
  }: {
    title: string;
    action: string;
    onconfirm: () => void;
    oncancel: () => void;
    children: Snippet;
  } = $props();
  let dialog = $state<HTMLDialogElement>();
  $effect(() => {
    dialog?.showModal();
    return () => dialog?.close();
  });
</script>

<dialog
  bind:this={dialog}
  oncancel={(event) => {
    event.preventDefault();
    oncancel();
  }}
  aria-labelledby="confirm-title"
>
  <h2 id="confirm-title">{title}</h2>
  <div class="dialog-body">{@render children()}</div>
  <div class="actions">
    <button onclick={oncancel}>Keep registration</button>
    <button class="danger" onclick={onconfirm}>{action}</button>
  </div>
</dialog>

<style>
  .dialog-body {
    margin: 20px 0;
  }
</style>
