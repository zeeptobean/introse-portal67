<script lang="ts">
  import type { FormMessage } from "#lib/form.ts";
  import { fieldId } from "#lib/form.ts";
  import Notice from "#lib/components/Notice.svelte";
  let {
    message,
    errors = [],
    onSelectError,
  }: {
    message?: FormMessage;
    errors?: { path: string; messages: string[] }[];
    onSelectError?: (path: string) => void | Promise<void>;
  } = $props();
</script>

{#if message}
  <Notice tone={message.tone ?? "info"}>
    {message.text}
    {#if message.uncertain}
      <p>Do not repeat the submission yet. Refresh the relevant records to check its outcome.</p>
    {/if}
  </Notice>
{/if}
{#if errors.length}
  <div class="notice error" role="alert" tabindex="-1">
    <strong>Check the following fields</strong>
    <ul>
      {#each errors as error}
        <li>
          <a
            href="#{fieldId(error.path)}"
            onclick={onSelectError
              ? (event) => {
                  event.preventDefault();
                  void onSelectError(error.path);
                }
              : undefined}
          >
            {error.messages.join(" ")}
          </a>
        </li>
      {/each}
    </ul>
  </div>
{/if}
